package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InstallmentSchedule struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	LoanID       uuid.UUID `gorm:"type:uuid;not null" json:"loan_id"`
	PeriodNo     int       `gorm:"not null" json:"period_no"`
	DueDate      time.Time `gorm:"type:date;not null" json:"due_date"`
	PrincipalDue int64     `gorm:"not null" json:"principal_due"`
	InterestDue  int64     `gorm:"not null" json:"interest_due"`
	TotalDue     int64     `gorm:"not null" json:"total_due"`
	Status       string    `gorm:"not null;default:belum_bayar" json:"status"`
}

func (i *InstallmentSchedule) BeforeCreate(_ *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
