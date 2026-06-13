package controller

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type InstallmentController struct {
	uc       usecase.InstallmentUsecase
	validate *validator.Validate
}

func NewInstallmentController(uc usecase.InstallmentUsecase, validate *validator.Validate) *InstallmentController {
	return &InstallmentController{uc: uc, validate: validate}
}

func (ctrl *InstallmentController) Pay(c *fiber.Ctx) error {
	scheduleID := c.Params("id")
	var req model.PayInstallmentRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	userID, _ := c.Locals("user_id").(string)
	resp, err := ctrl.uc.Pay(c.Context(), scheduleID, &req, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrInstallmentNotFound) {
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "angsuran tidak ditemukan")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}

