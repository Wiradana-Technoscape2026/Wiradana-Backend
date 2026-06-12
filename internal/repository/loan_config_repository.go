package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var ErrLoanConfigNotFound = errors.New("konfigurasi pinjaman tidak ditemukan")

type LoanConfigRepository interface {
	FindByCoopID(ctx context.Context, coopID string) (*entity.LoanConfig, error)
	Upsert(ctx context.Context, lc *entity.LoanConfig) error
}

type loanConfigRepository struct{ db *gorm.DB }

func NewLoanConfigRepository(db *gorm.DB) LoanConfigRepository {
	return &loanConfigRepository{db: db}
}

func (r *loanConfigRepository) FindByCoopID(ctx context.Context, coopID string) (*entity.LoanConfig, error) {
	var lc entity.LoanConfig
	err := r.db.WithContext(ctx).Where("cooperative_id = ?", coopID).First(&lc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrLoanConfigNotFound
	}
	return &lc, err
}

func (r *loanConfigRepository) Upsert(ctx context.Context, lc *entity.LoanConfig) error {
	if lc.ID == uuid.Nil {
		lc.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Save(lc).Error
}
