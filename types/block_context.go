package types

import (
	"bxs/logger"
	"bxs/repository/orm"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

type HeightTime struct {
	Height       uint64
	Timestamp    uint64
	HeightBigInt *big.Int
	Time         time.Time
}

func GetBlockHeightTime(header *ethtypes.Header) *HeightTime {
	return &HeightTime{
		HeightBigInt: header.Number,
		Height:       header.Number.Uint64(),
		Timestamp:    header.Time,
		Time:         time.Unix((int64)(header.Time), 0).UTC(),
	}
}

type BlockContext struct {
	HeightTime       *HeightTime
	TransactionsLen  uint
	Transactions     []*ethtypes.Transaction
	Receipts         []*ethtypes.Receipt
	NativeTokenPrice decimal.Decimal
	Senders          []common.Address
	TxResults        []*TxResult
}

func (c *BlockContext) getSender(transactionIndex uint) common.Address {
	if transactionIndex >= c.TransactionsLen {
		logger.G.Sugar().Fatalf("get tx sender out of range, transactionIndex: %d, TransactionsLen: %d", transactionIndex, c.TransactionsLen)
	}
	return c.Senders[transactionIndex]
}

func (c *BlockContext) NewTxResult(transactionIndex uint) *TxResult {
	tr := &TxResult{
		Sender:            c.getSender(transactionIndex),
		SwapEvents:        make([]Event, 0, 16),
		PairCreatedEvents: make([]Event, 0, 8),
		PoolUpdates:       make([]*PoolUpdate, 0, 16),
		Pairs:             make([]*Pair, 0, 8),
		Tokens:            make([]*Token, 0, 8),
		MigratedPools:     make([]*MigratedPool, 0, 4),
		Actions:           make([]*orm.Action, 0, 4),
	}
	c.TxResults[transactionIndex] = tr
	return tr
}

func (c *BlockContext) DecorateEvent(event Event) {
	event.SetBlockTime(c.HeightTime.Time)
}

func (c *BlockContext) SetTxResult(transactionIndex uint, tr *TxResult) {
	if transactionIndex >= c.TransactionsLen {
		logger.G.Sugar().Fatalf("SetTxResult out of range, transactionIndex: %d, TransactionsLen: %d", transactionIndex, c.TransactionsLen)
	}
	c.TxResults[transactionIndex] = tr
}

func mergePoolUpdates(poolUpdates []*PoolUpdate) []*PoolUpdate {
	pairAddress2PoolUpdate := make(map[string]*PoolUpdate)
	for _, poolUpdate := range poolUpdates {
		poolUpdate_, ok := pairAddress2PoolUpdate[poolUpdate.Address]
		if ok {
			if poolUpdate.LogIndex > poolUpdate_.LogIndex {
				pairAddress2PoolUpdate[poolUpdate.Address] = poolUpdate
			}
		} else {
			pairAddress2PoolUpdate[poolUpdate.Address] = poolUpdate
		}
	}
	poolUpdatesMerged := make([]*PoolUpdate, 0, len(pairAddress2PoolUpdate))
	for _, pu := range pairAddress2PoolUpdate {
		poolUpdatesMerged = append(poolUpdatesMerged, pu)
	}
	return poolUpdatesMerged
}

func (c *BlockContext) GetKafkaMsg() *KafkaMsg {
	poolUpdates := make([]*PoolUpdate, 0, len(c.TxResults))
	txs := make([]*orm.Tx, 0, len(c.TxResults))
	migratedPools := make([]*MigratedPool, 0, len(c.TxResults))
	actions := make([]*orm.Action, 0, len(c.TxResults))
	ormPairs := make([]*orm.Pair, 0, len(c.TxResults))
	ormTokens := make([]*orm.Token, 0, len(c.TxResults))

	for _, txResult := range c.TxResults {
		poolUpdates = append(poolUpdates, txResult.PoolUpdates...)
		for _, event := range txResult.SwapEvents {
			if event.CanGetTx() {
				txs = append(txs, event.GetTx(c.NativeTokenPrice))
			}
		}

		for _, pair := range txResult.Pairs {
			ormPairs = append(ormPairs, pair.GetOrmPair())
		}

		for _, token := range txResult.Tokens {
			ormTokens = append(ormTokens, token.GetOrmToken())
		}

		migratedPools = append(migratedPools, txResult.MigratedPools...)
		actions = append(actions, txResult.Actions...)
	}

	block := &KafkaMsg{
		Height:           c.HeightTime.Height,
		Timestamp:        c.HeightTime.Timestamp,
		NativeTokenPrice: c.NativeTokenPrice.String(),
		Txs:              txs,
		MigratedPools:    migratedPools,
		Actions:          actions,
		NewTokens:        ormTokens,
		NewPairs:         ormPairs,
		PoolUpdates:      mergePoolUpdates(poolUpdates),
	}

	return block
}

func (c *BlockContext) GetSequence() uint64 {
	return c.HeightTime.Height
}
