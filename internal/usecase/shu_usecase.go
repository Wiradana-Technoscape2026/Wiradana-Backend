package usecase

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/repository"
)

var (
	ErrPeriodNotDraft = errors.New("periode SHU bukan dalam status draft")
)

type ShuUsecase interface {
	CreatePeriod(ctx context.Context, coopID string, req *model.CreateShuPeriodRequest) (*model.ShuPeriodResponse, error)
	ListPeriods(ctx context.Context, coopID string) ([]model.ShuPeriodResponse, error)
	Calculate(ctx context.Context, coopID, periodID string) (*model.CalculateShuResponse, error)
	GetMemberDistributions(ctx context.Context, memberID string) ([]model.ShuDistributionDetail, error)
}

type shuUsecase struct {
	shuRepo repository.ShuRepository
}

func NewShuUsecase(shuRepo repository.ShuRepository) ShuUsecase {
	return &shuUsecase{shuRepo: shuRepo}
}

func (u *shuUsecase) CreatePeriod(ctx context.Context, coopID string, req *model.CreateShuPeriodRequest) (*model.ShuPeriodResponse, error) {
	coopUUID, err := uuid.Parse(coopID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}

	period := &entity.ShuPeriod{
		CooperativeID: coopUUID,
		Year:          req.Year,
		TotalShu:      req.TotalShu,
		PctJasaModal:  req.PctJasaModal,
		PctJasaUsaha:  req.PctJasaUsaha,
		Status:        "draft",
	}

	if err := u.shuRepo.CreatePeriod(ctx, period); err != nil {
		return nil, err
	}

	resp := converter.ToShuPeriodResponse(period)
	return &resp, nil
}

func (u *shuUsecase) ListPeriods(ctx context.Context, coopID string) ([]model.ShuPeriodResponse, error) {
	periods, err := u.shuRepo.FindPeriods(ctx, coopID)
	if err != nil {
		return nil, err
	}

	result := make([]model.ShuPeriodResponse, 0, len(periods))
	for _, p := range periods {
		result = append(result, converter.ToShuPeriodResponse(&p))
	}
	return result, nil
}

func (u *shuUsecase) Calculate(ctx context.Context, coopID, periodID string) (*model.CalculateShuResponse, error) {
	period, err := u.shuRepo.FindPeriodByID(ctx, coopID, periodID)
	if err != nil {
		if errors.Is(err, repository.ErrShuPeriodNotFound) {
			return nil, errors.New("periode SHU tidak ditemukan")
		}
		return nil, err
	}

	if period.Status != "draft" {
		return nil, ErrPeriodNotDraft
	}

	// Calculate pools
	shuModalPool := int64(math.Round(float64(period.TotalShu) * period.PctJasaModal / 100))
	shuUsahaPool := int64(math.Round(float64(period.TotalShu) * period.PctJasaUsaha / 100))

	// Get totals
	totalSimpanan, _ := u.shuRepo.GetTotalSimpananAktif(ctx, coopID)
	totalJasa, _ := u.shuRepo.GetTotalJasaPinjaman(ctx, coopID, period.Year)

	// Get active members
	members, err := u.shuRepo.GetActiveMembers(ctx, coopID)
	if err != nil {
		return nil, err
	}

	distributions := make([]entity.ShuDistribution, 0, len(members))

	for _, member := range members {
		simpananAnggota, _ := u.shuRepo.GetSimpananMember(ctx, member.ID.String())
		jasaAnggota, _ := u.shuRepo.GetJasaPinjamanMember(ctx, member.ID.String(), period.Year)

		var jasaModal int64
		var jasaUsaha int64

		if totalSimpanan > 0 {
			jasaModal = int64(math.Round(float64(simpananAnggota) / float64(totalSimpanan) * float64(shuModalPool)))
		}

		if totalJasa > 0 && jasaAnggota > 0 {
			jasaUsaha = int64(math.Round(float64(jasaAnggota) / float64(totalJasa) * float64(shuUsahaPool)))
		}

		dist := entity.ShuDistribution{
			ShuPeriodID: period.ID,
			MemberID:    member.ID,
			JasaModal:   jasaModal,
			JasaUsaha:   jasaUsaha,
			TotalShu:    jasaModal + jasaUsaha,
		}
		distributions = append(distributions, dist)
	}

	// Bulk insert distributions
	if err := u.shuRepo.CreateDistributions(ctx, distributions); err != nil {
		return nil, err
	}

	// Update period status to final
	if err := u.shuRepo.UpdatePeriodStatus(ctx, periodID, "final"); err != nil {
		return nil, err
	}
	period.Status = "final"

	// Fetch distributions with member names
	distWithNames, err := u.shuRepo.FindDistributionsByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}

	distResponses := make([]model.ShuDistributionDetail, 0, len(distWithNames))
	for _, d := range distWithNames {
		distResponses = append(distResponses, model.ShuDistributionDetail{
			ID:          d.ID.String(),
			ShuPeriodID: d.ShuPeriodID.String(),
			MemberID:    d.MemberID.String(),
			MemberName:  d.MemberName,
			JasaModal:   d.JasaModal,
			JasaUsaha:   d.JasaUsaha,
			TotalShu:    d.TotalShu,
		})
	}

	return &model.CalculateShuResponse{
		Period:        converter.ToShuPeriodResponse(period),
		Distributions: distResponses,
	}, nil
}

func (u *shuUsecase) GetMemberDistributions(ctx context.Context, memberID string) ([]model.ShuDistributionDetail, error) {
	dists, err := u.shuRepo.FindDistributionsByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}

	result := make([]model.ShuDistributionDetail, 0, len(dists))
	for _, d := range dists {
		result = append(result, model.ShuDistributionDetail{
			ID:          d.ID.String(),
			ShuPeriodID: d.ShuPeriodID.String(),
			MemberID:    d.MemberID.String(),
			MemberName:  d.MemberName,
			JasaModal:   d.JasaModal,
			JasaUsaha:   d.JasaUsaha,
			TotalShu:    d.TotalShu,
		})
	}
	return result, nil
}
