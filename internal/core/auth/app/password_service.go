package app

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password hashing and validation
type PasswordService struct {
	bcryptCost      int
	minLength       int
	requireUpper    bool
	requireLower    bool
	requireNumber   bool
	requireSpecial  bool
	commonPasswords map[string]bool
}

// PasswordConfig holds password service configuration
type PasswordConfig struct {
	BcryptCost     int  // Default: 12
	MinLength      int  // Default: 8
	RequireUpper   bool // Default: true
	RequireLower   bool // Default: true
	RequireNumber  bool // Default: true
	RequireSpecial bool // Default: true
}

// NewPasswordService creates a new password service
func NewPasswordService(config *PasswordConfig) *PasswordService {
	if config == nil {
		config = &PasswordConfig{
			BcryptCost:     12,
			MinLength:      8,
			RequireUpper:   true,
			RequireLower:   true,
			RequireNumber:  true,
			RequireSpecial: true,
		}
	}

	// Set defaults if not provided
	if config.BcryptCost == 0 {
		config.BcryptCost = 12
	}
	if config.MinLength == 0 {
		config.MinLength = 8
	}

	return &PasswordService{
		bcryptCost:      config.BcryptCost,
		minLength:       config.MinLength,
		requireUpper:    config.RequireUpper,
		requireLower:    config.RequireLower,
		requireNumber:   config.RequireNumber,
		requireSpecial:  config.RequireSpecial,
		commonPasswords: loadCommonPasswords(),
	}
}

// HashPassword hashes a password using bcrypt
func (s *PasswordService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func (s *PasswordService) VerifyPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return fmt.Errorf("invalid password")
	}
	if err != nil {
		return fmt.Errorf("failed to verify password: %w", err)
	}
	return nil
}

// ValidatePasswordStrength validates password strength
func (s *PasswordService) ValidatePasswordStrength(password string) error {
	// Check minimum length
	if len(password) < s.minLength {
		return fmt.Errorf("password must be at least %d characters long", s.minLength)
	}

	// Check maximum length (prevent DoS)
	if len(password) > 128 {
		return fmt.Errorf("password must not exceed 128 characters")
	}

	// Check for uppercase
	if s.requireUpper && !hasUpperCase(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	// Check for lowercase
	if s.requireLower && !hasLowerCase(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	// Check for number
	if s.requireNumber && !hasNumber(password) {
		return fmt.Errorf("password must contain at least one number")
	}

	// Check for special character
	if s.requireSpecial && !hasSpecialChar(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	// Check against common passwords
	if s.isCommonPassword(password) {
		return fmt.Errorf("password is too common, please choose a stronger password")
	}

	return nil
}

// CalculatePasswordStrength returns password strength score (0-100)
func (s *PasswordService) CalculatePasswordStrength(password string) int {
	score := 0

	// Length score (up to 40 points)
	length := len(password)
	switch {
	case length >= 16:
		score += 40
	case length >= 12:
		score += 30
	case length >= 8:
		score += 20
	default:
		score += 10
	}

	// Complexity score (up to 40 points)
	complexity := 0
	if hasUpperCase(password) {
		complexity += 10
	}
	if hasLowerCase(password) {
		complexity += 10
	}
	if hasNumber(password) {
		complexity += 10
	}
	if hasSpecialChar(password) {
		complexity += 10
	}
	score += complexity

	// Diversity score (up to 20 points)
	uniqueChars := countUniqueChars(password)
	diversityScore := (uniqueChars * 20) / length
	if diversityScore > 20 {
		diversityScore = 20
	}
	score += diversityScore

	// Penalty for common passwords
	if s.isCommonPassword(password) {
		score = score / 2
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

// Helper functions

func hasUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func hasLowerCase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func hasNumber(s string) bool {
	for _, r := range s {
		if unicode.IsNumber(r) {
			return true
		}
	}
	return false
}

func hasSpecialChar(s string) bool {
	// Special characters: !@#$%^&*()_+-=[]{}|;:,.<>?
	specialChars := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`)
	return specialChars.MatchString(s)
}

func countUniqueChars(s string) int {
	chars := make(map[rune]bool)
	for _, r := range s {
		chars[r] = true
	}
	return len(chars)
}

func (s *PasswordService) isCommonPassword(password string) bool {
	// Check case-insensitive
	return s.commonPasswords[strings.ToLower(password)]
}

// loadCommonPasswords returns a set of common passwords to reject
func loadCommonPasswords() map[string]bool {
	// Top 100 most common passwords
	commonPasswords := []string{
		"password", "123456", "12345678", "1234", "qwerty", "12345",
		"dragon", "pussy", "baseball", "football", "letmein", "monkey",
		"696969", "abc123", "mustang", "michael", "shadow", "master",
		"jennifer", "111111", "2000", "jordan", "superman", "harley",
		"1234567", "fuckme", "hunter", "fuckyou", "trustno1", "ranger",
		"buster", "thomas", "tigger", "robert", "soccer", "fuck",
		"batman", "test", "pass", "killer", "hockey", "george",
		"charlie", "andrew", "michelle", "love", "sunshine", "jessica",
		"asshole", "6969", "pepper", "daniel", "access", "123456789",
		"654321", "joshua", "maggie", "starwars", "silver", "william",
		"dallas", "yankees", "123123", "ashley", "666666", "hello",
		"amanda", "orange", "biteme", "freedom", "computer", "sexy",
		"thunder", "nicole", "ginger", "heather", "hammer", "summer",
		"corvette", "taylor", "fucker", "austin", "1111", "merlin",
		"matthew", "121212", "golfer", "cheese", "princess", "martin",
		"chelsea", "patrick", "richard", "diamond", "yellow", "bigdog",
		"secret", "admin", "administrator", "root", "welcome",
	}

	passwordMap := make(map[string]bool)
	for _, pwd := range commonPasswords {
		passwordMap[pwd] = true
	}

	return passwordMap
}

// PasswordStrengthResult holds password strength analysis
type PasswordStrengthResult struct {
	Score       int      `json:"score"`        // 0-100
	Strength    string   `json:"strength"`     // weak, fair, good, strong
	Suggestions []string `json:"suggestions"`  // Improvement suggestions
	Valid       bool     `json:"valid"`        // Meets minimum requirements
}

// AnalyzePassword provides detailed password analysis
func (s *PasswordService) AnalyzePassword(password string) *PasswordStrengthResult {
	result := &PasswordStrengthResult{
		Score:       s.CalculatePasswordStrength(password),
		Suggestions: []string{},
	}

	// Determine strength label
	switch {
	case result.Score >= 80:
		result.Strength = "strong"
	case result.Score >= 60:
		result.Strength = "good"
	case result.Score >= 40:
		result.Strength = "fair"
	default:
		result.Strength = "weak"
	}

	// Check if meets minimum requirements
	err := s.ValidatePasswordStrength(password)
	result.Valid = (err == nil)

	// Generate suggestions
	if len(password) < 12 {
		result.Suggestions = append(result.Suggestions, "Use at least 12 characters for better security")
	}
	if !hasUpperCase(password) {
		result.Suggestions = append(result.Suggestions, "Add uppercase letters")
	}
	if !hasLowerCase(password) {
		result.Suggestions = append(result.Suggestions, "Add lowercase letters")
	}
	if !hasNumber(password) {
		result.Suggestions = append(result.Suggestions, "Add numbers")
	}
	if !hasSpecialChar(password) {
		result.Suggestions = append(result.Suggestions, "Add special characters (!@#$%^&*)")
	}
	if s.isCommonPassword(password) {
		result.Suggestions = append(result.Suggestions, "This is a common password, choose something unique")
	}
	if countUniqueChars(password) < len(password)/2 {
		result.Suggestions = append(result.Suggestions, "Use more diverse characters")
	}

	return result
}
