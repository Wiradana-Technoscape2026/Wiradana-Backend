package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SavingsTransaction struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID uuid.UUID `gorm:"type:uuid;not null" json:"cooperative_id"`
	MemberID      uuid.UUID `gorm:"type:uuid;not null" json:"member_id"`
	SavingsType   string    `gorm:"not null" json:"savings_type"`
	Direction     string    `gorm:"not null" json:"direction"`
	Amount        int64     `gorm:"not null" json:"amount"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (s *SavingsTransaction) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
