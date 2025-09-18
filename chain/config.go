package chain

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/params"
)

var (
	Config *params.ChainConfig
)

const (
	configJson = `
{
    "chainId": 36888,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "muirGlacierBlock": 0,
    "berlinBlock": 0,
    "clique": {
      "period": 3,
      "epoch": 30000
    }
}`
)

func init() {
	var chainConfig params.ChainConfig
	err := json.Unmarshal([]byte(configJson), &chainConfig)
	if err != nil {
		panic(err)
	}
	Config = &chainConfig
}
