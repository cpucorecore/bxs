package event_parser

import (
	"bxs/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type EventParser interface {
	Parse(ethLog *ethtypes.Log) (types.Event, error)
}
