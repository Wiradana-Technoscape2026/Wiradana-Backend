package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Loan struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID   uuid.UUID `gorm:"type:uuid;not null" json:"cooperative_id"`
	ApplicationID   uuid.UUID `gorm:"type:uuid;not null" json:"application_id"`
	MemberID        uuid.UUID `gorm:"type:uuid;not null" json:"member_id"`
	Principal       int64     `gorm:"not null" json:"principal"`
	FlatRateMonthly float64   `gorm:"type:numeric(5,2);not null" json:"flat_rate_monthly"`
	TenorMonths     int       `gorm:"not null" json:"tenor_months"`
	Status          string    `gorm:"not null;default:aktif" json:"status"`
	DisbursedAt     time.Time `gorm:"type:date;not null" json:"disbursed_at"`
}

func (l *Loan) BeforeCreate(_ *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	if l.DisbursedAt.IsZero() {
		l.DisbursedAt = time.Now()
	}
	return nil
}
