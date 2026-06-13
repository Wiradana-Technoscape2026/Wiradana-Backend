package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LoanAuditLog struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID uuid.UUID      `gorm:"type:uuid;not null" json:"cooperative_id"`
	LoanID        uuid.UUID      `gorm:"type:uuid;not null" json:"loan_id"`
	Action        string         `gorm:"not null" json:"action"`
	PerformedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"performed_by"`
	PerformedAt   time.Time      `gorm:"not null;default:now()" json:"performed_at"`
	BeforeData    datatypes.JSON `gorm:"type:jsonb" json:"before_data"`
	AfterData     datatypes.JSON `gorm:"type:jsonb" json:"after_data"`
	Note          string         `json:"note"`
	IsFlagged     bool           `gorm:"not null;default:false" json:"is_flagged"`
	FlaggedByName string         `json:"flagged_by_name"`
	FlaggedAt     *time.Time     `json:"flagged_at"`
	FlaggedReason string         `json:"flagged_reason"`
}

func (l *LoanAuditLog) BeforeCreate(_ *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	if l.PerformedAt.IsZero() {
		l.PerformedAt = time.Now()
	}
	return nil
}
