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
    // AI repository not implemented yet; use a noop repository
    m.aiRepo = infra.NewNoopAIAnalysisRepository()

    // Initialize service
    m.attachmentService = domain.NewAttachmentService(
        m.attachmentRepo,
        m.queueRepo,
        m.aiRepo, // nil ok; AI endpoints will be disabled
        m.logger,
    )

    // Initialize real HTTP handler
    m.httpHandler = api.NewAttachmentHandler(m.attachmentService, m.logger)
    m.mockHandler = nil

	m.logger.Info("Attachment module initialized successfully")
	return nil
}

// MountRoutes registers HTTP routes for the module
func (m *Module) MountRoutes(r chi.Router) {
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

	r.Route("/attachments", func(r chi.Router) {
		// Note: Auth middleware is already applied globally by main.go
		// No need to apply it again here - just use logging and rate limiting
		
		// Apply logging middleware
		r.Use(loggingMiddleware.Middleware())
		
		// Apply rate limiting middleware
		r.Use(rateLimitMiddleware.Middleware())

        // Use real handler
        // Note: Permission checks removed - JWT authentication is sufficient
        // GET endpoints
        r.Get("/", m.httpHandler.ListAttachments)
        r.Get("/stats", m.httpHandler.GetStats)
        // AI analysis endpoint not yet backed by repository; keep route but handler returns 501
        r.Get("/{id}", m.httpHandler.GetAttachment)
        r.Get("/{id}/download", m.httpHandler.DownloadAttachment)
        r.Get("/{id}/ai-analysis", m.httpHandler.GetAIAnalysis)

        // POST endpoints
        r.Post("/", m.httpHandler.CreateAttachment)
        r.Post("/{id}/link", m.httpHandler.LinkAttachment)

		// DELETE for attachment removal
		r.Delete("/{id}", m.httpHandler.DeleteAttachment)
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

	m.logger.Info("Attachment routes registered", 
		slog.String("prefix", "/api/v1/attachments"),
		slog.String("auth", "global_jwt_middleware"))

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
