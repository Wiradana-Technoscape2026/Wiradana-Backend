package route

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wiradana/backend/internal/config"
	"github.com/wiradana/backend/internal/delivery/http/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, validate *validator.Validate, log *logrus.Logger) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "data": "ok"})
	})

	// Public
	auth := app.Group("/auth")
	_ = auth // POST /auth/login registered here later

	// Authenticated
	api := app.Group("/", middleware.Auth(cfg.JWT.Secret))

	// Pengurus endpoints
	pengurus := api.Group("/", middleware.RequireRole("pengurus"))
	_ = pengurus

	// Anggota portal endpoints
	portal := api.Group("/portal", middleware.RequireRole("anggota"))
	_ = portal

	// Integration (demo) endpoints
	integrations := api.Group("/integrations", middleware.RequireRole("pengurus"))
	_ = integrations

	// Inventory (Tier 3) — guarded by module
	inventory := api.Group("/inventory",
		middleware.RequireRole("pengurus"),
		middleware.RequireModule(db, "inventory"),
	)
	_ = inventory

	log.Info("routes registered")
}
