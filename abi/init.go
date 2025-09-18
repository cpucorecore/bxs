package abi

import (
	"bxs/abi/xlaunch"
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
)

var Topic2ProtocolIds = map[common.Hash][]int{}
var FactoryAddress2ProtocolId = map[common.Address]int{}
var Topic2FactoryAddresses = map[common.Hash]map[common.Address]struct{}{}

func mapTopicToProtocolId(topic common.Hash, protocolId int) {
	protocolIds, ok := Topic2ProtocolIds[topic]
	if !ok {
		protocolIds = []int{}
	}
	protocolIds = append(protocolIds, protocolId)
	Topic2ProtocolIds[topic] = protocolIds
}

func mapTopicToFactoryAddress(topic common.Hash, factoryAddress common.Address) {
	factoryAddresses, ok := Topic2FactoryAddresses[topic]
	if !ok {
		factoryAddresses = make(map[common.Address]struct{})
	}
	factoryAddresses[factoryAddress] = struct{}{}
	Topic2FactoryAddresses[topic] = factoryAddresses
}

func init() {
	FactoryAddress2ProtocolId[xlaunch.FactoryAddress] = types.ProtocolIdXLaunch
	mapTopicToFactoryAddress(xlaunch.CreatedTopic0, xlaunch.FactoryAddress)
	mapTopicToProtocolId(xlaunch.BuyTopic0, types.ProtocolIdXLaunch)
}
