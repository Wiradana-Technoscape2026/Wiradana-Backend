package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type InventoryFieldDef struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID uuid.UUID      `gorm:"type:uuid;not null" json:"cooperative_id"`
	FieldKey      string         `gorm:"not null" json:"field_key"`
	Label         string         `gorm:"not null" json:"label"`
	DataType      string         `gorm:"not null;default:text" json:"data_type"`
	Options       datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"options"`
	Required      bool           `gorm:"not null;default:false" json:"required"`
	SortOrder     int            `gorm:"not null;default:0" json:"sort_order"`
}

func (i *InventoryFieldDef) BeforeCreate(_ *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

type InventoryProduct struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID    uuid.UUID      `gorm:"type:uuid;not null" json:"cooperative_id"`
	SKU              string         `gorm:"not null" json:"sku"`
	Name             string         `gorm:"not null" json:"name"`
	Unit             string         `gorm:"not null" json:"unit"`
	CustomAttributes datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"custom_attributes"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (i *InventoryProduct) BeforeCreate(_ *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

type InventoryMovement struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID    uuid.UUID      `gorm:"type:uuid;not null" json:"cooperative_id"`
	ProductID        uuid.UUID      `gorm:"type:uuid;not null" json:"product_id"`
	Direction        string         `gorm:"not null" json:"direction"`
	Quantity         float64        `gorm:"not null" json:"quantity"`
	Note             *string        `json:"note"`
	CustomAttributes datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"custom_attributes"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

func (i *InventoryMovement) BeforeCreate(_ *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
