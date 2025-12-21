package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		
		// Prevent MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		
		// Enable XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		// Enforce HTTPS (in production)
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		// Control referrer information
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy (adjust as needed)
		w.Header().Set("Content-Security-Policy", 
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data: https:; "+
			"font-src 'self' data:; "+
			"connect-src 'self'")
		
		// Permissions Policy (formerly Feature-Policy)
		w.Header().Set("Permissions-Policy", 
			"camera=(), microphone=(), geolocation=()")
		
		next.ServeHTTP(w, r)
	})
}

// RateLimitByIP creates a rate limiter based on IP address
// limit: number of requests allowed
// window: time window for the limit
func RateLimitByIP(limit int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.LimitByIP(limit, window)
}

// RateLimitByUser creates a rate limiter based on authenticated user
// This should be applied after authentication middleware
func RateLimitByUser(limit int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.Limit(
		limit,
		window,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			// Try to get user ID from context (set by auth middleware)
			userID := r.Context().Value("user_id")
			if userID != nil {
				return userID.(string), nil
			}
			// Fallback to IP if not authenticated
			return httprate.KeyByIP(r)
		}),
	)
}

// RequestID adds a unique request ID to each request
func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

// RealIP sets the real IP address from proxy headers
func RealIP(next http.Handler) http.Handler {
	return middleware.RealIP(next)
}

// Recoverer recovers from panics and returns 500
func Recoverer(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

// Compress compresses responses
func Compress(level int) func(http.Handler) http.Handler {
	return middleware.Compress(level)
}
