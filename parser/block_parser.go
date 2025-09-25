package parser

import (
	"bxs/cache"
	"bxs/config"
	"bxs/log"
	"bxs/metrics"
	"bxs/repository/orm"
	"bxs/sequencer"
	"bxs/service"
	"bxs/types"
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"sync"
	"time"
)

type BlockParser interface {
	Start(*sync.WaitGroup)
	Stop()
	ParseBlockAsync(bw *types.ParseBlockContext)
}

type blockParser struct {
	inputQueue   chan *types.ParseBlockContext
	workPool     *ants.Pool
	cache        cache.Cache
	sequencer    sequencer.Sequencer
	outputQueue  chan *types.ParseBlockContext
	priceService service.PriceService
	pairService  service.PairService
	topicRouter  TopicRouter
	kafkaSender  service.KafkaSender
	dbService    service.DBService
	parseTxPool  *ants.Pool
}

func NewBlockParser(
	cache cache.Cache,
	sequencer sequencer.Sequencer,
	priceService service.PriceService,
	pairService service.PairService,
	topicRouter TopicRouter,
	kafkaSender service.KafkaSender,
	dbService service.DBService,
) BlockParser {
	workPool, err := ants.NewPool(config.G.BlockHandler.PoolSize)
	if err != nil {
		log.Logger.Fatal("ants.NewPool err", zap.Error(err))
	}

	parseTxPool, err := ants.NewPool(config.G.BlockHandler.ParseTxPoolSize)
	if err != nil {
		log.Logger.Fatal("ants.NewPool err", zap.Error(err))
	}

	return &blockParser{
		inputQueue:   make(chan *types.ParseBlockContext, config.G.BlockHandler.QueueSize),
		workPool:     workPool,
		cache:        cache,
		sequencer:    sequencer,
		outputQueue:  make(chan *types.ParseBlockContext, config.G.BlockHandler.QueueSize),
		priceService: priceService,
		pairService:  pairService,
		topicRouter:  topicRouter,
		kafkaSender:  kafkaSender,
		dbService:    dbService,
		parseTxPool:  parseTxPool,
	}
}

func (p *blockParser) Commit(x sequencer.Sequenceable) {
	p.outputQueue <- x.(*types.ParseBlockContext)
}

func (p *blockParser) Start(waitGroup *sync.WaitGroup) {
	p.startHandleBlockResult(waitGroup)

	go func() {
		wg := &sync.WaitGroup{}
	tagFor:
		for {
			select {
			case pbc, ok := <-p.inputQueue:
				if !ok {
					log.Logger.Info("block handler inputQueue is closed")
					break tagFor
				}

				wg.Add(1)
				p.workPool.Submit(func() {
					defer wg.Done()
					p.parseBlock(pbc)
				})
			}
		}

		wg.Wait()
		log.Logger.Info("all block parse task finish")
		p.doStop()
	}()
}

func (p *blockParser) Stop() {
	close(p.inputQueue)
}

func (p *blockParser) ParseBlockAsync(bw *types.ParseBlockContext) {
	p.inputQueue <- bw
}

func (p *blockParser) waitForNativeTokenPrice(blockNumber *big.Int, blockTimestamp uint64) decimal.Decimal {
	for {
		price, err := p.priceService.GetPrice(blockNumber)
		if err != nil {
			log.Logger.Error("get price err", zap.Error(err), zap.Any("blockNumber", blockNumber), zap.Any("blockTimestamp", blockTimestamp))
			time.Sleep(time.Second)
			continue
		}
		return price
	}
}

func (p *blockParser) parseTxReceipt(pbc *types.ParseBlockContext, txReceipt *ethtypes.Receipt) *types.TxResult {
	txSender, err := pbc.GetTxSender(txReceipt.TransactionIndex)
	if err != nil {
		log.Logger.Fatal("Err: get tx sender err",
			zap.Error(err),
			zap.Any("pbc", pbc),
			zap.Any("txReceipt", txReceipt),
		)
	}

	txResult := types.NewTxResult(txSender, pbc.HeightTime.Time)
	pairs := make([]*types.Pair, 0, 2)
	tokens := make([]*types.Token, 0, 2)
	migratedPools := make([]*types.MigratedPool, 0, 2)
	for _, ethLog := range txReceipt.Logs {
		if len(ethLog.Topics) == 0 {
			continue
		}

		event, parseErr := p.topicRouter.Parse(ethLog)
		if parseErr != nil {
			continue
		}

		if event.IsCreated() {
			event.SetBlockTime(pbc.HeightTime.Time)

			pair := event.GetPair()
			token0 := event.GetToken0()
			p.cache.SetPair(pair)
			p.cache.SetToken(token0)
			pairs = append(pairs, pair)
			tokens = append(tokens, token0)
			txResult.AddPoolUpdate(event.GetPoolUpdate())
			continue
		}

		if event.IsPairCreated() {
			txResult.AddPairCreatedEvent(event)
			pairs = append(pairs, event.GetPair())
			continue
		}

		// swap event: buy or sell
		pairAddr := event.GetPairAddress()
		pair, ok := p.cache.GetPair(pairAddr)
		if !ok {
			log.Logger.Sugar().Warnf("pool %s not cached, it should be", pairAddr)
			pw := p.pairService.GetPair(pairAddr)
			pair = pw.Pair
			if pair == nil {
				log.Logger.Sugar().Fatalf("get pair %s fail", pairAddr)
			}

			if pair.Filtered {
				log.Logger.Sugar().Warnf("pair %s is filtered, filter code %d", pairAddr, pair.FilterCode)
				continue
			}

			pairs = append(pairs, pair)
			tokens = append(tokens, &types.Token{
				Address:  pair.Token0.Address,
				Symbol:   pair.Token0.Symbol,
				Decimals: pair.Token0.Decimals,
				Program:  types.ProtocolNameXLaunch,
			})
		}
		event.SetPair(pair)

		if event.IsMigrated() {
			p.cache.SetMigrateToken(event.GetPair().Token0Core.Address)
			migratedPools = append(migratedPools, &types.MigratedPool{
				Pool:  event.GetPair().Address,
				Token: event.GetPair().Token0Core.Address,
			})
		}

		txResult.AddPoolUpdate(event.GetPoolUpdate())
		txResult.AddSwapEvent(event)
	}

	txResult.SetPairs(pairs)
	txResult.SetTokens(tokens)
	txResult.SetMigratedPools(migratedPools)

	p.processTxPairCreatedEvents(txResult)
	return txResult
}

func (p *blockParser) processTxPairCreatedEvents(txResult *types.TxResult) *types.TxResult {
	actions := make([]*orm.Action, 0, 2)
	for _, event := range txResult.PairCreatedEvents {
		nonWBNBToken := event.GetNonWBNBToken()
		if p.cache.MigrateTokenExist(nonWBNBToken) {
			actions = append(actions, event.GetAction())
			p.cache.DelMigrateToken(nonWBNBToken)
		}
	}
	txResult.SetActions(actions)
	return txResult
}

func (p *blockParser) parseBlock(pbc *types.ParseBlockContext) {
	pbc.NativeTokenPrice = p.waitForNativeTokenPrice(pbc.HeightTime.HeightBigInt, pbc.HeightTime.Timestamp)

	now := time.Now()
	br := types.NewBlockResult(pbc.HeightTime.Height, pbc.HeightTime.Timestamp, pbc.NativeTokenPrice)

	for _, txReceipt := range pbc.BlockReceipts {
		if txReceipt.Status != 1 {
			continue
		}
		br.AddTxResult(p.parseTxReceipt(pbc, txReceipt))
	}

	duration := time.Since(now)
	metrics.ParseBlockDurationMs.Observe(float64(duration.Milliseconds()))
	log.Logger.Info(fmt.Sprintf("parse block %d duration %dms", pbc.HeightTime.HeightBigInt, duration.Milliseconds()))

	pbc.BlockResult = br
	p.sequencer.CommitWithSequence(pbc, p)
}

func (p *blockParser) commitBlockResult(blockResult *types.BlockResult) {
	blockInfo := blockResult.GetKafkaMessage()

	now := time.Now()
	var err error
	if len(blockInfo.NewTokens) > 0 {
		err = p.dbService.AddTokens(blockInfo.NewTokens)
		if err != nil {
			log.Logger.Fatal("add tokens err", zap.Any("height", blockInfo.Height), zap.Error(err))
		}
	}

	if len(blockInfo.NewPairs) > 0 {
		err = p.dbService.AddPairs(blockInfo.NewPairs)
		if err != nil {
			log.Logger.Fatal("add pairs err", zap.Any("height", blockInfo.Height), zap.Error(err))
		}
	}

	if len(blockInfo.Txs) > 0 {
		err = p.dbService.AddTxs(blockInfo.Txs)
		if err != nil {
			log.Logger.Fatal("add txs err", zap.Any("height", blockInfo.Height), zap.Error(err))
		}
	}

	if len(blockInfo.Actions) > 0 {
		err = p.dbService.AddActions(blockInfo.Actions)
		if err != nil {
			log.Logger.Fatal("add actions err", zap.Any("height", blockInfo.Height), zap.Error(err))
		}

		for _, action := range blockInfo.Actions {
			log.Logger.Sugar().Infof("add action: pair:%s, token:%s", action.Pair, action.Token)
			if action.Pair == "" {
				log.Logger.Warn("action pair is empty", zap.String("token", action.Token))
				continue
			}

			for {
				err = p.dbService.UpdateToken(action.Token, action.Pair)
				if err != nil {
					log.Logger.Error("update token main pair err", zap.Error(err), zap.String("token", action.Token), zap.String("pair", action.Pair))
					time.Sleep(time.Millisecond * 100)
					continue
				}
				break
			}
		}
	}

	duration := time.Since(now)
	metrics.DbOperationDurationMs.Observe(float64(duration.Milliseconds()))
	log.Logger.Sugar().Debugf("block %d native token price %s", blockInfo.Height, blockInfo.NativeTokenPrice)
	if blockInfo.UsefulInfo() {
		log.Logger.Info("summary",
			zap.Uint64("block", blockResult.Height),
			zap.Float64("db duration", duration.Seconds()),
			zap.Int("tokens", len(blockInfo.NewTokens)),
			zap.Int("pairs", len(blockInfo.NewPairs)),
			zap.Int("actions", len(blockInfo.Actions)),
			zap.Int("txs", len(blockInfo.Txs)))
	}

	err = p.kafkaSender.Send(blockInfo)
	if err != nil {
		log.Logger.Fatal("kafka send msg err", zap.Error(err), zap.Any("block", blockResult.Height))
	}

	p.cache.SetFinishedBlock(blockResult.Height)
	metrics.CurrentHeight.Set(float64(blockResult.Height))
	metrics.TxCntByBlock.Set(float64(len(blockInfo.Txs)))
}

func (p *blockParser) startHandleBlockResult(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			blockContext, ok := <-p.outputQueue
			if !ok {
				log.Logger.Info("commitBlockResultOld - output queue closed")
				return
			}

			p.commitBlockResult(blockContext.BlockResult)
		}
	}()
}

func (p *blockParser) doStop() {
	p.workPool.Release()
	close(p.outputQueue)
}
