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

func (o *CreatedEventParser) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	if !types.IsSameAddress(ethLog.Address, chain_params.G.XLaunchFactoryAddress) {
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
		URL:                 eventInput[5].(string),
		Description:         eventInput[6].(string),
		Telegram:            eventInput[7].(string),
		Twitter:             eventInput[8].(string),
		Website:             eventInput[9].(string),
	}

	createdEvent.FormatString()

	createdEvent.Pair = createdEvent.DoGetPair()
	return createdEvent, nil
}
