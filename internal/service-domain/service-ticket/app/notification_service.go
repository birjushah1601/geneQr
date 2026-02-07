package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
)

// EmailService is a placeholder interface for email functionality
type EmailService interface {
	Send(to, subject, body string) error
}

// NotificationService handles ticket notifications
type NotificationService struct {
	ticketRepo       domain.TicketRepository
	notificationRepo domain.NotificationRepository
	emailService     EmailService
	logger           *slog.Logger
	baseURL          string
}

// NotificationConfig holds notification configuration
type NotificationConfig struct {
	BaseURL              string
	TokenExpiryDays      int
	EnableCreatedEmail   bool
	EnableDailyDigest    bool
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	ticketRepo domain.TicketRepository,
	notificationRepo domain.NotificationRepository,
	emailService EmailService,
	logger *slog.Logger,
	config *NotificationConfig,
) *NotificationService {
	if config.TokenExpiryDays == 0 {
		config.TokenExpiryDays = 30
	}
	if config.BaseURL == "" {
		config.BaseURL = os.Getenv("TICKET_TRACKING_BASE_URL")
		if config.BaseURL == "" {
			config.BaseURL = "https://servqr.com/track"
		}
	}

	return &NotificationService{
		ticketRepo:       ticketRepo,
		notificationRepo: notificationRepo,
		emailService:     emailService,
		logger:           logger.With(slog.String("component", "notification_service")),
		baseURL:          config.BaseURL,
	}
}

// SendManualEmail sends a manual email notification for a ticket
// TODO: Implement when email service and customer_email field are available
func (s *NotificationService) SendManualEmail(ctx context.Context, ticketID string, includeComments bool) error {
	s.logger.Warn("Email functionality not yet implemented",
		slog.String("ticket_id", ticketID))
	return fmt.Errorf("email functionality not yet implemented")
}

// SendTicketCreatedEmail sends an email when a ticket is created
// TODO: Implement when email service and customer_email field are available
func (s *NotificationService) SendTicketCreatedEmail(ctx context.Context, ticketID string) error {
	s.logger.Debug("Ticket created email not yet implemented", slog.String("ticket_id", ticketID))
	return nil
}

// SendDailyDigest sends daily digest emails for all active tickets
// TODO: Implement when email service is available
func (s *NotificationService) SendDailyDigest(ctx context.Context) error {
	s.logger.Debug("Daily digest not yet implemented")
	return nil
}

// GetOrCreateTrackingToken gets an existing token or creates a new one
func (s *NotificationService) GetOrCreateTrackingToken(ticketID string) (string, error) {
	// Generate secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	tokenStr := hex.EncodeToString(tokenBytes)

	// Set expiry to 5 years in the future (effectively permanent for ticket tracking)
	// Tokens should remain valid as long as the ticket exists
	// Customer may need to reference it months/years later for warranty/records
	expiryDays := 365 * 5 // 5 years
	if envDays := os.Getenv("TRACKING_TOKEN_EXPIRY_DAYS"); envDays != "" {
		if days, err := strconv.Atoi(envDays); err == nil {
			expiryDays = days
		}
	}

	// Create tracking token
	token := &domain.TrackingToken{
		TicketID:  ticketID,
		Token:     tokenStr,
		ExpiresAt: time.Now().AddDate(0, 0, expiryDays),
	}

	err := s.notificationRepo.CreateTrackingToken(token)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// GetPublicTicketView retrieves public-safe ticket information by tracking token
func (s *NotificationService) GetPublicTicketView(ctx context.Context, token string) (*domain.PublicTicketView, error) {
	// Get tracking token
	trackingToken, err := s.notificationRepo.GetTrackingToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired tracking token")
	}

	// Get ticket
	ticket, err := s.ticketRepo.GetByID(ctx, trackingToken.TicketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// TODO: Get comments (requires comment repository method)
	// For now, return empty comments array
	publicComments := []domain.PublicComment{}

	// Build public view
	publicView := &domain.PublicTicketView{
		TicketNumber:     ticket.TicketNumber,
		Status:           string(ticket.Status),
		Priority:         string(ticket.Priority),
		EquipmentName:    ticket.EquipmentName,
		IssueDescription: ticket.IssueDescription,
		CreatedAt:        ticket.CreatedAt,
		UpdatedAt:        ticket.UpdatedAt,
		PublicComments:   publicComments,
		AssignedEngineer: ticket.AssignedEngineerName,
	}

	return publicView, nil
}

// TODO: Email helper methods will be implemented when email service is integrated
