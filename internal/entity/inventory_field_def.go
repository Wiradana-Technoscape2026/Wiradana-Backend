package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type InventoryFieldDef struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID uuid.UUID      `gorm:"type:uuid;not null" json:"cooperative_id"`
	FieldKey      string         `gorm:"not null" json:"field_key"`
	Label         string         `gorm:"not null" json:"label"`
	DataType      string         `gorm:"not null;default:text" json:"data_type"`
	Options       datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'" json:"options"`
	Required      bool           `gorm:"not null;default:false" json:"required"`
	SortOrder     int            `gorm:"not null;default:0" json:"sort_order"`
}

func (f *InventoryFieldDef) BeforeCreate(_ *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
