package whatsapp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MediaMessage represents a WhatsApp media message
type MediaMessage struct {
	ID            string    `json:"id"`
	From          string    `json:"from"`
	To            string    `json:"to"`
	MediaURL      string    `json:"media_url"`
	MediaType     string    `json:"media_type"` // image, video, document, audio
	MimeType      string    `json:"mime_type"`  // image/jpeg, video/mp4, etc.
	Filename      string    `json:"filename,omitempty"`
	Caption       string    `json:"caption,omitempty"`
	FileSize      int64     `json:"file_size,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}

// AttachmentInfo represents processed attachment metadata
type AttachmentInfo struct {
	ID               uuid.UUID `json:"id"`
	TicketID         string    `json:"ticket_id"`
	Filename         string    `json:"filename"`
	OriginalFilename string    `json:"original_filename"`
	FileType         string    `json:"file_type"`
	FileSizeBytes    int64     `json:"file_size_bytes"`
	StoragePath      string    `json:"storage_path"`
	Source           string    `json:"source"`
	SourceMessageID  string    `json:"source_message_id"`
	Category         string    `json:"attachment_category"`
}

// MediaHandler handles WhatsApp media processing and storage
type MediaHandler struct {
	logger      *slog.Logger
	storageRoot string
	maxFileSize int64 // Maximum file size in bytes
}

// NewMediaHandler creates a new media handler
func NewMediaHandler(storageRoot string, maxFileSize int64, logger *slog.Logger) *MediaHandler {
	return &MediaHandler{
		logger:      logger.With(slog.String("component", "media_handler")),
		storageRoot: storageRoot,
		maxFileSize: maxFileSize,
	}
}

// ProcessMediaMessage processes a WhatsApp media message and stores the file
func (h *MediaHandler) ProcessMediaMessage(ctx context.Context, ticketID string, msg *MediaMessage) (*AttachmentInfo, error) {
	h.logger.Info("Processing WhatsApp media message",
		slog.String("message_id", msg.ID),
		slog.String("media_type", msg.MediaType),
		slog.String("mime_type", msg.MimeType),
		slog.Int64("file_size", msg.FileSize),
	)

	// Validate file size
	if msg.FileSize > h.maxFileSize {
		return nil, fmt.Errorf("file size %d exceeds maximum allowed size %d", msg.FileSize, h.maxFileSize)
	}

	// Download media file
	localPath, err := h.downloadMediaFile(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to download media file: %w", err)
	}

	// Generate secure filename
	filename := h.generateSecureFilename(msg)
	
	// Determine attachment category based on media type and caption
	category := h.determineAttachmentCategory(msg)

	// Create attachment info
	attachmentInfo := &AttachmentInfo{
		ID:               uuid.New(),
		TicketID:         ticketID,
		Filename:         filename,
		OriginalFilename: msg.Filename,
		FileType:         msg.MimeType,
		FileSizeBytes:    msg.FileSize,
		StoragePath:      localPath,
		Source:           "whatsapp",
		SourceMessageID:  msg.ID,
		Category:         category,
	}

	h.logger.Info("Media processing completed",
		slog.String("attachment_id", attachmentInfo.ID.String()),
		slog.String("filename", filename),
		slog.String("category", category),
		slog.String("storage_path", localPath),
	)

	return attachmentInfo, nil
}

// downloadMediaFile downloads media from WhatsApp and stores locally
func (h *MediaHandler) downloadMediaFile(ctx context.Context, msg *MediaMessage) (string, error) {
	// Create HTTP request to download media
	req, err := http.NewRequestWithContext(ctx, "GET", msg.MediaURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create download request: %w", err)
	}

	// Add authentication headers (WhatsApp Business API requires auth)
	// TODO: Add proper WhatsApp API authentication
	req.Header.Set("Authorization", "Bearer "+os.Getenv("WHATSAPP_ACCESS_TOKEN"))

	// Execute request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download media: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("media download failed with status: %d", resp.StatusCode)
	}

	// Create directory structure
	datePath := time.Now().Format("2006/01/02")
	storageDir := filepath.Join(h.storageRoot, "whatsapp", datePath)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Generate unique filename
	filename := h.generateSecureFilename(msg)
	filePath := filepath.Join(storageDir, filename)

	// Create and write file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy media content with size limit
	written, err := io.CopyN(file, resp.Body, h.maxFileSize+1)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to write media file: %w", err)
	}
	
	if written > h.maxFileSize {
		os.Remove(filePath) // Clean up
		return "", fmt.Errorf("file size exceeds limit")
	}

	h.logger.Info("Media file downloaded successfully",
		slog.String("file_path", filePath),
		slog.Int64("bytes_written", written),
	)

	return filePath, nil
}

// generateSecureFilename generates a secure filename for the media file
func (h *MediaHandler) generateSecureFilename(msg *MediaMessage) string {
	// Generate random hex string
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)

	// Get file extension from mime type or original filename
	var ext string
	if msg.Filename != "" {
		ext = filepath.Ext(msg.Filename)
	} else {
		exts, _ := mime.ExtensionsByType(msg.MimeType)
		if len(exts) > 0 {
			ext = exts[0]
		}
	}

	// Create secure filename: timestamp_mediatype_random.ext
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s_%s%s", timestamp, msg.MediaType, randomHex, ext)
	
	return filename
}

// determineAttachmentCategory determines the category based on media type and caption
func (h *MediaHandler) determineAttachmentCategory(msg *MediaMessage) string {
	caption := strings.ToLower(msg.Caption)
	
	switch msg.MediaType {
	case "image":
		// Analyze caption for specific categories
		if strings.Contains(caption, "before") || strings.Contains(caption, "पहले") {
			return "equipment_photo"
		}
		if strings.Contains(caption, "after") || strings.Contains(caption, "बाद") {
			return "repair_photo"
		}
		if strings.Contains(caption, "error") || strings.Contains(caption, "issue") || 
		   strings.Contains(caption, "problem") || strings.Contains(caption, "गलती") {
			return "issue_photo"
		}
		return "equipment_photo" // Default for images
		
	case "video":
		return "video"
		
	case "document":
		return "document"
		
	case "audio":
		return "audio"
		
	default:
		return "other"
	}
}

// Enhanced WhatsApp message structure to include media
type EnhancedWhatsAppMessage struct {
	ID        string         `json:"id"`
	From      string         `json:"from"`
	To        string         `json:"to"`
	Type      string         `json:"type"` // text, image, video, document, audio
	Text      string         `json:"text,omitempty"`
	Media     *MediaMessage  `json:"media,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// ProcessEnhancedMessage processes both text and media messages
// TODO: Implement after WhatsAppHandler is defined
/*func (h *WhatsAppHandler) ProcessEnhancedMessage(ctx context.Context, msg *EnhancedWhatsAppMessage) error {
	h.logger.Info("Processing enhanced WhatsApp message",
		slog.String("from", msg.From),
		slog.String("type", msg.Type),
		slog.String("id", msg.ID),
	)

	// Extract QR code from text or media caption
	var qrCode string
	if msg.Type == "text" && msg.Text != "" {
		qrCode = h.extractQRCode(msg.Text)
	} else if msg.Media != nil && msg.Media.Caption != "" {
		qrCode = h.extractQRCode(msg.Media.Caption)
	}

	if qrCode == "" {
		h.logger.Warn("No QR code found in message", slog.String("message_id", msg.ID))
		h.sendHelpMessage(ctx, msg.From)
		return nil
	}

	// Lookup equipment
	equipment, err := h.equipmentService.GetEquipmentByQR(ctx, qrCode)
	if err != nil {
		h.logger.Error("Equipment not found", slog.String("qr_code", qrCode), slog.String("error", err.Error()))
		h.sendErrorMessage(ctx, msg.From, "Equipment not found. Please check the QR code and try again.")
		return err
	}

	// Extract issue description
	var issueDescription string
	if msg.Type == "text" {
		issueDescription = h.extractIssueDescription(msg.Text, qrCode)
	} else if msg.Media != nil {
		if msg.Media.Caption != "" {
			issueDescription = h.extractIssueDescription(msg.Media.Caption, qrCode)
		} else {
			issueDescription = fmt.Sprintf("%s attachment received for equipment issue", msg.Media.MediaType)
		}
	}

	// Determine priority
	priority := h.determinePriority(issueDescription)

	// Create service ticket
	ticket, err := h.createEnhancedTicketFromWhatsApp(ctx, equipment, msg, issueDescription, priority)
	if err != nil {
		h.logger.Error("Failed to create ticket", slog.String("error", err.Error()))
		h.sendErrorMessage(ctx, msg.From, "Failed to create service ticket. Please try again.")
		return err
	}

	h.logger.Info("Ticket created with media support",
		slog.String("ticket_number", ticket.TicketNumber),
		slog.String("has_media", fmt.Sprintf("%t", msg.Media != nil)),
	)

	// Process media if present
	if msg.Media != nil {
		go h.processMediaAttachment(context.Background(), ticket.ID, msg.Media)
	}

	// Send confirmation
	h.sendTicketConfirmation(ctx, msg.From, ticket)
	
	return nil
}

// processMediaAttachment processes media attachment asynchronously
func (h *WhatsAppHandler) processMediaAttachment(ctx context.Context, ticketID string, media *MediaMessage) {
	mediaHandler := NewMediaHandler(
		os.Getenv("STORAGE_ROOT"),
		50*1024*1024, // 50MB max file size
		h.logger,
	)

	// Process and store media
	attachmentInfo, err := mediaHandler.ProcessMediaMessage(ctx, ticketID, media)
	if err != nil {
		h.logger.Error("Failed to process media attachment",
			slog.String("ticket_id", ticketID),
			slog.String("media_id", media.ID),
			slog.String("error", err.Error()),
		)
		return
	}

	// Store attachment in database
	if err := h.storeAttachmentInDB(ctx, attachmentInfo); err != nil {
		h.logger.Error("Failed to store attachment in database",
			slog.String("attachment_id", attachmentInfo.ID.String()),
			slog.String("error", err.Error()),
		)
		return
	}

	// Queue for AI analysis if it's an image
	if strings.HasPrefix(attachmentInfo.FileType, "image/") {
		if err := h.queueForAIAnalysis(ctx, attachmentInfo.ID); err != nil {
			h.logger.Error("Failed to queue for AI analysis",
				slog.String("attachment_id", attachmentInfo.ID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			h.logger.Info("Attachment queued for AI analysis",
				slog.String("attachment_id", attachmentInfo.ID.String()),
				slog.String("file_type", attachmentInfo.FileType),
			)
		}
	}
}

// storeAttachmentInDB stores attachment metadata in database
func (h *WhatsAppHandler) storeAttachmentInDB(ctx context.Context, attachment *AttachmentInfo) error {
	// TODO: Implement database storage using repository pattern
	h.logger.Info("Storing attachment in database",
		slog.String("attachment_id", attachment.ID.String()),
		slog.String("ticket_id", attachment.TicketID),
		slog.String("filename", attachment.Filename),
	)
	return nil
}

// queueForAIAnalysis queues attachment for AI vision analysis
func (h *WhatsAppHandler) queueForAIAnalysis(ctx context.Context, attachmentID uuid.UUID) error {
	// TODO: Implement queue system for AI analysis
	h.logger.Info("Queuing attachment for AI analysis",
		slog.String("attachment_id", attachmentID.String()),
	)
	return nil
}

// Enhanced ticket creation to handle media attachments
func (h *WhatsAppHandler) createEnhancedTicketFromWhatsApp(
	ctx context.Context,
	equipment interface{}, // domain.Equipment
	msg *EnhancedWhatsAppMessage,
	issueDescription string,
	priority interface{}, // ticketDomain.TicketPriority
) (interface{}, error) { // Returns *ticketDomain.ServiceTicket
	// This would create ticket with awareness of attached media
	// The actual implementation would use your existing ticket creation logic
	// but with additional metadata about attachments
	
	h.logger.Info("Creating enhanced ticket with media support",
		slog.String("message_type", msg.Type),
		slog.Bool("has_media", msg.Media != nil),
		slog.String("description", issueDescription),
	)
	
	// TODO: Implement actual ticket creation with media awareness
	return nil, fmt.Errorf("enhanced ticket creation not yet implemented")
}*/