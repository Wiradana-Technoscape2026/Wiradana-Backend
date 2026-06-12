package main

import (
	"log"

	"github.com/wiradana/backend/internal/config"
	"github.com/wiradana/backend/internal/delivery/http/route"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := config.NewLogger(cfg)
	logger.Info("config loaded")

	db, err := config.NewDatabase(cfg, logger)
	if err != nil {
		logger.Fatalf("failed to connect database: %v", err)
	}

	validate := config.NewValidator()
	app := config.NewFiber(cfg)

	route.RegisterRoutes(app, db, cfg, validate, logger)

	logger.Infof("starting server on :%s", cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}
