package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type PortalController struct {
	memberUC    usecase.MemberUsecase
	savingsUC   usecase.SavingsUsecase
	shuUC       usecase.ShuUsecase
	authUC      usecase.AuthUsecase
	dashboardUC usecase.DashboardUsecase
}

func NewPortalController(
	memberUC usecase.MemberUsecase,
	savingsUC usecase.SavingsUsecase,
	shuUC usecase.ShuUsecase,
	authUC usecase.AuthUsecase,
	dashboardUC usecase.DashboardUsecase,
) *PortalController {
	return &PortalController{
		memberUC:    memberUC,
		savingsUC:   savingsUC,
		shuUC:       shuUC,
		authUC:      authUC,
		dashboardUC: dashboardUC,
	}
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
		EstimatedShu int64                         `json:"estimated_shu"`
		History      []model.ShuDistributionDetail `json:"history"`
	}

	if dists == nil {
		dists = []model.ShuDistributionDetail{}
	}

	return OK(c, shuResponse{
		EstimatedShu: estimatedShu,
		History:      dists,
	})
}

// Cooperatives mengembalikan semua koperasi yang dimiliki user yang sedang login.
func (ctrl *PortalController) Cooperatives(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	memberships, err := ctrl.authUC.GetUserCooperatives(c.Context(), userID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal mengambil daftar koperasi")
	}

	type coopItem struct {
		CooperativeID   string `json:"cooperative_id"`
		CooperativeName string `json:"cooperative_name"`
		MemberID        string `json:"member_id"`
	}
	result := make([]coopItem, len(memberships))
	for i, m := range memberships {
		result[i] = coopItem{
			CooperativeID:   m.CooperativeID,
			CooperativeName: m.CooperativeName,
			MemberID:        m.MemberID,
		}
	}
	return OKList(c, result, int64(len(result)))
}

// MemberDashboard mengembalikan data dashboard anggota untuk koperasi tertentu.
// Query param: cooperative_id (wajib). Membership diverifikasi dari DB, bukan dari JWT scope.
func (ctrl *PortalController) MemberDashboard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	coopID := c.Query("cooperative_id")
	if coopID == "" {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "cooperative_id wajib diisi")
	}

	memberships, err := ctrl.authUC.GetUserCooperatives(c.Context(), userID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal memverifikasi keanggotaan")
	}

	var memberID string
	for _, m := range memberships {
		if m.CooperativeID == coopID {
			memberID = m.MemberID
			break
		}
	}
	if memberID == "" {
		return Fail(c, fiber.StatusForbidden, "FORBIDDEN", "Anda bukan anggota koperasi ini")
	}

	resp, err := ctrl.dashboardUC.GetMemberDashboard(c.Context(), coopID, memberID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal mengambil data dashboard")
	}

	return OK(c, resp)
}
