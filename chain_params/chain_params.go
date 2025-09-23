package chain_params

import (
	pancakev2 "bxs/abi/pancake/v2"
	"bxs/chain"
	"bxs/chain/v1_5_17"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

type ChainParams struct {
	ChainID                 int
	ChainConfig             *params.ChainConfig
	PancakeV2FactoryAddress common.Address
	WBNBAddress             common.Address
}

const (
	WBNBAddressHex        = "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"
	WBNBAddressTestnetHex = "0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd"
)

var (
	WBNBAddress        = common.HexToAddress(WBNBAddressHex)
	WBNBAddressTestnet = common.HexToAddress(WBNBAddressTestnetHex)

	mainnetParams = ChainParams{
		ChainID:                 chain.BSCMainnetID,
		ChainConfig:             v1_5_17.BSCChainConfig,
		PancakeV2FactoryAddress: pancakev2.FactoryAddress,
		WBNBAddress:             WBNBAddress,
	}

	testnetParams = ChainParams{
		ChainID:                 chain.BSCTestnetID,
		ChainConfig:             v1_5_17.ChapelChainConfig,
		PancakeV2FactoryAddress: pancakev2.FactoryAddressTestnet,
		WBNBAddress:             WBNBAddressTestnet,
	}

	G ChainParams = mainnetParams
)

func LoadNetwork(testnet bool) {
	if testnet {
		G = testnetParams
	}
}
