package event_parser

import (
	"bxs/parser/event_parser/event"
	"bxs/types"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrWrongFactory = errors.New("wrong factory")
	ErrNotWBNBPair  = errors.New("not bnb pair")
)

type PairCreatedEventParser struct {
	TopicUnpacker
	FactoryAddress common.Address
}

func (o *PairCreatedEventParser) CheckFactoryAddress(address common.Address) bool {
	return types.IsSameAddress(o.FactoryAddress, address)
}

func (o *PairCreatedEventParser) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	if !o.CheckFactoryAddress(ethLog.Address) {
		return nil, ErrWrongFactory
	}

	eventInput, err := o.unpacker.Unpack(ethLog)
	if err != nil {
		return nil, err
	}

	pairCreatedEvent := &event.PairCreatedEvent{
		EventCommon: types.EventCommonFromEthLog(ethLog),
		Token0Addr:  common.BytesToAddress(ethLog.Topics[1].Bytes()[12:]),
		Token1Addr:  common.BytesToAddress(ethLog.Topics[2].Bytes()[12:]),
		PairAddr:    eventInput[0].(common.Address),
	}

	if !pairCreatedEvent.IsWBNBPair() {
		return nil, ErrNotWBNBPair
	}

	return pairCreatedEvent, nil
}
