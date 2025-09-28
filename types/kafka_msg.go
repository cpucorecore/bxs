package types

import (
	"bxs/repository/orm"
	"github.com/ethereum/go-ethereum/common"
)

type MigratedPool struct {
	Pool  common.Address `json:"pool"`
	Token common.Address `json:"token"`
}

type BlockInfo struct {
	Height           uint64          `json:"height"`
	Timestamp        uint64          `json:"timestamp"`
	NativeTokenPrice string          `json:"native_token_price"`
	Txs              []*orm.Tx       `json:"txs"`
	MigratedPools    []*MigratedPool `json:"migrated_pools"`
	Actions          []*orm.Action   `json:"actions"`
	NewTokens        []*orm.Token    `json:"new_tokens"`
	NewPairs         []*orm.Pair     `json:"new_pairs"`
	PoolUpdates      []*PoolUpdate   `json:"pool_updates"`
}

func (bi *BlockInfo) UsefulInfo() bool {
	return len(bi.Txs) != 0 ||
		len(bi.NewTokens) != 0 ||
		len(bi.MigratedPools) != 0 ||
		len(bi.Actions) != 0 ||
		len(bi.NewPairs) != 0
}
