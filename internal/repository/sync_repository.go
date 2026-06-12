package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"gorm.io/gorm"
)

var ErrIdempotencyKeyNotFound = errors.New("idempotency key tidak ditemukan")

type SyncRepository interface {
	FindKey(ctx context.Context, key uuid.UUID) (*entity.IdempotencyKey, error)
	SaveKey(ctx context.Context, k *entity.IdempotencyKey) error
	PullDelta(ctx context.Context, coopID uuid.UUID, since time.Time, role string, memberID *uuid.UUID) (*model.SyncPullResponse, error)
}

type syncRepository struct {
	db *gorm.DB
}

func NewSyncRepository(db *gorm.DB) SyncRepository {
	return &syncRepository{db: db}
}

func (r *syncRepository) FindKey(ctx context.Context, key uuid.UUID) (*entity.IdempotencyKey, error) {
	var k entity.IdempotencyKey
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&k).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIdempotencyKeyNotFound
		}
		return nil, err
	}
	return &k, nil
}

func (r *syncRepository) SaveKey(ctx context.Context, k *entity.IdempotencyKey) error {
	return r.db.WithContext(ctx).Save(k).Error
}

func (r *syncRepository) PullDelta(ctx context.Context, coopID uuid.UUID, since time.Time, role string, memberID *uuid.UUID) (*model.SyncPullResponse, error) {
	cursor := time.Now()
	resp := &model.SyncPullResponse{
		Cursor:               cursor,
		Members:              []model.MemberResponse{},
		SavingsTransactions:  []model.SavingsTransactionResponse{},
		LoanApplications:     []model.LoanApplicationResponse{},
		Loans:                []model.LoanResponse{},
		InstallmentSchedules: []model.InstallmentResponse{},
		Payments:             []model.PaymentResponse{},
		ShuPeriods:           []model.ShuPeriodResponse{},
		ShuDistributions:     []model.ShuDistributionResponse{},
	}

	sinceIsZero := since.IsZero()

	// Build member name map (used for loan/application responses)
	memberNameMap := map[uuid.UUID]string{}
	{
		var members []entity.Member
		r.db.WithContext(ctx).Where("cooperative_id = ?", coopID).Find(&members)
		for _, m := range members {
			memberNameMap[m.ID] = m.FullName
		}
	}

	// ── Members (pengurus only) ──────────────────────────────────────────────
	if role == "pengurus" {
		var members []entity.Member
		q := r.db.WithContext(ctx).Where("cooperative_id = ?", coopID)
		if !sinceIsZero {
			q = q.Where("joined_at > ?", since)
		}
		q.Find(&members)
		for _, m := range members {
			summary, _ := r.getMemberSavingsSummary(ctx, m.ID.String(), coopID.String())
			resp.Members = append(resp.Members, converter.ToMemberResponse(&m, summary))
		}
	}

	// ── Savings Transactions ─────────────────────────────────────────────────
	{
		var savings []entity.SavingsTransaction
		q := r.db.WithContext(ctx).Where("cooperative_id = ?", coopID)
		if !sinceIsZero {
			q = q.Where("created_at > ?", since)
		}
		if role == "anggota" && memberID != nil {
			q = q.Where("member_id = ?", memberID)
		}
		q.Find(&savings)
		for _, s := range savings {
			resp.SavingsTransactions = append(resp.SavingsTransactions, converter.ToSavingsResponse(&s))
		}
	}

	// ── Loan Applications ────────────────────────────────────────────────────
	{
		var apps []entity.LoanApplication
		q := r.db.WithContext(ctx).Where("cooperative_id = ?", coopID)
		if !sinceIsZero {
			q = q.Where("created_at > ?", since)
		}
		if role == "anggota" && memberID != nil {
			q = q.Where("member_id = ?", memberID)
		}
		q.Find(&apps)

		if len(apps) > 0 {
			appIDs := make([]uuid.UUID, len(apps))
			for i, a := range apps {
				appIDs[i] = a.ID
			}
			assessmentMap := r.fetchAssessmentMap(ctx, appIDs)
			for _, a := range apps {
				resp.LoanApplications = append(resp.LoanApplications,
					converter.ToLoanApplicationResponse(&a, memberNameMap[a.MemberID], assessmentMap[a.ID]))
			}
		}
	}

	// ── Loans ────────────────────────────────────────────────────────────────
	{
		var loans []entity.Loan
		q := r.db.WithContext(ctx).Where("cooperative_id = ?", coopID)
		if !sinceIsZero {
			q = q.Where("disbursed_at > ?", since)
		}
		if role == "anggota" && memberID != nil {
			q = q.Where("member_id = ?", memberID)
		}
		q.Find(&loans)
		for _, l := range loans {
			lCopy := l
			resp.Loans = append(resp.Loans,
				converter.ToLoanResponse(&lCopy, memberNameMap[l.MemberID], 0, nil))
		}
	}

	// ── Installment Schedules ────────────────────────────────────────────────
	{
		var schedules []entity.InstallmentSchedule
		q := r.db.WithContext(ctx).
			Joins("JOIN loan ON loan.id = installment_schedule.loan_id").
			Where("loan.cooperative_id = ?", coopID)
		if role == "anggota" && memberID != nil {
			q = q.Where("loan.member_id = ?", memberID)
		}
		q.Find(&schedules)

		if len(schedules) > 0 {
			schedIDs := make([]uuid.UUID, len(schedules))
			for i, s := range schedules {
				schedIDs[i] = s.ID
			}
			paidMap := r.fetchPaidAmountMap(ctx, schedIDs)
			for _, s := range schedules {
				sCopy := s
				resp.InstallmentSchedules = append(resp.InstallmentSchedules,
					converter.ToInstallmentResponse(&sCopy, paidMap[s.ID]))
			}
		}
	}

	// ── Payments ─────────────────────────────────────────────────────────────
	{
		var payments []entity.Payment
		q := r.db.WithContext(ctx).
			Joins("JOIN installment_schedule ON installment_schedule.id = payment.schedule_id").
			Joins("JOIN loan ON loan.id = installment_schedule.loan_id").
			Where("loan.cooperative_id = ?", coopID)
		if !sinceIsZero {
			q = q.Where("payment.paid_at > ?", since)
		}
		if role == "anggota" && memberID != nil {
			q = q.Where("loan.member_id = ?", memberID)
		}
		q.Find(&payments)
		for _, p := range payments {
			pCopy := p
			resp.Payments = append(resp.Payments, converter.ToPaymentResponse(&pCopy))
		}
	}

	// ── SHU Periods ──────────────────────────────────────────────────────────
	{
		var periods []entity.ShuPeriod
		r.db.WithContext(ctx).Where("cooperative_id = ?", coopID).Find(&periods)
		for _, p := range periods {
			pCopy := p
			resp.ShuPeriods = append(resp.ShuPeriods, converter.ToShuPeriodResponse(&pCopy))
		}
	}

	// ── SHU Distributions ────────────────────────────────────────────────────
	{
		var dists []entity.ShuDistribution
		q := r.db.WithContext(ctx).
			Joins("JOIN shu_period ON shu_period.id = shu_distribution.shu_period_id").
			Where("shu_period.cooperative_id = ?", coopID)
		if role == "anggota" && memberID != nil {
			q = q.Where("shu_distribution.member_id = ?", memberID)
		}
		q.Find(&dists)
		for _, d := range dists {
			dCopy := d
			resp.ShuDistributions = append(resp.ShuDistributions, converter.ToShuDistributionResponse(&dCopy))
		}
	}

	// ── Loan Config (pengurus only) ──────────────────────────────────────────
	if role == "pengurus" {
		var lc entity.LoanConfig
		err := r.db.WithContext(ctx).Where("cooperative_id = ?", coopID).First(&lc).Error
		if err == nil {
			r := converter.ToLoanConfigResponse(&lc)
			resp.LoanConfig = &r
		}
	}

	return resp, nil
}

func (r *syncRepository) getMemberSavingsSummary(ctx context.Context, memberID, coopID string) (*model.SavingsSummary, error) {
	type row struct {
		Pokok    int64
		Wajib    int64
		Sukarela int64
	}
	var result row
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(amount) FILTER (WHERE savings_type='pokok'    AND direction='setor'), 0) AS pokok,
			COALESCE(SUM(amount) FILTER (WHERE savings_type='wajib'    AND direction='setor'), 0)
			- COALESCE(SUM(amount) FILTER (WHERE savings_type='wajib'  AND direction='tarik'), 0) AS wajib,
			COALESCE(SUM(amount) FILTER (WHERE savings_type='sukarela' AND direction='setor'), 0)
			- COALESCE(SUM(amount) FILTER (WHERE savings_type='sukarela' AND direction='tarik'), 0) AS sukarela
		FROM savings_transaction
		WHERE member_id = ? AND cooperative_id = ?
	`, memberID, coopID).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	s := &model.SavingsSummary{
		Pokok:    result.Pokok,
		Wajib:    result.Wajib,
		Sukarela: result.Sukarela,
		Total:    result.Pokok + result.Wajib + result.Sukarela,
	}
	return s, nil
}

func (r *syncRepository) fetchAssessmentMap(ctx context.Context, appIDs []uuid.UUID) map[uuid.UUID]*entity.CreditAssessment {
	var assessments []entity.CreditAssessment
	r.db.WithContext(ctx).Where("application_id IN ?", appIDs).Find(&assessments)
	m := make(map[uuid.UUID]*entity.CreditAssessment, len(assessments))
	for i := range assessments {
		m[assessments[i].ApplicationID] = &assessments[i]
	}
	return m
}

func (r *syncRepository) fetchPaidAmountMap(ctx context.Context, schedIDs []uuid.UUID) map[uuid.UUID]int64 {
	type row struct {
		ScheduleID uuid.UUID
		PaidAmount int64
	}
	var rows []row
	r.db.WithContext(ctx).Raw(
		`SELECT schedule_id, COALESCE(SUM(amount), 0) AS paid_amount FROM payment WHERE schedule_id IN ? GROUP BY schedule_id`,
		schedIDs,
	).Scan(&rows)
	m := make(map[uuid.UUID]int64, len(rows))
	for _, row := range rows {
		m[row.ScheduleID] = row.PaidAmount
	}
	return m
}
