package parser

import (
	"bxs/cache"
	"bxs/chain_params"
	"bxs/config"
	"bxs/logger"
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
	ParseBlockAsync(bw *types.BlockCtx)
}

type blockParser struct {
	inputQueue   chan *types.BlockCtx
	workPool     *ants.Pool
	cache        cache.Cache
	sequencer    sequencer.Sequencer
	outputQueue  chan *types.BlockCtx
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
		logger.G.Fatal("ants.NewPool err", zap.Error(err))
	}

	parseTxPool, err := ants.NewPool(config.G.BlockHandler.ParseTxPoolSize)
	if err != nil {
		logger.G.Fatal("ants.NewPool err", zap.Error(err))
	}

	return &blockParser{
		inputQueue:   make(chan *types.BlockCtx, config.G.BlockHandler.QueueSize),
		workPool:     workPool,
		cache:        cache,
		sequencer:    sequencer,
		outputQueue:  make(chan *types.BlockCtx, config.G.BlockHandler.QueueSize),
		priceService: priceService,
		pairService:  pairService,
		topicRouter:  topicRouter,
		kafkaSender:  kafkaSender,
		dbService:    dbService,
		parseTxPool:  parseTxPool,
	}
}

func (p *blockParser) Commit(x sequencer.Sequenceable) {
	p.outputQueue <- x.(*types.BlockCtx)
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
					logger.G.Info("block handler inputQueue is closed")
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
		logger.G.Info("all block parse task finish")
		p.doStop()
	}()
}

func (p *blockParser) Stop() {
	close(p.inputQueue)
}

func (p *blockParser) ParseBlockAsync(bw *types.BlockCtx) {
	p.inputQueue <- bw
}

func (p *blockParser) waitForNativeTokenPrice(blockNumber *big.Int, blockTimestamp uint64) decimal.Decimal {
	for {
		price, err := p.priceService.GetPrice(blockNumber)
		if err != nil {
			logger.G.Error("get price err", zap.Error(err), zap.Any("blockNumber", blockNumber), zap.Any("blockTimestamp", blockTimestamp))
			time.Sleep(time.Second)
			continue
		}
		return price
	}
}
func (p *blockParser) parseTxReceipt(pbc *types.BlockCtx, receipt *ethtypes.Receipt) *types.TxResult {
	sender, err := pbc.GetTxSender(receipt.TransactionIndex)
	if err != nil {
		logger.G.Fatal("get tx sender err",
			zap.Error(err),
			zap.Any("pbc", pbc),
			zap.Any("receipt", receipt),
		)
	}

	txResult := types.NewTxResult(sender, pbc.HeightTime.Time)
	pairs := make([]*types.Pair, 0, 2)
	tokens := make([]*types.Token, 0, 2)
	migratedPools := make([]*types.MigratedPool, 0, 2)
	for _, log := range receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}

		event, parseErr := p.topicRouter.Parse(log)
		if parseErr != nil {
			continue
		}

		if event.IsCreated() {
			// xlaunch
			// save token
			// save pair
			// pool update
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
			// pancakev2
			// get pair
			// if filtered: continue

			// query token cache(pair.token0)
			// if token not exist, pair filtered, and update pair cache, continue

			pair := event.GetPair()
			token0, ok := p.cache.GetToken(pair.Token0.Address)
			if !ok {
				logger.G.Sugar().Infof("pair %s have no xlaunch token, ignore it, token0 %s", pair.Address, pair.Token0.Address)
				pair.Filtered = true
				pair.FilterCode = types.FilterCodeNoXLaunchToken
			} else {
				pair.Token0 = &types.TokenTinyInfo{
					Address: token0.Address,
					Symbol:  token0.Symbol,
					Decimal: token0.Decimals,
				}
				pair.Token1 = &types.TokenTinyInfo{
					Address: chain_params.G.WBNBAddress,
					Symbol:  "WBNB",
					Decimal: 18,
				}
				txResult.AddPairCreatedEvent(event)
				pairs = append(pairs, pair)
			}

			p.cache.SetPair(pair)
			continue
		}

		// xlaunch buy/sell
		if event.IsBuyOrSell() {
			pairAddr := event.GetPairAddress()
			pair, ok := p.cache.GetPair(pairAddr)
			if !ok {
				logger.G.Sugar().Warnf("pool %s not cached, it should be", pairAddr)
				pw := p.pairService.GetPair(pairAddr)
				pair = pw.Pair
				if pair == nil {
					logger.G.Sugar().Fatalf("get pair %s fail", pairAddr)
				}

				if pair.Filtered {
					logger.G.Sugar().Warnf("pair %s is filtered, filter code %d", pairAddr, pair.FilterCode)
					continue
				}

				pairs = append(pairs, pair)
				tokens = append(tokens, &types.Token{
					Address:  pair.Token0.Address,
					Symbol:   pair.Token0.Symbol,
					Decimals: pair.Token0.Decimal,
					Program:  types.ProtocolNameXLaunch,
				})
			}
			event.SetPair(pair)

			if event.IsMigrated() {
				p.cache.SetMigrateToken(event.GetPair().Token0.Address)
				migratedPools = append(migratedPools, &types.MigratedPool{
					Pool:  event.GetPair().Address,
					Token: event.GetPair().Token0.Address,
				})
			}

			txResult.AddPoolUpdate(event.GetPoolUpdate())
			txResult.AddSwapEvent(event)
		} else if event.IsSwap() {
			pair, ok := p.cache.GetPair(event.GetPairAddress())
			if !ok {
				logger.G.Sugar().Warnf("pair %s not cached, ignore it, tx hash %s", event.GetPairAddress(), event.GetTxHash())
				continue
			}

			if pair.Filtered {
				logger.G.Sugar().Infof("pair %s is filtered, filter code %d, tx hash %s", event.GetPairAddress(), pair.FilterCode, event.GetTxHash())
				continue
			}

			event.SetPair(pair)
			txResult.AddSwapEvent(event)
		} else if event.IsSync() {
			pair, ok := p.cache.GetPair(event.GetPairAddress())
			if !ok {
				logger.G.Sugar().Warnf("pair %s not cached, ignore it, tx hash %s", event.GetPairAddress(), event.GetTxHash())
				continue
			}

			if pair.Filtered {
				logger.G.Sugar().Infof("pair %s is filtered, filter code %d, tx hash %s", event.GetPairAddress(), pair.FilterCode, event.GetTxHash())
				continue
			}

			event.SetPair(pair)
			txResult.AddPoolUpdate(event.GetPoolUpdate())
		}
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

func (p *blockParser) parseBlock(pbc *types.BlockCtx) {
	pbc.NativeTokenPrice = p.waitForNativeTokenPrice(pbc.HeightTime.HeightBigInt, pbc.HeightTime.Timestamp)

	now := time.Now()
	br := types.NewBlockResult(pbc.HeightTime.Height, pbc.HeightTime.Timestamp, pbc.NativeTokenPrice)

	for _, receipt := range pbc.Receipts {
		if receipt.Status != 1 {
			continue
		}
		br.AddTxResult(p.parseTxReceipt(pbc, receipt))
	}

	duration := time.Since(now)
	metrics.ParseBlockDurationMs.Observe(float64(duration.Milliseconds()))
	logger.G.Info(fmt.Sprintf("parse block %d duration %dms", pbc.HeightTime.HeightBigInt, duration.Milliseconds()))

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
			logger.G.Fatal("add tokens err", zap.Any("height", blockInfo.Height), zap.Error(err))
		}
	}

	if len(blockInfo.NewPairs) > 0 {
		err = p.dbService.AddPairs(blockInfo.NewPairs)
		if err != nil {
			logger.G.Fatal("add pairs err", zap.Any("height", blockInfo.Height), zap.Error(err))
		}
	}

	if len(blockInfo.Txs) > 0 {
		err = p.dbService.AddTxs(blockInfo.Txs)
		if err != nil {
			logger.G.Fatal("add txs err", zap.Any("height", blockInfo.Height), zap.Error(err))
		}
	}

	if len(blockInfo.Actions) > 0 {
		err = p.dbService.AddActions(blockInfo.Actions)
		if err != nil {
			logger.G.Fatal("add actions err", zap.Any("height", blockInfo.Height), zap.Error(err))
		}

		for _, action := range blockInfo.Actions {
			logger.G.Sugar().Infof("add action: pair:%s, token:%s", action.Pair, action.Token)
			if action.Pair == "" {
				logger.G.Warn("action pair is empty", zap.String("token", action.Token))
				continue
			}

			for {
				err = p.dbService.UpdateToken(action.Token, action.Pair)
				if err != nil {
					logger.G.Error("update token main pair err", zap.Error(err), zap.String("token", action.Token), zap.String("pair", action.Pair))
					time.Sleep(time.Millisecond * 100)
					continue
				}
				break
			}
		}
	}

	duration := time.Since(now)
	metrics.DbOperationDurationMs.Observe(float64(duration.Milliseconds()))
	logger.G.Sugar().Debugf("block %d native token price %s", blockInfo.Height, blockInfo.NativeTokenPrice)
	if blockInfo.UsefulInfo() {
		logger.G.Info("summary",
			zap.Uint64("block", blockResult.Height),
			zap.Float64("db duration", duration.Seconds()),
			zap.Int("tokens", len(blockInfo.NewTokens)),
			zap.Int("pairs", len(blockInfo.NewPairs)),
			zap.Int("actions", len(blockInfo.Actions)),
			zap.Int("txs", len(blockInfo.Txs)))
	}

	err = p.kafkaSender.Send(blockInfo)
	if err != nil {
		logger.G.Fatal("kafka send msg err", zap.Error(err), zap.Any("block", blockResult.Height))
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
				logger.G.Info("commitBlockResultOld - output queue closed")
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
