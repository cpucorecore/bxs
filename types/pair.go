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

type TokenCore struct {
	Address  common.Address `json:"-"`
	Symbol   string
	Decimals int8
}

func (t *TokenCore) MarshalJSON() ([]byte, error) {
	type Alias TokenCore
	return json.Marshal(&struct {
		AddressString string `json:"Address"`
		*Alias
	}{
		AddressString: t.Address.String(),
		Alias:         (*Alias)(t),
	})
}

func (t *TokenCore) UnmarshalJSON(data []byte) error {
	type Alias TokenCore
	aux := &struct {
		AddressString string `json:"Address"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Address = common.HexToAddress(aux.AddressString)
	return nil
}

func (t *TokenCore) Equal(token *TokenCore) bool {
	if !IsSameAddress(t.Address, token.Address) {
		return false
	}

	if t.Symbol != token.Symbol {
		return false
	}

	if t.Decimals != token.Decimals {
		return false
	}

	return true
}

type Pair struct {
	Address          common.Address `json:"-"`
	TokenReversed    bool
	Token0Core       *TokenCore
	Token1Core       *TokenCore
	Token0           *Token `json:"-"`
	Token1           *Token `json:"-"`
	Token0InitAmount decimal.Decimal
	Token1InitAmount decimal.Decimal
	Block            uint64
	BlockAt          time.Time
	ProtocolId       int
	Filtered         bool
	FilterCode       int
	Timestamp        time.Time
}

func (p *Pair) String() string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

func (p *Pair) MarshalBinary() ([]byte, error) {
	type Alias Pair
	return json.Marshal(&struct {
		AddressString string `json:"Address"`
		*Alias
	}{
		AddressString: p.Address.String(),
		Alias:         (*Alias)(p),
	})
}

func (p *Pair) UnmarshalBinary(data []byte) error {
	type Alias Pair
	aux := &struct {
		AddressString string `json:"Address"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.Address = common.HexToAddress(aux.AddressString)
	return nil
}

func (p *Pair) swapToken0Token1() {
	p.Token0Core, p.Token1Core = p.Token1Core, p.Token0Core
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
	if !p.Token0Core.Equal(pair.Token0Core) {
		return false
	}
	if !p.Token1Core.Equal(pair.Token1Core) {
		return false
	}
	if !p.Token0InitAmount.Equal(pair.Token0InitAmount) {
		return false
	}
	if !p.Token1InitAmount.Equal(pair.Token1InitAmount) {
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

func (p *Pair) IsFiltered() bool {
	return p.Filtered
}

func (p *Pair) FilterByToken0AndToken1() bool {
	if !IsBaseToken(p.Token0Core.Address) && !IsBaseToken(p.Token1Core.Address) {
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
		Name:     getPairName(p.Token0Core.Symbol, p.Token1Core.Symbol),
		Address:  p.Address.String(),
		Token0:   p.Token0Core.Address.String(),
		Token1:   p.Token1Core.Address.String(),
		Reserve0: p.Token0InitAmount.Mul(decimal.New(1, int32(p.Token0Core.Decimals))), // for db type is numeric(78)
		Reserve1: p.Token1InitAmount.Mul(decimal.New(1, int32(p.Token1Core.Decimals))), // for db type is numeric(78)
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
