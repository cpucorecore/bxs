package parser

import (
	"bxs/log"
	"bxs/parser/event_parser"
	"bxs/types"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrParserNotFound = errors.New("parser not found")
)

type TopicRouter interface {
	Parse(ethLog *ethtypes.Log) (types.Event, error)
}

type topicRouter struct {
	topic2EventParser map[common.Hash]event_parser.EventParser
}

func NewTopicRouter() TopicRouter {
	return &topicRouter{
		topic2EventParser: event_parser.Topic2EventParser,
	}
}

func (p *topicRouter) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	topic0 := ethLog.Topics[0]
	if topic0.String() == "0x08e034a062383b8ce2fb4c25ef30ade713769d90b38db5207a7eb29b64340ef3" {
		log.Logger.Info(fmt.Sprintf("%s", ethLog.Topics[0].String()))
	}
	eventParser, ok := p.topic2EventParser[ethLog.Topics[0]]
	if !ok {
		return nil, ErrParserNotFound
	}

	return eventParser.Parse(ethLog)
}
