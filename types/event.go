package types

import (
	"bxs/repository/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"time"
)

type Event interface {
	CanGetPair() bool
	GetPair() *Pair
	GetPairAddress() common.Address
	SetPair(pair *Pair)
	CanGetToken0() bool
	GetToken0() *Token
	SetMaker(maker common.Address)
	SetBlockTime(blockTime time.Time)

	CanGetTx() bool
	GetTx(nativeTokenPrice decimal.Decimal) *orm.Tx

	CanGetPoolUpdate() bool
	GetPoolUpdate() *PoolUpdate

	IsCreatePair() bool
	IsMigrated() bool
}
