package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type InventoryMovement struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID    uuid.UUID      `gorm:"type:uuid;not null" json:"cooperative_id"`
	ProductID        uuid.UUID      `gorm:"type:uuid;not null" json:"product_id"`
	Direction        string         `gorm:"not null" json:"direction"` // masuk | keluar
	Quantity         int64          `gorm:"not null" json:"quantity"`
	Note             string         `json:"note"`
	CustomAttributes datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"custom_attributes"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

func (m *InventoryMovement) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
