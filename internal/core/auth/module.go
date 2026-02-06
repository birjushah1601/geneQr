package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/aby-med/medical-platform/internal/core/auth/api"
	"github.com/aby-med/medical-platform/internal/core/auth/app"
	"github.com/aby-med/medical-platform/internal/core/auth/domain"
	"github.com/aby-med/medical-platform/internal/core/auth/infra"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Module represents the auth module with all dependencies
type Module struct {
	Handler *api.AuthHandler
}

// Config holds configuration for the auth module
type Config struct {
	// JWT
	JWTPrivateKeyPath  string
	JWTPublicKeyPath   string
	JWTAccessExpiry    time.Duration
	JWTRefreshExpiry   time.Duration
	JWTIssuer          string

	// OTP
	OTPLength           int
	OTPExpiryMinutes    int
	OTPMaxAttempts      int
	OTPRateLimitPerHour int
	OTPCooldownSeconds  int

	// Password
	PasswordBcryptCost int
	PasswordMinLength  int

	// Auth
	MaxFailedAttempts int
	LockoutDuration   time.Duration
	AllowRegistration bool

	// External services
	EmailSender EmailSender
	SMSSender   SMSSender
}

// EmailSender interface for sending emails
type EmailSender interface {
	SendOTP(ctx context.Context, to, otp string) error
}

// SMSSender interface for sending SMS
type SMSSender interface {
	SendOTP(ctx context.Context, to, otp string) error
	SendWhatsAppOTP(ctx context.Context, to, otp string) error
}

// NewModule creates and wires up the auth module
func NewModule(db *sqlx.DB, config *Config) (*Module, error) {
	// Set defaults
	if config.JWTAccessExpiry == 0 {
		config.JWTAccessExpiry = 15 * time.Minute
	}
	if config.JWTRefreshExpiry == 0 {
		config.JWTRefreshExpiry = 7 * 24 * time.Hour
	}
	if config.JWTIssuer == "" {
		config.JWTIssuer = "servqr-platform"
	}
	if config.OTPLength == 0 {
		config.OTPLength = 6
	}
	if config.OTPExpiryMinutes == 0 {
		config.OTPExpiryMinutes = 5
	}
	if config.OTPMaxAttempts == 0 {
		config.OTPMaxAttempts = 3
	}
	if config.OTPRateLimitPerHour == 0 {
		config.OTPRateLimitPerHour = 3
	}
	if config.OTPCooldownSeconds == 0 {
		config.OTPCooldownSeconds = 60
	}
	if config.PasswordBcryptCost == 0 {
		config.PasswordBcryptCost = 12
	}
	if config.PasswordMinLength == 0 {
		config.PasswordMinLength = 8
	}
	if config.MaxFailedAttempts == 0 {
		config.MaxFailedAttempts = 5
	}
	if config.LockoutDuration == 0 {
		config.LockoutDuration = 30 * time.Minute
	}

	// Load JWT keys
	privateKey, err := loadRSAPrivateKey(config.JWTPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load JWT private key: %w", err)
	}

	publicKey, err := loadRSAPublicKey(config.JWTPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load JWT public key: %w", err)
	}

	// Initialize repositories
	userRepo := infra.NewUserRepository(db)
	otpRepo := infra.NewOTPRepository(db)
	refreshTokenRepo := infra.NewRefreshTokenRepository(db)
	auditRepo := infra.NewAuditRepository(db)
	orgRepo := infra.NewOrganizationRepository(db)

	// Initialize services
	otpService := app.NewOTPService(
		otpRepo,
		auditRepo,
		config.EmailSender,
		config.SMSSender,
		&app.OTPConfig{
			Length:           config.OTPLength,
			ExpiryMinutes:    config.OTPExpiryMinutes,
			MaxAttempts:      config.OTPMaxAttempts,
			RateLimitPerHour: config.OTPRateLimitPerHour,
			CooldownSeconds:  config.OTPCooldownSeconds,
		},
	)

	// Create adapter for refresh token repository
	refreshTokenRepoAdapter := &refreshTokenRepoAdapter{repo: refreshTokenRepo}

	jwtService := app.NewJWTService(
		&app.JWTConfig{
			PrivateKey:         privateKey,
			PublicKey:          publicKey,
			AccessTokenExpiry:  config.JWTAccessExpiry,
			RefreshTokenExpiry: config.JWTRefreshExpiry,
			Issuer:             config.JWTIssuer,
		},
		refreshTokenRepoAdapter,
	)

	passwordService := app.NewPasswordService(&app.PasswordConfig{
		BcryptCost:     config.PasswordBcryptCost,
		MinLength:      config.PasswordMinLength,
		RequireUpper:   true,
		RequireLower:   true,
		RequireNumber:  true,
		RequireSpecial: true,
	})

	authService := app.NewAuthService(
		userRepo,
		otpService,
		jwtService,
		passwordService,
		auditRepo,
		orgRepo,
		&app.AuthConfig{
			MaxFailedAttempts: config.MaxFailedAttempts,
			LockoutDuration:   config.LockoutDuration,
			AllowRegistration: config.AllowRegistration,
		},
	)

	// Initialize handler
	handler := api.NewAuthHandler(authService, jwtService)

	return &Module{
		Handler: handler,
	}, nil
}

// RegisterRoutes registers auth routes
func (m *Module) RegisterRoutes(r chi.Router) {
	m.Handler.RegisterRoutes(r)
}

// loadRSAPrivateKey loads RSA private key from file
func loadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Try PKCS8 first (modern format with "BEGIN PRIVATE KEY")
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Fallback to PKCS1 (old format with "BEGIN RSA PRIVATE KEY")
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		return privateKey, nil
	}

	// PKCS8 can contain different key types, ensure it's RSA
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}

	return privateKey, nil
}

// loadRSAPublicKey loads RSA public key from file
func loadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPublicKey, nil
}

// refreshTokenRepoAdapter adapts domain.RefreshTokenRepository to app.RefreshTokenRepository
type refreshTokenRepoAdapter struct {
	repo domain.RefreshTokenRepository
}

func (a *refreshTokenRepoAdapter) Create(ctx context.Context, token *app.RefreshToken) error {
	ipAddr := token.IPAddress
	domainToken := &domain.RefreshToken{
		ID:         token.ID,
		UserID:     token.UserID,
		TokenHash:  token.TokenHash,
		DeviceInfo: token.DeviceInfo,
		IPAddress:  &ipAddr,
		ExpiresAt:  token.ExpiresAt,
	}
	return a.repo.Create(ctx, domainToken)
}

func (a *refreshTokenRepoAdapter) GetByTokenHash(ctx context.Context, tokenHash string) (*app.RefreshToken, error) {
	domainToken, err := a.repo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	ipAddr := ""
	if domainToken.IPAddress != nil {
		ipAddr = *domainToken.IPAddress
	}
	return &app.RefreshToken{
		ID:         domainToken.ID,
		UserID:     domainToken.UserID,
		TokenHash:  domainToken.TokenHash,
		ExpiresAt:  domainToken.ExpiresAt,
		DeviceInfo: domainToken.DeviceInfo,
		IPAddress:  ipAddr,
	}, nil
}

func (a *refreshTokenRepoAdapter) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	return a.repo.UpdateLastUsed(ctx, id)
}

func (a *refreshTokenRepoAdapter) Revoke(ctx context.Context, id uuid.UUID, reason string) error {
	return a.repo.Revoke(ctx, id, reason)
}

func (a *refreshTokenRepoAdapter) RevokeAllForUser(ctx context.Context, userID uuid.UUID, reason string) error {
	return a.repo.RevokeAllForUser(ctx, userID, reason)
}
