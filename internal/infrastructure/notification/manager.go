package notification

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aby-med/internal/infrastructure/config"
	"github.com/aby-med/internal/infrastructure/email"
)

// Manager handles all notification types with feature flag support
type Manager struct {
	emailService  *email.NotificationService
	featureFlags  *config.FeatureFlags
	logger        *slog.Logger
	adminEmail    string // Default admin email for notifications
}

// NewManager creates a new notification manager
func NewManager(
	emailService *email.NotificationService,
	featureFlags *config.FeatureFlags,
	logger *slog.Logger,
	adminEmail string,
) *Manager {
	return &Manager{
		emailService: emailService,
		featureFlags: featureFlags,
		logger:       logger,
		adminEmail:   adminEmail,
	}
}

// SendTicketCreatedNotifications sends all enabled notifications for ticket creation
func (m *Manager) SendTicketCreatedNotifications(ctx context.Context, data TicketCreatedData) error {
	var errors []error

	// Email notifications
	if m.featureFlags.ShouldSendEmailNotification("ticket_created") {
		m.logger.Info("Sending ticket created email notifications",
			slog.String("ticket", data.TicketNumber),
			slog.String("customer", data.CustomerEmail),
		)

		// Set admin email if not provided
		if data.AdminEmail == "" {
			data.AdminEmail = m.adminEmail
		}

		// Convert to email data structure
		emailData := email.TicketCreatedData{
			TicketNumber:  data.TicketNumber,
			CustomerName:  data.CustomerName,
			CustomerEmail: data.CustomerEmail,
			CustomerPhone: data.CustomerPhone,
			EquipmentName: data.EquipmentName,
			Description:   data.Description,
			Priority:      data.Priority,
			AdminEmail:    data.AdminEmail,
		}

		if err := m.emailService.SendTicketCreatedNotification(ctx, emailData); err != nil {
			m.logger.Error("Failed to send ticket created email",
				slog.String("ticket", data.TicketNumber),
				slog.String("error", err.Error()),
			)
			errors = append(errors, fmt.Errorf("email notification failed: %w", err))
		} else {
			m.logger.Info("Ticket created email notifications sent successfully",
				slog.String("ticket", data.TicketNumber),
			)
		}
	} else {
		m.logger.Debug("Ticket created email notifications disabled by feature flag",
			slog.String("ticket", data.TicketNumber),
		)
	}

	// SMS notifications (future)
	if m.featureFlags.ShouldSendSMSNotification("ticket_created") {
		m.logger.Info("SMS notification for ticket created (not implemented yet)",
			slog.String("ticket", data.TicketNumber),
		)
		// TODO: Implement SMS sending
	}

	// WhatsApp notifications (future)
	if m.featureFlags.ShouldSendWhatsAppNotification("ticket_created") {
		m.logger.Info("WhatsApp notification for ticket created (not implemented yet)",
			slog.String("ticket", data.TicketNumber),
		)
		// TODO: Implement WhatsApp sending
	}

	// Return error if any notification failed
	if len(errors) > 0 {
		return fmt.Errorf("some notifications failed: %v", errors)
	}

	return nil
}

// SendEngineerAssignedNotifications sends all enabled notifications for engineer assignment
func (m *Manager) SendEngineerAssignedNotifications(ctx context.Context, data EngineerAssignedData) error {
	var errors []error

	// Email notifications
	if m.featureFlags.ShouldSendEmailNotification("engineer_assigned") {
		m.logger.Info("Sending engineer assigned email notifications",
			slog.String("ticket", data.TicketNumber),
			slog.String("engineer", data.EngineerName),
		)

		// Convert to email data structure
		emailData := email.TicketAssignedData{
			TicketNumber:  data.TicketNumber,
			CustomerName:  data.CustomerName,
			CustomerEmail: data.CustomerEmail,
			EngineerName:  data.EngineerName,
			EngineerEmail: data.EngineerEmail,
			EngineerPhone: data.EngineerPhone,
			EquipmentName: data.EquipmentName,
			Description:   data.Description,
			Priority:      data.Priority,
		}

		if err := m.emailService.SendTicketAssignedNotification(ctx, emailData); err != nil {
			m.logger.Error("Failed to send engineer assigned email",
				slog.String("ticket", data.TicketNumber),
				slog.String("error", err.Error()),
			)
			errors = append(errors, fmt.Errorf("email notification failed: %w", err))
		} else {
			m.logger.Info("Engineer assigned email notifications sent successfully",
				slog.String("ticket", data.TicketNumber),
			)
		}
	} else {
		m.logger.Debug("Engineer assigned email notifications disabled by feature flag",
			slog.String("ticket", data.TicketNumber),
		)
	}

	// SMS notifications (future)
	if m.featureFlags.ShouldSendSMSNotification("engineer_assigned") {
		m.logger.Info("SMS notification for engineer assigned (not implemented yet)",
			slog.String("ticket", data.TicketNumber),
		)
		// TODO: Implement SMS sending
	}

	// WhatsApp notifications (future)
	if m.featureFlags.ShouldSendWhatsAppNotification("engineer_assigned") {
		m.logger.Info("WhatsApp notification for engineer assigned (not implemented yet)",
			slog.String("ticket", data.TicketNumber),
		)
		// TODO: Implement WhatsApp sending
	}

	// Return error if any notification failed
	if len(errors) > 0 {
		return fmt.Errorf("some notifications failed: %v", errors)
	}

	return nil
}

// SendStatusChangedNotifications sends all enabled notifications for status change
func (m *Manager) SendStatusChangedNotifications(ctx context.Context, data StatusChangedData) error {
	var errors []error

	// Email notifications
	if m.featureFlags.ShouldSendEmailNotification("status_changed") {
		m.logger.Info("Sending status changed email notifications",
			slog.String("ticket", data.TicketNumber),
			slog.String("old_status", data.OldStatus),
			slog.String("new_status", data.NewStatus),
		)

		// Set admin email if not provided
		if data.AdminEmail == "" {
			data.AdminEmail = m.adminEmail
		}

		// Convert to email data structure
		emailData := email.TicketStatusChangedData{
			TicketNumber:  data.TicketNumber,
			CustomerName:  data.CustomerName,
			CustomerEmail: data.CustomerEmail,
			OldStatus:     data.OldStatus,
			NewStatus:     data.NewStatus,
			EquipmentName: data.EquipmentName,
			UpdatedBy:     data.UpdatedBy,
			AdminEmail:    data.AdminEmail,
		}

		if err := m.emailService.SendTicketStatusChangedNotification(ctx, emailData); err != nil {
			m.logger.Error("Failed to send status changed email",
				slog.String("ticket", data.TicketNumber),
				slog.String("error", err.Error()),
			)
			errors = append(errors, fmt.Errorf("email notification failed: %w", err))
		} else {
			m.logger.Info("Status changed email notifications sent successfully",
				slog.String("ticket", data.TicketNumber),
			)
		}
	} else {
		m.logger.Debug("Status changed email notifications disabled by feature flag",
			slog.String("ticket", data.TicketNumber),
		)
	}

	// SMS notifications (future)
	if m.featureFlags.ShouldSendSMSNotification("status_changed") {
		m.logger.Info("SMS notification for status changed (not implemented yet)",
			slog.String("ticket", data.TicketNumber),
		)
		// TODO: Implement SMS sending
	}

	// WhatsApp notifications (future)
	if m.featureFlags.ShouldSendWhatsAppNotification("status_changed") {
		m.logger.Info("WhatsApp notification for status changed (not implemented yet)",
			slog.String("ticket", data.TicketNumber),
		)
		// TODO: Implement WhatsApp sending
	}

	// Return error if any notification failed
	if len(errors) > 0 {
		return fmt.Errorf("some notifications failed: %v", errors)
	}

	return nil
}

// GetFeatureStatus returns the current status of all notification features
func (m *Manager) GetFeatureStatus() map[string]bool {
	return m.featureFlags.GetFeatureFlagsStatus()
}

// Data structures for notification events

// TicketCreatedData contains data for ticket creation notification
type TicketCreatedData struct {
	TicketNumber  string
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	EquipmentName string
	Description   string
	Priority      string
	AdminEmail    string // Optional, will use default if not provided
}

// EngineerAssignedData contains data for engineer assignment notification
type EngineerAssignedData struct {
	TicketNumber  string
	CustomerName  string
	CustomerEmail string
	EngineerName  string
	EngineerEmail string
	EngineerPhone string
	EquipmentName string
	Description   string
	Priority      string
}

// StatusChangedData contains data for status change notification
type StatusChangedData struct {
	TicketNumber  string
	CustomerName  string
	CustomerEmail string
	OldStatus     string
	NewStatus     string
	EquipmentName string
	UpdatedBy     string
	AdminEmail    string // Optional, will use default if not provided
}
