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

// ── mock SyncRepository ──────────────────────────────────────────────────────

type mockSyncRepo struct {
	keys map[uuid.UUID]*entity.IdempotencyKey
}

func newMockSyncRepo() *mockSyncRepo {
	return &mockSyncRepo{keys: make(map[uuid.UUID]*entity.IdempotencyKey)}
}

func (m *mockSyncRepo) FindKey(_ context.Context, key uuid.UUID) (*entity.IdempotencyKey, error) {
	if k, ok := m.keys[key]; ok {
		return k, nil
	}
	return nil, repository.ErrIdempotencyKeyNotFound
}

func (m *mockSyncRepo) SaveKey(_ context.Context, k *entity.IdempotencyKey) error {
	m.keys[k.Key] = k
	return nil
}

func (m *mockSyncRepo) PullDelta(_ context.Context, _ uuid.UUID, since time.Time, _ string, _ *uuid.UUID) (*model.SyncPullResponse, error) {
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

// ── stub usecases ────────────────────────────────────────────────────────────

type stubMemberUC struct{}

func (s *stubMemberUC) Create(_ context.Context, _ string, req *model.CreateMemberRequest) (*model.MemberResponse, error) {
	id := uuid.New().String()
	return &model.MemberResponse{ID: id, FullName: req.FullName}, nil
}
func (s *stubMemberUC) FindByID(_ context.Context, _, _ string) (*model.MemberResponse, error) {
	return &model.MemberResponse{ID: uuid.New().String()}, nil
}
func (s *stubMemberUC) FindAll(_ context.Context, _, _, _ string) ([]model.MemberResponse, error) {
	return nil, nil
}
func (s *stubMemberUC) Update(_ context.Context, _, _ string, _ *model.UpdateMemberRequest) (*model.MemberResponse, error) {
	return nil, errors.New("member_id diperlukan")
}

type stubSavingsUC struct{}

func (s *stubSavingsUC) Record(_ context.Context, _, _ string, _ *model.CreateSavingsRequest) (*model.SavingsTransactionResponse, error) {
	return nil, errors.New("stub error")
}
func (s *stubSavingsUC) FindByMember(_ context.Context, _, _ string) ([]model.SavingsTransactionResponse, error) {
	return nil, nil
}
func (s *stubSavingsUC) GetCoopSummary(_ context.Context, _ string) (*model.SavingsSummaryResponse, error) {
	return &model.SavingsSummaryResponse{}, nil
}
func (s *stubSavingsUC) FindAllRecent(_ context.Context, _, _, _ string, _, _ int) ([]model.SavingsTransactionWithMemberResponse, int64, error) {
	return nil, 0, nil
}

type stubLoanAppUC struct{}

func (s *stubLoanAppUC) Create(_ context.Context, _ string, _ *model.CreateLoanApplicationRequest) (*model.LoanApplicationResponse, error) {
	return nil, errors.New("stub error")
}
func (s *stubLoanAppUC) CreateForMember(_ context.Context, _, _ string, _ *model.CreatePortalLoanApplicationRequest) (*model.LoanApplicationResponse, error) {
	return nil, errors.New("stub error")
}
func (s *stubLoanAppUC) List(_ context.Context, _, _ string) ([]model.LoanApplicationResponse, error) {
	return nil, nil
}
func (s *stubLoanAppUC) ListForMember(_ context.Context, _ string) ([]model.LoanApplicationResponse, error) {
	return nil, nil
}
func (s *stubLoanAppUC) Approve(_ context.Context, _, _, _ string) (*model.ApproveApplicationResponse, error) {
	return nil, errors.New("stub error")
}
func (s *stubLoanAppUC) Reject(_ context.Context, _, _ string) (*model.LoanApplicationResponse, error) {
	return nil, errors.New("stub error")
}

type stubInstallmentUC struct{}

func (s *stubInstallmentUC) Pay(_ context.Context, _ string, _ *model.PayInstallmentRequest) (*model.PayInstallmentResponse, error) {
	return nil, errors.New("stub error")
}

type stubLoanConfigUC struct{}

func (s *stubLoanConfigUC) Get(_ context.Context, _ string) (*model.LoanConfigResponse, error) {
	return &model.LoanConfigResponse{}, nil
}
func (s *stubLoanConfigUC) Update(_ context.Context, _ string, _ *model.UpdateLoanConfigRequest) (*model.LoanConfigResponse, error) {
	return &model.LoanConfigResponse{ID: uuid.New().String()}, nil
}

// ── helper ───────────────────────────────────────────────────────────────────

func newSyncUC(repo *mockSyncRepo) usecase.SyncUsecase {
	return usecase.NewSyncUsecase(
		repo,
		&stubMemberUC{},
		&stubSavingsUC{},
		&stubLoanAppUC{},
		&stubInstallmentUC{},
		&stubLoanConfigUC{},
		nil, // inventoryUC — not needed for these tests
	)
}

const (
	syncCoopID = "00000000-0000-0000-0000-000000000010"
	syncUserID = "00000000-0000-0000-0000-000000000011"
)

// ── tests ────────────────────────────────────────────────────────────────────

func TestSync_InvalidMutationID_Rejected(t *testing.T) {
	uc := newSyncUC(newMockSyncRepo())
	resp, err := uc.Push(context.Background(), syncCoopID, syncUserID, []model.MutationRequest{
		{ID: "not-a-uuid", Type: "create_member", Payload: map[string]interface{}{}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("want 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].Status != "rejected" {
		t.Errorf("want rejected, got %s", resp.Results[0].Status)
	}
}

func TestSync_UnknownMutationType_Rejected(t *testing.T) {
	uc := newSyncUC(newMockSyncRepo())
	mutID := uuid.New().String()
	resp, err := uc.Push(context.Background(), syncCoopID, syncUserID, []model.MutationRequest{
		{ID: mutID, Type: "unknown_type_xyz", Payload: map[string]interface{}{}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Results[0].Status != "rejected" {
		t.Errorf("want rejected, got %s", resp.Results[0].Status)
	}
}

func TestSync_DuplicateIdempotencyKey_ReturnsDuplicate(t *testing.T) {
	repo := newMockSyncRepo()
	existingID := uuid.New()
	resultID := uuid.New()

	// Pre-populate the key as already applied
	repo.keys[existingID] = &entity.IdempotencyKey{
		Key:      existingID,
		ResultID: &resultID,
		Status:   "applied",
	}

	uc := newSyncUC(repo)
	resp, err := uc.Push(context.Background(), syncCoopID, syncUserID, []model.MutationRequest{
		{ID: existingID.String(), Type: "create_member", Payload: map[string]interface{}{}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Results[0].Status != "duplicate" {
		t.Errorf("want duplicate, got %s", resp.Results[0].Status)
	}
	if resp.Results[0].ResultID == nil || *resp.Results[0].ResultID != resultID.String() {
		t.Errorf("want resultID %s, got %v", resultID, resp.Results[0].ResultID)
	}
}

func TestSync_MissingField_UpdateMember_Rejected(t *testing.T) {
	// update_member without member_id in payload → rejected
	uc := newSyncUC(newMockSyncRepo())
	mutID := uuid.New().String()
	resp, err := uc.Push(context.Background(), syncCoopID, syncUserID, []model.MutationRequest{
		{ID: mutID, Type: "update_member", Payload: map[string]interface{}{
			// member_id intentionally missing
			"full_name": "Test",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Results[0].Status != "rejected" {
		t.Errorf("want rejected for missing member_id, got %s", resp.Results[0].Status)
	}
}

func TestSync_Pull_ReturnsCursor(t *testing.T) {
	uc := newSyncUC(newMockSyncRepo())
	before := time.Now()
	resp, err := uc.Pull(context.Background(), syncCoopID, nil, "pengurus", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Cursor.Before(before) {
		t.Errorf("cursor should be >= before pull time")
	}
	// Members, etc. are empty slices (not nil) per mock
	if resp.Members == nil {
		t.Error("Members should be initialized (not nil)")
	}
}
