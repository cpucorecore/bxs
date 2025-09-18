package service

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestContractCaller_CallContract(t *testing.T) {
	cc := GetTestContext().ContractCaller
	address := common.HexToAddress("0x4200000000000000000000000000000000000006")
	req := &CallContractReq{
		Address: &address,
	}

	// call erc20 contract with a method not exist, should return non err and empty bytes
	req.Data = Name2Data["getReserves"]
	bytes, err := cc.CallContract(req)
	require.Nil(t, err)
	require.Equal(t, 0, len(bytes))

	// call erc20 contract with a method exist, should return non err and non-empty bytes
	req.Data = Name2Data["name"]
	bytes, err = cc.CallContract(req)
	require.Nil(t, err)
	require.True(t, len(bytes) > 0)
}

func TestContractCaller_queryValues(t *testing.T) {
	cc := GetTestContext().ContractCaller
	pairAddress := common.HexToAddress("0xc9034c3E7F58003E6ae0C8438e7c8f4598d5ACAA")

	// call pair contract with a method not exist, should return err and empty values
	values, err := cc.queryValues(&pairAddress, "name", 1)
	require.Equal(t, ErrOutputEmpty, err)
	require.Equal(t, 0, len(values))

	// call pair contract with a method exist, should return non err and non-empty values
	values, err = cc.queryValues(&pairAddress, "token0", 1)
	require.Nil(t, err)
	require.True(t, len(values) > 0)
}

func TestContractCaller_CallXX(t *testing.T) {
	cc := GetTestContext().ContractCaller

	tests := []struct {
		address        string
		expectName     string
		expectSymbol   string
		expectDecimals int
	}{
		{
			address:        "0x4200000000000000000000000000000000000006",
			expectName:     "Wrapped Ether",
			expectSymbol:   "WETH",
			expectDecimals: 18,
		},
		{
			address:        "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
			expectName:     "USD Coin",
			expectSymbol:   "USDC",
			expectDecimals: 6,
		},
		{
			address:        "0xC2DC84144f625B4feC5a21B888028EBD6c95E38d",
			expectName:     "nightmare.exe",
			expectSymbol:   "nightmare.exe",
			expectDecimals: 18,
		},
		{
			address:        "0xB4bD1fE69dCAA0C64fDb34075a0CC2b332Bd015e",
			expectName:     "PB64",
			expectSymbol:   "PB64",
			expectDecimals: 18,
		},
	}

	for _, test := range tests {
		address := common.HexToAddress(test.address)
		name, callNameErr := cc.CallName(&address)
		require.Nil(t, callNameErr)
		require.Equal(t, test.expectName, name)
		symbol, callSymbolErr := cc.CallSymbol(&address)
		require.Nil(t, callSymbolErr)
		require.Equal(t, test.expectSymbol, symbol)
		decimals, callDecimalsErr := cc.CallDecimals(&address)
		require.Nil(t, callDecimalsErr)
		require.Equal(t, test.expectDecimals, decimals)
		totalSupply, callTotalSupplyErr := cc.CallTotalSupply(&address)
		require.Nil(t, callTotalSupplyErr)
		t.Log(address, totalSupply)
	}
}

func TestContractCaller_CallToken0AndCallToken1(t *testing.T) {
	cc := GetTestContext().ContractCaller

	tests := []struct {
		pairAddress   common.Address
		token0Address common.Address
		token1Address common.Address
	}{
		{
			pairAddress:   pairUniswapV2.address,
			token0Address: pairUniswapV2.token0.address,
			token1Address: pairUniswapV2.token1.address,
		},
		{
			pairAddress:   pairUniswapV3.address,
			token0Address: pairUniswapV3.token0.address,
			token1Address: pairUniswapV3.token1.address,
		},
	}

	for _, test := range tests {
		token0Address, err0 := cc.CallToken0(&test.pairAddress)
		require.Nil(t, err0)
		require.Equal(t, test.token0Address, token0Address)
		token1Address, err1 := cc.CallToken1(&test.pairAddress)
		require.Nil(t, err1)
		require.Equal(t, test.token1Address, token1Address)
	}
}

func TestContractCaller_CallFee(t *testing.T) {
	cc := GetTestContext().ContractCaller

	tests := []struct {
		callErr     bool
		pairAddress common.Address
		expectFee   *big.Int
	}{
		{
			callErr:     true,
			pairAddress: pairUniswapV2.address,
		},
		{
			callErr:     false,
			pairAddress: pairUniswapV3.address,
			expectFee:   pairUniswapV3.fee,
		},
	}

	for _, test := range tests {
		fee, err := cc.CallFee(&test.pairAddress)
		if test.callErr {
			require.NotNil(t, err, test.pairAddress)
		} else {
			require.Nil(t, err, test.pairAddress)
			require.Equal(t, test.expectFee.String(), fee.String(), test.pairAddress)
		}
	}
}
