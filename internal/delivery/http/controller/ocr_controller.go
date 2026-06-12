package controller

import (
	"errors"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wiradana/backend/internal/gateway/adins"
	"github.com/wiradana/backend/internal/usecase"
)

type OCRController struct {
	ocrUC usecase.OCRUsecase
	log   *logrus.Logger
}

func NewOCRController(ocrUC usecase.OCRUsecase, log *logrus.Logger) *OCRController {
	return &OCRController{ocrUC: ocrUC, log: log}
}

func (ctrl *OCRController) ExtractKTP(c *fiber.Ctx) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "field 'file' harus disertakan")
	}

	f, err := fh.Open()
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal membaca file")
	}
	defer f.Close()

	imgBytes, err := io.ReadAll(f)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal membaca file")
	}

	result, err := ctrl.ocrUC.ExtractKTP(c.Context(), imgBytes, fh.Filename)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrOCRFileTooLarge):
			return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case errors.Is(err, usecase.ErrOCRInvalidFileType):
			return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case errors.Is(err, adins.ErrOCRNotReadable):
			return Fail(c, fiber.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		case errors.Is(err, adins.ErrOCRInsufficientCredits):
			return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
		case errors.Is(err, adins.ErrOCRUnauthorized):
			ctrl.log.Errorf("ocr/ktp configuration error: %v", err)
			return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
		default:
			ctrl.log.Errorf("ocr/ktp error: %v", err)
			return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan saat memproses KTP")
		}
	}

	return OK(c, result)
}
