package parser

import (
	"bxs/cache"
	"bxs/config"
	"bxs/log"
	"bxs/metrics"
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

type TxResultAndPairWrap struct {
	TxResult  *types.TxResult
	PairWraps []*types.PairWrap
	Pairs     []*types.Pair
	Tokens    []*types.Token
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

	r := types.NewTxResult(txSender, pbc.HeightTime.Time)
	pairs := make([]*types.Pair, 0, 2)
	tokens := make([]*types.Token, 0, 2)
	for _, ethLog := range txReceipt.Logs {
		if len(ethLog.Topics) == 0 {
			continue
		}

		event, parseErr := p.topicRouter.Parse(ethLog)
		if parseErr != nil {
			continue
		}

		if event.IsCreated() {
			pair := event.GetPair()
			token0 := event.GetToken0()
			p.cache.SetPair(pair)
			p.cache.SetToken(token0)
			pairs = append(pairs, pair)
			tokens = append(tokens, token0)
			continue
		} else {
			pairAddr := event.GetPairAddress()
			pair, ok := p.cache.GetPair(pairAddr)
			if !ok {
				log.Logger.Sugar().Warnf("get Pool %s err", pairAddr)
				pw := p.pairService.GetPair(pairAddr)
				pair = pw.Pair
				if pair == nil {
					log.Logger.Sugar().Fatalf("get pair %s fail", pairAddr)
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
		}

		r.SetPairs(pairs)
		r.SetTokens(tokens)
		r.AddEvent(event)
	}

	return r
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
	err := p.dbService.AddTokens(blockInfo.NewTokens)
	if err != nil {
		log.Logger.Fatal("add tokens err", zap.Any("height", blockInfo.Height), zap.Error(err))
	}

	err = p.dbService.AddPairs(blockInfo.NewPairs)
	if err != nil {
		log.Logger.Fatal("add pairs err", zap.Any("height", blockInfo.Height), zap.Error(err))
	}

	err = p.dbService.AddTxs(blockInfo.Txs)
	if err != nil {
		log.Logger.Fatal("add txs err", zap.Any("height", blockInfo.Height), zap.Error(err))
	}
	err = p.dbService.AddActions(blockInfo.Actions)
	if err != nil {
		log.Logger.Fatal("add actions err", zap.Any("height", blockInfo.Height), zap.Error(err))
	}

	duration := time.Since(now)
	metrics.DbOperationDurationMs.Observe(float64(duration.Milliseconds()))
	log.Logger.Sugar().Debugf("block %d native token price %s", blockInfo.Height, blockInfo.NativeTokenPrice)
	if blockInfo.CatchInfo() {
		log.Logger.Info("db operation duration",
			zap.Uint64("block", blockResult.Height),
			zap.Float64("duration", duration.Seconds()),
			zap.Int("new tokens", len(blockInfo.NewTokens)),
			zap.Int("new pairs", len(blockInfo.NewPairs)),
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
