package main

import (
	"log/slog"
	"os"

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

	// Initialize auth module
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
}

// Helper function to get boolean from env
func getEnvBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val == "true" || val == "1" || val == "yes" || val == "on"
}
