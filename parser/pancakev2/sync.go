package event_parser

import (
	"bxs/types"
	"math/big"
)

type SyncEvent struct {
	*types.EventCommon
	amount0Wei *big.Int
	amount1Wei *big.Int
}

func (e *SyncEvent) CanGetPoolUpdate() bool {
	return true
}

func (e *SyncEvent) GetPoolUpdate() *types.PoolUpdate {
	pu := &types.PoolUpdate{
		LogIndex: e.LogIndex,
		Address:  e.ContractAddress.String(),
		Token0:   e.Pair.Token0.Address.String(),
		Token1:   e.Pair.Token1.Address.String(),
	}

	pu.Amount0, pu.Amount1 = types.ParseAmountsByPair(e.amount0Wei, e.amount1Wei, e.Pair)
	return pu
}

func (e *SyncEvent) IsSync() bool {
	return true
}
