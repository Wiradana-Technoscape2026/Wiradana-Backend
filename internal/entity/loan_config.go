package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoanConfig struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID   uuid.UUID `gorm:"type:uuid;not null" json:"cooperative_id"`
	FlatRateMonthly float64   `gorm:"type:numeric(5,2);not null;default:1.5" json:"flat_rate_monthly"`
	MaxPlafond      int64     `gorm:"not null;default:20000000" json:"max_plafond"`
	PenaltyDaily    int64     `gorm:"not null;default:0" json:"penalty_daily"`
}

func (l *LoanConfig) BeforeCreate(_ *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
