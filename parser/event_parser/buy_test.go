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
	// https://testnet.bscscan.com/tx/0xd7b86b1409b41cae1ccd839f190939bd24bfd41071e459c3ff317ce23c818fe5#eventlog#3
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0xd7b86b1409b41cae1ccd839f190939bd24bfd41071e459c3ff317ce23c818fe5", 2)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), nil)
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(service.MockNativeTokenPrice)
	expectAmt0, _ := decimal.NewFromString("102795326.423086443244433004")
	expectAmt1, _ := decimal.NewFromString("1.804950495049504950")
	expectTx := &orm.Tx{
		TxHash:        "0xd7b86b1409b41cae1ccd839f190939bd24bfd41071e459c3ff317ce23c818fe5",
		Event:         types.Buy,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: "0xFA4dA14E995408Fd456928F4a0512AC348de1794",
		Token1Address: types.ZeroAddress.String(),
		Block:         65764034,
		BlockIndex:    1,
		TxIndex:       3,
		PairAddress:   "0x87485818145cEC5017a6466AAD2Ef5FEeA99aaae",
		Program:       types.ProtocolNameXLaunch,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
