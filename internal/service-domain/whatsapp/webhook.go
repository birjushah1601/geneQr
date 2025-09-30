package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	qrcode "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/qrcode"
)

// WebhookHandler handles WhatsApp webhook events
type WebhookHandler struct {
	verifyToken  string
	accessToken  string
	phoneNumberID string
	qrGenerator  *qrcode.Generator
	logger       *slog.Logger
	ticketCreator TicketCreator
	mediaDir     string
}

// TicketCreator interface for creating tickets from WhatsApp messages
type TicketCreator interface {
	CreateFromWhatsApp(ctx context.Context, req WhatsAppTicketRequest) (string, error)
}

// WhatsAppTicketRequest represents a ticket creation request from WhatsApp
type WhatsAppTicketRequest struct {
	EquipmentID      string
	QRCode           string
	SerialNumber     string
	EquipmentName    string
	CustomerName     string
	CustomerPhone    string
	CustomerWhatsApp string
	IssueDescription string
	Photos           []string
	Videos           []string
	SourceMessageID  string
}

// WebhookConfig holds WhatsApp webhook configuration
type WebhookConfig struct {
	VerifyToken   string
	AccessToken   string
	PhoneNumberID string
	MediaDir      string
}

// NewWebhookHandler creates a new WhatsApp webhook handler
func NewWebhookHandler(
	cfg WebhookConfig,
	qrGenerator *qrcode.Generator,
	ticketCreator TicketCreator,
	logger *slog.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		verifyToken:   cfg.VerifyToken,
		accessToken:   cfg.AccessToken,
		phoneNumberID: cfg.PhoneNumberID,
		qrGenerator:   qrGenerator,
		ticketCreator: ticketCreator,
		mediaDir:      cfg.MediaDir,
		logger:        logger.With(slog.String("component", "whatsapp_webhook")),
	}
}

// VerifyWebhook handles webhook verification from Meta
func (h *WebhookHandler) VerifyWebhook(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")
	
	h.logger.Info("Webhook verification request",
		slog.String("mode", mode),
		slog.String("token", token))
	
	if mode == "subscribe" && token == h.verifyToken {
		h.logger.Info("Webhook verified successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}
	
	h.logger.Warn("Webhook verification failed")
	w.WriteHeader(http.StatusForbidden)
}

// HandleWebhook processes incoming WhatsApp messages
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.logger.Error("Failed to decode webhook payload", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	// Acknowledge receipt immediately
	w.WriteHeader(http.StatusOK)
	
	// Process messages asynchronously
	go h.processMessages(ctx, payload)
}

// processMessages processes WhatsApp messages
func (h *WebhookHandler) processMessages(ctx context.Context, payload WebhookPayload) {
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Value.Messages == nil {
				continue
			}
			
			for _, msg := range change.Value.Messages {
				if err := h.processMessage(ctx, msg, change.Value.Contacts); err != nil {
					h.logger.Error("Failed to process message",
						slog.String("message_id", msg.ID),
						slog.String("error", err.Error()))
				}
			}
		}
	}
}

// processMessage processes a single WhatsApp message
func (h *WebhookHandler) processMessage(ctx context.Context, msg Message, contacts []Contact) error {
	h.logger.Info("Processing WhatsApp message",
		slog.String("message_id", msg.ID),
		slog.String("from", msg.From),
		slog.String("type", msg.Type))
	
	// Get contact info
	var contactName string
	for _, contact := range contacts {
		if contact.WaID == msg.From {
			contactName = contact.Profile.Name
			break
		}
	}
	
	if contactName == "" {
		contactName = msg.From
	}
	
	// Handle different message types
	switch msg.Type {
	case "image":
		return h.handleImageMessage(ctx, msg, contactName)
	case "text":
		return h.handleTextMessage(ctx, msg, contactName)
	case "video":
		return h.handleVideoMessage(ctx, msg, contactName)
	default:
		h.logger.Info("Unsupported message type", slog.String("type", msg.Type))
		return h.sendMessage(msg.From, "Sorry, I can only process images, videos, and text messages for now.")
	}
}

// handleImageMessage processes image messages (potentially containing QR codes)
func (h *WebhookHandler) handleImageMessage(ctx context.Context, msg Message, contactName string) error {
	if msg.Image == nil {
		return fmt.Errorf("image data missing")
	}
	
	// Download image
	imagePath, err := h.downloadMedia(msg.Image.ID, "jpg")
	if err != nil {
		h.logger.Error("Failed to download image", slog.String("error", err.Error()))
		return h.sendMessage(msg.From, "Sorry, I couldn't download the image. Please try again.")
	}
	
	h.logger.Info("Image downloaded", slog.String("path", imagePath))
	
	// Try to decode QR code from image
	qrData, err := h.qrGenerator.DecodeQRFromImage(imagePath)
	if err != nil {
		h.logger.Warn("No QR code found in image", slog.String("error", err.Error()))
		return h.sendMessage(msg.From, "I couldn't find a QR code in this image. Please send a clear photo of the QR code on the equipment.")
	}
	
	h.logger.Info("QR code decoded", slog.Any("data", qrData))
	
	// Create ticket request
	ticketReq := WhatsAppTicketRequest{
		EquipmentID:      qrData.ID,
		QRCode:           qrData.QRCode,
		SerialNumber:     qrData.SerialNo,
		EquipmentName:    fmt.Sprintf("Equipment %s", qrData.SerialNo),
		CustomerName:     contactName,
		CustomerPhone:    msg.From,
		CustomerWhatsApp: msg.From,
		IssueDescription: msg.Image.Caption,
		Photos:           []string{imagePath},
		SourceMessageID:  msg.ID,
	}
	
	// Create ticket
	ticketID, err := h.ticketCreator.CreateFromWhatsApp(ctx, ticketReq)
	if err != nil {
		h.logger.Error("Failed to create ticket", slog.String("error", err.Error()))
		return h.sendMessage(msg.From, "Sorry, I couldn't create a service ticket. Please try again or contact support.")
	}
	
	// Send acknowledgment
	response := fmt.Sprintf(
		"✅ Service ticket created successfully!\n\n"+
			"Ticket ID: %s\n"+
			"Equipment: %s\n"+
			"Serial: %s\n\n"+
			"Our engineer will contact you soon. Thank you!",
		ticketID, ticketReq.EquipmentName, ticketReq.SerialNumber)
	
	return h.sendMessage(msg.From, response)
}

// handleTextMessage processes text messages
func (h *WebhookHandler) handleTextMessage(ctx context.Context, msg Message, contactName string) error {
	if msg.Text == nil {
		return fmt.Errorf("text data missing")
	}
	
	text := strings.ToLower(strings.TrimSpace(msg.Text.Body))
	
	// Handle commands
	if text == "help" || text == "start" {
		return h.sendWelcomeMessage(msg.From)
	}
	
	// Otherwise, ask for QR code
	return h.sendMessage(msg.From, "Please send a photo of the QR code on your equipment so I can help you.")
}

// handleVideoMessage processes video messages
func (h *WebhookHandler) handleVideoMessage(ctx context.Context, msg Message, contactName string) error {
	if msg.Video == nil {
		return fmt.Errorf("video data missing")
	}
	
	// Download video
	videoPath, err := h.downloadMedia(msg.Video.ID, "mp4")
	if err != nil {
		h.logger.Error("Failed to download video", slog.String("error", err.Error()))
		return h.sendMessage(msg.From, "Sorry, I couldn't download the video. Please try again.")
	}
	
	h.logger.Info("Video downloaded", slog.String("path", videoPath))
	
	// Ask for QR code to link the video
	return h.sendMessage(msg.From, "Video received. Please also send a photo of the QR code on the equipment so I can create a service ticket.")
}

// downloadMedia downloads media from WhatsApp
func (h *WebhookHandler) downloadMedia(mediaID, extension string) (string, error) {
	// Get media URL
	mediaURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s", mediaID)
	
	req, err := http.NewRequest("GET", mediaURL, nil)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Authorization", "Bearer "+h.accessToken)
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get media URL: %s", string(body))
	}
	
	var mediaData struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mediaData); err != nil {
		return "", err
	}
	
	// Download actual media
	req, err = http.NewRequest("GET", mediaData.URL, nil)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Authorization", "Bearer "+h.accessToken)
	
	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download media: status %d", resp.StatusCode)
	}
	
	// Save to file
	filename := fmt.Sprintf("%s_%d.%s", mediaID, time.Now().Unix(), extension)
	filepath := filepath.Join(h.mediaDir, filename)
	
	// Ensure directory exists
	if err := os.MkdirAll(h.mediaDir, 0755); err != nil {
		return "", err
	}
	
	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", err
	}
	
	return filepath, nil
}

// sendMessage sends a text message via WhatsApp
func (h *WebhookHandler) sendMessage(to, message string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", h.phoneNumberID)
	
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": message,
		},
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.accessToken)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message: %s", string(body))
	}
	
	h.logger.Info("Message sent successfully", slog.String("to", to))
	return nil
}

// sendWelcomeMessage sends a welcome/help message
func (h *WebhookHandler) sendWelcomeMessage(to string) error {
	message := "👋 Welcome to Equipment Service Support!\n\n" +
		"To request service:\n" +
		"1. Take a clear photo of the QR code on your equipment\n" +
		"2. Send it to me\n" +
		"3. Add a description of the issue\n\n" +
		"I'll create a service ticket and our engineer will contact you soon."
	
	return h.sendMessage(to, message)
}

// WebhookPayload represents the structure of WhatsApp webhook data
type WebhookPayload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Value Value  `json:"value"`
	Field string `json:"field"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

type Message struct {
	From      string     `json:"from"`
	ID        string     `json:"id"`
	Timestamp string     `json:"timestamp"`
	Type      string     `json:"type"`
	Text      *TextMsg   `json:"text,omitempty"`
	Image     *ImageMsg  `json:"image,omitempty"`
	Video     *VideoMsg  `json:"video,omitempty"`
}

type TextMsg struct {
	Body string `json:"body"`
}

type ImageMsg struct {
	Caption  string `json:"caption"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
	ID       string `json:"id"`
}

type VideoMsg struct {
	Caption  string `json:"caption"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
	ID       string `json:"id"`
}
