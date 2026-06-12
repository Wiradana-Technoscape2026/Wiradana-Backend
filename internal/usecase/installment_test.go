package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
	"github.com/wiradana/backend/internal/usecase"
)

// ── mock InstallmentRepository ─────────────────────────────────────────────

type instRepoMock struct {
	inst       *repository.InstallmentWithPaid
	findErr    error
	paidAmount int64
	hasOverdue bool
	allPaid    bool
	// captured status from UpdateInstallmentStatus
	updatedStatus string
}

func (m *instRepoMock) FindByID(_ context.Context, _ string) (*repository.InstallmentWithPaid, error) {
	return m.inst, m.findErr
}
func (m *instRepoMock) FindByLoanID(_ context.Context, _ string) ([]repository.InstallmentWithPaid, error) {
	return nil, nil
}
func (m *instRepoMock) CreatePayment(_ context.Context, _ *entity.Payment) error { return nil }
func (m *instRepoMock) GetPaidAmount(_ context.Context, _ string) (int64, error) {
	return m.paidAmount, nil
}
func (m *instRepoMock) UpdateInstallmentStatus(_ context.Context, _ string, status string) error {
	m.updatedStatus = status
	return nil
}
func (m *instRepoMock) HasOverdueInstallments(_ context.Context, _ string) (bool, error) {
	return m.hasOverdue, nil
}
func (m *instRepoMock) AllInstallmentsPaid(_ context.Context, _ string) (bool, error) {
	return m.allPaid, nil
}

// ── mock LoanRepository (minimal, for installment tests) ───────────────────

type loanRepoMock struct {
	capturedLoanStatus string
}

func (m *loanRepoMock) Create(_ context.Context, _ *entity.Loan) error { return nil }
func (m *loanRepoMock) BulkCreateSchedule(_ context.Context, _ []entity.InstallmentSchedule) error {
	return nil
}
func (m *loanRepoMock) FindAll(_ context.Context, _, _ string) ([]*repository.LoanWithMeta, error) {
	return nil, nil
}
func (m *loanRepoMock) FindAllByMember(_ context.Context, _ string) ([]*repository.LoanWithMeta, error) {
	return nil, nil
}
func (m *loanRepoMock) FindByID(_ context.Context, _, _ string) (*repository.LoanWithMeta, error) {
	return nil, nil
}
func (m *loanRepoMock) UpdateStatus(_ context.Context, _ string, status string) error {
	m.capturedLoanStatus = status
	return nil
}

// ── helpers ────────────────────────────────────────────────────────────────

func makeInstallment(totalDue int64, dueDate time.Time, currentStatus string) *repository.InstallmentWithPaid {
	loanID := uuid.New()
	return &repository.InstallmentWithPaid{
		InstallmentSchedule: entity.InstallmentSchedule{
			ID:       uuid.New(),
			LoanID:   loanID,
			TotalDue: totalDue,
			DueDate:  dueDate,
			Status:   currentStatus,
		},
		PaidAmount: 0,
	}
}

func newInstUC(ir *instRepoMock, lr *loanRepoMock) usecase.InstallmentUsecase {
	return usecase.NewInstallmentUsecase(ir, lr)
}

// ── tests ──────────────────────────────────────────────────────────────────

func TestInstallment_FullPayment_StatusLunas(t *testing.T) {
	inst := makeInstallment(491666, time.Now().AddDate(0, 1, 0), "belum_bayar")
	ir := &instRepoMock{
		inst:       inst,
		paidAmount: 491666, // full amount paid after insert
		hasOverdue: false,
		allPaid:    false,
	}
	lr := &loanRepoMock{}
	uc := newInstUC(ir, lr)

	_, err := uc.Pay(context.Background(), inst.ID.String(), &model.PayInstallmentRequest{Amount: 491666})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ir.updatedStatus != "lunas" {
		t.Errorf("want installment status lunas, got %q", ir.updatedStatus)
	}
}

func TestInstallment_AllPaid_LoanStatusLunas(t *testing.T) {
	inst := makeInstallment(491666, time.Now().AddDate(0, 1, 0), "belum_bayar")
	ir := &instRepoMock{
		inst:       inst,
		paidAmount: 491666,
		hasOverdue: false,
		allPaid:    true, // all installments are lunas
	}
	lr := &loanRepoMock{}
	uc := newInstUC(ir, lr)

	_, err := uc.Pay(context.Background(), inst.ID.String(), &model.PayInstallmentRequest{Amount: 491666})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lr.capturedLoanStatus != "lunas" {
		t.Errorf("want loan status lunas, got %q", lr.capturedLoanStatus)
	}
}

func TestInstallment_OverdueExists_LoanStatusMenunggak(t *testing.T) {
	inst := makeInstallment(491666, time.Now().AddDate(0, 1, 0), "belum_bayar")
	ir := &instRepoMock{
		inst:       inst,
		paidAmount: 491666,
		hasOverdue: true, // another installment is terlambat
		allPaid:    false,
	}
	lr := &loanRepoMock{}
	uc := newInstUC(ir, lr)

	_, err := uc.Pay(context.Background(), inst.ID.String(), &model.PayInstallmentRequest{Amount: 491666})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lr.capturedLoanStatus != "menunggak" {
		t.Errorf("want loan status menunggak, got %q", lr.capturedLoanStatus)
	}
}

func TestInstallment_NotFound_ReturnsError(t *testing.T) {
	ir := &instRepoMock{findErr: repository.ErrInstallmentNotFound}
	lr := &loanRepoMock{}
	uc := newInstUC(ir, lr)

	_, err := uc.Pay(context.Background(), uuid.New().String(), &model.PayInstallmentRequest{Amount: 100000})
	if !errors.Is(err, usecase.ErrInstallmentNotFound) {
		t.Errorf("want ErrInstallmentNotFound, got %v", err)
	}
}
