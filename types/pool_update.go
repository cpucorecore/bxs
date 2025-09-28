package types

import (
	"bxs/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type PoolUpdate struct {
	LogIndex uint            `json:"log_index"`
	Address  common.Address  `json:"address"`
	Token0   common.Address  `json:"token0"`
	Token1   common.Address  `json:"token1"`
	Amount0  decimal.Decimal `json:"amount0"`
	Amount1  decimal.Decimal `json:"amount1"`
}

func (u *PoolUpdate) Equal(tx *PoolUpdate) bool {
	if u.LogIndex != tx.LogIndex {
		return false
	}
	if u.Address != tx.Address {
		return false
	}
	if u.Token0 != tx.Token0 {
		return false
	}
	if u.Token1 != tx.Token1 {
		return false
	}
	if !util.DecimalEqual(u.Amount0, tx.Amount0) {
		return false
	}
	if !util.DecimalEqual(u.Amount1, tx.Amount1) {
		return false
	}
	return true
}
