package event_parser

import (
	"bxs/abi"
	"bxs/parser/event_parser/event"
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type CreatedEventParser struct {
	FactoryEventParser
}

func (o *CreatedEventParser) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	pair := &types.Pair{}

	_, ok := o.PossibleFactoryAddresses[ethLog.Address]
	if !ok {
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

	e := &event.CreatedEvent{
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

	pair.Address = e.PoolAddress
	pair.Token0Core = &types.TokenCore{
		Address:  e.TokenAddress,
		Symbol:   e.Symbol,
		Decimals: 18,
	}
	pair.Token1Core = types.NativeTokenCore
	pair.Block = ethLog.BlockNumber
	pair.ProtocolId = abi.FactoryAddress2ProtocolId[ethLog.Address]

	pair.FilterByToken0AndToken1()

	e.Pair = pair

	return e, nil
}
