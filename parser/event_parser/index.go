package event_parser

import (
	pancakev2 "bxs/abi/pancake/v2"
	"bxs/abi/xlaunch"
	"github.com/ethereum/go-ethereum/common"
)

var (
	createdEventParser = &CreatedEventParser{
		TopicUnpacker{
			topic: xlaunch.CreatedTopic0,
			unpacker: EthLogUnpacker{
				AbiEvent:      xlaunch.CreatedEvent,
				TopicLen:      4,
				DataUnpackLen: 6,
			},
		},
	}

	buyEventParser = &BuyEventParser{
		TopicUnpacker{
			topic: xlaunch.BuyTopic0,
			unpacker: EthLogUnpacker{
				AbiEvent:      xlaunch.BuyEvent,
				TopicLen:      2,
				DataUnpackLen: 6,
			},
		},
	}

	sellEventParser = &SellEventParser{
		TopicUnpacker{
			topic: xlaunch.SellTopic0,
			unpacker: EthLogUnpacker{
				AbiEvent:      xlaunch.SellEvent,
				TopicLen:      2,
				DataUnpackLen: 5,
			},
		},
	}

	pairCreatedEventParser = &PairCreatedEventParser{
		TopicUnpacker: TopicUnpacker{
			topic: pancakev2.PairCreatedTopic0,
			unpacker: EthLogUnpacker{
				AbiEvent:      pancakev2.PairCreatedEvent,
				TopicLen:      3,
				DataUnpackLen: 2,
			},
		},
		FactoryAddress: pancakev2.FactoryAddressTestnet,
	}

	Topic2EventParser = map[common.Hash]EventParser{
		xlaunch.CreatedTopic0:       createdEventParser,
		xlaunch.BuyTopic0:           buyEventParser,
		xlaunch.SellTopic0:          sellEventParser,
		pancakev2.PairCreatedTopic0: pairCreatedEventParser,
	}
)
