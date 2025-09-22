package event_parser

import (
	"bxs/abi/xlaunch"
	"bxs/parser/event_parser/event"
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

type CreatedEventParser struct {
	FactoryEventParser
}

func (o *CreatedEventParser) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	pair := &types.Pair{}

	if !types.IsSameAddress(ethLog.Address, xlaunch.FactoryAddress) {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeWrongFactory
		return nil, ErrWrongFactoryAddress
	}

	input, err := o.LogUnpacker.Unpack(ethLog)
	if err != nil {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeUnpackDataErr
		return nil, err
	}

	createdEvent := &event.CreatedEvent{
		EventCommon:         types.EventCommonFromEthLog(ethLog),
		PoolAddress:         common.BytesToAddress(ethLog.Topics[1].Bytes()[12:]),
		Creator:             common.BytesToAddress(ethLog.Topics[2].Bytes()[12:]),
		TokenAddress:        common.BytesToAddress(ethLog.Topics[3].Bytes()[12:]),
		BaseTokenInitAmount: input[0].(*big.Int),
		TokenInitAmount:     input[1].(*big.Int),
		Name:                input[2].(string),
		Symbol:              input[3].(string),
		URL:                 input[4].(string),
		Description:         input[5].(string),
	}

	pair.Address = createdEvent.PoolAddress
	pair.Token0Core = &types.TokenCore{
		Address:  createdEvent.TokenAddress,
		Symbol:   createdEvent.Symbol,
		Decimals: types.DefaultDecimals,
	}
	pair.Token1Core = types.NativeTokenCore
	pair.Block = ethLog.BlockNumber
	pair.ProtocolId = types.ProtocolIdXLaunch

	pair.FilterByToken0AndToken1()

	createdEvent.EventCommon.Pair.Token0 = &types.Token{
		Address:     createdEvent.TokenAddress,
		Creator:     createdEvent.Creator,
		Name:        createdEvent.Name,
		Symbol:      createdEvent.Symbol,
		Decimals:    types.DefaultDecimals,
		BlockNumber: createdEvent.BlockNumber,
		BlockTime:   createdEvent.BlockTime,
		Program:     types.ProtocolNameXLaunch,
		Filtered:    false,
		Timestamp:   time.Now(),
		URL:         createdEvent.URL,
		Description: createdEvent.Description,
	}

	createdEvent.Pair = pair

	return createdEvent, nil
}
