package route

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wiradana/backend/internal/config"
	"github.com/wiradana/backend/internal/delivery/http/controller"
	"github.com/wiradana/backend/internal/delivery/http/middleware"
	"github.com/wiradana/backend/internal/gateway/adins"
	gwnot "github.com/wiradana/backend/internal/gateway/notification"
	"github.com/wiradana/backend/internal/repository"
	"github.com/wiradana/backend/internal/scheduler"
	"github.com/wiradana/backend/internal/usecase"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, validate *validator.Validate, log *logrus.Logger) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "data": "ok"})
	})

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	savingsRepo := repository.NewSavingsRepository(db)
	loanConfigRepo := repository.NewLoanConfigRepository(db)
	loanAppRepo := repository.NewLoanApplicationRepository(db)
	loanRepo := repository.NewLoanRepository(db)
	installmentRepo := repository.NewInstallmentRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	shuRepo := repository.NewShuRepository(db)
	moduleRepo := repository.NewModuleRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// ── Gateways ──────────────────────────────────────────────────────────────
	var ocrGateway adins.KTPOCRGateway
	if cfg.OCR.APIKey == "" {
		log.Warn("OCR_API_KEY tidak dikonfigurasi — menggunakan mock LITEDMS (MOCK_LITEDMS)")
		ocrGateway = &adins.MockKTPOCRGateway{}
	} else {
		ocrGateway = adins.NewAPICoIDGateway(cfg.OCR.APIKey, cfg.OCR.BaseURL)
	}
	scoringGateway := adins.NewMockScoringGateway()

	var notifGateway gwnot.Gateway
	if cfg.WhatsApp.Mode == "sandbox" {
		notifGateway = gwnot.NewWhatsAppGateway(cfg.WhatsApp.Token, cfg.WhatsApp.PhoneID)
	} else {
		notifGateway = gwnot.NewMockGateway(log)
	}

	// ── Usecases ──────────────────────────────────────────────────────────────
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	memberUC := usecase.NewMemberUsecase(memberRepo, userRepo)
	savingsUC := usecase.NewSavingsUsecase(savingsRepo, memberRepo)
	ocrUC := usecase.NewOCRUsecase(ocrGateway)
	loanConfigUC := usecase.NewLoanConfigUsecase(loanConfigRepo)
	loanAppUC := usecase.NewLoanApplicationUsecase(loanAppRepo, loanConfigRepo, memberRepo, loanRepo, scoringGateway)
	loanUC := usecase.NewLoanUsecase(loanRepo)
	installmentUC := usecase.NewInstallmentUsecase(installmentRepo, loanRepo)
	dashboardUC := usecase.NewDashboardUsecase(dashboardRepo)
	shuUC := usecase.NewShuUsecase(shuRepo)
	moduleUC := usecase.NewModuleUsecase(moduleRepo)
	inventoryUC := usecase.NewInventoryUsecase(inventoryRepo)
	notifUC := usecase.NewNotificationUsecase(notifRepo, memberRepo, notifGateway)
	reportUC := usecase.NewReportUsecase(reportRepo)

	// ── Controllers ───────────────────────────────────────────────────────────
	authCtrl := controller.NewAuthController(authUC, validate)
	memberCtrl := controller.NewMemberController(memberUC, validate)
	savingsCtrl := controller.NewSavingsController(savingsUC, notifUC, validate)
	ocrCtrl := controller.NewOCRController(ocrUC, log)
	loanConfigCtrl := controller.NewLoanConfigController(loanConfigUC, validate)
	loanAppCtrl := controller.NewLoanApplicationController(loanAppUC, notifUC, validate)
	loanCtrl := controller.NewLoanController(loanUC)
	installmentCtrl := controller.NewInstallmentController(installmentUC, validate)
	scoringCtrl := controller.NewScoringController(scoringGateway)
	portalLoanCtrl := controller.NewPortalLoanController(loanAppUC, loanUC, validate)
	dashboardCtrl := controller.NewDashboardController(dashboardUC)
	shuCtrl := controller.NewShuController(shuUC, validate)
	moduleCtrl := controller.NewModuleController(moduleUC, validate)
	inventoryCtrl := controller.NewInventoryController(inventoryUC, validate)
	portalCtrl := controller.NewPortalController(memberUC, savingsUC, shuUC)
	notifCtrl := controller.NewNotificationController(notifUC)
	reportCtrl := controller.NewReportController(reportUC)

	// ── Middleware shorthands ─────────────────────────────────────────────────
	// In Fiber v2, handlers passed to Group() are registered via app.Use(prefix),
	// meaning a Group("", m1, m2) on /api/v1 applies m1+m2 to ALL /api/v1/* routes
	// including /portal/*. Fix: apply middleware per-route, not per-group.
	authM := middleware.Auth(cfg.JWT.Secret)
	asP := middleware.RequireRole("pengurus") // assert pengurus
	asA := middleware.RequireRole("anggota")  // assert anggota

	p := func(h fiber.Handler) []fiber.Handler { return []fiber.Handler{authM, asP, h} }
	a := func(h fiber.Handler) []fiber.Handler { return []fiber.Handler{authM, asA, h} }

	api := app.Group("/api/v1")

	// ── Public ────────────────────────────────────────────────────────────────
	api.Post("/auth/login", authCtrl.Login)
	api.Post("/auth/select-cooperative", authCtrl.SelectCooperative)
	api.Post("/auth/register/pengurus", authCtrl.RegisterPengurus)

	// ── Pengurus ──────────────────────────────────────────────────────────────
	// Groups used only for URL prefix — no middleware in Group() to avoid global Use leak.
	pengurus := api.Group("")
	pengurus.Post("/members", p(memberCtrl.Create)...)
	pengurus.Get("/members", p(memberCtrl.List)...)
	pengurus.Get("/members/:id", p(memberCtrl.GetByID)...)
	pengurus.Put("/members/:id", p(memberCtrl.Update)...)
	pengurus.Get("/members/:id/savings", p(savingsCtrl.List)...)
	pengurus.Post("/members/:id/savings", p(savingsCtrl.Record)...)
	pengurus.Get("/savings/summary", p(savingsCtrl.Summary)...)
	pengurus.Get("/savings/transactions", p(savingsCtrl.ListAll)...)

	pengurus.Get("/loan-config", p(loanConfigCtrl.Get)...)
	pengurus.Put("/loan-config", p(loanConfigCtrl.Update)...)

	pengurus.Get("/loan-applications", p(loanAppCtrl.List)...)
	pengurus.Post("/loan-applications", p(loanAppCtrl.Create)...)
	pengurus.Post("/loan-applications/:id/approve", p(loanAppCtrl.Approve)...)
	pengurus.Post("/loan-applications/:id/reject", p(loanAppCtrl.Reject)...)

	pengurus.Get("/loans", p(loanCtrl.List)...)
	pengurus.Get("/loans/:id", p(loanCtrl.GetByID)...)
	pengurus.Get("/loans/:id/export", p(loanCtrl.ExportExcel)...)

	pengurus.Post("/installments/:id/pay", p(installmentCtrl.Pay)...)

	pengurus.Get("/dashboard", p(dashboardCtrl.Get)...)

	pengurus.Get("/shu-periods", p(shuCtrl.List)...)
	pengurus.Post("/shu-periods", p(shuCtrl.Create)...)
	pengurus.Get("/shu-periods/:id", p(shuCtrl.GetByID)...)
	pengurus.Post("/shu-periods/:id/calculate", p(shuCtrl.Calculate)...)

	pengurus.Get("/modules", p(moduleCtrl.List)...)
	pengurus.Put("/modules/:key", p(moduleCtrl.Update)...)

	pengurus.Get("/reports/summary", p(reportCtrl.Summary)...)

	pengurus.Get("/notifications/logs", p(notifCtrl.ListLogs)...)
	pengurus.Post("/notifications/trigger", p(notifCtrl.Trigger)...)

	// ── Portal Anggota ────────────────────────────────────────────────────────
	portal := api.Group("/portal")
	portal.Get("/me", a(portalCtrl.Me)...)
	portal.Get("/shu", a(portalCtrl.SHU)...)
	portal.Get("/loan-applications", a(portalLoanCtrl.ListApplications)...)
	portal.Post("/loan-applications", a(portalLoanCtrl.Apply)...)
	portal.Get("/loans", a(portalLoanCtrl.ListLoans)...)
	portal.Get("/loans/:id", a(portalLoanCtrl.GetLoan)...)

	// ── Integration (demo endpoints) ──────────────────────────────────────────
	integrations := api.Group("/integrations")
	integrations.Post("/adins/ocr/ktp", p(ocrCtrl.ExtractKTP)...)
	integrations.Post("/adins/credit-scoring", p(scoringCtrl.Score)...)

	// ── Sync (offline-first) ─────────────────────────────────────────────────
	syncRepo := repository.NewSyncRepository(db)
	syncUC := usecase.NewSyncUsecase(syncRepo, memberUC, savingsUC, loanAppUC, installmentUC, loanConfigUC, inventoryUC)
	syncCtrl := controller.NewSyncController(syncUC, validate)
	syncGroup := api.Group("/sync")
	syncGroup.Get("/pull", authM, syncCtrl.Pull)
	syncGroup.Post("/push", p(syncCtrl.Push)...)

	// ── Inventory (Tier 3) ────────────────────────────────────────────────────
	invM := middleware.RequireModule(db, "inventory")
	inv := func(h fiber.Handler) []fiber.Handler { return []fiber.Handler{authM, asP, invM, h} }

	inventory := api.Group("/inventory")
	inventory.Get("/field-defs", inv(inventoryCtrl.ListFieldDefs)...)
	inventory.Post("/field-defs", inv(inventoryCtrl.CreateFieldDef)...)
	inventory.Delete("/field-defs/:id", inv(inventoryCtrl.DeleteFieldDef)...)
	inventory.Get("/products", inv(inventoryCtrl.ListProducts)...)
	inventory.Post("/products", inv(inventoryCtrl.CreateProduct)...)
	inventory.Put("/products/:id", inv(inventoryCtrl.UpdateProduct)...)
	inventory.Post("/products/:id/movements", inv(inventoryCtrl.RecordMovement)...)

	scheduler.StartNotificationScheduler(context.Background(), notifUC, log)

	log.Info("routes registered")
}
