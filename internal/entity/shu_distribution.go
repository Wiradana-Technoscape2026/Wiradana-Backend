package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShuDistribution struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ShuPeriodID uuid.UUID `gorm:"type:uuid;not null" json:"shu_period_id"`
	MemberID    uuid.UUID `gorm:"type:uuid;not null" json:"member_id"`
	JasaModal   int64     `gorm:"not null" json:"jasa_modal"`
	JasaUsaha   int64     `gorm:"not null" json:"jasa_usaha"`
	TotalShu    int64     `gorm:"not null" json:"total_shu"`
}

func (s *ShuDistribution) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
