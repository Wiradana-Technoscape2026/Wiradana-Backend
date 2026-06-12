package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
	"github.com/wiradana/backend/internal/usecase"
)

// ── Stub implementations ────────────────────────────────────────────────────

type stubSyncRepo struct {
	keys map[uuid.UUID]*entity.IdempotencyKey
}

func newStubSyncRepo() *stubSyncRepo {
	return &stubSyncRepo{keys: map[uuid.UUID]*entity.IdempotencyKey{}}
}

func (s *stubSyncRepo) FindKey(_ context.Context, key uuid.UUID) (*entity.IdempotencyKey, error) {
	k, ok := s.keys[key]
	if !ok {
		return nil, repository.ErrIdempotencyKeyNotFound
	}
	return k, nil
}

func (s *stubSyncRepo) SaveKey(_ context.Context, k *entity.IdempotencyKey) error {
	s.keys[k.Key] = k
	return nil
}

func (s *stubSyncRepo) PullDelta(_ context.Context, _ uuid.UUID, _ time.Time, _ string, _ *uuid.UUID) (*model.SyncPullResponse, error) {
	return &model.SyncPullResponse{
		Cursor:               time.Now(),
		Members:              []model.MemberResponse{},
		SavingsTransactions:  []model.SavingsTransactionResponse{},
		LoanApplications:     []model.LoanApplicationResponse{},
		Loans:                []model.LoanResponse{},
		InstallmentSchedules: []model.InstallmentResponse{},
		Payments:             []model.PaymentResponse{},
		ShuPeriods:           []model.ShuPeriodResponse{},
		ShuDistributions:     []model.ShuDistributionResponse{},
	}, nil
}

type stubSavingsUC struct{}

func (s *stubSavingsUC) Record(_ context.Context, _, _ string, _ *model.CreateSavingsRequest) (*model.SavingsTransactionResponse, error) {
	id := uuid.New().String()
	return &model.SavingsTransactionResponse{ID: id}, nil
}

func (s *stubSavingsUC) FindByMember(_ context.Context, _, _ string) ([]model.SavingsTransactionResponse, error) {
	return nil, nil
}

// ── Tests ───────────────────────────────────────────────────────────────────


func TestPush_Duplicate(t *testing.T) {
	repo := newStubSyncRepo()
	uc := usecase.NewSyncUsecase(repo, nil, &stubSavingsUC{}, nil, nil, nil)

	coopID := uuid.New().String()
	userID := uuid.New().String()
	mutID := uuid.New().String()

	mutations := []model.MutationRequest{{
		ID:   mutID,
		Type: "create_savings_transaction",
		Payload: map[string]interface{}{
			"member_id":    uuid.New().String(),
			"savings_type": "wajib",
			"direction":    "setor",
			"amount":       float64(100000),
		},
	}}

	// First push — should be applied
	resp1, err := uc.Push(context.Background(), coopID, userID, mutations)
	if err != nil {
		t.Fatalf("push 1 error: %v", err)
	}
	if resp1.Results[0].Status != "applied" {
		t.Fatalf("first push: want applied, got %s", resp1.Results[0].Status)
	}

	// Second push with same idempotency key — must be duplicate
	resp2, err := uc.Push(context.Background(), coopID, userID, mutations)
	if err != nil {
		t.Fatalf("push 2 error: %v", err)
	}
	if resp2.Results[0].Status != "duplicate" {
		t.Fatalf("second push: want duplicate, got %s", resp2.Results[0].Status)
	}
	if resp2.Results[0].ResultID == nil {
		t.Fatal("second push: expected result_id in duplicate response")
	}
	if resp1.Results[0].ResultID == nil || *resp1.Results[0].ResultID != *resp2.Results[0].ResultID {
		t.Errorf("duplicate result_id should match first push: %v != %v", resp1.Results[0].ResultID, resp2.Results[0].ResultID)
	}
}

func TestPush_UnknownMutationType(t *testing.T) {
	repo := newStubSyncRepo()
	uc := usecase.NewSyncUsecase(repo, nil, &stubSavingsUC{}, nil, nil, nil)

	resp, err := uc.Push(context.Background(), uuid.New().String(), uuid.New().String(), []model.MutationRequest{{
		ID:      uuid.New().String(),
		Type:    "delete_everything",
		Payload: map[string]interface{}{},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Results[0].Status != "rejected" {
		t.Fatalf("want rejected, got %s", resp.Results[0].Status)
	}
	if resp.Results[0].Error == nil {
		t.Fatal("expected error message for unknown type")
	}
}

func TestPush_MissingMemberID(t *testing.T) {
	repo := newStubSyncRepo()
	uc := usecase.NewSyncUsecase(repo, nil, &stubSavingsUC{}, nil, nil, nil)

	resp, err := uc.Push(context.Background(), uuid.New().String(), uuid.New().String(), []model.MutationRequest{{
		ID:      uuid.New().String(),
		Type:    "create_savings_transaction",
		Payload: map[string]interface{}{"amount": float64(100000)}, // missing member_id
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Results[0].Status != "rejected" {
		t.Fatalf("want rejected, got %s", resp.Results[0].Status)
	}
}

func TestPush_InvalidMutationID(t *testing.T) {
	repo := newStubSyncRepo()
	uc := usecase.NewSyncUsecase(repo, nil, &stubSavingsUC{}, nil, nil, nil)

	resp, err := uc.Push(context.Background(), uuid.New().String(), uuid.New().String(), []model.MutationRequest{{
		ID:      "not-a-uuid",
		Type:    "create_savings_transaction",
		Payload: map[string]interface{}{},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Results[0].Status != "rejected" {
		t.Fatalf("want rejected, got %s", resp.Results[0].Status)
	}
}

func TestPull_ReturnsCursor(t *testing.T) {
	repo := newStubSyncRepo()
	uc := usecase.NewSyncUsecase(repo, nil, &stubSavingsUC{}, nil, nil, nil)

	resp, err := uc.Pull(context.Background(), uuid.New().String(), nil, "pengurus", nil)
	if err != nil {
		t.Fatalf("pull error: %v", err)
	}
	if resp.Cursor.IsZero() {
		t.Fatal("cursor should not be zero")
	}
}
