package service

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
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
