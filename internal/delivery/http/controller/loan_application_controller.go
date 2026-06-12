package controller

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type LoanApplicationController struct {
	uc       usecase.LoanApplicationUsecase
	notifUC  usecase.NotificationUsecase
	validate *validator.Validate
}

func NewLoanApplicationController(uc usecase.LoanApplicationUsecase, notifUC usecase.NotificationUsecase, validate *validator.Validate) *LoanApplicationController {
	return &LoanApplicationController{uc: uc, notifUC: notifUC, validate: validate}
}

func (ctrl *LoanApplicationController) List(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	status := c.Query("status")
	apps, err := ctrl.uc.List(c.Context(), coopID, status)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OKList(c, apps, int64(len(apps)))
}

func (ctrl *LoanApplicationController) Create(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	var req model.CreateLoanApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	app, err := ctrl.uc.Create(c.Context(), coopID, &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrMemberNotFound):
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "anggota tidak ditemukan")
		case errors.Is(err, usecase.ErrMemberNotAktif):
			return Fail(c, fiber.StatusConflict, "CONFLICT", err.Error())
		case errors.Is(err, usecase.ErrAmountExceedsPlafond):
			return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		default:
			return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
		}
	}
	return OK(c, app)
}

func (ctrl *LoanApplicationController) Approve(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	userID := c.Locals("user_id").(string)
	appID := c.Params("id")
	resp, err := ctrl.uc.Approve(c.Context(), coopID, appID, userID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrLoanApplicationNotFound):
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "pengajuan tidak ditemukan")
		case errors.Is(err, usecase.ErrApplicationNotPending):
			return Fail(c, fiber.StatusConflict, "CONFLICT", err.Error())
		default:
			return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
		}
	}

	if len(resp.Loan.Schedule) > 0 {
		if firstDue, parseErr := time.Parse("2006-01-02", resp.Loan.Schedule[0].DueDate); parseErr == nil {
			go ctrl.notifUC.SendLoanDisbursed(context.Background(), coopID, resp.Loan.MemberID, resp.Loan.ID, resp.Loan.Principal, firstDue)
		}
	}

	return OK(c, resp)
}

func (ctrl *LoanApplicationController) Reject(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	appID := c.Params("id")
	resp, err := ctrl.uc.Reject(c.Context(), coopID, appID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrLoanApplicationNotFound):
			return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "pengajuan tidak ditemukan")
		case errors.Is(err, usecase.ErrApplicationNotPending):
			return Fail(c, fiber.StatusConflict, "CONFLICT", err.Error())
		default:
			return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
		}
	}
	return OK(c, resp)
}
