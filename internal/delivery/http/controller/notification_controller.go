package controller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/usecase"
)

type NotificationController struct {
	uc usecase.NotificationUsecase
}

func NewNotificationController(uc usecase.NotificationUsecase) *NotificationController {
	return &NotificationController{uc: uc}
}

func (ctrl *NotificationController) ListLogs(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	eventType := c.Query("event_type")
	limit := c.QueryInt("limit", 20)
	page := c.QueryInt("page", 1)

	logs, total, err := ctrl.uc.ListLogs(c.Context(), coopID, eventType, page, limit)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "terjadi kesalahan server")
	}
	return OKList(c, logs, total)
}

func (ctrl *NotificationController) Trigger(c *fiber.Ctx) error {
	go func() {
		ctx := context.Background()
		ctrl.uc.SendDueReminders(ctx, 3)
		ctrl.uc.SendDueReminders(ctx, 1)
		ctrl.uc.SendDueTodayAlerts(ctx)
		ctrl.uc.SendOverdueAlerts(ctx)
	}()
	return OK(c, fiber.Map{"message": "notifikasi dijadwalkan untuk dikirim"})
}
