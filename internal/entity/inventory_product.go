package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type InventoryProduct struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID    uuid.UUID      `gorm:"type:uuid;not null" json:"cooperative_id"`
	SKU              string         `gorm:"column:sku;not null" json:"sku"`
	Name             string         `gorm:"not null" json:"name"`
	Unit             string         `gorm:"not null;default:pcs" json:"unit"`
	CustomAttributes datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"custom_attributes"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *InventoryProduct) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
