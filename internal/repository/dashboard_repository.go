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
	GetMemberStats(ctx context.Context, cooperativeID, memberID string) (*model.MemberDashboardStats, error)
	GetMemberUpcomingInstallments(ctx context.Context, cooperativeID, memberID string) ([]model.MemberUpcomingInstallment, error)
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

func (r *dashboardRepository) GetMemberStats(ctx context.Context, cooperativeID, memberID string) (*model.MemberDashboardStats, error) {
	var stats model.MemberDashboardStats
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			m.full_name AS member_name,
			c.name AS cooperative_name,
			COALESCE(SUM(st.amount) FILTER (WHERE st.savings_type = 'pokok' AND st.direction = 'setor'), 0) AS pokok,
			COALESCE(SUM(st.amount) FILTER (WHERE st.savings_type = 'wajib' AND st.direction = 'setor'), 0)
			  - COALESCE(SUM(st.amount) FILTER (WHERE st.savings_type = 'wajib' AND st.direction = 'tarik'), 0) AS wajib,
			COALESCE(SUM(st.amount) FILTER (WHERE st.savings_type = 'sukarela' AND st.direction = 'setor'), 0)
			  - COALESCE(SUM(st.amount) FILTER (WHERE st.savings_type = 'sukarela' AND st.direction = 'tarik'), 0) AS sukarela,
			(SELECT COUNT(*) FROM loan
			 WHERE member_id = ? AND cooperative_id = ? AND status = 'aktif') AS active_loans,
			(SELECT COALESCE(SUM(l.principal - COALESCE(paid.sum_principal, 0)), 0)
			 FROM loan l
			 LEFT JOIN (
			     SELECT i2.loan_id, SUM(i2.principal_due) AS sum_principal
			     FROM installment_schedule i2
			     WHERE i2.status = 'lunas'
			     GROUP BY i2.loan_id
			 ) paid ON paid.loan_id = l.id
			 WHERE l.member_id = ? AND l.cooperative_id = ? AND l.status = 'aktif') AS outstanding_amount,
			(SELECT COUNT(*) FROM installment_schedule i
			 JOIN loan l ON l.id = i.loan_id
			 WHERE l.member_id = ? AND l.cooperative_id = ? AND i.status = 'terlambat') AS overdue_installments,
			(SELECT COALESCE(SUM(total_shu), 0) FROM shu_distribution WHERE member_id = ?) AS estimated_shu
		FROM member m
		JOIN cooperative c ON c.id = m.cooperative_id
		LEFT JOIN savings_transaction st ON st.member_id = m.id AND st.cooperative_id = m.cooperative_id
		WHERE m.id = ? AND m.cooperative_id = ?
		GROUP BY m.full_name, c.name
	`, memberID, cooperativeID, memberID, cooperativeID, memberID, cooperativeID, memberID, memberID, cooperativeID).Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *dashboardRepository) GetMemberUpcomingInstallments(ctx context.Context, cooperativeID, memberID string) ([]model.MemberUpcomingInstallment, error) {
	type row struct {
		InstallmentID string `gorm:"column:installment_id"`
		LoanID        string `gorm:"column:loan_id"`
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
			i.period_no,
			i.due_date::text AS due_date,
			i.total_due,
			i.status
		FROM installment_schedule i
		JOIN loan l ON l.id = i.loan_id
		WHERE l.member_id = ? AND l.cooperative_id = ?
		  AND (
		    i.status = 'terlambat'
		    OR (i.status = 'belum_bayar' AND i.due_date <= (NOW() + INTERVAL '30 days')::date)
		  )
		ORDER BY
		  CASE WHEN i.status = 'terlambat' THEN 0 ELSE 1 END ASC,
		  i.due_date ASC
		LIMIT 10
	`, memberID, cooperativeID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]model.MemberUpcomingInstallment, 0, len(rows))
	for _, r := range rows {
		result = append(result, model.MemberUpcomingInstallment{
			InstallmentID: r.InstallmentID,
			LoanID:        r.LoanID,
			PeriodNo:      r.PeriodNo,
			DueDate:       r.DueDate,
			TotalDue:      r.TotalDue,
			Status:        r.Status,
		})
	}
	return result, nil
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
