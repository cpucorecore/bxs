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
	MintEvent           types.Event
}

func (e *CreatedEvent) CanGetPair() bool {
	return true
}

func (e *CreatedEvent) GetPair() *types.Pair {
	if e.MintEvent != nil {
		e.Pair.Token0InitAmount, e.Pair.Token1InitAmount = e.MintEvent.GetMintAmount()
	}
	e.Pair.BlockAt = e.BlockTime
	return e.Pair
}

func (e *CreatedEvent) IsCreatePair() bool {
	return true
}

func (e *CreatedEvent) CanGetPoolUpdate() bool {
	return true
}

func (e *CreatedEvent) GetPoolUpdate() *types.PoolUpdate {
	a0, a1 := ParseAmountsByPair(e.TokenInitAmount, e.BaseTokenInitAmount, e.Pair)
	return &types.PoolUpdate{
		Program:       types.ProtocolNameXLaunch,
		LogIndex:      e.EventCommon.LogIndex,
		Address:       e.EventCommon.Pair.Address,
		Token0Address: e.EventCommon.Pair.Token0Core.Address,
		Token1Address: e.EventCommon.Pair.Token1Core.Address,
		Token0Amount:  a0,
		Token1Amount:  a1,
	}
}
