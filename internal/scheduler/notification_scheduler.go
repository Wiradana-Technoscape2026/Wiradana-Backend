package scheduler

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wiradana/backend/internal/usecase"
)

func StartNotificationScheduler(ctx context.Context, uc usecase.NotificationUsecase, log *logrus.Logger) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				wib := time.FixedZone("WIB", 7*3600)
				now := t.In(wib)
				if now.Hour() == 8 && now.Minute() == 0 {
					log.Info("scheduler: menjalankan pengiriman notifikasi harian")
					runDailyNotifications(ctx, uc, log)
				}
			}
		}
	}()
}

func runDailyNotifications(ctx context.Context, uc usecase.NotificationUsecase, log *logrus.Logger) {
	if err := uc.SendDueReminders(ctx, 3); err != nil {
		log.Errorf("scheduler: reminder_3d gagal: %v", err)
	}
	if err := uc.SendDueReminders(ctx, 1); err != nil {
		log.Errorf("scheduler: reminder_1d gagal: %v", err)
	}
	if err := uc.SendDueTodayAlerts(ctx); err != nil {
		log.Errorf("scheduler: due_today gagal: %v", err)
	}
	if err := uc.SendOverdueAlerts(ctx); err != nil {
		log.Errorf("scheduler: overdue gagal: %v", err)
	}
}
