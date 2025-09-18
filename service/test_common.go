package service

import (
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type TestToken struct {
	address  common.Address
	name     string
	symbol   string
	decimals int8
}

type TestPair struct {
	protocolId    int
	address       common.Address
	tokenReversed bool
	token0        *TestToken
	token1        *TestToken
	fee           *big.Int
}

func (tp *TestPair) GetPairWithoutTokenInfo() *types.Pair {
	pair := &types.Pair{
		Address:        tp.address,
		TokensReversed: tp.tokenReversed,
		Token0Core: &types.TokenCore{
			Address: tp.token0.address,
		},
		Token1Core: &types.TokenCore{
			Address: tp.token1.address,
		},
		ProtocolId: tp.protocolId,
	}
	return pair
}

func (tp *TestPair) GetExpectedPair() *types.Pair {
	pair := &types.Pair{
		Address:        tp.address,
		TokensReversed: tp.tokenReversed,
		Token0Core: &types.TokenCore{
			Address:  tp.token0.address,
			Symbol:   tp.token0.symbol,
			Decimals: tp.token0.decimals,
		},
		Token1Core: &types.TokenCore{
			Address:  tp.token1.address,
			Symbol:   tp.token1.symbol,
			Decimals: tp.token1.decimals,
		},
		ProtocolId: tp.protocolId,
	}

	pair.OrderToken0Token1()
	return pair
}

var (
	tokenWETH = &TestToken{
		address:  types.WETHAddress,
		name:     "Wrapped Ether",
		symbol:   "WETH",
		decimals: 18,
	}
	tokenAERO = &TestToken{
		address:  common.HexToAddress("0x940181a94A35A4569E4529A3CDfB74e38FD98631"),
		name:     "Aerodrome",
		symbol:   "AERO",
		decimals: 18,
	}
	tokenUSDC = &TestToken{
		address:  common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"),
		name:     "USD Coin",
		symbol:   "USDC",
		decimals: 6,
	}
	tokenPEPE = &TestToken{
		address:  common.HexToAddress("0x52b492a33E447Cdb854c7FC19F1e57E8BfA1777D"),
		name:     "BasedPepe",
		symbol:   "PEPE",
		decimals: 18,
	}
	tokenCAKE = &TestToken{
		address:  common.HexToAddress("0x3055913c90Fcc1A6CE9a358911721eEb942013A1"),
		name:     "PancakeSwap Token",
		symbol:   "Cake",
		decimals: 18,
	}
	tokenDEGEN = &TestToken{
		address:  common.HexToAddress("0x4ed4E862860beD51a9570b96d89aF5E1B0Efefed"),
		name:     "Degen",
		symbol:   "DEGEN",
		decimals: 18,
	}
)

var (
	pairUniswapV2 = &TestPair{
		protocolId: types.ProtocolIdNewSwap,
		address:    common.HexToAddress("0x88A43bbDF9D098eEC7bCEda4e2494615dfD9bB9C"),
		token0:     tokenWETH,
		token1:     tokenUSDC,
	}
	pairUniswapV3 = &TestPair{
		protocolId:    types.ProtocolIdUniswapV3,
		address:       common.HexToAddress("0x0FB597D6cFE5bE0d5258A7f017599C2A4Ece34c7"),
		tokenReversed: true,
		token0:        tokenWETH,
		token1:        tokenPEPE,
		fee:           big.NewInt(10000),
	}
)

var (
	possibleProtocolIds = []int{types.ProtocolIdNewSwap, types.ProtocolIdUniswapV3}
)
