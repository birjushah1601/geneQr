# Ticket Creation Enhancements - Implementation Guide

## 📋 Overview

Complete implementation guide for ticket creation enhancements:
1. **Priority field** - Remove from creation, admin-only updates
2. **Audio file uploads** - Support audio in ticket creation  
3. **WhatsApp audio messages** - Attach file and/or convert to text

## ✅ Phase 1: Frontend Changes (COMPLETE)

### Changes Made
- ✅ Removed priority selection field from ticket creation form
- ✅ Set default priority to 'medium' for all new tickets
- ✅ Added audio/* to file upload accept types
- ✅ Updated UI labels to include "Audio"
- ✅ Updated form state to remove priority

### Files Modified
- `admin-ui/src/app/service-request/page.tsx` (Phase 1 complete)

### Testing
```bash
# Manual testing:
1. Navigate to /service-request?qr=TEST-QR-001
2. Verify priority field is NOT visible
3. Try uploading an audio file (.mp3, .wav, .m4a)
4. Verify file is accepted and listed
5. Submit ticket and verify it's created with priority='medium'
```

## 🔧 Phase 2: Backend - Admin-Only Priority Updates

### Implementation Steps

#### 1. Create Priority Update API Endpoint

**File:** `internal/service-domain/service-ticket/api/handler.go`

Add new handler method:

```go
// UpdateTicketPriority handles PATCH /tickets/{id}/priority (admin-only)
func (h *TicketHandler) UpdateTicketPriority(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "id")
	
	// Extract user role from context (set by auth middleware)
	userRole := ctx.Value("user_role").(string)
	
	// Only admins can update priority
	if userRole != "admin" && userRole != "super_admin" {
		h.respondError(w, http.StatusForbidden, "Only admins can update ticket priority")
		return
	}
	
	var req struct {
		Priority string `json:"priority"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate priority
	validPriorities := map[string]bool{
		"critical": true,
		"high":     true,
		"medium":   true,
		"low":      true,
	}
	
	if !validPriorities[req.Priority] {
		h.respondError(w, http.StatusBadRequest, "Invalid priority value")
		return
	}
	
	// Update ticket priority in database
	if err := h.service.UpdatePriority(ctx, ticketID, domain.TicketPriority(req.Priority)); err != nil {
		h.logger.Error("Failed to update ticket priority",
			slog.String("ticket_id", ticketID),
			slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to update priority")
		return
	}
	
	// Log audit event
	if h.auditLogger != nil {
		h.auditLogger.LogAsync(ctx, &audit.AuditEvent{
			EventType:     "ticket_priority_updated",
			EventCategory: audit.CategoryTicket,
			EventAction:   audit.ActionUpdate,
			EventStatus:   audit.StatusSuccess,
			ResourceType:  stringPtr("ticket"),
			ResourceID:    &ticketID,
			Metadata: map[string]interface{}{
				"new_priority": req.Priority,
				"updated_by":   ctx.Value("user_id"),
			},
		})
	}
	
	h.respondJSON(w, http.StatusOK, map[string]string{
		"status":   "success",
		"priority": req.Priority,
	})
}
```

#### 2. Add Service Method

**File:** `internal/service-domain/service-ticket/app/service.go`

Add method to TicketService:

```go
// UpdatePriority updates the priority of a ticket (admin-only)
func (s *TicketService) UpdatePriority(ctx context.Context, ticketID string, priority domain.TicketPriority) error {
	s.logger.Info("Updating ticket priority",
		slog.String("ticket_id", ticketID),
		slog.String("priority", string(priority)))
	
	// Get existing ticket
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}
	
	// Update priority
	ticket.Priority = priority
	
	// Save to database
	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return fmt.Errorf("failed to update ticket priority: %w", err)
	}
	
	// Create event for audit trail
	event := &domain.TicketEvent{
		TicketID:    ticketID,
		EventType:   "priority_updated",
		Description: fmt.Sprintf("Priority updated to %s", priority),
		Timestamp:   time.Now(),
		Metadata: map[string]interface{}{
			"priority": priority,
		},
	}
	
	if err := s.eventRepo.Create(ctx, event); err != nil {
		s.logger.Warn("Failed to create priority update event", slog.String("error", err.Error()))
	}
	
	return nil
}
```

#### 3. Add Route

**File:** `internal/service-domain/service-ticket/module.go`

Add route in `MountRoutes`:

```go
r.Route("/tickets", func(r chi.Router) {
	// ... existing routes ...
	
	// Admin-only: Update priority
	r.With(authMiddleware.RequireAdmin).Patch("/{id}/priority", m.ticketHandler.UpdateTicketPriority)
	
	// ... rest of routes ...
})
```

#### 4. Create Auth Middleware

**File:** `internal/shared/middleware/auth.go` (create if doesn't exist)

```go
package middleware

import (
	"context"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	logger *slog.Logger
}

func NewAuthMiddleware(logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{logger: logger}
}

// RequireAdmin ensures user has admin role
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract role from JWT token or session
		// For now, check Authorization header
		authHeader := r.Header.Get("Authorization")
		
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		// TODO: Validate JWT token and extract role
		// For now, assume role is in header as "Bearer admin:token"
		parts := strings.Split(authHeader, ":")
		if len(parts) < 1 {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}
		
		role := parts[0]
		if role != "admin" && role != "super_admin" {
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}
		
		// Set role in context
		ctx := context.WithValue(r.Context(), "user_role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

### Testing Phase 2

```bash
# Test priority update (should fail for non-admin)
curl -X PATCH http://localhost:8081/api/v1/tickets/{ticket_id}/priority \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer user:token" \
  -d '{"priority": "high"}'

# Expected: 403 Forbidden

# Test priority update (should succeed for admin)
curl -X PATCH http://localhost:8081/api/v1/tickets/{ticket_id}/priority \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer admin:token" \
  -d '{"priority": "critical"}'

# Expected: 200 OK
```

## 🎙️ Phase 3: WhatsApp Audio Message Handling

### Implementation Steps

#### 1. Add Audio Message Detection

**File:** `internal/service-domain/whatsapp/handler.go`

Update `handleIncomingMessage` to detect audio:

```go
func (h *WhatsAppHandler) handleIncomingMessage(ctx context.Context, msg *WhatsAppMessage) {
	h.logger.Info("Received WhatsApp message",
		slog.String("from", msg.From),
		slog.String("type", msg.Type),
	)
	
	// Handle audio messages
	if msg.Type == "audio" {
		h.handleAudioMessage(ctx, msg)
		return
	}
	
	// ... existing text message handling ...
}
```

#### 2. Implement Audio Message Handler

Add new method to WhatsAppHandler:

```go
// handleAudioMessage processes WhatsApp audio messages
func (h *WhatsAppHandler) handleAudioMessage(ctx context.Context, msg *WhatsAppMessage) {
	h.logger.Info("Processing audio message",
		slog.String("from", msg.From),
		slog.String("media_url", msg.MediaURL),
	)
	
	// Extract QR code from caption (if provided)
	qrCode := ""
	if msg.Caption != "" {
		qrCode = h.extractQRCode(msg.Caption)
	}
	
	if qrCode == "" {
		h.sendHelpMessage(ctx, msg.From)
		return
	}
	
	// Lookup equipment
	equipment, err := h.equipmentService.GetEquipmentByQR(ctx, qrCode)
	if err != nil {
		h.logger.Error("Equipment not found", slog.String("qr_code", qrCode))
		h.sendErrorMessage(ctx, msg.From, "Equipment not found")
		return
	}
	
	// Download audio file
	audioFile, err := h.downloadAudioFile(ctx, msg.MediaURL)
	if err != nil {
		h.logger.Error("Failed to download audio", slog.String("error", err.Error()))
		h.sendErrorMessage(ctx, msg.From, "Failed to process audio message")
		return
	}
	
	// Convert audio to text using Whisper API
	transcription, err := h.transcribeAudio(ctx, audioFile)
	if err != nil {
		h.logger.Warn("Audio transcription failed", slog.String("error", err.Error()))
		// Continue without transcription
	}
	
	// Create ticket with audio attachment
	issueDescription := msg.Caption
	if transcription != "" {
		issueDescription = fmt.Sprintf("%s\n\n[Audio transcription]: %s", msg.Caption, transcription)
	} else {
		issueDescription = fmt.Sprintf("%s\n\n[Audio message received - transcription unavailable]", msg.Caption)
	}
	
	priority := h.determinePriority(issueDescription)
	
	// Create ticket
	ticket, err := h.createTicketFromWhatsApp(ctx, equipment, msg, issueDescription, priority)
	if err != nil {
		h.logger.Error("Failed to create ticket", slog.String("error", err.Error()))
		h.sendErrorMessage(ctx, msg.From, "Failed to create service ticket")
		return
	}
	
	// Attach audio file to ticket
	go h.attachAudioToTicket(context.Background(), ticket.ID, audioFile, transcription)
	
	// Send confirmation
	h.sendTicketConfirmation(ctx, msg.From, ticket)
}
```

#### 3. Implement Audio Transcription (Whisper API)

Add transcription method:

```go
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
		return "", err
	}
	io.Copy(part, file)
	
	// Add model
	writer.WriteField("model", "whisper-1")
	
	// Add language (optional, auto-detect if not specified)
	writer.WriteField("language", "en") // or "hi" for Hindi
	
	writer.Close()
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.openai.com/v1/audio/transcriptions",
		body)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Execute request
	client := &http.Client{Timeout: 30 * time.Second}
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
	
	h.logger.Info("Audio transcribed successfully",
		slog.Int("text_length", len(result.Text)))
	
	return result.Text, nil
}
```

#### 4. Attach Audio to Ticket

```go
// attachAudioToTicket saves audio file and transcription as ticket attachment
func (h *WhatsAppHandler) attachAudioToTicket(ctx context.Context, ticketID string, audioFilePath string, transcription string) {
	// Create attachment record
	attachment := &domain.Attachment{
		ID:          uuid.New().String(),
		TicketID:    ticketID,
		EntityType:  "ticket",
		EntityID:    ticketID,
		Filename:    filepath.Base(audioFilePath),
		FileType:    "audio/mpeg",
		StoragePath: audioFilePath,
		Category:    "audio_message",
		Source:      "whatsapp",
		UploadedBy:  "whatsapp_system",
		UploadedAt:  time.Now(),
		Metadata: map[string]interface{}{
			"transcription": transcription,
			"has_transcription": transcription != "",
		},
	}
	
	// Save to database
	if err := h.attachmentRepo.Create(ctx, attachment); err != nil {
		h.logger.Error("Failed to create attachment record",
			slog.String("ticket_id", ticketID),
			slog.String("error", err.Error()))
		return
	}
	
	h.logger.Info("Audio attachment created",
		slog.String("ticket_id", ticketID),
		slog.String("attachment_id", attachment.ID))
}
```

### Environment Configuration

Add to `.env`:

```bash
# OpenAI Whisper API for audio transcription
OPENAI_API_KEY=sk-your-api-key-here

# WhatsApp Media Configuration
WHATSAPP_MEDIA_DIR=./storage/whatsapp_media
WHATSAPP_MAX_AUDIO_SIZE=16777216  # 16MB (WhatsApp limit)
```

### Testing Phase 3

```bash
# Test audio message handling
# 1. Send WhatsApp audio message with QR code in caption
# 2. Verify ticket is created
# 3. Verify audio file is attached
# 4. Verify transcription is included (if available)
# 5. Check logs for transcription success/failure
```

## 📊 Phase 4: Testing & Validation

### Test Scenarios

#### 1. Priority Field Removal
- [ ] Ticket creation form doesn't show priority field
- [ ] All tickets created have default priority='medium'
- [ ] Priority can be viewed on ticket details page
- [ ] Only admins can update priority

#### 2. Audio File Upload
- [ ] .mp3 files can be uploaded
- [ ] .wav files can be uploaded
- [ ] .m4a files can be uploaded
- [ ] .ogg files can be uploaded
- [ ] Files appear in selected files list
- [ ] Files are attached to ticket after creation
- [ ] Audio files can be played/downloaded

#### 3. WhatsApp Audio Messages
- [ ] Audio message with QR code creates ticket
- [ ] Audio file is downloaded and stored
- [ ] Transcription is attempted
- [ ] Ticket description includes transcription (if successful)
- [ ] Attachment record is created
- [ ] User receives confirmation message

#### 4. Admin Priority Updates
- [ ] Non-admin cannot update priority (403 Forbidden)
- [ ] Admin can update priority (200 OK)
- [ ] Priority update is logged in audit trail
- [ ] Priority change is visible in ticket details

### Performance Testing

```bash
# Test concurrent audio uploads
# Send 10 audio messages simultaneously
# Verify all tickets are created
# Check for any failures or timeouts
```

## 📝 Documentation Updates

### API Documentation

Update `docs/API.md` with:

```markdown
### Update Ticket Priority (Admin Only)

**Endpoint:** `PATCH /api/v1/tickets/{id}/priority`

**Authorization:** Admin role required

**Request:**
```json
{
  "priority": "critical"  // critical, high, medium, low
}
```

**Response:**
```json
{
  "status": "success",
  "priority": "critical"
}
```

**Errors:**
- 401: Unauthorized (no auth token)
- 403: Forbidden (not admin)
- 400: Invalid priority value
- 404: Ticket not found
```

## 🚀 Deployment Checklist

- [ ] Environment variables configured (OPENAI_API_KEY)
- [ ] WhatsApp media directory created
- [ ] Auth middleware configured
- [ ] Priority update endpoint deployed
- [ ] Audio transcription tested
- [ ] Error monitoring enabled
- [ ] Audit logging verified
- [ ] Documentation updated
- [ ] Team trained on new features

## 🔍 Monitoring

### Metrics to Track

1. **Audio Transcription Success Rate**
   - Successful transcriptions / Total audio messages
   - Target: >90%

2. **Priority Update Requests**
   - Track who is updating priorities
   - Audit admin actions

3. **Audio Message Volume**
   - Number of audio messages per day
   - Transcription API costs

4. **Default Priority Distribution**
   - Verify most tickets are medium
   - Track admin priority adjustments

### Alerts

- Alert if transcription success rate <80%
- Alert if OPENAI_API_KEY is missing
- Alert on repeated 403 errors (unauthorized priority updates)

## 📞 Support

For issues or questions:
- Check logs in `./logs/ticket-enhancements.log`
- Review audit trail for priority changes
- Test transcription API manually with curl
- Verify WhatsApp webhook is receiving messages

---

**Status:** Phase 1 Complete | Phase 2-4 Pending

**Last Updated:** December 23, 2025

**Next Steps:** Implement Phase 2 (Admin Priority API)
