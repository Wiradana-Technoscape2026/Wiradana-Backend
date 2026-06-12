package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppUser struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID uuid.UUID  `gorm:"type:uuid;not null" json:"cooperative_id"`
	MemberID      *uuid.UUID `gorm:"type:uuid" json:"member_id"`
	Email         string     `gorm:"not null" json:"email"`
	PasswordHash  string     `gorm:"not null" json:"-"`
	Role          string     `gorm:"not null" json:"role"`
}

func (u *AppUser) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
