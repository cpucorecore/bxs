package service

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

/*
curl -X 'GET' 'https://api.coingecko.com/api/v3/simple/price?ids=newton-project&vs_currencies=usd' -H 'accept: application/json' | jq

	{
	  "newton-project": {
		"usd": 0.01562833
	  }
	}
*/
type coingeckoResp struct {
	NewtonProject struct {
		USD float64 `json:"usd"`
	} `json:"newton-project"`
}

type PriceGetterCoingecko struct{}

func (pg *PriceGetterCoingecko) GetLatest() (price decimal.Decimal, timestampSec int64, err error) {
	timestamp := time.Now().Unix()
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=newton-project&vs_currencies=usd")
	if err != nil {
		return decimal.Zero, 0, err
	}
	defer resp.Body.Close()

	var result coingeckoResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return decimal.Zero, 0, err
	}

	return decimal.NewFromFloat(result.NewtonProject.USD), timestamp, nil
}
