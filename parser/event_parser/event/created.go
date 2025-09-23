package event

import (
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
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
		Program:       types.ProtocolNameXLaunch,
		LogIndex:      e.EventCommon.LogIndex,
		Address:       e.EventCommon.Pair.Address,
		Token0Address: e.EventCommon.Pair.Token0Core.Address,
		Token1Address: e.EventCommon.Pair.Token1Core.Address,
	}
	u.Token0Amount, u.Token1Amount = ParseAmountsByPair(e.TokenInitAmount, e.BaseTokenInitAmount, e.Pair)
	return u
}
