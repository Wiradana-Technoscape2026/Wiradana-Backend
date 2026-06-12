package controller

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type SavingsController struct {
	uc       usecase.SavingsUsecase
	validate *validator.Validate
}

func NewSavingsController(uc usecase.SavingsUsecase, validate *validator.Validate) *SavingsController {
	return &SavingsController{uc: uc, validate: validate}
}

func (ctrl *SavingsController) Record(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	memberID := c.Params("id")

	var req model.CreateSavingsRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	result, err := ctrl.uc.Record(c.Context(), coopID, memberID, &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrPokokAlreadyRecorded):
			return Fail(c, fiber.StatusConflict, "CONFLICT", "simpanan pokok sudah pernah disetor")
		case errors.Is(err, usecase.ErrCannotWithdrawMandatory):
			return Fail(c, fiber.StatusConflict, "CONFLICT", "simpanan pokok dan wajib tidak dapat ditarik")
		case errors.Is(err, usecase.ErrInsufficientSukarela):
			return Fail(c, fiber.StatusConflict, "CONFLICT", "saldo sukarela tidak mencukupi")
		default:
			return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
		}
	}
	return OK(c, result)
}

func (ctrl *SavingsController) List(c *fiber.Ctx) error {
	memberID := c.Params("id")
	txs, err := ctrl.uc.ListByMember(c.Context(), memberID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OKList(c, txs, int64(len(txs)))
}
