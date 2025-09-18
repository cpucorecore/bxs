package service

import (
	"bxs/chain"
	"bxs/config"
	"bxs/log"
	"bxs/metrics"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math"
	"strconv"
	"time"
)

const (
	priceSourceCoingecko = "coingecko"
	priceSourceBitget    = "bitget"
)

var (
	priceKey = fmt.Sprintf("%d:P", chain.ID)
)

var (
	errPriceNotExist = fmt.Errorf("no price found in the specified range, try add toleranceSec")
)

type PriceService interface {
	// Start starts the price service to periodically fetch and save prices.
	Start(ctx context.Context)
	StartApiServer(port int)
	// GetClosestPriceByTimestamp retrieves the native token price by timestamp with a time tolerance.
	// Parameters:
	//   - timestampSec: The timestamp for which the price is requested.
	//   - toleranceSec: The tolerance in seconds to find the closest price. [timestampSec-toleranceSec, timestampSec+toleranceSec]
	GetClosestPriceByTimestamp(timestampSec, toleranceSec int64) (decimal.Decimal, error)
	// GetLatestPrice retrieves the latest native token price and timestamp.
	GetLatestPrice() (price decimal.Decimal, timestampSec int64, err error)
}

type priceService struct {
	ctx              context.Context
	priceSource      string
	getPriceInterval time.Duration
	priceGetter      PriceGetter
	redis            *redis.Client
}

func getPriceGetter(priceSource string) PriceGetter {
	if priceSource == priceSourceCoingecko {
		return &PriceGetterCoingecko{}
	} else {
		return NewPriceGetterBitget()
	}
}

func NewPriceService(priceSource string, redis *redis.Client) PriceService {
	if priceSource != priceSourceCoingecko && priceSource != priceSourceBitget {
		log.Logger.Fatal(
			fmt.Sprintf("invalid price source, must %s or %s", priceSourceCoingecko, priceSourceBitget),
			zap.String("source", priceSource))
	}

	return &priceService{
		ctx:              context.Background(),
		priceSource:      priceSource,
		getPriceInterval: time.Second * time.Duration(config.G.PriceService.GetPriceIntervalSec),
		priceGetter:      getPriceGetter(priceSource),
		redis:            redis,
	}
}

func (ps *priceService) getAndSavePrice() (decimal.Decimal, int64, error) {
	now := time.Now()
	price, timestamp, err := ps.priceGetter.GetLatest()
	if err != nil {
		metrics.GetPriceResult.With(prometheus.Labels{"result": "fail"}).Inc()
		return decimal.Zero, 0, err
	}
	metrics.GetPriceDurationMs.Observe(float64(time.Since(now).Milliseconds()))
	metrics.Price.Set(price.InexactFloat64())
	metrics.GetPriceResult.With(prometheus.Labels{"result": "success"}).Inc()

	if err := ps.savePriceToRedis(price, timestamp); err != nil {
		return decimal.Zero, 0, err
	}

	return price, timestamp, nil
}

func (ps *priceService) Start(ctx context.Context) {
	ticker := time.NewTicker(ps.getPriceInterval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				price, timestamp, err := ps.getAndSavePrice()
				if err != nil {
					log.Logger.Error("getAndSavePrice error", zap.Error(err))
				} else {
					log.Logger.Info("price fetched and saved", zap.String("source", ps.priceSource), zap.String("price", price.String()), zap.Int64("timestamp", timestamp))
				}
			case <-ctx.Done():
				log.Logger.Info("PriceService stopped")
				return
			}
		}
	}()
}

func (ps *priceService) getPriceFromRedis(timestampSec, toleranceSec int64) (decimal.Decimal, error) {
	minTime := timestampSec - toleranceSec
	maxTime := timestampSec + toleranceSec

	pricesWithScores, err := ps.redis.ZRangeByScoreWithScores(ps.ctx, priceKey, &redis.ZRangeBy{
		Min: strconv.FormatInt(minTime, 10),
		Max: strconv.FormatInt(maxTime, 10),
	}).Result()

	if err != nil {
		return decimal.Zero, err
	}

	if len(pricesWithScores) == 0 {
		return decimal.Zero, errPriceNotExist
	}

	var closestPrice decimal.Decimal
	var minDiff int64 = math.MaxInt64

	for _, priceWithScore := range pricesWithScores {
		priceTime := int64(priceWithScore.Score)
		diff := abs(priceTime - timestampSec)

		if diff < minDiff {
			minDiff = diff
			price, err := decimal.NewFromString(priceWithScore.Member.(string))
			if err != nil {
				return decimal.Zero, err
			}
			closestPrice = price
		}
	}

	return closestPrice, nil
}

func (ps *priceService) savePriceToRedis(price decimal.Decimal, timestamp int64) error {
	err := ps.redis.ZAdd(ps.ctx, priceKey, &redis.Z{
		Score:  float64(timestamp),
		Member: price.String(),
	}).Err()
	if err != nil {
		log.Logger.Error("savePriceToRedis error", zap.Error(err))
		return err
	}
	return nil
}

func (ps *priceService) GetClosestPriceByTimestamp(timestampSec, toleranceSec int64) (decimal.Decimal, error) {
	return ps.getPriceFromRedis(timestampSec, toleranceSec)
}

func (ps *priceService) GetLatestPrice() (price decimal.Decimal, timestampSec int64, err error) {
	return ps.getAndSavePrice()
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
