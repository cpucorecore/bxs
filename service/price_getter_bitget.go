package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

/*
https://www.bitget.com/zh-CN/api-doc/spot/market/Get-Tickers
curl "https://api.bitget.com/api/v2/spot/market/tickers?symbol=ABUSDT" | jq

	{
	  "code": "00000",
	  "msg": "success",
	  "requestTime": 1750142199572,
	  "data": [
	    {
	      "open": "0.015264",
	      "symbol": "ABUSDT",
	      "high24h": "0.01604",
	      "low24h": "0.0103",
	      "lastPr": "0.015647",
	      "quoteVolume": "91356498.3",
	      "baseVolume": "6008761391.8",
	      "usdtVolume": "91356498.291362518172",
	      "ts": "1750142198529",
	      "bidPr": "0.015624",
	      "askPr": "0.015646",
	      "bidSz": "10000",
	      "askSz": "12754.28",
	      "openUtc": "0.015806",
	      "changeUtc24h": "-0.01006",
	      "change24h": "0.02509"
	    }
	  ]
	}
*/

const (
	bitgetAPIUrl = "https://api.bitget.com/api/v2/spot/market/tickers?symbol=BNBUSDT"
)

type bitgetResp struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
	Data        []struct {
		LastPr string `json:"lastPr"`
		Ts     string `json:"ts"`
	} `json:"data"`
}

type PriceGetterBitget struct {
	httpClient *resty.Client

	// History price getter is not implemented yet
	// https://www.bitget.com/zh-CN/api-doc/spot/market/Get-History-Candle-Data
}

func NewPriceGetterBitget() *PriceGetterBitget {
	httpClient := resty.New()

	httpClient.SetTimeout(time.Second * 10)
	httpClient.SetRetryCount(10)
	httpClient.SetRetryWaitTime(time.Millisecond * 100)
	httpClient.SetRetryMaxWaitTime(time.Millisecond * 500)

	return &PriceGetterBitget{
		httpClient: httpClient,
	}
}

func (pg *PriceGetterBitget) GetLatest() (price decimal.Decimal, timestampSec int64, err error) {
	resp, err := pg.httpClient.R().Get(bitgetAPIUrl)
	if err != nil {
		return decimal.Zero, 0, err
	}

	var result bitgetResp
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return decimal.Zero, 0, err
	}

	if len(result.Data) == 0 {
		return decimal.Zero, 0, fmt.Errorf("no data returned from bitget")
	}

	timestamp, err := strconv.ParseInt(result.Data[0].Ts, 10, 64)
	if err != nil {
		return decimal.Zero, 0, err
	}
	timestamp = timestamp / 1000

	price, err = decimal.NewFromString(result.Data[0].LastPr)
	if err != nil {
		return decimal.Zero, 0, err
	}

	return price, timestamp, nil
}
