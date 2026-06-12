package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationLog struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CooperativeID uuid.UUID  `gorm:"type:uuid;not null" json:"cooperative_id"`
	MemberID      uuid.UUID  `gorm:"type:uuid;not null" json:"member_id"`
	PhoneNumber   string     `gorm:"not null" json:"phone_number"`
	Channel       string     `gorm:"not null;default:whatsapp" json:"channel"`
	EventType     string     `gorm:"not null" json:"event_type"`
	RefID         *uuid.UUID `gorm:"type:uuid" json:"ref_id,omitempty"`
	Message       string     `gorm:"type:text;not null" json:"message"`
	Status        string     `gorm:"not null;default:sent" json:"status"`
	Source        string     `gorm:"not null;default:MOCK_WA" json:"source"`
	ErrorMsg      *string    `gorm:"type:text" json:"error_msg,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (n *NotificationLog) BeforeCreate(_ *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
