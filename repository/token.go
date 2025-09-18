package repository

import (
	"bxs/chain"
	"bxs/repository/orm"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TokenRepository struct {
	*BaseRepository[orm.Token]
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	baseRepo := NewBaseRepository[orm.Token](db)
	return &TokenRepository{BaseRepository: baseRepo}
}

func (r *TokenRepository) GetByAddressAndChainId(address string) (*orm.Token, error) {
	var token orm.Token
	err := r.db.Where("address = ? AND chain_id = ?", address, chain.ID).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepository) UpdateMainPair(address string, mainPair string) error {
	return r.db.Model(&orm.Token{}).
		Where("address = ? AND chain_id = ?", address, chain.ID).
		Update("main_pair", mainPair).Error
}

func (r *TokenRepository) DeleteByAddressAndChainId(address string) error {
	return r.db.Where("address = ? AND chain_id = ?", address, chain.ID).Delete(&orm.Token{}).Error
}
