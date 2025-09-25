package event_parser

import (
	pcommon "bxs/parser/common"
	"bxs/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type SwapEventParser struct {
	pcommon.TopicUnpacker
}

func (o *SwapEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	eventInput, err := o.Unpacker.Unpack(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &SwapEvent{
		EventCommon:   types.EventCommonFromEthLog(receiptLog),
		Amount0InWei:  eventInput[0].(*big.Int),
		Amount1InWei:  eventInput[1].(*big.Int),
		Amount0OutWei: eventInput[2].(*big.Int),
		Amount1OutWei: eventInput[3].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address:    receiptLog.Address,
		ProtocolId: protocolId,
	}

	return e, nil
}
