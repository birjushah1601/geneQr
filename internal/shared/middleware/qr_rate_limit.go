package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// QRRateLimiter tracks ticket creation attempts per QR code
type QRRateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int           // Max requests
	window   time.Duration // Time window
	logger   *slog.Logger
}

// NewQRRateLimiter creates a new QR-based rate limiter
// Example: NewQRRateLimiter(5, 1*time.Hour) = 5 requests per hour
func NewQRRateLimiter(limit int, window time.Duration, logger *slog.Logger) *QRRateLimiter {
	limiter := &QRRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		logger:   logger,
	}

	// Cleanup old entries every 10 minutes
	go limiter.cleanup()

	return limiter
}

// Allow checks if a request for the given QR code is allowed
func (rl *QRRateLimiter) Allow(qrCode string) bool {
	if qrCode == "" {
		return true // Allow if no QR code (will be caught by validation)
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Get existing requests for this QR code
	requests, exists := rl.requests[qrCode]
	if !exists {
		requests = []time.Time{}
	}

	// Filter out old requests
	validRequests := []time.Time{}
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Check if limit exceeded
	if len(validRequests) >= rl.limit {
		rl.logger.Warn("QR rate limit exceeded",
			slog.String("qr_code", qrCode),
			slog.Int("requests_in_window", len(validRequests)),
			slog.Int("limit", rl.limit))
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[qrCode] = validRequests

	return true
}

// Middleware creates HTTP middleware for QR-based rate limiting
func (rl *QRRateLimiter) Middleware(qrCodeExtractor func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			qrCode := qrCodeExtractor(r)

			if qrCode != "" && !rl.Allow(qrCode) {
				rl.logger.Warn("QR rate limit exceeded - request denied",
					slog.String("qr_code", qrCode),
					slog.String("ip", r.RemoteAddr),
					slog.String("path", r.URL.Path))
				
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(rl.window.Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(fmt.Sprintf(`{
					"error": "Rate limit exceeded",
					"message": "Too many requests for this equipment. Please try again later.",
					"retry_after_seconds": %d
				}`, int(rl.window.Seconds()))))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// cleanup removes old entries periodically
func (rl *QRRateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)

		// Remove QR codes with no recent requests
		for qrCode, requests := range rl.requests {
			validRequests := []time.Time{}
			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}

			if len(validRequests) == 0 {
				delete(rl.requests, qrCode)
			} else {
				rl.requests[qrCode] = validRequests
			}
		}

		rl.logger.Debug("QR rate limiter cleanup completed",
			slog.Int("active_qr_codes", len(rl.requests)))
		rl.mu.Unlock()
	}
}

// GetStats returns current rate limiter statistics
func (rl *QRRateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"active_qr_codes": len(rl.requests),
		"limit":           rl.limit,
		"window_minutes":  int(rl.window.Minutes()),
	}
}
