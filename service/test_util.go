package service

import (
	"bxs/cache"
	"bxs/chain_params"
	"bxs/config"
	"context"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"math/big"
)

var (
	MockNativeTokenPrice = decimal.NewFromInt(1)
	Wei18, _             = decimal.NewFromString("1000000000000000000")
	Wei6, _              = decimal.NewFromString("1000000")
)

type TestContext struct {
	ethClient      *ethclient.Client
	Cache          cache.Cache
	ContractCaller *ContractCaller
	PairService    PairService
}

func GetTestContext() *TestContext {
	ethClient, err := ethclient.Dial("https://bsc-testnet-dataseed.bnbchain.org")
	if err != nil {
		panic(err)
	}

	factoryAddress := common.HexToAddress("0xdb9d8574AB68AA9ccd5CAedE1C641BC42a995a54")
	chain_params.LoadNetwork(true, factoryAddress)

	contractCaller := NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())
	cache := cache.NewMockCache()
	pairService_ := NewPairService(cache, contractCaller)

	return &TestContext{
		ethClient:      ethClient,
		Cache:          cache,
		ContractCaller: contractCaller,
		PairService:    pairService_,
	}
}

func (g *TestContext) GetEthLog(txHashStr string, logIndex int) *ethtypes.Log {
	txHash := common.HexToHash(txHashStr)
	txReceipt, apiErr := g.ethClient.TransactionReceipt(context.Background(), txHash)
	if apiErr != nil {
		panic(apiErr)
	}

	return txReceipt.Logs[logIndex]
}

func (g *TestContext) GetBlockTimestamp(blockNumber uint64) uint64 {
	blockHeader, err := g.ethClient.HeaderByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		panic(err)
	}
	return blockHeader.Time
}
