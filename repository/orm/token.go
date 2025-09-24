package orm

import (
	"time"
)

type Token struct {
	Address     string
	Creator     string
	Name        string
	Symbol      string
	Decimal     int8
	TotalSupply string
	ChainId     int
	Block       uint64
	BlockAt     time.Time
	Program     string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	MainPair    string
	Description string
	Url         string
	Telegram    string
	Twitter     string
	Website     string
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
