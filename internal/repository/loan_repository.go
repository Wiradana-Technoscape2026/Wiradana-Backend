package repository

import (
	"context"
	"errors"

	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var ErrLoanNotFound = errors.New("pinjaman tidak ditemukan")

type LoanWithMeta struct {
	entity.Loan
	MemberName  string
	Outstanding int64
	Schedule    []InstallmentWithPaid
}

type LoanRepository interface {
	Create(ctx context.Context, loan *entity.Loan) error
	BulkCreateSchedule(ctx context.Context, schedules []entity.InstallmentSchedule) error
	FindAll(ctx context.Context, coopID, status string) ([]*LoanWithMeta, error)
	FindAllByMember(ctx context.Context, memberID string) ([]*LoanWithMeta, error)
	FindByID(ctx context.Context, coopID, loanID string) (*LoanWithMeta, error)
	UpdateStatus(ctx context.Context, loanID, status string) error
}

type loanRepository struct{ db *gorm.DB }

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (r *loanRepository) Create(ctx context.Context, loan *entity.Loan) error {
	return r.db.WithContext(ctx).Create(loan).Error
}

func (r *loanRepository) BulkCreateSchedule(ctx context.Context, schedules []entity.InstallmentSchedule) error {
	return r.db.WithContext(ctx).Create(&schedules).Error
}

type loanRow struct {
	entity.Loan
	MemberName  string `gorm:"column:member_name"`
	Outstanding int64  `gorm:"column:outstanding"`
}

const loanSelect = `loan.*,
	member.full_name as member_name,
	loan.principal - COALESCE(
		(SELECT SUM(i.principal_due) FROM installment_schedule i WHERE i.loan_id = loan.id AND i.status = 'lunas'),
	0) as outstanding`

func (r *loanRepository) scanLoans(ctx context.Context, q *gorm.DB) ([]*LoanWithMeta, error) {
	var rows []loanRow
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*LoanWithMeta, len(rows))
	for i, row := range rows {
		result[i] = &LoanWithMeta{Loan: row.Loan, MemberName: row.MemberName, Outstanding: row.Outstanding}
	}
	return result, nil
}

func (r *loanRepository) FindAll(ctx context.Context, coopID, status string) ([]*LoanWithMeta, error) {
	q := r.db.WithContext(ctx).Table("loan").Select(loanSelect).
		Joins("JOIN member ON member.id = loan.member_id").
		Where("loan.cooperative_id = ?", coopID)
	if status != "" {
		q = q.Where("loan.status = ?", status)
	}
	return r.scanLoans(ctx, q.Order("loan.disbursed_at DESC"))
}

func (r *loanRepository) FindAllByMember(ctx context.Context, memberID string) ([]*LoanWithMeta, error) {
	q := r.db.WithContext(ctx).Table("loan").Select(loanSelect).
		Joins("JOIN member ON member.id = loan.member_id").
		Where("loan.member_id = ?", memberID)
	return r.scanLoans(ctx, q.Order("loan.disbursed_at DESC"))
}

func (r *loanRepository) FindByID(ctx context.Context, coopID, loanID string) (*LoanWithMeta, error) {
	q := r.db.WithContext(ctx).Table("loan").Select(loanSelect).
		Joins("JOIN member ON member.id = loan.member_id").
		Where("loan.id = ?", loanID)
	if coopID != "" {
		q = q.Where("loan.cooperative_id = ?", coopID)
	}
	loans, err := r.scanLoans(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(loans) == 0 {
		return nil, ErrLoanNotFound
	}
	loan := loans[0]
	var schedules []InstallmentWithPaid
	r.db.WithContext(ctx).Raw(`
		SELECT i.*,
			COALESCE((SELECT SUM(p.amount) FROM payment p WHERE p.schedule_id = i.id), 0) as paid_amount
		FROM installment_schedule i WHERE i.loan_id = ? ORDER BY i.period_no`, loanID).Scan(&schedules)
	loan.Schedule = schedules
	return loan, nil
}

func (r *loanRepository) UpdateStatus(ctx context.Context, loanID, status string) error {
	return r.db.WithContext(ctx).Model(&entity.Loan{}).
		Where("id = ?", loanID).Update("status", status).Error
}
