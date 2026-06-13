package adins

import (
	"context"
	"math"
)

type MockScoringGateway struct{}

func NewMockScoringGateway() ScoringGateway {
	return &MockScoringGateway{}
}

func (g *MockScoringGateway) Score(_ context.Context, in ScoringInput) (ScoringResult, error) {
	f := in.Features
	raw := 100 * (
		0.35*f["ketepatan_bayar"] +
			0.25*math.Min(f["rasio_simpanan_pinjaman"], 1) +
			0.15*math.Min(f["lama_keanggotaan_hari"]/1080, 1) +
			0.15*f["konsistensi_simpanan"] +
			0.10*(1-math.Min(f["rasio_beban_angsuran"], 1)))
	score := int(math.Round(raw))

	grade, rec := scoreToGrade(score)
	limit := limitFromGrade(grade, in.TotalSimpanan, in.MaxPlafond)
	reasons := generateReasons(f)

	return ScoringResult{
		Score: score, Grade: grade, Recommendation: rec,
		LimitRekomendasi: limit, Reasons: reasons, Source: "MOCK_ADINS_SCORING",
	}, nil
}

func scoreToGrade(score int) (string, string) {
	switch {
	case score >= 80:
		return "A", "approve"
	case score >= 70:
		return "B", "approve"
	case score >= 60:
		return "C", "review"
	default:
		return "D", "reject"
	}
}

func limitFromGrade(grade string, totalSimpanan, maxPlafond int64) int64 {
	factors := map[string]float64{"A": 3.0, "B": 2.5, "C": 1.5, "D": 0}
	limit := int64(math.Round(factors[grade] * float64(totalSimpanan)))
	if maxPlafond > 0 && limit > maxPlafond {
		limit = maxPlafond
	}
	return limit
}

func generateReasons(f map[string]float64) []string {
	var reasons []string
	if f["ketepatan_bayar"] >= 0.9 {
		reasons = append(reasons, "riwayat pembayaran sangat baik")
	} else if f["ketepatan_bayar"] < 0.6 {
		reasons = append(reasons, "riwayat pembayaran perlu diperhatikan")
	}
	if f["rasio_simpanan_pinjaman"] >= 0.5 {
		reasons = append(reasons, "simpanan memadai terhadap pinjaman")
	} else if f["rasio_simpanan_pinjaman"] < 0.2 {
		reasons = append(reasons, "simpanan relatif kecil dibanding pinjaman")
	}
	if f["lama_keanggotaan_hari"] >= 720 {
		reasons = append(reasons, "anggota lama dan loyal")
	}
	if f["rasio_beban_angsuran"] > 0.7 {
		reasons = append(reasons, "beban angsuran tinggi")
	}
	if len(reasons) == 0 {
		reasons = []string{"profil risiko cukup baik"}
	}
	return reasons
}
