package repository

import (
	"context"
	"time"

	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

type SavingsRepository interface {
	Create(ctx context.Context, tx *entity.SavingsTransaction) error
	FindByMember(ctx context.Context, cooperativeID, memberID string) ([]*entity.SavingsTransaction, error)
	CountPokok(ctx context.Context, memberID, cooperativeID string) (int64, error)
	GetSukarelaBalance(ctx context.Context, memberID, cooperativeID string) (int64, error)
	GetCoopSummary(ctx context.Context, cooperativeID string) (*CoopSavingsSummaryRow, error)
	FindAllRecent(ctx context.Context, cooperativeID, savingsType, direction string, limit, offset int) ([]*SavingsTxWithMemberRow, int64, error)
}

type CoopSavingsSummaryRow struct {
	Pokok    int64 `gorm:"column:pokok"`
	Wajib    int64 `gorm:"column:wajib"`
	Sukarela int64 `gorm:"column:sukarela"`
}

type SavingsTxWithMemberRow struct {
	ID          string    `gorm:"column:id"`
	MemberID    string    `gorm:"column:member_id"`
	MemberName  string    `gorm:"column:member_name"`
	SavingsType string    `gorm:"column:savings_type"`
	Direction   string    `gorm:"column:direction"`
	Amount      int64     `gorm:"column:amount"`
	CreatedAt   time.Time `gorm:"column:created_at"`
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

func (r *savingsRepository) GetCoopSummary(ctx context.Context, cooperativeID string) (*CoopSavingsSummaryRow, error) {
	var row CoopSavingsSummaryRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(amount) FILTER (WHERE savings_type = 'pokok'    AND direction = 'setor'), 0) AS pokok,
			COALESCE(SUM(amount) FILTER (WHERE savings_type = 'wajib'    AND direction = 'setor'), 0)
		  - COALESCE(SUM(amount) FILTER (WHERE savings_type = 'wajib'    AND direction = 'tarik'), 0) AS wajib,
			COALESCE(SUM(amount) FILTER (WHERE savings_type = 'sukarela' AND direction = 'setor'), 0)
		  - COALESCE(SUM(amount) FILTER (WHERE savings_type = 'sukarela' AND direction = 'tarik'), 0) AS sukarela
		FROM savings_transaction
		WHERE cooperative_id = ?
	`, cooperativeID).Scan(&row).Error
	return &row, err
}

func (r *savingsRepository) FindAllRecent(ctx context.Context, cooperativeID, savingsType, direction string, limit, offset int) ([]*SavingsTxWithMemberRow, int64, error) {
	var rows []*SavingsTxWithMemberRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT st.id, st.member_id, m.full_name AS member_name,
		       st.savings_type, st.direction, st.amount, st.created_at
		FROM savings_transaction st
		JOIN member m ON m.id = st.member_id
		WHERE st.cooperative_id = ?
		  AND (? = '' OR st.savings_type = ?)
		  AND (? = '' OR st.direction = ?)
		ORDER BY st.created_at DESC
		LIMIT ? OFFSET ?
	`, cooperativeID, savingsType, savingsType, direction, direction, limit, offset).Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM savings_transaction
		WHERE cooperative_id = ?
		  AND (? = '' OR savings_type = ?)
		  AND (? = '' OR direction = ?)
	`, cooperativeID, savingsType, savingsType, direction, direction).Scan(&total).Error
	return rows, total, err
}
