package controller

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type SavingsController struct {
	savingsUC usecase.SavingsUsecase
	notifUC   usecase.NotificationUsecase
	validate  *validator.Validate
}

func NewSavingsController(savingsUC usecase.SavingsUsecase, notifUC usecase.NotificationUsecase, validate *validator.Validate) *SavingsController {
	return &SavingsController{savingsUC: savingsUC, notifUC: notifUC, validate: validate}
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
}

func (ctrl *SavingsController) Record(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	userID := c.Locals("user_id").(string)
	memberID := c.Params("id")

	var req model.CreateSavingsRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	tx, err := ctrl.savingsUC.Record(c.Context(), coopID, memberID, userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrMemberNotFound):
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "anggota tidak ditemukan")
		case errors.Is(err, usecase.ErrPokokAlreadyRecorded):
			return Fail(c, fiber.StatusConflict, "CONFLICT", "simpanan pokok sudah pernah disetor")
		case errors.Is(err, usecase.ErrCannotWithdrawMandatory):
			return Fail(c, fiber.StatusConflict, "CONFLICT", "simpanan pokok dan wajib tidak dapat ditarik")
		case errors.Is(err, usecase.ErrInsufficientBalance):
			return Fail(c, fiber.StatusConflict, "CONFLICT", "saldo sukarela tidak mencukupi")
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}

	if req.Direction == "setor" {
		go ctrl.notifUC.SendSavingsConfirmation(context.Background(), coopID, memberID, tx.ID, tx.Amount, tx.SavingsType)
	}

	return OK(c, tx)
}

func (ctrl *SavingsController) Summary(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	resp, err := ctrl.savingsUC.GetCoopSummary(c.Context(), coopID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}
	return OK(c, resp)
}

func (ctrl *SavingsController) ListAll(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	savingsType := c.Query("type")
	direction := c.Query("direction")
	limit := c.QueryInt("limit", 20)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	txs, total, err := ctrl.savingsUC.FindAllRecent(c.Context(), coopID, savingsType, direction, limit, offset)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}
	return OKList(c, txs, total)
}
