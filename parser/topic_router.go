package parser

import (
	pcommon "bxs/parser/common"
	ppancakev2 "bxs/parser/pancakev2"
	pxlaunch "bxs/parser/xlaunch"
	"bxs/types"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrParserNotFound = errors.New("parser not found")
)

type TopicRouter interface {
	Route(ethLog *ethtypes.Log) (types.Event, error)
}

type topicRouter struct {
	topic2EventParser map[common.Hash]pcommon.EventParser
}

func NewTopicRouter() TopicRouter {
	r := &topicRouter{
		topic2EventParser: make(map[common.Hash]pcommon.EventParser),
	}

	ppancakev2.Reg(r)
	pxlaunch.Reg(r)

	return r
}

func (p *topicRouter) Route(ethLog *ethtypes.Log) (types.Event, error) {
	eventParser, ok := p.topic2EventParser[ethLog.Topics[0]]
	if !ok {
		return nil, ErrParserNotFound
	}

	return eventParser.Parse(ethLog)
}

func (p *topicRouter) Register(commonHash common.Hash, eventParser pcommon.EventParser) {
	p.topic2EventParser[commonHash] = eventParser
}
