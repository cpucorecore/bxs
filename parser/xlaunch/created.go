package event_parser

import (
	"bxs/types"
	"bxs/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

const (
	MaxLenName        = 128
	MaxLenSymbol      = 64
	MaxLenCid         = 255
	MaxLenDescription = 1024
	MaxLenTelegram    = 255
	MaxLenTwitter     = 255
	MaxLenWebsite     = 255
	MaxLenTid         = 64
)

type CreatedEvent struct {
	*types.EventCommon
	PoolAddress         common.Address
	Creator             common.Address
	TokenAddress        common.Address
	BaseTokenInitAmount *big.Int
	TokenInitAmount     *big.Int
	TotalSupply         *big.Int
	Name                string
	Symbol              string
	Cid                 string
	Tid                 string
	Description         string
	Telegram            string
	Twitter             string
	Website             string
}

func (e *CreatedEvent) FormatString() {
	e.Name = util.TruncateToMaxChars(util.FormatUTF8(e.Name), MaxLenName)
	e.Symbol = util.TruncateToMaxChars(util.FormatUTF8(e.Symbol), MaxLenSymbol)
	e.Cid = util.TruncateToMaxChars(util.FormatUTF8(e.Cid), MaxLenCid)
	e.Tid = util.TruncateToMaxChars(util.FormatUTF8(e.Tid), MaxLenTid)
	e.Description = util.TruncateToMaxChars(util.FormatUTF8(e.Description), MaxLenDescription)
	e.Telegram = util.TruncateToMaxChars(util.FormatUTF8(e.Telegram), MaxLenTelegram)
	e.Twitter = util.TruncateToMaxChars(util.FormatUTF8(e.Twitter), MaxLenTwitter)
	e.Website = util.TruncateToMaxChars(util.FormatUTF8(e.Website), MaxLenWebsite)
}

func (e *CreatedEvent) GetPairAddress() common.Address {
	return e.PoolAddress
}

func (e *CreatedEvent) getPair() *types.Pair {
	pair := &types.Pair{
		Address: e.PoolAddress,
		Token0: &types.TokenTinyInfo{
			Address: e.TokenAddress,
			Symbol:  e.Symbol,
			Decimal: types.Decimal18,
		},
		Token1:     types.NativeTokenTinyInfo,
		Block:      e.BlockNumber,
		ProtocolId: protocolId,
	}

	pair.InitAmount0, pair.InitAmount1 = types.ParseAmount(e.TokenInitAmount, e.BaseTokenInitAmount, pair)
	e.Pair = pair
	return e.Pair
}

func (e *CreatedEvent) GetToken0() *types.Token {
	return &types.Token{
		Address:     e.TokenAddress,
		Creator:     e.Creator,
		Name:        e.Name,
		Symbol:      e.Symbol,
		Decimals:    xLaunchTokenDecimal,
		TotalSupply: decimal.NewFromBigInt(e.TotalSupply, -int32(xLaunchTokenDecimal)),
		BlockNumber: e.BlockNumber,
		Program:     protocolName,
		Filtered:    false,
		Cid:         e.Cid,
		Tid:         e.Tid,
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
		Address:  e.EventCommon.Pair.Address.String(),
		Token0:   e.EventCommon.Pair.Token0.Address.String(),
		Token1:   e.EventCommon.Pair.Token1.Address.String(),
	}
	u.Amount0, u.Amount1 = types.ParseAmount(e.TokenInitAmount, e.BaseTokenInitAmount, e.Pair)
	return u
}

func (e *CreatedEvent) SetBlockTime(blockTime time.Time) {
	e.BlockTime = blockTime
	e.Pair.BlockAt = e.BlockTime
}
