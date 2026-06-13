package usecase

import (
	"context"

	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
)

type DashboardUsecase interface {
	Get(ctx context.Context, cooperativeID string) (*model.DashboardResponse, error)
	GetMemberDashboard(ctx context.Context, cooperativeID, memberID string) (*model.MemberDashboardResponse, error)
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

func (u *dashboardUsecase) GetMemberDashboard(ctx context.Context, cooperativeID, memberID string) (*model.MemberDashboardResponse, error) {
	stats, err := u.dashboardRepo.GetMemberStats(ctx, cooperativeID, memberID)
	if err != nil {
		return nil, err
	}

	upcoming, err := u.dashboardRepo.GetMemberUpcomingInstallments(ctx, cooperativeID, memberID)
	if err != nil {
		return nil, err
	}
	if upcoming == nil {
		upcoming = []model.MemberUpcomingInstallment{}
	}

	total := stats.Pokok + stats.Wajib + stats.Sukarela

	return &model.MemberDashboardResponse{
		CooperativeID:   cooperativeID,
		CooperativeName: stats.CooperativeName,
		MemberID:        memberID,
		MemberName:      stats.MemberName,
		SavingsSummary: model.SavingsSummary{
			Pokok:    stats.Pokok,
			Wajib:    stats.Wajib,
			Sukarela: stats.Sukarela,
			Total:    total,
		},
		ActiveLoans:          stats.ActiveLoans,
		OutstandingAmount:    stats.OutstandingAmount,
		OverdueInstallments:  stats.OverdueInstallments,
		UpcomingInstallments: upcoming,
		EstimatedShu:         stats.EstimatedShu,
	}, nil
}
