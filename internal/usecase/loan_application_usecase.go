package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/gateway/adins"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/repository"
	"gorm.io/datatypes"
)

var (
	ErrMemberNotAktif        = errors.New("anggota tidak aktif")
	ErrAmountExceedsPlafond  = errors.New("jumlah melebihi batas maksimal pinjaman")
	ErrApplicationNotPending = errors.New("pengajuan bukan dalam status pending")
	ErrLoanApplicationNotFound = errors.New("pengajuan tidak ditemukan")
)

type LoanApplicationUsecase interface {
	Create(ctx context.Context, coopID string, req *model.CreateLoanApplicationRequest) (*model.LoanApplicationResponse, error)
	CreateForMember(ctx context.Context, memberID, coopID string, req *model.CreatePortalLoanApplicationRequest) (*model.LoanApplicationResponse, error)
	List(ctx context.Context, coopID, status string) ([]model.LoanApplicationResponse, error)
	ListForMember(ctx context.Context, memberID string) ([]model.LoanApplicationResponse, error)
	Approve(ctx context.Context, coopID, appID, approvedByUserID string) (*model.ApproveApplicationResponse, error)
	Reject(ctx context.Context, coopID, appID string) (*model.LoanApplicationResponse, error)
}

type loanApplicationUsecase struct {
	appRepo    repository.LoanApplicationRepository
	configRepo repository.LoanConfigRepository
	memberRepo repository.MemberRepository
	loanRepo   repository.LoanRepository
	scoring    adins.ScoringGateway
	auditRepo  repository.LoanAuditRepository
}

func NewLoanApplicationUsecase(
	appRepo repository.LoanApplicationRepository,
	configRepo repository.LoanConfigRepository,
	memberRepo repository.MemberRepository,
	loanRepo repository.LoanRepository,
	scoring adins.ScoringGateway,
	auditRepo repository.LoanAuditRepository,
) LoanApplicationUsecase {
	return &loanApplicationUsecase{
		appRepo: appRepo, configRepo: configRepo,
		memberRepo: memberRepo, loanRepo: loanRepo, scoring: scoring,
		auditRepo: auditRepo,
	}
}

func (u *loanApplicationUsecase) buildAndScore(ctx context.Context, coopID, memberID string, amount int64, tenor int) (*entity.CreditAssessment, error) {
	member, err := u.memberRepo.FindByID(ctx, coopID, memberID)
	if err != nil {
		return nil, ErrMemberNotFound
	}
	if member.Status != "aktif" {
		return nil, ErrMemberNotAktif
	}

	cfg, _ := u.configRepo.FindByCoopID(ctx, coopID)
	maxPlafond := int64(20_000_000)
	rate := 1.5
	if cfg != nil {
		maxPlafond = cfg.MaxPlafond
		rate = cfg.FlatRateMonthly
	}
	if amount > maxPlafond {
		return nil, ErrAmountExceedsPlafond
	}

	totalSimpanan, _ := u.appRepo.GetTotalSavings(ctx, memberID)
	ketepatan, _ := u.appRepo.GetKetepatanBayar(ctx, memberID)
	konsistensi, _ := u.appRepo.GetKonsistensiSimpanan(ctx, memberID, member.JoinedAt)
	hasPrior, _ := u.appRepo.HasPriorLoans(ctx, memberID)

	lamaHari := int(time.Since(member.JoinedAt).Hours() / 24)
	angsuranPerBulan := int64(math.Round(float64(amount)*rate/100)) + amount/int64(tenor)
	kapasitas := float64(totalSimpanan)/12 + 1
	beban := math.Min(float64(angsuranPerBulan)/kapasitas, 1.0)

	noPrior := 0.0
	if !hasPrior {
		noPrior = 1.0
	}
	features := map[string]float64{
		"ketepatan_bayar":         ketepatan,
		"rasio_simpanan_pinjaman": math.Min(float64(totalSimpanan)/float64(amount), 1.0),
		"lama_keanggotaan_hari":   float64(lamaHari),
		"konsistensi_simpanan":    konsistensi,
		"rasio_beban_angsuran":    beban,
		"no_prior_loans":          noPrior,
	}

	result, err := u.scoring.Score(ctx, adins.ScoringInput{
		MemberID: memberID, Features: features,
		JumlahDiajukan: amount, TenorBulan: tenor,
		TotalSimpanan: totalSimpanan, MaxPlafond: maxPlafond,
	})
	if err != nil {
		return nil, err
	}

	featuresJSON, _ := json.Marshal(features)
	reasonsJSON, _ := json.Marshal(result.Reasons)

	return &entity.CreditAssessment{
		Score: result.Score, Grade: result.Grade,
		Recommendation: result.Recommendation, LimitSuggested: result.LimitRekomendasi,
		Features: featuresJSON, Reasons: reasonsJSON, Source: result.Source,
	}, nil
}

func (u *loanApplicationUsecase) Create(ctx context.Context, coopID string, req *model.CreateLoanApplicationRequest) (*model.LoanApplicationResponse, error) {
	ca, err := u.buildAndScore(ctx, coopID, req.MemberID, req.Amount, req.TenorMonths)
	if err != nil {
		return nil, err
	}

	coopUUID, _ := uuid.Parse(coopID)
	memberUUID, _ := uuid.Parse(req.MemberID)
	app := &entity.LoanApplication{
		CooperativeID: coopUUID, MemberID: memberUUID,
		Amount: req.Amount, TenorMonths: req.TenorMonths, Status: "pending",
	}
	if req.Purpose != "" {
		app.Purpose = &req.Purpose
	}
	if err := u.appRepo.Create(ctx, app); err != nil {
		return nil, err
	}

	ca.ApplicationID = app.ID
	if err := u.appRepo.CreateAssessment(ctx, ca); err != nil {
		return nil, err
	}

	// Re-fetch to get member name
	meta, err2 := u.appRepo.FindByID(ctx, coopID, app.ID.String())
	if err2 == nil {
		resp := converter.ToLoanApplicationResponse(&meta.LoanApplication, meta.MemberName, meta.Assessment)
		return &resp, nil
	}
	resp := converter.ToLoanApplicationResponse(app, req.MemberID, ca)
	return &resp, nil
}

func (u *loanApplicationUsecase) CreateForMember(ctx context.Context, memberID, coopID string, req *model.CreatePortalLoanApplicationRequest) (*model.LoanApplicationResponse, error) {
	return u.Create(ctx, coopID, &model.CreateLoanApplicationRequest{
		MemberID: memberID, Amount: req.Amount,
		TenorMonths: req.TenorMonths, Purpose: req.Purpose,
	})
}

func (u *loanApplicationUsecase) List(ctx context.Context, coopID, status string) ([]model.LoanApplicationResponse, error) {
	apps, err := u.appRepo.FindAll(ctx, coopID, status)
	if err != nil {
		return nil, err
	}
	result := make([]model.LoanApplicationResponse, len(apps))
	for i, a := range apps {
		result[i] = converter.ToLoanApplicationResponse(&a.LoanApplication, a.MemberName, a.Assessment)
	}
	return result, nil
}

func (u *loanApplicationUsecase) ListForMember(ctx context.Context, memberID string) ([]model.LoanApplicationResponse, error) {
	apps, err := u.appRepo.FindAllByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}
	result := make([]model.LoanApplicationResponse, len(apps))
	for i, a := range apps {
		result[i] = converter.ToLoanApplicationResponse(&a.LoanApplication, a.MemberName, a.Assessment)
	}
	return result, nil
}

func (u *loanApplicationUsecase) Approve(ctx context.Context, coopID, appID, approvedByUserID string) (*model.ApproveApplicationResponse, error) {
	meta, err := u.appRepo.FindByID(ctx, coopID, appID)
	if errors.Is(err, repository.ErrLoanApplicationNotFound) {
		return nil, ErrLoanApplicationNotFound
	}
	if err != nil {
		return nil, err
	}
	if meta.Status != "pending" {
		return nil, ErrApplicationNotPending
	}

	approvedByUUID, _ := uuid.Parse(approvedByUserID)
	if err := u.appRepo.UpdateStatus(ctx, appID, "approved", &approvedByUUID); err != nil {
		return nil, err
	}

	cfg, _ := u.configRepo.FindByCoopID(ctx, coopID)
	rate := 1.5
	if cfg != nil {
		rate = cfg.FlatRateMonthly
	}

	coopUUID, _ := uuid.Parse(coopID)
	loan := &entity.Loan{
		CooperativeID:   coopUUID,
		ApplicationID:   meta.ID,
		MemberID:        meta.MemberID,
		Principal:       meta.Amount,
		FlatRateMonthly: rate,
		TenorMonths:     meta.TenorMonths,
		Status:          "aktif",
		DisbursedAt:     time.Now(),
	}
	if err := u.loanRepo.Create(ctx, loan); err != nil {
		return nil, err
	}

	// Create audit log for loan disbursement
	afterDataMap := map[string]interface{}{
		"id":                 loan.ID.String(),
		"cooperative_id":     loan.CooperativeID.String(),
		"application_id":     loan.ApplicationID.String(),
		"member_id":          loan.MemberID.String(),
		"principal":          loan.Principal,
		"flat_rate_monthly":  loan.FlatRateMonthly,
		"tenor_months":       loan.TenorMonths,
		"status":             loan.Status,
		"disbursed_at":       loan.DisbursedAt.Format("2006-01-02"),
	}
	afterDataJSON, _ := json.Marshal(afterDataMap)

	auditLog := &entity.LoanAuditLog{
		CooperativeID: loan.CooperativeID,
		LoanID:        loan.ID,
		Action:        "disburse",
		PerformedBy:   approvedByUUID,
		BeforeData:    datatypes.JSON("{}"),
		AfterData:     datatypes.JSON(afterDataJSON),
		Note:          "Pencairan pinjaman disetujui",
	}
	_ = u.auditRepo.CreateLog(ctx, auditLog)

	schedules := GenerateSchedule(loan.ID, loan.Principal, loan.FlatRateMonthly, loan.TenorMonths, loan.DisbursedAt)
	if err := u.loanRepo.BulkCreateSchedule(ctx, schedules); err != nil {
		return nil, err
	}

	updatedApp, _ := u.appRepo.FindByID(ctx, coopID, appID)
	loanWithMeta, _ := u.loanRepo.FindByID(ctx, coopID, loan.ID.String())

	instResponses := make([]model.InstallmentResponse, 0)
	if loanWithMeta != nil {
		instResponses = make([]model.InstallmentResponse, len(loanWithMeta.Schedule))
		for i, s := range loanWithMeta.Schedule {
			instResponses[i] = converter.ToInstallmentResponse(&s.InstallmentSchedule, s.PaidAmount)
		}
	}

	var appResp model.LoanApplicationResponse
	if updatedApp != nil {
		appResp = converter.ToLoanApplicationResponse(&updatedApp.LoanApplication, updatedApp.MemberName, updatedApp.Assessment)
	}

	var loanResp model.LoanResponse
	if loanWithMeta != nil {
		loanResp = converter.ToLoanResponse(&loanWithMeta.Loan, loanWithMeta.MemberName, loanWithMeta.Outstanding, instResponses)
	}

	return &model.ApproveApplicationResponse{Application: appResp, Loan: loanResp}, nil
}

func (u *loanApplicationUsecase) Reject(ctx context.Context, coopID, appID string) (*model.LoanApplicationResponse, error) {
	meta, err := u.appRepo.FindByID(ctx, coopID, appID)
	if errors.Is(err, repository.ErrLoanApplicationNotFound) {
		return nil, ErrLoanApplicationNotFound
	}
	if err != nil {
		return nil, err
	}
	if meta.Status != "pending" {
		return nil, ErrApplicationNotPending
	}
	if err := u.appRepo.UpdateStatus(ctx, appID, "rejected", nil); err != nil {
		return nil, err
	}
	meta.Status = "rejected"
	resp := converter.ToLoanApplicationResponse(&meta.LoanApplication, meta.MemberName, meta.Assessment)
	return &resp, nil
}
