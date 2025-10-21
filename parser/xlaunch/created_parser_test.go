package event_parser

import (
	"bxs/service"
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreated(t *testing.T) {
	// https://testnet.bscscan.com/tx/0xff34b651bc1cf2b5cdd57cdb72dbe84ca953d4a1cd833f83441f2ac834d7cffc#eventlog#5
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0xff34b651bc1cf2b5cdd57cdb72dbe84ca953d4a1cd833f83441f2ac834d7cffc", 5)
	blockTimestamp := tc.GetBlockTimestamp(ethLog.BlockNumber)

	event, pErr := topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	pair := event.GetPair()

	token0InitAmount, err := decimal.NewFromString("10666666.666666666666666666")
	require.NoError(t, err)
	token1InitAmount, err := decimal.NewFromString("0.069333333333333333")
	require.NoError(t, err)
	expectBlockTime := time.Unix(int64(blockTimestamp), 0)
	expectPair := &types.Pair{
		Address:       common.HexToAddress("0x9182A7b564C43dbc2EE58Da9B270Fe13D1dd976e"),
		TokenReversed: false,
		Token0: &types.TokenTinyInfo{
			Address: common.HexToAddress("0xDA519FB1b564A0CE2e10E48CEDbe5BFEd490623D"),
			Symbol:  "zz",
			Decimal: 18,
		},
		Token1:      types.NativeTokenTinyInfo,
		InitAmount0: token0InitAmount,
		InitAmount1: token1InitAmount,
		Block:       69460293,
		BlockAt:     expectBlockTime,
		ProtocolId:  protocolId,
		Filtered:    false,
		FilterCode:  0,
	}

	require.True(t, pair.Equal(expectPair), "expect: %v, actual: %v", expectPair, pair)

	token0 := event.GetToken0()
	expectToken0TotalSupply, _ := decimal.NewFromString("10000000")
	expectToken0 := &types.Token{
		Address:     common.HexToAddress("0xDA519FB1b564A0CE2e10E48CEDbe5BFEd490623D"),
		Creator:     common.HexToAddress("0x866925e79c447352711bF740183AA3Cc67371E16"),
		Name:        "223",
		Symbol:      "zz",
		Decimals:    18,
		TotalSupply: expectToken0TotalSupply,
		BlockNumber: 69460293,
		BlockTime:   expectBlockTime,
		Program:     protocolName,
		Filtered:    false,
		Cid:         "bafkreibut3ftcrldii42ffjpbm3tvlvgvtw7nhdribgpgkrm5xpr5z67qm",
		Tid:         "10392",
		Description: "ss",
		Telegram:    "",
		Twitter:     "",
		Website:     "",
	}

	require.True(t, token0.Equal(expectToken0), "expect: %v, actual: %v", expectToken0, token0)
}
