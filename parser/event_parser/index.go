package event_parser

import (
	"bxs/abi"
	"bxs/abi/xlaunch"
	"github.com/ethereum/go-ethereum/common"
)

var (
	createdEventParser = &CreatedEventParser{
		FactoryEventParser{
			Topic:                    xlaunch.CreatedTopic0,
			PossibleFactoryAddresses: abi.Topic2FactoryAddresses[xlaunch.CreatedTopic0],
			LogUnpacker: EthLogUnpacker{
				AbiEvent:      xlaunch.CreatedEvent,
				TopicLen:      4,
				DataUnpackLen: 6,
			},
		},
	}

	buyEventParser = &BuyEventParser{
		PoolEventParser{
			Topic:               xlaunch.BuyTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[xlaunch.BuyTopic0],
			ethLogUnpacker: EthLogUnpacker{
				AbiEvent:      xlaunch.BuyEvent,
				TopicLen:      2,
				DataUnpackLen: 6,
			},
		},
	}

	sellEventParser = &SellEventParser{
		PoolEventParser{
			Topic:               xlaunch.SellTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[xlaunch.SellTopic0],
			ethLogUnpacker: EthLogUnpacker{
				AbiEvent:      xlaunch.SellEvent,
				TopicLen:      2,
				DataUnpackLen: 5,
			},
		},
	}

	Topic2EventParser = map[common.Hash]EventParser{
		xlaunch.CreatedTopic0: createdEventParser,
		xlaunch.BuyTopic0:     buyEventParser,
		xlaunch.SellTopic0:    sellEventParser,
	}
)
