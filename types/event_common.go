package types

import (
	"bxs/repository/orm"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"time"
)

type EventCommon struct {
	Pair            *Pair
	ContractAddress common.Address
	BlockNumber     uint64
	BlockTime       time.Time
	TxHash          common.Hash
	Maker           common.Address
	TxIndex         uint
	LogIndex        uint
}

var _ Event = &EventCommon{}

func (e *EventCommon) GetPairAddress() common.Address {
	return e.ContractAddress
}

func (e *EventCommon) GetPair() *Pair {
	return e.Pair
}

func (e *EventCommon) SetPair(pair *Pair) {
	e.Pair = pair
}

func (e *EventCommon) GetToken0() *Token {
	return nil
}

func (e *EventCommon) CanGetTx() bool {
	return false
}

func (e *EventCommon) GetTx(nativeTokenPrice decimal.Decimal) *orm.Tx {
	return nil
}

func (e *EventCommon) CanGetPoolUpdate() bool {
	return false
}

func (e *EventCommon) GetPoolUpdate() *PoolUpdate {
	return nil
}

func (e *EventCommon) IsCreated() bool {
	return false
}

func (e *EventCommon) SetMaker(maker common.Address) {
	e.Maker = maker
}

func (e *EventCommon) SetBlockTime(blockTime time.Time) {
	e.BlockTime = blockTime
}

func (e *EventCommon) IsMigrated() bool {
	return false
}

func (e *EventCommon) GetAction() *orm.Action {
	return nil
}

func (e *EventCommon) GetNonWBNBToken() common.Address {
	return ZeroAddress
}

func (e *EventCommon) IsPairCreated() bool {
	return false
}

func EventCommonFromEthLog(ethLog *ethtypes.Log) *EventCommon {
	return &EventCommon{
		ContractAddress: ethLog.Address,
		BlockNumber:     ethLog.BlockNumber,
		TxHash:          ethLog.TxHash,
		TxIndex:         ethLog.TxIndex,
		LogIndex:        ethLog.Index,
	}
}
