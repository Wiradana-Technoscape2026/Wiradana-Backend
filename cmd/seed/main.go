package main

import (
	"log"
	"os"

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

	log.Println("Seeding final demo database...")

	// 1. Wipe all tables first for a clean seed
	log.Println("Wiping all tables for clean seed...")
	err = db.Transaction(func(tx *gorm.DB) error {
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
			"TRUNCATE TABLE coop_module CASCADE",
			"TRUNCATE TABLE loan_config CASCADE",
			"TRUNCATE TABLE app_user CASCADE",
			"TRUNCATE TABLE member CASCADE",
			"TRUNCATE TABLE cooperative CASCADE",
		}
		for _, q := range queries {
			if err := tx.Exec(q).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to wipe database: %v", err)
	}

	// 2. Load and run seed_dummy.sql
	log.Println("Loading and running db/seed_dummy.sql...")
	sqlBytes, err := os.ReadFile("db/seed_dummy.sql")
	if err != nil {
		log.Fatalf("failed to read db/seed_dummy.sql: %v", err)
	}

	if err := db.Exec(string(sqlBytes)).Error; err != nil {
		log.Fatalf("failed to run seed_dummy.sql: %v", err)
	}

	// 3. Run transparency backfill query
	log.Println("Backfilling transparency fields...")
	err = db.Transaction(func(tx *gorm.DB) error {
		// Ahmad Fauzi for Padiwangi savings
		q1 := "UPDATE savings_transaction SET recorded_by = '99999999-9999-9999-9999-999999999991' WHERE cooperative_id = '11111111-1111-1111-1111-111111111111'"
		if err := tx.Exec(q1).Error; err != nil {
			return err
		}

		// Bambang Wijaya for Sawargi savings
		q2 := "UPDATE savings_transaction SET recorded_by = '99999999-9999-9999-9999-999999999994' WHERE cooperative_id = '22222222-2222-2222-2222-222222222222'"
		if err := tx.Exec(q2).Error; err != nil {
			return err
		}

		// Approved/rejected loan applications approved_at
		q3 := "UPDATE loan_application SET approved_at = created_at + INTERVAL '2 days' WHERE status IN ('approved', 'rejected') AND approved_at IS NULL"
		if err := tx.Exec(q3).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Failed to backfill transparency fields: %v", err)
	}

	log.Println("Demo seed completed successfully! Database is ready for video demo.")
}
