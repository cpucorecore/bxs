package repository

import (
	"bxs/chain"
	"bxs/repository/orm"
	"gorm.io/gorm"
)

type PairRepository struct {
	*BaseRepository[orm.Pair]
}

func NewPairRepository(db *gorm.DB) *PairRepository {
	baseRepo := NewBaseRepository[orm.Pair](db)
	return &PairRepository{BaseRepository: baseRepo}
}

func (r *PairRepository) GetByAddressAndChainId(address string) (*orm.Pair, error) {
	var pair orm.Pair
	err := r.db.Where("address = ? AND chain_id = ?", address, chain.ID).First(&pair).Error
	if err != nil {
		return nil, err
	}
	return &pair, nil
}

func (r *PairRepository) DeleteByAddressAndChainId(address string) error {
	return r.db.Where("address = ? AND chain_id = ?", address, chain.ID).Delete(&orm.Pair{}).Error
}
