package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
)

func ParseAmountsByPair(token0AmountWei, token1AmountWei *big.Int, pair *Pair) (token0Amount, token1Amount decimal.Decimal) {
	if !pair.TokenReversed {
		token0Amount = decimal.NewFromBigInt(token0AmountWei, -(int32)(pair.Token0.Decimal))
		token1Amount = decimal.NewFromBigInt(token1AmountWei, -(int32)(pair.Token1.Decimal))
	} else {
		token0Amount = decimal.NewFromBigInt(token1AmountWei, -(int32)(pair.Token0.Decimal))
		token1Amount = decimal.NewFromBigInt(token0AmountWei, -(int32)(pair.Token1.Decimal))
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
