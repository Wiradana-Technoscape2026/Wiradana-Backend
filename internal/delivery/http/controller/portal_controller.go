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

func NewPortalController(memberUC usecase.MemberUsecase, savingsUC usecase.SavingsUsecase, shuUC usecase.ShuUsecase) *PortalController {
	return &PortalController{memberUC: memberUC, savingsUC: savingsUC, shuUC: shuUC}
}

func (ctrl *PortalController) Me(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	memberID := c.Locals("member_id").(string)

	member, err := ctrl.memberUC.FindByID(c.Context(), coopID, memberID)
	if err != nil {
		return Fail(c, fiber.StatusNotFound, "NOT_FOUND", "anggota tidak ditemukan")
	}

	recentTxs, _ := ctrl.savingsUC.FindByMember(c.Context(), coopID, memberID)
	if recentTxs == nil {
		recentTxs = []model.SavingsTransactionResponse{}
	}

	if len(recentTxs) > 5 {
		recentTxs = recentTxs[:5]
	}

	type meResponse struct {
		Member         *model.MemberResponse              `json:"member"`
		SavingsSummary model.SavingsSummary               `json:"savings_summary"`
		SavingsRecent  []model.SavingsTransactionResponse `json:"savings_recent"`
	}

	return OK(c, meResponse{
		Member:         member,
		SavingsSummary: member.SavingsSummary,
		SavingsRecent:  recentTxs,
	})
}

func (ctrl *PortalController) SHU(c *fiber.Ctx) error {
	memberID := c.Locals("member_id").(string)

	dists, err := ctrl.shuUC.GetMemberDistributions(c.Context(), memberID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal mengambil data SHU")
	}

	var estimatedShu int64
	for _, d := range dists {
		estimatedShu += d.TotalShu
	}

	type shuResponse struct {
		EstimatedShu int64                          `json:"estimated_shu"`
		History      []model.ShuDistributionDetail  `json:"history"`
	}

	if dists == nil {
		dists = []model.ShuDistributionDetail{}
	}

	return OK(c, shuResponse{
		EstimatedShu: estimatedShu,
		History:      dists,
	})
}
