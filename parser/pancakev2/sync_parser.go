package event_parser

import (
	pcommon "bxs/parser/common"
	"bxs/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type SyncEventParser struct {
	pcommon.TopicUnpacker
}

func (o *SyncEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	input, err := o.Unpacker.Unpack(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &SyncEvent{
		EventCommon: types.EventCommonFromEthLog(receiptLog),
		amount0Wei:  input[0].(*big.Int),
		amount1Wei:  input[1].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address:    receiptLog.Address,
		ProtocolId: protocolId,
	}

	return e, nil
}
