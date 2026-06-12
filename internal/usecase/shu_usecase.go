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
	ErrShuPeriodNotFound = errors.New("periode SHU tidak ditemukan")
	ErrShuPeriodNotDraft = errors.New("periode SHU bukan dalam status draft")
)

type ShuUsecase interface {
	CreatePeriod(ctx context.Context, coopID string, req *model.CreateShuPeriodRequest) (*model.ShuPeriodResponse, error)
	ListPeriods(ctx context.Context, coopID string) ([]model.ShuPeriodResponse, error)
	Calculate(ctx context.Context, coopID, periodID string) (*model.CalculateShuResponse, error)
	GetForMember(ctx context.Context, memberID string) (*model.PortalSHUResponse, error)
}

type shuUsecase struct {
	shuRepo repository.ShuRepository
}

func NewShuUsecase(shuRepo repository.ShuRepository) ShuUsecase {
	return &shuUsecase{shuRepo: shuRepo}
}

func (u *shuUsecase) CreatePeriod(ctx context.Context, coopID string, req *model.CreateShuPeriodRequest) (*model.ShuPeriodResponse, error) {
	coopUUID, _ := uuid.Parse(coopID)
	sp := &entity.ShuPeriod{
		CooperativeID: coopUUID,
		Year:          req.Year,
		TotalShu:      req.TotalShu,
		PctJasaModal:  req.PctJasaModal,
		PctJasaUsaha:  req.PctJasaUsaha,
		Status:        "draft",
	}
	if err := u.shuRepo.CreatePeriod(ctx, sp); err != nil {
		return nil, err
	}
	resp := converter.ToShuPeriodResponse(sp)
	return &resp, nil
}

func (u *shuUsecase) ListPeriods(ctx context.Context, coopID string) ([]model.ShuPeriodResponse, error) {
	periods, err := u.shuRepo.FindPeriods(ctx, coopID)
	if err != nil {
		return nil, err
	}
	result := make([]model.ShuPeriodResponse, len(periods))
	for i, p := range periods {
		result[i] = converter.ToShuPeriodResponse(&p)
	}
	return result, nil
}

func (u *shuUsecase) Calculate(ctx context.Context, coopID, periodID string) (*model.CalculateShuResponse, error) {
	period, err := u.shuRepo.FindPeriodByID(ctx, periodID)
	if err != nil {
		if errors.Is(err, repository.ErrShuPeriodNotFound) {
			return nil, ErrShuPeriodNotFound
		}
		return nil, err
	}
	if period.Status != "draft" {
		return nil, ErrShuPeriodNotDraft
	}

	// Calculate pools
	shuModalPool := int64(math.Round(float64(period.TotalShu) * period.PctJasaModal / 100))
	shuUsahaPool := int64(math.Round(float64(period.TotalShu) * period.PctJasaUsaha / 100))

	// Get total simpanan aktif (pokok + wajib)
	totalSimpanan, _ := u.shuRepo.GetTotalSimpananAktif(ctx, coopID)
	totalJasaPinjaman, _ := u.shuRepo.GetTotalJasaPinjaman(ctx, coopID)

	// Get per-member data
	simpananMap, _ := u.shuRepo.GetSimpananPerMember(ctx, coopID)
	jasaPinjamanMap, _ := u.shuRepo.GetJasaPinjamanPerMember(ctx, coopID)

	// Build lookup maps
	simpananLookup := make(map[string]int64, len(simpananMap))
	for _, r := range simpananMap {
		simpananLookup[r.MemberID] = r.TotalSimpanan
	}
	jasaLookup := make(map[string]int64, len(jasaPinjamanMap))
	for _, r := range jasaPinjamanMap {
		jasaLookup[r.MemberID] = r.TotalJasaPinjaman
	}

	// Get all active members
	memberIDs, _ := u.shuRepo.GetActiveMemberIDs(ctx, coopID)

	// Calculate distribution per member
	periodUUID, _ := uuid.Parse(periodID)
	distributions := make([]entity.ShuDistribution, 0, len(memberIDs))

	for _, mid := range memberIDs {
		simpananAnggota := simpananLookup[mid]
		jasaAnggota := jasaLookup[mid]

		var jasaModal int64
		if totalSimpanan > 0 {
			jasaModal = int64(math.Round(float64(simpananAnggota) / float64(totalSimpanan) * float64(shuModalPool)))
		}

		var jasaUsaha int64
		if totalJasaPinjaman > 0 {
			jasaUsaha = int64(math.Round(float64(jasaAnggota) / float64(totalJasaPinjaman) * float64(shuUsahaPool)))
		}

		totalShu := jasaModal + jasaUsaha
		memberUUID, _ := uuid.Parse(mid)

		distributions = append(distributions, entity.ShuDistribution{
			ShuPeriodID: periodUUID,
			MemberID:    memberUUID,
			JasaModal:   jasaModal,
			JasaUsaha:   jasaUsaha,
			TotalShu:    totalShu,
		})
	}

	// Bulk insert
	if err := u.shuRepo.CreateDistributions(ctx, distributions); err != nil {
		return nil, err
	}

	// Update period status to final
	if err := u.shuRepo.UpdatePeriodStatus(ctx, periodID, "final"); err != nil {
		return nil, err
	}

	// Build response with member names (fetched via JOIN)
	distsWithMeta, _ := u.shuRepo.FindDistributionsByPeriod(ctx, periodID)
	distResponses := make([]model.ShuDistributionResponse, len(distsWithMeta))
	for i, d := range distsWithMeta {
		distResponses[i] = converter.ToShuDistributionResponse(&d.ShuDistribution, d.MemberName)
	}

	updatedPeriod, _ := u.shuRepo.FindPeriodByID(ctx, periodID)
	periodResp := converter.ToShuPeriodResponse(updatedPeriod)

	return &model.CalculateShuResponse{
		Period:        periodResp,
		Distributions: distResponses,
	}, nil
}

func (u *shuUsecase) GetForMember(ctx context.Context, memberID string) (*model.PortalSHUResponse, error) {
	dists, err := u.shuRepo.FindDistributionsByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}

	history := make([]model.ShuDistributionResponse, len(dists))
	var estimatedShu int64
	for i, d := range dists {
		history[i] = converter.ToShuDistributionResponse(&d.ShuDistribution, d.MemberName)
	}

	// If no history, estimated_shu = 0
	if len(history) > 0 {
		estimatedShu = history[len(history)-1].TotalShu
	}

	return &model.PortalSHUResponse{
		EstimatedShu: estimatedShu,
		History:      history,
	}, nil
}
