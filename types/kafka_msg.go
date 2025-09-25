package types

import (
	"bxs/repository/orm"
	"github.com/ethereum/go-ethereum/common"
)

type MigratedPool struct {
	Pool  common.Address
	Token common.Address
}

type BlockInfo struct {
	Height           uint64
	Timestamp        uint64
	NativeTokenPrice string
	Txs              []*orm.Tx
	MigratedPools    []*MigratedPool
	Actions          []*orm.Action
	NewTokens        []*orm.Token
	NewPairs         []*orm.Pair
	PoolUpdates      []*PoolUpdate
}

func (bi *BlockInfo) UsefulInfo() bool {
	return len(bi.Txs) != 0 ||
		len(bi.NewTokens) != 0 ||
		len(bi.MigratedPools) != 0 ||
		len(bi.Actions) != 0 ||
		len(bi.NewPairs) != 0
}
