package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerHour    int           `json:"requests_per_hour"`
	CleanupInterval    time.Duration `json:"cleanup_interval"`
	BurstSize          int           `json:"burst_size"`
	BlockDuration      time.Duration `json:"block_duration"`
	WhitelistPaths     []string      `json:"whitelist_paths"`
	WhitelistIPs       []string      `json:"whitelist_ips"`
}

// ClientRateInfo tracks rate limiting for a client
type ClientRateInfo struct {
	Requests    []time.Time   `json:"requests"`
	BlockedUntil *time.Time   `json:"blocked_until,omitempty"`
	TotalRequests int64       `json:"total_requests"`
	LastRequest   time.Time   `json:"last_request"`
}

// RateLimitMiddleware handles rate limiting based on API keys and IP addresses
type RateLimitMiddleware struct {
	config      *RateLimitConfig
	logger      *slog.Logger
	clients     map[string]*ClientRateInfo
	mutex       sync.RWMutex
	stopCleanup chan struct{}
}

// NewRateLimitMiddleware creates a new rate limiting middleware
func NewRateLimitMiddleware(config *RateLimitConfig, logger *slog.Logger) *RateLimitMiddleware {
	if config.RequestsPerHour == 0 {
		config.RequestsPerHour = 1000 // Default: 1000 requests per hour
	}
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 5 * time.Minute // Cleanup every 5 minutes
	}
	if config.BurstSize == 0 {
		config.BurstSize = 100 // Allow burst of 100 requests
	}
	if config.BlockDuration == 0 {
		config.BlockDuration = 10 * time.Minute // Block for 10 minutes after limit exceeded
	}

	rl := &RateLimitMiddleware{
		config:      config,
		logger:      logger.With(slog.String("component", "rate_limit_middleware")),
		clients:     make(map[string]*ClientRateInfo),
		stopCleanup: make(chan struct{}),
	}

	// Start cleanup goroutine
	go rl.cleanupLoop()

	return rl
}

// DefaultRateLimitConfig returns default rate limiting configuration
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		RequestsPerHour:    1000,
		CleanupInterval:    5 * time.Minute,
		BurstSize:          100,
		BlockDuration:      10 * time.Minute,
		WhitelistPaths:     []string{"/health", "/metrics", "/ping"},
		WhitelistIPs:       []string{"127.0.0.1", "::1"},
	}
}

// Middleware returns the HTTP middleware function
func (rl *RateLimitMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if path is whitelisted
			if rl.isWhitelistedPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Get client identifier (API key or IP address)
			clientID := rl.getClientIdentifier(r)

			// Check if client IP is whitelisted
			if rl.isWhitelistedIP(r.RemoteAddr) {
				next.ServeHTTP(w, r)
				return
			}

			// Check rate limit
			allowed, rateInfo := rl.checkRateLimit(clientID)
			if !allowed {
				rl.writeRateLimitExceeded(w, rateInfo)
				return
			}

			// Record request
			rl.recordRequest(clientID)

			// Add rate limit headers
			rl.addRateLimitHeaders(w, clientID, rateInfo)

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIdentifier returns a unique identifier for the client
func (rl *RateLimitMiddleware) getClientIdentifier(r *http.Request) string {
	// Try to get API key first
	if authCtx, exists := GetAuthContext(r); exists {
		return "key:" + authCtx.APIKey[:min(len(authCtx.APIKey), 8)] // Use first 8 chars of API key
	}

	// Fall back to IP address
	return "ip:" + r.RemoteAddr
}

// checkRateLimit checks if a client has exceeded the rate limit
func (rl *RateLimitMiddleware) checkRateLimit(clientID string) (bool, *ClientRateInfo) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	info, exists := rl.clients[clientID]
	if !exists {
		info = &ClientRateInfo{
			Requests:      []time.Time{},
			TotalRequests: 0,
		}
		rl.clients[clientID] = info
	}

	// Check if client is currently blocked
	if info.BlockedUntil != nil && now.Before(*info.BlockedUntil) {
		return false, info
	}

	// Clear expired block
	if info.BlockedUntil != nil && now.After(*info.BlockedUntil) {
		info.BlockedUntil = nil
	}

	// Remove old requests (outside the 1-hour window)
	oneHourAgo := now.Add(-time.Hour)
	validRequests := []time.Time{}
	for _, reqTime := range info.Requests {
		if reqTime.After(oneHourAgo) {
			validRequests = append(validRequests, reqTime)
		}
	}
	info.Requests = validRequests

	// Get rate limit for this client
	limit := rl.getRateLimitForClient(clientID)

	// Check if client exceeds rate limit
	if len(info.Requests) >= limit {
		// Block the client
		blockUntil := now.Add(rl.config.BlockDuration)
		info.BlockedUntil = &blockUntil

		rl.logger.Warn("Rate limit exceeded, blocking client",
			slog.String("client_id", clientID),
			slog.Int("requests", len(info.Requests)),
			slog.Int("limit", limit),
			slog.Time("blocked_until", blockUntil))

		return false, info
	}

	return true, info
}

// recordRequest records a request for rate limiting
func (rl *RateLimitMiddleware) recordRequest(clientID string) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	info := rl.clients[clientID]
	
	info.Requests = append(info.Requests, now)
	info.TotalRequests++
	info.LastRequest = now

	// Log if client is approaching limit
	limit := rl.getRateLimitForClient(clientID)
	if len(info.Requests) > limit*8/10 { // 80% of limit
		rl.logger.Debug("Client approaching rate limit",
			slog.String("client_id", clientID),
			slog.Int("requests", len(info.Requests)),
			slog.Int("limit", limit),
			slog.Float64("percentage", float64(len(info.Requests))/float64(limit)*100))
	}
}

// getRateLimitForClient gets the rate limit for a specific client
func (rl *RateLimitMiddleware) getRateLimitForClient(clientID string) int {
	// If this is an API key, try to get its specific rate limit
	if authCtx, _ := getAuthContextFromClientID(clientID); authCtx != nil {
		if authCtx.KeyInfo.RateLimit > 0 {
			return authCtx.KeyInfo.RateLimit
		}
	}

	// Default rate limit
	return rl.config.RequestsPerHour
}

// getAuthContextFromClientID attempts to get auth context from client ID
// This is a simplified approach - in production, you might want to cache this
func getAuthContextFromClientID(clientID string) (*AuthContext, bool) {
	// This is a placeholder - in a real implementation, you'd look up the API key
	// from the client ID and get its rate limit from configuration
	return nil, false
}

// addRateLimitHeaders adds rate limiting headers to the response
func (rl *RateLimitMiddleware) addRateLimitHeaders(w http.ResponseWriter, clientID string, info *ClientRateInfo) {
	limit := rl.getRateLimitForClient(clientID)
	remaining := limit - len(info.Requests)
	if remaining < 0 {
		remaining = 0
	}

	// Add standard rate limit headers
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
	
	// Reset time (next hour)
	resetTime := time.Now().Add(time.Hour).Unix()
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))

	// If blocked, add retry-after header
	if info.BlockedUntil != nil {
		retryAfter := int(time.Until(*info.BlockedUntil).Seconds())
		w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
	}
}

// isWhitelistedPath checks if a path should bypass rate limiting
func (rl *RateLimitMiddleware) isWhitelistedPath(path string) bool {
	for _, whitelistPath := range rl.config.WhitelistPaths {
		if path == whitelistPath || (len(path) > len(whitelistPath) && 
			path[:len(whitelistPath)] == whitelistPath && path[len(whitelistPath)] == '/') {
			return true
		}
	}
	return false
}

// isWhitelistedIP checks if an IP should bypass rate limiting
func (rl *RateLimitMiddleware) isWhitelistedIP(addr string) bool {
	for _, whitelistIP := range rl.config.WhitelistIPs {
		if addr == whitelistIP {
			return true
		}
	}
	return false
}

// writeRateLimitExceeded writes a rate limit exceeded response
func (rl *RateLimitMiddleware) writeRateLimitExceeded(w http.ResponseWriter, info *ClientRateInfo) {
	w.Header().Set("Content-Type", "application/json")
	
	if info.BlockedUntil != nil {
		retryAfter := int(time.Until(*info.BlockedUntil).Seconds())
		w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
	}

	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte(`{"success": false, "error": "Rate limit exceeded", "details": {"code": "RATE_LIMIT_EXCEEDED", "message": "Too many requests, please try again later"}}`))
}

// cleanupLoop periodically cleans up old rate limit data
func (rl *RateLimitMiddleware) cleanupLoop() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.stopCleanup:
			return
		}
	}
}

// cleanup removes old rate limit data
func (rl *RateLimitMiddleware) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	oneHourAgo := now.Add(-time.Hour)
	cleanupThreshold := now.Add(-2 * time.Hour) // Remove clients not seen for 2 hours

	clientsToRemove := []string{}

	for clientID, info := range rl.clients {
		// Remove old requests
		validRequests := []time.Time{}
		for _, reqTime := range info.Requests {
			if reqTime.After(oneHourAgo) {
				validRequests = append(validRequests, reqTime)
			}
		}
		info.Requests = validRequests

		// Remove clients that haven't made requests recently and have no active blocks
		if len(info.Requests) == 0 && 
		   info.LastRequest.Before(cleanupThreshold) && 
		   (info.BlockedUntil == nil || now.After(*info.BlockedUntil)) {
			clientsToRemove = append(clientsToRemove, clientID)
		}
	}

	// Remove inactive clients
	for _, clientID := range clientsToRemove {
		delete(rl.clients, clientID)
	}

	if len(clientsToRemove) > 0 {
		rl.logger.Debug("Cleaned up rate limit data",
			slog.Int("clients_removed", len(clientsToRemove)),
			slog.Int("active_clients", len(rl.clients)))
	}
}

// GetStats returns rate limiting statistics
func (rl *RateLimitMiddleware) GetStats() map[string]interface{} {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	now := time.Now()
	stats := map[string]interface{}{
		"total_clients": len(rl.clients),
		"active_clients": 0,
		"blocked_clients": 0,
		"total_requests": int64(0),
	}

	oneHourAgo := now.Add(-time.Hour)
	for _, info := range rl.clients {
		stats["total_requests"] = stats["total_requests"].(int64) + info.TotalRequests
		
		// Count active clients (made request in last hour)
		if info.LastRequest.After(oneHourAgo) {
			stats["active_clients"] = stats["active_clients"].(int) + 1
		}

		// Count blocked clients
		if info.BlockedUntil != nil && now.Before(*info.BlockedUntil) {
			stats["blocked_clients"] = stats["blocked_clients"].(int) + 1
		}
	}

	return stats
}

// Stop stops the cleanup goroutine
func (rl *RateLimitMiddleware) Stop() {
	close(rl.stopCleanup)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}