package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/repository"
)

var ErrInstallmentNotFound = errors.New("angsuran tidak ditemukan")

type InstallmentUsecase interface {
	Pay(ctx context.Context, scheduleID string, req *model.PayInstallmentRequest) (*model.PayInstallmentResponse, error)
}

type installmentUsecase struct {
	instRepo repository.InstallmentRepository
	loanRepo repository.LoanRepository
}

func NewInstallmentUsecase(instRepo repository.InstallmentRepository, loanRepo repository.LoanRepository) InstallmentUsecase {
	return &installmentUsecase{instRepo: instRepo, loanRepo: loanRepo}
}

func (u *installmentUsecase) Pay(ctx context.Context, scheduleID string, req *model.PayInstallmentRequest) (*model.PayInstallmentResponse, error) {
	inst, err := u.instRepo.FindByID(ctx, scheduleID)
	if errors.Is(err, repository.ErrInstallmentNotFound) {
		return nil, ErrInstallmentNotFound
	}
	if err != nil {
		return nil, err
	}

	schedUUID, _ := uuid.Parse(scheduleID)
	payment := &entity.Payment{
		ScheduleID: schedUUID,
		Amount:     req.Amount,
		Penalty:    req.Penalty,
		PaidAt:     time.Now(),
	}
	if err := u.instRepo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	paidAmount, _ := u.instRepo.GetPaidAmount(ctx, scheduleID)
	newStatus := inst.Status
	if paidAmount >= inst.TotalDue {
		newStatus = "lunas"
	} else if time.Now().After(inst.DueDate) && paidAmount > 0 {
		newStatus = "terlambat"
	}
	if newStatus != inst.Status {
		_ = u.instRepo.UpdateInstallmentStatus(ctx, scheduleID, newStatus)
		inst.Status = newStatus
	}

	loanID := inst.LoanID.String()
	overdue, _ := u.instRepo.HasOverdueInstallments(ctx, loanID)
	allPaid, _ := u.instRepo.AllInstallmentsPaid(ctx, loanID)

	loanStatus := "aktif"
	if allPaid {
		loanStatus = "lunas"
	} else if overdue {
		loanStatus = "menunggak"
	}
	_ = u.loanRepo.UpdateStatus(ctx, loanID, loanStatus)

	loanWithMeta, _ := u.loanRepo.FindByID(ctx, "", loanID)
	loanInstResponses := make([]model.InstallmentResponse, 0)
	if loanWithMeta != nil {
		loanInstResponses = make([]model.InstallmentResponse, len(loanWithMeta.Schedule))
		for i, s := range loanWithMeta.Schedule {
			loanInstResponses[i] = converter.ToInstallmentResponse(&s.InstallmentSchedule, s.PaidAmount)
		}
	}

	instResp := converter.ToInstallmentResponse(&inst.InstallmentSchedule, paidAmount)
	payResp := converter.ToPaymentResponse(payment)
	loanResp := model.LoanResponse{}
	if loanWithMeta != nil {
		loanResp = converter.ToLoanResponse(&loanWithMeta.Loan, loanWithMeta.MemberName, loanWithMeta.Outstanding, loanInstResponses)
	}

	return &model.PayInstallmentResponse{
		Payment:     payResp,
		Installment: instResp,
		Loan:        loanResp,
	}, nil
}
