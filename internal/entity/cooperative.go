package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cooperative struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null;default:simpan_pinjam" json:"type"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (c *Cooperative) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
