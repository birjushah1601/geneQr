# Option 3 - WhatsApp Integration Implementation Guide

**Date:** December 21, 2025  
**Status:** ðŸ“‹ IMPLEMENTATION GUIDE  
**Existing Code:** Skeleton exists, needs activation  

---

## ðŸŽ¯ **OVERVIEW**

WhatsApp integration allows customers to create service tickets by sending messages to your business WhatsApp number. The system will:
1. Receive WhatsApp messages
2. Parse QR codes or equipment details
3. Create service tickets automatically
4. Send confirmations back to customers

---

## âœ… **WHAT ALREADY EXISTS**

### **Backend Files (Disabled):**

1. **`internal/service-domain/whatsapp/handler.go`** (330 lines)
   - WhatsApp webhook handler
   - Message parsing logic
   - Ticket creation from messages
   - **Status:** Has `//go:build ignore` - needs activation

2. **`internal/service-domain/whatsapp/webhook.go`**
   - Webhook registration
   - Twilio integration
   - **Status:** Needs review

3. **`internal/service-domain/whatsapp/media_handler.go`**
   - Media (images, documents) processing
   - QR code scanning from images
   - **Status:** Needs review

### **Database Schema (Exists):**

Tables already created:
- `whatsapp_conversations` - Track customer conversations
- `whatsapp_messages` - Store all messages
- Schema defined in `service-ticket/infra/schema.go`

### **Twilio Integration (Exists):**

**File:** `internal/infrastructure/sms/twilio.go`
- Twilio client setup
- SMS sending
- **Needs:** WhatsApp-specific methods

---

## ðŸš€ **IMPLEMENTATION STEPS**

### **Phase 1: Activate Existing Code** (2-3 hours)

#### **Step 1: Remove Build Ignore**

**File:** `internal/service-domain/whatsapp/handler.go`

```go
// Remove this line:
//go:build ignore

// Keep the rest of the file as is
```

**File:** `internal/service-domain/whatsapp/webhook.go`

```go
// Remove build ignore if present
```

**File:** `internal/service-domain/whatsapp/media_handler.go`

```go
// Remove build ignore if present
```

#### **Step 2: Create WhatsApp Module**

**New File:** `internal/service-domain/whatsapp/module.go`

```go
package whatsapp

import (
	"log/slog"
	
	"github.com/ServQR/medical-platform/internal/service-domain/equipment-registry/app"
	ticketApp "github.com/ServQR/medical-platform/internal/service-domain/service-ticket/app"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WhatsAppModule struct {
	handler *WhatsAppHandler
	logger  *slog.Logger
}

func NewWhatsAppModule(
	db *pgxpool.Pool,
	equipmentService *app.EquipmentService,
	ticketService *ticketApp.ServiceTicketService,
	twilioAccountSID string,
	twilioAuthToken string,
	twilioWhatsAppNumber string,
	logger *slog.Logger,
) *WhatsAppModule {
	// Create WhatsApp service
	whatsappService := NewWhatsAppService(
		twilioAccountSID,
		twilioAuthToken,
		twilioWhatsAppNumber,
		logger,
	)
	
	// Create handler
	handler := NewWhatsAppHandler(
		equipmentService,
		ticketService,
		whatsappService,
		logger,
	)
	
	return &WhatsAppModule{
		handler: handler,
		logger:  logger,
	}
}

func (m *WhatsAppModule) MountRoutes(r chi.Router) {
	r.Post("/whatsapp/webhook", m.handler.HandleWebhook)
	r.Get("/whatsapp/webhook", m.handler.HandleVerification) // For Twilio verification
}
```

#### **Step 3: Register Module in Main**

**File:** `cmd/platform/main.go`

```go
// In initializeModules function, add:

// Initialize WhatsApp module
if os.Getenv("ENABLE_WHATSAPP") == "true" {
	logger.Info("Initializing WhatsApp integration")
	
	whatsappModule := whatsapp.NewWhatsAppModule(
		dbPool,
		equipmentService, // Pass existing service
		ticketService,    // Pass existing service
		os.Getenv("TWILIO_ACCOUNT_SID"),
		os.Getenv("TWILIO_AUTH_TOKEN"),
		os.Getenv("TWILIO_WHATSAPP_NUMBER"),
		logger,
	)
	
	whatsappModule.MountRoutes(apiRouter)
	logger.Info("âœ… WhatsApp integration initialized")
}
```

---

### **Phase 2: Configure Twilio WhatsApp** (1-2 hours)

#### **Step 1: Twilio Setup**

1. **Go to Twilio Console:** https://console.twilio.com
2. **Navigate to:** Messaging â†’ Try it out â†’ Send a WhatsApp message
3. **Get Credentials:**
   - Account SID
   - Auth Token
   - WhatsApp-enabled number (starts with `whatsapp:+...`)

#### **Step 2: Configure Webhook**

In Twilio Console:
1. **Go to:** Messaging â†’ Settings â†’ WhatsApp sandbox
2. **Set Webhook URL:** `https://yourdomain.com/api/v1/whatsapp/webhook`
3. **Method:** POST
4. **Events:** Select "Incoming messages"

#### **Step 3: Environment Variables**

**File:** `.env`

```bash
# WhatsApp Configuration
ENABLE_WHATSAPP=true
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=your_auth_token_here
TWILIO_WHATSAPP_NUMBER=whatsapp:+14155238886

# WhatsApp Settings
WHATSAPP_VERIFY_TOKEN=your-verify-token-123
WHATSAPP_MEDIA_DIR=./data/whatsapp
```

---

### **Phase 3: Implement WhatsApp Service** (2-3 hours)

#### **Create WhatsApp Service**

**New File:** `internal/service-domain/whatsapp/service.go`

```go
package whatsapp

import (
	"context"
	"fmt"
	"log/slog"
	
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type WhatsAppService struct {
	client           *twilio.RestClient
	whatsappNumber   string
	logger           *slog.Logger
}

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

func (s *WhatsAppService) SendMessage(ctx context.Context, to, message string) error {
	s.logger.Info("Sending WhatsApp message",
		slog.String("to", to),
		slog.String("message_preview", message[:min(50, len(message))]))
	
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.whatsappNumber)
	params.SetBody(message)
	
	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		s.logger.Error("Failed to send WhatsApp message",
			slog.String("error", err.Error()))
		return err
	}
	
	s.logger.Info("WhatsApp message sent successfully",
		slog.String("message_sid", *resp.Sid))
	
	return nil
}

func (s *WhatsAppService) SendTicketConfirmation(ctx context.Context, to, ticketNumber string) error {
	message := fmt.Sprintf(
		"âœ… Service ticket created successfully!\n\n"+
		"Ticket Number: %s\n\n"+
		"Our engineer will contact you soon.\n\n"+
		"Reply to this message for updates.",
		ticketNumber,
	)
	
	return s.SendMessage(ctx, to, message)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

---

### **Phase 4: Test WhatsApp Integration** (1-2 hours)

#### **Test 1: Webhook Verification**

```bash
# Test that webhook responds
curl -X GET "http://localhost:8081/api/v1/whatsapp/webhook"

# Expected: 200 OK
```

#### **Test 2: Send Test Message**

From your phone:
1. Join Twilio WhatsApp sandbox (send "join <code>" to sandbox number)
2. Send message: "Help with equipment EQ-001"
3. Check backend logs for message processing
4. Verify ticket created in database
5. Should receive confirmation message back

#### **Test 3: QR Code Scanning**

1. Take photo of equipment QR code
2. Send image via WhatsApp
3. System should:
   - Download image
   - Scan QR code
   - Create ticket
   - Send confirmation

---

## ðŸ“Š **COMPLETE FLOW DIAGRAM**

```
Customer                 WhatsApp                Backend                Database
   |                        |                       |                      |
   |-- Send Message ------->|                       |                      |
   |                        |-- Webhook ----------->|                      |
   |                        |                       |                      |
   |                        |                       |-- Parse Message ---->|
   |                        |                       |                      |
   |                        |                       |-- Extract QR/Info -->|
   |                        |                       |                      |
   |                        |                       |-- Get Equipment ---->|
   |                        |                       |<-- Equipment --------|
   |                        |                       |                      |
   |                        |                       |-- Create Ticket ---->|
   |                        |                       |<-- Ticket Created ---|
   |                        |                       |                      |
   |                        |<-- Send Confirmation--|                      |
   |<-- Confirmation -------|                       |                      |
```

---

## ðŸŽ¯ **FEATURES TO IMPLEMENT**

### **Message Parsing:**

**Patterns to recognize:**
```
1. QR Code: "EQ-001", "QR:EQ-001"
2. Equipment ID: "Equipment ID: EQ-001"
3. Serial Number: "Serial: ABC123"
4. Issue Description: Free text
```

### **Conversation State:**

```go
// Track conversation state in database
type ConversationState struct {
	Step         int    // 1=equipment, 2=issue, 3=confirm
	EquipmentID  string
	IssueText    string
	CustomerInfo string
}
```

### **Multi-Step Flow:**

```
Bot: "Please send equipment QR code or ID"
User: "EQ-001"
Bot: "Found: Siemens MRI Scanner at Apollo Hospital. What's the issue?"
User: "Not starting, showing error E02"
Bot: "Creating ticket..."
Bot: "âœ… Ticket #T-12345 created! Engineer will contact you."
```

---

## ðŸ“‹ **TESTING CHECKLIST**

### **Phase 1 - Basic Integration:**
- [ ] WhatsApp module compiles
- [ ] Webhook endpoint responds
- [ ] Twilio credentials valid
- [ ] Can receive messages
- [ ] Can send messages

### **Phase 2 - Message Processing:**
- [ ] Parse text messages
- [ ] Extract QR codes
- [ ] Find equipment by ID
- [ ] Create tickets
- [ ] Send confirmations

### **Phase 3 - Media Handling:**
- [ ] Receive images
- [ ] Download media
- [ ] Scan QR codes from images
- [ ] Process documents

### **Phase 4 - Conversation Management:**
- [ ] Track conversation state
- [ ] Multi-step flows
- [ ] Handle errors gracefully
- [ ] Session timeouts

---

## ðŸš¨ **COMMON ISSUES & SOLUTIONS**

### **Issue 1: Webhook Not Receiving Messages**

**Symptoms:** Messages sent but webhook not called

**Solutions:**
1. Check Twilio webhook configuration
2. Ensure public URL is accessible
3. Verify ngrok/tunnel if using localhost
4. Check Twilio debugger logs

### **Issue 2: Can't Send Messages**

**Symptoms:** Error: "Message failed to send"

**Solutions:**
1. Verify Twilio credentials
2. Check WhatsApp number format (`whatsapp:+...`)
3. Ensure recipient joined sandbox (for testing)
4. Check Twilio account balance

### **Issue 3: QR Code Not Detected**

**Symptoms:** Image received but QR not scanned

**Solutions:**
1. Install QR code library: `go get github.com/makiuchi-d/gozxing`
2. Improve image preprocessing
3. Check image quality/resolution
4. Add manual fallback

---

## ðŸ’° **COST ESTIMATION**

### **Twilio WhatsApp Pricing:**

**Conversations (24-hour window):**
- Business-initiated: $0.005 - $0.02 per conversation
- User-initiated: $0.005 - $0.02 per conversation

**Monthly Estimate:**
- 100 tickets/month: ~$2-4
- 500 tickets/month: ~$10-20
- 1,000 tickets/month: ~$20-40

**Much cheaper than SMS!**

---

## âœ… **DELIVERABLES**

### **Code Files:**
1. âœ… WhatsApp handler (exists, needs activation)
2. âœ… WhatsApp webhook (exists, needs activation)
3. âœ… Media handler (exists, needs activation)
4. [ ] WhatsApp service (needs creation)
5. [ ] WhatsApp module (needs creation)

### **Configuration:**
6. [ ] Twilio account setup
7. [ ] Webhook configuration
8. [ ] Environment variables
9. [ ] Testing checklist

### **Documentation:**
10. âœ… Implementation guide (this document)
11. [ ] User guide (how customers use it)
12. [ ] Troubleshooting guide

---

## ðŸŽ¯ **ESTIMATED TIME**

**Total: 1-2 days**

- Phase 1 (Activation): 2-3 hours
- Phase 2 (Configuration): 1-2 hours
- Phase 3 (Service Implementation): 2-3 hours
- Phase 4 (Testing): 1-2 hours
- Documentation: 1 hour

**Most code already exists! Just needs activation and testing.**

---

## ðŸ“ž **SUPPORT RESOURCES**

**Twilio Documentation:**
- https://www.twilio.com/docs/whatsapp
- https://www.twilio.com/docs/whatsapp/quickstart/go

**Twilio Go SDK:**
- https://github.com/twilio/twilio-go

**WhatsApp Business API:**
- https://developers.facebook.com/docs/whatsapp

---

**Document:** Option 3 WhatsApp Implementation Guide  
**Last Updated:** December 21, 2025  
**Status:** ðŸ“‹ READY TO IMPLEMENT  
**Existing Code:** 80% complete, needs activation  
**Estimated Time:** 1-2 days  
**Business Value:** HIGH - Modern customer experience
