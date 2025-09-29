package types

import (
	"bxs/repository/orm"
	"github.com/ethereum/go-ethereum/common"
)

type TxResult struct {
	Sender            common.Address
	SwapEvents        []Event
	PairCreatedEvents []Event
	PoolUpdates       []*PoolUpdate
	Pairs             []*Pair
	Tokens            []*Token
	MigratedPools     []*MigratedPool
	Actions           []*orm.Action
}

func (r *TxResult) decorateEvent(event Event) {
	event.SetMaker(r.Sender)
}

func (r *TxResult) AddSwapEvent(event Event) {
	r.decorateEvent(event)
	r.SwapEvents = append(r.SwapEvents, event)
}

func (r *TxResult) AddPairCreatedEvent(event Event) {
	r.decorateEvent(event)
	r.PairCreatedEvents = append(r.PairCreatedEvents, event)
}

func (r *TxResult) AddPoolUpdate(poolUpdate *PoolUpdate) {
	r.PoolUpdates = append(r.PoolUpdates, poolUpdate)
}

func (r *TxResult) AddPair(pair *Pair) {
	r.Pairs = append(r.Pairs, pair)
}

func (r *TxResult) AddToken(token *Token) {
	r.Tokens = append(r.Tokens, token)
}

func (r *TxResult) AddMigratedPool(pool *MigratedPool) {
	r.MigratedPools = append(r.MigratedPools, pool)
}

func (r *TxResult) AddAction(action *orm.Action) {
	r.Actions = append(r.Actions, action)
}
