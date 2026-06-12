package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Member struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID    uuid.UUID      `gorm:"type:uuid;not null" json:"cooperative_id"`
	NIK              string         `gorm:"column:nik;not null" json:"nik"`
	FullName         string         `gorm:"not null" json:"full_name"`
	Address          *string        `json:"address"`
	BirthDate        *time.Time     `gorm:"type:date" json:"birth_date"`
	Status           string         `gorm:"not null;default:aktif" json:"status"`
	CustomAttributes datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"custom_attributes"`
	JoinedAt         time.Time      `gorm:"autoCreateTime" json:"joined_at"`
}

func (m *Member) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
