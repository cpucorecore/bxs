package types

import (
	"bxs/repository/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"time"
)

type Event interface {
	GetPairAddress() common.Address
	GetPair() *Pair
	SetPair(pair *Pair)
	GetToken0() *Token
	SetMaker(maker common.Address)
	SetBlockTime(blockTime time.Time)

	CanGetTx() bool
	GetTx(nativeTokenPrice decimal.Decimal) *orm.Tx

	CanGetPoolUpdate() bool
	GetPoolUpdate() *PoolUpdate

	IsCreated() bool
	IsMigrated() bool
	GetAction() *orm.Action

	IsPairCreated() bool
	GetNonWBNBToken() common.Address
}
