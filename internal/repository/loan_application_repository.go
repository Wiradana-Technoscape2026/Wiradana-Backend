package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var ErrLoanApplicationNotFound = errors.New("pengajuan pinjaman tidak ditemukan")

type LoanApplicationWithMeta struct {
	entity.LoanApplication
	MemberName     string
	ApprovedByName *string
	Assessment     *entity.CreditAssessment
}

type LoanApplicationRepository interface {
	Create(ctx context.Context, app *entity.LoanApplication) error
	CreateAssessment(ctx context.Context, ca *entity.CreditAssessment) error
	FindByID(ctx context.Context, coopID, appID string) (*LoanApplicationWithMeta, error)
	FindAll(ctx context.Context, coopID, status string) ([]*LoanApplicationWithMeta, error)
	FindAllByMember(ctx context.Context, memberID string) ([]*LoanApplicationWithMeta, error)
	UpdateStatus(ctx context.Context, appID, status string, approvedBy *uuid.UUID) error
	GetTotalSavings(ctx context.Context, memberID string) (int64, error)
	GetKetepatanBayar(ctx context.Context, memberID string) (float64, error)
	GetKonsistensiSimpanan(ctx context.Context, memberID string, joinedAt time.Time) (float64, error)
}

type loanApplicationRepository struct{ db *gorm.DB }

func NewLoanApplicationRepository(db *gorm.DB) LoanApplicationRepository {
	return &loanApplicationRepository{db: db}
}

func (r *loanApplicationRepository) Create(ctx context.Context, app *entity.LoanApplication) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *loanApplicationRepository) CreateAssessment(ctx context.Context, ca *entity.CreditAssessment) error {
	return r.db.WithContext(ctx).Create(ca).Error
}

func (r *loanApplicationRepository) loadAssessments(ctx context.Context, apps []*LoanApplicationWithMeta) {
	if len(apps) == 0 {
		return
	}
	ids := make([]uuid.UUID, len(apps))
	for i, a := range apps {
		ids[i] = a.ID
	}
	var cas []entity.CreditAssessment
	r.db.WithContext(ctx).Where("application_id IN ?", ids).Find(&cas)
	m := make(map[uuid.UUID]*entity.CreditAssessment, len(cas))
	for i := range cas {
		m[cas[i].ApplicationID] = &cas[i]
	}
	for _, a := range apps {
		a.Assessment = m[a.ID]
	}
}

type appRow struct {
	entity.LoanApplication
	MemberName     string  `gorm:"column:member_name"`
	ApprovedByName *string `gorm:"column:approved_by_name"`
}

func (r *loanApplicationRepository) FindAll(ctx context.Context, coopID, status string) ([]*LoanApplicationWithMeta, error) {
	var rows []appRow
	q := r.db.WithContext(ctx).
		Table("loan_application").
		Select("loan_application.*, member.full_name as member_name, au.full_name as approved_by_name").
		Joins("JOIN member ON member.id = loan_application.member_id").
		Joins("LEFT JOIN app_user au ON au.id = loan_application.approved_by").
		Where("loan_application.cooperative_id = ?", coopID)
	if status != "" {
		q = q.Where("loan_application.status = ?", status)
	}
	if err := q.Order("loan_application.created_at DESC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*LoanApplicationWithMeta, len(rows))
	for i, row := range rows {
		result[i] = &LoanApplicationWithMeta{LoanApplication: row.LoanApplication, MemberName: row.MemberName, ApprovedByName: row.ApprovedByName}
	}
	r.loadAssessments(ctx, result)
	return result, nil
}

func (r *loanApplicationRepository) FindAllByMember(ctx context.Context, memberID string) ([]*LoanApplicationWithMeta, error) {
	var rows []appRow
	if err := r.db.WithContext(ctx).
		Table("loan_application").
		Select("loan_application.*, member.full_name as member_name, au.full_name as approved_by_name").
		Joins("JOIN member ON member.id = loan_application.member_id").
		Joins("LEFT JOIN app_user au ON au.id = loan_application.approved_by").
		Where("loan_application.member_id = ?", memberID).
		Order("loan_application.created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*LoanApplicationWithMeta, len(rows))
	for i, row := range rows {
		result[i] = &LoanApplicationWithMeta{LoanApplication: row.LoanApplication, MemberName: row.MemberName, ApprovedByName: row.ApprovedByName}
	}
	r.loadAssessments(ctx, result)
	return result, nil
}

func (r *loanApplicationRepository) FindByID(ctx context.Context, coopID, appID string) (*LoanApplicationWithMeta, error) {
	var row appRow
	q := r.db.WithContext(ctx).
		Table("loan_application").
		Select("loan_application.*, member.full_name as member_name, au.full_name as approved_by_name").
		Joins("JOIN member ON member.id = loan_application.member_id").
		Joins("LEFT JOIN app_user au ON au.id = loan_application.approved_by").
		Where("loan_application.id = ?", appID)
	if coopID != "" {
		q = q.Where("loan_application.cooperative_id = ?", coopID)
	}
	if err := q.Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.ID == uuid.Nil {
		return nil, ErrLoanApplicationNotFound
	}
	meta := &LoanApplicationWithMeta{LoanApplication: row.LoanApplication, MemberName: row.MemberName, ApprovedByName: row.ApprovedByName}
	r.loadAssessments(ctx, []*LoanApplicationWithMeta{meta})
	return meta, nil
}

func (r *loanApplicationRepository) UpdateStatus(ctx context.Context, appID, status string, approvedBy *uuid.UUID) error {
	updates := map[string]interface{}{"status": status}
	if approvedBy != nil {
		updates["approved_by"] = approvedBy
		now := time.Now()
		updates["approved_at"] = &now
	}
	return r.db.WithContext(ctx).Model(&entity.LoanApplication{}).
		Where("id = ?", appID).Updates(updates).Error
}

func (r *loanApplicationRepository) GetTotalSavings(ctx context.Context, memberID string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(CASE WHEN direction='setor' THEN amount ELSE -amount END), 0)
		FROM savings_transaction WHERE member_id = ?`, memberID).Scan(&total).Error
	return total, err
}

func (r *loanApplicationRepository) GetKetepatanBayar(ctx context.Context, memberID string) (float64, error) {
	var totalDue, totalLunas int64
	r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM installment_schedule i
		JOIN loan l ON l.id = i.loan_id
		WHERE l.member_id = ? AND i.status IN ('lunas','terlambat')`, memberID).Scan(&totalDue)
	r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM installment_schedule i
		JOIN loan l ON l.id = i.loan_id
		WHERE l.member_id = ? AND i.status = 'lunas'`, memberID).Scan(&totalLunas)
	if totalDue == 0 {
		return 1.0, nil
	}
	return float64(totalLunas) / float64(totalDue), nil
}

func (r *loanApplicationRepository) GetKonsistensiSimpanan(ctx context.Context, memberID string, joinedAt time.Time) (float64, error) {
	monthsMember := int(time.Since(joinedAt).Hours() / 24 / 30)
	if monthsMember <= 0 {
		return 1.0, nil
	}
	var monthsWithWajib int64
	r.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT DATE_TRUNC('month', created_at))
		FROM savings_transaction
		WHERE member_id = ? AND savings_type = 'wajib' AND direction = 'setor'`, memberID).Scan(&monthsWithWajib)
	cap := monthsMember
	if cap > 12 {
		cap = 12
	}
	return float64(monthsWithWajib) / float64(cap), nil
}
