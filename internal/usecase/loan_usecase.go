package usecase

import (
	"context"
	"errors"

	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/repository"
)

type LoanUsecase interface {
	List(ctx context.Context, coopID, status string) ([]model.LoanResponse, error)
	ListForMember(ctx context.Context, memberID string) ([]model.LoanResponse, error)
	GetByID(ctx context.Context, coopID, loanID string) (*model.LoanResponse, error)
	GetByIDForMember(ctx context.Context, memberID, loanID string) (*model.LoanResponse, error)
}

type loanUsecase struct{ repo repository.LoanRepository }

func NewLoanUsecase(repo repository.LoanRepository) LoanUsecase {
	return &loanUsecase{repo: repo}
}

func (u *loanUsecase) toResponses(loans []*repository.LoanWithMeta) []model.LoanResponse {
	result := make([]model.LoanResponse, len(loans))
	for i, l := range loans {
		result[i] = converter.ToLoanResponse(&l.Loan, l.MemberName, l.Outstanding, nil)
	}
	return result
}

func (u *loanUsecase) List(ctx context.Context, coopID, status string) ([]model.LoanResponse, error) {
	loans, err := u.repo.FindAll(ctx, coopID, status)
	if err != nil {
		return nil, err
	}
	return u.toResponses(loans), nil
}

func (u *loanUsecase) ListForMember(ctx context.Context, memberID string) ([]model.LoanResponse, error) {
	loans, err := u.repo.FindAllByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}
	return u.toResponses(loans), nil
}

func (u *loanUsecase) GetByID(ctx context.Context, coopID, loanID string) (*model.LoanResponse, error) {
	loan, err := u.repo.FindByID(ctx, coopID, loanID)
	if errors.Is(err, repository.ErrLoanNotFound) {
		return nil, errors.New("pinjaman tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	instResponses := make([]model.InstallmentResponse, len(loan.Schedule))
	for i, s := range loan.Schedule {
		instResponses[i] = converter.ToInstallmentResponse(&s.InstallmentSchedule, s.PaidAmount)
	}
	resp := converter.ToLoanResponse(&loan.Loan, loan.MemberName, loan.Outstanding, instResponses)
	return &resp, nil
}

func (u *loanUsecase) GetByIDForMember(ctx context.Context, memberID, loanID string) (*model.LoanResponse, error) {
	loan, err := u.repo.FindByID(ctx, "", loanID)
	if errors.Is(err, repository.ErrLoanNotFound) {
		return nil, errors.New("pinjaman tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	if loan.MemberID.String() != memberID {
		return nil, errors.New("forbidden")
	}
	instResponses := make([]model.InstallmentResponse, len(loan.Schedule))
	for i, s := range loan.Schedule {
		instResponses[i] = converter.ToInstallmentResponse(&s.InstallmentSchedule, s.PaidAmount)
	}
	resp := converter.ToLoanResponse(&loan.Loan, loan.MemberName, loan.Outstanding, instResponses)
	return &resp, nil
}
