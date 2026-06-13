package controller

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type LoanAuditController struct {
	uc       usecase.LoanAuditUsecase
	validate *validator.Validate
}

func NewLoanAuditController(uc usecase.LoanAuditUsecase, validate *validator.Validate) *LoanAuditController {
	return &LoanAuditController{uc: uc, validate: validate}
}

func (ctrl *LoanAuditController) CreateToken(c *fiber.Ctx) error {
	loanID := c.Params("id")
	coopID, _ := c.Locals("cooperative_id").(string)
	userID, _ := c.Locals("user_id").(string)

	var req model.CreateAuditTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	resp, err := ctrl.uc.CreateToken(c.Context(), coopID, loanID, userID, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}

func (ctrl *LoanAuditController) GetAuditDetails(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "token audit diperlukan")
	}

	resp, err := ctrl.uc.GetAuditDetails(c.Context(), token)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidOrExpiredToken) {
			return Fail(c, fiber.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, resp)
}

func (ctrl *LoanAuditController) FlagLog(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "token audit diperlukan")
	}

	var req model.FlagAuditLogRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	err := ctrl.uc.FlagLog(c.Context(), token, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidOrExpiredToken) {
			return Fail(c, fiber.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		}
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, map[string]interface{}{"message": "audit log berhasil ditandai"})
}
