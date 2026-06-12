package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/gateway/adins"
	"github.com/wiradana/backend/internal/model"
)

type ScoringController struct {
	gw adins.ScoringGateway
}

func NewScoringController(gw adins.ScoringGateway) *ScoringController {
	return &ScoringController{gw: gw}
}

func (ctrl *ScoringController) Score(c *fiber.Ctx) error {
	var req model.CreditScoringRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	result, err := ctrl.gw.Score(c.Context(), adins.ScoringInput{
		MemberID:       req.MemberID,
		Features:       req.Features,
		JumlahDiajukan: req.JumlahDiajukan,
		TenorBulan:     req.TenorBulan,
		TotalSimpanan:  req.TotalSimpanan,
		MaxPlafond:     req.MaxPlafond,
	})
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OK(c, model.CreditScoringResponse{
		Score:            result.Score,
		Grade:            result.Grade,
		Recommendation:   result.Recommendation,
		LimitRekomendasi: result.LimitRekomendasi,
		Reasons:          result.Reasons,
		Source:           result.Source,
	})
}
