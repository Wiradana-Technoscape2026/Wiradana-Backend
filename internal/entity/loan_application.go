package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoanApplication struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID uuid.UUID  `gorm:"type:uuid;not null" json:"cooperative_id"`
	MemberID      uuid.UUID  `gorm:"type:uuid;not null" json:"member_id"`
	ApprovedBy    *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt    *time.Time `gorm:"column:approved_at" json:"approved_at"`
	Amount        int64      `gorm:"not null" json:"amount"`
	TenorMonths   int        `gorm:"not null" json:"tenor_months"`
	Purpose       *string    `json:"purpose"`
	Status        string     `gorm:"not null;default:pending" json:"status"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (l *LoanApplication) BeforeCreate(_ *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
