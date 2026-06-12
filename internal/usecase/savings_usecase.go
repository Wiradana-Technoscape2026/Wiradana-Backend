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
	ErrPokokAlreadyRecorded   = errors.New("simpanan pokok sudah pernah disetor")
	ErrCannotWithdrawMandatory = errors.New("simpanan pokok dan wajib tidak dapat ditarik")
	ErrInsufficientBalance    = errors.New("saldo sukarela tidak mencukupi")
)

type SavingsUsecase interface {
	Record(ctx context.Context, coopID, memberID string, req *model.CreateSavingsRequest) (*model.SavingsTransactionResponse, error)
	FindByMember(ctx context.Context, coopID, memberID string) ([]model.SavingsTransactionResponse, error)
}

type savingsUsecase struct {
	savingsRepo repository.SavingsRepository
	memberRepo  repository.MemberRepository
}

func NewSavingsUsecase(savingsRepo repository.SavingsRepository, memberRepo repository.MemberRepository) SavingsUsecase {
	return &savingsUsecase{savingsRepo: savingsRepo, memberRepo: memberRepo}
}

func (u *savingsUsecase) Record(ctx context.Context, coopID, memberID string, req *model.CreateSavingsRequest) (*model.SavingsTransactionResponse, error) {
	if _, err := u.memberRepo.FindByID(ctx, coopID, memberID); err != nil {
		if errors.Is(err, repository.ErrMemberNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	switch {
	case req.SavingsType == "pokok" && req.Direction == "setor":
		count, err := u.savingsRepo.CountPokok(ctx, memberID, coopID)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, ErrPokokAlreadyRecorded
		}

	case (req.SavingsType == "pokok" || req.SavingsType == "wajib") && req.Direction == "tarik":
		return nil, ErrCannotWithdrawMandatory

	case req.SavingsType == "sukarela" && req.Direction == "tarik":
		balance, err := u.savingsRepo.GetSukarelaBalance(ctx, memberID, coopID)
		if err != nil {
			return nil, err
		}
		if balance < req.Amount {
			return nil, ErrInsufficientBalance
		}
	}

	coopUUID, err := uuid.Parse(coopID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return nil, errors.New("member_id tidak valid")
	}

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

	resp := converter.ToSavingsResponse(tx)
	return &resp, nil
}

func (u *savingsUsecase) FindByMember(ctx context.Context, coopID, memberID string) ([]model.SavingsTransactionResponse, error) {
	if _, err := u.memberRepo.FindByID(ctx, coopID, memberID); err != nil {
		if errors.Is(err, repository.ErrMemberNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	txs, err := u.savingsRepo.FindByMember(ctx, coopID, memberID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.SavingsTransactionResponse, 0, len(txs))
	for _, tx := range txs {
		responses = append(responses, converter.ToSavingsResponse(tx))
	}
	return responses, nil
}
