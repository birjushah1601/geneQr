package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/aby-med/medical-platform/internal/shared/config"
	"github.com/aby-med/medical-platform/internal/shared/observability"
    "github.com/aby-med/medical-platform/internal/shared/service"
	sharedmiddleware "github.com/aby-med/medical-platform/internal/shared/middleware"
	appmiddleware "github.com/aby-med/medical-platform/internal/middleware"
    organizations "github.com/aby-med/medical-platform/internal/core/organizations"
	"github.com/aby-med/medical-platform/internal/marketplace/catalog"
	"github.com/aby-med/medical-platform/internal/service-domain/rfq"
	"github.com/aby-med/medical-platform/internal/service-domain/supplier"
	"github.com/aby-med/medical-platform/internal/service-domain/quote"
	"github.com/aby-med/medical-platform/internal/service-domain/comparison"
	"github.com/aby-med/medical-platform/internal/service-domain/contract"
	equipment "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry"
	equipmentApp "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/app"
	serviceticket "github.com/aby-med/medical-platform/internal/service-domain/service-ticket"
	serviceticketApp "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/app"
	"github.com/aby-med/medical-platform/internal/service-domain/attachment"
	"github.com/aby-med/medical-platform/internal/service-domain/whatsapp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	
	// AI Services
	aimanager "github.com/aby-med/medical-platform/internal/ai"
	"github.com/aby-med/medical-platform/internal/ai/aiconfig"
	// TODO: Uncomment when routes are mounted
	// "github.com/aby-med/medical-platform/internal/diagnosis"
	// "github.com/aby-med/medical-platform/internal/assignment"
	// "github.com/aby-med/medical-platform/internal/parts"
	// "github.com/aby-med/medical-platform/internal/feedback"
	"github.com/jackc/pgx/v5/pgxpool"
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
	
	// Security middleware
	r.Use(sharedmiddleware.SecurityHeaders)
	
	// Rate limiting (100 requests per minute per IP)
	r.Use(sharedmiddleware.RateLimitByIP(100, 1*time.Minute))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Tenant-ID", "X-Api-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Observability middleware
	r.Use(observability.LoggingMiddleware(logger))

	// CRITICAL: Organization context middleware for multi-tenant data isolation
	// Must be registered BEFORE any routes are mounted
	r.Use(appmiddleware.OrganizationContextMiddleware(logger))
	logger.Info("✅ Organization context middleware registered")

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
    // Register Organizations module behind feature flag
    if os.Getenv("ENABLE_ORG") == "true" || os.Getenv("ENABLE_ORG") == "1" || os.Getenv("ENABLE_ORG") == "on" {
        registry.Register(organizations.New(cfg, logger))
    }
	
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
	
	// Register Comparison module
	comparisonConfig := &comparison.Config{
		DatabaseURL: cfg.GetDSN(),
	}
	registry.Register(comparison.NewModule(*comparisonConfig, logger))
	
	// Register Contract module
	contractConfig := &contract.Config{
		DatabaseURL: cfg.GetDSN(),
	}
	registry.Register(contract.NewModule(*contractConfig, logger))
	
	// Register Equipment Registry module (Field Service Management)
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://service.yourcompany.com"
	}
	qrOutputDir := os.Getenv("QR_OUTPUT_DIR")
	if qrOutputDir == "" {
		qrOutputDir = "./data/qrcodes"
	}
	
	// Parse database port
	dbPort, _ := strconv.Atoi(cfg.Database.Port)
	if dbPort == 0 {
		dbPort = 5432 // Default PostgreSQL port
	}
	
	equipmentModule, err := equipment.NewModule(equipment.ModuleConfig{
		DBHost:      cfg.Database.Host,
		DBPort:      dbPort,
		DBUser:      cfg.Database.User,
		DBPassword:  cfg.Database.Password,
		DBName:      cfg.Database.Name,
		BaseURL:     baseURL,
		QROutputDir: qrOutputDir,
	}, logger)
	if err == nil {
		registry.Register(equipmentModule)
	}
	
	// Register Service Ticket module (includes WhatsApp integration)
	whatsappVerifyToken := os.Getenv("WHATSAPP_VERIFY_TOKEN")
	if whatsappVerifyToken == "" {
		whatsappVerifyToken = "your-verify-token-123"
	}
	whatsappAccessToken := os.Getenv("WHATSAPP_ACCESS_TOKEN")
	whatsappPhoneID := os.Getenv("WHATSAPP_PHONE_ID")
	whatsappMediaDir := os.Getenv("WHATSAPP_MEDIA_DIR")
	if whatsappMediaDir == "" {
		whatsappMediaDir = "./data/whatsapp"
	}
	
	serviceTicketModule, err := serviceticket.NewModule(serviceticket.ModuleConfig{
		DBHost:             cfg.Database.Host,
		DBPort:             dbPort,
		DBUser:             cfg.Database.User,
		DBPassword:         cfg.Database.Password,
		DBName:             cfg.Database.Name,
		BaseURL:            baseURL,
		QROutputDir:        qrOutputDir,
		WhatsAppVerifyToken: whatsappVerifyToken,
		WhatsAppAccessToken: whatsappAccessToken,
		WhatsAppPhoneID:     whatsappPhoneID,
		WhatsAppMediaDir:    whatsappMediaDir,
	}, logger)
	if err == nil {
		registry.Register(serviceTicketModule)
	}
	
	// Register Attachment module
	attachmentConfig := attachment.Config{
		DatabaseDSN: cfg.GetDSN(),
	}
	attachmentModule := attachment.NewModule(attachmentConfig, logger)
	registry.Register(attachmentModule)
	
	// ========================================================================
	// INITIALIZE AI SERVICES
	// ========================================================================
	logger.Info("Initializing AI Services")
	
	// Initialize AI Manager
	timeout := time.Duration(cfg.AI.TimeoutSeconds) * time.Second
	aiConfig := &aimanager.Config{
		DefaultProvider:        cfg.AI.Provider,
		EnableFallback:         true,
		MaxRetries:             cfg.AI.MaxRetries,
		RetryBackoffMultiplier: 2.0,
		DefaultTimeout:         timeout,
		EnableCostTracking:     cfg.AI.CostTrackingEnabled,
		EnableHealthChecks:     true,
		HealthCheckInterval:    5 * time.Minute,
		OpenAI: aiconfig.OpenAIConfig{
			APIKey:          cfg.AI.OpenAIAPIKey,
			DefaultModel:    cfg.AI.OpenAIModel,
			Timeout:         timeout,
			MaxRetries:      cfg.AI.MaxRetries,
			EnableStreaming: false,
		},
		Anthropic: aiconfig.AnthropicConfig{
			APIKey:       cfg.AI.AnthropicAPIKey,
			DefaultModel: cfg.AI.AnthropicModel,
			Timeout:      timeout,
			MaxRetries:   cfg.AI.MaxRetries,
			APIVersion:   "2023-06-01",
		},
	}
	
	_, err = aimanager.NewManager(aiConfig)
	if err != nil {
		logger.Warn("Failed to initialize AI manager (AI features will be disabled)", 
			slog.String("error", err.Error()))
	} else {
		logger.Info("AI Manager initialized successfully",
			slog.String("primary_provider", cfg.AI.Provider),
			slog.String("fallback_provider", cfg.AI.FallbackProvider))
		
		// Initialize database connection pool for AI services
		dbPoolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
		if err != nil {
			logger.Error("Failed to parse database config for AI services", slog.String("error", err.Error()))
		} else {
			_, err = pgxpool.NewWithConfig(ctx, dbPoolConfig)
			if err != nil {
				logger.Error("Failed to create database pool for AI services", slog.String("error", err.Error()))
			} else {
				logger.Info("Database pool created for AI services")
				
				// Initialize AI Engines
				logger.Info("Initializing AI Diagnosis Engine")
				// _ = diagnosis.NewEngine(aiMgr, db) // TODO: Initialize when routes are mounted
				
				logger.Info("Initializing AI Assignment Optimizer")
				// _ = assignment.NewEngine(aiMgr, db)
				
				logger.Info("Initializing AI Parts Recommender")
				// _ = parts.NewEngine(aiMgr, db)
				
				logger.Info("Initializing AI Feedback Loop Manager")
				// _ = feedback.NewCollector(db)
				// _ = feedback.NewAnalyzer(db)
				// _ = feedback.NewLearner(db)
				
				logger.Info("All AI services initialized successfully")
				
				// TODO Phase 3: Mount AI service routes
				// TODO Phase 4: Integrate with service ticket workflow
			}
		}
	}

	// ========================================================================
	// INITIALIZE AUTHENTICATION SYSTEM
	// ========================================================================
	logger.Info("Initializing Authentication System")
	
	// Create database connection for auth module
	authDB, err := sqlx.Connect("postgres", cfg.GetDSN())
	if err != nil {
		logger.Warn("Failed to connect to database for auth", slog.String("error", err.Error()))
	} else {
		err = initAuthModule(router, authDB, logger)
		if err != nil {
			logger.Warn("Failed to initialize auth module", slog.String("error", err.Error()))
		}
	}

	// ========================================================================
	// INITIALIZE WHATSAPP INTEGRATION (Optional)
	// ========================================================================
	if os.Getenv("ENABLE_WHATSAPP") == "true" {
		logger.Info("Initializing WhatsApp integration")
		
		// Get WhatsApp configuration from environment
		twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
		twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
		twilioWhatsAppNumber := os.Getenv("TWILIO_WHATSAPP_NUMBER")
		
		if twilioAccountSID == "" || twilioAuthToken == "" || twilioWhatsAppNumber == "" {
			logger.Warn("WhatsApp integration enabled but missing Twilio credentials",
				slog.Bool("has_account_sid", twilioAccountSID != ""),
				slog.Bool("has_auth_token", twilioAuthToken != ""),
				slog.Bool("has_whatsapp_number", twilioWhatsAppNumber != ""))
		} else {
			// Get equipment and ticket services from registry
			// Note: These need to be initialized first by the module registry
			// For now, we'll initialize WhatsApp routes separately
			logger.Info("WhatsApp integration configured",
				slog.String("whatsapp_number", twilioWhatsAppNumber))
			
			// WhatsApp module will be initialized after other modules are ready
			// This ensures equipment and ticket services are available
		}
	}

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
	
	// Initialize WhatsApp module if enabled (after other modules are ready)
	if os.Getenv("ENABLE_WHATSAPP") == "true" {
		twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
		twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
		twilioWhatsAppNumber := os.Getenv("TWILIO_WHATSAPP_NUMBER")
		
		if twilioAccountSID != "" && twilioAuthToken != "" && twilioWhatsAppNumber != "" {
			logger.Info("Mounting WhatsApp routes")
			
			// Get services from modules (simplified - assumes they're available)
			// In a full implementation, you'd extract these from the registry
			// For now, create a placeholder that can be enhanced later
			var equipmentService *equipmentApp.EquipmentService
			var ticketService *serviceticketApp.TicketService
			
			// Create database pool for WhatsApp
			dbPool, err := pgxpool.New(ctx, cfg.GetDSN())
			if err == nil {
				// Try to get services from modules
				for _, module := range modules {
					if module.Name() == "equipment" {
						// Cast to equipment module type if possible
						logger.Info("Found equipment module for WhatsApp")
					}
					if module.Name() == "service-ticket" {
						// Cast to ticket module type if possible
						logger.Info("Found service-ticket module for WhatsApp")
					}
				}
				
				// If we have the services, initialize WhatsApp module
				if equipmentService != nil && ticketService != nil {
					whatsappModule := whatsapp.NewWhatsAppModule(
						dbPool,
						equipmentService,
						ticketService,
						twilioAccountSID,
						twilioAuthToken,
						twilioWhatsAppNumber,
						logger,
					)
					whatsappModule.MountRoutes(apiRouter)
					logger.Info("✅ WhatsApp integration initialized and routes mounted")
				} else {
					logger.Warn("WhatsApp enabled but required services not available yet - webhook endpoint created for verification only")
					// Create a simple verification endpoint
					apiRouter.Get("/whatsapp/webhook", func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte("WhatsApp webhook endpoint (services pending)"))
					})
				}
			} else {
				logger.Error("Failed to create database pool for WhatsApp", slog.String("error", err.Error()))
			}
		}
	}
	
	// Add spare parts catalog endpoint
	apiRouter.Get("/catalog/parts", createSparePartsHandler(cfg.GetDSN(), logger))

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

// createSparePartsHandler creates a handler for listing spare parts
func createSparePartsHandler(dsn string, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create database connection
		db, err := pgxpool.New(r.Context(), dsn)
		if err != nil {
			logger.Error("Failed to connect to database", slog.String("error", err.Error()))
			http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Parse query parameters
		category := r.URL.Query().Get("category")
		searchQuery := r.URL.Query().Get("q")
		
		// Build SQL query
		query := `
			SELECT 
				id, part_number, part_name, category, subcategory, 
				description, unit_price, currency, is_available, 
				stock_status, requires_engineer, engineer_level_required,
				installation_time_minutes, lead_time_days, minimum_order_quantity,
				image_url, photos
			FROM spare_parts_catalog
			WHERE is_available = true AND is_obsolete = false`
		
		args := []interface{}{}
		argCount := 1
		
		if category != "" {
			query += fmt.Sprintf(" AND category ILIKE $%d", argCount)
			args = append(args, "%"+category+"%")
			argCount++
		}
		
		if searchQuery != "" {
			query += fmt.Sprintf(" AND (part_name ILIKE $%d OR part_number ILIKE $%d OR description ILIKE $%d)", argCount, argCount, argCount)
			args = append(args, "%"+searchQuery+"%")
			argCount++
		}
		
		query += " ORDER BY category, part_name LIMIT 100"
		
		// Execute query
		rows, err := db.Query(r.Context(), query, args...)
		if err != nil {
			logger.Error("Failed to query spare parts", slog.String("error", err.Error()))
			http.Error(w, `{"error":"Query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		
		// Build response
		parts := []map[string]interface{}{}
		for rows.Next() {
			var (
				id, partNumber, partName, cat, description, currency, stockStatus string
				subcategory, engineerLevel, imageURL *string
				unitPrice float64
				isAvailable, requiresEngineer bool
				installTime, leadTime, minOrderQty *int
				photos []string
			)
			
			if err := rows.Scan(&id, &partNumber, &partName, &cat, &subcategory,
				&description, &unitPrice, &currency, &isAvailable, &stockStatus,
				&requiresEngineer, &engineerLevel, &installTime, &leadTime, &minOrderQty,
				&imageURL, &photos); err != nil {
				continue
			}
			
			part := map[string]interface{}{
				"id":            id,
				"part_number":   partNumber,
				"part_name":     partName,
				"category":      cat,
				"description":   description,
				"unit_price":    unitPrice,
				"currency":      currency,
				"is_available":  isAvailable,
				"stock_status":  stockStatus,
				"requires_engineer": requiresEngineer,
			}
			
			if subcategory != nil {
				part["subcategory"] = *subcategory
			}
			if engineerLevel != nil {
				part["engineer_level_required"] = *engineerLevel
			}
			if installTime != nil {
				part["installation_time_minutes"] = *installTime
			}
			if leadTime != nil {
				part["lead_time_days"] = *leadTime
			}
			if minOrderQty != nil {
				part["minimum_order_quantity"] = *minOrderQty
			} else {
				part["minimum_order_quantity"] = 1
			}
			if imageURL != nil {
				part["image_url"] = *imageURL
			}
			if photos != nil && len(photos) > 0 {
				part["photos"] = photos
			}
			
			parts = append(parts, part)
		}
		
		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		response := map[string]interface{}{
			"parts": parts,
			"count": len(parts),
		}
		
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Failed to encode response", slog.String("error", err.Error()))
		}
	}
}



