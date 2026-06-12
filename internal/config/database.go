package config

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDatabase(cfg *Config, log *logrus.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.Database.DSN,
		PreferSimpleProtocol: true, // required for Supabase transaction pooler (pgbouncer)
	}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // match SQL table names: cooperative, member, loan, etc.
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Info("database connection established")
	return db, nil
}
