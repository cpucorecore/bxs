package types

import (
	"bxs/chain_params"
	"bxs/logger"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"time"
)

var (
	txIndexOutOfRange = errors.New("txIndex out of range")
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

type BlockCtx struct {
	// input
	HeightTime       *HeightTime
	NativeTokenPrice decimal.Decimal
	TransactionsLen  uint
	Transactions     []*ethtypes.Transaction
	Receipts         []*ethtypes.Receipt
	TxSenders        []*common.Address

	// output
	BlockResult *BlockResult
}

func (c *BlockCtx) GetSequence() uint64 {
	return c.HeightTime.Height
}

func (c *BlockCtx) GetTxSender(txIndex uint) (common.Address, error) {
	if c.TxSenders[txIndex] != nil {
		return *c.TxSenders[txIndex], nil
	}

	if txIndex >= c.TransactionsLen {
		logger.G.Info("Waring: txIndex out of range",
			zap.Uint64("height", c.HeightTime.Height),
			zap.Any("transactions length", c.TransactionsLen),
			zap.Uint("txIndex", txIndex),
		)
		return ZeroAddress, txIndexOutOfRange
	}

	signer := ethtypes.MakeSigner(chain_params.G.ChainConfig, c.HeightTime.HeightBigInt, c.HeightTime.Timestamp)
	sender, err := ethtypes.Sender(signer, c.Transactions[txIndex])
	if err != nil {
		return ZeroAddress, err
	}

	c.TxSenders[txIndex] = &sender
	return sender, nil
}
