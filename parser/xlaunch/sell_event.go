package event_parser

import (
	"bxs/repository/orm"
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
)

type SellEvent struct {
	*types.EventCommon
	Seller            common.Address
	NativeTokenAmount *big.Int
	TokenAmount       *big.Int
	NativeTokenRaised *big.Int
	TokensSold        *big.Int
	Fee               *big.Int
}

func (e *SellEvent) CanGetTx() bool {
	return true
}

func (e *SellEvent) GetTx(nativeTokenPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Event:         types.Sell,
		Maker:         e.Seller.String(),
		Token0Address: e.Pair.Token0.Address.String(),
		Token1Address: e.Pair.Token1.Address.String(),
		Block:         e.BlockNumber,
		BlockAt:       e.BlockTime,
		BlockIndex:    e.TxIndex,
		TxIndex:       e.LogIndex,
		PairAddress:   e.Pair.Address.String(),
		Program:       types.ProtocolNameXLaunch,
	}

	tx.Token0Amount, tx.Token1Amount = types.ParseAmountsByPair(e.TokenAmount, e.NativeTokenAmount, e.Pair)
	tx.AmountUsd, tx.PriceUsd = types.CalcAmountAndPrice(nativeTokenPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1.Address)
	return tx
}

func (e *SellEvent) CanGetPoolUpdate() bool {
	return true
}

func (e *SellEvent) GetPoolUpdate() *types.PoolUpdate {
	a0, a1 := types.ParseAmountsByPair(e.TokensSold, e.NativeTokenRaised, e.Pair)
	return &types.PoolUpdate{
		LogIndex: e.EventCommon.LogIndex,
		Address:  e.EventCommon.Pair.Address,
		Token0:   e.EventCommon.Pair.Token0.Address,
		Token1:   e.EventCommon.Pair.Token1.Address,
		Amount0:  a0,
		Amount1:  a1,
	}
}

func (e *SellEvent) IsBuyOrSell() bool {
	return true
}
