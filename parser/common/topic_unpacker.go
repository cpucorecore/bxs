package event_parser

import "github.com/ethereum/go-ethereum/common"

type TopicUnpacker struct {
	Topic    common.Hash
	Unpacker EthLogUnpacker
}
