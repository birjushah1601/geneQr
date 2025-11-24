package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/aby-med/medical-platform/internal/shared/monitoring"
	"github.com/google/uuid"
)

// LoggingMiddleware provides structured logging for HTTP requests
type LoggingMiddleware struct {
	logger           *monitoring.StructuredLogger
	config           LoggingConfig
	skipPaths        map[string]bool
}

// LoggingConfig holds configuration for logging middleware
type LoggingConfig struct {
	LogLevel         monitoring.LogLevel `json:"log_level"`
	LogRequestBody   bool                `json:"log_request_body"`
	LogResponseBody  bool                `json:"log_response_body"`
	LogHeaders       bool                `json:"log_headers"`
	EnableTracing    bool                `json:"enable_tracing"`
	SkipPaths        []string           `json:"skip_paths"`
	MaxBodySize      int64              `json:"max_body_size"` // Max size in bytes to log
}

// DefaultLoggingConfig returns default logging middleware configuration
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		LogLevel:         monitoring.LogLevelInfo,
		LogRequestBody:   false, // Generally unsafe for production
		LogResponseBody:  false, // Generally unsafe for production
		LogHeaders:       true,
		EnableTracing:    true,
		SkipPaths:        []string{"/health", "/metrics", "/ping"},
		MaxBodySize:      1024, // 1KB max for logging
	}
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(baseLogger *slog.Logger, config LoggingConfig) *LoggingMiddleware {
	logConfig := monitoring.DefaultLogConfig()
	structuredLogger := monitoring.NewStructuredLogger(baseLogger, logConfig)
	
	// Create skip paths map for faster lookup
	skipPaths := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}
	
	return &LoggingMiddleware{
		logger:    structuredLogger,
		config:    config,
		skipPaths: skipPaths,
	}
}

// Middleware returns the HTTP middleware function
func (m *LoggingMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip logging for configured paths
			if m.skipPaths[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}
			
			start := time.Now()
			
			// Generate request ID and trace ID
			requestID := uuid.New().String()
			traceID := uuid.New().String()
			
			// Add IDs to context
			ctx := r.Context()
			if m.config.EnableTracing {
				ctx = context.WithValue(ctx, "request_id", requestID)
				ctx = context.WithValue(ctx, "trace_id", traceID)
				r = r.WithContext(ctx)
			}
			
			// Create response writer wrapper to capture status code
			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // Default status
			}
			
			// Process request
			next.ServeHTTP(wrapper, r)
			
			duration := time.Since(start)
			
			// Create log context
			logCtx := monitoring.LogContext{
				TraceID:    traceID,
				RequestID:  requestID,
				Component:  "http_middleware",
				Operation:  "request_processing",
			}
			
			// Extract API key for logging (masked)
			if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
				logCtx.APIKey = apiKey
			}
			
			// Log the request
			m.logger.LogHTTPRequest(r, wrapper.statusCode, duration, logCtx)
			
			// Log additional details for errors
			if wrapper.statusCode >= 400 {
				m.logErrorDetails(r, wrapper.statusCode, duration, logCtx)
			}
		})
	}
}

// logErrorDetails logs additional details for error responses
func (m *LoggingMiddleware) logErrorDetails(r *http.Request, statusCode int, duration time.Duration, logCtx monitoring.LogContext) {
	logCtx.ErrorType = m.determineErrorType(statusCode)
	
	details := map[string]any{
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
		"path":        r.URL.Path,
		"query":       r.URL.RawQuery,
		"referer":     r.Header.Get("Referer"),
	}
	
	severity := "medium"
	if statusCode >= 500 {
		severity = "high"
	}
	
	m.logger.LogSecurityEvent(
		r.Context(),
		"http_error",
		severity,
		details,
		logCtx,
	)
}

// determineErrorType determines error type based on HTTP status code
func (m *LoggingMiddleware) determineErrorType(statusCode int) monitoring.ErrorType {
	switch {
	case statusCode == 401:
		return monitoring.ErrorTypeAuthentication
	case statusCode == 403:
		return monitoring.ErrorTypeAuthorization
	case statusCode == 429:
		return monitoring.ErrorTypeRateLimit
	case statusCode >= 400 && statusCode < 500:
		return monitoring.ErrorTypeValidation
	case statusCode >= 500:
		return monitoring.ErrorTypeInternal
	default:
		return monitoring.ErrorTypeInternal
	}
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write ensures WriteHeader is called
func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

// RequestLogger creates a logger with request context
func RequestLogger(ctx context.Context, baseLogger *monitoring.StructuredLogger) *monitoring.StructuredLogger {
	// This could be extended to create a logger with request-specific context
	// For now, just return the base logger
	return baseLogger
}

// LogMiddlewareError logs middleware-specific errors
func LogMiddlewareError(logger *monitoring.StructuredLogger, ctx context.Context, middleware string, err error, details map[string]any) {
	logCtx := monitoring.WithContext(ctx)
	logCtx.Component = middleware
	logCtx.ErrorType = monitoring.ErrorTypeInternal
	logCtx.Metadata = details
	
	logger.LogError(ctx, err, "Middleware error occurred", logCtx)
}

// LogMiddlewareWarning logs middleware-specific warnings
func LogMiddlewareWarning(logger *monitoring.StructuredLogger, ctx context.Context, middleware string, message string, details map[string]any) {
	logCtx := monitoring.WithContext(ctx)
	logCtx.Component = middleware
	logCtx.Metadata = details
	
	logger.LogWarning(ctx, message, logCtx)
}