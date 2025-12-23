# Week 1 Implementation Guide - Authentication Integration

**Duration:** 5-7 days  
**Priority:** CRITICAL 🔥  
**Goal:** Fully functional authentication system integrated with existing platform  

---

## 📅 **DAY-BY-DAY BREAKDOWN**

---

## **DAY 1: Deploy & Test Authentication** ✅

### **Morning (2-3 hours)**

#### **Step 1: Run Complete Setup**
```powershell
# Run the automated setup script
.\scripts\complete-setup.ps1
```

**This will:**
- ✅ Generate JWT keys
- ✅ Apply database migrations
- ✅ Install Go dependencies
- ✅ Build backend
- ✅ Configure frontend
- ✅ Create startup scripts

#### **Step 2: Start the System**
```powershell
# Start both backend and frontend
.\start-platform.ps1
```

Or manually:
```powershell
# Terminal 1: Backend
.\start-backend.bat

# Terminal 2: Frontend  
cd admin-ui
npm run dev
```

#### **Step 3: Verify Health**
```bash
# Test backend health
curl http://localhost:8080/health

# Should return: {"status":"ok"}
```

### **Afternoon (2-3 hours)**

#### **Step 4: Test Authentication Endpoints**

**Test Registration:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "test@example.com",
    "full_name": "Test User",
    "password": "SecurePass123!"
  }'
```

**Expected Response:**
```json
{
  "user_id": "uuid-here",
  "requires_otp": true,
  "otp_sent_to": "t***@example.com",
  "expires_in": 300
}
```

**Check Backend Logs** - You should see:
```
📧 MOCK EMAIL to=test@example.com otp=123456
```

**Test OTP Verification:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "test@example.com",
    "code": "123456"
  }'
```

**Expected Response:**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

**Test Token Validation:**
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### **Step 5: Test Frontend**

1. **Open:** http://localhost:3000/register
2. **Fill form:**
   - Full Name: Test User
   - Email: test@example.com
   - Password: (optional)
3. **Click "Create Account"**
4. **Check backend logs** for OTP code
5. **Enter OTP code**
6. **Should redirect to dashboard**

**Test Login:**
1. **Open:** http://localhost:3000/login
2. **Enter:** test@example.com
3. **Click "Send OTP"**
4. **Check logs** for OTP
5. **Enter OTP**
6. **Should login successfully**

### **Evening (1-2 hours)**

#### **Step 6: Document What Works**

Create a test results file:
```markdown
# Authentication Test Results - Day 1

## ✅ Working
- [x] Registration with email
- [x] OTP generation and delivery (mock)
- [x] OTP verification
- [x] JWT token generation
- [x] Token validation
- [x] Frontend registration
- [x] Frontend login
- [x] Token refresh

## ⏳ To Test
- [ ] Password login
- [ ] Password reset
- [ ] Multiple sessions
- [ ] Token expiry handling
```

**Deliverable:** Authentication system running and tested ✅

---

## **DAY 2-3: Integrate with Existing System**

### **Goal:** Protect all existing API routes with authentication

### **Task 1: Add Auth Check to Existing Endpoints** (4-6 hours)

#### **Option A: Quick Integration (Recommended for Week 1)**

Add this to `cmd/platform/main.go` after module initialization:

```go
// Add at line ~380 (after AI services init)

// ========================================================================
// INITIALIZE AUTHENTICATION
// ========================================================================
logger.Info("Initializing Authentication System")

// Create simple database connection for auth
import "github.com/aby-med/medical-platform/internal/core/auth"

authDB, err := sqlx.Connect("postgres", cfg.GetDSN())
if err != nil {
	logger.Warn("Failed to connect auth database", slog.String("error", err.Error()))
} else {
	err = auth.IntegrateAuthModule(router, authDB, logger)
	if err != nil {
		logger.Warn("Failed to initialize auth module", slog.String("error", err.Error()))
	} else {
		logger.Info("✅ Authentication system initialized")
	}
}
```

#### **Option B: Full Integration with Auth Middleware**

Create protected route group:

```go
// In initializeModules function, wrap protected routes:

// Get auth handler for middleware
authHandler := getAuthHandler() // Helper function you create

// Create protected routes group
r.Route("/api/v1", func(r chi.Router) {
	// Public routes (no auth)
	r.Get("/health", healthHandler)
	
	// Protected routes (require auth)
	r.Group(func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)
		
		// All protected endpoints here
		r.Mount("/equipment", equipmentModule.Routes())
		r.Mount("/tickets", ticketModule.Routes())
		r.Mount("/engineers", engineerModule.Routes())
		// ... etc
	})
})
```

### **Task 2: Update Frontend API Client** (2-3 hours)

#### **Update API helper to include auth headers:**

Create/update `admin-ui/src/lib/api/client.ts`:

```typescript
import { useAuth } from '@/contexts/AuthContext';

export function createAPIClient() {
  const { accessToken } = useAuth();
  
  return {
    get: async (url: string) => {
      const response = await fetch(url, {
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
        },
      });
      
      if (response.status === 401) {
        // Token expired, redirect to login
        window.location.href = '/login';
        throw new Error('Unauthorized');
      }
      
      return response.json();
    },
    
    post: async (url: string, data: any) => {
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });
      
      if (response.status === 401) {
        window.location.href = '/login';
        throw new Error('Unauthorized');
      }
      
      return response.json();
    },
    // ... put, delete, etc
  };
}
```

#### **Update all API calls to use client:**

```typescript
// Before
const response = await fetch('/api/v1/equipment');

// After
import { createAPIClient } from '@/lib/api/client';
const api = createAPIClient();
const data = await api.get('/api/v1/equipment');
```

### **Task 3: Add User Context to UI** (1-2 hours)

#### **Update layout to show logged-in user:**

```typescript
// admin-ui/src/app/layout.tsx or dashboard layout
import { useAuth } from '@/contexts/AuthContext';

export function DashboardLayout({ children }) {
  const { user, logout } = useAuth();
  
  return (
    <div>
      <header>
        <div className="user-info">
          <span>Welcome, {user?.name}</span>
          <button onClick={logout}>Logout</button>
        </div>
      </header>
      <main>{children}</main>
    </div>
  );
}
```

**Deliverable:** All routes protected, frontend using auth ✅

---

## **DAY 4-5: Production Configuration**

### **Goal:** Configure external services and security

### **Task 1: Configure Twilio (SMS/WhatsApp)** (1-2 hours)

1. **Sign up:** https://www.twilio.com/try-twilio
2. **Get credentials:**
   - Account SID
   - Auth Token
   - Phone Number
   - WhatsApp Number

3. **Add to `.env.local`:**
```bash
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=your_auth_token_here
TWILIO_PHONE_NUMBER=+1234567890
TWILIO_WHATSAPP_NUMBER=+1234567890
```

4. **Test real SMS:**
```bash
# Register with real phone number
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "+1234567890",
    "full_name": "Test User"
  }'

# Check your phone for SMS!
```

### **Task 2: Configure SendGrid (Email)** (1-2 hours)

1. **Sign up:** https://signup.sendgrid.com
2. **Create API key:** Settings → API Keys
3. **Verify sender:** Settings → Sender Authentication

4. **Add to `.env.local`:**
```bash
SENDGRID_API_KEY=SG.xxxxxxxxxxx
SENDGRID_FROM_EMAIL=noreply@yourdomain.com
SENDGRID_FROM_NAME=ABY-MED Platform
```

5. **Test real email:**
```bash
# Register with real email
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "your.email@example.com",
    "full_name": "Test User"
  }'

# Check your email inbox!
```

### **Task 3: Security Hardening** (2-3 hours)

#### **Enable HTTPS (Production):**

```go
// For production, use TLS
server := &http.Server{
	Addr:      ":443",
	Handler:   router,
	TLSConfig: &tls.Config{
		MinVersion: tls.VersionTLS13,
	},
}

err := server.ListenAndServeTLS("cert.pem", "key.pem")
```

#### **Rate Limiting:**

```go
import "github.com/go-chi/httprate"

// Add to router setup
r.Use(httprate.LimitByIP(100, 1*time.Minute))
```

#### **Security Headers:**

```go
r.Use(func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		next.ServeHTTP(w, r)
	})
})
```

**Deliverable:** Production-ready security configuration ✅

---

## **DAY 6-7: Testing & Verification**

### **Goal:** Comprehensive testing of authentication system

### **Test Scenarios:**

#### **1. Registration Flow**
- [ ] Register with email
- [ ] Register with phone
- [ ] Receive OTP via email
- [ ] Receive OTP via SMS
- [ ] Verify OTP successfully
- [ ] Handle invalid OTP
- [ ] Handle expired OTP
- [ ] Rate limiting works (3/hour)

#### **2. Login Flow**
- [ ] Login with OTP (email)
- [ ] Login with OTP (phone)
- [ ] Login with password
- [ ] Switch between OTP/password
- [ ] Handle failed login attempts
- [ ] Account locks after 5 attempts
- [ ] Account unlocks after 30 minutes

#### **3. Token Management**
- [ ] Access token works
- [ ] Refresh token works
- [ ] Token rotation on refresh
- [ ] Logout revokes tokens
- [ ] Expired tokens handled

#### **4. Security**
- [ ] Weak passwords rejected
- [ ] Common passwords blocked
- [ ] SQL injection protected
- [ ] XSS protected
- [ ] CSRF protected
- [ ] Rate limiting active

#### **5. Audit Logging**
- [ ] Registration logged
- [ ] Login attempts logged
- [ ] Failed logins logged
- [ ] Token refresh logged
- [ ] Logout logged

**Deliverable:** Fully tested authentication system ✅

---

## **WEEK 1 SUCCESS CRITERIA**

### **Must Have:**
- [x] Authentication system deployed
- [x] All 12 endpoints working
- [x] Frontend login/register working
- [x] JWT tokens working
- [x] Basic security active

### **Should Have:**
- [x] External services configured (Twilio/SendGrid)
- [x] Existing routes protected
- [x] Frontend using auth headers
- [x] Comprehensive testing done

### **Nice to Have:**
- [ ] Advanced security features
- [ ] Performance optimization
- [ ] Monitoring dashboards
- [ ] Documentation updates

---

## **TROUBLESHOOTING**

### **Issue: JWT keys not found**
```bash
# Regenerate keys
.\scripts\generate-jwt-keys.ps1
```

### **Issue: Database tables don't exist**
```bash
# Reapply migrations
go run scripts/apply-auth-migrations.go
```

### **Issue: OTP not showing in logs**
```
# Check if auth module initialized
# Look for: "✅ Authentication module initialized"
```

### **Issue: CORS errors**
```go
// Update CORS config in main.go
AllowedOrigins: []string{"http://localhost:3000"}
```

---

## **END OF WEEK 1**

**Status Check:**
- ✅ Authentication deployed and working
- ✅ System secured with JWT
- ✅ External services configured
- ✅ Comprehensive testing done

**Ready for Week 2:** Dashboard completion ✅

---

**Next:** See `WEEK2-IMPLEMENTATION-GUIDE.md`
