package event_parser

import (
	"bxs/parser/event_parser/event"
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type SellEventParser struct {
	PoolEventParser
}

func (o *SellEventParser) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	eventInput, err := o.ethLogUnpacker.Unpack(ethLog)
	if err != nil {
		return nil, err
	}

	e := &event.SellEvent{
		EventCommon:       types.EventCommonFromEthLog(ethLog),
		Seller:            common.BytesToAddress(ethLog.Topics[1].Bytes()[12:]),
		NativeTokenAmount: eventInput[0].(*big.Int),
		TokenAmount:       eventInput[1].(*big.Int),
		NativeTokenRaised: eventInput[2].(*big.Int),
		TokensSold:        eventInput[3].(*big.Int),
		Fee:               eventInput[4].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address: ethLog.Address,
	}

	e.PossibleProtocolIds = o.PossibleProtocolIds

	return e, nil
}
