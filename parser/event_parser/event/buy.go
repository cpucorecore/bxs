package event

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

func (e *BuyEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Event:         types.Buy,
		Maker:         e.Buyer.String(),
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
