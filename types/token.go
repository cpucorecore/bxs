package types

import (
	"bxs/chain_params"
	"bxs/repository/orm"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"time"
)

var (
	NativeToken = &Token{
		Address:  ZeroAddress,
		Creator:  ZeroAddress,
		Symbol:   NativeTokenSymbol,
		Decimals: Decimals18,
	}

	NativeTokenCore = &TokenCore{
		Address:  ZeroAddress,
		Symbol:   NativeTokenSymbol,
		Decimals: Decimals18,
	}
)

func IsSameAddress(address1, address2 common.Address) bool {
	return address1.Cmp(address2) == 0
}

func IsWBNB(address common.Address) bool {
	return IsSameAddress(address, chain_params.G.WBNBAddress)
}

func IsNativeToken(address common.Address) bool {
	return IsSameAddress(address, ZeroAddress)
}

func IsBaseToken(address common.Address) bool {
	if IsNativeToken(address) {
		return true
	}
	return false
}

type Token struct {
	Address     common.Address `json:"-"`
	Creator     common.Address `json:"-"`
	Name        string
	Symbol      string
	Decimals    int8
	TotalSupply decimal.Decimal
	BlockNumber uint64
	BlockTime   time.Time
	Program     string
	Filtered    bool
	Timestamp   time.Time
	URL         string
	Description string
}

func (t *Token) MarshalBinary() ([]byte, error) {
	type Alias Token
	return json.Marshal(&struct {
		AddressString string `json:"Address"`
		CreatorString string `json:"Creator"`
		*Alias
	}{
		AddressString: t.Address.String(),
		CreatorString: t.Creator.String(),
		Alias:         (*Alias)(t),
	})
}

func (t *Token) UnmarshalBinary(data []byte) error {
	type Alias Token
	aux := &struct {
		AddressString string `json:"Address"`
		CreatorString string `json:"Creator"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Address = common.HexToAddress(aux.AddressString)
	t.Creator = common.HexToAddress(aux.CreatorString)
	return nil
}

func (t *Token) Equal(token *Token) bool {
	if !IsSameAddress(t.Address, token.Address) {
		return false
	}
	if !IsSameAddress(t.Creator, token.Creator) {
		return false
	}
	if t.Name != token.Name {
		return false
	}
	if t.Symbol != token.Symbol {
		return false
	}
	if t.Decimals != token.Decimals {
		return false
	}
	if t.BlockNumber != token.BlockNumber {
		return false
	}
	if t.Program != token.Program {
		return false
	}

	return true
}

func (t *Token) GetOrmToken() *orm.Token {
	ormToken := &orm.Token{
		Address:     t.Address.String(),
		Creator:     t.Creator.String(),
		Name:        t.Name,
		Symbol:      t.Symbol,
		Decimal:     t.Decimals,
		TotalSupply: t.TotalSupply.String(),
		ChainId:     chain_params.G.ChainID,
		Block:       t.BlockNumber,
		BlockAt:     t.BlockTime,
		Program:     t.Program,
	}

	return ormToken.Normalize()
}
