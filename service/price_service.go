package service

import (
	"bxs/cache"
	"bxs/log"
	"bxs/metrics"
	"bxs/types"
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"time"
)

type PriceService interface {
	Start(startBlockNumber uint64)
	GetPrice(blockNumber *big.Int) (decimal.Decimal, error)
}

type priceService struct {
	cache          cache.Cache
	contractCaller *ContractCaller
	workPoolSize   int
	workPool       *ants.Pool
	ethClient      *ethclient.Client
}

func NewPriceService(
	cache cache.Cache,
	contractCaller *ContractCaller,
	ethClient *ethclient.Client,
	poolSize int,
) PriceService {
	var workPool *ants.Pool
	var err error
	if poolSize > 0 {
		workPool, err = ants.NewPool(poolSize)
		if err != nil {
			log.Logger.Fatal("ants pool(BlockGetter) init err", zap.Error(err))
		}
	}

	return &priceService{
		cache:          cache,
		contractCaller: contractCaller,
		workPoolSize:   poolSize,
		workPool:       workPool,
		ethClient:      ethClient,
	}
}

func (ps *priceService) Start(startBlockNumber uint64) {
	if ps.workPoolSize <= 0 {
		return
	}

	go func() {
		for {
			headerBlockNumber, err := ps.ethClient.BlockNumber(context.Background())
			if err != nil {
				log.Logger.Error("ethClient.BigIntHeight", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}

			for startBlockNumber <= headerBlockNumber {
				ps.workPool.Submit(func() {
					ps.GetPrice(big.NewInt(int64(startBlockNumber)))
					startBlockNumber++
				})
			}
		}
	}()
}

func (ps *priceService) GetPrice(blockNumber *big.Int) (decimal.Decimal, error) {
	cachePrice, ok := ps.cache.GetPrice(blockNumber)
	if ok {
		return cachePrice, nil
	}

	return ps.getPrice(blockNumber)
}

func (ps *priceService) getPrice(blockNumber *big.Int) (decimal.Decimal, error) {
	now := time.Now()

	bnbPrice, err := ps.contractCaller.GetPriceByBlockNumber(blockNumber)
	if err != nil {
		log.Logger.Error("GetPriceByBlockNumber err", zap.Error(err), zap.Uint64("blockNumber", blockNumber.Uint64()))
		return types.ZeroDecimal, err
	}

	metrics.CallContractForBNBPrice.Observe(time.Since(now).Seconds())
	ps.cache.SetPrice(blockNumber, bnbPrice)

	return bnbPrice, nil
}
