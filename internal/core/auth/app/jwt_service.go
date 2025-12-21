package app

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService handles JWT token generation and validation
type JWTService struct {
	privateKey           *rsa.PrivateKey
	publicKey            *rsa.PublicKey
	accessTokenExpiry    time.Duration
	refreshTokenExpiry   time.Duration
	issuer               string
	refreshTokenRepo     RefreshTokenRepository
}

// RefreshTokenRepository interface for token persistence
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error
	Revoke(ctx context.Context, id uuid.UUID, reason string) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID, reason string) error
}

// RefreshToken model (simplified for JWT service)
type RefreshToken struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TokenHash  string
	ExpiresAt  time.Time
	DeviceInfo map[string]interface{}
	IPAddress  string
}

// JWTConfig holds JWT service configuration
type JWTConfig struct {
	PrivateKey         *rsa.PrivateKey
	PublicKey          *rsa.PublicKey
	AccessTokenExpiry  time.Duration // Default: 15 minutes
	RefreshTokenExpiry time.Duration // Default: 7 days
	Issuer             string
}

// Claims represents JWT claims
type Claims struct {
	UserID         string                 `json:"sub"`
	Email          string                 `json:"email,omitempty"`
	Name           string                 `json:"name,omitempty"`
	OrganizationID string                 `json:"org_id,omitempty"`
	Role           string                 `json:"role,omitempty"`
	Permissions    []string               `json:"permissions,omitempty"`
	TokenType      string                 `json:"type"` // "access" or "refresh"
	TokenID        string                 `json:"jti,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service
func NewJWTService(config *JWTConfig, refreshTokenRepo RefreshTokenRepository) *JWTService {
	if config.AccessTokenExpiry == 0 {
		config.AccessTokenExpiry = 15 * time.Minute
	}
	if config.RefreshTokenExpiry == 0 {
		config.RefreshTokenExpiry = 7 * 24 * time.Hour
	}
	if config.Issuer == "" {
		config.Issuer = "aby-med-platform"
	}

	return &JWTService{
		privateKey:           config.PrivateKey,
		publicKey:            config.PublicKey,
		accessTokenExpiry:    config.AccessTokenExpiry,
		refreshTokenExpiry:   config.RefreshTokenExpiry,
		issuer:               config.Issuer,
		refreshTokenRepo:     refreshTokenRepo,
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (s *JWTService) GenerateTokenPair(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
	now := time.Now()
	
	// Generate access token
	accessClaims := &Claims{
		UserID:         req.UserID.String(),
		Email:          req.Email,
		Name:           req.Name,
		OrganizationID: req.OrganizationID,
		Role:           req.Role,
		Permissions:    req.Permissions,
		TokenType:      "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   req.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExpiry)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshTokenID := uuid.New()
	refreshClaims := &Claims{
		UserID:    req.UserID.String(),
		TokenType: "refresh",
		TokenID:   refreshTokenID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   req.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenExpiry)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	// Hash and store refresh token
	tokenHash := s.hashToken(refreshTokenString)
	refreshTokenRecord := &RefreshToken{
		ID:         refreshTokenID,
		UserID:     req.UserID,
		TokenHash:  tokenHash,
		ExpiresAt:  now.Add(s.refreshTokenExpiry),
		DeviceInfo: req.DeviceInfo,
		IPAddress:  req.IPAddress,
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshTokenRecord); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessTokenExpiry.Seconds()),
	}, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token using refresh token
func (s *JWTService) RefreshAccessToken(ctx context.Context, refreshTokenString string) (*TokenResponse, error) {
	// Validate refresh token
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("not a refresh token")
	}

	// Check if token exists and is not revoked
	tokenHash := s.hashToken(refreshTokenString)
	storedToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or revoked: %w", err)
	}

	// Check expiry
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Update last used
	s.refreshTokenRepo.UpdateLastUsed(ctx, storedToken.ID)

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	// Generate new token pair (token rotation)
	newTokens, err := s.GenerateTokenPair(ctx, &TokenRequest{
		UserID:     userID,
		Email:      claims.Email,
		Name:       claims.Name,
		DeviceInfo: storedToken.DeviceInfo,
		IPAddress:  storedToken.IPAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// Revoke old refresh token (token rotation security)
	s.refreshTokenRepo.Revoke(ctx, storedToken.ID, "rotated")

	return newTokens, nil
}

// RevokeToken revokes a refresh token
func (s *JWTService) RevokeToken(ctx context.Context, tokenString string) error {
	// Validate token
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	if claims.TokenType != "refresh" {
		return fmt.Errorf("can only revoke refresh tokens")
	}

	// Get token ID from claims
	tokenID, err := uuid.Parse(claims.TokenID)
	if err != nil {
		return fmt.Errorf("invalid token ID: %w", err)
	}

	// Revoke token
	return s.refreshTokenRepo.Revoke(ctx, tokenID, "user_logout")
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (s *JWTService) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return s.refreshTokenRepo.RevokeAllForUser(ctx, userID, "user_logout_all")
}

// hashToken creates a SHA-256 hash of the token
func (s *JWTService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ParseUnverified parses token without validation (for debugging)
func (s *JWTService) ParseUnverified(tokenString string) (*Claims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	return claims, nil
}

// Request/Response types

type TokenRequest struct {
	UserID         uuid.UUID
	Email          string
	Name           string
	OrganizationID string
	Role           string
	Permissions    []string
	DeviceInfo     map[string]interface{}
	IPAddress      string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}
