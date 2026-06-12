package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/usecase"
)

type LoanController struct {
	uc usecase.LoanUsecase
}

func NewLoanController(uc usecase.LoanUsecase) *LoanController {
	return &LoanController{uc: uc}
}

func (ctrl *LoanController) List(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	status := c.Query("status")
	loans, err := ctrl.uc.List(c.Context(), coopID, status)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OKList(c, loans, int64(len(loans)))
}

func (ctrl *LoanController) GetByID(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	loanID := c.Params("id")
	loan, err := ctrl.uc.GetByID(c.Context(), coopID, loanID)
	if err != nil {
		return Fail(c, fiber.StatusNotFound, "NOT_FOUND", err.Error())
	}
	return OK(c, loan)
}
