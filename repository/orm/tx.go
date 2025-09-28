package orm

import (
	"bxs/logger"
	"bxs/util"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Tx struct {
	Id            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;readonly"`
	TxHash        string
	Event         string
	Token0Amount  decimal.Decimal
	Token1Amount  decimal.Decimal
	Maker         string
	Token0Address string
	Token1Address string
	AmountUsd     decimal.Decimal
	PriceUsd      decimal.Decimal
	Block         uint64
	BlockAt       time.Time
	BlockIndex    uint
	TxIndex       uint
	PairAddress   string
	Program       string
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (t *Tx) Equal(tx *Tx) bool {
	if t.TxHash != tx.TxHash {
		return false
	}
	if t.Event != tx.Event {
		return false
	}
	//if t.Maker != tx.Maker {
	//	return false
	//}
	if t.Token0Address != tx.Token0Address {
		return false
	}
	if t.Token1Address != tx.Token1Address {
		return false
	}
	if t.Block != tx.Block {
		return false
	}
	if t.BlockIndex != tx.BlockIndex {
		return false
	}
	if t.TxIndex != tx.TxIndex {
		return false
	}
	if t.PairAddress != tx.PairAddress {
		return false
	}
	if t.Program != tx.Program {
		return false
	}
	if !util.DecimalEqual(tx.Token0Amount, t.Token0Amount) {
		return false
	}
	if !util.DecimalEqual(tx.Token1Amount, t.Token1Amount) {
		return false
	}

	// ignore priceUsd and amountUSd
	return true
}

func (t *Tx) Diff(tx *Tx) {
	if t.TxHash != tx.TxHash {
		logger.G.Sugar().Infof("diff hash[%s/%s]", t.TxHash, tx.TxHash)
	}
	if t.Event != tx.Event {
		logger.G.Sugar().Infof("diff event[%s/%s]", t.Event, tx.Event)
	}
	if t.Token0Address != tx.Token0Address {
		logger.G.Sugar().Infof("diff token0[%s/%s]", t.Token0Address, tx.Token0Address)
	}
	if t.Token1Address != tx.Token1Address {
		logger.G.Sugar().Infof("diff token1[%s/%s]", t.Token1Address, tx.Token1Address)
	}
	if t.Block != tx.Block {
		logger.G.Sugar().Infof("diff block[%d/%d]", t.Block, tx.Block)
	}
	if t.BlockIndex != tx.BlockIndex {
		logger.G.Sugar().Infof("diff block index[%d/%d]", t.BlockIndex, tx.BlockIndex)
	}
	if t.TxIndex != tx.TxIndex {
		logger.G.Sugar().Infof("diff tx index[%d/%d]", t.TxIndex, tx.TxIndex)
	}
	if t.PairAddress != tx.PairAddress {
		logger.G.Sugar().Infof("diff pair[%s/%s]", t.PairAddress, tx.PairAddress)
	}
	if t.Program != tx.Program {
		logger.G.Sugar().Infof("diff program[%s/%s]", t.Program, tx.Program)
	}
	if !util.DecimalEqual(tx.Token0Amount, t.Token0Amount) {
		logger.G.Sugar().Infof("diff amt0[%s/%s]", t.Token0Amount, tx.Token1Amount)
	}
	if !util.DecimalEqual(tx.Token1Amount, t.Token1Amount) {
		logger.G.Sugar().Infof("diff amt1[%s/%s]", t.Token1Amount, tx.Token0Amount)
	}
}

func (t *Tx) TableName() string {
	return "tx"
}
