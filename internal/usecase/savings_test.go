package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
	"github.com/wiradana/backend/internal/usecase"
)

// ── mock SavingsRepository ──────────────────────────────────────────────────

type mockSavingsRepo struct {
	pokokCount      int64
	sukarelaBalance int64
}

func (m *mockSavingsRepo) Create(_ context.Context, _ *entity.SavingsTransaction) error { return nil }
func (m *mockSavingsRepo) FindByMember(_ context.Context, _, _ string) ([]*entity.SavingsTransaction, error) {
	return nil, nil
}
func (m *mockSavingsRepo) CountPokok(_ context.Context, _, _ string) (int64, error) {
	return m.pokokCount, nil
}
func (m *mockSavingsRepo) GetSukarelaBalance(_ context.Context, _, _ string) (int64, error) {
	return m.sukarelaBalance, nil
}
func (m *mockSavingsRepo) GetCoopSummary(_ context.Context, _ string) (*repository.CoopSavingsSummaryRow, error) {
	return &repository.CoopSavingsSummaryRow{}, nil
}
func (m *mockSavingsRepo) FindAllRecent(_ context.Context, _, _, _ string, _, _ int) ([]*repository.SavingsTxWithMemberRow, int64, error) {
	return nil, 0, nil
}

// ── mock MemberRepository ───────────────────────────────────────────────────

type mockMemberRepo struct{}

func (m *mockMemberRepo) Create(_ context.Context, _ *entity.Member) error { return nil }
func (m *mockMemberRepo) FindByID(_ context.Context, _, _ string) (*entity.Member, error) {
	return &entity.Member{FullName: "Test Member", Status: "aktif"}, nil
}
func (m *mockMemberRepo) FindByNIK(_ context.Context, _, _ string) (*entity.Member, error) {
	return nil, repository.ErrMemberNotFound
}
func (m *mockMemberRepo) FindAll(_ context.Context, _, _, _ string) ([]*entity.Member, error) {
	return nil, nil
}
func (m *mockMemberRepo) Update(_ context.Context, _ *entity.Member) error { return nil }
func (m *mockMemberRepo) GetSavingsSummary(_ context.Context, _ string) (*model.SavingsSummary, error) {
	return &model.SavingsSummary{}, nil
}

// ── helpers ─────────────────────────────────────────────────────────────────

func newSavingsUC(pokokCount int64, sukarelaBalance int64) usecase.SavingsUsecase {
	return usecase.NewSavingsUsecase(
		&mockSavingsRepo{pokokCount: pokokCount, sukarelaBalance: sukarelaBalance},
		&mockMemberRepo{},
	)
}

const (
	testCoopID   = "00000000-0000-0000-0000-000000000001"
	testMemberID = "00000000-0000-0000-0000-000000000002"
)

// ── tests ────────────────────────────────────────────────────────────────────

func TestSavings_PokokDuplicate_Returns409Error(t *testing.T) {
	// CountPokok returns 1 (already recorded) → ErrPokokAlreadyRecorded
	uc := newSavingsUC(1, 0)
	_, err := uc.Record(context.Background(), testCoopID, testMemberID, &model.CreateSavingsRequest{
		SavingsType: "pokok",
		Direction:   "setor",
		Amount:      500_000,
	})
	if !errors.Is(err, usecase.ErrPokokAlreadyRecorded) {
		t.Errorf("want ErrPokokAlreadyRecorded, got %v", err)
	}
}

func TestSavings_TarikWajib_Returns409Error(t *testing.T) {
	uc := newSavingsUC(0, 0)
	_, err := uc.Record(context.Background(), testCoopID, testMemberID, &model.CreateSavingsRequest{
		SavingsType: "wajib",
		Direction:   "tarik",
		Amount:      100_000,
	})
	if !errors.Is(err, usecase.ErrCannotWithdrawMandatory) {
		t.Errorf("want ErrCannotWithdrawMandatory, got %v", err)
	}
}

func TestSavings_TarikPokok_Returns409Error(t *testing.T) {
	uc := newSavingsUC(0, 0)
	_, err := uc.Record(context.Background(), testCoopID, testMemberID, &model.CreateSavingsRequest{
		SavingsType: "pokok",
		Direction:   "tarik",
		Amount:      500_000,
	})
	if !errors.Is(err, usecase.ErrCannotWithdrawMandatory) {
		t.Errorf("want ErrCannotWithdrawMandatory, got %v", err)
	}
}

func TestSavings_TarikSukarela_InsufficientBalance(t *testing.T) {
	// balance = 50_000, tarik 100_000 → ErrInsufficientBalance
	uc := newSavingsUC(0, 50_000)
	_, err := uc.Record(context.Background(), testCoopID, testMemberID, &model.CreateSavingsRequest{
		SavingsType: "sukarela",
		Direction:   "tarik",
		Amount:      100_000,
	})
	if !errors.Is(err, usecase.ErrInsufficientBalance) {
		t.Errorf("want ErrInsufficientBalance, got %v", err)
	}
}

func TestSavings_Setor_OK(t *testing.T) {
	uc := newSavingsUC(0, 0)
	resp, err := uc.Record(context.Background(), testCoopID, testMemberID, &model.CreateSavingsRequest{
		SavingsType: "wajib",
		Direction:   "setor",
		Amount:      150_000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Amount != 150_000 {
		t.Errorf("want amount 150000, got %d", resp.Amount)
	}
	_ = time.Now() // suppress unused import
}
