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
	// https://testnet.bscscan.com/tx/0x6144c88e93d9e8f055cd0abff16e1be204482096ffe9f64e1e30fa2154003ed9#eventlog#13
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0x6144c88e93d9e8f055cd0abff16e1be204482096ffe9f64e1e30fa2154003ed9", 5)
	blockTimestamp := tc.GetBlockTimestamp(ethLog.BlockNumber)

	event, pErr := topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	pair := event.GetPair()
	pairWrap := tc.PairService.GetPairTokens(pair)
	event.SetPair(pairWrap.Pair)

	token0InitAmount, err := decimal.NewFromString("1066666666.666666666666666666")
	require.NoError(t, err)
	token1InitAmount, err := decimal.NewFromString("6.933333333333333333")
	require.NoError(t, err)
	expectPair := &types.Pair{
		Address:       common.HexToAddress("0x87485818145cEC5017a6466AAD2Ef5FEeA99aaae"),
		TokenReversed: false,
		Token0Core: &types.TokenCore{
			Address:  common.HexToAddress("0xFA4dA14E995408Fd456928F4a0512AC348de1794"),
			Symbol:   "T",
			Decimals: 18,
		},
		Token1Core:  types.NativeTokenTinyInfo,
		InitAmount0: token0InitAmount,
		InitAmount1: token1InitAmount,
		Block:       65762817,
		BlockAt:     time.Unix(int64(blockTimestamp), 0),
		ProtocolId:  protocolId,
		Filtered:    false,
		FilterCode:  0,
	}

	require.True(t, pairWrap.Pair.Equal(expectPair), "expect: %v, actual: %v", expectPair, pairWrap.Pair)
}
