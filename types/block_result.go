package types

import (
	"bxs/repository/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"time"
)

type BlockResult struct {
	Height           uint64
	Timestamp        uint64
	BlockTime        time.Time
	NativeTokenPrice decimal.Decimal
	NewPairs         map[common.Address]*Pair
	NewTokens        map[common.Address]*Token
	TxResults        []*TxResult
}

func NewBlockResult(height, Timestamp uint64, nativeTokenPrice decimal.Decimal) *BlockResult {
	return &BlockResult{
		Height:           height,
		Timestamp:        Timestamp,
		BlockTime:        time.Unix(int64(Timestamp), 0),
		NativeTokenPrice: nativeTokenPrice,
		NewPairs:         make(map[common.Address]*Pair),
		NewTokens:        make(map[common.Address]*Token),
		TxResults:        make([]*TxResult, 0, 200),
	}
}

func (br *BlockResult) AddTxResult(txResult *TxResult) {
	br.TxResults = append(br.TxResults, txResult)
}

func (br *BlockResult) getAllEvents() []Event {
	events := make([]Event, 0, 500)
	for _, txResult := range br.TxResults {
		for _, txPairEvent := range txResult.PairAddress2TxPairEvent {
			events = append(events, txPairEvent.UniswapV2...)
			events = append(events, txPairEvent.UniswapV3...)
			events = append(events, txPairEvent.PancakeV2...)
			events = append(events, txPairEvent.PancakeV3...)
			events = append(events, txPairEvent.Aerodrome...)
			events = append(events, txPairEvent.XLaunch...)
		}
	}
	return events
}

func mergePoolUpdates(poolUpdates []*PoolUpdate) []*PoolUpdate {
	pairAddress2PoolUpdate := make(map[common.Address]*PoolUpdate)
	for _, poolUpdate := range poolUpdates {
		poolUpdate_, ok := pairAddress2PoolUpdate[poolUpdate.Address]
		if ok {
			if poolUpdate.LogIndex > poolUpdate_.LogIndex {
				pairAddress2PoolUpdate[poolUpdate.Address] = poolUpdate
			}
		} else {
			pairAddress2PoolUpdate[poolUpdate.Address] = poolUpdate
		}
	}
	poolUpdatesMerged := make([]*PoolUpdate, 0, len(pairAddress2PoolUpdate))
	for _, pu := range pairAddress2PoolUpdate {
		poolUpdatesMerged = append(poolUpdatesMerged, pu)
	}
	return poolUpdatesMerged
}

func (br *BlockResult) GetKafkaMessage() *BlockInfo {
	events := br.getAllEvents()

	txs := make([]*orm.Tx, 0, len(events))
	newPairs := make([]*Pair, 0, 10)
	poolUpdates := make([]*PoolUpdate, 0, 40)
	for _, event := range events {
		if event.IsCreatePair() {
			newPairs = append(newPairs, event.GetPair())
			continue
		}

		if event.CanGetTx() {
			txs = append(txs, event.GetTx(br.NativeTokenPrice))
		}

		if event.CanGetPoolUpdate() {
			poolUpdates = append(poolUpdates, event.GetPoolUpdate())
		}
	}

	for _, pair := range newPairs {
		br.NewPairs[pair.Address] = pair
	}

	ormTokens := make([]*orm.Token, 0, len(br.NewTokens))
	for _, token := range br.NewTokens {
		ormTokens = append(ormTokens, token.GetOrmToken())
	}

	ormPairs := make([]*orm.Pair, 0, len(newPairs))
	for _, pair := range br.NewPairs {
		ormPairs = append(ormPairs, pair.GetOrmPair())
	}

	poolUpdatesMerged := mergePoolUpdates(poolUpdates)

	block := &BlockInfo{
		Height:           br.Height,
		Timestamp:        br.Timestamp,
		NativeTokenPrice: br.NativeTokenPrice.String(),
		Txs:              txs,
		NewTokens:        ormTokens,
		NewPairs:         ormPairs,
		PoolUpdates:      poolUpdatesMerged,
	}

	return block
}
