package usecase_test

import (
	"math"
	"testing"
)

// shuJasaModal menghitung jasa modal SHU sesuai formula be_implementation.md §5.6.
// Fungsi ini tidak bergantung pada DB — memvalidasi rumus matematis saja.
func shuJasaModal(totalShu int64, pctJasaModal float64, simpananAnggota, totalSimpanan int64) int64 {
	shuModalPool := int64(math.Round(float64(totalShu) * pctJasaModal / 100))
	if totalSimpanan == 0 {
		return 0
	}
	return int64(math.Round(float64(simpananAnggota) / float64(totalSimpanan) * float64(shuModalPool)))
}

// TestSHU_JasaModal_ContohBeImplementation memvalidasi contoh angka di be_implementation.md §5.6:
// total_shu=40jt, modal 20%, simpanan anggota 8jt dari total 200jt → jasa_modal = 320.000
func TestSHU_JasaModal_ContohBeImplementation(t *testing.T) {
	got := shuJasaModal(40_000_000, 20, 8_000_000, 200_000_000)
	want := int64(320_000)
	if got != want {
		t.Errorf("jasa_modal: want %d, got %d", want, got)
	}
}

func TestSHU_JasaModal_ModalPool(t *testing.T) {
	// Validasi pool calculation: total_shu=10jt, pct=30% → pool=3jt
	pool := int64(math.Round(float64(10_000_000) * 30.0 / 100))
	if pool != 3_000_000 {
		t.Errorf("shu_modal_pool: want 3000000, got %d", pool)
	}
}

func TestSHU_JasaModal_TotalSimpananZero_ReturnsZero(t *testing.T) {
	got := shuJasaModal(40_000_000, 20, 8_000_000, 0)
	if got != 0 {
		t.Errorf("want 0 when totalSimpanan=0, got %d", got)
	}
}

func TestSHU_JasaUsaha_Formula(t *testing.T) {
	// jasa_usaha = round(jasaAnggota / totalJasa * shuUsahaPool)
	totalShu := int64(40_000_000)
	pctUsaha := 30.0
	jasaAnggota := int64(500_000)
	totalJasa := int64(5_000_000)

	shuUsahaPool := int64(math.Round(float64(totalShu) * pctUsaha / 100))
	got := int64(math.Round(float64(jasaAnggota) / float64(totalJasa) * float64(shuUsahaPool)))
	// 500_000/5_000_000 * 12_000_000 = 0.1 * 12_000_000 = 1_200_000
	want := int64(1_200_000)
	if got != want {
		t.Errorf("jasa_usaha: want %d, got %d", want, got)
	}
}
