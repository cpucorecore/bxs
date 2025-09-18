package service

import "github.com/shopspring/decimal"

type PriceGetter interface {
	GetLatest() (price decimal.Decimal, timestampSec int64, err error)
}
