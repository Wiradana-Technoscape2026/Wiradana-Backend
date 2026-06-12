package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type LoanConfigController struct {
	uc       usecase.LoanConfigUsecase
	validate *validator.Validate
}

func NewLoanConfigController(uc usecase.LoanConfigUsecase, validate *validator.Validate) *LoanConfigController {
	return &LoanConfigController{uc: uc, validate: validate}
}

func (ctrl *LoanConfigController) Get(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	resp, err := ctrl.uc.Get(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}

func (ctrl *LoanConfigController) Update(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	var req model.UpdateLoanConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	resp, err := ctrl.uc.Update(c.Context(), coopID, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}
