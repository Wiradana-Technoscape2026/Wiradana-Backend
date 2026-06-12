package usecase

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
)

// GenerateSchedule creates flat-rate installment schedule.
// Last period absorbs rounding remainder so total principal sums exactly.
func GenerateSchedule(loanID uuid.UUID, principal int64, rateMonthly float64, tenorMonths int, disbursedAt time.Time) []entity.InstallmentSchedule {
	interestPerPeriod := int64(math.Round(float64(principal) * rateMonthly / 100))
	principalPerPeriod := principal / int64(tenorMonths)

	schedules := make([]entity.InstallmentSchedule, tenorMonths)
	for i := 0; i < tenorMonths; i++ {
		p := principalPerPeriod
		if i == tenorMonths-1 {
			p = principal - principalPerPeriod*int64(tenorMonths-1)
		}
		schedules[i] = entity.InstallmentSchedule{
			LoanID:       loanID,
			PeriodNo:     i + 1,
			DueDate:      disbursedAt.AddDate(0, i+1, 0),
			PrincipalDue: p,
			InterestDue:  interestPerPeriod,
			TotalDue:     p + interestPerPeriod,
			Status:       "belum_bayar",
		}
	}
	return schedules
}
