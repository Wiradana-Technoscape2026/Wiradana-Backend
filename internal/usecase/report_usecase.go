package usecase

import (
	"context"

	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/repository"
)

type ReportUsecase interface {
	GetSummary(ctx context.Context, cooperativeID string) (*model.ReportSummary, error)
}

type reportUsecase struct {
	reportRepo repository.ReportRepository
}

func NewReportUsecase(reportRepo repository.ReportRepository) ReportUsecase {
	return &reportUsecase{reportRepo: reportRepo}
}

func (u *reportUsecase) GetSummary(ctx context.Context, cooperativeID string) (*model.ReportSummary, error) {
	return u.reportRepo.GetSummary(ctx, cooperativeID)
}
