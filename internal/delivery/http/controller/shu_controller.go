package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type ShuController struct {
	shuUC    usecase.ShuUsecase
	validate *validator.Validate
}

func NewShuController(shuUC usecase.ShuUsecase, validate *validator.Validate) *ShuController {
	return &ShuController{shuUC: shuUC, validate: validate}
}

func (ctrl *ShuController) List(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)

	periods, err := ctrl.shuUC.ListPeriods(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal mengambil periode SHU")
	}
	return OKList(c, periods, int64(len(periods)))
}

func (ctrl *ShuController) Create(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)

	var req model.CreateShuPeriodRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	resp, err := ctrl.shuUC.CreatePeriod(c.Context(), coopID, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal membuat periode SHU")
	}
	return OK(c, resp)
}

func (ctrl *ShuController) Calculate(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	periodID := c.Params("id")

	resp, err := ctrl.shuUC.Calculate(c.Context(), coopID, periodID)
	if err != nil {
		if err == usecase.ErrPeriodNotDraft {
			return Fail(c, fiber.StatusConflict, "CONFLICT", "periode SHU bukan dalam status draft")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal menghitung SHU")
	}
	return OK(c, resp)
}
