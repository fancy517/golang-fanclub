package app

import (
	"fanclub/internal/env"
	"fmt"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	chain string

	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string

	allowOrigin string

	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}

	house struct {
		address    string
		privateKey string
		minDeposit float64
		earningFee float64
	}

	// SCAN
	blockScanInterval int

	// Admin contact email
	siteContactEmail string
}

const (
	DefaultHost              = "localhost"
	DefaultPort              = 3306
	DefaultUser              = "root"
	DefaultPassword          = "root"
	DefaultDBName            = "fanclub"
	DefaultBlockScanInterval = 1

	DefaultSMTPPort = 587
)

func (app *application) loadConfig() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("failed to get load env; %w", err)
	}

	app.config.chain = env.GetString("NETWORK", "Mainnet")

	app.config.allowOrigin = env.GetString("ALLOW_ORIGIN", "http://localhost:3000")

	// database
	app.config.dbHost = env.GetString("DB_HOST", DefaultHost)
	app.config.dbPort = env.GetInt("DB_PORT", DefaultPort)
	app.config.dbUser = env.GetString("DB_USER", DefaultUser)
	app.config.dbPassword = env.GetString("DB_PASSWORD", DefaultPassword)
	app.config.dbName = env.GetString("DB_NAME", DefaultDBName)

	// smtp
	app.config.smtp.host = env.GetString("SMTP_HOST", "")
	app.config.smtp.port = env.GetInt("SMTP_PORT", DefaultSMTPPort)
	app.config.smtp.username = env.GetString("SMTP_USER", "user")
	app.config.smtp.password = env.GetString("SMTP_PASSWORD", "")
	app.config.smtp.sender = env.GetString("SMTP_SENDER", "noreply@example.com")

	// house wallet
	app.config.house.address = env.GetString("HOUSE_WALLET_ADDRESS", "")
	app.config.house.privateKey = env.GetString("HOUSE_WALLET_PRIVATE", "")
	app.config.house.minDeposit = float64(env.GetInt("HOUSE_MIN_DEPOSIT", 10))
	app.config.house.earningFee = float64(env.GetFloat("HOUSE_EARNING_FEE", 5.0))

	// scanner
	app.config.blockScanInterval = env.GetInt("BLOCK_SCAN_INTERVAL", DefaultBlockScanInterval)

	// site contact
	app.config.siteContactEmail = env.GetString("SITE_CONTACT_EMAIL_ADDRESS", "admin@fanclub.io")

	return nil
}

func (c AppConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", c.dbUser, c.dbPassword, c.dbHost, c.dbName)
}
