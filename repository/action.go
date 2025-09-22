package repository

import (
	"bxs/repository/orm"
	"gorm.io/gorm"
)

type ActionRepository struct {
	*BaseRepository[orm.Action]
}

func NewActionRepository(db *gorm.DB) *ActionRepository {
	baseRepo := NewBaseRepository[orm.Action](db)
	return &ActionRepository{BaseRepository: baseRepo}
}

func (r *ActionRepository) GetById(id string) (*orm.Action, error) {
	var action orm.Action
	err := r.db.Where("id = ?", id).First(&action).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func (r *ActionRepository) DeleteById(id string) error {
	return r.db.Where("id = ?", id).Delete(&orm.Action{}).Error
}
