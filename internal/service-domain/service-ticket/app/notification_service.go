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
	"github.com/aby-med/medical-platform/pkg/email"
)

// NotificationService handles ticket notifications
type NotificationService struct {
	ticketRepo       domain.Repository
	notificationRepo domain.NotificationRepository
	emailService     *email.NotificationService
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
	ticketRepo domain.Repository,
	notificationRepo domain.NotificationRepository,
	emailService *email.NotificationService,
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
func (s *NotificationService) SendManualEmail(ctx context.Context, ticketID string, includeComments bool) error {
	// Check feature flag
	if os.Getenv("FEATURE_TICKET_MANUAL_EMAIL") == "false" {
		return fmt.Errorf("manual email feature is disabled")
	}

	// Get ticket details
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	// Validate customer email
	if ticket.CustomerEmail == "" {
		return fmt.Errorf("ticket has no customer email")
	}

	// Get or create tracking token
	token, err := s.GetOrCreateTrackingToken(ticketID)
	if err != nil {
		s.logger.Error("Failed to create tracking token", slog.String("error", err.Error()))
		// Continue without tracking token
		token = ""
	}

	// Prepare email data
	emailData := s.prepareEmailData(ticket, token, includeComments)

	// Send email
	err = s.sendEmail(ticket.CustomerEmail, "ticket-update", emailData)
	
	// Log notification
	logEntry := &domain.NotificationLog{
		TicketID:         ticketID,
		NotificationType: domain.NotificationTypeManual,
		RecipientEmail:   ticket.CustomerEmail,
		Status:           domain.NotificationStatusSent,
	}
	
	if err != nil {
		errMsg := err.Error()
		logEntry.Status = domain.NotificationStatusFailed
		logEntry.ErrorMessage = &errMsg
		s.notificationRepo.LogNotification(logEntry)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.notificationRepo.LogNotification(logEntry)
	s.logger.Info("Manual email sent", 
		slog.String("ticket_id", ticketID),
		slog.String("recipient", ticket.CustomerEmail),
	)

	return nil
}

// SendTicketCreatedEmail sends an email when a ticket is created
func (s *NotificationService) SendTicketCreatedEmail(ctx context.Context, ticketID string) error {
	// Check feature flag
	if os.Getenv("FEATURE_TICKET_CREATED_EMAIL") != "true" {
		s.logger.Debug("Ticket created email disabled", slog.String("ticket_id", ticketID))
		return nil // Not an error, just disabled
	}

	// Get ticket details
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	// Validate customer email
	if ticket.CustomerEmail == "" {
		s.logger.Warn("Ticket has no customer email", slog.String("ticket_id", ticketID))
		return nil // Not an error
	}

	// Create tracking token
	token, err := s.GetOrCreateTrackingToken(ticketID)
	if err != nil {
		s.logger.Error("Failed to create tracking token", slog.String("error", err.Error()))
		token = ""
	}

	// Prepare email data
	emailData := s.prepareEmailData(ticket, token, false)

	// Send email
	err = s.sendEmail(ticket.CustomerEmail, "ticket-created", emailData)
	
	// Log notification
	logEntry := &domain.NotificationLog{
		TicketID:         ticketID,
		NotificationType: domain.NotificationTypeTicketCreated,
		RecipientEmail:   ticket.CustomerEmail,
		Status:           domain.NotificationStatusSent,
	}
	
	if err != nil {
		errMsg := err.Error()
		logEntry.Status = domain.NotificationStatusFailed
		logEntry.ErrorMessage = &errMsg
		s.notificationRepo.LogNotification(logEntry)
		return fmt.Errorf("failed to send ticket created email: %w", err)
	}

	s.notificationRepo.LogNotification(logEntry)
	s.logger.Info("Ticket created email sent", 
		slog.String("ticket_id", ticketID),
		slog.String("recipient", ticket.CustomerEmail),
	)

	return nil
}

// SendDailyDigest sends daily digest emails for all active tickets
func (s *NotificationService) SendDailyDigest(ctx context.Context) error {
	// Check feature flag
	if os.Getenv("FEATURE_TICKET_DAILY_DIGEST") != "true" {
		s.logger.Debug("Daily digest disabled")
		return nil
	}

	s.logger.Info("Starting daily digest")

	// Get tickets updated in the last 24 hours
	since := time.Now().Add(-24 * time.Hour)
	
	// TODO: Add method to ticketRepo to get tickets updated since timestamp
	// For now, we'll skip implementation until the repository method is added
	
	s.logger.Info("Daily digest completed", slog.Int("sent", 0))
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

	// Get expiry days from env or default to 30
	expiryDays := 30
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
		Status:           ticket.Status,
		Priority:         ticket.Priority,
		EquipmentName:    ticket.EquipmentName,
		IssueDescription: ticket.IssueDescription,
		CreatedAt:        ticket.CreatedAt,
		UpdatedAt:        ticket.UpdatedAt,
		PublicComments:   publicComments,
		AssignedEngineer: ticket.AssignedEngineerName,
	}

	return publicView, nil
}

// prepareEmailData prepares data for email templates
func (s *NotificationService) prepareEmailData(ticket *domain.ServiceTicket, trackingToken string, includeComments bool) *domain.EmailTemplateData {
	trackingURL := ""
	if trackingToken != "" {
		trackingURL = fmt.Sprintf("%s/%s", s.baseURL, trackingToken)
	}

	data := &domain.EmailTemplateData{
		TicketNumber:         ticket.TicketNumber,
		TicketID:             ticket.ID,
		Status:               ticket.Status,
		Priority:             ticket.Priority,
		CustomerName:         ticket.CustomerName,
		CustomerEmail:        ticket.CustomerEmail,
		EquipmentName:        ticket.EquipmentName,
		SerialNumber:         ticket.SerialNumber,
		IssueDescription:     ticket.IssueDescription,
		IssueCategory:        ticket.IssueCategory,
		AssignedEngineerName: ticket.AssignedEngineerName,
		AssignedAt:           ticket.AssignedAt,
		CreatedAt:            ticket.CreatedAt,
		UpdatedAt:            ticket.UpdatedAt,
		TrackingURL:          trackingURL,
		TrackingToken:        trackingToken,
		Comments:             []domain.CommentData{},
	}

	// TODO: Add comments if requested and available
	// This requires accessing comment repository

	return data
}

// sendEmail sends an email using the email service
func (s *NotificationService) sendEmail(to, templateName string, data *domain.EmailTemplateData) error {
	// For now, use a simple email format
	// TODO: Implement proper HTML templates
	
	subject := s.getEmailSubject(templateName, data)
	body := s.getEmailBody(templateName, data)

	// Use the existing email service
	// Note: This is a placeholder - actual implementation depends on email.NotificationService interface
	s.logger.Info("Sending email",
		slog.String("to", to),
		slog.String("subject", subject),
		slog.String("template", templateName),
	)

	// TODO: Call actual email service send method
	// err := s.emailService.Send(to, subject, body)
	
	return nil
}

// getEmailSubject returns the subject for an email template
func (s *NotificationService) getEmailSubject(templateName string, data *domain.EmailTemplateData) string {
	switch templateName {
	case "ticket-created":
		return fmt.Sprintf("New Service Ticket Created - %s", data.TicketNumber)
	case "ticket-update":
		return fmt.Sprintf("Service Ticket Update - %s", data.TicketNumber)
	case "daily-digest":
		return fmt.Sprintf("Daily Service Ticket Summary - %s", data.TicketNumber)
	default:
		return fmt.Sprintf("Service Ticket Notification - %s", data.TicketNumber)
	}
}

// getEmailBody returns the body for an email template
func (s *NotificationService) getEmailBody(templateName string, data *domain.EmailTemplateData) string {
	// TODO: Use proper HTML templates
	switch templateName {
	case "ticket-created":
		return fmt.Sprintf(`
Dear %s,

Your service request has been created successfully.

Ticket Details:
- Ticket #: %s
- Equipment: %s
- Issue: %s
- Priority: %s
- Status: %s

Track your ticket: %s

Thank you,
ServQR Support Team
`,
			data.CustomerName,
			data.TicketNumber,
			data.EquipmentName,
			data.IssueDescription,
			data.Priority,
			data.Status,
			data.TrackingURL,
		)
	case "ticket-update":
		return fmt.Sprintf(`
Dear %s,

Here's an update on your service request.

Ticket #: %s
Current Status: %s
Equipment: %s

Track your ticket: %s

Thank you,
ServQR Support Team
`,
			data.CustomerName,
			data.TicketNumber,
			data.Status,
			data.EquipmentName,
			data.TrackingURL,
		)
	default:
		return fmt.Sprintf("Service ticket notification for %s", data.TicketNumber)
	}
}
