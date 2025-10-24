package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	pgConn string
	BotConfig
	PaymentConfig PaymentConfig
	JwtConfig     JwtConfig
	SMTPConfig    SMTPConfig
}

type BotConfig struct {
	token      string
	webhookURL string
	useWebhook bool
	isActive   bool
}

type JwtConfig struct {
	JwtSecret     string
	RefreshSecret string
}

type SMTPConfig struct {
	Address string
	PassKey string
}

func (c *Config) PgConn() string {
	return c.pgConn
}

func (c *Config) Token() string {
	return c.token
}

func (c *Config) WebhookURL() string {
	return c.webhookURL
}

func (c *Config) UseWebhook() bool {
	return c.useWebhook
}

func (c *Config) BotIsActive() bool {
	return c.isActive
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		BotConfig: BotConfig{
			webhookURL: os.Getenv("WEBHOOK_URL"),
			token:      os.Getenv("BOT_TOKEN"),
			useWebhook: os.Getenv("USE_WEBHOOK") == "true",
			isActive:   os.Getenv("BOT_IS_ACTIVE") == "true",
		},
		pgConn: fmt.Sprintf(
			"user=%s password=%s host=%s dbname=%s sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_NAME"),
		),
		JwtConfig: JwtConfig{
			JwtSecret:     os.Getenv("JWT_SECRET_KEY"),
			RefreshSecret: os.Getenv("REFRESH_JWT_SECRET_KEY"),
		},
		SMTPConfig: SMTPConfig{
			Address: os.Getenv("ADMIN_EMAIL"),
			PassKey: os.Getenv("PASS_KEY"),
		},
	}
}
