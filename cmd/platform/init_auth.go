package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/aby-med/medical-platform/internal/core/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// initAuthModule initializes the authentication module if enabled
// Returns the auth module for accessing middleware, or nil if disabled/failed
func initAuthModule(router *chi.Mux, db *sqlx.DB, logger *slog.Logger) (*auth.Module, error) {
	// Check if auth module is enabled
	enableAuth := getEnvBool("ENABLE_AUTH", true) // Default: enabled
	
	if !enableAuth {
		logger.Info("Authentication module is disabled")
		return nil, nil
	}

	logger.Info("Initializing authentication module...")

	// Initialize auth module (without registering routes if router is nil)
	if router != nil {
		// Full initialization with route registration
		authModule, err := auth.IntegrateAuthModuleWithReturn(router, db, logger)
		if err != nil {
			logger.Error("Failed to initialize auth module", slog.String("error", err.Error()))
			return nil, err
		}

		logger.Info("✅ Authentication module initialized successfully",
			slog.String("endpoints", "12 auth endpoints"),
			slog.Bool("jwt_enabled", true),
			slog.Bool("otp_enabled", true))

		return authModule, nil
	} else {
		// Initialize module without route registration
		authModule, err := auth.NewModule(db, &auth.Config{
			JWTPrivateKeyPath:   getEnvOrDefault("JWT_PRIVATE_KEY_PATH", "./keys/jwt-private.pem"),
			JWTPublicKeyPath:    getEnvOrDefault("JWT_PUBLIC_KEY_PATH", "./keys/jwt-public.pem"),
			JWTAccessExpiry:     15 * time.Minute,
			JWTRefreshExpiry:    7 * 24 * time.Hour,
			JWTIssuer:           "aby-med-platform",
		})
		if err != nil {
			logger.Error("Failed to create auth module", slog.String("error", err.Error()))
			return nil, err
		}

		logger.Info("✅ Authentication module created (routes will be registered later)")
		return authModule, nil
	}
}

// Helper function to get boolean from env
func getEnvBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val == "true" || val == "1" || val == "yes" || val == "on"
}

// Helper function to get string from env with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
