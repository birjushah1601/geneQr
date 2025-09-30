package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aby-med/medical-platform/internal/shared/config"
	"github.com/aby-med/medical-platform/internal/shared/observability"
	"github.com/aby-med/medical-platform/internal/shared/service"
	"github.com/aby-med/medical-platform/internal/marketplace/catalog"
	"github.com/aby-med/medical-platform/internal/service-domain/rfq"
	"github.com/aby-med/medical-platform/internal/service-domain/supplier"
	"github.com/aby-med/medical-platform/internal/service-domain/quote"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

const (
	defaultPort        = "8081"
	defaultShutdownSec = 30
	allModulesWildcard = "*"
)

func main() {
	// Initialize context with cancellation for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Load environment variables
	if err := loadEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load environment: %v\n", err)
		os.Exit(1)
	}

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize observability (logging, metrics, tracing)
	obsCfg := observability.Config{
		LogLevel:       cfg.Observability.LogLevel,
		TracingEnabled: cfg.Observability.TracingEnabled,
		TracingURL:     cfg.Observability.TracingURL,
		SamplingRate:   cfg.Observability.SamplingRate,
		MetricsEnabled: cfg.Observability.MetricsEnabled,
		ServiceName:    "medical-platform",
		Environment:    cfg.Environment,
	}

	logger, tracer, _, err := observability.Setup(ctx, obsCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize observability: %v\n", err)
		os.Exit(1)
	}
	defer observability.Shutdown(ctx)

	// Log startup information
	logger.Info("Starting medical platform",
		slog.String("version", cfg.Version),
		slog.String("environment", cfg.Environment),
		slog.String("enabled_modules", cfg.EnabledModulesString))

	// Parse enabled modules
	enabledModules := parseEnabledModules(cfg.EnabledModulesString)
	logger.Info("Modules enabled", slog.Any("modules", enabledModules))

	// Initialize router with middleware
	router := setupRouter(cfg, logger, tracer)

	// Initialize modules (non-blocking)
	modules, modulesCtx, err := initializeModules(ctx, router, enabledModules, cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize modules", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Start HTTP server
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("HTTP server listening", slog.String("port", cfg.Port))
		serverErrors <- server.ListenAndServe()
	}()

	// Start all module background processes
	go func() {
		if err := startModuleBackgroundProcesses(modulesCtx, modules, logger); err != nil {
			logger.Error("Module background processes failed", slog.String("error", err.Error()))
			cancel() // Trigger shutdown
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", slog.String("error", err.Error()))
		}
	case <-ctx.Done():
		logger.Info("Shutdown signal received")
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeoutSec)*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}

// loadEnv loads environment variables from .env files
func loadEnv() error {
	// Try loading from .env.local first, then fall back to .env
	if err := godotenv.Load(".env.local"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			// Not finding any .env file is not a fatal error
			return nil
		}
	}
	return nil
}

// parseEnabledModules parses the ENABLED_MODULES environment variable
func parseEnabledModules(modulesStr string) []string {
	if modulesStr == "" || modulesStr == allModulesWildcard {
		// Return empty slice to indicate all modules
		return []string{allModulesWildcard}
	}

	// Split by comma and trim spaces
	modules := strings.Split(modulesStr, ",")
	for i, m := range modules {
		modules[i] = strings.TrimSpace(m)
	}
	return modules
}

// setupRouter initializes the HTTP router with middleware
func setupRouter(cfg *config.Config, logger *slog.Logger, tracer observability.Tracer) *chi.Mux {
	r := chi.NewRouter()
	// tracer is currently unused (otelchi middleware removed); keep reference to avoid linter warnings
	_ = tracer

	// Core middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Tenant-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Observability middleware
	r.Use(observability.LoggingMiddleware(logger))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Metrics endpoint
	r.Get("/metrics", observability.MetricsHandler())

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// Module routes will be mounted here
		})
	})

	return r
}

// initializeModules initializes all enabled modules and mounts their routes (non-blocking)
func initializeModules(ctx context.Context, router *chi.Mux, enabledModules []string, cfg *config.Config, logger *slog.Logger) ([]service.Module, context.Context, error) {
	// Get all available modules
	registry := service.NewRegistry(cfg, logger)
	// Register individual modules here
	registry.Register(catalog.New(cfg, logger))
	
	// Register RFQ module
	rfqConfig := &rfq.Config{
		DatabaseDSN:  cfg.GetDSN(),
		KafkaBrokers: cfg.Kafka.Brokers,
	}
	registry.Register(rfq.NewModule(rfqConfig, logger))
	
	// Register Supplier module
	supplierConfig := &supplier.Config{
		DatabaseDSN: cfg.GetDSN(),
	}
	registry.Register(supplier.NewModule(supplierConfig, logger))
	
	// Register Quote module
	quoteConfig := &quote.Config{
		DatabaseURL: cfg.GetDSN(),
	}
	registry.Register(quote.NewModule(*quoteConfig, logger))

	modules, err := registry.GetModules(enabledModules)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get modules: %w", err)
	}

	if len(modules) == 0 {
		return nil, nil, fmt.Errorf("no modules enabled")
	}

	// Initialize each module first (sets up DB, handlers, etc.)
	for _, module := range modules {
		moduleName := module.Name()
		logger.Info("Initializing module", slog.String("module", moduleName))
		if err := module.Initialize(ctx); err != nil {
			return nil, nil, fmt.Errorf("module %s failed to initialize: %w", moduleName, err)
		}
	}

	// Mount routes for each module
	apiRouter := router.Route("/api/v1", func(r chi.Router) {})
	for _, module := range modules {
		moduleName := module.Name()
		logger.Info("Mounting routes for module", slog.String("module", moduleName))
		module.MountRoutes(apiRouter)
	}

	return modules, ctx, nil
}

// startModuleBackgroundProcesses starts all module background processes (blocking)
func startModuleBackgroundProcesses(ctx context.Context, modules []service.Module, logger *slog.Logger) error {
	// Start each module in parallel
	g, ctx := errgroup.WithContext(ctx)
	for _, module := range modules {
		m := module // Capture variable for goroutine
		g.Go(func() error {
			moduleName := m.Name()
			logger.Info("Starting module background processes", slog.String("module", moduleName))
			if err := m.Start(ctx); err != nil {
				return fmt.Errorf("module %s failed to start: %w", moduleName, err)
			}
			return nil
		})
	}

	// Wait for all modules to start or for first error
	return g.Wait()
}
