package types

import (
	"bxs/repository/orm"
	"github.com/shopspring/decimal"
	"time"
)

type BlockResult struct {
	Height           uint64
	Timestamp        uint64
	BlockTime        time.Time
	NativeTokenPrice decimal.Decimal
	TxResults        []*TxResult
}

func NewBlockResult(height, Timestamp uint64, nativeTokenPrice decimal.Decimal) *BlockResult {
	return &BlockResult{
		Height:           height,
		Timestamp:        Timestamp,
		BlockTime:        time.Unix(int64(Timestamp), 0),
		NativeTokenPrice: nativeTokenPrice,
		TxResults:        make([]*TxResult, 0, 300),
	}
}

func (br *BlockResult) AddTxResult(txResult *TxResult) {
	br.TxResults = append(br.TxResults, txResult)
}

func mergePoolUpdates(poolUpdates []*PoolUpdate) []*PoolUpdate {
	pairAddress2PoolUpdate := make(map[string]*PoolUpdate)
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
	poolUpdates := make([]*PoolUpdate, 0, len(br.TxResults))
	txs := make([]*orm.Tx, 0, len(br.TxResults))
	migratedPools := make([]*MigratedPool, 0, len(br.TxResults))
	actions := make([]*orm.Action, 0, len(br.TxResults))
	ormPairs := make([]*orm.Pair, 0, len(br.TxResults))
	ormTokens := make([]*orm.Token, 0, len(br.TxResults))

	for _, txResult := range br.TxResults {
		poolUpdates = append(poolUpdates, txResult.PoolUpdates...)
		for _, event := range txResult.SwapEvents {
			if event.CanGetTx() {
				txs = append(txs, event.GetTx(br.NativeTokenPrice))
			}
		}

		for _, pair := range txResult.Pairs {
			ormPairs = append(ormPairs, pair.GetOrmPair())
		}

		for _, token := range txResult.Tokens {
			ormTokens = append(ormTokens, token.GetOrmToken())
		}

		migratedPools = append(migratedPools, txResult.MigratedPools...)
		actions = append(actions, txResult.Actions...)
	}

	block := &BlockInfo{
		Height:           br.Height,
		Timestamp:        br.Timestamp,
		NativeTokenPrice: br.NativeTokenPrice.String(),
		Txs:              txs,
		MigratedPools:    migratedPools,
		Actions:          actions,
		NewTokens:        ormTokens,
		NewPairs:         ormPairs,
		PoolUpdates:      mergePoolUpdates(poolUpdates),
	}

	return block
}
