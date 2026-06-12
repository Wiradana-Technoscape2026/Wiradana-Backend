package repository

import (
	"context"

	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

type SavingsRepository interface {
	Create(ctx context.Context, tx *entity.SavingsTransaction) error
	FindByMember(ctx context.Context, cooperativeID, memberID string) ([]*entity.SavingsTransaction, error)
	CountPokok(ctx context.Context, memberID, cooperativeID string) (int64, error)
	GetSukarelaBalance(ctx context.Context, memberID, cooperativeID string) (int64, error)
}

type savingsRepository struct {
	db *gorm.DB
}

func NewSavingsRepository(db *gorm.DB) SavingsRepository {
	return &savingsRepository{db: db}
}

func (r *savingsRepository) Create(ctx context.Context, tx *entity.SavingsTransaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *savingsRepository) FindByMember(ctx context.Context, cooperativeID, memberID string) ([]*entity.SavingsTransaction, error) {
	var txs []*entity.SavingsTransaction
	err := r.db.WithContext(ctx).
		Where("member_id = ? AND cooperative_id = ?", memberID, cooperativeID).
		Order("created_at DESC").
		Find(&txs).Error
	return txs, err
}

func (r *savingsRepository) CountPokok(ctx context.Context, memberID, cooperativeID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.SavingsTransaction{}).
		Where("member_id = ? AND cooperative_id = ? AND savings_type = 'pokok' AND direction = 'setor'", memberID, cooperativeID).
		Count(&count).Error
	return count, err
}

func (r *savingsRepository) GetSukarelaBalance(ctx context.Context, memberID, cooperativeID string) (int64, error) {
	var balance int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(amount) FILTER (WHERE direction = 'setor'), 0)
		  - COALESCE(SUM(amount) FILTER (WHERE direction = 'tarik'), 0)
		FROM savings_transaction
		WHERE member_id = ? AND cooperative_id = ? AND savings_type = 'sukarela'
	`, memberID, cooperativeID).Scan(&balance).Error
	return balance, err
}
