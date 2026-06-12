package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type PortalLoanController struct {
	loanAppUC usecase.LoanApplicationUsecase
	loanUC    usecase.LoanUsecase
	validate  *validator.Validate
}

func NewPortalLoanController(loanAppUC usecase.LoanApplicationUsecase, loanUC usecase.LoanUsecase, validate *validator.Validate) *PortalLoanController {
	return &PortalLoanController{loanAppUC: loanAppUC, loanUC: loanUC, validate: validate}
}

func (ctrl *PortalLoanController) ListApplications(c *fiber.Ctx) error {
	memberID := c.Locals("member_id").(string)
	apps, err := ctrl.loanAppUC.ListForMember(c.Context(), memberID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OKList(c, apps, int64(len(apps)))
}

func (ctrl *PortalLoanController) Apply(c *fiber.Ctx) error {
	memberID := c.Locals("member_id").(string)
	coopID := c.Locals("cooperative_id").(string)
	var req model.CreatePortalLoanApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}
	app, err := ctrl.loanAppUC.CreateForMember(c.Context(), memberID, coopID, &req)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, app)
}

func (ctrl *PortalLoanController) ListLoans(c *fiber.Ctx) error {
	memberID := c.Locals("member_id").(string)
	loans, err := ctrl.loanUC.ListForMember(c.Context(), memberID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OKList(c, loans, int64(len(loans)))
}

func (ctrl *PortalLoanController) GetLoan(c *fiber.Ctx) error {
	memberID := c.Locals("member_id").(string)
	loanID := c.Params("id")
	loan, err := ctrl.loanUC.GetByIDForMember(c.Context(), memberID, loanID)
	if err != nil {
		msg := err.Error()
		if msg == "forbidden" {
			return Fail(c, fiber.StatusForbidden, "FORBIDDEN", "bukan pinjaman Anda")
		}
		return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "pinjaman tidak ditemukan")
	}
	return OK(c, loan)
}
