package usecase

import (
	"context"

	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
)

type DashboardUsecase interface {
	Get(ctx context.Context, coopID string) (*model.DashboardResponse, error)
}

type dashboardUsecase struct {
	repo repository.DashboardRepository
}

func NewDashboardUsecase(repo repository.DashboardRepository) DashboardUsecase {
	return &dashboardUsecase{repo: repo}
}

func (u *dashboardUsecase) Get(ctx context.Context, coopID string) (*model.DashboardResponse, error) {
	return u.repo.GetStats(ctx, coopID)
}
