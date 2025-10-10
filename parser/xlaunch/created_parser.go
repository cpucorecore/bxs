package event_parser

import (
	"bxs/chain_params"
	pcommon "bxs/parser/common"
	"bxs/types"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

var (
	ErrWrongFactoryAddress = errors.New("wrong factory address")
)

type CreatedEventParser struct {
	pcommon.TopicUnpacker
}

func checkFactoryAddr(addr common.Address) bool {
	return types.IsSameAddress(addr, chain_params.G.XLaunchFactoryAddress)
}

func (o *CreatedEventParser) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	if !checkFactoryAddr(ethLog.Address) {
		return nil, ErrWrongFactoryAddress
	}

	eventInput, err := o.Unpacker.Unpack(ethLog)
	if err != nil {
		return nil, err
	}

	createdEvent := &CreatedEvent{
		EventCommon:         types.EventCommonFromEthLog(ethLog),
		PoolAddress:         common.BytesToAddress(ethLog.Topics[1].Bytes()[12:]),
		Creator:             common.BytesToAddress(ethLog.Topics[2].Bytes()[12:]),
		TokenAddress:        common.BytesToAddress(ethLog.Topics[3].Bytes()[12:]),
		BaseTokenInitAmount: eventInput[0].(*big.Int),
		TokenInitAmount:     eventInput[1].(*big.Int),
		TotalSupply:         eventInput[2].(*big.Int),
		Name:                eventInput[3].(string),
		Symbol:              eventInput[4].(string),
		Cid:                 eventInput[5].(string),
		Tid:                 eventInput[6].(string),
		Description:         eventInput[7].(string),
		Telegram:            eventInput[8].(string),
		Twitter:             eventInput[9].(string),
		Website:             eventInput[10].(string),
	}

	createdEvent.FormatString()

	createdEvent.Pair = createdEvent.getPair()
	return createdEvent, nil
}
