package types

import (
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
	BlockTime time.Time
	Maker     common.Address
	Events    []Event
	Pairs     []*Pair
	Tokens    []*Token
}

func NewTxResult(maker common.Address, blockTime time.Time) *TxResult {
	return &TxResult{
		Maker:     maker,
		BlockTime: blockTime,
		Events:    make([]Event, 0, 32),
	}
}

func (r *TxResult) AddEvent(event Event) {
	event.SetMaker(r.Maker)
	event.SetBlockTime(r.BlockTime)
	r.Events = append(r.Events, event)
}

func (r *TxResult) SetPairs(pairs []*Pair) {
	r.Pairs = pairs
}

func (r *TxResult) SetTokens(tokens []*Token) {
	r.Tokens = tokens
}
