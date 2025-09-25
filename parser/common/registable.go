package event_parser

import (
	"github.com/ethereum/go-ethereum/common"
)

type Registrable interface {
	Register(common.Hash, EventParser)
}
