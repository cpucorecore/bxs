package event

import (
	"bxs/types"
	"bxs/util"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	MaxNameLen        = 256
	MaxSymbolLen      = 256
	MaxURLLen         = 256
	MaxDescriptionLen = 512
	MaxTelegramLen    = 256
	MaxTwitterLen     = 256
	MaxWebsiteLen     = 256
)

type CreatedEvent struct {
	*types.EventCommon
	PoolAddress         common.Address
	Creator             common.Address
	TokenAddress        common.Address
	BaseTokenInitAmount *big.Int
	TokenInitAmount     *big.Int
	Name                string
	Symbol              string
	URL                 string
	Description         string
	Telegram            string
	Twitter             string
	Website             string
}

func (e *CreatedEvent) FormatString() {
	e.Name = util.TruncateToMaxChars(util.FormatUTF8(e.Name), MaxNameLen)
	e.Symbol = util.TruncateToMaxChars(util.FormatUTF8(e.Symbol), MaxSymbolLen)
	e.URL = util.TruncateToMaxChars(util.FormatUTF8(e.URL), MaxURLLen)
	e.Description = util.TruncateToMaxChars(util.FormatUTF8(e.Description), MaxDescriptionLen)
	e.Telegram = util.TruncateToMaxChars(util.FormatUTF8(e.Telegram), MaxTelegramLen)
	e.Twitter = util.TruncateToMaxChars(util.FormatUTF8(e.Twitter), MaxTwitterLen)
	e.Website = util.TruncateToMaxChars(util.FormatUTF8(e.Website), MaxWebsiteLen)
}

func (e *CreatedEvent) GetPairAddress() common.Address {
	return e.PoolAddress
}

func (e *CreatedEvent) GetPair() *types.Pair {
	return e.Pair
}

func (e *CreatedEvent) DoGetPair() *types.Pair {
	if e.Pair != nil {
		return e.Pair
	}

	pair := &types.Pair{
		Address: e.PoolAddress,
		Token0Core: &types.TokenCore{
			Address:  e.TokenAddress,
			Symbol:   e.Symbol,
			Decimals: types.Decimals18,
		},
		Token1Core: types.NativeTokenCore,
		Block:      e.BlockNumber,
		BlockAt:    e.BlockTime,
		ProtocolId: types.ProtocolIdXLaunch,
	}

	pair.Token0InitAmount, pair.Token1InitAmount = ParseAmountsByPair(e.TokenInitAmount, e.BaseTokenInitAmount, pair)
	pair.Token0 = e.DoGetToken0()
	e.Pair = pair
	return e.Pair
}

func (e *CreatedEvent) GetToken0() *types.Token {
	return e.Pair.Token0
}

func (e *CreatedEvent) DoGetToken0() *types.Token {
	return &types.Token{
		Address:     e.TokenAddress,
		Creator:     e.Creator,
		Name:        e.Name,
		Symbol:      e.Symbol,
		Decimals:    types.Decimals18,
		BlockNumber: e.BlockNumber,
		BlockTime:   e.BlockTime,
		Program:     types.ProtocolNameXLaunch,
		Filtered:    false,
		URL:         e.URL,
		Description: e.Description,
		Telegram:    e.Telegram,
		Twitter:     e.Twitter,
		Website:     e.Website,
	}
}

func (e *CreatedEvent) IsCreated() bool {
	return true
}

func (e *CreatedEvent) CanGetPoolUpdate() bool {
	return true
}

func (e *CreatedEvent) GetPoolUpdate() *types.PoolUpdate {
	u := &types.PoolUpdate{
		LogIndex: e.EventCommon.LogIndex,
		Address:  e.EventCommon.Pair.Address,
		Token0:   e.EventCommon.Pair.Token0Core.Address,
		Token1:   e.EventCommon.Pair.Token1Core.Address,
	}
	u.Amount0, u.Amount1 = ParseAmountsByPair(e.TokenInitAmount, e.BaseTokenInitAmount, e.Pair)
	return u
}
