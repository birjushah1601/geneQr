package attachment

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aby-med/medical-platform/internal/service-domain/attachment/api"
	"github.com/aby-med/medical-platform/internal/service-domain/attachment/domain"
	"github.com/aby-med/medical-platform/internal/service-domain/attachment/infra"
	"github.com/aby-med/medical-platform/internal/shared/middleware"
	"github.com/aby-med/medical-platform/internal/shared/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Module represents the attachment service module
type Module struct {
	config          *Config
	logger          *slog.Logger
	db              *pgxpool.Pool
	attachmentRepo  domain.AttachmentRepository
	queueRepo       domain.ProcessingQueueRepository
	aiRepo          domain.AIAnalysisRepository
	attachmentService *domain.AttachmentService
	httpHandler     *api.AttachmentHandler
	mockHandler     *api.MockAttachmentHandler
	healthHandler   *api.HealthHandler
	rateLimitMiddleware *middleware.RateLimitMiddleware
	loggingMiddleware *middleware.LoggingMiddleware
}

// Config holds configuration for the attachment module
type Config struct {
	DatabaseDSN string
}

// NewModule creates a new attachment module instance
func NewModule(config Config, logger *slog.Logger) *Module {
	return &Module{
		config: &config,
		logger: logger.With(slog.String("module", "attachment")),
	}
}

// Name returns the module name
func (m *Module) Name() string {
	return "attachment"
}

// Initialize sets up the module dependencies
func (m *Module) Initialize(ctx context.Context) error {
	m.logger.Info("Initializing attachment module")

	// Connect to database
	db, err := pgxpool.New(ctx, m.config.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	m.db = db

	// Test connection
	if err := m.db.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize repositories
	m.attachmentRepo = infra.NewPostgresAttachmentRepository(m.db)
	m.queueRepo = infra.NewPostgresProcessingQueueRepository(m.db)
	// TODO: Initialize AI repository when needed

	// Initialize service (skip stats for now due to repository issues)
	// m.attachmentService = domain.NewAttachmentService(
	//	m.attachmentRepo,
	//	m.queueRepo,
	//	m.aiRepo, // nil for now
	//	m.logger,
	// )

	// Temporarily skip service initialization to test API routes
	m.logger.Info("Skipping attachment service initialization for testing")

	// Initialize mock HTTP handler for testing
	m.mockHandler = api.NewMockAttachmentHandler(m.logger)

	m.logger.Info("Attachment module initialized successfully")
	return nil
}

// MountRoutes registers HTTP routes for the module
func (m *Module) MountRoutes(r chi.Router) {
	// Create authentication middleware
	authConfig := middleware.DefaultAuthConfig()
	authMiddleware := middleware.NewAuthMiddleware(authConfig, m.logger)

	// Create rate limiting middleware
	rateLimitConfig := middleware.DefaultRateLimitConfig()
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(rateLimitConfig, m.logger)
	m.rateLimitMiddleware = rateLimitMiddleware

	// Create logging middleware
	loggingConfig := middleware.DefaultLoggingConfig()
	loggingMiddleware := middleware.NewLoggingMiddleware(m.logger, loggingConfig)
	m.loggingMiddleware = loggingMiddleware

	// Create health handler
	m.healthHandler = api.NewHealthHandler(m.logger, m.db, rateLimitMiddleware)

	r.Route("/api/v1/attachments", func(r chi.Router) {
		// Apply logging middleware first
		r.Use(loggingMiddleware.Middleware())
		
		// Apply authentication middleware second
		r.Use(authMiddleware.Middleware())
		
		// Then apply rate limiting middleware
		r.Use(rateLimitMiddleware.Middleware())

		// Use mock handler for testing
		if m.mockHandler != nil {
			// GET endpoints require 'read' permission
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission("read", m.logger))
				r.Get("/", m.mockHandler.ListAttachments)
				r.Get("/stats", m.mockHandler.GetStats)
				r.Get("/{id}", m.mockHandler.GetAttachment)
				r.Get("/{id}/ai-analysis", m.mockHandler.GetAIAnalysis)
			})

			// POST endpoints require 'upload' permission
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission("upload", m.logger))
				r.Post("/", m.mockHandler.CreateAttachment)
			})
		} else {
			// GET endpoints require 'read' permission
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission("read", m.logger))
				r.Get("/", m.httpHandler.ListAttachments)
				r.Get("/stats", m.httpHandler.GetStats)
				r.Get("/{id}", m.httpHandler.GetAttachment)
				r.Get("/{id}/ai-analysis", m.httpHandler.GetAIAnalysis)
			})

			// POST endpoints require 'upload' permission
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequirePermission("upload", m.logger))
				r.Post("/", m.httpHandler.CreateAttachment)
			})
		}

		// TODO: Add DELETE for attachment removal (requires 'delete' permission)
		// r.Group(func(r chi.Router) {
		//     r.Use(middleware.RequirePermission("delete", m.logger))
		//     r.Delete("/{id}", m.httpHandler.DeleteAttachment)
		// })
	})

	// Health check endpoints (no auth required)
	r.Route("/health", func(r chi.Router) {
		r.Get("/attachments", m.healthHandler.Health)
	})

	r.Route("/metrics", func(r chi.Router) {
		r.Get("/attachments", m.healthHandler.Metrics)
	})

	r.Route("/status", func(r chi.Router) {
		r.Get("/ai-analysis", m.healthHandler.AIAnalysisStatusHandler)
	})

	m.logger.Info("Attachment routes registered with authentication", 
		slog.String("prefix", "/api/v1/attachments"),
		slog.String("auth", "enabled"),
		slog.Int("api_keys", len(authConfig.APIKeys)))

	m.logger.Info("Health check routes registered", 
		slog.String("health_endpoint", "/health/attachments"),
		slog.String("metrics_endpoint", "/metrics/attachments"),
		slog.String("ai_status_endpoint", "/status/ai-analysis"))
}

// Start starts background processes (if any)
func (m *Module) Start(ctx context.Context) error {
	m.logger.Info("Attachment module started")
	// TODO: Start any background workers (queue processor, etc.)
	return nil
}

// Stop gracefully stops the module
func (m *Module) Stop(ctx context.Context) error {
	m.logger.Info("Shutting down attachment module")
	
	// Stop rate limiting middleware
	if m.rateLimitMiddleware != nil {
		m.rateLimitMiddleware.Stop()
	}
	
	if m.db != nil {
		m.db.Close()
	}
	
	return nil
}

// Health returns the health status
func (m *Module) Health(ctx context.Context) error {
	// TODO: Check database connection
	return nil
}

// Ensure Module implements the service.Module interface
var _ service.Module = (*Module)(nil)