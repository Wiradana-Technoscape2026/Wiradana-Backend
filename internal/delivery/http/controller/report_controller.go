package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/usecase"
)

type ReportController struct {
	reportUC usecase.ReportUsecase
}

func NewReportController(reportUC usecase.ReportUsecase) *ReportController {
	return &ReportController{reportUC: reportUC}
}

func (ctrl *ReportController) Summary(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)

	summary, err := ctrl.reportUC.GetSummary(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal mengambil ringkasan laporan")
	}
	return OK(c, summary)
}
