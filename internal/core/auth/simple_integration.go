package auth

import (
	"log/slog"
	"os"

	"github.com/aby-med/medical-platform/internal/infrastructure/email"
	"github.com/aby-med/medical-platform/internal/infrastructure/sms"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// SimpleIntegration provides an easy way to add auth to existing system
func SimpleIntegration(router chi.Router, dbPool *pgxpool.Pool, logger *slog.Logger) error {
	// Convert pgxpool to sqlx for auth module
	// Auth module uses sqlx, but main app uses pgxpool
	// Solution: Create a new sqlx connection for auth
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Construct from individual env vars
		host := getEnvOrDefault("DB_HOST", "localhost")
		port := getEnvOrDefault("DB_PORT", "5430")
		user := getEnvOrDefault("DB_USER", "postgres")
		password := getEnvOrDefault("DB_PASSWORD", "postgres")
		dbname := getEnvOrDefault("DB_NAME", "med_platform")
		dsn = "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"
	}

	logger.Info("Connecting auth module to database", slog.String("dsn", maskPassword(dsn)))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}

	// Initialize external services
	var emailSender EmailSender
	var smsSender SMSSender

	sendgridKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridKey == "" {
		logger.Warn("SendGrid not configured, using mock email service")
		emailSender = &MockEmailSender{}
	} else {
		emailSender = email.NewSendGridSender(
			sendgridKey,
			getEnvOrDefault("SENDGRID_FROM_EMAIL", "noreply@ServQR.com"),
			getEnvOrDefault("SENDGRID_FROM_NAME", "ServQR Platform"),
		)
	}

	twilioSID := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if twilioSID == "" || twilioToken == "" {
		logger.Warn("Twilio not configured, using mock SMS service")
		smsSender = &MockSMSSender{}
	} else {
		smsSender = sms.NewTwilioSender(
			twilioSID,
			twilioToken,
			getEnvOrDefault("TWILIO_PHONE_NUMBER", ""),
			getEnvOrDefault("TWILIO_WHATSAPP_NUMBER", ""),
		)
	}

	// Mark as used (SimpleIntegration is WIP - these will be used later)
	_, _ = emailSender, smsSender

	// Create auth module
	return IntegrateAuthModule(router, db, logger)
}

func maskPassword(dsn string) string {
	// Mask password in DSN for logging
	if len(dsn) > 20 {
		return dsn[:10] + "***" + dsn[len(dsn)-10:]
	}
	return "***"
}
