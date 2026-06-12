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
	BirthDate  string  `json:"birth_date"` // "YYYY-MM-DD"
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"`
}

type KTPOCRGateway interface {
	ExtractKTP(ctx context.Context, image []byte, filename string) (*KTPOCRResult, error)
}
