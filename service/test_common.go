package service

import (
	"github.com/ethereum/go-ethereum/common"
)

type TestToken struct {
	address  common.Address
	name     string
	symbol   string
	decimals int8
}
