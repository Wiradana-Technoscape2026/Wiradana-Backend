package converter

import (
	"encoding/json"
	"time"

	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
)

func ToLoanConfigResponse(lc *entity.LoanConfig) model.LoanConfigResponse {
	return model.LoanConfigResponse{
		ID:              lc.ID.String(),
		FlatRateMonthly: lc.FlatRateMonthly,
		MaxPlafond:      lc.MaxPlafond,
		PenaltyDaily:    lc.PenaltyDaily,
	}
}

func ToCreditAssessmentResponse(ca *entity.CreditAssessment) *model.CreditAssessmentResponse {
	if ca == nil {
		return nil
	}
	var features map[string]float64
	var reasons []string
	_ = json.Unmarshal(ca.Features, &features)
	_ = json.Unmarshal(ca.Reasons, &reasons)
	if features == nil {
		features = map[string]float64{}
	}
	if reasons == nil {
		reasons = []string{}
	}
	return &model.CreditAssessmentResponse{
		ID:             ca.ID.String(),
		ApplicationID:  ca.ApplicationID.String(),
		Score:          ca.Score,
		Grade:          ca.Grade,
		Recommendation: ca.Recommendation,
		LimitSuggested: ca.LimitSuggested,
		Features:       features,
		Reasons:        reasons,
		Source:         ca.Source,
	}
}

func ToLoanApplicationResponse(app *entity.LoanApplication, memberName string, approvedByName *string, assessment *entity.CreditAssessment) model.LoanApplicationResponse {
	purpose := ""
	if app.Purpose != nil {
		purpose = *app.Purpose
	}
	var approvedBy *string
	if app.ApprovedBy != nil {
		s := app.ApprovedBy.String()
		approvedBy = &s
	}
	var approvedAt *string
	if app.ApprovedAt != nil {
		s := app.ApprovedAt.Format(time.RFC3339)
		approvedAt = &s
	}
	return model.LoanApplicationResponse{
		ID:             app.ID.String(),
		MemberID:       app.MemberID.String(),
		MemberName:     memberName,
		Amount:         app.Amount,
		TenorMonths:    app.TenorMonths,
		Purpose:        purpose,
		Status:         app.Status,
		ApprovedBy:     approvedBy,
		ApprovedByName: approvedByName,
		ApprovedAt:     approvedAt,
		CreatedAt:      app.CreatedAt.Format(time.RFC3339),
		Assessment:     ToCreditAssessmentResponse(assessment),
	}
}

func ToInstallmentResponse(inst *entity.InstallmentSchedule, paidAmount int64) model.InstallmentResponse {
	return model.InstallmentResponse{
		ID:           inst.ID.String(),
		LoanID:       inst.LoanID.String(),
		PeriodNo:     inst.PeriodNo,
		DueDate:      inst.DueDate.Format("2006-01-02"),
		PrincipalDue: inst.PrincipalDue,
		InterestDue:  inst.InterestDue,
		TotalDue:     inst.TotalDue,
		PaidAmount:   paidAmount,
		Status:       inst.Status,
	}
}

func ToLoanResponse(loan *entity.Loan, memberName string, outstanding int64, schedule []model.InstallmentResponse) model.LoanResponse {
	return model.LoanResponse{
		ID:              loan.ID.String(),
		ApplicationID:   loan.ApplicationID.String(),
		MemberID:        loan.MemberID.String(),
		MemberName:      memberName,
		Principal:       loan.Principal,
		FlatRateMonthly: loan.FlatRateMonthly,
		TenorMonths:     loan.TenorMonths,
		Status:          loan.Status,
		DisbursedAt:     loan.DisbursedAt.Format("2006-01-02"),
		Outstanding:     outstanding,
		Schedule:        schedule,
	}
}

func ToPaymentResponse(p *entity.Payment) model.PaymentResponse {
	return model.PaymentResponse{
		ID:         p.ID.String(),
		ScheduleID: p.ScheduleID.String(),
		Amount:     p.Amount,
		Penalty:    p.Penalty,
		PaidAt:     p.PaidAt.Format(time.RFC3339),
	}
}

func ToLoanAuditLogResponse(l *entity.LoanAuditLog, performedByEmail string) model.LoanAuditLogResponse {
	flaggedAt := ""
	if l.FlaggedAt != nil {
		flaggedAt = l.FlaggedAt.Format(time.RFC3339)
	}
	return model.LoanAuditLogResponse{
		ID:               l.ID.String(),
		LoanID:           l.LoanID.String(),
		Action:           l.Action,
		PerformedByEmail: performedByEmail,
		PerformedAt:      l.PerformedAt.Format(time.RFC3339),
		BeforeData:       string(l.BeforeData),
		AfterData:        string(l.AfterData),
		Note:             l.Note,
		IsFlagged:        l.IsFlagged,
		FlaggedByName:    l.FlaggedByName,
		FlaggedAt:        flaggedAt,
		FlaggedReason:    l.FlaggedReason,
	}
}

