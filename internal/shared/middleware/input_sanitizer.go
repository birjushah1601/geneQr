package middleware

import (
	"encoding/json"
	"html"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

// InputSanitizer handles input sanitization for requests
type InputSanitizer struct {
	logger *slog.Logger
}

// NewInputSanitizer creates a new input sanitizer
func NewInputSanitizer(logger *slog.Logger) *InputSanitizer {
	return &InputSanitizer{
		logger: logger.With(slog.String("component", "input_sanitizer")),
	}
}

// SanitizeConfig holds sanitization configuration
type SanitizeConfig struct {
	// Size limits
	MaxDescriptionLength int // Max length for description fields
	MaxNameLength        int // Max length for name fields
	MaxPhoneLength       int // Max length for phone fields
	MaxBodySize          int // Max request body size in bytes

	// Sanitization options
	StripHTML           bool // Strip all HTML tags
	StripScripts        bool // Strip script tags and event handlers
	EscapeSpecialChars  bool // Escape special characters
	TrimWhitespace      bool // Trim leading/trailing whitespace
	AllowedHTMLTags     []string // Allowed HTML tags (if StripHTML is false)
}

// DefaultSanitizeConfig returns default configuration
func DefaultSanitizeConfig() *SanitizeConfig {
	return &SanitizeConfig{
		MaxDescriptionLength: 5000,
		MaxNameLength:        200,
		MaxPhoneLength:       50,
		MaxBodySize:          1024 * 1024, // 1 MB
		StripHTML:            true,
		StripScripts:         true,
		EscapeSpecialChars:   true,
		TrimWhitespace:       true,
		AllowedHTMLTags:      []string{}, // No HTML allowed by default
	}
}

// Middleware creates HTTP middleware for input sanitization
func (s *InputSanitizer) Middleware(config *SanitizeConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = DefaultSanitizeConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check content type
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				next.ServeHTTP(w, r)
				return
			}

			// Check body size limit
			if r.ContentLength > int64(config.MaxBodySize) {
				s.logger.Warn("Request body too large",
					slog.Int64("size", r.ContentLength),
					slog.Int("max", config.MaxBodySize))
				
				http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
				return
			}

			// Read body
			body, err := io.ReadAll(io.LimitReader(r.Body, int64(config.MaxBodySize)))
			if err != nil {
				s.logger.Error("Failed to read request body", slog.Any("error", err))
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			r.Body.Close()

			// Always log body info at INFO level for debugging
			s.logger.Info("📦 Input Sanitizer - Received request body",
				slog.Int("size", len(body)),
				slog.String("preview", string(body)[:min(len(body), 200)]),
				slog.String("path", r.URL.Path))

			// Parse JSON
			var data map[string]interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				s.logger.Error("❌ Failed to parse JSON", 
					slog.Any("error", err),
					slog.String("body", string(body)[:min(len(body), 500)]))
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			s.logger.Info("✅ JSON parsed successfully", slog.Int("fields", len(data)))

			// Sanitize the data
			sanitized := s.sanitizeMap(data, config)

			// Marshal back to JSON
			sanitizedBody, err := json.Marshal(sanitized)
			if err != nil {
				s.logger.Error("Failed to marshal sanitized data", slog.Any("error", err))
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			s.logger.Info("🔄 Creating new body reader",
				slog.Int("sanitized_size", len(sanitizedBody)),
				slog.String("sanitized_preview", string(sanitizedBody)[:min(len(sanitizedBody), 200)]))

			// Create new request with sanitized body
			r.Body = io.NopCloser(strings.NewReader(string(sanitizedBody)))
			r.ContentLength = int64(len(sanitizedBody))

			s.logger.Info("✅ Body reader created, calling next handler")

			next.ServeHTTP(w, r)
		})
	}
}

// sanitizeMap sanitizes all string values in a map
func (s *InputSanitizer) sanitizeMap(data map[string]interface{}, config *SanitizeConfig) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = s.sanitizeString(v, key, config)
		case map[string]interface{}:
			result[key] = s.sanitizeMap(v, config)
		case []interface{}:
			result[key] = s.sanitizeArray(v, config)
		default:
			result[key] = value
		}
	}

	return result
}

// sanitizeArray sanitizes all values in an array
func (s *InputSanitizer) sanitizeArray(data []interface{}, config *SanitizeConfig) []interface{} {
	result := make([]interface{}, len(data))

	for i, value := range data {
		switch v := value.(type) {
		case string:
			result[i] = s.sanitizeString(v, "", config)
		case map[string]interface{}:
			result[i] = s.sanitizeMap(v, config)
		case []interface{}:
			result[i] = s.sanitizeArray(v, config)
		default:
			result[i] = value
		}
	}

	return result
}

// sanitizeString sanitizes a single string value
func (s *InputSanitizer) sanitizeString(value string, fieldName string, config *SanitizeConfig) string {
	// Trim whitespace
	if config.TrimWhitespace {
		value = strings.TrimSpace(value)
	}

	// Apply field-specific length limits
	value = s.applyLengthLimit(value, fieldName, config)

	// Strip HTML tags
	if config.StripHTML {
		value = s.stripHTML(value, config.AllowedHTMLTags)
	}

	// Strip script tags and event handlers
	if config.StripScripts {
		value = s.stripScripts(value)
	}

	// Escape special characters
	if config.EscapeSpecialChars {
		value = html.EscapeString(value)
	}

	return value
}

// applyLengthLimit applies field-specific length limits
func (s *InputSanitizer) applyLengthLimit(value string, fieldName string, config *SanitizeConfig) string {
	var maxLength int

	// Determine max length based on field name
	fieldLower := strings.ToLower(fieldName)
	switch {
	case strings.Contains(fieldLower, "description") || strings.Contains(fieldLower, "comment") || strings.Contains(fieldLower, "notes"):
		maxLength = config.MaxDescriptionLength
	case strings.Contains(fieldLower, "name"):
		maxLength = config.MaxNameLength
	case strings.Contains(fieldLower, "phone"):
		maxLength = config.MaxPhoneLength
	default:
		maxLength = config.MaxDescriptionLength // Default to description length
	}

	// Truncate if too long
	if len(value) > maxLength {
		s.logger.Warn("Truncating field due to length limit",
			slog.String("field", fieldName),
			slog.Int("original_length", len(value)),
			slog.Int("max_length", maxLength))
		return value[:maxLength]
	}

	return value
}

// stripHTML removes HTML tags from the string
func (s *InputSanitizer) stripHTML(value string, allowedTags []string) string {
	// If no tags are allowed, strip all HTML
	if len(allowedTags) == 0 {
		// Remove all HTML tags
		re := regexp.MustCompile(`<[^>]*>`)
		value = re.ReplaceAllString(value, "")
		return value
	}

	// TODO: Implement whitelist-based HTML tag filtering if needed
	// For now, strip all HTML for security
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(value, "")
}

// stripScripts removes script tags and event handlers
func (s *InputSanitizer) stripScripts(value string) string {
	// Remove script tags and content
	scriptRe := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	value = scriptRe.ReplaceAllString(value, "")

	// Remove inline event handlers (onclick, onerror, etc.)
	eventRe := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*["'][^"']*["']`)
	value = eventRe.ReplaceAllString(value, "")

	// Remove javascript: protocol
	jsProtocolRe := regexp.MustCompile(`(?i)javascript:\s*`)
	value = jsProtocolRe.ReplaceAllString(value, "")

	return value
}

// SanitizeString is a helper function to sanitize a single string
func SanitizeString(value string) string {
	sanitizer := &InputSanitizer{}
	config := DefaultSanitizeConfig()
	return sanitizer.sanitizeString(value, "", config)
}

// StripHTML is a helper function to strip HTML from a string
func StripHTML(value string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(value, "")
}

// StripScripts is a helper function to strip scripts from a string
func StripScripts(value string) string {
	sanitizer := &InputSanitizer{}
	return sanitizer.stripScripts(value)
}
