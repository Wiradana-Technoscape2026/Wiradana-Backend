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

	notifs, err := u.dashboardRepo.GetUpcomingNotifications(ctx, cooperativeID)
	if err != nil {
		return nil, err
	}

	if notifs == nil {
		notifs = []model.DashboardNotification{}
	}

	return &model.DashboardResponse{
		TotalMembers:  stats.TotalMembers,
		TotalSavings:  stats.TotalSavings,
		ActiveLoans:   stats.ActiveLoans,
		OverdueLoans:  stats.OverdueLoans,
		Notifications: notifs,
	}, nil
}
