package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/usecase"
)

type DashboardController struct {
	dashboardUC usecase.DashboardUsecase
}

func NewDashboardController(dashboardUC usecase.DashboardUsecase) *DashboardController {
	return &DashboardController{dashboardUC: dashboardUC}
}

func (ctrl *DashboardController) Get(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)

	resp, err := ctrl.dashboardUC.Get(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}

	return OK(c, resp)
}
