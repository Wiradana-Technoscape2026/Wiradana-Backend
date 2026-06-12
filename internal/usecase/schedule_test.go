package usecase_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/usecase"
)

func TestGenerateSchedule_FlatExample(t *testing.T) {
	// be_implementation §5.4: principal=5_000_000, rate=1.5%, tenor=12
	loanID := uuid.New()
	disbursed := time.Date(2026, 6, 6, 0, 0, 0, 0, time.UTC)
	sched := usecase.GenerateSchedule(loanID, 5_000_000, 1.5, 12, disbursed)

	if len(sched) != 12 {
		t.Fatalf("want 12 periods, got %d", len(sched))
	}

	var totalPrincipal int64
	for i, s := range sched {
		if s.InterestDue != 75_000 {
			t.Errorf("period %d: want interest 75000, got %d", i+1, s.InterestDue)
		}
		totalPrincipal += s.PrincipalDue
	}

	if totalPrincipal != 5_000_000 {
		t.Errorf("total principal want 5000000, got %d", totalPrincipal)
	}

	// Period 1 due date = +1 month
	if !sched[0].DueDate.Equal(time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("period 1 due_date wrong: %v", sched[0].DueDate)
	}

	// All statuses are belum_bayar
	for i, s := range sched {
		if s.Status != "belum_bayar" {
			t.Errorf("period %d: want status belum_bayar, got %s", i+1, s.Status)
		}
	}
}

func TestGenerateSchedule_TotalDue(t *testing.T) {
	loanID := uuid.New()
	sched := usecase.GenerateSchedule(loanID, 5_000_000, 1.5, 12, time.Now())
	for i, s := range sched {
		if s.TotalDue != s.PrincipalDue+s.InterestDue {
			t.Errorf("period %d: total_due=%d != principal_due=%d + interest_due=%d", i+1, s.TotalDue, s.PrincipalDue, s.InterestDue)
		}
	}
}
