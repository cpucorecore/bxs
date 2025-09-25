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
		Address:       tp.address,
		TokenReversed: tp.tokenReversed,
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
		Address:       tp.address,
		TokenReversed: tp.tokenReversed,
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

	return pair
}
