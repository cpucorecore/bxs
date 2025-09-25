package event_parser

import (
	pancakev2 "bxs/abi/pancake/v2"
	pcommon "bxs/parser/common"
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	protocolId   = types.ProtocolIdPancakeV2
	protocolName = types.ProtocolNamePancakeV2
)

var (
	pairCreatedEventParser = &PairCreatedEventParser{
		TopicUnpacker: pcommon.TopicUnpacker{
			Topic: pancakev2.PairCreatedTopic0,
			Unpacker: pcommon.EthLogUnpacker{
				AbiEvent:      pancakev2.PairCreatedEvent,
				TopicLen:      3,
				DataUnpackLen: 2,
			},
		},
	}

	swapEventParser = &SwapEventParser{
		TopicUnpacker: pcommon.TopicUnpacker{
			Topic: pancakev2.SwapTopic0,
			Unpacker: pcommon.EthLogUnpacker{
				AbiEvent:      pancakev2.SwapEvent,
				TopicLen:      3,
				DataUnpackLen: 4,
			},
		},
	}

	syncEventParser = &SyncEventParser{
		TopicUnpacker: pcommon.TopicUnpacker{
			Topic: pancakev2.SyncTopic0,
			Unpacker: pcommon.EthLogUnpacker{
				AbiEvent:      pancakev2.SyncEvent,
				TopicLen:      1,
				DataUnpackLen: 2,
			},
		},
	}

	topic2EventParser = map[common.Hash]pcommon.EventParser{
		pancakev2.PairCreatedTopic0: pairCreatedEventParser,
		pancakev2.SwapTopic0:        swapEventParser,
		pancakev2.SyncTopic0:        syncEventParser,
	}
)

func Reg(registrable pcommon.Registrable) {
	for k, v := range topic2EventParser {
		registrable.Register(k, v)
	}
}
