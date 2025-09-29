package types

import (
	"bxs/repository/orm"
)

type MigratedPool struct {
	Pool  string `json:"pool"`
	Token string `json:"token"`
}

type KafkaMsg struct {
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

func (bi *KafkaMsg) UsefulInfo() bool {
	return len(bi.Txs) != 0 ||
		len(bi.NewTokens) != 0 ||
		len(bi.MigratedPools) != 0 ||
		len(bi.Actions) != 0 ||
		len(bi.NewPairs) != 0
}
