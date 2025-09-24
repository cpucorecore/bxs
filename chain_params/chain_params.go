package chain_params

import (
	"bxs/chain"
	"bxs/chain/v1_5_17"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

type ChainParams struct {
	ChainID                      int
	ChainConfig                  *params.ChainConfig
	WBNBAddress                  common.Address
	PancakeV2FactoryAddress      common.Address
	PancakeV2BusdWbnbPairAddress common.Address
	XLaunchFactoryAddress        common.Address
}

const (
	WBNBAddressHex                    = "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"
	WBNBAddressTestnetHex             = "0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd"
	PancakeV2FactoryAddressHex        = "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"
	PancakeV2FactoryAddressTestnetHex = "0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"
	PancakeV2BusdWbnbPairHex          = "0x58F876857a02D6762E0101bb5C46A8c1ED44Dc16"
	PancakeV2BusdWbnbPairTestnetHex   = "0x85EcDcdd01EbE0BfD0Aba74B81Ca6d7F4A53582b"
)

var (
	WBNBAddress                     = common.HexToAddress(WBNBAddressHex)
	WBNBAddressTestnet              = common.HexToAddress(WBNBAddressTestnetHex)
	PancakeV2FactoryAddress         = common.HexToAddress(PancakeV2FactoryAddressHex)
	PancakeV2FactoryAddressTestnet  = common.HexToAddress(PancakeV2FactoryAddressTestnetHex)
	PancakeV2BusdWbnbAddress        = common.HexToAddress(PancakeV2BusdWbnbPairHex)
	PancakeV2BusdWbnbAddressTestnet = common.HexToAddress(PancakeV2BusdWbnbPairTestnetHex)

	mainnetParams = &ChainParams{
		ChainID:                      chain.BSCMainnetID,
		ChainConfig:                  v1_5_17.BSCChainConfig,
		PancakeV2FactoryAddress:      PancakeV2FactoryAddress,
		PancakeV2BusdWbnbPairAddress: PancakeV2BusdWbnbAddress,
		WBNBAddress:                  WBNBAddress,
	}

	testnetParams = &ChainParams{
		ChainID:                      chain.BSCTestnetID,
		ChainConfig:                  v1_5_17.ChapelChainConfig,
		PancakeV2FactoryAddress:      PancakeV2FactoryAddressTestnet,
		PancakeV2BusdWbnbPairAddress: PancakeV2BusdWbnbAddressTestnet,
		WBNBAddress:                  WBNBAddressTestnet,
	}

	G *ChainParams
)

func LoadNetwork(testnet bool, xLaunchFactoryAddress common.Address) {
	if testnet {
		G = testnetParams
	} else {
		G = mainnetParams
	}
	G.XLaunchFactoryAddress = xLaunchFactoryAddress
}
