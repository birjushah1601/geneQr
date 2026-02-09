package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aby-med/medical-platform/internal/infrastructure/email"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// InvitationHandler handles user invitation operations
type InvitationHandler struct {
	db            *sqlx.DB
	emailService  *email.NotificationService
	logger        *slog.Logger
	appURL        string
}

// NewInvitationHandler creates a new invitation handler
func NewInvitationHandler(db *sqlx.DB, emailService *email.NotificationService, logger *slog.Logger) *InvitationHandler {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		// Default to production HTTPS URL
		appURL = "https://servqr.com"
	}

	return &InvitationHandler{
		db:           db,
		emailService: emailService,
		logger:       logger.With(slog.String("component", "invitation_handler")),
		appURL:       appURL,
	}
}

// CreateInvitationRequest represents the request to create an invitation
type CreateInvitationRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Role  string `json:"role"` // admin, manager, viewer
}

// CreateInvitationResponse represents the response
type CreateInvitationResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ValidateInvitationResponse represents the validation response
type ValidateInvitationResponse struct {
	Valid            bool   `json:"valid"`
	Email            string `json:"email,omitempty"`
	OrganizationID   string `json:"organization_id,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
	Role             string `json:"role,omitempty"`
	Name             string `json:"name,omitempty"`
	Error            string `json:"error,omitempty"`
}

// AcceptInvitationRequest represents the request to accept an invitation
type AcceptInvitationRequest struct {
	FullName        string `json:"full_name"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// AcceptInvitationResponse represents the response
type AcceptInvitationResponse struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	OrganizationID string `json:"organization_id"`
	Message        string `json:"message"`
}

// CreateInvitation creates a new user invitation
func (h *InvitationHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgID := chi.URLParam(r, "orgId")

	var req CreateInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Role == "" {
		h.respondError(w, http.StatusBadRequest, "email and role are required")
		return
	}

	// Validate role
	validRoles := map[string]bool{"admin": true, "manager": true, "viewer": true, "engineer": true}
	if !validRoles[req.Role] {
		h.respondError(w, http.StatusBadRequest, "invalid role. Must be: admin, manager, viewer, or engineer")
		return
	}

	// Check if user already exists
	var existingUserID string
	err := h.db.GetContext(ctx, &existingUserID, "SELECT id FROM users WHERE email = $1", req.Email)
	if err == nil {
		h.respondError(w, http.StatusConflict, "user with this email already exists")
		return
	}

	// Check if invitation already exists
	var existingInviteID string
	err = h.db.GetContext(ctx, &existingInviteID,
		"SELECT id FROM invitations WHERE email = $1 AND organization_id = $2 AND status = 'pending'",
		req.Email, orgID)
	if err == nil {
		h.respondError(w, http.StatusConflict, "pending invitation already exists for this email")
		return
	}

	// Generate secure token
	token := generateInvitationToken()

	// Set expiry (7 days from now)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Get current user from context (inviter)
	inviterID := r.Context().Value("user_id")
	var inviterName string
	if inviterID != nil {
		h.db.GetContext(ctx, &inviterName, "SELECT full_name FROM users WHERE id = $1", inviterID)
	}
	if inviterName == "" {
		inviterName = "Administrator"
	}

	// Get organization name
	var orgName string
	err = h.db.GetContext(ctx, &orgName, "SELECT name FROM organizations WHERE id = $1", orgID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "organization not found")
		return
	}

	// Create invitation
	inviteID := uuid.New().String()
	metadata := map[string]string{"name": req.Name, "phone": req.Phone}
	metadataJSON, _ := json.Marshal(metadata)

	_, err = h.db.ExecContext(ctx, `
		INSERT INTO invitations (id, email, organization_id, role, invited_by, token, expires_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, inviteID, req.Email, orgID, req.Role, inviterID, token, expiresAt, metadataJSON)

	if err != nil {
		h.logger.Error("failed to create invitation", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "failed to create invitation")
		return
	}

	// Send invitation email
	inviteURL := fmt.Sprintf("%s/invite/accept?token=%s", h.appURL, token)
	emailData := email.InvitationData{
		InviteeName:      req.Name,
		InviteeEmail:     req.Email,
		InviterName:      inviterName,
		OrganizationName: orgName,
		Role:             req.Role,
		InviteURL:        inviteURL,
		ExpiresAt:        expiresAt.Format("January 2, 2006 at 3:04 PM"),
	}

	if err := h.emailService.SendInvitationEmail(ctx, emailData); err != nil {
		h.logger.Error("failed to send invitation email", slog.String("error", err.Error()))
		// Don't fail the request - invitation was created
	}

	// Respond
	resp := CreateInvitationResponse{
		ID:        inviteID,
		Email:     req.Email,
		Role:      req.Role,
		Status:    "pending",
		ExpiresAt: expiresAt,
	}

	h.respondJSON(w, http.StatusCreated, resp)
}

// ValidateInvitation validates an invitation token
func (h *InvitationHandler) ValidateInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")

	var invitation struct {
		ID               string    `db:"id"`
		Email            string    `db:"email"`
		OrganizationID   string    `db:"organization_id"`
		OrganizationName string    `db:"organization_name"`
		Role             string    `db:"role"`
		Status           string    `db:"status"`
		ExpiresAt        time.Time `db:"expires_at"`
		Metadata         string    `db:"metadata"`
	}

	err := h.db.GetContext(ctx, &invitation, `
		SELECT 
			i.id, i.email, i.organization_id, o.name as organization_name,
			i.role, i.status, i.expires_at, i.metadata
		FROM invitations i
		JOIN organizations o ON i.organization_id = o.id
		WHERE i.token = $1
	`, token)

	if err != nil {
		h.respondJSON(w, http.StatusOK, ValidateInvitationResponse{
			Valid: false,
			Error: "Invalid invitation token",
		})
		return
	}

	// Check status
	if invitation.Status != "pending" {
		h.respondJSON(w, http.StatusOK, ValidateInvitationResponse{
			Valid: false,
			Error: "Invitation already accepted or expired",
		})
		return
	}

	// Check expiry
	if time.Now().After(invitation.ExpiresAt) {
		// Mark as expired
		h.db.ExecContext(ctx, "UPDATE invitations SET status = 'expired' WHERE id = $1", invitation.ID)

		h.respondJSON(w, http.StatusOK, ValidateInvitationResponse{
			Valid: false,
			Error: "Invitation has expired. Please request a new invitation.",
		})
		return
	}

	// Extract name from metadata
	var metadata map[string]string
	json.Unmarshal([]byte(invitation.Metadata), &metadata)

	h.respondJSON(w, http.StatusOK, ValidateInvitationResponse{
		Valid:            true,
		Email:            invitation.Email,
		OrganizationID:   invitation.OrganizationID,
		OrganizationName: invitation.OrganizationName,
		Role:             invitation.Role,
		Name:             metadata["name"],
	})
}

// AcceptInvitation accepts an invitation and creates user account
func (h *InvitationHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")

	var req AcceptInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate input
	if req.FullName == "" || req.Password == "" {
		h.respondError(w, http.StatusBadRequest, "full_name and password are required")
		return
	}

	if req.Password != req.PasswordConfirm {
		h.respondError(w, http.StatusBadRequest, "passwords do not match")
		return
	}

	// Validate password strength (min 8 chars)
	if len(req.Password) < 8 {
		h.respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	// Get invitation
	var invitation struct {
		ID             string `db:"id"`
		Email          string `db:"email"`
		OrganizationID string `db:"organization_id"`
		Role           string `db:"role"`
		Status         string `db:"status"`
		ExpiresAt      time.Time `db:"expires_at"`
	}

	err := h.db.GetContext(ctx, &invitation, `
		SELECT id, email, organization_id, role, status, expires_at
		FROM invitations
		WHERE token = $1
	`, token)

	if err != nil {
		h.respondError(w, http.StatusNotFound, "invalid invitation token")
		return
	}

	// Check status
	if invitation.Status != "pending" {
		h.respondError(w, http.StatusBadRequest, "invitation already accepted or expired")
		return
	}

	// Check expiry
	if time.Now().After(invitation.ExpiresAt) {
		h.db.ExecContext(ctx, "UPDATE invitations SET status = 'expired' WHERE id = $1", invitation.ID)
		h.respondError(w, http.StatusBadRequest, "invitation has expired")
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("failed to hash password", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	// Start transaction
	tx, err := h.db.BeginTxx(ctx, nil)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to start transaction")
		return
	}
	defer tx.Rollback()

	// Create user
	userID := uuid.New().String()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, full_name, phone, email_verified, status)
		VALUES ($1, $2, $3, $4, $5, TRUE, 'active')
	`, userID, invitation.Email, passwordHash, req.FullName, req.Phone)

	if err != nil {
		h.logger.Error("failed to create user", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "failed to create user account")
		return
	}

	// Link user to organization
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_organizations (id, user_id, organization_id, role, is_primary, status)
		VALUES ($1, $2, $3, $4, TRUE, 'active')
	`, uuid.New().String(), userID, invitation.OrganizationID, invitation.Role)

	if err != nil {
		h.logger.Error("failed to link user to organization", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "failed to link user to organization")
		return
	}

	// Mark invitation as accepted
	_, err = tx.ExecContext(ctx, `
		UPDATE invitations
		SET status = 'accepted', accepted_at = NOW(), accepted_by_user_id = $1
		WHERE id = $2
	`, userID, invitation.ID)

	if err != nil {
		h.logger.Error("failed to update invitation", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "failed to update invitation")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to commit transaction")
		return
	}

	// Respond
	resp := AcceptInvitationResponse{
		UserID:         userID,
		Email:          invitation.Email,
		OrganizationID: invitation.OrganizationID,
		Message:        "Account created successfully. You can now login.",
	}

	h.respondJSON(w, http.StatusCreated, resp)
}

// generateInvitationToken generates a cryptographically secure random token
func generateInvitationToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback - use timestamp (should never happen)
		return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// Helper methods
func (h *InvitationHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *InvitationHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
