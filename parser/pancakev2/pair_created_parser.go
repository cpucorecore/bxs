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
		Token0Addr:  common.BytesToAddress(ethLog.Topics[1].Bytes()[12:]),
		Token1Addr:  common.BytesToAddress(ethLog.Topics[2].Bytes()[12:]),
		PairAddr:    eventInput[0].(common.Address),
	}

	if !e.IsWBNBPair() {
		return nil, ErrNotWBNBPair
	}

	e.Token0Addr, e.Token1Addr, e.tokenReversed = types.OrderToken0Token1Address(e.Token0Addr, e.Token1Addr)
	e.Pair = &types.Pair{
		Address:       e.PairAddr,
		TokenReversed: e.tokenReversed,
		Token0Core: &types.TokenCore{
			Address: e.Token0Addr,
		},
		Token1Core: &types.TokenCore{
			Address: e.Token1Addr,
		},
		Block:      e.BlockNumber,
		BlockAt:    e.BlockTime,
		ProtocolId: protocolId,
	}

	return e, nil
}
