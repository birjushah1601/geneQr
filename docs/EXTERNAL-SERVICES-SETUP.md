# External Services Configuration Guide

**Date:** December 21, 2025  
**Purpose:** Configure Twilio (SMS/WhatsApp) and SendGrid (Email) for production  
**Prerequisites:** Authentication system already integrated  

---

## ðŸ“‹ **OVERVIEW**

The authentication system currently uses **mock services** for development. This guide walks you through configuring **real external services** for production.

### **Services to Configure:**
1. **Twilio** - SMS and WhatsApp OTP delivery
2. **SendGrid** - Email OTP delivery

### **Current Status:**
- âœ… Mock services work in development (no configuration needed)
- âœ… Code supports both mock and real services
- âœ… Automatic fallback to mock if credentials missing
- â³ Real services need API keys for production

---

## ðŸ”§ **TWILIO SETUP (SMS & WHATSAPP)**

### **Step 1: Create Twilio Account**

1. **Sign up:** https://www.twilio.com/try-twilio
2. **Get free trial:** $15.50 credit (good for ~500 SMS)
3. **Verify your phone number** (required for trial)

### **Step 2: Get Credentials**

1. **Go to Console:** https://console.twilio.com
2. **Find credentials:**
   - **Account SID:** `ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
   - **Auth Token:** Click "show" to reveal
3. **Copy both values**

### **Step 3: Get Phone Number**

**For SMS:**
1. Go to: **Phone Numbers â†’ Manage â†’ Buy a number**
2. Select your country
3. Check "SMS" capability
4. Click "Buy" (uses trial credit)
5. Copy the number: `+1234567890`

**For WhatsApp:**
1. Go to: **Messaging â†’ Try it out â†’ Try WhatsApp**
2. Follow setup wizard
3. Get WhatsApp-enabled number (different from SMS)
4. Copy WhatsApp number: `whatsapp:+1234567890`

### **Step 4: Configure Environment**

Create/update `.env` file:

```bash
# Twilio Configuration
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=your_auth_token_here
TWILIO_PHONE_NUMBER=+1234567890
TWILIO_WHATSAPP_NUMBER=whatsapp:+1234567890
```

### **Step 5: Verify Setup**

**Backend will automatically detect Twilio credentials:**

```
# If credentials found:
âœ… Twilio SMS service initialized

# If credentials missing:
âš ï¸ Twilio not configured, using mock SMS service
```

### **Step 6: Test Real SMS**

```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "+1234567890",
    "full_name": "Test User"
  }'
```

**Check your phone - you should receive real SMS!** ðŸ“±

### **Step 7: Test WhatsApp**

```bash
curl -X POST http://localhost:8081/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "+1234567890",
    "delivery_method": "whatsapp"
  }'
```

**Check WhatsApp - you should receive OTP!** ðŸ’¬

---

## ðŸ“§ **SENDGRID SETUP (EMAIL)**

### **Step 1: Create SendGrid Account**

1. **Sign up:** https://signup.sendgrid.com
2. **Free tier:** 100 emails/day forever
3. **Verify email address**

### **Step 2: Create API Key**

1. **Go to:** Settings â†’ API Keys
2. **Click:** "Create API Key"
3. **Name:** "ServQR Platform Production"
4. **Access:** "Full Access" (or "Mail Send" only)
5. **Create & Copy:** `SG.xxxxxxxxxxxxxxxxxxxxxxx`
6. **Save immediately** - shown only once!

### **Step 3: Verify Sender**

**Single Sender Verification (Quick):**
1. Go to: **Settings â†’ Sender Authentication**
2. Click: **"Verify a Single Sender"**
3. Fill in:
   - From Name: `ServQR Platform`
   - From Email: `noreply@yourdomain.com`
   - Company: `ServQR`
   - Address: Your company address
4. **Verify** - check email and click verification link

**Domain Authentication (Production):**
1. Go to: **Settings â†’ Sender Authentication**
2. Click: **"Authenticate Your Domain"**
3. Choose DNS host
4. Add DNS records (CNAME, TXT)
5. **Verify** - improves deliverability

### **Step 4: Configure Environment**

Update `.env` file:

```bash
# SendGrid Configuration
SENDGRID_API_KEY=SG.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
SENDGRID_FROM_EMAIL=noreply@yourdomain.com
SENDGRID_FROM_NAME=ServQR Platform
```

### **Step 5: Verify Setup**

**Backend will automatically detect SendGrid:**

```
# If API key found:
âœ… SendGrid email service initialized

# If API key missing:
âš ï¸ SendGrid not configured, using mock email service
```

### **Step 6: Test Real Email**

```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "your.email@example.com",
    "full_name": "Test User"
  }'
```

**Check your inbox - you should receive real email!** ðŸ“§

---

## ðŸ” **ENVIRONMENT CONFIGURATION**

### **Development (.env)**

```bash
# ============================================================================
# DEVELOPMENT ENVIRONMENT
# ============================================================================

# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5430/med_platform?sslmode=disable

# Server
PORT=8081
ENVIRONMENT=development

# JWT Keys
JWT_PRIVATE_KEY_PATH=./keys/jwt-private.pem
JWT_PUBLIC_KEY_PATH=./keys/jwt-public.pem

# External Services (Optional in dev - uses mocks if missing)
# TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
# TWILIO_AUTH_TOKEN=your_auth_token_here
# TWILIO_PHONE_NUMBER=+1234567890
# TWILIO_WHATSAPP_NUMBER=whatsapp:+1234567890

# SENDGRID_API_KEY=SG.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
# SENDGRID_FROM_EMAIL=noreply@yourdomain.com
# SENDGRID_FROM_NAME=ServQR Platform

# Authentication
ENABLE_AUTH=true
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d
OTP_EXPIRY=5m
OTP_MAX_ATTEMPTS=3
OTP_RATE_LIMIT_PER_HOUR=3
```

### **Production (.env.production)**

```bash
# ============================================================================
# PRODUCTION ENVIRONMENT
# ============================================================================

# Database (Use production DB URL)
DATABASE_URL=postgres://user:password@prod-db.amazonaws.com:5432/med_platform?sslmode=require

# Server
PORT=8081
ENVIRONMENT=production

# JWT Keys (Use production keys)
JWT_PRIVATE_KEY_PATH=/secrets/jwt-private.pem
JWT_PUBLIC_KEY_PATH=/secrets/jwt-public.pem

# External Services (REQUIRED in production)
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=your_production_auth_token
TWILIO_PHONE_NUMBER=+1234567890
TWILIO_WHATSAPP_NUMBER=whatsapp:+1234567890

SENDGRID_API_KEY=SG.production_key_here
SENDGRID_FROM_EMAIL=noreply@yourcompany.com
SENDGRID_FROM_NAME=ServQR Platform

# Authentication (Production values)
ENABLE_AUTH=true
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d
OTP_EXPIRY=5m
OTP_MAX_ATTEMPTS=3
OTP_RATE_LIMIT_PER_HOUR=3

# Security
ENABLE_HTTPS=true
ENABLE_RATE_LIMIT=true
RATE_LIMIT_PER_MINUTE=100
```

---

## ðŸ“Š **COST ESTIMATION**

### **Twilio Costs:**

**SMS:**
- US/Canada: $0.0079 per SMS
- International: $0.04 - $0.15 per SMS
- **1,000 OTPs/month:** ~$8/month

**WhatsApp:**
- Conversations: $0.005 - $0.02 per conversation
- **1,000 OTPs/month:** ~$5-20/month

**Total Twilio:** ~$13-28/month for 1,000 users

### **SendGrid Costs:**

**Free Tier:**
- 100 emails/day = 3,000/month
- **Perfect for starting!**

**Essentials ($19.95/month):**
- 50,000 emails/month
- Good for 1,600 OTPs/day

**Total SendGrid:** $0-20/month depending on volume

### **Combined Monthly Cost:**
- **Small scale (100 users/day):** ~$0-15/month
- **Medium scale (1,000 users/day):** ~$30-50/month
- **Large scale (10,000 users/day):** ~$200-300/month

---

## ðŸ§ª **TESTING CHECKLIST**

### **Development Testing (Mock Services):**
- [ ] Backend starts without external credentials
- [ ] Mock email logs OTP to console
- [ ] Mock SMS logs OTP to console
- [ ] Mock WhatsApp logs OTP to console
- [ ] OTP verification works correctly
- [ ] Users can register and login

### **Twilio SMS Testing:**
- [ ] Add Twilio credentials to .env
- [ ] Restart backend
- [ ] Register with real phone number
- [ ] Receive SMS within 30 seconds
- [ ] OTP is correct (6 digits)
- [ ] Verify OTP successfully

### **Twilio WhatsApp Testing:**
- [ ] Configure WhatsApp number
- [ ] Request OTP via WhatsApp
- [ ] Receive WhatsApp message
- [ ] OTP format is correct
- [ ] Verify OTP successfully

### **SendGrid Email Testing:**
- [ ] Add SendGrid API key to .env
- [ ] Verify sender email
- [ ] Restart backend
- [ ] Register with real email
- [ ] Receive email within 30 seconds
- [ ] Email formatting is correct
- [ ] OTP is visible and correct
- [ ] Verify OTP successfully

### **Production Readiness:**
- [ ] All credentials in production .env
- [ ] Secrets stored securely (not in git)
- [ ] Domain authentication for email
- [ ] Phone numbers verified
- [ ] Rate limiting tested
- [ ] Error handling verified
- [ ] Monitoring configured
- [ ] Alerts set up for failures

---

## ðŸš¨ **TROUBLESHOOTING**

### **Twilio: SMS Not Received**

**Problem:** Registered but no SMS
**Solutions:**
1. Check phone number format: `+1234567890` (include country code)
2. Verify phone is verified in Twilio console (trial accounts)
3. Check Twilio logs: Console â†’ Monitor â†’ Logs
4. Check backend logs for errors
5. Verify TWILIO_PHONE_NUMBER is correct
6. Check trial credit balance

### **Twilio: WhatsApp Not Working**

**Problem:** WhatsApp OTP not received
**Solutions:**
1. Join WhatsApp sandbox first (trial accounts)
2. Send "join <code>" to Twilio WhatsApp number
3. Use format: `whatsapp:+1234567890`
4. Check WhatsApp is enabled on number
5. Verify user's WhatsApp is active

### **SendGrid: Email Not Received**

**Problem:** Registered but no email
**Solutions:**
1. Check spam/junk folder
2. Verify sender email in SendGrid
3. Check SendGrid activity: Email Activity
4. Verify API key has "Mail Send" permission
5. Check from email is verified
6. Try different email provider (Gmail, Outlook)
7. Check backend logs for errors

### **General: Mock Services Not Working**

**Problem:** Console not showing mock logs
**Solutions:**
1. Verify backend is running
2. Check log level (should show INFO)
3. Remove external credentials to force mock mode
4. Restart backend
5. Check OTP service initialization logs

---

## ðŸ”’ **SECURITY BEST PRACTICES**

### **1. API Key Management:**
- âœ… **Never commit API keys to git**
- âœ… Use `.env` files (add to `.gitignore`)
- âœ… Use environment variables in production
- âœ… Rotate keys periodically (every 90 days)
- âœ… Use separate keys for dev/staging/prod

### **2. Phone Number Protection:**
- âœ… Validate phone format before sending
- âœ… Rate limit OTP requests (3 per hour)
- âœ… Implement cooldown (60 seconds between requests)
- âœ… Log all OTP attempts
- âœ… Block suspicious patterns

### **3. Email Protection:**
- âœ… Validate email format
- âœ… Check disposable email domains
- âœ… Rate limit email sending
- âœ… Use SPF/DKIM/DMARC
- âœ… Monitor bounce rates

### **4. OTP Security:**
- âœ… 6-digit codes (100,000 combinations)
- âœ… 5-minute expiry
- âœ… Single-use only
- âœ… SHA-256 hashing in database
- âœ… Account lockout after 5 failed attempts

---

## ðŸ“ˆ **MONITORING & ALERTS**

### **Key Metrics to Track:**

**Twilio:**
- SMS delivery rate (should be >95%)
- Average delivery time (should be <30s)
- Failed sends
- Cost per SMS

**SendGrid:**
- Email delivery rate (should be >98%)
- Bounce rate (should be <2%)
- Spam complaint rate (should be <0.1%)
- Open rate (informational)

**OTP System:**
- OTPs sent per hour
- Verification success rate
- Average time to verify
- Failed verification attempts
- Account lockouts

### **Alert Thresholds:**

**Critical Alerts:**
- Twilio API errors
- SendGrid API errors
- OTP delivery failure rate >10%
- Account lockout spike

**Warning Alerts:**
- SMS cost spike
- Email bounce rate >5%
- Unusual OTP request patterns

---

## âœ… **COMPLETION CHECKLIST**

### **Development Setup:**
- [ ] Backend runs with mock services
- [ ] Console shows mock OTP logs
- [ ] OTP verification works
- [ ] Users can register/login

### **Twilio Setup:**
- [ ] Twilio account created
- [ ] Credentials obtained
- [ ] Phone number purchased
- [ ] WhatsApp configured
- [ ] Environment variables set
- [ ] Real SMS tested
- [ ] Real WhatsApp tested

### **SendGrid Setup:**
- [ ] SendGrid account created
- [ ] API key created
- [ ] Sender email verified
- [ ] Environment variables set
- [ ] Real email tested
- [ ] Email formatting verified

### **Production Ready:**
- [ ] All services tested
- [ ] Monitoring configured
- [ ] Alerts set up
- [ ] Documentation complete
- [ ] Team trained
- [ ] Backup credentials stored
- [ ] Incident response plan ready

---

## ðŸŽ¯ **NEXT STEPS**

1. **Now:** Keep using mock services for development
2. **Before user testing:** Configure Twilio SMS
3. **Before beta launch:** Configure SendGrid email
4. **Before production:** Configure WhatsApp, monitoring, alerts
5. **After launch:** Monitor metrics, optimize costs

---

**Document:** External Services Setup Guide  
**Last Updated:** December 21, 2025  
**Status:** Ready for configuration when needed  
**Estimated Time:** 1-2 hours to configure both services
