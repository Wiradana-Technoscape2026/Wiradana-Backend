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
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/repository"
	"gorm.io/datatypes"
)

var ErrInstallmentNotFound = errors.New("angsuran tidak ditemukan")

type InstallmentUsecase interface {
	Pay(ctx context.Context, scheduleID string, req *model.PayInstallmentRequest, userID string) (*model.PayInstallmentResponse, error)
}

type installmentUsecase struct {
	instRepo  repository.InstallmentRepository
	loanRepo  repository.LoanRepository
	auditRepo repository.LoanAuditRepository
}

func NewInstallmentUsecase(instRepo repository.InstallmentRepository, loanRepo repository.LoanRepository, auditRepo repository.LoanAuditRepository) InstallmentUsecase {
	return &installmentUsecase{instRepo: instRepo, loanRepo: loanRepo, auditRepo: auditRepo}
}

func (u *installmentUsecase) Pay(ctx context.Context, scheduleID string, req *model.PayInstallmentRequest, userID string) (*model.PayInstallmentResponse, error) {
	inst, err := u.instRepo.FindByID(ctx, scheduleID)
	if errors.Is(err, repository.ErrInstallmentNotFound) {
		return nil, ErrInstallmentNotFound
	}
	if err != nil {
		return nil, err
	}

	loanID := inst.LoanID.String()
	loanBefore, _ := u.loanRepo.FindByID(ctx, "", loanID)

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


	overdue, _ := u.instRepo.HasOverdueInstallments(ctx, loanID)
	allPaid, _ := u.instRepo.AllInstallmentsPaid(ctx, loanID)

	loanStatus := "aktif"
	if allPaid {
		loanStatus = "lunas"
	} else if overdue {
		loanStatus = "menunggak"
	}
	_ = u.loanRepo.UpdateStatus(ctx, loanID, loanStatus)

	loanAfter, _ := u.loanRepo.FindByID(ctx, "", loanID)

	// Create audit log for payment
	var beforeJSON, afterJSON []byte
	var coopUUID uuid.UUID
	if loanBefore != nil {
		coopUUID = loanBefore.CooperativeID
		bm := map[string]interface{}{
			"status":      loanBefore.Status,
			"outstanding": loanBefore.Outstanding,
		}
		beforeJSON, _ = json.Marshal(bm)
	}
	if loanAfter != nil {
		coopUUID = loanAfter.CooperativeID
		am := map[string]interface{}{
			"status":      loanAfter.Status,
			"outstanding": loanAfter.Outstanding,
		}
		afterJSON, _ = json.Marshal(am)
	}

	userUUID, _ := uuid.Parse(userID)
	auditLog := &entity.LoanAuditLog{
		CooperativeID: coopUUID,
		LoanID:        inst.LoanID,
		Action:        "pay",
		PerformedBy:   userUUID,
		BeforeData:    datatypes.JSON(beforeJSON),
		AfterData:     datatypes.JSON(afterJSON),
		Note:          fmt.Sprintf("Pembayaran angsuran periode %d sebesar Rp %d", inst.PeriodNo, payment.Amount),
	}
	_ = u.auditRepo.CreateLog(ctx, auditLog)

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
