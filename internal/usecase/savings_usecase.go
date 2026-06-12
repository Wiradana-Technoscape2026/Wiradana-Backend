package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/repository"
)

var (
	ErrPokokAlreadyRecorded  = errors.New("simpanan pokok sudah pernah disetor")
	ErrCannotWithdrawMandatory = errors.New("simpanan pokok dan wajib tidak dapat ditarik")
	ErrInsufficientSukarela  = errors.New("saldo sukarela tidak mencukupi untuk penarikan")
)

type SavingsUsecase interface {
	Record(ctx context.Context, coopID, memberID string, req *model.CreateSavingsRequest) (*model.SavingsTransactionResponse, error)
	ListByMember(ctx context.Context, memberID string) ([]model.SavingsTransactionResponse, error)
}

type savingsUsecase struct {
	savingsRepo repository.SavingsRepository
}

func NewSavingsUsecase(savingsRepo repository.SavingsRepository) SavingsUsecase {
	return &savingsUsecase{savingsRepo: savingsRepo}
}

func (u *savingsUsecase) Record(ctx context.Context, coopID, memberID string, req *model.CreateSavingsRequest) (*model.SavingsTransactionResponse, error) {
	// Rule: pokok hanya boleh setor sekali
	if req.SavingsType == "pokok" && req.Direction == "setor" {
		count, _ := u.savingsRepo.CountPokokSetoran(ctx, memberID)
		if count > 0 {
			return nil, ErrPokokAlreadyRecorded
		}
	}

	// Rule: pokok & wajib tidak bisa ditarik
	if (req.SavingsType == "pokok" || req.SavingsType == "wajib") && req.Direction == "tarik" {
		return nil, ErrCannotWithdrawMandatory
	}

	// Rule: sukarela tarik tidak boleh melebihi saldo
	if req.SavingsType == "sukarela" && req.Direction == "tarik" {
		saldo, _ := u.savingsRepo.GetSukarelaSaldo(ctx, memberID)
		if req.Amount > saldo {
			return nil, ErrInsufficientSukarela
		}
	}

	coopUUID, _ := uuid.Parse(coopID)
	memberUUID, _ := uuid.Parse(memberID)

	tx := &entity.SavingsTransaction{
		CooperativeID: coopUUID,
		MemberID:      memberUUID,
		SavingsType:   req.SavingsType,
		Direction:     req.Direction,
		Amount:        req.Amount,
	}

	if err := u.savingsRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	resp := converter.ToSavingsTransactionResponse(tx)
	return &resp, nil
}

func (u *savingsUsecase) ListByMember(ctx context.Context, memberID string) ([]model.SavingsTransactionResponse, error) {
	txs, err := u.savingsRepo.FindByMemberID(ctx, memberID)
	if err != nil {
		return nil, err
	}
	result := make([]model.SavingsTransactionResponse, len(txs))
	for i, tx := range txs {
		result[i] = converter.ToSavingsTransactionResponse(&tx)
	}
	return result, nil
}
