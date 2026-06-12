package usecase

import (
	"context"

	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
)

type DashboardUsecase interface {
	Get(ctx context.Context, cooperativeID string) (*model.DashboardResponse, error)
}

type dashboardUsecase struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardUsecase(dashboardRepo repository.DashboardRepository) DashboardUsecase {
	return &dashboardUsecase{dashboardRepo: dashboardRepo}
}

func (u *dashboardUsecase) Get(ctx context.Context, cooperativeID string) (*model.DashboardResponse, error) {
	stats, err := u.dashboardRepo.GetStats(ctx, cooperativeID)
	if err != nil {
		return nil, err
	}

	upcomingCount, err := u.dashboardRepo.GetUpcomingInstallmentsCount(ctx, cooperativeID)
	if err != nil {
		return nil, err
	}

	upcoming, err := u.dashboardRepo.GetUpcomingInstallments(ctx, cooperativeID)
	if err != nil {
		return nil, err
	}

	pendingCount, err := u.dashboardRepo.GetPendingApplicationsCount(ctx, cooperativeID)
	if err != nil {
		return nil, err
	}

	pending, err := u.dashboardRepo.GetPendingApplications(ctx, cooperativeID)
	if err != nil {
		return nil, err
	}

	if upcoming == nil {
		upcoming = []model.UpcomingInstallment{}
	}
	if pending == nil {
		pending = []model.PendingApplication{}
	}

	return &model.DashboardResponse{
		ActiveMembers:             stats.ActiveMembers,
		TotalMembers:              stats.TotalMembers,
		TotalSavings:              stats.TotalSavings,
		ActiveLoans:               stats.ActiveLoans,
		ActiveLoansOutstanding:    stats.ActiveLoansOutstanding,
		OverdueLoans:              stats.OverdueLoans,
		UpcomingInstallmentsCount: upcomingCount,
		UpcomingInstallments:      upcoming,
		PendingApplicationsCount:  pendingCount,
		PendingApplications:       pending,
	}, nil
}
