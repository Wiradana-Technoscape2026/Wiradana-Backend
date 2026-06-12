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
<<<<<<< HEAD
	ErrPokokAlreadyRecorded  = errors.New("simpanan pokok sudah pernah disetor")
	ErrCannotWithdrawMandatory = errors.New("simpanan pokok dan wajib tidak dapat ditarik")
	ErrInsufficientSukarela  = errors.New("saldo sukarela tidak mencukupi untuk penarikan")
=======
	ErrPokokAlreadyRecorded   = errors.New("simpanan pokok sudah pernah disetor")
	ErrCannotWithdrawMandatory = errors.New("simpanan pokok dan wajib tidak dapat ditarik")
	ErrInsufficientBalance    = errors.New("saldo sukarela tidak mencukupi")
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
)

type SavingsUsecase interface {
	Record(ctx context.Context, coopID, memberID string, req *model.CreateSavingsRequest) (*model.SavingsTransactionResponse, error)
<<<<<<< HEAD
	ListByMember(ctx context.Context, memberID string) ([]model.SavingsTransactionResponse, error)
=======
	FindByMember(ctx context.Context, coopID, memberID string) ([]model.SavingsTransactionResponse, error)
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
}

type savingsUsecase struct {
	savingsRepo repository.SavingsRepository
<<<<<<< HEAD
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
=======
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
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed

	tx := &entity.SavingsTransaction{
		CooperativeID: coopUUID,
		MemberID:      memberUUID,
		SavingsType:   req.SavingsType,
		Direction:     req.Direction,
		Amount:        req.Amount,
	}
<<<<<<< HEAD

=======
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
	if err := u.savingsRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

<<<<<<< HEAD
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
=======
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
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
}
