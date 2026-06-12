package controller

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type ShuController struct {
	uc       usecase.ShuUsecase
	validate *validator.Validate
}

func NewShuController(uc usecase.ShuUsecase, validate *validator.Validate) *ShuController {
	return &ShuController{uc: uc, validate: validate}
}

func (ctrl *ShuController) ListPeriods(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	periods, err := ctrl.uc.ListPeriods(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OKList(c, periods, int64(len(periods)))
}

func (ctrl *ShuController) CreatePeriod(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	var req model.CreateShuPeriodRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	resp, err := ctrl.uc.CreatePeriod(c.Context(), coopID, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}

func (ctrl *ShuController) Calculate(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	periodID := c.Params("id")
	resp, err := ctrl.uc.Calculate(c.Context(), coopID, periodID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrShuPeriodNotFound):
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", err.Error())
		case errors.Is(err, usecase.ErrShuPeriodNotDraft):
			return Fail(c, fiber.StatusConflict, "CONFLICT", err.Error())
		default:
			return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
		}
	}
	return OK(c, resp)
}
