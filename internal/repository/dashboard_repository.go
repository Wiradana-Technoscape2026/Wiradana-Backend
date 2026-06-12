package repository

import (
	"context"

	"github.com/wiradana/backend/internal/model"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetStats(ctx context.Context, cooperativeID string) (*dashboardStats, error)
	GetUpcomingNotifications(ctx context.Context, cooperativeID string) ([]model.DashboardNotification, error)
}

type dashboardStats struct {
	TotalMembers int64 `gorm:"column:total_members"`
	TotalSavings int64 `gorm:"column:total_savings"`
	ActiveLoans  int64 `gorm:"column:active_loans"`
	OverdueLoans int64 `gorm:"column:overdue_loans"`
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStats(ctx context.Context, cooperativeID string) (*dashboardStats, error) {
	var stats dashboardStats
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			(SELECT COUNT(*) FROM member
			 WHERE cooperative_id = ? AND status = 'aktif') AS total_members,
			(SELECT COALESCE(SUM(amount) FILTER (WHERE direction='setor'), 0)
			      - COALESCE(SUM(amount) FILTER (WHERE direction='tarik'), 0)
			 FROM savings_transaction WHERE cooperative_id = ?) AS total_savings,
			(SELECT COUNT(*) FROM loan
			 WHERE cooperative_id = ? AND status = 'aktif') AS active_loans,
			(SELECT COUNT(*) FROM loan
			 WHERE cooperative_id = ? AND status = 'menunggak') AS overdue_loans
	`, cooperativeID, cooperativeID, cooperativeID, cooperativeID).Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

type notificationRow struct {
	FullName string `gorm:"column:full_name"`
	DueDate  string `gorm:"column:due_date"`
}

func (r *dashboardRepository) GetUpcomingNotifications(ctx context.Context, cooperativeID string) ([]model.DashboardNotification, error) {
	var rows []notificationRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT m.full_name, i.due_date::text AS due_date
		FROM installment_schedule i
		JOIN loan l ON l.id = i.loan_id
		JOIN member m ON m.id = l.member_id
		WHERE l.cooperative_id = ?
		  AND i.status = 'belum_bayar'
		  AND i.due_date <= (NOW() + INTERVAL '3 days')::date
		ORDER BY i.due_date ASC
		LIMIT 20
	`, cooperativeID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	notifs := make([]model.DashboardNotification, 0, len(rows))
	for _, row := range rows {
		notifs = append(notifs, model.DashboardNotification{
			Type:       "jatuh_tempo",
			MemberName: row.FullName,
			DueDate:    row.DueDate,
		})
	}
	return notifs, nil
}
