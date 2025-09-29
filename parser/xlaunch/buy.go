package event_parser

import (
	"bxs/repository/orm"
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
)

type BuyEvent struct {
	*types.EventCommon
	Buyer             common.Address
	NativeTokenAmount *big.Int
	TokenAmount       *big.Int
	NativeTokenRaised *big.Int
	TokensSold        *big.Int
	Fee               *big.Int
	Migrated          bool
}

func (e *BuyEvent) CanGetTx() bool {
	return true
}

func (e *BuyEvent) GetTx(nativeTokenPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Event:         types.Buy,
		Maker:         e.Buyer.String(),
		Token0Address: e.Pair.Token0.Address.String(),
		Token1Address: e.Pair.Token1.Address.String(),
		Block:         e.BlockNumber,
		BlockAt:       e.BlockTime,
		BlockIndex:    e.TxIndex,
		TxIndex:       e.LogIndex,
		PairAddress:   e.Pair.Address.String(),
		Program:       protocolName,
	}

	tx.Token0Amount, tx.Token1Amount = types.ParseAmount(e.TokenAmount, e.NativeTokenAmount, e.Pair)
	tx.AmountUsd, tx.PriceUsd = types.CalcAmountAndPrice(nativeTokenPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1.Address)
	return tx
}

func (e *BuyEvent) CanGetPoolUpdate() bool {
	return true
}

func (e *BuyEvent) GetPoolUpdate() *types.PoolUpdate {
	u := &types.PoolUpdate{
		LogIndex: e.EventCommon.LogIndex,
		Address:  e.EventCommon.Pair.Address.String(),
		Token0:   e.EventCommon.Pair.Token0.Address.String(),
		Token1:   e.EventCommon.Pair.Token1.Address.String(),
	}
	u.Amount0, u.Amount1 = types.ParseAmount(e.TokensSold, e.NativeTokenRaised, e.Pair)
	return u
}

func (e *BuyEvent) IsMigrated() bool {
	return e.Migrated
}

func (e *BuyEvent) GetAction() *orm.Action {
	action := &orm.Action{
		Maker:   e.Buyer.String(),
		Token:   e.Pair.Token0.Address.String(),
		Action:  "on-uniswap",
		TxHash:  e.TxHash.String(),
		Creator: e.Buyer.String(),
		Block:   e.BlockNumber,
		BlockAt: e.BlockTime,
	}
	action.Token0Amount, action.Token1Amount = types.ParseAmount(e.TokenAmount, e.NativeTokenAmount, e.Pair)
	return action
}

func (e *BuyEvent) IsBuyOrSell() bool {
	return true
}
