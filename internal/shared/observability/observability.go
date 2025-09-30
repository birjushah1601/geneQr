package observability

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config holds observability configuration
type Config struct {
	LogLevel       string
	TracingEnabled bool
	TracingURL     string
	SamplingRate   float64
	MetricsEnabled bool
	ServiceName    string
	Environment    string
}

// Tracer is a simple wrapper around OpenTelemetry tracer
type Tracer interface {
	Start(ctx context.Context, name string) (context.Context, interface{})
}

// noopTracer is used when tracing is disabled or not compiled in.
type noopTracer struct{}

func (t noopTracer) Start(ctx context.Context, name string) (context.Context, interface{}) {
	return ctx, struct{}{}
}

var (
	// Global variables for observability components
	globalLogger   *slog.Logger
	globalTracer   Tracer

	// HTTP request metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// ensures collectors are only registered once
	registerMetricsOnce sync.Once

	// custom registry to avoid duplicate global registrations
	metricsRegistry *prometheus.Registry
)

// Setup initializes all observability components
func Setup(_ context.Context, cfg Config) (*slog.Logger, Tracer, interface{}, error) {
	// Initialize logger
	logger := setupLogger(cfg.LogLevel)
	globalLogger = logger

	// Set up basic Prometheus metrics (always enabled)
	registerMetricsOnce.Do(func() {
		metricsRegistry = prometheus.NewRegistry()
		metricsRegistry.MustRegister(collectors.NewGoCollector())
		metricsRegistry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		metricsRegistry.MustRegister(httpRequestsTotal)
		metricsRegistry.MustRegister(httpRequestDuration)
	})

	globalTracer = noopTracer{}

	logger.Info("Observability setup complete",
		slog.Bool("tracing_enabled", cfg.TracingEnabled),
		slog.Bool("metrics_enabled", true),
		slog.String("log_level", cfg.LogLevel))

	return logger, globalTracer, nil, nil
}

// setupLogger initializes the structured logger
func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	return logger
}

// LoggingMiddleware creates a middleware that logs HTTP requests
func LoggingMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Create a response writer wrapper to capture the status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			
			// Extract request ID from context or create a new one
			requestID := middleware.GetReqID(r.Context())
			
			// Create a child context with the logger containing request information
			ctx := r.Context()
			requestLogger := logger.With(
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)
			
			// Store the logger in the context
			ctx = context.WithValue(ctx, "logger", requestLogger)
			
			// Call the next handler with the enriched context
			next.ServeHTTP(ww, r.WithContext(ctx))
			
			// Calculate request duration
			duration := time.Since(start)
			
			// Log the request completion
			requestLogger.Info("HTTP request completed",
				slog.Int("status", ww.Status()),
				slog.Duration("duration", duration),
				slog.Int("bytes", ww.BytesWritten()),
			)
			
			// Record metrics if enabled
			statusCode := fmt.Sprintf("%d", ww.Status())
			httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, statusCode).Inc()
			httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, statusCode).Observe(duration.Seconds())
		})
	}
}

// SetupRouter configures a Chi router with standard middleware instead of OpenTelemetry
func SetupRouter(logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	// Use standard Chi middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(LoggingMiddleware(logger)) // Use our custom logging middleware instead of otelchi
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Add metrics endpoint
	r.Get("/metrics", MetricsHandler())
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r
}

// MetricsHandler returns an HTTP handler for the /metrics endpoint
func MetricsHandler() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}

// Shutdown gracefully shuts down observability components
func Shutdown(ctx context.Context) {
	// nothing to clean up in the simplified implementation
}

// GetLogger retrieves the logger from the context or returns the global logger
func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		return logger
	}
	return globalLogger
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, name string) (context.Context, interface{}) {
	return globalTracer.Start(ctx, name)
}

// AddSpanAttributes adds attributes to the current span
func AddSpanAttributes(_ context.Context, _ ...interface{}) {
	// no-op
}

// RecordError records an error in the current span
func RecordError(_ context.Context, _ error, _ ...interface{}) {
	// no-op
}

// ExtractTraceInfo extracts trace and span IDs from context for logging
func ExtractTraceInfo(ctx context.Context) (string, string) {
	// tracing disabled – always return empty IDs
	return "", ""
}

// WithSpan wraps a function with a new span
func WithSpan(ctx context.Context, name string, fn func(context.Context) error) error {
	// execute function directly (no tracing)
	ctx, _ = StartSpan(ctx, name)
	return fn(ctx)
}
