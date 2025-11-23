package monitoring

import (
	"context"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity of log entries
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// ErrorType categorizes different types of errors
type ErrorType string

const (
	ErrorTypeValidation    ErrorType = "validation_error"
	ErrorTypeDatabase      ErrorType = "database_error"
	ErrorTypeNetwork       ErrorType = "network_error"
	ErrorTypeAuthentication ErrorType = "auth_error"
	ErrorTypeAuthorization ErrorType = "authorization_error"
	ErrorTypeRateLimit     ErrorType = "rate_limit_error"
	ErrorTypeFileSystem    ErrorType = "filesystem_error"
	ErrorTypeAI            ErrorType = "ai_service_error"
	ErrorTypeInternal      ErrorType = "internal_error"
	ErrorTypeExternal      ErrorType = "external_error"
)

// LogContext holds contextual information for logging
type LogContext struct {
	TraceID       string            `json:"trace_id,omitempty"`
	UserID        string            `json:"user_id,omitempty"`
	APIKey        string            `json:"api_key,omitempty"`
	RequestID     string            `json:"request_id,omitempty"`
	Operation     string            `json:"operation,omitempty"`
	Component     string            `json:"component,omitempty"`
	Method        string            `json:"method,omitempty"`
	URL           string            `json:"url,omitempty"`
	UserAgent     string            `json:"user_agent,omitempty"`
	RemoteAddr    string            `json:"remote_addr,omitempty"`
	Duration      time.Duration     `json:"duration,omitempty"`
	StatusCode    int               `json:"status_code,omitempty"`
	ErrorType     ErrorType         `json:"error_type,omitempty"`
	ErrorCode     string            `json:"error_code,omitempty"`
	Metadata      map[string]any    `json:"metadata,omitempty"`
}

// StructuredLogger provides enhanced logging capabilities with structured data
type StructuredLogger struct {
	logger *slog.Logger
	config LogConfig
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level           LogLevel          `json:"level"`
	Format          string            `json:"format"` // "json" or "text"
	EnableStackTrace bool             `json:"enable_stack_trace"`
	EnableCaller     bool             `json:"enable_caller"`
	TimestampFormat  string           `json:"timestamp_format"`
	FieldMapping     map[string]string `json:"field_mapping"`
	SensitiveFields  []string         `json:"sensitive_fields"`
}

// DefaultLogConfig returns default logging configuration
func DefaultLogConfig() LogConfig {
	return LogConfig{
		Level:           LogLevelInfo,
		Format:          "json",
		EnableStackTrace: true,
		EnableCaller:     true,
		TimestampFormat:  time.RFC3339,
		FieldMapping: map[string]string{
			"timestamp": "timestamp",
			"level":     "level",
			"message":   "message",
			"component": "component",
			"trace_id":  "trace_id",
		},
		SensitiveFields: []string{"password", "token", "api_key", "secret"},
	}
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger(baseLogger *slog.Logger, config LogConfig) *StructuredLogger {
	return &StructuredLogger{
		logger: baseLogger,
		config: config,
	}
}

// LogError logs an error with structured context
func (l *StructuredLogger) LogError(ctx context.Context, err error, message string, logCtx LogContext) {
	attrs := l.buildLogAttributes(logCtx)
	
	// Add error information
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	if logCtx.ErrorType != "" {
		attrs = append(attrs, slog.String("error_type", string(logCtx.ErrorType)))
	}
	
	// Add stack trace if enabled
	if l.config.EnableStackTrace {
		attrs = append(attrs, slog.String("stack_trace", l.getStackTrace()))
	}
	
	// Convert to []any for slog
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	l.logger.ErrorContext(ctx, message, args...)
}

// LogWarning logs a warning with structured context
func (l *StructuredLogger) LogWarning(ctx context.Context, message string, logCtx LogContext) {
	attrs := l.buildLogAttributes(logCtx)
	// Convert to []any for slog
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	l.logger.WarnContext(ctx, message, args...)
}

// LogInfo logs an info message with structured context
func (l *StructuredLogger) LogInfo(ctx context.Context, message string, logCtx LogContext) {
	attrs := l.buildLogAttributes(logCtx)
	// Convert to []any for slog
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	l.logger.InfoContext(ctx, message, args...)
}

// LogDebug logs a debug message with structured context
func (l *StructuredLogger) LogDebug(ctx context.Context, message string, logCtx LogContext) {
	attrs := l.buildLogAttributes(logCtx)
	// Convert to []any for slog
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	l.logger.DebugContext(ctx, message, args...)
}

// LogHTTPRequest logs HTTP request details
func (l *StructuredLogger) LogHTTPRequest(r *http.Request, statusCode int, duration time.Duration, logCtx LogContext) {
	// Sanitize sensitive headers
	headers := l.sanitizeHeaders(r.Header)
	
	logCtx.Method = r.Method
	logCtx.URL = r.URL.String()
	logCtx.UserAgent = r.UserAgent()
	logCtx.RemoteAddr = r.RemoteAddr
	logCtx.Duration = duration
	logCtx.StatusCode = statusCode
	
	if logCtx.Metadata == nil {
		logCtx.Metadata = make(map[string]any)
	}
	logCtx.Metadata["headers"] = headers
	logCtx.Metadata["content_length"] = r.ContentLength
	logCtx.Metadata["protocol"] = r.Proto
	
	message := "HTTP request processed"
	if statusCode >= 400 {
		logCtx.ErrorType = l.determineErrorType(statusCode)
		l.LogWarning(r.Context(), message, logCtx)
	} else {
		l.LogInfo(r.Context(), message, logCtx)
	}
}

// LogDatabaseOperation logs database operation details
func (l *StructuredLogger) LogDatabaseOperation(ctx context.Context, operation string, table string, duration time.Duration, err error, logCtx LogContext) {
	logCtx.Operation = operation
	logCtx.Component = "database"
	logCtx.Duration = duration
	
	if logCtx.Metadata == nil {
		logCtx.Metadata = make(map[string]any)
	}
	logCtx.Metadata["table"] = table
	
	message := "Database operation completed"
	
	if err != nil {
		logCtx.ErrorType = ErrorTypeDatabase
		l.LogError(ctx, err, message, logCtx)
	} else {
		l.LogDebug(ctx, message, logCtx)
	}
}

// LogAIOperation logs AI service operation details
func (l *StructuredLogger) LogAIOperation(ctx context.Context, provider string, model string, tokensUsed int, cost float64, duration time.Duration, err error, logCtx LogContext) {
	logCtx.Operation = "ai_analysis"
	logCtx.Component = "ai_service"
	logCtx.Duration = duration
	
	if logCtx.Metadata == nil {
		logCtx.Metadata = make(map[string]any)
	}
	logCtx.Metadata["provider"] = provider
	logCtx.Metadata["model"] = model
	logCtx.Metadata["tokens_used"] = tokensUsed
	logCtx.Metadata["cost_usd"] = cost
	
	message := "AI operation completed"
	
	if err != nil {
		logCtx.ErrorType = ErrorTypeAI
		l.LogError(ctx, err, message, logCtx)
	} else {
		l.LogInfo(ctx, message, logCtx)
	}
}

// LogFileOperation logs file system operation details
func (l *StructuredLogger) LogFileOperation(ctx context.Context, operation string, filePath string, fileSize int64, duration time.Duration, err error, logCtx LogContext) {
	logCtx.Operation = operation
	logCtx.Component = "filesystem"
	logCtx.Duration = duration
	
	if logCtx.Metadata == nil {
		logCtx.Metadata = make(map[string]any)
	}
	logCtx.Metadata["file_path"] = filePath
	logCtx.Metadata["file_size"] = fileSize
	
	message := "File operation completed"
	
	if err != nil {
		logCtx.ErrorType = ErrorTypeFileSystem
		l.LogError(ctx, err, message, logCtx)
	} else {
		l.LogDebug(ctx, message, logCtx)
	}
}

// LogSecurityEvent logs security-related events
func (l *StructuredLogger) LogSecurityEvent(ctx context.Context, eventType string, severity string, details map[string]any, logCtx LogContext) {
	logCtx.Component = "security"
	logCtx.ErrorType = ErrorTypeAuthentication
	
	if logCtx.Metadata == nil {
		logCtx.Metadata = make(map[string]any)
	}
	logCtx.Metadata["event_type"] = eventType
	logCtx.Metadata["severity"] = severity
	for k, v := range details {
		logCtx.Metadata[k] = v
	}
	
	message := "Security event detected"
	
	switch severity {
	case "critical", "high":
		l.LogError(ctx, nil, message, logCtx)
	case "medium":
		l.LogWarning(ctx, message, logCtx)
	default:
		l.LogInfo(ctx, message, logCtx)
	}
}

// buildLogAttributes converts LogContext to slog attributes
func (l *StructuredLogger) buildLogAttributes(logCtx LogContext) []slog.Attr {
	var attrs []slog.Attr
	
	if logCtx.TraceID != "" {
		attrs = append(attrs, slog.String("trace_id", logCtx.TraceID))
	}
	if logCtx.UserID != "" {
		attrs = append(attrs, slog.String("user_id", logCtx.UserID))
	}
	if logCtx.APIKey != "" {
		attrs = append(attrs, slog.String("api_key", l.maskSensitiveData(logCtx.APIKey)))
	}
	if logCtx.RequestID != "" {
		attrs = append(attrs, slog.String("request_id", logCtx.RequestID))
	}
	if logCtx.Operation != "" {
		attrs = append(attrs, slog.String("operation", logCtx.Operation))
	}
	if logCtx.Component != "" {
		attrs = append(attrs, slog.String("component", logCtx.Component))
	}
	if logCtx.Method != "" {
		attrs = append(attrs, slog.String("method", logCtx.Method))
	}
	if logCtx.URL != "" {
		attrs = append(attrs, slog.String("url", logCtx.URL))
	}
	if logCtx.UserAgent != "" {
		attrs = append(attrs, slog.String("user_agent", logCtx.UserAgent))
	}
	if logCtx.RemoteAddr != "" {
		attrs = append(attrs, slog.String("remote_addr", logCtx.RemoteAddr))
	}
	if logCtx.Duration > 0 {
		attrs = append(attrs, slog.Duration("duration", logCtx.Duration))
	}
	if logCtx.StatusCode > 0 {
		attrs = append(attrs, slog.Int("status_code", logCtx.StatusCode))
	}
	if logCtx.ErrorType != "" {
		attrs = append(attrs, slog.String("error_type", string(logCtx.ErrorType)))
	}
	if logCtx.ErrorCode != "" {
		attrs = append(attrs, slog.String("error_code", logCtx.ErrorCode))
	}
	
	// Add metadata as nested attributes
	if logCtx.Metadata != nil {
		for key, value := range logCtx.Metadata {
			attrs = append(attrs, slog.Any(key, value))
		}
	}
	
	return attrs
}

// sanitizeHeaders removes or masks sensitive headers
func (l *StructuredLogger) sanitizeHeaders(headers http.Header) map[string]string {
	sanitized := make(map[string]string)
	
	for name, values := range headers {
		lowerName := strings.ToLower(name)
		
		// Skip or mask sensitive headers
		if l.isSensitiveField(lowerName) {
			sanitized[name] = "[MASKED]"
		} else {
			sanitized[name] = strings.Join(values, ", ")
		}
	}
	
	return sanitized
}

// maskSensitiveData masks sensitive data for logging
func (l *StructuredLogger) maskSensitiveData(data string) string {
	if len(data) <= 8 {
		return strings.Repeat("*", len(data))
	}
	return data[:4] + strings.Repeat("*", len(data)-8) + data[len(data)-4:]
}

// isSensitiveField checks if a field contains sensitive data
func (l *StructuredLogger) isSensitiveField(fieldName string) bool {
	for _, sensitive := range l.config.SensitiveFields {
		if strings.Contains(fieldName, sensitive) {
			return true
		}
	}
	return false
}

// getStackTrace returns the current stack trace
func (l *StructuredLogger) getStackTrace() string {
	buf := make([]byte, 1024*10)
	runtime.Stack(buf, false)
	return string(buf)
}

// determineErrorType determines error type based on HTTP status code
func (l *StructuredLogger) determineErrorType(statusCode int) ErrorType {
	switch {
	case statusCode == 401:
		return ErrorTypeAuthentication
	case statusCode == 403:
		return ErrorTypeAuthorization
	case statusCode == 429:
		return ErrorTypeRateLimit
	case statusCode >= 400 && statusCode < 500:
		return ErrorTypeValidation
	case statusCode >= 500:
		return ErrorTypeInternal
	default:
		return ErrorTypeInternal
	}
}

// WithContext adds trace ID and request context to LogContext
func WithContext(ctx context.Context) LogContext {
	logCtx := LogContext{}
	
	// Try to extract trace ID from context (if available)
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		logCtx.TraceID = traceID
	}
	
	// Try to extract request ID from context
	if requestID, ok := ctx.Value("request_id").(string); ok {
		logCtx.RequestID = requestID
	}
	
	// Try to extract user ID from context
	if userID, ok := ctx.Value("user_id").(string); ok {
		logCtx.UserID = userID
	}
	
	return logCtx
}