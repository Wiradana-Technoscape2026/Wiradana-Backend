package entity

import (
	"time"

	"github.com/google/uuid"
)

type IdempotencyKey struct {
	Key           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	CooperativeID uuid.UUID  `gorm:"type:uuid;not null"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null"`
	MutationType  string     `gorm:"not null"`
	ResultID      *uuid.UUID `gorm:"type:uuid"`
	Status        string     `gorm:"not null"`
	ErrorMessage  *string
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	ProcessedAt   *time.Time
}
