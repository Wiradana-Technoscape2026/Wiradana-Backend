package repository

import (
	"context"
	"time"

	"github.com/wiradana/backend/internal/model"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetStats(ctx context.Context, coopID string) (*model.DashboardResponse, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStats(ctx context.Context, coopID string) (*model.DashboardResponse, error) {
	resp := &model.DashboardResponse{}

	// Total members aktif
	r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM member WHERE cooperative_id = ? AND status = 'aktif'`, coopID).Scan(&resp.TotalMembers)

	// Total simpanan (semua tipe: setor - tarik)
	r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(CASE WHEN direction='setor' THEN amount ELSE -amount END), 0)
		FROM savings_transaction WHERE cooperative_id = ?`, coopID).Scan(&resp.TotalSavings)

	// Active loans
	r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM loan WHERE cooperative_id = ? AND status = 'aktif'`, coopID).Scan(&resp.ActiveLoans)

	// Overdue loans
	r.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT l.id) FROM loan l
		JOIN installment_schedule i ON i.loan_id = l.id
		WHERE l.cooperative_id = ? AND i.status = 'terlambat'`, coopID).Scan(&resp.OverdueLoans)

	// Notifications: installments due within 3 days and belum_bayar
	threeDaysFromNow := time.Now().AddDate(0, 0, 3).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	type notifRow struct {
		MemberName string `gorm:"column:member_name"`
		DueDate    string `gorm:"column:due_date"`
		PeriodNo   int    `gorm:"column:period_no"`
		TotalDue   int64  `gorm:"column:total_due"`
		LoanID     string `gorm:"column:loan_id"`
	}
	var rows []notifRow
	r.db.WithContext(ctx).Raw(`
		SELECT m.full_name as member_name,
		       i.due_date::text,
		       i.period_no,
		       i.total_due,
		       i.loan_id::text
		FROM installment_schedule i
		JOIN loan l ON l.id = i.loan_id
		JOIN member m ON m.id = l.member_id
		WHERE l.cooperative_id = ?
		  AND i.status = 'belum_bayar'
		  AND i.due_date::date <= ?::date
		  AND i.due_date::date >= ?::date
		ORDER BY i.due_date ASC
		LIMIT 20`, coopID, threeDaysFromNow, today).Scan(&rows)

	resp.Notifications = make([]model.Notification, len(rows))
	for i, row := range rows {
		resp.Notifications[i] = model.Notification{
			Type:       "jatuh_tempo",
			MemberName: row.MemberName,
			DueDate:    row.DueDate,
			PeriodNo:   row.PeriodNo,
			TotalDue:   row.TotalDue,
			LoanID:     row.LoanID,
		}
	}

	return resp, nil
}
