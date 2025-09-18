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

func (e *SellEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
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
		Program:       types.GetProtocolName(e.Pair.ProtocolId),
	}

	tx.Token0Amount, tx.Token1Amount = ParseAmountsByPair(e.NativeTokenAmount, e.TokenAmount, e.Pair)
	tx.AmountUsd, tx.PriceUsd = CalcAmountAndPrice(bnbPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1Core.Address)
	return tx
}
