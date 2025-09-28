package cache

import (
	"bxs/chain_params"
	"bxs/logger"
	"bxs/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"time"
)

type PriceCache interface {
	SetPrice(blockNumber *big.Int, price decimal.Decimal)
	GetPrice(blockNumber *big.Int) (decimal.Decimal, bool)
}

type TokenCache interface {
	SetToken(token *types.Token)
	GetToken(address common.Address) (*types.Token, bool)
	DelToken(address common.Address)
}

type PairCache interface {
	SetPair(pair *types.Pair)
	GetPair(address common.Address) (*types.Pair, bool)
	PairExist(address common.Address) bool
	DelPair(address common.Address)
}

type BlockCache interface {
	SetFinishedBlock(blockNumber uint64)
	GetFinishedBlock() uint64
}

type MigrateTokenCache interface {
	SetMigrateToken(address common.Address)
	MigrateTokenExist(address common.Address) bool
	DelMigrateToken(address common.Address)
}

type Cache interface {
	PriceCache
	TokenCache
	PairCache
	BlockCache
	MigrateTokenCache
}

type twoTierCache struct {
	ctx    context.Context
	memory *cache.Cache
	redis  *redis.Client
}

func NewTwoTierCache(redis *redis.Client) Cache {
	return &twoTierCache{
		ctx:    context.Background(),
		memory: cache.New(time.Hour*24, time.Hour),
		redis:  redis,
	}
}

func PriceCacheKey(blockNumber *big.Int) string {
	return fmt.Sprintf("%d:P:%s", chain_params.G.ChainID, blockNumber.String())
}

func TokenCacheKey(address common.Address) string {
	return fmt.Sprintf("%d:t:%s", chain_params.G.ChainID, address.Hex())
}

func PairCacheKey(address common.Address) string {
	return fmt.Sprintf("%d:p:%s", chain_params.G.ChainID, address.Hex())
}

func MigrateTokenCacheKey(address common.Address) string {
	return fmt.Sprintf("%d:m:%s", chain_params.G.ChainID, address.Hex())
}

func FbKey() string {
	return fmt.Sprintf("%d:fb", chain_params.G.ChainID)
}

func (c *twoTierCache) SetPrice(blockNumber *big.Int, price decimal.Decimal) {
	k := PriceCacheKey(blockNumber)
	c.memory.Set(k, price, cache.DefaultExpiration)
	err := c.redis.Set(c.ctx, k, price.String(), 0).Err()
	if err != nil {
		logger.G.Error("save price failed", zap.Error(err))
	}
}

func (c *twoTierCache) GetPrice(blockNumber *big.Int) (decimal.Decimal, bool) {
	k := PriceCacheKey(blockNumber)
	price, ok := c.memory.Get(k)
	if ok {
		return price.(decimal.Decimal), true
	}

	v, err := c.redis.Get(c.ctx, k).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logger.G.Error("get price failed", zap.Error(err))
		}
		return decimal.Zero, false
	}

	decimalPrice, err := decimal.NewFromString(v)
	if err != nil {
		return decimal.Decimal{}, false
	}
	c.memory.Set(k, decimalPrice, 0)
	return decimalPrice, true
}

func (c *twoTierCache) SetToken(token *types.Token) {
	token.UpdateTs = time.Now()
	k := TokenCacheKey(token.Address)
	c.memory.Set(k, token, cache.DefaultExpiration)

	bytes, err := json.Marshal(token)
	if err != nil {
		logger.G.Sugar().Fatalf("json marshal token [%s] err [%s]", token.Address.String(), err.Error())
	}

	err = c.redis.Set(c.ctx, k, bytes, 0).Err()
	if err != nil {
		logger.G.Sugar().Fatalf("redis set token [%s] err [%s]", string(bytes), err.Error())
	}
}

func (c *twoTierCache) GetToken(address common.Address) (*types.Token, bool) {
	k := TokenCacheKey(address)
	obj, ok := c.memory.Get(k)
	if ok {
		return obj.(*types.Token), true
	}

	bytes, err := c.redis.Get(c.ctx, k).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logger.G.Sugar().Errorf("redis get token [%s] err [%s]", k, err.Error())
		}
		return nil, false
	}

	token := &types.Token{}
	err = json.Unmarshal(bytes, token)
	if err != nil {
		logger.G.Sugar().Errorf("json unmarshal token [%s] err [%s]", string(bytes), err.Error())
		return nil, false
	}

	c.memory.Set(k, token, cache.DefaultExpiration)
	return token, true
}

func (c *twoTierCache) DelToken(address common.Address) {
	k := TokenCacheKey(address)
	c.memory.Delete(k)
	err := c.redis.Del(c.ctx, k).Err()
	if err != nil {
		logger.G.Sugar().Errorf("redis del token [%s] err [%s]", k, err.Error())
	}
}

func (c *twoTierCache) SetPair(pair *types.Pair) {
	pair.UpdateTs = time.Now()
	k := PairCacheKey(pair.Address)
	c.memory.Set(k, pair, cache.DefaultExpiration)

	bytes, err := json.Marshal(pair)
	if err != nil {
		logger.G.Sugar().Fatalf("json marshal pair [%s] err [%s]", pair.Address.String(), err.Error())
	}

	err = c.redis.Set(c.ctx, k, bytes, 0).Err()
	if err != nil {
		logger.G.Sugar().Fatalf("redis set pair [%s] err [%s]", string(bytes), err.Error())
	}
}

func (c *twoTierCache) GetPair(address common.Address) (*types.Pair, bool) {
	k := PairCacheKey(address)
	obj, ok := c.memory.Get(k)
	if ok {
		return obj.(*types.Pair), true
	}

	bytes, err := c.redis.Get(c.ctx, k).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logger.G.Sugar().Errorf("redis get k [%s] err [%s]", k, err.Error())
		}
		return nil, false
	}

	pair := &types.Pair{}
	err = json.Unmarshal(bytes, pair)
	if err != nil {
		logger.G.Sugar().Warnf("json unmarshal bytes [%s] err [%v]", string(bytes), err.Error())
		return nil, false
	}

	c.memory.Set(k, pair, cache.DefaultExpiration)
	return pair, true
}

func (c *twoTierCache) PairExist(address common.Address) bool {
	_, ok := c.GetPair(address)
	return ok
}

func (c *twoTierCache) DelPair(address common.Address) {
	k := PairCacheKey(address)
	c.memory.Delete(k)
	err := c.redis.Del(c.ctx, k).Err()
	if err != nil {
		logger.G.Sugar().Errorf("redis del [%s] err [%s]", k, err.Error())
	}
}

func (c *twoTierCache) SetFinishedBlock(blockNumber uint64) {
	c.redis.Set(c.ctx, FbKey(), blockNumber, 0)
}

func (c *twoTierCache) GetFinishedBlock() uint64 {
	v, err := c.redis.Get(c.ctx, FbKey()).Uint64()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logger.G.Error("redis get err", zap.Error(err))
		}
		return 0
	}
	return v
}

func (c *twoTierCache) SetMigrateToken(address common.Address) {
	k := MigrateTokenCacheKey(address)
	c.memory.Set(k, true, cache.DefaultExpiration)
	err := c.redis.Set(c.ctx, k, true, 0).Err()
	if err != nil {
		logger.G.Sugar().Errorf("redis set migrate token [%s] err [%s]", k, err.Error())
	}
}

func (c *twoTierCache) MigrateTokenExist(address common.Address) bool {
	k := MigrateTokenCacheKey(address)
	_, ok := c.memory.Get(k)
	if ok {
		return true
	}

	_, err := c.redis.Get(c.ctx, k).Bool()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logger.G.Sugar().Errorf("redis get migrate token [%s] err [%s]", k, err.Error())
		}
		return false
	}

	c.memory.Set(k, true, 0)
	return true
}

func (c *twoTierCache) DelMigrateToken(address common.Address) {
	k := MigrateTokenCacheKey(address)
	c.memory.Delete(k)
	err := c.redis.Del(c.ctx, k).Err()
	if err != nil {
		logger.G.Sugar().Fatalf("redis del migrate token [%s] err[%s]", k, err.Error())
	}
}
