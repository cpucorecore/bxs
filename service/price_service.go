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
	"sync"
	"time"
)

type PriceService interface {
	Start(startBlockNumber uint64)
	GetPrice(blockNumber *big.Int) (decimal.Decimal, error)
}

type priceService struct {
	fromChain      bool
	cache          cache.Cache
	contractCaller *ContractCaller
	workPoolSize   int
	workPool       *ants.Pool
	ethClient      *ethclient.Client
	priceGetter    *PriceGetterBitget
	price          decimal.Decimal
	timestampSec   int64
	lock           sync.RWMutex
}

func NewPriceService(
	fromChain bool,
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

	ps := &priceService{
		fromChain:      fromChain,
		cache:          cache,
		contractCaller: contractCaller,
		workPoolSize:   poolSize,
		workPool:       workPool,
		ethClient:      ethClient,
		priceGetter:    NewPriceGetterBitget(),
	}

	if !fromChain {
		p, ts, err := ps.priceGetter.GetLatest()
		if err != nil {
			log.Logger.Fatal("get latest price", zap.Error(err))
		}
		ps.updatePrice(p, ts)
	}

	return ps
}

func (ps *priceService) updatePrice(price decimal.Decimal, timestamp int64) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.price = price
	ps.timestampSec = timestamp
}

func (ps *priceService) Start(startBlockNumber uint64) {
	if !ps.fromChain {
		go func() {
			for {
				p, t, err := ps.priceGetter.GetLatest()
				if err != nil {
					log.Logger.Error("get latest price err", zap.Error(err))
					time.Sleep(time.Second)
					continue
				}
				ps.lock.Lock()
				ps.price = p
				ps.timestampSec = t
				ps.lock.Unlock()
			}
		}()
		return
	}

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
	if !ps.fromChain {
		ps.lock.RLock()
		defer ps.lock.RUnlock()
		if time.Now().Unix()-ps.timestampSec > 600 {
			log.Logger.Sugar().Fatal("price get slow 10min")
		}
		return ps.price, nil
	}

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
