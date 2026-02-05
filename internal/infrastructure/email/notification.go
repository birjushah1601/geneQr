package email

import (
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// NotificationService handles sending notification emails
type NotificationService struct {
	apiKey    string
	fromEmail string
	fromName  string
}

// NewNotificationService creates a new notification service
func NewNotificationService(apiKey, fromEmail, fromName string) *NotificationService {
	return &NotificationService{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

// InvitationData contains data for user invitation email
type InvitationData struct {
	InviteeName        string
	InviteeEmail       string
	InviterName        string
	OrganizationName   string
	Role               string
	InviteURL          string
	ExpiresAt          string
}

// TicketCreatedData contains data for ticket creation notification
type TicketCreatedData struct {
	TicketNumber    string
	CustomerName    string
	CustomerEmail   string
	CustomerPhone   string
	EquipmentName   string
	Description     string
	Priority        string
	AdminEmail      string
}

// TicketAssignedData contains data for ticket assignment notification
type TicketAssignedData struct {
	TicketNumber    string
	CustomerName    string
	CustomerEmail   string
	EngineerName    string
	EngineerEmail   string
	EngineerPhone   string
	EquipmentName   string
	Description     string
	Priority        string
}

// TicketStatusChangedData contains data for status change notification
type TicketStatusChangedData struct {
	TicketNumber    string
	CustomerName    string
	CustomerEmail   string
	OldStatus       string
	NewStatus       string
	EquipmentName   string
	UpdatedBy       string
	AdminEmail      string
}

// SendTicketCreatedNotification sends email when a ticket is created
func (s *NotificationService) SendTicketCreatedNotification(ctx context.Context, data TicketCreatedData) error {
	// Email to customer
	if data.CustomerEmail != "" {
		if err := s.sendCustomerTicketCreated(ctx, data); err != nil {
			return fmt.Errorf("failed to send customer notification: %w", err)
		}
	}

	// Email to admin
	if data.AdminEmail != "" {
		if err := s.sendAdminTicketCreated(ctx, data); err != nil {
			return fmt.Errorf("failed to send admin notification: %w", err)
		}
	}

	return nil
}

// SendTicketAssignedNotification sends email when engineer is assigned
func (s *NotificationService) SendTicketAssignedNotification(ctx context.Context, data TicketAssignedData) error {
	// Email to customer
	if data.CustomerEmail != "" {
		if err := s.sendCustomerEngineerAssigned(ctx, data); err != nil {
			return fmt.Errorf("failed to send customer notification: %w", err)
		}
	}

	// Email to engineer
	if data.EngineerEmail != "" {
		if err := s.sendEngineerAssigned(ctx, data); err != nil {
			return fmt.Errorf("failed to send engineer notification: %w", err)
		}
	}

	return nil
}

// SendTicketStatusChangedNotification sends email when status changes
func (s *NotificationService) SendTicketStatusChangedNotification(ctx context.Context, data TicketStatusChangedData) error {
	// Email to customer
	if data.CustomerEmail != "" {
		if err := s.sendCustomerStatusChanged(ctx, data); err != nil {
			return fmt.Errorf("failed to send customer notification: %w", err)
		}
	}

	// Email to admin for important status changes
	if data.AdminEmail != "" && (data.NewStatus == "resolved" || data.NewStatus == "closed" || data.NewStatus == "cancelled") {
		if err := s.sendAdminStatusChanged(ctx, data); err != nil {
			return fmt.Errorf("failed to send admin notification: %w", err)
		}
	}

	return nil
}

// sendCustomerTicketCreated sends ticket creation email to customer
func (s *NotificationService) sendCustomerTicketCreated(ctx context.Context, data TicketCreatedData) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	subject := fmt.Sprintf("Service Ticket Created - %s", data.TicketNumber)
	to := mail.NewEmail(data.CustomerName, data.CustomerEmail)

	plainText := fmt.Sprintf(`
Dear %s,

Your service ticket has been successfully created.

Ticket Details:
--------------------------------------------------
Ticket Number: %s
Equipment: %s
Priority: %s
Description: %s

Contact Information:
Phone: %s
Email: %s

What Happens Next:
1. Our team will review your request
2. An engineer will be assigned shortly
3. You'll receive updates via email and SMS

Thank you for contacting us.

Best regards,
ServQR Service Team
`, data.CustomerName, data.TicketNumber, data.EquipmentName, data.Priority, data.Description, data.CustomerPhone, data.CustomerEmail)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2563eb; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 5px 5px; }
        .ticket-info { background-color: #fff; border-left: 4px solid #2563eb; padding: 15px; margin: 20px 0; }
        .info-row { padding: 8px 0; border-bottom: 1px solid #e5e7eb; }
        .info-label { font-weight: bold; color: #6b7280; }
        .info-value { color: #111827; }
        .priority-high { color: #dc2626; font-weight: bold; }
        .priority-medium { color: #f59e0b; font-weight: bold; }
        .priority-low { color: #10b981; font-weight: bold; }
        .steps { background-color: #eff6ff; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>âœ… Service Ticket Created</h1>
        </div>
        <div class="content">
            <p>Dear <strong>%s</strong>,</p>
            <p>Your service ticket has been successfully created. We'll get back to you shortly.</p>
            
            <div class="ticket-info">
                <div class="info-row">
                    <span class="info-label">Ticket Number:</span>
                    <span class="info-value"><strong>%s</strong></span>
                </div>
                <div class="info-row">
                    <span class="info-label">Equipment:</span>
                    <span class="info-value">%s</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Priority:</span>
                    <span class="info-value priority-%s">%s</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Description:</span>
                    <span class="info-value">%s</span>
                </div>
            </div>

            <div class="steps">
                <h3>ðŸ“‹ What Happens Next:</h3>
                <ol>
                    <li>Our team will review your request</li>
                    <li>An engineer will be assigned shortly</li>
                    <li>You'll receive updates via email and SMS</li>
                </ol>
            </div>

            <p>Thank you for contacting ServQR Service Team.</p>
        </div>
        <div class="footer">
            <p>This is an automated message from ServQR Platform</p>
            <p>&copy; 2025 ServQR. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, data.CustomerName, data.TicketNumber, data.EquipmentName, data.Priority, data.Priority, data.Description)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: status %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

// sendAdminTicketCreated sends ticket creation email to admin
func (s *NotificationService) sendAdminTicketCreated(ctx context.Context, data TicketCreatedData) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	subject := fmt.Sprintf("New Service Ticket - %s (%s Priority)", data.TicketNumber, data.Priority)
	to := mail.NewEmail("Admin", data.AdminEmail)

	plainText := fmt.Sprintf(`
New Service Ticket Created

Ticket Details:
--------------------------------------------------
Ticket Number: %s
Customer: %s
Phone: %s
Email: %s
Equipment: %s
Priority: %s
Description: %s

Action Required:
Please review and assign an engineer.

ServQR Admin System
`, data.TicketNumber, data.CustomerName, data.CustomerPhone, data.CustomerEmail, data.EquipmentName, data.Priority, data.Description)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #dc2626; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 5px 5px; }
        .ticket-info { background-color: #fff; border-left: 4px solid #dc2626; padding: 15px; margin: 20px 0; }
        .info-row { padding: 8px 0; border-bottom: 1px solid #e5e7eb; }
        .info-label { font-weight: bold; color: #6b7280; }
        .action-required { background-color: #fef2f2; border: 2px solid #dc2626; padding: 15px; border-radius: 5px; margin: 20px 0; text-align: center; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>ðŸ”” New Service Ticket</h1>
        </div>
        <div class="content">
            <p><strong>A new service ticket requires your attention.</strong></p>
            
            <div class="ticket-info">
                <div class="info-row">
                    <span class="info-label">Ticket Number:</span>
                    <span><strong>%s</strong></span>
                </div>
                <div class="info-row">
                    <span class="info-label">Customer:</span>
                    <span>%s</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Phone:</span>
                    <span>%s</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Email:</span>
                    <span>%s</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Equipment:</span>
                    <span>%s</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Priority:</span>
                    <span><strong>%s</strong></span>
                </div>
                <div class="info-row">
                    <span class="info-label">Description:</span>
                    <span>%s</span>
                </div>
            </div>

            <div class="action-required">
                <h3>âš ï¸ Action Required</h3>
                <p>Please review and assign an engineer to this ticket.</p>
            </div>
        </div>
        <div class="footer">
            <p>ServQR Admin Notification System</p>
        </div>
    </div>
</body>
</html>
`, data.TicketNumber, data.CustomerName, data.CustomerPhone, data.CustomerEmail, data.EquipmentName, data.Priority, data.Description)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: status %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

// sendCustomerEngineerAssigned sends engineer assignment email to customer
func (s *NotificationService) sendCustomerEngineerAssigned(ctx context.Context, data TicketAssignedData) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	subject := fmt.Sprintf("Engineer Assigned - Ticket %s", data.TicketNumber)
	to := mail.NewEmail(data.CustomerName, data.CustomerEmail)

	plainText := fmt.Sprintf(`
Dear %s,

Great news! An engineer has been assigned to your service ticket.

Ticket: %s
Equipment: %s

Assigned Engineer:
--------------------------------------------------
Name: %s
Phone: %s
Email: %s

The engineer will contact you shortly to schedule the service visit.

Best regards,
ServQR Service Team
`, data.CustomerName, data.TicketNumber, data.EquipmentName, data.EngineerName, data.EngineerPhone, data.EngineerEmail)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #10b981; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 5px 5px; }
        .engineer-info { background-color: #fff; border-left: 4px solid #10b981; padding: 15px; margin: 20px 0; }
        .info-row { padding: 8px 0; border-bottom: 1px solid #e5e7eb; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>ðŸ‘· Engineer Assigned</h1>
        </div>
        <div class="content">
            <p>Dear <strong>%s</strong>,</p>
            <p>Great news! An engineer has been assigned to your service ticket <strong>%s</strong>.</p>
            
            <div class="engineer-info">
                <h3>Assigned Engineer:</h3>
                <div class="info-row">
                    <strong>Name:</strong> %s
                </div>
                <div class="info-row">
                    <strong>Phone:</strong> %s
                </div>
                <div class="info-row">
                    <strong>Email:</strong> %s
                </div>
            </div>

            <p>The engineer will contact you shortly to schedule the service visit.</p>
            <p>Thank you for your patience.</p>
        </div>
        <div class="footer">
            <p>This is an automated message from ServQR Platform</p>
            <p>&copy; 2025 ServQR. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, data.CustomerName, data.TicketNumber, data.EngineerName, data.EngineerPhone, data.EngineerEmail)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: status %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

// sendEngineerAssigned sends assignment email to engineer
func (s *NotificationService) sendEngineerAssigned(ctx context.Context, data TicketAssignedData) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	subject := fmt.Sprintf("New Ticket Assigned - %s", data.TicketNumber)
	to := mail.NewEmail(data.EngineerName, data.EngineerEmail)

	plainText := fmt.Sprintf(`
Dear %s,

A new service ticket has been assigned to you.

Ticket Details:
--------------------------------------------------
Ticket Number: %s
Customer: %s
Phone: %s
Email: %s
Equipment: %s
Priority: %s
Description: %s

Action Required:
Please contact the customer to schedule a service visit.

ServQR Service System
`, data.EngineerName, data.TicketNumber, data.CustomerName, data.CustomerEmail, data.CustomerEmail, data.EquipmentName, data.Priority, data.Description)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f59e0b; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 5px 5px; }
        .ticket-info { background-color: #fff; border-left: 4px solid #f59e0b; padding: 15px; margin: 20px 0; }
        .info-row { padding: 8px 0; border-bottom: 1px solid #e5e7eb; }
        .action-box { background-color: #fffbeb; border: 2px solid #f59e0b; padding: 15px; border-radius: 5px; margin: 20px 0; text-align: center; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>ðŸ”§ New Ticket Assigned</h1>
        </div>
        <div class="content">
            <p>Dear <strong>%s</strong>,</p>
            <p>A new service ticket has been assigned to you.</p>
            
            <div class="ticket-info">
                <div class="info-row">
                    <strong>Ticket:</strong> %s
                </div>
                <div class="info-row">
                    <strong>Customer:</strong> %s
                </div>
                <div class="info-row">
                    <strong>Phone:</strong> %s
                </div>
                <div class="info-row">
                    <strong>Email:</strong> %s
                </div>
                <div class="info-row">
                    <strong>Equipment:</strong> %s
                </div>
                <div class="info-row">
                    <strong>Priority:</strong> <span style="color: #dc2626; font-weight: bold;">%s</span>
                </div>
                <div class="info-row">
                    <strong>Issue:</strong> %s
                </div>
            </div>

            <div class="action-box">
                <h3>ðŸ“ž Action Required</h3>
                <p>Please contact the customer to schedule a service visit.</p>
            </div>
        </div>
        <div class="footer">
            <p>ServQR Service Notification System</p>
        </div>
    </div>
</body>
</html>
`, data.EngineerName, data.TicketNumber, data.CustomerName, data.CustomerEmail, data.CustomerEmail, data.EquipmentName, data.Priority, data.Description)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: status %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

// sendCustomerStatusChanged sends status change email to customer
func (s *NotificationService) sendCustomerStatusChanged(ctx context.Context, data TicketStatusChangedData) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	subject := fmt.Sprintf("Ticket Status Updated - %s", data.TicketNumber)
	to := mail.NewEmail(data.CustomerName, data.CustomerEmail)

	plainText := fmt.Sprintf(`
Dear %s,

Your service ticket status has been updated.

Ticket: %s
Equipment: %s

Status Change:
%s â†’ %s

Updated by: %s

Thank you for your patience.

Best regards,
ServQR Service Team
`, data.CustomerName, data.TicketNumber, data.EquipmentName, data.OldStatus, data.NewStatus, data.UpdatedBy)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #6366f1; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 5px 5px; }
        .status-box { background-color: #fff; padding: 20px; margin: 20px 0; text-align: center; border: 2px solid #6366f1; border-radius: 5px; }
        .status-arrow { font-size: 24px; margin: 0 10px; }
        .old-status { color: #6b7280; }
        .new-status { color: #10b981; font-weight: bold; font-size: 18px; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>ðŸ“Š Status Updated</h1>
        </div>
        <div class="content">
            <p>Dear <strong>%s</strong>,</p>
            <p>Your service ticket <strong>%s</strong> status has been updated.</p>
            
            <div class="status-box">
                <div class="old-status">%s</div>
                <div class="status-arrow">â¬‡ï¸</div>
                <div class="new-status">%s</div>
            </div>

            <p><strong>Equipment:</strong> %s</p>
            <p><strong>Updated by:</strong> %s</p>
            
            <p>Thank you for your patience.</p>
        </div>
        <div class="footer">
            <p>This is an automated message from ServQR Platform</p>
            <p>&copy; 2025 ServQR. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, data.CustomerName, data.TicketNumber, data.OldStatus, data.NewStatus, data.EquipmentName, data.UpdatedBy)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: status %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

// sendAdminStatusChanged sends status change email to admin
func (s *NotificationService) sendAdminStatusChanged(ctx context.Context, data TicketStatusChangedData) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	subject := fmt.Sprintf("Ticket %s - Status: %s", data.TicketNumber, data.NewStatus)
	to := mail.NewEmail("Admin", data.AdminEmail)

	plainText := fmt.Sprintf(`
Ticket Status Updated

Ticket: %s
Customer: %s
Equipment: %s
Status: %s â†’ %s
Updated by: %s

ServQR Admin System
`, data.TicketNumber, data.CustomerName, data.EquipmentName, data.OldStatus, data.NewStatus, data.UpdatedBy)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #6b7280; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 5px 5px; }
        .info-box { background-color: #fff; padding: 15px; margin: 20px 0; border-left: 4px solid #6b7280; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Ticket Status Updated</h1>
        </div>
        <div class="content">
            <div class="info-box">
                <p><strong>Ticket:</strong> %s</p>
                <p><strong>Customer:</strong> %s</p>
                <p><strong>Equipment:</strong> %s</p>
                <p><strong>Status Change:</strong> %s â†’ <strong>%s</strong></p>
                <p><strong>Updated by:</strong> %s</p>
            </div>
        </div>
        <div class="footer">
            <p>ServQR Admin Notification System</p>
        </div>
    </div>
</body>
</html>
`, data.TicketNumber, data.CustomerName, data.EquipmentName, data.OldStatus, data.NewStatus, data.UpdatedBy)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: status %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

// SendInvitationEmail sends invitation email to new team member
func (s *NotificationService) SendInvitationEmail(ctx context.Context, data InvitationData) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	subject := fmt.Sprintf("You've been invited to %s", data.OrganizationName)
	to := mail.NewEmail(data.InviteeName, data.InviteeEmail)

	// Capitalize role for display
	roleDisplay := data.Role
	if len(roleDisplay) > 0 {
		roleDisplay = string(roleDisplay[0]-32) + roleDisplay[1:]
	}

	plainText := fmt.Sprintf(`
Dear %s,

%s has invited you to join %s as a %s.

Click the link below to accept the invitation and create your account:

%s

This invitation expires on %s.

What you'll be able to do:
- Manage equipment and parts catalog
- View and manage service tickets
- Oversee operations and team
- Generate reports and analytics

If you didn't expect this invitation, you can safely ignore this email.

Best regards,
%s Team
`, data.InviteeName, data.InviterName, data.OrganizationName, roleDisplay, data.InviteURL, data.ExpiresAt, data.OrganizationName)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2563eb; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 5px 5px; }
        .invite-box { background-color: #eff6ff; border-left: 4px solid #2563eb; padding: 20px; margin: 20px 0; }
        .info-row { padding: 8px 0; }
        .info-label { font-weight: bold; color: #6b7280; }
        .cta-button { display: inline-block; background-color: #2563eb; color: white; padding: 15px 30px; text-decoration: none; border-radius: 5px; margin: 20px 0; font-weight: bold; }
        .cta-button:hover { background-color: #1d4ed8; }
        .features { background-color: #fff; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .features ul { list-style-type: none; padding: 0; }
        .features li { padding: 8px 0; padding-left: 25px; position: relative; }
        .features li:before { content: "âœ“"; position: absolute; left: 0; color: #10b981; font-weight: bold; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #6b7280; }
        .expiry { color: #dc2626; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>ðŸ‘¥ You're Invited!</h1>
        </div>
        <div class="content">
            <p>Dear <strong>%s</strong>,</p>
            <p><strong>%s</strong> has invited you to join <strong>%s</strong> as a <strong>%s</strong>.</p>
            
            <div class="invite-box">
                <div class="info-row">
                    <span class="info-label">Organization:</span>
                    <span>%s</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Role:</span>
                    <span>%s</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Invited by:</span>
                    <span>%s</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Expires:</span>
                    <span class="expiry">%s</span>
                </div>
            </div>

            <div style="margin: 30px 0; padding: 20px; background-color: #f8f9fa; border-radius: 8px;">
                <p style="margin: 0 0 15px 0; font-size: 14px; color: #666;">
                    Click the link below to validate your account and set up your password:
                </p>
                <div style="text-align: center; margin: 15px 0;">
                    <a href="%s" style="color: #0066cc; font-size: 14px; word-break: break-all;">%s</a>
                </div>
                <p style="margin: 15px 0 0 0; font-size: 12px; color: #999;">
                    This link will expire on %s. If you didn't request this invitation, you can safely ignore this email.
                </p>
            </div>

            <div class="features">
                <h3>What you'll be able to do:</h3>
                <ul>
                    <li>Manage equipment and parts catalog</li>
                    <li>View and manage service tickets</li>
                    <li>Oversee operations and team</li>
                    <li>Generate reports and analytics</li>
                </ul>
            </div>

            <p style="font-size: 12px; color: #6b7280;">If you didn't expect this invitation, you can safely ignore this email.</p>
        </div>
        <div class="footer">
            <p>This is an automated message from ServQR Platform</p>
            <p>&copy; 2026 ServQR. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, data.InviteeName, data.InviterName, data.OrganizationName, roleDisplay, data.OrganizationName, roleDisplay, data.InviterName, data.ExpiresAt, data.InviteURL, data.InviteURL, data.ExpiresAt)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send invitation email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: status %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}
