package types

import (
	"bxs/repository/orm"
)

type BlockInfo struct {
	Height           uint64
	Timestamp        uint64
	NativeTokenPrice string
	Txs              []*orm.Tx
	Actions          []*orm.Action
	NewTokens        []*orm.Token
	NewPairs         []*orm.Pair
	PoolUpdates      []*PoolUpdate
}

func (bi *BlockInfo) CatchInfo() bool {
	return len(bi.Txs) != 0 || len(bi.NewTokens) != 0 || len(bi.Actions) != 0 || len(bi.NewPairs) != 0
}
