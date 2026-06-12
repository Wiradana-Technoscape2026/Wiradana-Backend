package adins

import (
	"context"
	"errors"
)

var (
	ErrOCRNotReadable         = errors.New("KTP tidak terbaca, upload foto yang lebih jelas")
	ErrOCRInsufficientCredits = errors.New("OCR credits habis, hubungi administrator")
	ErrOCRUnauthorized        = errors.New("API key OCR tidak valid atau belum dikonfigurasi")
)

type KTPOCRResult struct {
	NIK        string  `json:"nik"`
	FullName   string  `json:"full_name"`
	Address    string  `json:"address"`
	Pekerjaan  string  `json:"pekerjaan"`
	BirthDate  string  `json:"birth_date"` // "YYYY-MM-DD"
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"`
}

type KTPOCRGateway interface {
	ExtractKTP(ctx context.Context, image []byte, filename string) (*KTPOCRResult, error)
}

type ScoringInput struct {
	MemberID       string             `json:"member_id"`
	Features       map[string]float64 `json:"features"`
	JumlahDiajukan int64              `json:"jumlah_diajukan"`
	TenorBulan     int                `json:"tenor_bulan"`
	TotalSimpanan  int64              `json:"total_simpanan"`
	MaxPlafond     int64              `json:"max_plafond"`
}

type ScoringResult struct {
	Score            int      `json:"score"`
	Grade            string   `json:"grade"`
	Recommendation   string   `json:"recommendation"`
	LimitRekomendasi int64    `json:"limit_rekomendasi"`
	Reasons          []string `json:"reasons"`
	Source           string   `json:"source"`
}

type ScoringGateway interface {
	Score(ctx context.Context, in ScoringInput) (ScoringResult, error)
}
