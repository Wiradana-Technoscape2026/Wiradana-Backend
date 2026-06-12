package repository

import (
	"context"
	"errors"

	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var ErrModuleNotFound = errors.New("modul tidak ditemukan")

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
	var modules []entity.CoopModule
	err := r.db.WithContext(ctx).
		Where("cooperative_id = ?", coopID).
		Order("module_key ASC").
		Find(&modules).Error
	return modules, err
}

func (r *moduleRepository) FindByKey(ctx context.Context, coopID, key string) (*entity.CoopModule, error) {
	var m entity.CoopModule
	err := r.db.WithContext(ctx).
		Where("cooperative_id = ? AND module_key = ?", coopID, key).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModuleNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *moduleRepository) Upsert(ctx context.Context, m *entity.CoopModule) error {
	var existing entity.CoopModule
	err := r.db.WithContext(ctx).
		Where("cooperative_id = ? AND module_key = ?", m.CooperativeID, m.ModuleKey).
		First(&existing).Error
	if err == nil {
		m.ID = existing.ID
	}
	return r.db.WithContext(ctx).Save(m).Error
}
