package repository

import (
	"context"
<<<<<<< HEAD
	"errors"

	"github.com/google/uuid"
=======

>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

<<<<<<< HEAD
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
=======
type SavingsRepository interface {
	Create(ctx context.Context, tx *entity.SavingsTransaction) error
	FindByMember(ctx context.Context, cooperativeID, memberID string) ([]*entity.SavingsTransaction, error)
	CountPokok(ctx context.Context, memberID, cooperativeID string) (int64, error)
	GetSukarelaBalance(ctx context.Context, memberID, cooperativeID string) (int64, error)
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
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

<<<<<<< HEAD
func (r *savingsRepository) FindByMemberID(ctx context.Context, memberID string) ([]entity.SavingsTransaction, error) {
	var txs []entity.SavingsTransaction
	err := r.db.WithContext(ctx).
		Where("member_id = ?", memberID).
=======
func (r *savingsRepository) FindByMember(ctx context.Context, cooperativeID, memberID string) ([]*entity.SavingsTransaction, error) {
	var txs []*entity.SavingsTransaction
	err := r.db.WithContext(ctx).
		Where("member_id = ? AND cooperative_id = ?", memberID, cooperativeID).
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
		Order("created_at DESC").
		Find(&txs).Error
	return txs, err
}

<<<<<<< HEAD
func (r *savingsRepository) CountPokokSetoran(ctx context.Context, memberID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.SavingsTransaction{}).
		Where("member_id = ? AND savings_type = ? AND direction = ?", memberID, "pokok", "setor").
=======
func (r *savingsRepository) CountPokok(ctx context.Context, memberID, cooperativeID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.SavingsTransaction{}).
		Where("member_id = ? AND cooperative_id = ? AND savings_type = 'pokok' AND direction = 'setor'", memberID, cooperativeID).
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
		Count(&count).Error
	return count, err
}

<<<<<<< HEAD
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
=======
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
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
}
