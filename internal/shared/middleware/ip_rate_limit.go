package middleware

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

// IPRateLimiter tracks request rates per IP address
type IPRateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int           // Max requests
	window   time.Duration // Time window
	logger   *slog.Logger
}

// NewIPRateLimiter creates a new IP-based rate limiter
// Example: NewIPRateLimiter(20, 1*time.Hour) = 20 requests per hour per IP
func NewIPRateLimiter(limit int, window time.Duration, logger *slog.Logger) *IPRateLimiter {
	limiter := &IPRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		logger:   logger,
	}

	// Cleanup old entries every 10 minutes
	go limiter.cleanup()

	return limiter
}

// Allow checks if a request from the given IP is allowed
func (rl *IPRateLimiter) Allow(ipAddress string) bool {
	if ipAddress == "" {
		return true // Allow if no IP (should not happen)
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Get existing requests for this IP
	requests, exists := rl.requests[ipAddress]
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
		rl.logger.Warn("IP rate limit exceeded",
			slog.String("ip_address", ipAddress),
			slog.Int("requests_in_window", len(validRequests)),
			slog.Int("limit", rl.limit))
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[ipAddress] = validRequests

	return true
}

// Middleware creates HTTP middleware for IP-based rate limiting
func (rl *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := extractIPAddress(r)

		if !rl.Allow(ipAddress) {
			rl.logger.Warn("IP rate limit exceeded - request denied",
				slog.String("ip", ipAddress),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method))

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(rl.window.Seconds())))
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
			w.Header().Set("X-RateLimit-Window", fmt.Sprintf("%d", int(rl.window.Seconds())))
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(fmt.Sprintf(`{
				"error": "Rate limit exceeded",
				"message": "Too many requests from your IP address. Please try again later.",
				"retry_after_seconds": %d,
				"limit": %d,
				"window_seconds": %d
			}`, int(rl.window.Seconds()), rl.limit, int(rl.window.Seconds()))))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// cleanup removes old entries periodically
func (rl *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)

		// Remove IPs with no recent requests
		for ip, requests := range rl.requests {
			validRequests := []time.Time{}
			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}

			if len(validRequests) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = validRequests
			}
		}

		rl.logger.Debug("IP rate limiter cleanup completed",
			slog.Int("active_ips", len(rl.requests)))
		rl.mu.Unlock()
	}
}

// GetStats returns current rate limiter statistics
func (rl *IPRateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"active_ips":     len(rl.requests),
		"limit":          rl.limit,
		"window_minutes": int(rl.window.Minutes()),
	}
}

// extractIPAddress gets the client IP from the request
func extractIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take first IP if multiple
		ips := splitIPs(xff)
		if len(ips) > 0 {
			return ips[0]
		}
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	return r.RemoteAddr
}

// splitIPs splits a comma-separated list of IPs
func splitIPs(s string) []string {
	var ips []string
	for _, ip := range splitAndTrim(s, ",") {
		if ip != "" {
			ips = append(ips, ip)
		}
	}
	return ips
}

// splitAndTrim splits a string and trims whitespace
func splitAndTrim(s string, sep string) []string {
	parts := []string{}
	for _, part := range split(s, sep) {
		trimmed := trim(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

// split is a simple string split helper
func split(s string, sep string) []string {
	result := []string{}
	current := ""
	
	for _, char := range s {
		if string(char) == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		result = append(result, current)
	}
	
	return result
}

// trim removes leading/trailing whitespace
func trim(s string) string {
	start := 0
	end := len(s)
	
	// Trim leading whitespace
	for start < end && isWhitespace(s[start]) {
		start++
	}
	
	// Trim trailing whitespace
	for end > start && isWhitespace(s[end-1]) {
		end--
	}
	
	return s[start:end]
}

// isWhitespace checks if a byte is whitespace
func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
