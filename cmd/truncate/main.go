package main

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/wiradana/backend/internal/config"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logrus.New()
	db, err := config.NewDatabase(cfg, logger)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Clearing database transactional data...")
	err = db.Transaction(func(tx *gorm.DB) error {
		// TRUNCATE operational tables CASCADE (keeps app_user, member, cooperative, configs)
		queries := []string{
			"TRUNCATE TABLE payment CASCADE",
			"TRUNCATE TABLE installment_schedule CASCADE",
			"TRUNCATE TABLE loan_audit_token CASCADE",
			"TRUNCATE TABLE loan_audit_log CASCADE",
			"TRUNCATE TABLE loan CASCADE",
			"TRUNCATE TABLE credit_assessment CASCADE",
			"TRUNCATE TABLE loan_application CASCADE",
			"TRUNCATE TABLE savings_transaction CASCADE",
			"TRUNCATE TABLE shu_distribution CASCADE",
			"TRUNCATE TABLE shu_period CASCADE",
			"TRUNCATE TABLE inventory_movement CASCADE",
			"TRUNCATE TABLE inventory_product CASCADE",
			"TRUNCATE TABLE notification_log CASCADE",
		}
		for _, q := range queries {
			if err := tx.Exec(q).Error; err != nil {
				return fmt.Errorf("failed executing %q: %w", q, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to clear database: %v", err)
	}
	log.Println("Database transactional data cleared successfully (users, members, and coops preserved)!")
}
