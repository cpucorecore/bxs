package event_parser

import (
	"bxs/repository/orm"
	"bxs/service"
	"bxs/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuy(t *testing.T) {
	// https://testnet.bscscan.com/tx/0xb93f156a59a1f9c92a0af06f430fa942a08392c46f126de104c24fd9d8fb75c9#eventlog#2
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0xb93f156a59a1f9c92a0af06f430fa942a08392c46f126de104c24fd9d8fb75c9", 2)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(service.MockNativeTokenPrice)
	expectAmt0, _ := decimal.NewFromString("44711.015496156927491764")
	expectAmt1, _ := decimal.NewFromString("0.004573267326732715")
	expectTx := &orm.Tx{
		TxHash:        "0xb93f156a59a1f9c92a0af06f430fa942a08392c46f126de104c24fd9d8fb75c9",
		Event:         types.Buy,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: "0xFA4dA14E995408Fd456928F4a0512AC348de1794",
		Token1Address: types.ZeroAddress.String(),
		Block:         65764034,
		BlockIndex:    1,
		TxIndex:       3,
		PairAddress:   "0xCe32c1326450C7AC8D9698E65d3303efB4F211c0",
		Program:       types.ProtocolNameXLaunch,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
