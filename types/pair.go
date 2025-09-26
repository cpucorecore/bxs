package types

import (
	"bxs/chain_params"
	"bxs/repository/orm"
	"bxs/util"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"time"
)

const (
	FilterCodeGetToken = iota + 1
	FilterCodeVerifyFailed
	FilterCodeNoBaseToken
	FilterCodeNoXLaunchToken
)

type TokenTinyInfo struct {
	Address common.Address
	Symbol  string
	Decimal int8
}

func (t *TokenTinyInfo) IsBaseToken() bool {
	return IsBaseToken(t.Address)
}

type Pair struct {
	Address       common.Address
	Token0        *TokenTinyInfo
	Token1        *TokenTinyInfo
	TokenReversed bool
	InitAmount0   decimal.Decimal
	InitAmount1   decimal.Decimal
	Block         uint64
	BlockAt       time.Time
	ProtocolId    int
	Filtered      bool
	FilterCode    int
	UpdateTs      time.Time
}

func (p *Pair) String() string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

func (p *Pair) swapToken0Token1() {
	p.Token0, p.Token1 = p.Token1, p.Token0
	p.TokenReversed = true
}

func (p *Pair) Equal(pair *Pair) bool {
	if !IsSameAddress(p.Address, pair.Address) {
		return false
	}
	if p.TokenReversed != pair.TokenReversed {
		return false
	}
	if !IsSameAddress(p.Token0.Address, pair.Token0.Address) {
		return false
	}
	if !IsSameAddress(p.Token1.Address, pair.Token1.Address) {
		return false
	}
	if !p.InitAmount0.Equal(pair.InitAmount0) {
		return false
	}
	if !p.InitAmount1.Equal(pair.InitAmount1) {
		return false
	}
	if p.Block != pair.Block {
		return false
	}
	if p.ProtocolId != pair.ProtocolId {
		return false
	}
	if p.Filtered != pair.Filtered {
		return false
	}
	if p.FilterCode != pair.FilterCode {
		return false
	}
	return true
}

func (p *Pair) FilterNoBaseToken() bool {
	if !p.Token0.IsBaseToken() && !p.Token1.IsBaseToken() {
		p.Filtered = true
		p.FilterCode = FilterCodeNoBaseToken
	}

	return p.Filtered
}

/*
getPairName constructs a pair name from two token symbols.
- token0Symbol has a maximum of 64 characters
- token1Symbol has a maximum of 63 characters
- The combined format is `token0Symbol/token1Symbol` (total ≤ 128 characters)
to match PostgreSQL's varchar(128) character limit.
*/
func getPairName(token0Symbol, token1Symbol string) string {
	// Truncate by character count (not bytes)
	token0Symbol = util.TruncateToMaxChars(token0Symbol, 64)
	token1Symbol = util.TruncateToMaxChars(token1Symbol, 63)
	return token0Symbol + "/" + token1Symbol
}

func (p *Pair) GetOrmPair() *orm.Pair {
	return &orm.Pair{
		Name:     getPairName(p.Token0.Symbol, p.Token1.Symbol),
		Address:  p.Address.String(),
		Token0:   p.Token0.Address.String(),
		Token1:   p.Token1.Address.String(),
		Reserve0: p.InitAmount0,
		Reserve1: p.InitAmount1,
		ChainId:  chain_params.G.ChainID,
		Block:    p.Block,
		BlockAt:  p.BlockAt,
		Program:  GetProtocolName(p.ProtocolId),
	}
}

type PairWrap struct {
	Pair      *Pair
	NewPair   bool
	NewToken0 bool
	NewToken1 bool
}
