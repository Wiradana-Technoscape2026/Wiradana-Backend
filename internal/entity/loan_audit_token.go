package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoanAuditToken struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID uuid.UUID `gorm:"type:uuid;not null" json:"cooperative_id"`
	LoanID        uuid.UUID `gorm:"type:uuid;not null" json:"loan_id"`
	TokenHash     string    `gorm:"not null;unique" json:"token_hash"`
	ExpiresAt     time.Time `gorm:"not null" json:"expires_at"`
	CreatedBy     uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt     time.Time `gorm:"not null;default:now()" json:"created_at"`
	Revoked       bool      `gorm:"not null;default:false" json:"revoked"`
}

func (t *LoanAuditToken) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	return nil
}
