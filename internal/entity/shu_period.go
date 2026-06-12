package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShuPeriod struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID uuid.UUID `gorm:"type:uuid;not null" json:"cooperative_id"`
	Year          int       `gorm:"not null" json:"year"`
	TotalShu      int64     `gorm:"not null" json:"total_shu"`
	PctJasaModal  float64   `gorm:"type:numeric(5,2);not null" json:"pct_jasa_modal"`
	PctJasaUsaha  float64   `gorm:"type:numeric(5,2);not null" json:"pct_jasa_usaha"`
	Status        string    `gorm:"not null;default:draft" json:"status"`
}

func (s *ShuPeriod) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
