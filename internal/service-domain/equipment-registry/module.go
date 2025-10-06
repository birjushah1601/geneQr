package equipment

import (
	"context"
	"log/slog"

	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/api"
	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/app"
	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/infra"
	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/qrcode"
	"github.com/go-chi/chi/v5"
)

// Module represents the equipment registry module
type Module struct {
	config  ModuleConfig
	handler *api.EquipmentHandler
	logger  *slog.Logger
}

// ModuleConfig holds module configuration
type ModuleConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	BaseURL    string
	QROutputDir string
}

// NewModule creates a new equipment registry module
func NewModule(cfg ModuleConfig, logger *slog.Logger) (*Module, error) {
	return &Module{
		config: cfg,
		logger: logger.With(slog.String("module", "equipment-registry")),
	}, nil
}

// Initialize initializes the module (database connections, etc.)
func (m *Module) Initialize(ctx context.Context) error {
	m.logger.Info("Initializing Equipment Registry module")

	// Create database connection pool
	dbConfig := infra.PostgresConfig{
		Host:     m.config.DBHost,
		Port:     m.config.DBPort,
		User:     m.config.DBUser,
		Password: m.config.DBPassword,
		Database: m.config.DBName,
		SSLMode:  "disable",
	}

    pool, err := infra.NewPostgresPool(ctx, dbConfig)
	if err != nil {
		return err
	}

	// Create repository
	repo := infra.NewEquipmentRepository(pool)

	// Create QR code generator
	qrGenerator := qrcode.NewGenerator(m.config.BaseURL, m.config.QROutputDir)

	// Create application service
	service := app.NewEquipmentService(repo, qrGenerator, m.logger, m.config.BaseURL)

	// Create HTTP handler
	m.handler = api.NewEquipmentHandler(service, m.logger)

    // Ensure schema is compatible with application expectations
    if err := infra.EnsureEquipmentSchema(ctx, pool); err != nil {
        return err
    }

	m.logger.Info("Equipment Registry module initialized successfully")
	return nil
}

// MountRoutes mounts module routes to the router
func (m *Module) MountRoutes(r chi.Router) {
	m.logger.Info("Mounting Equipment Registry routes")

	// Equipment management routes
	r.Route("/equipment", func(r chi.Router) {
		// Routes without {id} parameter first
		r.Post("/", m.handler.RegisterEquipment)          // Register equipment
		r.Get("/", m.handler.ListEquipment)               // List equipment
		r.Post("/import", m.handler.ImportCSV)            // CSV import
		r.Post("/qr/bulk-generate", m.handler.BulkGenerateQRCodes) // Bulk generate QR codes
		r.Get("/qr/image/{id}", m.handler.GetQRCodeImage) // Get QR code image (different pattern to avoid conflict)
		r.Get("/qr/{qr_code}", m.handler.GetEquipmentByQR) // Get by QR code
		r.Get("/serial/{serial}", m.handler.GetEquipmentBySerial) // Get by serial
		
		// {id} sub-routes
		r.Get("/{id}/qr/pdf", m.handler.DownloadQRLabel)   // Download PDF label
		r.Post("/{id}/qr", m.handler.GenerateQRCode)       // Generate QR code
		r.Post("/{id}/service", m.handler.RecordService)   // Record service
		
		// Base /{id} routes LAST
		r.Get("/{id}", m.handler.GetEquipment)            // Get by ID
		r.Patch("/{id}", m.handler.UpdateEquipment)       // Update equipment
	})

	m.logger.Info("Equipment Registry routes mounted successfully")
}

// Start starts background tasks (if any)
func (m *Module) Start(ctx context.Context) error {
	m.logger.Info("Equipment Registry module started")
	return nil
}

// Stop gracefully stops the module
func (m *Module) Stop(ctx context.Context) error {
	m.logger.Info("Equipment Registry module stopped")
	return nil
}

// Health returns the health status
func (m *Module) Health(ctx context.Context) error {
	// TODO: Check database connection
	return nil
}

// Name returns the module name
func (m *Module) Name() string {
	return "equipment-registry"
}
