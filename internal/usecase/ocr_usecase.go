package usecase

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/wiradana/backend/internal/gateway/adins"
	"github.com/wiradana/backend/internal/model"
)

var (
	ErrOCRInvalidFileType = errors.New("hanya file JPG atau PNG yang diizinkan")
	ErrOCRFileTooLarge    = errors.New("ukuran file melebihi batas 5 MB")
)

const maxOCRFileSize = 5 * 1024 * 1024 // 5 MB

type OCRUsecase interface {
	ExtractKTP(ctx context.Context, image []byte, filename string) (*model.KTPOCRResponse, error)
}

type ocrUsecase struct {
	gateway adins.KTPOCRGateway
}

func NewOCRUsecase(gateway adins.KTPOCRGateway) OCRUsecase {
	return &ocrUsecase{gateway: gateway}
}

func (u *ocrUsecase) ExtractKTP(ctx context.Context, image []byte, filename string) (*model.KTPOCRResponse, error) {
	if len(image) > maxOCRFileSize {
		return nil, ErrOCRFileTooLarge
	}
	if !isValidImageType(image) || !isValidImageExtension(filename) {
		return nil, ErrOCRInvalidFileType
	}

	result, err := u.gateway.ExtractKTP(ctx, image, filename)
	if err != nil {
		return nil, err
	}

	return &model.KTPOCRResponse{
		NIK:        result.NIK,
		FullName:   result.FullName,
		Address:    result.Address,
		BirthDate:  result.BirthDate,
		Confidence: result.Confidence,
		Source:     result.Source,
	}, nil
}

func isValidImageType(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	// JPEG magic: FF D8 FF
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return true
	}
	// PNG magic: 89 50 4E 47
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return true
	}
	return false
}

func isValidImageExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}
