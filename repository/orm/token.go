package orm

import (
	"time"
)

type Token struct {
	Address     string    `json:"address"`
	Creator     string    `json:"creator"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	Decimal     int8      `json:"decimal"`
	TotalSupply string    `json:"total_supply"`
	ChainId     int       `json:"chain_id"`
	Block       uint64    `json:"block"`
	BlockAt     time.Time `json:"block_at"`
	Program     string    `json:"program"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	MainPair    string    `json:"main_pair"`
	Telegram    string    `json:"telegram"`
	Twitter     string    `json:"twitter"`
	Website     string    `json:"website"`
}

func (t *Token) TableName() string {
	return "token"
}

func (t *Token) Equal(t2 *Token) bool {
	if t.Address != t2.Address {
		return false
	}
	if t.Name != t2.Name {
		return false
	}
	if t.Symbol != t2.Symbol {
		return false
	}
	if t.Decimal != t2.Decimal {
		return false
	}
	if t.TotalSupply != t2.TotalSupply {
		return false
	}
	if t.Block != t2.Block {
		return false
	}
	return true
}
