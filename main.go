package main

import (
	"bxs/block_getter"
	"bxs/cache"
	"bxs/chain_params"
	"bxs/config"
	"bxs/log"
	"bxs/metrics"
	"bxs/parser"
	"bxs/repository"
	"bxs/sequencer"
	"bxs/service"
	"bxs/types"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func createDBService() service.DBService {
	var (
		txDb             *gorm.DB
		txDbErr          error
		tokenPairDb      *gorm.DB
		tokenPairDbErr   error
		tokenRepository  *repository.TokenRepository
		pairRepository   *repository.PairRepository
		txRepository     *repository.TxRepository
		actionRepository *repository.ActionRepository
	)

	if config.G.TxDatabase.Enabled {
		txDb, txDbErr = gorm.Open(postgres.Open(config.G.TxDatabase.DBDatasource.GetPostgresDsn()))
		if txDbErr != nil {
			log.Logger.Fatal("failed to connect to tx db", zap.Error(txDbErr))
		}

		txRepository = repository.NewTxRepository(txDb)
		actionRepository = repository.NewActionRepository(txDb)
	}

	if config.G.TokenPairDatabase.Enabled {
		tokenPairDb, tokenPairDbErr = gorm.Open(postgres.Open(config.G.TokenPairDatabase.DBDatasource.GetPostgresDsn()))
		if tokenPairDbErr != nil {
			log.Logger.Fatal("failed to connect to token_pair db", zap.Error(tokenPairDbErr))
		}

		tokenRepository = repository.NewTokenRepository(tokenPairDb)
		pairRepository = repository.NewPairRepository(tokenPairDb)
	}

	return service.NewDBService(tokenRepository, pairRepository, txRepository, actionRepository)
}

func main() {
	time.Local = time.UTC

	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "show version information")
	var configFile string
	flag.StringVar(&configFile, "c", "config.json", "config file")
	flag.Parse()

	if showVersion {
		fmt.Println(GetVersion())
		os.Exit(0)
	}

	log.Logger.Info(GetVersion().String())
	log.Logger.Info("config", zap.String("file path", configFile))
	loadConfigErr := config.LoadConfigFile(configFile)
	if loadConfigErr != nil {
		log.Logger.Fatal("load config file err", zap.Error(loadConfigErr))
	}

	chain_params.LoadNetwork(config.G.TestNet, config.G.XLaunchFactoryAddress)
	log.InitLogger()
	metrics.Init(config.G.MetricsPort)

	ethClient, dialEthErr := ethclient.Dial(config.G.Chain.Endpoint)
	if dialEthErr != nil {
		log.Logger.Fatal("Failed to connect to the chain(http): %v", zap.Error(dialEthErr))
	}

	wsEthClient, err := ethclient.Dial(config.G.Chain.WsEndpoint)
	if err != nil {
		log.Logger.Fatal("Failed to connect to the chain(ws): %v", zap.Error(err))
	}

	ethClientArchive, dialEthErrArchive := ethclient.Dial(config.G.Chain.EndpointArchive)
	if dialEthErrArchive != nil {
		log.Logger.Fatal("Failed to connect to the bsc archive(http): %v", zap.Error(dialEthErrArchive))
	}

	contractCallerArchive := service.NewContractCaller(ethClientArchive, config.G.ContractCaller.Retry.GetRetryParams())

	redisCli := redis.NewClient(&redis.Options{
		Addr:     config.G.Redis.Addr,
		Username: config.G.Redis.Username,
		Password: config.G.Redis.Password,
	})
	cache := cache.NewTwoTierCache(redisCli)

	contractCaller := service.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())

	pairService := service.NewPairService(cache, contractCaller)
	priceService := service.NewPriceService(config.G.PriceService.FromChain, cache, contractCallerArchive, ethClientArchive, config.G.PriceService.PoolSize)

	sequencerForBlockHandler := sequencer.NewSequencer()

	topicRouter := parser.NewTopicRouter()
	kafkaSender := service.NewKafkaSender(config.G.Kafka)

	blockParser := parser.NewBlockParser(
		cache,
		sequencerForBlockHandler,
		priceService,
		pairService,
		topicRouter,
		kafkaSender,
		createDBService(),
	)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	blockParser.Start(wg)

	sequencerForBlockGetter := sequencer.NewSequencer()
	blockGetter := block_getter.NewBlockGetter(config.G.BlockGetter.SubHeader, wsEthClient, cache, sequencerForBlockGetter, config.G.BlockGetter.Retry.GetRetryParams())
	startBlockNumber := blockGetter.GetStartBlockNumber(config.G.BlockGetter.StartBlockNumber)
	if startBlockNumber == 0 {
		log.Logger.Fatal("start block number is zero")
	}

	sequencerForBlockGetter.Init(startBlockNumber)
	sequencerForBlockHandler.Init(startBlockNumber)

	priceService.Start(startBlockNumber)
	blockGetter.Start()
	blockGetter.StartDispatch(startBlockNumber)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Logger.Info("receive signal", zap.String("signal", sig.String()))
		blockGetter.Stop()
	}()

	var blockCtx *types.ParseBlockContext
	for {
		blockCtx = blockGetter.Next()
		if blockCtx == nil {
			log.Logger.Info("no more block to parse")
			blockParser.Stop()
			break
		}
		blockParser.ParseBlockAsync(blockCtx)
	}

	log.Logger.Info("wait all block commited")
	wg.Wait()
	log.Logger.Info("all block commited")
	log.Logger.Sync()
}
