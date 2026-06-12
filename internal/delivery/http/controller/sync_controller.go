package controller

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/usecase"
)

type SyncController struct {
	syncUC   usecase.SyncUsecase
	validate *validator.Validate
}

func NewSyncController(syncUC usecase.SyncUsecase, validate *validator.Validate) *SyncController {
	return &SyncController{syncUC: syncUC, validate: validate}
}

func (ctrl *SyncController) Pull(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	role := c.Locals("role").(string)
	memberID, _ := c.Locals("member_id").(string)

	var since *time.Time
	if sinceStr := c.Query("since"); sinceStr != "" {
		t, err := time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "since harus format RFC3339 (contoh: 2026-06-01T00:00:00Z)")
		}
		since = &t
	}

	resp, err := ctrl.syncUC.Pull(c.Context(), coopID, since, role, &memberID)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal mengambil data sync")
	}
	return OK(c, resp)
}

func (ctrl *SyncController) Push(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	userID := c.Locals("user_id").(string)

	var req model.SyncPushRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "request body tidak valid")
	}
	if err := ctrl.validate.Struct(req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	resp, err := ctrl.syncUC.Push(c.Context(), coopID, userID, req.Mutations)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal memproses sync push")
	}
	return OK(c, resp)
}
