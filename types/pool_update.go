package types

import (
	"bxs/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type PoolUpdate struct {
	LogIndex uint
	Address  common.Address
	Token0   common.Address
	Token1   common.Address
	Amount0  decimal.Decimal
	Amount1  decimal.Decimal
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
