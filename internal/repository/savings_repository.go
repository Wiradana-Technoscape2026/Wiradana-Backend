package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var (
	ErrPokokAlreadyRecorded  = errors.New("simpanan pokok sudah pernah disetor")
	ErrCannotWithdrawMandatory = errors.New("simpanan pokok dan wajib tidak dapat ditarik")
	ErrInsufficientSukarela  = errors.New("saldo sukarela tidak mencukupi untuk penarikan")
)

type SavingsRepository interface {
	Create(ctx context.Context, tx *entity.SavingsTransaction) error
	FindByMemberID(ctx context.Context, memberID string) ([]entity.SavingsTransaction, error)
	CountPokokSetoran(ctx context.Context, memberID string) (int64, error)
	GetSukarelaSaldo(ctx context.Context, memberID string) (int64, error)
	FindRecentByMember(ctx context.Context, memberID string, limit int) ([]entity.SavingsTransaction, error)
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

func (r *savingsRepository) FindByMemberID(ctx context.Context, memberID string) ([]entity.SavingsTransaction, error) {
	var txs []entity.SavingsTransaction
	err := r.db.WithContext(ctx).
		Where("member_id = ?", memberID).
		Order("created_at DESC").
		Find(&txs).Error
	return txs, err
}

func (r *savingsRepository) CountPokokSetoran(ctx context.Context, memberID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.SavingsTransaction{}).
		Where("member_id = ? AND savings_type = ? AND direction = ?", memberID, "pokok", "setor").
		Count(&count).Error
	return count, err
}

func (r *savingsRepository) GetSukarelaSaldo(ctx context.Context, memberID string) (int64, error) {
	if _, err := uuid.Parse(memberID); err != nil {
		return 0, nil
	}
	var saldo int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(
			SUM(CASE WHEN direction='setor' THEN amount ELSE -amount END), 0)
		FROM savings_transaction
		WHERE member_id = ? AND savings_type = 'sukarela'`, memberID).Scan(&saldo).Error
	return saldo, err
}

func (r *savingsRepository) FindRecentByMember(ctx context.Context, memberID string, limit int) ([]entity.SavingsTransaction, error) {
	var txs []entity.SavingsTransaction
	err := r.db.WithContext(ctx).
		Where("member_id = ?", memberID).
		Order("created_at DESC").
		Limit(limit).
		Find(&txs).Error
	return txs, err
}
