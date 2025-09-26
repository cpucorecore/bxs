package event_parser

import (
	"bxs/repository/orm"
	"bxs/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	action = "on-pancake"
)

type PairCreatedEvent struct {
	*types.EventCommon
	Address       common.Address
	Token0        common.Address
	Token1        common.Address
	tokenReversed bool
}

func (e *PairCreatedEvent) IsWBNBPair() bool {
	return types.IsWBNB(e.Token0) || types.IsWBNB(e.Token1)
}

func (e *PairCreatedEvent) GetNonWBNBToken() common.Address {
	if types.IsWBNB(e.Token0) {
		return e.Token1
	} else {
		return e.Token0
	}
}

func (e *PairCreatedEvent) IsPairCreated() bool {
	return true
}

func (e *PairCreatedEvent) GetAction() *orm.Action {
	return &orm.Action{
		Maker:        e.Maker.String(),
		Token:        e.GetNonWBNBToken().String(),
		Pair:         e.Address.String(),
		Action:       action,
		TxHash:       e.TxHash.String(),
		Creator:      e.Maker.String(),
		Block:        e.BlockNumber,
		BlockAt:      e.BlockTime,
		Token0Amount: types.ZeroDecimal, // TODO fixme
		Token1Amount: types.ZeroDecimal, // TODO fixme
	}
}

func (e *PairCreatedEvent) IsTokenReverse() bool {
	return e.tokenReversed
}
