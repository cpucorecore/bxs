package event_parser

import (
	"bxs/abi/xlaunch"
	pcommon "bxs/parser/common"
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	protocolId   = types.ProtocolIdXLaunch
	protocolName = types.ProtocolNameXLaunch
)

var (
	xLaunchTokenDecimal = int8(18)
)

var (
	createdEventParser = &CreatedEventParser{
		pcommon.TopicUnpacker{
			Topic: xlaunch.CreatedTopic0,
			Unpacker: pcommon.EthLogUnpacker{
				AbiEvent:      xlaunch.CreatedEvent,
				TopicLen:      4,
				DataUnpackLen: 10,
			},
		},
	}

	buyEventParser = &BuyEventParser{
		pcommon.TopicUnpacker{
			Topic: xlaunch.BuyTopic0,
			Unpacker: pcommon.EthLogUnpacker{
				AbiEvent:      xlaunch.BuyEvent,
				TopicLen:      2,
				DataUnpackLen: 6,
			},
		},
	}

	sellEventParser = &SellEventParser{
		pcommon.TopicUnpacker{
			Topic: xlaunch.SellTopic0,
			Unpacker: pcommon.EthLogUnpacker{
				AbiEvent:      xlaunch.SellEvent,
				TopicLen:      2,
				DataUnpackLen: 5,
			},
		},
	}

	topic2EventParser = map[common.Hash]pcommon.EventParser{
		xlaunch.CreatedTopic0: createdEventParser,
		xlaunch.BuyTopic0:     buyEventParser,
		xlaunch.SellTopic0:    sellEventParser,
	}
)

func Reg(registrable pcommon.Registrable) {
	for k, v := range topic2EventParser {
		registrable.Register(k, v)
	}
}
