package types

import (
	"github.com/ethereum/go-ethereum/common"
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
	Maker                   common.Address
	PairCreatedEvents       []Event
	PairAddress2TxPairEvent map[common.Address]*TxPairEvent
}

func NewTxResult(maker common.Address) *TxResult {
	return &TxResult{
		Maker:                   maker,
		PairCreatedEvents:       make([]Event, 0, 10),
		PairAddress2TxPairEvent: make(map[common.Address]*TxPairEvent),
	}
}

func (tr *TxResult) AddEvent(event Event) {
	event.SetMaker(tr.Maker)
	if event.IsCreatePair() {
		tr.PairCreatedEvents = append(tr.PairCreatedEvents, event)
	}

	pairAddress := event.GetPairAddress()
	txPairEvent, ok := tr.PairAddress2TxPairEvent[pairAddress]
	if ok {
		txPairEvent.AddEvent(event)
		return
	}

	txPairEvent = &TxPairEvent{}
	txPairEvent.AddEvent(event)
	tr.PairAddress2TxPairEvent[pairAddress] = txPairEvent
}
