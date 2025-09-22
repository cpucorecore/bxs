package block_getter

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func Test_GetBlock(t *testing.T) {
	ethCli, err := ethclient.Dial("https://bsc-dataseed.bnbchain.org")
	require.NoError(t, err)
	block, getBlockErr := ethCli.BlockByNumber(context.Background(), big.NewInt(62032287))
	require.NoError(t, getBlockErr)
	t.Log(block)
}

func Test_GetBlockReceipt(t *testing.T) {
	ethCli, err := ethclient.Dial("https://bsc-dataseed.bnbchain.org")
	require.NoError(t, err)
	blockReceipts, getErr := ethCli.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(62032287))
	require.NoError(t, getErr)
	t.Log(blockReceipts)
}
