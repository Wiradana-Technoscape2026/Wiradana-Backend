package controller

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type SavingsController struct {
<<<<<<< HEAD
	uc       usecase.SavingsUsecase
	validate *validator.Validate
}

func NewSavingsController(uc usecase.SavingsUsecase, validate *validator.Validate) *SavingsController {
	return &SavingsController{uc: uc, validate: validate}
=======
	savingsUC usecase.SavingsUsecase
	validate  *validator.Validate
}

func NewSavingsController(savingsUC usecase.SavingsUsecase, validate *validator.Validate) *SavingsController {
	return &SavingsController{savingsUC: savingsUC, validate: validate}
}

func (ctrl *SavingsController) List(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	memberID := c.Params("id")

	txs, err := ctrl.savingsUC.FindByMember(c.Context(), coopID, memberID)
	if err != nil {
		if errors.Is(err, usecase.ErrMemberNotFound) {
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "anggota tidak ditemukan")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}

	return OKList(c, txs, int64(len(txs)))
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
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

<<<<<<< HEAD
	result, err := ctrl.uc.Record(c.Context(), coopID, memberID, &req)
	if err != nil {
		switch {
=======
	tx, err := ctrl.savingsUC.Record(c.Context(), coopID, memberID, &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrMemberNotFound):
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "anggota tidak ditemukan")
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
		case errors.Is(err, usecase.ErrPokokAlreadyRecorded):
			return Fail(c, fiber.StatusConflict, "CONFLICT", "simpanan pokok sudah pernah disetor")
		case errors.Is(err, usecase.ErrCannotWithdrawMandatory):
			return Fail(c, fiber.StatusConflict, "CONFLICT", "simpanan pokok dan wajib tidak dapat ditarik")
<<<<<<< HEAD
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
=======
		case errors.Is(err, usecase.ErrInsufficientBalance):
			return Fail(c, fiber.StatusConflict, "CONFLICT", "saldo sukarela tidak mencukupi")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}

	return OK(c, tx)
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
}
