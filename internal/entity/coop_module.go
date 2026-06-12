package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CoopModule struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID uuid.UUID `gorm:"type:uuid;not null" json:"cooperative_id"`
	ModuleKey     string    `gorm:"not null" json:"module_key"`
	Enabled       bool      `gorm:"not null;default:false" json:"enabled"`
}

func (c *CoopModule) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
