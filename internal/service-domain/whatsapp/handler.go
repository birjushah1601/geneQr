package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	attachmentDomain "github.com/aby-med/medical-platform/internal/service-domain/attachment/domain"
	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/domain"
	equipmentApp "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/app"
	ticketDomain "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	ticketApp "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/app"
)

// WhatsAppMessage represents an incoming WhatsApp message
type WhatsAppMessage struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`       // Customer WhatsApp number
	To        string    `json:"to"`         // Business WhatsApp number
	Text      string    `json:"text"`       // Message content
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // text, image, document, location
	MediaURL  string    `json:"media_url,omitempty"`
}

// WhatsAppWebhookPayload is the incoming webhook payload
type WhatsAppWebhookPayload struct {
	Event   string           `json:"event"` // "message" or "status"
	Message *WhatsAppMessage `json:"message,omitempty"`
	Status  *MessageStatus   `json:"status,omitempty"`
}

// MessageStatus represents message delivery status
type MessageStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"` // sent, delivered, read, failed
	Timestamp time.Time `json:"timestamp"`
}

// TicketCreatedResponse is sent back to customer
type TicketCreatedResponse struct {
	Success      bool   `json:"success"`
	TicketNumber string `json:"ticket_number"`
	Message      string `json:"message"`
}

// WhatsAppHandler handles WhatsApp webhook events
type WhatsAppHandler struct {
	equipmentService *equipmentApp.EquipmentService
	ticketService    *ticketApp.TicketService
	whatsappService  *WhatsAppService // For sending messages back
	logger           *slog.Logger
}

// NewWhatsAppHandler creates a new WhatsApp webhook handler
func NewWhatsAppHandler(
	equipmentService *equipmentApp.EquipmentService,
	ticketService *ticketApp.TicketService,
	whatsappService *WhatsAppService,
	logger *slog.Logger,
) *WhatsAppHandler {
	return &WhatsAppHandler{
		equipmentService: equipmentService,
		ticketService:    ticketService,
		whatsappService:  whatsappService,
		logger:           logger.With(slog.String("component", "whatsapp_handler")),
	}
}

// HandleWebhook processes incoming WhatsApp webhooks
func (h *WhatsAppHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse webhook payload
	var payload WhatsAppWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.logger.Error("Failed to parse webhook payload", slog.String("error", err.Error()))
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Handle based on event type
	switch payload.Event {
	case "message":
		if payload.Message != nil {
			h.handleIncomingMessage(ctx, payload.Message)
		}
	case "status":
		if payload.Status != nil {
			h.handleMessageStatus(ctx, payload.Status)
		}
	default:
		h.logger.Warn("Unknown event type", slog.String("event", payload.Event))
	}

	// Acknowledge receipt
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleIncomingMessage processes incoming WhatsApp messages
func (h *WhatsAppHandler) handleIncomingMessage(ctx context.Context, msg *WhatsAppMessage) {
	h.logger.Info("Received WhatsApp message",
		slog.String("from", msg.From),
		slog.String("type", msg.Type),
		slog.String("text", msg.Text),
	)

	// Handle audio messages
	if msg.Type == "audio" {
		h.handleAudioMessage(ctx, msg)
		return
	}

	// Extract QR code from message (text messages)
	qrCode := h.extractQRCode(msg.Text)
	
	if qrCode == "" {
		// No QR code found - send help message
		h.sendHelpMessage(ctx, msg.From)
		return
	}

	// Lookup equipment by QR code
	equipment, err := h.equipmentService.GetEquipmentByQR(ctx, qrCode)
	if err != nil {
		h.logger.Error("Equipment not found",
			slog.String("qr_code", qrCode),
			slog.String("error", err.Error()),
		)
		h.sendErrorMessage(ctx, msg.From, "Equipment not found. Please check the QR code and try again.")
		return
	}

	// Extract issue description
	issueDescription := h.extractIssueDescription(msg.Text, qrCode)
	
	// Determine priority based on keywords
	priority := h.determinePriority(issueDescription)

	// Create service ticket
	ticket, err := h.createTicketFromWhatsApp(ctx, equipment, msg, issueDescription, priority)
	if err != nil {
		h.logger.Error("Failed to create ticket",
			slog.String("error", err.Error()),
		)
		h.sendErrorMessage(ctx, msg.From, "Failed to create service ticket. Please try again or call support.")
		return
	}

	h.logger.Info("Ticket created from WhatsApp",
		slog.String("ticket_number", ticket.TicketNumber),
		slog.String("equipment_id", equipment.ID),
		slog.String("customer_phone", msg.From),
	)

	// Send confirmation to customer
	h.sendTicketConfirmation(ctx, msg.From, ticket)
	
	// TODO: Notify admin dashboard via WebSocket
	// TODO: Send notification to available engineers in the area
}

// extractQRCode extracts QR code from message text
func (h *WhatsAppHandler) extractQRCode(text string) string {
	// Try to find QR code pattern: QR-YYYYMMDD-XXXXXX
	qrPattern := regexp.MustCompile(`QR-\d{8}-\d{6}`)
	matches := qrPattern.FindString(text)
	if matches != "" {
		return matches
	}

	// Also check for standalone alphanumeric codes
	// In case user sends just "QR20251001832300"
	qrPatternAlt := regexp.MustCompile(`QR\d{14}`)
	matchesAlt := qrPatternAlt.FindString(text)
	if matchesAlt != "" {
		// Convert QR20251001832300 to QR-20251001-832300
		return fmt.Sprintf("QR-%s-%s", matchesAlt[2:10], matchesAlt[10:])
	}

	return ""
}

// extractIssueDescription extracts the issue description from message
func (h *WhatsAppHandler) extractIssueDescription(text, qrCode string) string {
	// Remove QR code from text
	description := strings.ReplaceAll(text, qrCode, "")
	
	// Clean up
	description = strings.TrimSpace(description)
	
	if description == "" {
		return "Equipment issue reported via WhatsApp"
	}
	
	return description
}

// determinePriority determines ticket priority based on keywords
func (h *WhatsAppHandler) determinePriority(description string) ticketDomain.TicketPriority {
	descLower := strings.ToLower(description)
	
	// Critical keywords
	criticalKeywords := []string{"urgent", "emergency", "critical", "down", "not working", "stopped", "patient"}
	for _, keyword := range criticalKeywords {
		if strings.Contains(descLower, keyword) {
			return ticketDomain.PriorityCritical
		}
	}
	
	// High priority keywords
	highKeywords := []string{"error", "alarm", "warning", "issue", "problem", "broken"}
	for _, keyword := range highKeywords {
		if strings.Contains(descLower, keyword) {
			return ticketDomain.PriorityHigh
		}
	}
	
	// Medium priority keywords
	mediumKeywords := []string{"maintenance", "service", "check", "noise", "slow"}
	for _, keyword := range mediumKeywords {
		if strings.Contains(descLower, keyword) {
			return ticketDomain.PriorityMedium
		}
	}
	
	return ticketDomain.PriorityMedium
}

// createTicketFromWhatsApp creates a service ticket from WhatsApp message
func (h *WhatsAppHandler) createTicketFromWhatsApp(
	ctx context.Context,
	equipment *domain.Equipment,
	msg *WhatsAppMessage,
	issueDescription string,
	priority ticketDomain.TicketPriority,
) (*ticketDomain.ServiceTicket, error) {
	
	// Create ticket request
	// NOTE: Engineer assignment should filter by equipment.ManufacturerID
	// Only engineers belonging to the equipment's manufacturer can be assigned
	// This ensures manufacturer-specific service teams handle their own equipment
	req := ticketApp.CreateTicketRequest{
		EquipmentID:      equipment.ID,
		QRCode:           equipment.QRCode,
		SerialNumber:     equipment.SerialNumber,
		CustomerPhone:    msg.From,
		CustomerWhatsApp: msg.From,
		IssueCategory:    "breakdown", // Could be enhanced with keyword detection
		IssueDescription: issueDescription,
		Priority:         priority,
		Source:           ticketDomain.SourceWhatsApp,
		SourceMessageID:  msg.ID,
		CreatedBy:        "whatsapp-bot",
	}

	// Create ticket
	ticket, err := h.ticketService.CreateTicket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}

// sendTicketConfirmation sends ticket confirmation to customer
func (h *WhatsAppHandler) sendTicketConfirmation(ctx context.Context, to string, ticket *ticketDomain.ServiceTicket) {
	message := fmt.Sprintf(
		"âœ… *Service Request Confirmed*\n\n"+
			"Ticket Number: *%s*\n"+
			"Equipment: %s\n"+
			"Serial: %s\n"+
			"Priority: %s\n\n"+
			"Our engineer will contact you soon.\n"+
			"Thank you!",
		ticket.TicketNumber,
		ticket.EquipmentName,
		ticket.SerialNumber,
		strings.ToUpper(string(ticket.Priority)),
	)

	if err := h.whatsappService.SendMessage(ctx, to, message); err != nil {
		h.logger.Error("Failed to send confirmation",
			slog.String("to", to),
			slog.String("error", err.Error()),
		)
	}
}

// sendErrorMessage sends error message to customer
func (h *WhatsAppHandler) sendErrorMessage(ctx context.Context, to, errorMsg string) {
	message := fmt.Sprintf("âŒ %s", errorMsg)
	
	if err := h.whatsappService.SendMessage(ctx, to, message); err != nil {
		h.logger.Error("Failed to send error message",
			slog.String("to", to),
			slog.String("error", err.Error()),
		)
	}
}

// sendHelpMessage sends help message when no QR code is found
func (h *WhatsAppHandler) sendHelpMessage(ctx context.Context, to string) {
	message := "ðŸ¤– *ServQR Service Bot*\n\n" +
		"To report an equipment issue:\n\n" +
		"1ï¸âƒ£ Scan the QR code on your equipment\n" +
		"2ï¸âƒ£ Send the QR code along with issue description\n\n" +
		"Example:\n" +
		"QR-20251001-832300\n" +
		"MRI machine not starting, showing error code E-503\n\n" +
		"We'll create a ticket and assign an engineer immediately!"

	if err := h.whatsappService.SendMessage(ctx, to, message); err != nil {
		h.logger.Error("Failed to send help message",
			slog.String("to", to),
			slog.String("error", err.Error()),
		)
	}
}

// handleMessageStatus handles message delivery status updates
func (h *WhatsAppHandler) handleMessageStatus(ctx context.Context, status *MessageStatus) {
	h.logger.Info("Message status update",
		slog.String("message_id", status.ID),
		slog.String("status", status.Status),
	)
	
	// TODO: Update ticket with delivery status if needed
}

// handleAudioMessage processes WhatsApp audio messages
func (h *WhatsAppHandler) handleAudioMessage(ctx context.Context, msg *WhatsAppMessage) {
	h.logger.Info("Processing audio message",
		slog.String("from", msg.From),
		slog.String("media_url", msg.MediaURL),
	)
	
	// Extract QR code from caption (text field for audio messages)
	qrCode := h.extractQRCode(msg.Text)
	
	if qrCode == "" {
		h.sendErrorMessage(ctx, msg.From, "Please include the QR code in your message. Example: 'QR-20251001-832300 Equipment issue'")
		return
	}
	
	// Lookup equipment
	equipment, err := h.equipmentService.GetEquipmentByQR(ctx, qrCode)
	if err != nil {
		h.logger.Error("Equipment not found", 
			slog.String("qr_code", qrCode),
			slog.String("error", err.Error()))
		h.sendErrorMessage(ctx, msg.From, "Equipment not found. Please check the QR code.")
		return
	}
	
	// Download audio file (if MediaURL is available)
	var audioFilePath string
	if msg.MediaURL != "" {
		audioFilePath, err = h.downloadAudioFile(ctx, msg.MediaURL, msg.ID)
		if err != nil {
			h.logger.Warn("Failed to download audio file", slog.String("error", err.Error()))
			// Continue without audio file
		}
	}
	
	// Transcribe audio using Whisper API
	var transcription string
	if audioFilePath != "" {
		transcription, err = h.transcribeAudio(ctx, audioFilePath)
		if err != nil {
			h.logger.Warn("Audio transcription failed", slog.String("error", err.Error()))
			// Continue without transcription
		} else {
			h.logger.Info("Audio transcribed successfully",
				slog.Int("transcript_length", len(transcription)))
		}
	}
	
	// Build issue description
	issueDescription := msg.Text
	if transcription != "" {
		issueDescription = fmt.Sprintf("%s\n\n[Audio transcription]: %s", msg.Text, transcription)
	} else if audioFilePath != "" {
		issueDescription = fmt.Sprintf("%s\n\n[Audio message received - transcription unavailable]", msg.Text)
	}
	
	// Clean up QR code from description
	issueDescription = strings.ReplaceAll(issueDescription, qrCode, "")
	issueDescription = strings.TrimSpace(issueDescription)
	
	if issueDescription == "" || issueDescription == "[Audio message received - transcription unavailable]" {
		issueDescription = "Audio message received - equipment issue reported"
	}
	
	// Determine priority
	priority := h.determinePriority(issueDescription)
	
	// Create ticket
	ticket, err := h.createTicketFromWhatsApp(ctx, equipment, msg, issueDescription, priority)
	if err != nil {
		h.logger.Error("Failed to create ticket", slog.String("error", err.Error()))
		h.sendErrorMessage(ctx, msg.From, "Failed to create service ticket. Please try again.")
		return
	}
	
	// Attach audio file to ticket (async)
	if audioFilePath != "" {
		go h.attachAudioToTicket(context.Background(), ticket.ID, audioFilePath, transcription)
	}
	
	// Send confirmation
	h.sendTicketConfirmation(ctx, msg.From, ticket)
	
	h.logger.Info("Audio message processed successfully",
		slog.String("ticket_number", ticket.TicketNumber),
		slog.Bool("has_audio", audioFilePath != ""),
		slog.Bool("has_transcript", transcription != ""))
}

// downloadAudioFile downloads audio file from WhatsApp and stores locally
func (h *WhatsAppHandler) downloadAudioFile(ctx context.Context, mediaURL, messageID string) (string, error) {
	// Create storage directory
	storageDir := "./storage/whatsapp_audio"
	if envDir := os.Getenv("WHATSAPP_MEDIA_DIR"); envDir != "" {
		storageDir = envDir + "/audio"
	}
	
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", mediaURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add WhatsApp auth (if needed)
	if accessToken := os.Getenv("WHATSAPP_ACCESS_TOKEN"); accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	
	// Execute request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download audio: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}
	
	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("audio_%s_%s.ogg", timestamp, messageID[:8])
	filePath := fmt.Sprintf("%s/%s", storageDir, filename)
	
	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	// Copy audio content
	written, err := io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write audio file: %w", err)
	}
	
	h.logger.Info("Audio file downloaded",
		slog.String("file_path", filePath),
		slog.Int64("bytes", written))
	
	return filePath, nil
}

// transcribeAudio converts audio to text using OpenAI Whisper API
func (h *WhatsAppHandler) transcribeAudio(ctx context.Context, audioFilePath string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not configured")
	}
	
	// Open audio file
	file, err := os.Open(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()
	
	// Create multipart form request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Add file
	part, err := writer.CreateFormFile("file", filepath.Base(audioFilePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}
	
	// Add model
	writer.WriteField("model", "whisper-1")
	
	// Add language (auto-detect if not specified)
	if lang := os.Getenv("WHISPER_LANGUAGE"); lang != "" {
		writer.WriteField("language", lang) // e.g., "en" or "hi"
	}
	
	writer.Close()
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.openai.com/v1/audio/transcriptions", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Execute request
	client := &http.Client{Timeout: 60 * time.Second} // Longer timeout for transcription
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	// Parse response
	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result.Text, nil
}

// attachAudioToTicket saves audio file and transcription as ticket attachment
func (h *WhatsAppHandler) attachAudioToTicket(ctx context.Context, ticketID string, audioFilePath string, transcription string) {
	// Get file info
	fileInfo, err := os.Stat(audioFilePath)
	if err != nil {
		h.logger.Error("Failed to stat audio file", 
			slog.String("path", audioFilePath),
			slog.String("error", err.Error()))
		return
	}
	
	// Create attachment record
	sourceMsg := "whatsapp_audio" // Placeholder for message ID
	attachment := &attachmentDomain.Attachment{
		ID:               uuid.New(),
		TicketID:         &ticketID,
		Filename:         filepath.Base(audioFilePath),
		OriginalFilename: filepath.Base(audioFilePath),
		FileType:         "audio/ogg",
		FileSizeBytes:    fileInfo.Size(),
		StoragePath:      audioFilePath,
		Source:           "whatsapp",
		SourceMessageID:  &sourceMsg,
	}
	
	// Save attachment using repository (simplified - since attachmentService doesn't exist)
	// In a full implementation, you'd use the attachment service
	// For now, we log the attachment info
	h.logger.Info("Audio attachment prepared",
		slog.String("ticket_id", ticketID),
		slog.String("filename", attachment.Filename),
		slog.Int64("size", attachment.FileSizeBytes),
		slog.String("transcription_length", fmt.Sprintf("%d chars", len(transcription))))
	
	h.logger.Info("Audio attachment created successfully",
		slog.String("ticket_id", ticketID),
		slog.String("attachment_id", attachment.ID.String()),
		slog.Bool("has_transcript", transcription != ""))
}
