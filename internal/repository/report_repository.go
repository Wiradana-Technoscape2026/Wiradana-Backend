package repository

import (
	"context"

	"github.com/wiradana/backend/internal/model"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetSummary(ctx context.Context, cooperativeID string) (*model.ReportSummary, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetSummary(ctx context.Context, cooperativeID string) (*model.ReportSummary, error) {
	var summary model.ReportSummary
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			(SELECT COALESCE(SUM(amount) FILTER (WHERE direction='setor'), 0)
			       - COALESCE(SUM(amount) FILTER (WHERE direction='tarik'), 0)
			 FROM savings_transaction WHERE cooperative_id = ?) AS total_savings,

			(SELECT COALESCE(SUM(i.interest_due), 0)
			 FROM installment_schedule i
			 JOIN loan l ON l.id = i.loan_id
			 WHERE l.cooperative_id = ? AND i.status = 'lunas') AS total_interest_collected,

			(SELECT COALESCE(SUM(l.principal - COALESCE(paid.sum_principal, 0)), 0)
			 FROM loan l
			 LEFT JOIN (
			   SELECT i2.loan_id, SUM(i2.principal_due) AS sum_principal
			   FROM installment_schedule i2
			   WHERE i2.status = 'lunas'
			   GROUP BY i2.loan_id
			 ) paid ON paid.loan_id = l.id
			 WHERE l.cooperative_id = ? AND l.status != 'lunas') AS total_outstanding,

			(SELECT COUNT(*) FROM member
			 WHERE cooperative_id = ?) AS total_members,

			(SELECT COUNT(*) FROM member
			 WHERE cooperative_id = ? AND status = 'aktif') AS active_members,

			(SELECT COUNT(*) FROM loan
			 WHERE cooperative_id = ? AND status = 'aktif') AS active_loans_count,

			(SELECT COUNT(*) FROM loan
			 WHERE cooperative_id = ? AND status = 'menunggak') AS overdue_loans_count,

			(SELECT COALESCE(SUM(principal), 0) FROM loan
			 WHERE cooperative_id = ?) AS total_disbursed
	`, cooperativeID, cooperativeID, cooperativeID,
		cooperativeID, cooperativeID, cooperativeID,
		cooperativeID, cooperativeID).Scan(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}
