package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var ErrInstallmentNotFound = errors.New("angsuran tidak ditemukan")

type InstallmentWithPaid struct {
	entity.InstallmentSchedule
	PaidAmount int64 `gorm:"column:paid_amount"`
}

type InstallmentRepository interface {
	FindByID(ctx context.Context, installmentID string) (*InstallmentWithPaid, error)
	FindByLoanID(ctx context.Context, loanID string) ([]InstallmentWithPaid, error)
	CreatePayment(ctx context.Context, p *entity.Payment) error
	GetPaidAmount(ctx context.Context, scheduleID string) (int64, error)
	UpdateInstallmentStatus(ctx context.Context, scheduleID, status string) error
	HasOverdueInstallments(ctx context.Context, loanID string) (bool, error)
	AllInstallmentsPaid(ctx context.Context, loanID string) (bool, error)
}

type installmentRepository struct{ db *gorm.DB }

func NewInstallmentRepository(db *gorm.DB) InstallmentRepository {
	return &installmentRepository{db: db}
}

func (r *installmentRepository) FindByID(ctx context.Context, scheduleID string) (*InstallmentWithPaid, error) {
	var result InstallmentWithPaid
	err := r.db.WithContext(ctx).Raw(`
		SELECT i.*,
			COALESCE((SELECT SUM(p.amount) FROM payment p WHERE p.schedule_id = i.id), 0) as paid_amount
		FROM installment_schedule i WHERE i.id = ?`, scheduleID).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == uuid.Nil {
		return nil, ErrInstallmentNotFound
	}
	return &result, nil
}

func (r *installmentRepository) FindByLoanID(ctx context.Context, loanID string) ([]InstallmentWithPaid, error) {
	var result []InstallmentWithPaid
	err := r.db.WithContext(ctx).Raw(`
		SELECT i.*,
			COALESCE((SELECT SUM(p.amount) FROM payment p WHERE p.schedule_id = i.id), 0) as paid_amount
		FROM installment_schedule i WHERE i.loan_id = ? ORDER BY i.period_no`, loanID).Scan(&result).Error
	return result, err
}

func (r *installmentRepository) CreatePayment(ctx context.Context, p *entity.Payment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *installmentRepository) GetPaidAmount(ctx context.Context, scheduleID string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(
		`SELECT COALESCE(SUM(amount), 0) FROM payment WHERE schedule_id = ?`, scheduleID).Scan(&total).Error
	return total, err
}

func (r *installmentRepository) UpdateInstallmentStatus(ctx context.Context, scheduleID, status string) error {
	return r.db.WithContext(ctx).Model(&entity.InstallmentSchedule{}).
		Where("id = ?", scheduleID).Update("status", status).Error
}

func (r *installmentRepository) HasOverdueInstallments(ctx context.Context, loanID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.InstallmentSchedule{}).
		Where("loan_id = ? AND status = 'terlambat'", loanID).Count(&count).Error
	return count > 0, err
}

func (r *installmentRepository) AllInstallmentsPaid(ctx context.Context, loanID string) (bool, error) {
	var total, paid int64
	r.db.WithContext(ctx).Model(&entity.InstallmentSchedule{}).Where("loan_id = ?", loanID).Count(&total)
	r.db.WithContext(ctx).Model(&entity.InstallmentSchedule{}).Where("loan_id = ? AND status = 'lunas'", loanID).Count(&paid)
	return total > 0 && total == paid, nil
}
