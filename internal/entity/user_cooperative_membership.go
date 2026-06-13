package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserCooperativeMembership struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID    `gorm:"type:uuid;not null" json:"user_id"`
	CooperativeID uuid.UUID    `gorm:"type:uuid;not null" json:"cooperative_id"`
	MemberID      *uuid.UUID   `gorm:"type:uuid" json:"member_id"`
	CreatedAt     time.Time    `json:"created_at"`
	Cooperative   *Cooperative `gorm:"foreignKey:CooperativeID" json:"-"`
}

func (m *UserCooperativeMembership) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
