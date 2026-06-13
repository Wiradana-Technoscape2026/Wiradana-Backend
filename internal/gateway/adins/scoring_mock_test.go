package adins_test

import (
	"context"
	"testing"

	"github.com/wiradana/backend/internal/gateway/adins"
)

func TestMockScoringGateway_Example(t *testing.T) {
	gw := adins.NewMockScoringGateway()
	result, err := gw.Score(context.Background(), adins.ScoringInput{
		MemberID:       "test",
		JumlahDiajukan: 5000000,
		TenorBulan:     12,
		TotalSimpanan:  3000000,
		MaxPlafond:     20000000,
		Features: map[string]float64{
			"ketepatan_bayar":         0.95,
			"rasio_simpanan_pinjaman": 0.6,
			"lama_keanggotaan_hari":  840,
			"konsistensi_simpanan":    0.9,
			"rasio_beban_angsuran":    0.3,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Score < 78 || result.Score > 82 {
		t.Errorf("want score≈80, got %d", result.Score)
	}
	if result.Grade != "A" {
		t.Errorf("want grade A, got %s", result.Grade)
	}
	if result.Recommendation != "approve" {
		t.Errorf("want approve, got %s", result.Recommendation)
	}
	if result.LimitRekomendasi != 9_000_000 {
		t.Errorf("want limit 9000000, got %d", result.LimitRekomendasi)
	}
}

func TestMockScoringGateway_GradeD(t *testing.T) {
	gw := adins.NewMockScoringGateway()
	result, _ := gw.Score(context.Background(), adins.ScoringInput{
		TotalSimpanan: 100000,
		MaxPlafond:    20000000,
		Features: map[string]float64{
			"ketepatan_bayar":         0.2,
			"rasio_simpanan_pinjaman": 0.05,
			"lama_keanggotaan_hari":  30,
			"konsistensi_simpanan":    0.1,
			"rasio_beban_angsuran":    0.95,
		},
	})
	if result.Grade != "D" {
		t.Errorf("want grade D, got %s", result.Grade)
	}
	if result.Recommendation != "reject" {
		t.Errorf("want reject, got %s", result.Recommendation)
	}
}
