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
