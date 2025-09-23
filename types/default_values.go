package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

var (
	ZeroAddress       = common.Address{}
	ZeroDecimal       = decimal.NewFromInt(0)
	Decimals18        = int8(18)
	NativeTokenSymbol = "BNB"
)
