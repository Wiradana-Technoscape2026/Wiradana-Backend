package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type PortalController struct {
	memberUC  usecase.MemberUsecase
	savingsUC usecase.SavingsUsecase
	shuUC     usecase.ShuUsecase
}

func NewPortalController(
	memberUC usecase.MemberUsecase,
	savingsUC usecase.SavingsUsecase,
	shuUC usecase.ShuUsecase,
) *PortalController {
	return &PortalController{
		memberUC:  memberUC,
		savingsUC: savingsUC,
		shuUC:     shuUC,
	}
}

// PortalMeResponse — GET /portal/me (api_planning §3.8)
type PortalMeResponse struct {
	Member         model.MemberResponse              `json:"member"`
	SavingsSummary model.SavingsSummary              `json:"savings_summary"`
	SavingsRecent  []model.SavingsTransactionResponse `json:"savings_recent"`
}

func (ctrl *PortalController) Me(c *fiber.Ctx) error {
	memberID := c.Locals("member_id").(string)
	coopID := c.Locals("cooperative_id").(string)

	member, err := ctrl.memberUC.FindByID(c.Context(), coopID, memberID)
	if err != nil {
		return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "data anggota tidak ditemukan")
	}

	recentTxs, _ := ctrl.savingsUC.ListByMember(c.Context(), memberID)
	if len(recentTxs) > 5 {
		recentTxs = recentTxs[:5]
	}

	return OK(c, PortalMeResponse{
		Member:         *member,
		SavingsSummary: member.SavingsSummary,
		SavingsRecent:  recentTxs,
	})
}

func (ctrl *PortalController) SHU(c *fiber.Ctx) error {
	memberID := c.Locals("member_id").(string)
	result, err := ctrl.shuUC.GetForMember(c.Context(), memberID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, result)
}
