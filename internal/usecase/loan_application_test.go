package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/gateway/adins"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
	"github.com/wiradana/backend/internal/usecase"
)

// ── mock LoanApplicationRepository ────────────────────────────────────────

type loanAppRepoMock struct {
	app *repository.LoanApplicationWithMeta
}

func (m *loanAppRepoMock) Create(_ context.Context, _ *entity.LoanApplication) error { return nil }
func (m *loanAppRepoMock) CreateAssessment(_ context.Context, _ *entity.CreditAssessment) error {
	return nil
}
func (m *loanAppRepoMock) FindByID(_ context.Context, _, _ string) (*repository.LoanApplicationWithMeta, error) {
	if m.app == nil {
		return nil, repository.ErrLoanApplicationNotFound
	}
	return m.app, nil
}
func (m *loanAppRepoMock) FindAll(_ context.Context, _, _ string) ([]*repository.LoanApplicationWithMeta, error) {
	return nil, nil
}
func (m *loanAppRepoMock) FindAllByMember(_ context.Context, _ string) ([]*repository.LoanApplicationWithMeta, error) {
	return nil, nil
}
func (m *loanAppRepoMock) UpdateStatus(_ context.Context, _, _ string, _ *uuid.UUID) error {
	return nil
}
func (m *loanAppRepoMock) GetTotalSavings(_ context.Context, _ string) (int64, error) {
	return 1_000_000, nil
}
func (m *loanAppRepoMock) GetKetepatanBayar(_ context.Context, _ string) (float64, error) {
	return 1.0, nil
}
func (m *loanAppRepoMock) GetKonsistensiSimpanan(_ context.Context, _ string, _ time.Time) (float64, error) {
	return 0.8, nil
}

// ── mock LoanConfigRepository ──────────────────────────────────────────────

type loanConfigRepoMock struct {
	plafond int64
	rate    float64
}

func (m *loanConfigRepoMock) FindByCoopID(_ context.Context, _ string) (*entity.LoanConfig, error) {
	return &entity.LoanConfig{
		MaxPlafond:      m.plafond,
		FlatRateMonthly: m.rate,
		PenaltyDaily:    5000,
	}, nil
}
func (m *loanConfigRepoMock) Upsert(_ context.Context, _ *entity.LoanConfig) error { return nil }

// ── mock MemberRepository (with controllable status) ──────────────────────

type memberRepoForApp struct {
	status string
}

func (m *memberRepoForApp) Create(_ context.Context, _ *entity.Member) error { return nil }
func (m *memberRepoForApp) FindByID(_ context.Context, _, _ string) (*entity.Member, error) {
	return &entity.Member{
		FullName: "Test Member",
		Status:   m.status,
		JoinedAt: time.Now().AddDate(0, -12, 0),
	}, nil
}
func (m *memberRepoForApp) FindByNIK(_ context.Context, _, _ string) (*entity.Member, error) {
	return nil, repository.ErrMemberNotFound
}
func (m *memberRepoForApp) FindAll(_ context.Context, _, _, _ string) ([]*entity.Member, error) {
	return nil, nil
}
func (m *memberRepoForApp) Update(_ context.Context, _ *entity.Member) error { return nil }
func (m *memberRepoForApp) GetSavingsSummary(_ context.Context, _ string) (*model.SavingsSummary, error) {
	return &model.SavingsSummary{}, nil
}

// ── mock ScoringGateway ────────────────────────────────────────────────────

type scoringGatewayMock struct{}

func (m *scoringGatewayMock) Score(_ context.Context, in adins.ScoringInput) (adins.ScoringResult, error) {
	return adins.ScoringResult{
		Score: 75, Grade: "B", Recommendation: "approve",
		LimitRekomendasi: in.JumlahDiajukan, Reasons: []string{"ok"},
		Source: "MOCK_ADINS_SCORING",
	}, nil
}

// ── mock LoanRepository (for app tests) ───────────────────────────────────

type loanRepoForApp struct{}

func (m *loanRepoForApp) Create(_ context.Context, _ *entity.Loan) error { return nil }
func (m *loanRepoForApp) BulkCreateSchedule(_ context.Context, _ []entity.InstallmentSchedule) error {
	return nil
}
func (m *loanRepoForApp) FindAll(_ context.Context, _, _ string) ([]*repository.LoanWithMeta, error) {
	return nil, nil
}
func (m *loanRepoForApp) FindAllByMember(_ context.Context, _ string) ([]*repository.LoanWithMeta, error) {
	return nil, nil
}
func (m *loanRepoForApp) FindByID(_ context.Context, _, _ string) (*repository.LoanWithMeta, error) {
	return nil, nil
}
func (m *loanRepoForApp) UpdateStatus(_ context.Context, _, _ string) error { return nil }

// ── helper ─────────────────────────────────────────────────────────────────

func newLoanAppUC(memberStatus string, plafond int64) usecase.LoanApplicationUsecase {
	return usecase.NewLoanApplicationUsecase(
		&loanAppRepoMock{},
		&loanConfigRepoMock{plafond: plafond, rate: 1.5},
		&memberRepoForApp{status: memberStatus},
		&loanRepoForApp{},
		&scoringGatewayMock{},
	)
}

// ── tests ──────────────────────────────────────────────────────────────────

func TestLoanApplication_MemberNonAktif_ReturnsError(t *testing.T) {
	uc := newLoanAppUC("nonaktif", 20_000_000)
	_, err := uc.Create(context.Background(), "coop-id", &model.CreateLoanApplicationRequest{
		MemberID:    "member-id",
		Amount:      5_000_000,
		TenorMonths: 12,
		Purpose:     "Test",
	})
	if !errors.Is(err, usecase.ErrMemberNotAktif) {
		t.Errorf("want ErrMemberNotAktif, got %v", err)
	}
}

func TestLoanApplication_AmountExceedsPlafond_ReturnsError(t *testing.T) {
	uc := newLoanAppUC("aktif", 10_000_000) // plafond 10jt
	_, err := uc.Create(context.Background(), "coop-id", &model.CreateLoanApplicationRequest{
		MemberID:    "member-id",
		Amount:      15_000_000, // melebihi plafond
		TenorMonths: 12,
		Purpose:     "Test",
	})
	if !errors.Is(err, usecase.ErrAmountExceedsPlafond) {
		t.Errorf("want ErrAmountExceedsPlafond, got %v", err)
	}
}

func TestLoanApplication_ValidRequest_ReturnsWithAssessment(t *testing.T) {
	uc := newLoanAppUC("aktif", 20_000_000)
	resp, err := uc.Create(context.Background(), "coop-id", &model.CreateLoanApplicationRequest{
		MemberID:    "member-id",
		Amount:      5_000_000,
		TenorMonths: 12,
		Purpose:     "Modal tani",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Assessment == nil {
		t.Fatal("want credit assessment embedded in response, got nil")
	}
	if resp.Assessment.Source != "MOCK_ADINS_SCORING" {
		t.Errorf("want source MOCK_ADINS_SCORING, got %q", resp.Assessment.Source)
	}
	if resp.Status != "pending" {
		t.Errorf("want status pending, got %q", resp.Status)
	}
}

func TestLoanApplication_ApproveNonPending_ReturnsError(t *testing.T) {
	// Application exists but status is "rejected"
	appID := uuid.New().String()
	approvedApp := &repository.LoanApplicationWithMeta{
		LoanApplication: entity.LoanApplication{
			Status: "rejected",
		},
	}
	appRepo := &loanAppRepoMock{app: approvedApp}
	uc := usecase.NewLoanApplicationUsecase(
		appRepo,
		&loanConfigRepoMock{plafond: 20_000_000, rate: 1.5},
		&memberRepoForApp{status: "aktif"},
		&loanRepoForApp{},
		&scoringGatewayMock{},
	)
	_, err := uc.Approve(context.Background(), "coop-id", appID, "user-id")
	if !errors.Is(err, usecase.ErrApplicationNotPending) {
		t.Errorf("want ErrApplicationNotPending, got %v", err)
	}
}
