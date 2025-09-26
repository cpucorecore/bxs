package event_parser

import (
	"bxs/log"
	"bxs/repository/orm"
	"bxs/types"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
)

type SwapEvent struct {
	*types.EventCommon
	Amount0InWei  *big.Int
	Amount1InWei  *big.Int
	Amount0OutWei *big.Int
	Amount1OutWei *big.Int
}

func (e *SwapEvent) CanGetTx() bool {
	return true
}

func (e *SwapEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Maker:         e.Maker.String(),
		Token0Address: e.Pair.Token0.Address.String(),
		Token1Address: e.Pair.Token1.Address.String(),
		Block:         e.BlockNumber,
		BlockAt:       e.BlockTime,
		BlockIndex:    e.TxIndex,
		TxIndex:       e.LogIndex,
		PairAddress:   e.Pair.Address.String(),
		Program:       protocolName,
	}

	if e.Amount0InWei.Cmp(types.ZeroBigInt) > 0 {
		tx.Token0Amount, tx.Token1Amount = types.ParseAmountsByPair(e.Amount0InWei, e.Amount1OutWei, e.Pair)
		if !e.Pair.TokenReversed {
			tx.Event = types.Sell
		} else {
			tx.Event = types.Buy
		}
	} else if e.Amount1InWei.Cmp(types.ZeroBigInt) > 0 {
		tx.Token0Amount, tx.Token1Amount = types.ParseAmountsByPair(e.Amount0OutWei, e.Amount1InWei, e.Pair)
		if !e.Pair.TokenReversed {
			tx.Event = types.Buy
		} else {
			tx.Event = types.Sell
		}
	} else {
		log.Logger.Warn("wrong pancake v2 swap event", zap.Any("event", e))
	}

	tx.AmountUsd, tx.PriceUsd = types.CalcAmountAndPrice(bnbPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1.Address)
	return tx
}

func (e *SwapEvent) IsSwap() bool {
	return true
}
