package parser

import (
	"bxs/cache"
	"bxs/chain_params"
	"bxs/config"
	"bxs/logger"
	"bxs/metrics"
	"bxs/sequencer"
	"bxs/service"
	"bxs/types"
	"encoding/json"
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"sync"
	"time"
)

type BlockParser interface {
	Start(*sync.WaitGroup)
	Stop()
	ParseBlockAsync(bw *types.BlockContext)
}

type blockParser struct {
	cache        cache.Cache
	sequencer    sequencer.Sequencer
	priceService service.PriceService
	topicRouter  TopicRouter
	kafkaSender  service.KafkaSender
	dbService    service.DBService
	inputQueue   chan *types.BlockContext
	outputQueue  chan *types.BlockContext
}

func NewBlockParser(
	cache cache.Cache,
	sequencer sequencer.Sequencer,
	priceService service.PriceService,
	topicRouter TopicRouter,
	kafkaSender service.KafkaSender,
	dbService service.DBService,
) BlockParser {
	return &blockParser{
		cache:        cache,
		sequencer:    sequencer,
		priceService: priceService,
		topicRouter:  topicRouter,
		kafkaSender:  kafkaSender,
		dbService:    dbService,
		inputQueue:   make(chan *types.BlockContext, config.G.BlockHandler.QueueSize),
		outputQueue:  make(chan *types.BlockContext, config.G.BlockHandler.QueueSize),
	}
}

func (p *blockParser) Commit(x sequencer.Sequenceable) {
	p.outputQueue <- x.(*types.BlockContext)
}

func (p *blockParser) Start(waitGroup *sync.WaitGroup) {
	p.startCommitBlockResult(waitGroup)

	go func() {
	For:
		for {
			select {
			case bc, ok := <-p.inputQueue:
				if !ok {
					logger.G.Info("block parser inputQueue is closed")
					break For
				}
				p.parseBlock(bc)
			}
		}

		logger.G.Info("all block parse task finish")
		p.doStop()
	}()
}

func (p *blockParser) Stop() {
	close(p.inputQueue)
}

func (p *blockParser) ParseBlockAsync(bc *types.BlockContext) {
	p.inputQueue <- bc
}

func (p *blockParser) getNativeTokenPrice(blockNumber *big.Int, blockTimestamp uint64) decimal.Decimal {
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

func (p *blockParser) parseTxReceipt(bc *types.BlockContext, receipt *ethtypes.Receipt) *types.TxResult {
	tr := bc.NewTxResult(receipt.TransactionIndex)
	for _, log := range receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}

		event, err := p.topicRouter.Route(log)
		if err != nil {
			continue
		}

		bc.DecorateEvent(event)

		if event.IsCreated() { // xLaunch event: created
			pair := event.GetPair()
			token0 := event.GetToken0()
			p.cache.SetPair(pair)
			p.cache.SetToken(token0)
			tr.AddPair(pair)
			tr.AddToken(token0)
			tr.AddPoolUpdate(event.GetPoolUpdate())
			continue
		}

		if event.IsPairCreated() { // pancakeV2 event: PairCreated
			pair := event.GetPair()
			token0, ok := p.cache.GetToken(pair.Token0.Address)
			if !ok {
				logger.G.Sugar().Infof("pair %s has no xLaunch token, ignore it, token0 %s", pair.Address, pair.Token0.Address)
				pair.Filtered = true
				pair.FilterCode = types.FilterCodeNoXLaunchToken
				p.cache.SetPair(pair)
				continue
			}

			pair.Token0.Symbol = token0.Symbol
			pair.Token0.Decimal = token0.Decimals
			pair.Token1.Symbol = types.WBNBSymbol
			pair.Token1.Decimal = types.WBNBDecimal
			p.cache.SetPair(pair)
			tr.AddPairCreatedEvent(event)
			tr.AddPair(pair)
			continue
		}

		if event.IsBuyOrSell() { // xLaunch event: buy/sell
			pairAddr := event.GetPairAddress()
			pair, ok := p.cache.GetPair(pairAddr)
			if !ok {
				logger.G.Sugar().Warnf("pool %s not cached", pairAddr)
				continue
			}

			if event.IsMigrated() {
				p.cache.SetMigrateToken(pair.Token0.Address)
				tr.AddMigratedPool(&types.MigratedPool{
					Pool:  pairAddr.String(),
					Token: pair.Token0.Address.String(),
				})
			}

			event.SetPair(pair)
			tr.AddPoolUpdate(event.GetPoolUpdate())
			tr.AddSwapEvent(event)
			continue
		}

		if event.IsSwap() { // pancakeV2 event Swap
			pairAddr := event.GetPairAddress()
			pair, ok := p.cache.GetPair(pairAddr)
			if !ok {
				logger.G.Sugar().Warnf("pair %s not cached, ignore it, tx hash %s", pairAddr, event.GetTxHash())
				continue
			}

			if pair.Filtered {
				logger.G.Sugar().Infof("pair %s is filtered, filter code %d, tx hash %s", event.GetPairAddress(), pair.FilterCode, event.GetTxHash())
				continue
			}

			event.SetPair(pair)
			tr.AddSwapEvent(event)
			continue
		}

		if event.IsSync() { // pancakeV2 event Sync
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
			tr.AddPoolUpdate(event.GetPoolUpdate())
		}
	}

	p.processTxPairCreatedEvents(tr)
	return tr
}

func (p *blockParser) processTxPairCreatedEvents(tr *types.TxResult) *types.TxResult {
	for _, event := range tr.PairCreatedEvents {
		token0 := event.GetNonWBNBToken()
		if p.cache.MigrateTokenExist(token0) {
			tr.AddAction(event.GetAction())
			p.cache.DelMigrateToken(token0)
		}
	}
	return tr
}

func (p *blockParser) preParseBlock(bc *types.BlockContext) {
	bc.NativeTokenPrice = p.getNativeTokenPrice(bc.HeightTime.HeightBigInt, bc.HeightTime.Timestamp)

	signer := ethtypes.MakeSigner(chain_params.G.ChainConfig, bc.HeightTime.HeightBigInt, bc.HeightTime.Timestamp)
	for idx, tx := range bc.Transactions {
		sender, err := ethtypes.Sender(signer, tx)
		if err != nil {
			logger.G.Sugar().Fatalf("get tx [%v] sender err [%s]", tx, err)
		}
		bc.Senders[idx] = sender
	}
}

func (p *blockParser) parseBlock(bc *types.BlockContext) {
	p.preParseBlock(bc)

	now := time.Now()
	for _, receipt := range bc.Receipts {
		if receipt.Status != 1 {
			continue
		}
		bc.SetTxResult(receipt.TransactionIndex, p.parseTxReceipt(bc, receipt))
	}
	duration := time.Since(now).Milliseconds()
	metrics.ParseBlockDurationMs.Observe(float64(duration))
	logger.G.Info(fmt.Sprintf("parse block %d duration %dms", bc.HeightTime.HeightBigInt, duration))

	p.sequencer.CommitWithSequence(bc, p)
}

func (p *blockParser) commitBlockResult(bc *types.BlockContext) {
	blockInfo := bc.GetKafkaMsg()

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
			zap.Uint64("block", bc.HeightTime.Height),
			zap.Float64("db duration", duration.Seconds()),
			zap.Int("tokens", len(blockInfo.NewTokens)),
			zap.Int("pairs", len(blockInfo.NewPairs)),
			zap.Int("actions", len(blockInfo.Actions)),
			zap.Int("txs", len(blockInfo.Txs)))

		bytes, _ := json.Marshal(blockInfo)
		logger.G.Sugar().Debugf("%s", string(bytes))
	}

	err = p.kafkaSender.Send(blockInfo)
	if err != nil {
		logger.G.Fatal("send kafka msg err", zap.Error(err), zap.Any("block", bc.HeightTime.Height))
	}

	p.cache.SetFinishedBlock(bc.HeightTime.Height)
	metrics.CurrentHeight.Set(float64(bc.HeightTime.Height))
	metrics.TxCntByBlock.Set(float64(len(blockInfo.Txs)))
}

func (p *blockParser) startCommitBlockResult(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			bc, ok := <-p.outputQueue
			if !ok {
				logger.G.Info("commitBlockResult - output queue closed")
				return
			}

			p.commitBlockResult(bc)
		}
	}()
}

func (p *blockParser) doStop() {
	close(p.outputQueue)
}
