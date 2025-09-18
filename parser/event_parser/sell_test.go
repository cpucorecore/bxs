package event_parser

import (
	"bxs/repository/orm"
	"bxs/service"
	"bxs/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSell(t *testing.T) {
	// https://testnet.bscscan.com/tx/0x7cb0894568573d4bd590f185fa166fb73f64bbb827b362c0017de6473ad2849e#eventlog#2
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0x7cb0894568573d4bd590f185fa166fb73f64bbb827b362c0017de6473ad2849e", 2)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(service.MockNativeTokenPrice)
	expectAmt0, _ := decimal.NewFromString("102795326.423086443244433004")
	expectAmt1, _ := decimal.NewFromString("1.804950495049504950")
	expectTx := &orm.Tx{
		TxHash:        "0x7cb0894568573d4bd590f185fa166fb73f64bbb827b362c0017de6473ad2849e",
		Event:         types.Sell,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: "0xFA4dA14E995408Fd456928F4a0512AC348de1794",
		Token1Address: types.ZeroAddress.String(),
		Block:         65764330,
		BlockIndex:    1,
		TxIndex:       3,
		PairAddress:   "0x87485818145cEC5017a6466AAD2Ef5FEeA99aaae",
		Program:       types.ProtocolNameXLaunch,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
