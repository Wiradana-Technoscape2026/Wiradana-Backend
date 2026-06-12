package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OCR      OCRConfig
	WhatsApp WhatsAppConfig
}

type OCRConfig struct {
	APIKey  string
	BaseURL string
}

type WhatsAppConfig struct {
	Mode    string // "mock" | "sandbox"
	Token   string
	PhoneID string
}

type AppConfig struct {
	Port string
	Name string
}

type DatabaseConfig struct {
	DSN          string
	MigrationDSN string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("../..")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("APP_NAME", "Wiradana Backend")
	v.SetDefault("JWT_EXPIRATION_HOURS", 24)
	v.SetDefault("OCR_BASE_URL", "https://use.api.co.id")
	v.SetDefault("WA_MODE", "mock")

	cfg := &Config{
		App: AppConfig{
			Port: v.GetString("APP_PORT"),
			Name: v.GetString("APP_NAME"),
		},
		Database: DatabaseConfig{
			DSN:          v.GetString("DATABASE_DSN"),
			MigrationDSN: v.GetString("DATABASE_MIGRATION_DSN"),
		},
		JWT: JWTConfig{
			Secret:          v.GetString("JWT_SECRET"),
			ExpirationHours: v.GetInt("JWT_EXPIRATION_HOURS"),
		},
		OCR: OCRConfig{
			APIKey:  v.GetString("OCR_API_KEY"),
			BaseURL: v.GetString("OCR_BASE_URL"),
		},
		WhatsApp: WhatsAppConfig{
			Mode:    v.GetString("WA_MODE"),
			Token:   v.GetString("WA_TOKEN"),
			PhoneID: v.GetString("WA_PHONE_ID"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) validate() error {
	missing := []string{}

	required := map[string]string{
		"DATABASE_DSN":           cfg.Database.DSN,
		"DATABASE_MIGRATION_DSN": cfg.Database.MigrationDSN,
		"JWT_SECRET":             cfg.JWT.Secret,
		"OCR_API_KEY":            cfg.OCR.APIKey,
	}

	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}

	return nil
}
