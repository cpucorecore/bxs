package service

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPairService_GetPair_UniswapV2(t *testing.T) {
	tc := GetTestContext()
	testPair := pairUniswapV2

	pw := tc.PairService.GetPair(testPair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(testPair.GetExpectedPair()), "pair should be equal", pw.Pair, testPair.GetExpectedPair())

	pw = tc.PairService.GetPair(testPair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(testPair.GetExpectedPair()), "pair should be equal", pw.Pair, testPair.GetExpectedPair())
}

func TestPairService_GetPair_UniswapV3(t *testing.T) {
	tc := GetTestContext()
	testPair := pairUniswapV3

	pw := tc.PairService.GetPair(testPair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(testPair.GetExpectedPair()), "pair should be equal", pw.Pair, testPair.GetExpectedPair())

	pw = tc.PairService.GetPair(testPair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(testPair.GetExpectedPair()), "pair should be equal", pw.Pair, testPair.GetExpectedPair())
}

func TestPairService_GetGetPairTokens_UniswapV2(t *testing.T) {
	tc := GetTestContext()
	testPair := pairUniswapV2
	expectPair := testPair.GetExpectedPair()

	pairWithoutTokenInfo := testPair.GetPairWithoutTokenInfo()
	pw := tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)

	pairWithoutTokenInfo = testPair.GetPairWithoutTokenInfo()
	pw = tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)
}

func TestPairService_GetGetPairTokens_UniswapV3(t *testing.T) {
	tc := GetTestContext()
	testPair := pairUniswapV3
	expectPair := testPair.GetExpectedPair()

	pairWithoutTokenInfo := testPair.GetPairWithoutTokenInfo()
	pw := tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)

	pairWithoutTokenInfo = testPair.GetPairWithoutTokenInfo()
	pw = tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)
}
