package auth

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/aby-med/medical-platform/internal/infrastructure/email"
	"github.com/aby-med/medical-platform/internal/infrastructure/sms"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// IntegrateAuthModuleWithReturn adds authentication to the application and returns the module
func IntegrateAuthModuleWithReturn(router chi.Router, db *sqlx.DB, logger *slog.Logger) (*Module, error) {
	authModule, err := integrateAuthModuleInternal(router, db, logger)
	return authModule, err
}

// IntegrateAuthModule adds authentication to the application
func IntegrateAuthModule(router chi.Router, db *sqlx.DB, logger *slog.Logger) error {
	_, err := integrateAuthModuleInternal(router, db, logger)
	return err
}

// integrateAuthModuleInternal is the internal implementation
func integrateAuthModuleInternal(router chi.Router, db *sqlx.DB, logger *slog.Logger) (*Module, error) {
	logger.Info("Initializing authentication module")

	// Get configuration from environment
	jwtPrivateKeyPath := getEnvOrDefault("JWT_PRIVATE_KEY_PATH", "./keys/jwt-private.pem")
	jwtPublicKeyPath := getEnvOrDefault("JWT_PUBLIC_KEY_PATH", "./keys/jwt-public.pem")
	
	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")
	twilioWhatsAppNumber := os.Getenv("TWILIO_WHATSAPP_NUMBER")
	
	sendgridAPIKey := os.Getenv("SENDGRID_API_KEY")
	sendgridFromEmail := getEnvOrDefault("SENDGRID_FROM_EMAIL", "noreply@ServQR.com")
	sendgridFromName := getEnvOrDefault("SENDGRID_FROM_NAME", "ServQR Platform")

	// Initialize external services
	var emailSender EmailSender
	var smsSender SMSSender

	// Use mock services if credentials not provided (development)
	if sendgridAPIKey == "" {
		logger.Warn("SendGrid API key not configured, using mock email sender")
		emailSender = &MockEmailSender{}
	} else {
		emailSender = email.NewSendGridSender(sendgridAPIKey, sendgridFromEmail, sendgridFromName)
	}

	if twilioAccountSID == "" || twilioAuthToken == "" {
		logger.Warn("Twilio credentials not configured, using mock SMS sender")
		smsSender = &MockSMSSender{}
	} else {
		smsSender = sms.NewTwilioSender(twilioAccountSID, twilioAuthToken, twilioPhoneNumber, twilioWhatsAppNumber)
	}

	// Create auth module with configuration
	authModule, err := NewModule(db, &Config{
		JWTPrivateKeyPath:   jwtPrivateKeyPath,
		JWTPublicKeyPath:    jwtPublicKeyPath,
		JWTAccessExpiry:     15 * time.Minute,
		JWTRefreshExpiry:    7 * 24 * time.Hour,
		JWTIssuer:           "servqr-platform",
		OTPLength:           6,
		OTPExpiryMinutes:    5,
		OTPMaxAttempts:      3,
		OTPRateLimitPerHour: 3,
		OTPCooldownSeconds:  60,
		PasswordBcryptCost:  12,
		PasswordMinLength:   8,
		MaxFailedAttempts:   5,
		LockoutDuration:     30 * time.Minute,
		AllowRegistration:   true,
		EmailSender:         emailSender,
		SMSSender:           smsSender,
	})
	if err != nil {
		return nil, err
	}

	// Register routes
	authModule.RegisterRoutes(router)

	logger.Info("Authentication module initialized successfully")
	return authModule, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Mock services for development

type MockEmailSender struct{}

func (m *MockEmailSender) SendOTP(ctx context.Context, to, otp string) error {
	slog.Info("ðŸ“§ MOCK EMAIL", 
		slog.String("to", to), 
		slog.String("otp", otp),
		slog.String("message", "Would send OTP via email"))
	return nil
}

type MockSMSSender struct{}

func (m *MockSMSSender) SendOTP(ctx context.Context, to, otp string) error {
	slog.Info("ðŸ“± MOCK SMS",
		slog.String("to", to),
		slog.String("otp", otp),
		slog.String("message", "Would send OTP via SMS"))
	return nil
}

func (m *MockSMSSender) SendWhatsAppOTP(ctx context.Context, to, otp string) error {
	slog.Info("ðŸ’¬ MOCK WHATSAPP",
		slog.String("to", to),
		slog.String("otp", otp),
		slog.String("message", "Would send OTP via WhatsApp"))
	return nil
}
