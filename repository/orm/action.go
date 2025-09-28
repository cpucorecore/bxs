package orm

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Action struct {
	Id           uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id,omitempty"`
	Maker        string          `json:"maker"`
	Token        string          `json:"token"`
	Pair         string          `json:"pair"`
	Action       string          `json:"action"`
	TxHash       string          `json:"tx_hash"`
	Creator      string          `json:"creator"`
	Block        uint64          `json:"block"`
	BlockAt      time.Time       `json:"block_at"`
	Token0Amount decimal.Decimal `gorm:"-" json:"token0_amount,omitempty"`
	Token1Amount decimal.Decimal `gorm:"-" json:"token1_amount,omitempty"`
	CreatedAt    time.Time       `gorm:"autoCreateTime" json:"created_at,omitempty"`
}

func (a *Action) TableName() string {
	return "action"
}

func (a *Action) Equal(a2 *Action) bool {
	if a.Maker != a2.Maker {
		return false
	}
	if a.Token != a2.Token {
		return false
	}
	if a.Pair != a2.Pair {
		return false
	}
	if a.Action != a2.Action {
		return false
	}
	if a.TxHash != a2.TxHash {
		return false
	}
	if a.Creator != a2.Creator {
		return false
	}
	return true
}
