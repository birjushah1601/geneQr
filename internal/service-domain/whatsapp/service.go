package whatsapp

import (
	"context"
	"fmt"
	"log/slog"
	
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// WhatsAppService handles sending WhatsApp messages via Twilio
type WhatsAppService struct {
	client         *twilio.RestClient
	whatsappNumber string
	logger         *slog.Logger
}

// NewWhatsAppService creates a new WhatsApp service
func NewWhatsAppService(accountSID, authToken, whatsappNumber string, logger *slog.Logger) *WhatsAppService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	
	return &WhatsAppService{
		client:         client,
		whatsappNumber: whatsappNumber,
		logger:         logger.With(slog.String("component", "whatsapp_service")),
	}
}

// SendMessage sends a WhatsApp message to the specified number
func (s *WhatsAppService) SendMessage(ctx context.Context, to, message string) error {
	s.logger.Info("Sending WhatsApp message",
		slog.String("to", maskPhoneNumber(to)),
		slog.Int("message_length", len(message)))
	
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.whatsappNumber)
	params.SetBody(message)
	
	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		s.logger.Error("Failed to send WhatsApp message",
			slog.String("to", maskPhoneNumber(to)),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}
	
	s.logger.Info("WhatsApp message sent successfully",
		slog.String("message_sid", *resp.Sid),
		slog.String("to", maskPhoneNumber(to)))
	
	return nil
}

// SendTicketConfirmation sends a ticket creation confirmation message
func (s *WhatsAppService) SendTicketConfirmation(ctx context.Context, to, ticketNumber string) error {
	message := fmt.Sprintf(
		"✅ *Service Ticket Created Successfully!*\n\n"+
			"📋 Ticket Number: *%s*\n\n"+
			"Our engineer will contact you soon to address the issue.\n\n"+
			"💬 Reply to this message for updates or call our support team.\n\n"+
			"Thank you for using our service!",
		ticketNumber,
	)
	
	return s.SendMessage(ctx, to, message)
}

// SendTicketUpdate sends a ticket status update message
func (s *WhatsAppService) SendTicketUpdate(ctx context.Context, to, ticketNumber, status, message string) error {
	msg := fmt.Sprintf(
		"📋 *Ticket Update: %s*\n\n"+
			"Status: *%s*\n\n"+
			"%s\n\n"+
			"Reply for more information.",
		ticketNumber,
		status,
		message,
	)
	
	return s.SendMessage(ctx, to, msg)
}

// SendEngineerAssignment sends engineer assignment notification
func (s *WhatsAppService) SendEngineerAssignment(ctx context.Context, to, ticketNumber, engineerName, engineerPhone string) error {
	message := fmt.Sprintf(
		"👨‍🔧 *Engineer Assigned!*\n\n"+
			"📋 Ticket: *%s*\n"+
			"👤 Engineer: *%s*\n"+
			"📞 Contact: %s\n\n"+
			"The engineer will reach out to schedule a visit.\n\n"+
			"Thank you for your patience!",
		ticketNumber,
		engineerName,
		engineerPhone,
	)
	
	return s.SendMessage(ctx, to, message)
}

// SendErrorMessage sends an error message to the customer
func (s *WhatsAppService) SendErrorMessage(ctx context.Context, to, errorMsg string) error {
	message := fmt.Sprintf(
		"❌ *Unable to Process Request*\n\n"+
			"%s\n\n"+
			"Please try again or contact our support team for assistance.\n\n"+
			"📞 Support: +1-XXX-XXX-XXXX",
		errorMsg,
	)
	
	return s.SendMessage(ctx, to, message)
}

// SendHelpMessage sends a help message explaining available commands
func (s *WhatsAppService) SendHelpMessage(ctx context.Context, to string) error {
	message := `📱 *Medical Equipment Service - WhatsApp Support*

Welcome! Here's how you can create a service ticket:

1️⃣ *Send Equipment QR Code*
   • Take a photo of the equipment QR code
   • Send it to this number

2️⃣ *Send Equipment ID*
   • Type: EQ-001 (your equipment ID)
   • Describe the issue

3️⃣ *Get Updates*
   • Reply to any message for ticket status
   • We'll notify you of all updates

📞 Need help? Call: +1-XXX-XXX-XXXX

Let's get started! Send your equipment QR code or ID.`
	
	return s.SendMessage(ctx, to, message)
}

// maskPhoneNumber masks a phone number for logging (privacy)
func maskPhoneNumber(phone string) string {
	if len(phone) <= 4 {
		return "****"
	}
	return phone[:len(phone)-4] + "****"
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
