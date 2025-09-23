package service

import (
	"bxs/abi/bep20"
	"bxs/abi/ds_token"
	pancakev2 "bxs/abi/pancake/v2"
	"bxs/abi/xlaunch"
	"bxs/types"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
	"unicode/utf8"
)

type Unpacker interface {
	Unpack(method string, data []byte, length int) (values []interface{}, err error)
}

type unpacker struct {
	abis []*abi.ABI
}

func NewUnpacker(abis []*abi.ABI) Unpacker {
	return &unpacker{
		abis: abis,
	}
}

var (
	UnpackErr           = errors.New("unpack error")
	ErrWrongString      = errors.New("wrong string")
	ErrWrongIntType     = errors.New("wrong int type")
	ErrWrongBigIntType  = errors.New("wrong big int type")
	ErrWrongAddressType = errors.New("wrong address type")
	ErrWrongBoolType    = errors.New("wrong bool type")
)

func (u *unpacker) Unpack(method string, data []byte, length int) ([]interface{}, error) {
	for _, abi_ := range u.abis {
		values, err := abi_.Unpack(method, data)
		if err == nil && len(values) == length {
			return values, nil
		}
	}
	return nil, UnpackErr
}

var (
	TokenUnpacker = NewUnpacker([]*abi.ABI{
		bep20.Abi,
		ds_token.Abi,
	})

	PancakeV2PairUnpacker = NewUnpacker([]*abi.ABI{
		pancakev2.PairAbi,
	})

	PancakeV2FactoryUnpacker = NewUnpacker([]*abi.ABI{
		pancakev2.FactoryAbi,
	})

	XLaunchUnpacker = NewUnpacker([]*abi.ABI{
		xlaunch.PairAbi,
	})

	XLaunchFactoryUnpacker = NewUnpacker([]*abi.ABI{
		xlaunch.FactoryAbi,
	})

	Name2Unpacker = map[string]Unpacker{
		"name":        TokenUnpacker,
		"symbol":      TokenUnpacker,
		"decimals":    TokenUnpacker,
		"totalSupply": TokenUnpacker,
		"token0":      PancakeV2PairUnpacker,
		"token1":      PancakeV2PairUnpacker,
		"getReserves": PancakeV2PairUnpacker,
		"token":       XLaunchUnpacker,
	}
)

func sanitizeUTF8(s string) string {
	if !utf8.ValidString(s) {
		return strings.ToValidUTF8(s, "?")
	}
	return s
}

func ParseString(value interface{}) (string, error) {
	var str string
	var err error
	switch v := value.(type) {
	case string:
		str = v
	case [32]byte:
		str = string(v[:])
	default:
		err = ErrWrongString
	}
	if err != nil {
		return "", err
	}

	str = strings.ReplaceAll(str, "\x00", "") // for postgres db do not accept 0x00 as string char
	return sanitizeUTF8(str), nil
}

func ParseInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case uint8:
		return int(v), nil
	case *big.Int:
		return int(v.Int64()), nil
	default:
		return 0, ErrWrongIntType
	}
}

func ParseBigInt(value interface{}) (*big.Int, error) {
	if bigIntValue, ok := value.(*big.Int); ok {
		return bigIntValue, nil
	}
	return nil, ErrWrongBigIntType
}

func ParseAddress(value interface{}) (common.Address, error) {
	if address, ok := value.(common.Address); ok {
		return address, nil
	} else {
		return types.ZeroAddress, ErrWrongAddressType
	}
}

func ParseBool(value interface{}) (bool, error) {
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return false, ErrWrongBoolType
}
