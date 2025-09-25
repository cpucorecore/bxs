package types

import (
	"bxs/repository/orm"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type TxPairEvent struct {
	XLaunch []Event
}

func (tpe *TxPairEvent) AddEvent(event Event) {
	if tpe.XLaunch == nil {
		tpe.XLaunch = make([]Event, 0, 10)
	}
	tpe.XLaunch = append(tpe.XLaunch, event)
}

type TxResult struct {
	BlockTime         time.Time
	Sender            common.Address
	SwapEvents        []Event
	PairCreatedEvents []Event
	PoolUpdates       []*PoolUpdate
	Pairs             []*Pair
	Tokens            []*Token
	MigratedPools     []*MigratedPool
	Actions           []*orm.Action
}

func NewTxResult(sender common.Address, blockTime time.Time) *TxResult {
	return &TxResult{
		Sender:            sender,
		BlockTime:         blockTime,
		SwapEvents:        make([]Event, 0, 32),
		PairCreatedEvents: make([]Event, 0, 32),
		PoolUpdates:       make([]*PoolUpdate, 0, 32),
	}
}

func (r *TxResult) AddSwapEvent(event Event) {
	event.SetMaker(r.Sender)
	event.SetBlockTime(r.BlockTime)
	r.SwapEvents = append(r.SwapEvents, event)
}

func (r *TxResult) AddPairCreatedEvent(event Event) {
	event.SetMaker(r.Sender)
	event.SetBlockTime(r.BlockTime)
	r.PairCreatedEvents = append(r.PairCreatedEvents, event)
}

func (r *TxResult) SetPairs(pairs []*Pair) {
	r.Pairs = pairs
}

func (r *TxResult) SetTokens(tokens []*Token) {
	r.Tokens = tokens
}

func (r *TxResult) SetMigratedPools(pools []*MigratedPool) {
	r.MigratedPools = pools
}

func (r *TxResult) SetActions(actions []*orm.Action) {
	r.Actions = actions
}

func (r *TxResult) AddPoolUpdate(poolUpdate *PoolUpdate) {
	r.PoolUpdates = append(r.PoolUpdates, poolUpdate)
}
