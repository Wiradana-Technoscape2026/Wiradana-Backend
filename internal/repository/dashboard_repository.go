package repository

import (
	"context"
<<<<<<< HEAD
	"time"
=======
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed

	"github.com/wiradana/backend/internal/model"
	"gorm.io/gorm"
)

type DashboardRepository interface {
<<<<<<< HEAD
	GetStats(ctx context.Context, coopID string) (*model.DashboardResponse, error)
=======
	GetStats(ctx context.Context, cooperativeID string) (*dashboardStats, error)
	GetUpcomingNotifications(ctx context.Context, cooperativeID string) ([]model.DashboardNotification, error)
}

type dashboardStats struct {
	TotalMembers int64 `gorm:"column:total_members"`
	TotalSavings int64 `gorm:"column:total_savings"`
	ActiveLoans  int64 `gorm:"column:active_loans"`
	OverdueLoans int64 `gorm:"column:overdue_loans"`
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

<<<<<<< HEAD
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
=======
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
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
		FROM installment_schedule i
		JOIN loan l ON l.id = i.loan_id
		JOIN member m ON m.id = l.member_id
		WHERE l.cooperative_id = ?
		  AND i.status = 'belum_bayar'
<<<<<<< HEAD
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
=======
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
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
}
