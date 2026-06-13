package repository

import (
	"context"
	"errors"
	"time"

	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
)

var (
	ErrAuditTokenNotFound = errors.New("audit token tidak ditemukan")
	ErrAuditLogNotFound   = errors.New("audit log tidak ditemukan")
)

type LoanAuditLogWithUser struct {
	entity.LoanAuditLog
	PerformedByEmail string `gorm:"column:performed_by_email"`
}

type LoanAuditRepository interface {
	CreateLog(ctx context.Context, log *entity.LoanAuditLog) error
	CreateToken(ctx context.Context, token *entity.LoanAuditToken) error
	FindLogsByLoanID(ctx context.Context, loanID string) ([]*LoanAuditLogWithUser, error)
	FindTokenByHash(ctx context.Context, tokenHash string) (*entity.LoanAuditToken, error)
	FlagLog(ctx context.Context, logID string, flaggedByName string, flaggedReason string) error
}

type loanAuditRepository struct {
	db *gorm.DB
}

func NewLoanAuditRepository(db *gorm.DB) LoanAuditRepository {
	return &loanAuditRepository{db: db}
}

func (r *loanAuditRepository) CreateLog(ctx context.Context, log *entity.LoanAuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *loanAuditRepository) CreateToken(ctx context.Context, token *entity.LoanAuditToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *loanAuditRepository) FindLogsByLoanID(ctx context.Context, loanID string) ([]*LoanAuditLogWithUser, error) {
	var logs []*LoanAuditLogWithUser
	err := r.db.WithContext(ctx).Table("loan_audit_log").
		Select("loan_audit_log.*, app_user.email as performed_by_email").
		Joins("JOIN app_user ON app_user.id = loan_audit_log.performed_by").
		Where("loan_audit_log.loan_id = ?", loanID).
		Order("loan_audit_log.performed_at ASC").
		Scan(&logs).Error
	return logs, err
}


func (r *loanAuditRepository) FindTokenByHash(ctx context.Context, tokenHash string) (*entity.LoanAuditToken, error) {
	var token entity.LoanAuditToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked = false AND expires_at > ?", tokenHash, time.Now()).
		First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAuditTokenNotFound
	}
	return &token, err
}

func (r *loanAuditRepository) FlagLog(ctx context.Context, logID string, flaggedByName string, flaggedReason string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&entity.LoanAuditLog{}).
		Where("id = ?", logID).
		Updates(map[string]interface{}{
			"is_flagged":     true,
			"flagged_by_name": flaggedByName,
			"flagged_at":     &now,
			"flagged_reason": flaggedReason,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrAuditLogNotFound
	}
	return nil
}
