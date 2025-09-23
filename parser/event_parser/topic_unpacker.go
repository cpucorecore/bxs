package event_parser

import "github.com/ethereum/go-ethereum/common"

type TopicUnpacker struct {
	topic       common.Hash
	unpacker    EthLogUnpacker
	factoryAddr common.Address
}
