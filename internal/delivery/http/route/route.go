package route

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wiradana/backend/internal/config"
	"github.com/wiradana/backend/internal/delivery/http/controller"
	"github.com/wiradana/backend/internal/delivery/http/middleware"
	"github.com/wiradana/backend/internal/gateway/adins"
	"github.com/wiradana/backend/internal/repository"
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
<<<<<<< HEAD
	savingsRepo := repository.NewSavingsRepository(db)
	shuRepo := repository.NewShuRepository(db)
	moduleRepo := repository.NewModuleRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
=======
	dashboardRepo := repository.NewDashboardRepository(db)
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed

	// ── Gateways ──────────────────────────────────────────────────────────────
	ocrGateway := adins.NewAPICoIDGateway(cfg.OCR.APIKey, cfg.OCR.BaseURL)
	scoringGateway := adins.NewMockScoringGateway()

	// ── Usecases ──────────────────────────────────────────────────────────────
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	memberUC := usecase.NewMemberUsecase(memberRepo)
	savingsUC := usecase.NewSavingsUsecase(savingsRepo, memberRepo)
	ocrUC := usecase.NewOCRUsecase(ocrGateway)
	loanConfigUC := usecase.NewLoanConfigUsecase(loanConfigRepo)
	loanAppUC := usecase.NewLoanApplicationUsecase(loanAppRepo, loanConfigRepo, memberRepo, loanRepo, scoringGateway)
	loanUC := usecase.NewLoanUsecase(loanRepo)
	installmentUC := usecase.NewInstallmentUsecase(installmentRepo, loanRepo)
<<<<<<< HEAD
	savingsUC := usecase.NewSavingsUsecase(savingsRepo)
	shuUC := usecase.NewShuUsecase(shuRepo)
	moduleUC := usecase.NewModuleUsecase(moduleRepo)
	dashboardUC := usecase.NewDashboardUsecase(dashboardRepo)
	inventoryUC := usecase.NewInventoryUsecase(inventoryRepo)
=======
	dashboardUC := usecase.NewDashboardUsecase(dashboardRepo)
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed

	// ── Controllers ───────────────────────────────────────────────────────────
	authCtrl := controller.NewAuthController(authUC, validate)
	memberCtrl := controller.NewMemberController(memberUC, validate)
	savingsCtrl := controller.NewSavingsController(savingsUC, validate)
	ocrCtrl := controller.NewOCRController(ocrUC, log)
	loanConfigCtrl := controller.NewLoanConfigController(loanConfigUC, validate)
	loanAppCtrl := controller.NewLoanApplicationController(loanAppUC, validate)
	loanCtrl := controller.NewLoanController(loanUC)
	installmentCtrl := controller.NewInstallmentController(installmentUC, validate)
	scoringCtrl := controller.NewScoringController(scoringGateway)
	portalLoanCtrl := controller.NewPortalLoanController(loanAppUC, loanUC, validate)
<<<<<<< HEAD
	savingsCtrl := controller.NewSavingsController(savingsUC, validate)
	shuCtrl := controller.NewShuController(shuUC, validate)
	moduleCtrl := controller.NewModuleController(moduleUC, validate)
	dashboardCtrl := controller.NewDashboardController(dashboardUC)
	portalCtrl := controller.NewPortalController(memberUC, savingsUC, shuUC)
	inventoryCtrl := controller.NewInventoryController(inventoryUC, validate)
=======
	dashboardCtrl := controller.NewDashboardController(dashboardUC)
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed

	api := app.Group("/api/v1")

	// ── Public ────────────────────────────────────────────────────────────────
	api.Post("/auth/login", authCtrl.Login)
	api.Post("/auth/register/pengurus", authCtrl.RegisterPengurus)

	// ── Pengurus ──────────────────────────────────────────────────────────────
	pengurus := api.Group("", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole("pengurus"))
	pengurus.Post("/members", memberCtrl.Create)
	pengurus.Get("/members", memberCtrl.List)
	pengurus.Get("/members/:id", memberCtrl.GetByID)
	pengurus.Put("/members/:id", memberCtrl.Update)
	pengurus.Get("/members/:id/savings", savingsCtrl.List)
	pengurus.Post("/members/:id/savings", savingsCtrl.Record)

	pengurus.Get("/members/:id/savings", savingsCtrl.List)
	pengurus.Post("/members/:id/savings", savingsCtrl.Record)

	pengurus.Get("/loan-config", loanConfigCtrl.Get)
	pengurus.Put("/loan-config", loanConfigCtrl.Update)

	pengurus.Get("/loan-applications", loanAppCtrl.List)
	pengurus.Post("/loan-applications", loanAppCtrl.Create)
	pengurus.Post("/loan-applications/:id/approve", loanAppCtrl.Approve)
	pengurus.Post("/loan-applications/:id/reject", loanAppCtrl.Reject)

	pengurus.Get("/loans", loanCtrl.List)
	pengurus.Get("/loans/:id", loanCtrl.GetByID)

	pengurus.Post("/installments/:id/pay", installmentCtrl.Pay)

<<<<<<< HEAD
	pengurus.Get("/shu-periods", shuCtrl.ListPeriods)
	pengurus.Post("/shu-periods", shuCtrl.CreatePeriod)
	pengurus.Post("/shu-periods/:id/calculate", shuCtrl.Calculate)

	pengurus.Get("/modules", moduleCtrl.List)
	pengurus.Put("/modules/:key", moduleCtrl.Update)

=======
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
	pengurus.Get("/dashboard", dashboardCtrl.Get)

	// ── Portal Anggota ────────────────────────────────────────────────────────
	portal := api.Group("/portal", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole("anggota"))
	portal.Get("/me", portalCtrl.Me)
	portal.Get("/shu", portalCtrl.SHU)
	portal.Get("/loan-applications", portalLoanCtrl.ListApplications)
	portal.Post("/loan-applications", portalLoanCtrl.Apply)
	portal.Get("/loans", portalLoanCtrl.ListLoans)
	portal.Get("/loans/:id", portalLoanCtrl.GetLoan)

	// ── Integration (demo endpoints) ──────────────────────────────────────────
	integrations := api.Group("/integrations", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole("pengurus"))
	integrations.Post("/adins/ocr/ktp", ocrCtrl.ExtractKTP)
	integrations.Post("/adins/credit-scoring", scoringCtrl.Score)

<<<<<<< HEAD
	// ── Inventory (Tier 3 — guarded by RequireModule) ─────────────────────────
	inv := pengurus.Group("/inventory", middleware.RequireModule(db, "inventory"))
	inv.Get("/field-defs", inventoryCtrl.ListFieldDefs)
	inv.Post("/field-defs", inventoryCtrl.CreateFieldDef)
	inv.Delete("/field-defs/:id", inventoryCtrl.DeleteFieldDef)
	inv.Get("/products", inventoryCtrl.ListProducts)
	inv.Post("/products", inventoryCtrl.CreateProduct)
	inv.Put("/products/:id", inventoryCtrl.UpdateProduct)
	inv.Post("/products/:id/movements", inventoryCtrl.RecordMovement)
=======
	// ── Sync (offline-first) ─────────────────────────────────────────────────
	syncRepo := repository.NewSyncRepository(db)
	syncUC := usecase.NewSyncUsecase(syncRepo, memberUC, savingsUC, loanAppUC, installmentUC, loanConfigUC)
	syncCtrl := controller.NewSyncController(syncUC, validate)
	syncGroup := api.Group("/sync", middleware.Auth(cfg.JWT.Secret))
	syncGroup.Get("/pull", syncCtrl.Pull)
	syncGroup.Post("/push", middleware.RequireRole("pengurus"), syncCtrl.Push)

	// ── Inventory (Tier 3) ────────────────────────────────────────────────────
	inventory := api.Group("/inventory",
		middleware.Auth(cfg.JWT.Secret),
		middleware.RequireRole("pengurus"),
		middleware.RequireModule(db, "inventory"),
	)
	_ = inventory
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed

	log.Info("routes registered")
}
