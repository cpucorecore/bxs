package service

import (
	uniswapv2 "bxs/abi/uniswap/v2"
	uniswapv3 "bxs/abi/uniswap/v3"
	"bxs/abi/xlaunch"
	"bxs/cache"
	"bxs/log"
	"bxs/metrics"
	"bxs/types"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"math/big"
	"sync"
	"time"
)

type PairService interface {
	SetPair(pair *types.Pair)
	GetPairTokens(pair *types.Pair) *types.PairWrap
	GetPair(pairAddress common.Address, possibleProtocolIds []int) *types.PairWrap
}

type pairService struct {
	ctx            context.Context
	cache          cache.Cache
	contractCaller *ContractCaller
	group          singleflight.Group
}

func NewPairService(
	cache cache.Cache,
	contractCaller *ContractCaller,
) PairService {
	return &pairService{
		ctx:            context.Background(),
		cache:          cache,
		contractCaller: contractCaller,
	}
}

func (s *pairService) SetPair(pair *types.Pair) {
	s.cache.SetPair(pair)
}

func (s *pairService) doGetToken(tokenAddress common.Address) (*types.Token, error) {
	token := &types.Token{
		Address: tokenAddress,
	}

	var (
		nameRes struct {
			name string
			err  error
		}
		symbolRes struct {
			symbol string
			err    error
		}
		decimalsRes struct {
			decimals int
			err      error
		}
		supplyRes struct {
			supply *big.Int
			err    error
		}
	)

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		nameRes.name, nameRes.err = s.contractCaller.CallName(&tokenAddress)
	}()

	go func() {
		defer wg.Done()
		symbolRes.symbol, symbolRes.err = s.contractCaller.CallSymbol(&tokenAddress)
	}()

	go func() {
		defer wg.Done()
		decimalsRes.decimals, decimalsRes.err = s.contractCaller.CallDecimals(&tokenAddress)
	}()

	go func() {
		defer wg.Done()
		supplyRes.supply, supplyRes.err = s.contractCaller.CallTotalSupply(&tokenAddress)
	}()
	wg.Wait()

	if nameRes.err == nil {
		token.Name = nameRes.name
	}

	if symbolRes.err == nil {
		token.Symbol = symbolRes.symbol
	}

	if decimalsRes.err != nil {
		token.Filtered = true
		return token, decimalsRes.err
	}
	token.Decimals = int8(decimalsRes.decimals)

	if supplyRes.err == nil {
		token.TotalSupply = decimal.NewFromBigInt(supplyRes.supply, -int32(token.Decimals))
	}

	return token, nil
}

func (s *pairService) getToken(tokenAddress common.Address) (*types.Token, error, bool) {
	cacheToken, ok := s.cache.GetToken(tokenAddress)
	if ok {
		return cacheToken, nil, true
	}

	now := time.Now()
	doResult, err, _ := s.group.Do(tokenAddress.String(), func() (interface{}, error) {
		token, err := s.doGetToken(tokenAddress)
		s.cache.SetToken(token)
		return token, err
	})
	if err != nil {
		return nil, err, false
	}

	metrics.GetTokenDurationMs.Observe(float64(time.Since(now).Milliseconds()))
	return doResult.(*types.Token), nil, false
}

var (
	nativeToken = &types.Token{
		Address:  types.ZeroAddress,
		Creator:  types.ZeroAddress,
		Symbol:   "BNB",
		Decimals: 18,
	}
)

func (s *pairService) getNativeToken() (*types.Token, error) {
	return nativeToken, nil
}

func (s *pairService) getPairTokens(pair *types.Pair) *types.PairWrap {
	pairWrap := &types.PairWrap{
		Pair:    pair,
		NewPair: !s.cache.PairExist(pair.Address),
	}

	token0, err, fromCache := s.getToken(pair.Token0Core.Address)
	if err != nil {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken
		return pairWrap
	}

	pair.Token0 = token0
	pair.Token1 = nativeToken

	pair.Token0Core.Symbol = token0.Symbol
	pair.Token0Core.Decimals = token0.Decimals
	pair.Token1Core.Symbol = nativeToken.Symbol
	pair.Token1Core.Decimals = nativeToken.Decimals

	pairWrap.NewToken0 = !fromCache
	pairWrap.NewToken1 = false

	return pairWrap
}

func (s *pairService) GetPairTokens(pair *types.Pair) *types.PairWrap {
	doResult, _, _ := s.group.Do(pair.Address.String(), func() (interface{}, error) {
		pairWrap := s.getPairTokens(pair)
		s.SetPair(pair)
		return pairWrap, nil
	})

	return doResult.(*types.PairWrap)
}

func (s *pairService) getPair(pairAddress common.Address, possibleProtocolIds []int) *types.PairWrap {
	doResult, _, _ := s.group.Do(pairAddress.String()+"gp", func() (interface{}, error) {
		pair := s.doGetPair(pairAddress)
		if pair.Filtered {
			s.SetPair(pair)
			return &types.PairWrap{
				Pair:      pair,
				NewPair:   false,
				NewToken0: false,
				NewToken1: false,
			}, nil
		}

		if !s.verifyPair(pair, possibleProtocolIds) {
			s.SetPair(pair)
			return &types.PairWrap{
				Pair:      pair,
				NewPair:   false,
				NewToken0: false,
				NewToken1: false,
			}, nil
		}

		return s.GetPairTokens(pair), nil
	})

	return doResult.(*types.PairWrap)
}

func (s *pairService) GetPair(pairAddress common.Address, possibleProtocolIds []int) *types.PairWrap {
	cachePair, ok := s.cache.GetPair(pairAddress)
	if ok {
		return &types.PairWrap{
			Pair: cachePair,
		}
	}

	return s.getPair(pairAddress, possibleProtocolIds)
}

func (s *pairService) doGetPair(pairAddress common.Address) *types.Pair {
	pair := &types.Pair{
		Address: pairAddress,
	}

	now := time.Now()
	token0Addr, err := s.contractCaller.CallToken(&pairAddress)
	if err != nil {
		log.Logger.Info("Err: CallToken err, this pair will filtered",
			zap.Error(err),
			zap.String("pair address", pairAddress.String()),
		)
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken
		return pair
	}

	pair.Token0Core = &types.TokenCore{
		Address: token0Addr,
	}
	pair.Token1Core = &types.TokenCore{
		Address: types.ZeroAddress,
	}

	metrics.GetPairDurationMs.Observe(float64(time.Since(now).Milliseconds()))
	pair.FilterByToken0AndToken1()
	return pair
}

func (s *pairService) verifyPairV2(pairFactoryAddress common.Address, pair *types.Pair) bool {
	pairAddressQueried, err := s.contractCaller.CallGetPair(&pairFactoryAddress, &pair.Token0Core.Address, &pair.Token1Core.Address)
	if err != nil {
		return false
	}
	return types.IsSameAddress(pairAddressQueried, pair.Address)
}

func (s *pairService) verifyXLaunch(pairFactoryAddress common.Address, pair *types.Pair) bool {
	verified, err := s.contractCaller.CallGetLaunchByAddress(&pairFactoryAddress, &pair.Address)
	if err != nil {
		return false
	}
	return verified
}

func (s *pairService) verifyPairV3(pairFactoryAddress common.Address, pair *types.Pair) bool {
	fee, callFeeErr := s.contractCaller.CallFee(&pair.Address)
	if callFeeErr != nil {
		return false
	}

	pairAddressQueried, err := s.contractCaller.CallGetPool(&pairFactoryAddress, &pair.Token0Core.Address, &pair.Token1Core.Address, fee)
	if err != nil {
		return false
	}

	return types.IsSameAddress(pairAddressQueried, pair.Address)
}

func (s *pairService) verifyPair(pair *types.Pair, possibleProtocolIds []int) bool {
	now := time.Now()
	defer func() {
		duration := float64(time.Since(now).Milliseconds())
		metrics.VerifyPairDurationMs.Observe(duration)
	}()

	for _, protocolId := range possibleProtocolIds {
		switch protocolId {
		case types.ProtocolIdNewSwap:
			if s.verifyPairV2(uniswapv2.FactoryAddress, pair) {
				pair.ProtocolId = protocolId
				metrics.VerifyPairTotal.WithLabelValues("success").Inc()
				metrics.VerifyPairOkByProtocol.WithLabelValues("uniswap_v2").Inc()
				return true
			}

		case types.ProtocolIdUniswapV3:
			if s.verifyPairV3(uniswapv3.FactoryAddress, pair) {
				pair.ProtocolId = protocolId
				metrics.VerifyPairTotal.WithLabelValues("success").Inc()
				metrics.VerifyPairOkByProtocol.WithLabelValues("uniswap_v3").Inc()
				return true
			}

		case types.ProtocolIdXLaunch:
			if s.verifyXLaunch(xlaunch.FactoryAddress, pair) {
				pair.ProtocolId = protocolId
				metrics.VerifyPairTotal.WithLabelValues("success").Inc()
				metrics.VerifyPairOkByProtocol.WithLabelValues("xlaunch").Inc()
				return true
			}
		}
	}

	pair.Filtered = true
	pair.FilterCode = types.FilterCodeVerifyFailed
	metrics.VerifyPairTotal.WithLabelValues("failed").Inc()

	return false
}
