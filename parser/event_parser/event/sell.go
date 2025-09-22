package event

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
		Token0Address: e.Pair.Token0Core.Address.String(),
		Token1Address: e.Pair.Token1Core.Address.String(),
		Block:         e.BlockNumber,
		BlockAt:       e.BlockTime,
		BlockIndex:    e.TxIndex,
		TxIndex:       e.LogIndex,
		PairAddress:   e.Pair.Address.String(),
		Program:       types.ProtocolNameXLaunch,
	}

	tx.Token0Amount, tx.Token1Amount = ParseAmountsByPair(e.TokenAmount, e.NativeTokenAmount, e.Pair)
	tx.AmountUsd, tx.PriceUsd = CalcAmountAndPrice(nativeTokenPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1Core.Address)
	return tx
}

func (e *SellEvent) CanGetPoolUpdate() bool {
	return true
}

func (e *SellEvent) GetPoolUpdate() *types.PoolUpdate {
	a0, a1 := ParseAmountsByPair(e.TokensSold, e.NativeTokenRaised, e.Pair)
	return &types.PoolUpdate{
		Program:       types.ProtocolNameXLaunch,
		LogIndex:      e.EventCommon.LogIndex,
		Address:       e.EventCommon.Pair.Address,
		Token0Address: e.EventCommon.Pair.Token0Core.Address,
		Token1Address: e.EventCommon.Pair.Token1Core.Address,
		Token0Amount:  a0,
		Token1Amount:  a1,
	}
}
