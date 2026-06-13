package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
)

type SyncUsecase interface {
	Pull(ctx context.Context, coopID string, since *time.Time, role string, memberID *string) (*model.SyncPullResponse, error)
	Push(ctx context.Context, coopID string, userID string, mutations []model.MutationRequest) (*model.SyncPushResponse, error)
}

type syncUsecase struct {
	syncRepo      repository.SyncRepository
	memberUC      MemberUsecase
	savingsUC     SavingsUsecase
	loanAppUC     LoanApplicationUsecase
	installmentUC InstallmentUsecase
	loanConfigUC  LoanConfigUsecase
	inventoryUC   InventoryUsecase
}

func NewSyncUsecase(
	syncRepo repository.SyncRepository,
	memberUC MemberUsecase,
	savingsUC SavingsUsecase,
	loanAppUC LoanApplicationUsecase,
	installmentUC InstallmentUsecase,
	loanConfigUC LoanConfigUsecase,
	inventoryUC InventoryUsecase,
) SyncUsecase {
	return &syncUsecase{
		syncRepo:      syncRepo,
		memberUC:      memberUC,
		savingsUC:     savingsUC,
		loanAppUC:     loanAppUC,
		installmentUC: installmentUC,
		loanConfigUC:  loanConfigUC,
		inventoryUC:   inventoryUC,
	}
}

func (u *syncUsecase) Pull(ctx context.Context, coopID string, since *time.Time, role string, memberID *string) (*model.SyncPullResponse, error) {
	coopUUID, err := uuid.Parse(coopID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}

	sinceTime := time.Time{}
	if since != nil {
		sinceTime = *since
	}

	var memberUUID *uuid.UUID
	if memberID != nil && *memberID != "" {
		parsed, err := uuid.Parse(*memberID)
		if err == nil {
			memberUUID = &parsed
		}
	}

	return u.syncRepo.PullDelta(ctx, coopUUID, sinceTime, role, memberUUID)
}

func (u *syncUsecase) Push(ctx context.Context, coopID, userID string, mutations []model.MutationRequest) (*model.SyncPushResponse, error) {
	coopUUID, err := uuid.Parse(coopID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("user_id tidak valid")
	}

	results := make([]model.MutationResult, len(mutations))

	for i, mut := range mutations {
		keyUUID, err := uuid.Parse(mut.ID)
		if err != nil {
			errMsg := "mutation id bukan UUID valid"
			results[i] = model.MutationResult{ID: mut.ID, Status: "rejected", Error: &errMsg}
			continue
		}

		// Check idempotency
		existing, err := u.syncRepo.FindKey(ctx, keyUUID)
		if err == nil && existing != nil {
			results[i] = model.MutationResult{
				ID:       mut.ID,
				Status:   "duplicate",
				ResultID: uuidPtrToStr(existing.ResultID),
			}
			continue
		}

		// Dispatch mutation
		now := time.Now()
		resultID, dispatchErr := u.dispatchMutation(ctx, coopID, userID, mut)

		var keyEntry *entity.IdempotencyKey
		if dispatchErr != nil {
			errMsg := dispatchErr.Error()
			results[i] = model.MutationResult{ID: mut.ID, Status: "rejected", Error: &errMsg}
			keyEntry = &entity.IdempotencyKey{
				Key:           keyUUID,
				CooperativeID: coopUUID,
				UserID:        userUUID,
				MutationType:  mut.Type,
				Status:        "rejected",
				ErrorMessage:  &errMsg,
				ProcessedAt:   &now,
			}
		} else {
			results[i] = model.MutationResult{ID: mut.ID, Status: "applied", ResultID: &resultID}
			resultUUID, _ := uuid.Parse(resultID)
			keyEntry = &entity.IdempotencyKey{
				Key:           keyUUID,
				CooperativeID: coopUUID,
				UserID:        userUUID,
				MutationType:  mut.Type,
				ResultID:      &resultUUID,
				Status:        "applied",
				ProcessedAt:   &now,
			}
		}
		_ = u.syncRepo.SaveKey(ctx, keyEntry)
	}

	return &model.SyncPushResponse{Results: results}, nil
}

func (u *syncUsecase) dispatchMutation(ctx context.Context, coopID, userID string, mut model.MutationRequest) (string, error) {
	switch mut.Type {
	case "create_member":
		return u.dispatchCreateMember(ctx, coopID, mut.Payload)
	case "update_member":
		return u.dispatchUpdateMember(ctx, coopID, mut.Payload)
	case "create_savings_transaction":
		return u.dispatchCreateSavings(ctx, coopID, mut.Payload)
	case "create_loan_application":
		return u.dispatchCreateLoanApplication(ctx, coopID, mut.Payload)
	case "approve_loan_application":
		return u.dispatchApproveLoanApplication(ctx, coopID, userID, mut.Payload)
	case "reject_loan_application":
		return u.dispatchRejectLoanApplication(ctx, coopID, mut.Payload)
	case "create_payment":
		return u.dispatchCreatePayment(ctx, userID, mut.Payload)
	case "update_loan_config":
		return u.dispatchUpdateLoanConfig(ctx, coopID, mut.Payload)
	case "create_inventory_product":
		return u.dispatchCreateInventoryProduct(ctx, coopID, mut.Payload)
	case "update_inventory_product":
		return u.dispatchUpdateInventoryProduct(ctx, coopID, mut.Payload)
	case "record_inventory_movement":
		return u.dispatchRecordInventoryMovement(ctx, coopID, mut.Payload)
	default:
		return "", fmt.Errorf("unknown mutation type: %s", mut.Type)
	}
}

func (u *syncUsecase) dispatchCreateMember(ctx context.Context, coopID string, payload map[string]interface{}) (string, error) {
	var req model.CreateMemberRequest
	if err := remarshal(payload, &req); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}
	resp, err := u.memberUC.Create(ctx, coopID, &req)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (u *syncUsecase) dispatchUpdateMember(ctx context.Context, coopID string, payload map[string]interface{}) (string, error) {
	memberID, ok := strField(payload, "member_id")
	if !ok {
		return "", errors.New("member_id diperlukan")
	}
	var req model.UpdateMemberRequest
	if err := remarshal(payload, &req); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}
	resp, err := u.memberUC.Update(ctx, coopID, memberID, &req)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (u *syncUsecase) dispatchCreateSavings(ctx context.Context, coopID string, payload map[string]interface{}) (string, error) {
	memberID, ok := strField(payload, "member_id")
	if !ok {
		return "", errors.New("member_id diperlukan")
	}
	var req model.CreateSavingsRequest
	if err := remarshal(payload, &req); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}
	resp, err := u.savingsUC.Record(ctx, coopID, memberID, &req)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (u *syncUsecase) dispatchCreateLoanApplication(ctx context.Context, coopID string, payload map[string]interface{}) (string, error) {
	var req model.CreateLoanApplicationRequest
	if err := remarshal(payload, &req); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}
	resp, err := u.loanAppUC.Create(ctx, coopID, &req)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (u *syncUsecase) dispatchApproveLoanApplication(ctx context.Context, coopID, userID string, payload map[string]interface{}) (string, error) {
	appID, ok := strField(payload, "application_id")
	if !ok {
		return "", errors.New("application_id diperlukan")
	}
	resp, err := u.loanAppUC.Approve(ctx, coopID, appID, userID)
	if err != nil {
		return "", err
	}
	return resp.Application.ID, nil
}

func (u *syncUsecase) dispatchRejectLoanApplication(ctx context.Context, coopID string, payload map[string]interface{}) (string, error) {
	appID, ok := strField(payload, "application_id")
	if !ok {
		return "", errors.New("application_id diperlukan")
	}
	resp, err := u.loanAppUC.Reject(ctx, coopID, appID)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (u *syncUsecase) dispatchCreatePayment(ctx context.Context, userID string, payload map[string]interface{}) (string, error) {
	scheduleID, ok := strField(payload, "installment_id")
	if !ok {
		return "", errors.New("installment_id diperlukan")
	}
	var req model.PayInstallmentRequest
	if err := remarshal(payload, &req); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}
	resp, err := u.installmentUC.Pay(ctx, scheduleID, &req, userID)
	if err != nil {
		return "", err
	}
	return resp.Payment.ID, nil
}

func (u *syncUsecase) dispatchUpdateLoanConfig(ctx context.Context, coopID string, payload map[string]interface{}) (string, error) {
	var req model.UpdateLoanConfigRequest
	if err := remarshal(payload, &req); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}
	resp, err := u.loanConfigUC.Update(ctx, coopID, &req)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (u *syncUsecase) dispatchCreateInventoryProduct(ctx context.Context, coopID string, payload map[string]interface{}) (string, error) {
	var req model.CreateInventoryProductRequest
	if err := remarshal(payload, &req); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}
	resp, err := u.inventoryUC.CreateProduct(ctx, coopID, &req)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (u *syncUsecase) dispatchUpdateInventoryProduct(ctx context.Context, coopID string, payload map[string]interface{}) (string, error) {
	productID, ok := strField(payload, "product_id")
	if !ok {
		return "", errors.New("product_id diperlukan")
	}
	var req model.UpdateInventoryProductRequest
	if err := remarshal(payload, &req); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}
	resp, err := u.inventoryUC.UpdateProduct(ctx, coopID, productID, &req)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (u *syncUsecase) dispatchRecordInventoryMovement(ctx context.Context, coopID string, payload map[string]interface{}) (string, error) {
	productID, ok := strField(payload, "product_id")
	if !ok {
		return "", errors.New("product_id diperlukan")
	}
	var req model.CreateInventoryMovementRequest
	if err := remarshal(payload, &req); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}
	resp, err := u.inventoryUC.RecordMovement(ctx, coopID, productID, &req)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// remarshal JSON-encodes a map then decodes into the target struct.
func remarshal(src map[string]interface{}, dst interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

func strField(m map[string]interface{}, key string) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func uuidPtrToStr(u *uuid.UUID) *string {
	if u == nil {
		return nil
	}
	s := u.String()
	return &s
}
