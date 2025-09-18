package service

import (
	"bxs/config"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewContractCaller(t *testing.T) {
	t.Skip()
	err := config.LoadConfigFile("../config.json")
	require.Nil(t, err, err)

	redisCli := redis.NewClient(&redis.Options{
		Addr:     config.G.Redis.Addr,
		Username: config.G.Redis.Username,
		Password: config.G.Redis.Password,
	})
	ps := NewPriceService("bitget", redisCli)
	ps.Start(t.Context())
	// let the service run for a while to fetch prices
	time.Sleep(3 * time.Second)

	price, err := ps.GetClosestPriceByTimestamp(time.Now().Unix(), 600)
	require.Nil(t, err, err)
	t.Log(price)
}

func TestPriceServiceApi(t *testing.T) {
	t.Skip()
	err := config.LoadConfigFile("../config.json")
	require.Nil(t, err, err)

	redisCli := redis.NewClient(&redis.Options{
		Addr:     config.G.Redis.Addr,
		Username: config.G.Redis.Username,
		Password: config.G.Redis.Password,
	})
	ps := NewPriceService("bitget", redisCli)
	ps.Start(t.Context())
	ps.StartApiServer(config.G.PriceService.Port)
	// let the service run for a while to fetch prices
	time.Sleep(300 * time.Second)
}
