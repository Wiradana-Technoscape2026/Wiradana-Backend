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

	// Auth
	userRepo := repository.NewUserRepository(db)
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	authCtrl := controller.NewAuthController(authUC, validate)

	// Member
	memberRepo := repository.NewMemberRepository(db)
	memberUC := usecase.NewMemberUsecase(memberRepo)
	memberCtrl := controller.NewMemberController(memberUC, validate)

	// OCR
	ocrGateway := adins.NewAPICoIDGateway(cfg.OCR.APIKey, cfg.OCR.BaseURL)
	ocrUC := usecase.NewOCRUsecase(ocrGateway)
	ocrCtrl := controller.NewOCRController(ocrUC, log)

	api := app.Group("/api/v1")

	// Public
	api.Post("/auth/login", authCtrl.Login)
	api.Post("/auth/register/pengurus", authCtrl.RegisterPengurus)

	// Pengurus endpoints
	pengurus := api.Group("/", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole("pengurus"))
	pengurus.Post("/members", memberCtrl.Create)
	pengurus.Get("/members", memberCtrl.List)
	pengurus.Get("/members/:id", memberCtrl.GetByID)
	pengurus.Put("/members/:id", memberCtrl.Update)

	// Anggota portal endpoints
	portal := api.Group("/portal", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole("anggota"))
	_ = portal

	// Integration (demo) endpoints — pengurus only
	integrations := api.Group("/integrations", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole("pengurus"))
	integrations.Post("/adins/ocr/ktp", ocrCtrl.ExtractKTP)

	// Inventory (Tier 3) — guarded by module
	inventory := api.Group("/inventory",
		middleware.Auth(cfg.JWT.Secret),
		middleware.RequireRole("pengurus"),
		middleware.RequireModule(db, "inventory"),
	)
	_ = inventory

	log.Info("routes registered")
}
