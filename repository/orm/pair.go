package orm

import (
	"github.com/shopspring/decimal"
	"time"
)

type Pair struct {
	Name      string          `json:"name"`
	Address   string          `json:"address"`
	Token0    string          `json:"token0"`
	Token1    string          `json:"token1"`
	ChainId   int             `json:"chain_id"`
	Reserve0  decimal.Decimal `json:"reserve0"`
	Reserve1  decimal.Decimal `json:"reserve1"`
	Block     uint64          `json:"block"`
	BlockAt   time.Time       `json:"block_at"`
	Program   string          `json:"program"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at,omitempty"`
}

func (p *Pair) TableName() string {
	return "pair"
}

func (p *Pair) Equal(p2 *Pair) bool {
	if p.Name != p2.Name {
		return false
	}
	if p.Address != p2.Address {
		return false
	}
	if p.Token0 != p2.Token0 {
		return false
	}
	if p.Token1 != p2.Token1 {
		return false
	}
	if p.ChainId != p2.ChainId {
		return false
	}
	if !p.Reserve0.Equal(p2.Reserve0) {
		return false
	}
	if !p.Reserve1.Equal(p2.Reserve1) {
		return false
	}
	return true
}
