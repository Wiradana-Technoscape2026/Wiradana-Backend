package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CreditAssessment struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ApplicationID  uuid.UUID      `gorm:"type:uuid;not null" json:"application_id"`
	Score          int            `gorm:"not null" json:"score"`
	Grade          string         `gorm:"not null" json:"grade"`
	Recommendation string         `gorm:"not null" json:"recommendation"`
	LimitSuggested int64          `gorm:"not null" json:"limit_suggested"`
	Features       datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"features"`
	Reasons        datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'" json:"reasons"`
	Source         string         `gorm:"not null;default:MOCK_ADINS_SCORING" json:"source"`
}

func (c *CreditAssessment) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
