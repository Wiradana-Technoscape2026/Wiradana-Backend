package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/usecase"
)

type DashboardController struct {
	uc usecase.DashboardUsecase
}

func NewDashboardController(uc usecase.DashboardUsecase) *DashboardController {
	return &DashboardController{uc: uc}
}

func (ctrl *DashboardController) Get(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	resp, err := ctrl.uc.Get(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}
