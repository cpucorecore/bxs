package event_parser

import (
	"bxs/chain_params"
	pcommon "bxs/parser/common"
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
	pcommon.TopicUnpacker
}

func (o *PairCreatedEventParser) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	if !types.IsSameAddress(ethLog.Address, chain_params.G.PancakeV2FactoryAddress) {
		return nil, ErrWrongFactory
	}

	eventInput, err := o.Unpacker.Unpack(ethLog)
	if err != nil {
		return nil, err
	}

	e := &PairCreatedEvent{
		EventCommon: types.EventCommonFromEthLog(ethLog),
		Address:     eventInput[0].(common.Address),
		Token0:      common.BytesToAddress(ethLog.Topics[1].Bytes()[12:]),
		Token1:      common.BytesToAddress(ethLog.Topics[2].Bytes()[12:]),
	}

	if !e.IsWBNBPair() {
		return nil, ErrNotWBNBPair
	}

	e.Token0, e.Token1, e.tokenReversed = types.OrderToken0Token1Address(e.Token0, e.Token1)
	e.Pair = &types.Pair{
		Address: e.Address,
		Token0: &types.TokenTinyInfo{
			Address: e.Token0,
		},
		Token1: &types.TokenTinyInfo{
			Address: e.Token1,
		},
		TokenReversed: e.tokenReversed,
		Block:         e.BlockNumber,
		ProtocolId:    protocolId,
	}

	return e, nil
}
