package repository

import (
	"context"

	"github.com/wiradana/backend/internal/model"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetStats(ctx context.Context, cooperativeID string) (*model.DashboardStats, error)
	GetUpcomingInstallmentsCount(ctx context.Context, cooperativeID string) (int64, error)
	GetUpcomingInstallments(ctx context.Context, cooperativeID string) ([]model.UpcomingInstallment, error)
	GetPendingApplicationsCount(ctx context.Context, cooperativeID string) (int64, error)
	GetPendingApplications(ctx context.Context, cooperativeID string) ([]model.PendingApplication, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStats(ctx context.Context, cooperativeID string) (*model.DashboardStats, error) {
	var stats model.DashboardStats
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			(SELECT COUNT(*) FROM member
			 WHERE cooperative_id = ? AND status = 'aktif') AS active_members,
			(SELECT COUNT(*) FROM member
			 WHERE cooperative_id = ?) AS total_members,
			(SELECT COALESCE(SUM(amount) FILTER (WHERE direction='setor'), 0)
			      - COALESCE(SUM(amount) FILTER (WHERE direction='tarik'), 0)
			 FROM savings_transaction WHERE cooperative_id = ?) AS total_savings,
			(SELECT COUNT(*) FROM loan
			 WHERE cooperative_id = ? AND status = 'aktif') AS active_loans,
			(SELECT COALESCE(SUM(l.principal - COALESCE(paid.sum_principal, 0)), 0)
			 FROM loan l
			 LEFT JOIN (
			   SELECT i2.loan_id, SUM(i2.principal_due) AS sum_principal
			   FROM installment_schedule i2
			   WHERE i2.status = 'lunas'
			   GROUP BY i2.loan_id
			 ) paid ON paid.loan_id = l.id
			 WHERE l.cooperative_id = ? AND l.status = 'aktif') AS active_loans_outstanding,
			(SELECT COUNT(*) FROM loan
			 WHERE cooperative_id = ? AND status = 'menunggak') AS overdue_loans
	`, cooperativeID, cooperativeID, cooperativeID,
		cooperativeID, cooperativeID, cooperativeID).Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *dashboardRepository) GetUpcomingInstallmentsCount(ctx context.Context, cooperativeID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM installment_schedule i
		JOIN loan l ON l.id = i.loan_id
		WHERE l.cooperative_id = ?
		  AND (
		    i.status = 'terlambat'
		    OR (i.status = 'belum_bayar' AND i.due_date <= (NOW() + INTERVAL '30 days')::date)
		  )
	`, cooperativeID).Scan(&count).Error
	return count, err
}

func (r *dashboardRepository) GetUpcomingInstallments(ctx context.Context, cooperativeID string) ([]model.UpcomingInstallment, error) {
	type row struct {
		InstallmentID string `gorm:"column:installment_id"`
		LoanID        string `gorm:"column:loan_id"`
		MemberName    string `gorm:"column:member_name"`
		PeriodNo      int    `gorm:"column:period_no"`
		DueDate       string `gorm:"column:due_date"`
		TotalDue      int64  `gorm:"column:total_due"`
		Status        string `gorm:"column:status"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT
		  i.id        AS installment_id,
		  i.loan_id,
		  m.full_name AS member_name,
		  i.period_no,
		  i.due_date::text AS due_date,
		  i.total_due,
		  i.status
		FROM installment_schedule i
		JOIN loan l ON l.id = i.loan_id
		JOIN member m ON m.id = l.member_id
		WHERE l.cooperative_id = ?
		  AND (
		    i.status = 'terlambat'
		    OR (i.status = 'belum_bayar' AND i.due_date <= (NOW() + INTERVAL '30 days')::date)
		  )
		ORDER BY
		  CASE WHEN i.status = 'terlambat' THEN 0 ELSE 1 END ASC,
		  i.due_date ASC
		LIMIT 20
	`, cooperativeID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]model.UpcomingInstallment, 0, len(rows))
	for _, r := range rows {
		result = append(result, model.UpcomingInstallment{
			InstallmentID: r.InstallmentID,
			LoanID:        r.LoanID,
			MemberName:    r.MemberName,
			PeriodNo:      r.PeriodNo,
			DueDate:       r.DueDate,
			TotalDue:      r.TotalDue,
			Status:        r.Status,
		})
	}
	return result, nil
}

func (r *dashboardRepository) GetPendingApplicationsCount(ctx context.Context, cooperativeID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM loan_application
		WHERE cooperative_id = ? AND status = 'pending'
	`, cooperativeID).Scan(&count).Error
	return count, err
}

func (r *dashboardRepository) GetPendingApplications(ctx context.Context, cooperativeID string) ([]model.PendingApplication, error) {
	type row struct {
		ID          string `gorm:"column:id"`
		MemberName  string `gorm:"column:member_name"`
		Amount      int64  `gorm:"column:amount"`
		TenorMonths int    `gorm:"column:tenor_months"`
		Purpose     string `gorm:"column:purpose"`
		Grade       string `gorm:"column:grade"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT
		  la.id,
		  m.full_name    AS member_name,
		  la.amount,
		  la.tenor_months,
		  la.purpose,
		  COALESCE(ca.grade, '') AS grade
		FROM loan_application la
		JOIN member m ON m.id = la.member_id
		LEFT JOIN credit_assessment ca ON ca.application_id = la.id
		WHERE la.cooperative_id = ? AND la.status = 'pending'
		ORDER BY la.created_at DESC
		LIMIT 10
	`, cooperativeID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]model.PendingApplication, 0, len(rows))
	for _, r := range rows {
		result = append(result, model.PendingApplication{
			ID:          r.ID,
			MemberName:  r.MemberName,
			Amount:      r.Amount,
			TenorMonths: r.TenorMonths,
			Purpose:     r.Purpose,
			Grade:       r.Grade,
		})
	}
	return result, nil
}
