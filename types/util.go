package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
)

func ParseAmount(a0Wei, a1Wei *big.Int, pair *Pair) (a0, a1 decimal.Decimal) {
	if !pair.TokenReversed {
		a0 = decimal.NewFromBigInt(a0Wei, -(int32)(pair.Token0.Decimal))
		a1 = decimal.NewFromBigInt(a1Wei, -(int32)(pair.Token1.Decimal))
	} else {
		a0 = decimal.NewFromBigInt(a1Wei, -(int32)(pair.Token0.Decimal))
		a1 = decimal.NewFromBigInt(a0Wei, -(int32)(pair.Token1.Decimal))
	}
	return
}

func CalcAmountAndPrice(
	nativeTokenPrice decimal.Decimal,
	amount0, amount1 decimal.Decimal,
	token1 common.Address,
) (amountUSD, priceUSD decimal.Decimal) {
	if IsNativeToken(token1) || IsWBNB(token1) {
		amountUSD = amount1.Mul(nativeTokenPrice)
		if !amount0.IsZero() {
			priceUSD = amountUSD.Div(amount0)
		}
	}
	return
}
