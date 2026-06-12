package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ScheduleID uuid.UUID `gorm:"type:uuid;not null" json:"schedule_id"`
	Amount     int64     `gorm:"not null" json:"amount"`
	Penalty    int64     `gorm:"not null;default:0" json:"penalty"`
	PaidAt     time.Time `gorm:"autoCreateTime" json:"paid_at"`
}

func (p *Payment) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
