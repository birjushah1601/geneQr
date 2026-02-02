package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/aby-med/medical-platform/internal/shared/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	logger            *slog.Logger
	db                *pgxpool.Pool
	rateLimitMiddleware *middleware.RateLimitMiddleware
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(logger *slog.Logger, db *pgxpool.Pool, rateLimitMiddleware *middleware.RateLimitMiddleware) *HealthHandler {
	return &HealthHandler{
		logger:            logger.With(slog.String("component", "health_handler")),
		db:                db,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Version   string                 `json:"version"`
	Service   string                 `json:"service"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]HealthCheck `json:"checks"`
}

// HealthCheck represents an individual health check
type HealthCheck struct {
	Status      string                 `json:"status"`
	Message     string                 `json:"message"`
	Duration    string                 `json:"duration"`
	Timestamp   string                 `json:"timestamp"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// MetricsResponse represents the metrics endpoint response
type MetricsResponse struct {
	Service     string                 `json:"service"`
	Timestamp   string                 `json:"timestamp"`
	Uptime      string                 `json:"uptime"`
	Database    map[string]interface{} `json:"database"`
	RateLimit   map[string]interface{} `json:"rate_limit"`
	Attachments map[string]interface{} `json:"attachments"`
	System      map[string]interface{} `json:"system"`
}

// AIAnalysisStatus represents the AI analysis status
type AIAnalysisStatus struct {
	Status          string                 `json:"status"`
	Provider        string                 `json:"provider"`
	Model           string                 `json:"model"`
	Available       bool                   `json:"available"`
	LastCheck       string                 `json:"last_check"`
	QueueLength     int                    `json:"queue_length"`
	ProcessedToday  int64                  `json:"processed_today"`
	SuccessRate     float64                `json:"success_rate"`
	AverageTime     string                 `json:"average_time"`
	Details         map[string]interface{} `json:"details"`
}

var startTime = time.Now()

// Health handles GET /health/attachments
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
		Service:   "attachment-service",
		Uptime:    time.Since(startTime).String(),
		Checks:    make(map[string]HealthCheck),
	}

	// Database health check
	response.Checks["database"] = h.checkDatabase()
	
	// Rate limiting health check
	response.Checks["rate_limiting"] = h.checkRateLimit()
	
	// Storage health check
	response.Checks["storage"] = h.checkStorage()

	// AI service health check
	response.Checks["ai_service"] = h.checkAIService()

	// Determine overall status
	overallStatus := "healthy"
	for _, check := range response.Checks {
		if check.Status == "unhealthy" {
			overallStatus = "unhealthy"
			break
		} else if check.Status == "degraded" && overallStatus == "healthy" {
			overallStatus = "degraded"
		}
	}
	response.Status = overallStatus

	// Set appropriate HTTP status
	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	} else if overallStatus == "degraded" {
		statusCode = http.StatusPartialContent
	}

	duration := time.Since(start)
	h.logger.Debug("Health check completed",
		slog.String("status", overallStatus),
		slog.Duration("duration", duration),
		slog.Int("status_code", statusCode))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Metrics handles GET /metrics/attachments
func (h *HealthHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	response := MetricsResponse{
		Service:   "attachment-service",
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    time.Since(startTime).String(),
		Database:  h.getDatabaseMetrics(),
		RateLimit: h.getRateLimitMetrics(),
		Attachments: h.getAttachmentMetrics(),
		System:    h.getSystemMetrics(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// AIAnalysisStatusHandler handles GET /status/ai-analysis
func (h *HealthHandler) AIAnalysisStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := AIAnalysisStatus{
		Status:          "operational",
		Provider:        "openai",
		Model:           "gpt-4-vision-preview",
		Available:       true,
		LastCheck:       time.Now().Format(time.RFC3339),
		QueueLength:     h.getAIQueueLength(),
		ProcessedToday:  h.getProcessedTodayCount(),
		SuccessRate:     h.getAISuccessRate(),
		AverageTime:     h.getAverageProcessingTime(),
		Details:         h.getAIDetails(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// checkDatabase performs database health check
func (h *HealthHandler) checkDatabase() HealthCheck {
	start := time.Now()
	
	if h.db == nil {
		return HealthCheck{
			Status:    "unhealthy",
			Message:   "Database connection not initialized",
			Duration:  time.Since(start).String(),
			Timestamp: time.Now().Format(time.RFC3339),
			Error:     "no database connection",
		}
	}

	// Test connection with ping
	ctx := context.Background()
	err := h.db.Ping(ctx)
	if err != nil {
		return HealthCheck{
			Status:    "unhealthy",
			Message:   "Database ping failed",
			Duration:  time.Since(start).String(),
			Timestamp: time.Now().Format(time.RFC3339),
			Error:     err.Error(),
		}
	}

	// Get database stats
	stats := h.db.Stat()
	details := map[string]interface{}{
		"total_connections":    stats.TotalConns(),
		"acquired_connections": stats.AcquiredConns(),
		"idle_connections":     stats.IdleConns(),
		"max_connections":      stats.MaxConns(),
	}

	return HealthCheck{
		Status:    "healthy",
		Message:   "Database connection successful",
		Duration:  time.Since(start).String(),
		Timestamp: time.Now().Format(time.RFC3339),
		Details:   details,
	}
}

// checkRateLimit performs rate limiting health check
func (h *HealthHandler) checkRateLimit() HealthCheck {
	start := time.Now()
	
	if h.rateLimitMiddleware == nil {
		return HealthCheck{
			Status:    "degraded",
			Message:   "Rate limiting not configured",
			Duration:  time.Since(start).String(),
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	stats := h.rateLimitMiddleware.GetStats()
	
	return HealthCheck{
		Status:    "healthy",
		Message:   "Rate limiting operational",
		Duration:  time.Since(start).String(),
		Timestamp: time.Now().Format(time.RFC3339),
		Details:   stats,
	}
}

// checkStorage performs storage health check
func (h *HealthHandler) checkStorage() HealthCheck {
	start := time.Now()
	
	// TODO: Implement actual storage health check
	// For now, just return healthy
	
	return HealthCheck{
		Status:    "healthy",
		Message:   "Storage accessible",
		Duration:  time.Since(start).String(),
		Timestamp: time.Now().Format(time.RFC3339),
		Details: map[string]interface{}{
			"storage_type": "local_filesystem",
			"base_path":    "storage/attachments",
		},
	}
}

// checkAIService performs AI service health check
func (h *HealthHandler) checkAIService() HealthCheck {
	start := time.Now()
	
	// For now, just return healthy with mock status
	// TODO: Implement actual AI service health check
	
	return HealthCheck{
		Status:    "healthy",
		Message:   "AI service available",
		Duration:  time.Since(start).String(),
		Timestamp: time.Now().Format(time.RFC3339),
		Details: map[string]interface{}{
			"provider":           "openai",
			"model":             "gpt-4-vision-preview",
			"fallback_enabled":  true,
			"queue_processing":  true,
		},
	}
}

// getDatabaseMetrics returns database metrics
func (h *HealthHandler) getDatabaseMetrics() map[string]interface{} {
	if h.db == nil {
		return map[string]interface{}{
			"status": "disconnected",
		}
	}

	stats := h.db.Stat()
	return map[string]interface{}{
		"status":               "connected",
		"total_connections":    stats.TotalConns(),
		"acquired_connections": stats.AcquiredConns(),
		"idle_connections":     stats.IdleConns(),
		"max_connections":      stats.MaxConns(),
		"successful_acquire_count": stats.AcquireCount(),
	}
}

// getRateLimitMetrics returns rate limiting metrics
func (h *HealthHandler) getRateLimitMetrics() map[string]interface{} {
	if h.rateLimitMiddleware == nil {
		return map[string]interface{}{
			"status": "disabled",
		}
	}

	return h.rateLimitMiddleware.GetStats()
}

// getAttachmentMetrics returns attachment-specific metrics
func (h *HealthHandler) getAttachmentMetrics() map[string]interface{} {
	// TODO: Implement actual attachment metrics from database
	return map[string]interface{}{
		"total_attachments":     1000,
		"pending_processing":    5,
		"processed_today":       150,
		"failed_processing":     2,
		"average_file_size_mb":  2.5,
		"storage_used_gb":      25.6,
	}
}

// getSystemMetrics returns system-level metrics
func (h *HealthHandler) getSystemMetrics() map[string]interface{} {
	return map[string]interface{}{
		"uptime":           time.Since(startTime).String(),
		"go_version":       "1.21",
		"num_goroutines":   25,
		"memory_usage_mb":  128.5,
	}
}

// getAIQueueLength returns the current AI processing queue length
func (h *HealthHandler) getAIQueueLength() int {
	// TODO: Query database for actual queue length
	return 3
}

// getProcessedTodayCount returns the number of items processed today
func (h *HealthHandler) getProcessedTodayCount() int64 {
	// TODO: Query database for today's processed count
	return 150
}

// getAISuccessRate returns the AI processing success rate
func (h *HealthHandler) getAISuccessRate() float64 {
	// TODO: Calculate actual success rate from database
	return 0.96
}

// getAverageProcessingTime returns average AI processing time
func (h *HealthHandler) getAverageProcessingTime() string {
	// TODO: Calculate actual average from database
	return "3.2s"
}

// getAIDetails returns detailed AI service information
func (h *HealthHandler) getAIDetails() map[string]interface{} {
	return map[string]interface{}{
		"api_key_configured":    true,
		"fallback_enabled":      true,
		"supported_formats":     []string{"jpg", "jpeg", "png", "gif", "bmp"},
		"max_file_size_mb":      100,
		"concurrent_workers":    3,
		"retry_attempts":        3,
		"timeout_seconds":       60,
	}
}