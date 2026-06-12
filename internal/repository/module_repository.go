package repository

import (
	"context"

	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	FindAll(ctx context.Context, coopID string) ([]entity.CoopModule, error)
	FindByKey(ctx context.Context, coopID, key string) (*entity.CoopModule, error)
	Upsert(ctx context.Context, m *entity.CoopModule) error
}

type moduleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepository{db: db}
}

func (r *moduleRepository) FindAll(ctx context.Context, coopID string) ([]entity.CoopModule, error) {
	var mods []entity.CoopModule
	err := r.db.WithContext(ctx).
		Where("cooperative_id = ?", coopID).
		Find(&mods).Error
	return mods, err
}

func (r *moduleRepository) FindByKey(ctx context.Context, coopID, key string) (*entity.CoopModule, error) {
	var m entity.CoopModule
	err := r.db.WithContext(ctx).
		Where("cooperative_id = ? AND module_key = ?", coopID, key).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *moduleRepository) Upsert(ctx context.Context, m *entity.CoopModule) error {
	return r.db.WithContext(ctx).Save(m).Error
}
