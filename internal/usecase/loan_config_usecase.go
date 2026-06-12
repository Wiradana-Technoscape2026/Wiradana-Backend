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

type LoanConfigUsecase interface {
	Get(ctx context.Context, coopID string) (*model.LoanConfigResponse, error)
	Update(ctx context.Context, coopID string, req *model.UpdateLoanConfigRequest) (*model.LoanConfigResponse, error)
}

type loanConfigUsecase struct {
	repo repository.LoanConfigRepository
}

func NewLoanConfigUsecase(repo repository.LoanConfigRepository) LoanConfigUsecase {
	return &loanConfigUsecase{repo: repo}
}

func (u *loanConfigUsecase) Get(ctx context.Context, coopID string) (*model.LoanConfigResponse, error) {
	lc, err := u.repo.FindByCoopID(ctx, coopID)
	if errors.Is(err, repository.ErrLoanConfigNotFound) {
		coopUUID, _ := uuid.Parse(coopID)
		lc = &entity.LoanConfig{
			CooperativeID:   coopUUID,
			FlatRateMonthly: 1.5,
			MaxPlafond:      20_000_000,
			PenaltyDaily:    5000,
		}
		if err2 := u.repo.Upsert(ctx, lc); err2 != nil {
			return nil, err2
		}
	} else if err != nil {
		return nil, err
	}
	resp := converter.ToLoanConfigResponse(lc)
	return &resp, nil
}

func (u *loanConfigUsecase) Update(ctx context.Context, coopID string, req *model.UpdateLoanConfigRequest) (*model.LoanConfigResponse, error) {
	lc, err := u.repo.FindByCoopID(ctx, coopID)
	if errors.Is(err, repository.ErrLoanConfigNotFound) {
		coopUUID, _ := uuid.Parse(coopID)
		lc = &entity.LoanConfig{CooperativeID: coopUUID, FlatRateMonthly: 1.5, MaxPlafond: 20_000_000, PenaltyDaily: 5000}
	} else if err != nil {
		return nil, err
	}
	if req.FlatRateMonthly != nil {
		lc.FlatRateMonthly = *req.FlatRateMonthly
	}
	if req.MaxPlafond != nil {
		lc.MaxPlafond = *req.MaxPlafond
	}
	if req.PenaltyDaily != nil {
		lc.PenaltyDaily = *req.PenaltyDaily
	}
	if err := u.repo.Upsert(ctx, lc); err != nil {
		return nil, err
	}
	resp := converter.ToLoanConfigResponse(lc)
	return &resp, nil
}
