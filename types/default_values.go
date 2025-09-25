package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
)

var (
	ZeroAddress       = common.Address{}
	ZeroDecimal       = decimal.NewFromInt(0)
	ZeroBigInt        = new(big.Int)
	Decimals18        = int8(18)
	NativeTokenSymbol = "BNB"
)
