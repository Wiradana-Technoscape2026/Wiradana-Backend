package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DueInstallmentRow struct {
	InstallmentID uuid.UUID
	DueDate       time.Time
	TotalDue      int64
	CooperativeID uuid.UUID
	MemberID      uuid.UUID
	MemberName    string
	PhoneNumber   string
}

type NotificationLogWithMember struct {
	entity.NotificationLog
	MemberName string `gorm:"column:member_name"`
}

type NotificationRepository interface {
	FindDueInstallments(ctx context.Context, targetDate time.Time) ([]DueInstallmentRow, error)
	FindOverdueInstallments(ctx context.Context) ([]DueInstallmentRow, error)
	CreateLog(ctx context.Context, log *entity.NotificationLog) error
	ListLogs(ctx context.Context, coopID string, eventType string, limit, offset int) ([]NotificationLogWithMember, int64, error)
}

type notificationRepository struct{ db *gorm.DB }

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

const dueInstallmentQuery = `
SELECT
    i.id            AS installment_id,
    i.due_date,
    i.total_due,
    l.cooperative_id,
    l.member_id,
    m.full_name     AS member_name,
    m.phone_number  AS phone_number
FROM installment_schedule i
JOIN loan l ON l.id = i.loan_id
JOIN member m ON m.id = l.member_id
WHERE i.status = 'belum_bayar'
AND DATE(i.due_date) = ?
AND m.phone_number IS NOT NULL
AND m.phone_number != ''
`

func (r *notificationRepository) FindDueInstallments(ctx context.Context, targetDate time.Time) ([]DueInstallmentRow, error) {
	var rows []DueInstallmentRow
	err := r.db.WithContext(ctx).Raw(dueInstallmentQuery, targetDate.Format("2006-01-02")).Scan(&rows).Error
	return rows, err
}

func (r *notificationRepository) FindOverdueInstallments(ctx context.Context) ([]DueInstallmentRow, error) {
	var rows []DueInstallmentRow
	err := r.db.WithContext(ctx).Raw(`
SELECT
    i.id            AS installment_id,
    i.due_date,
    i.total_due,
    l.cooperative_id,
    l.member_id,
    m.full_name     AS member_name,
    m.phone_number  AS phone_number
FROM installment_schedule i
JOIN loan l ON l.id = i.loan_id
JOIN member m ON m.id = l.member_id
WHERE i.status = 'terlambat'
AND m.phone_number IS NOT NULL
AND m.phone_number != ''
`).Scan(&rows).Error
	return rows, err
}

// CreateLog inserts a notification log. Uses ON CONFLICT DO NOTHING for dedup index.
func (r *notificationRepository) CreateLog(ctx context.Context, log *entity.NotificationLog) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(log).Error
}

func (r *notificationRepository) ListLogs(ctx context.Context, coopID, eventType string, limit, offset int) ([]NotificationLogWithMember, int64, error) {
	query := r.db.WithContext(ctx).
		Table("notification_log nl").
		Select("nl.*, m.full_name AS member_name").
		Joins("JOIN member m ON m.id = nl.member_id").
		Where("nl.cooperative_id = ?", coopID)

	if eventType != "" {
		query = query.Where("nl.event_type = ?", eventType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []NotificationLogWithMember
	err := query.Order("nl.created_at DESC").Limit(limit).Offset(offset).Scan(&logs).Error
	return logs, total, err
}
