package types

import (
	"bxs/chain_params"
	"bxs/repository/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"time"
)

var (
	NativeTokenTinyInfo = &TokenTinyInfo{
		Address: ZeroAddress,
		Symbol:  NativeTokenSymbol,
		Decimal: Decimal18,
	}
	WBNBSymbol  = "WBNB"
	WBNBDecimal = Decimal18
)

func IsSameAddress(address1, address2 common.Address) bool {
	return address1.Cmp(address2) == 0
}

func IsWBNB(address common.Address) bool {
	return IsSameAddress(address, chain_params.G.WBNBAddress)
}

func OrderAddress(a0, a1 common.Address) (common.Address, common.Address, bool) {
	if IsSameAddress(a1, chain_params.G.WBNBAddress) {
		return a0, a1, false
	}
	return a1, a0, true
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
	Address     common.Address
	Creator     common.Address
	Name        string
	Symbol      string
	Decimals    int8
	TotalSupply decimal.Decimal
	BlockNumber uint64
	BlockTime   time.Time
	Program     string
	Filtered    bool
	Cid         string
	Tid         string
	Description string
	Telegram    string
	Twitter     string
	Website     string
	UpdateTs    time.Time
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
	return &orm.Token{
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
		Cid:         t.Cid,
		Tid:         t.Tid,
		Telegram:    t.Telegram,
		Twitter:     t.Twitter,
		Website:     t.Website,
		Description: t.Description,
	}
}

func (t *Token) GetTokenTinyInfo() *TokenTinyInfo {
	return &TokenTinyInfo{
		Address: t.Address,
		Symbol:  t.Symbol,
		Decimal: t.Decimals,
	}
}
