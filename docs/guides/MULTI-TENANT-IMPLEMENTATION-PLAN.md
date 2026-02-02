# Multi-Tenant Data Isolation - Implementation Plan

## ðŸ“‹ Overview

This document outlines the step-by-step plan to implement complete multi-tenant data isolation in the ServQR Platform. Each task is designed to be implemented incrementally and tested independently.

---

## ðŸŽ¯ Goals

1. **Data Isolation**: Ensure users only see data belonging to their organization
2. **Security**: Prevent unauthorized cross-organization access
3. **UI Adaptation**: Show organization-specific features and dashboards
4. **Testing**: Verify isolation with test users from different organizations

---

## ðŸ“Š Current Status

### âœ… What's Already Done
- [x] Multi-tenant database structure (organizations, user_organizations)
- [x] Authentication system with JWT tokens
- [x] JWT tokens include organization_id, role, permissions
- [x] Test users created for all organization types
- [x] User-organization relationships established

### ðŸš§ What Needs Implementation
- [ ] Backend: Organization context middleware
- [ ] Backend: Repository query filters
- [ ] Backend: API endpoint updates
- [ ] Frontend: Decode and store organization context
- [ ] Frontend: Organization-specific dashboards
- [ ] Frontend: Conditional feature rendering
- [ ] Testing: Data isolation verification

---

## ðŸ—ºï¸ Implementation Phases

### **Phase 1: Backend Foundation** (Priority: HIGH)
Set up the infrastructure to enforce organization-based data filtering.

### **Phase 2: API Data Filtering** (Priority: HIGH)  
Update all API endpoints to filter by organization.

### **Phase 3: Frontend Context** (Priority: MEDIUM)
Extract and use organization information from JWT tokens.

### **Phase 4: Organization-Specific UI** (Priority: MEDIUM)
Create different dashboards and features per organization type.

### **Phase 5: Testing & Validation** (Priority: HIGH)
Comprehensive testing to ensure data isolation works correctly.

---

# Phase 1: Backend Foundation

## Task 1.1: Create Organization Context Middleware â­

**Objective**: Extract organization information from JWT and inject into request context.

**Files to Create/Modify:**
- `internal/middleware/organization_context.go` (NEW)
- `cmd/platform/main.go` (MODIFY - register middleware)

**Implementation:**

```go
// internal/middleware/organization_context.go
package middleware

import (
    "context"
    "log/slog"
    "net/http"
    
    "github.com/google/uuid"
)

// Organization context keys
type contextKey string

const (
    OrganizationIDKey contextKey = "organization_id"
    OrganizationTypeKey contextKey = "organization_type"
    UserRoleKey contextKey = "user_role"
    UserPermissionsKey contextKey = "user_permissions"
)

// OrganizationContextMiddleware extracts organization info from JWT claims
// and injects it into the request context for downstream handlers
func OrganizationContextMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get JWT claims from context (set by auth middleware)
            claims, ok := r.Context().Value("claims").(map[string]interface{})
            if !ok {
                // No claims - might be public endpoint
                logger.Debug("No JWT claims found in request context")
                next.ServeHTTP(w, r)
                return
            }
            
            ctx := r.Context()
            
            // Extract organization_id
            if orgIDStr, ok := claims["organization_id"].(string); ok && orgIDStr != "" {
                if orgID, err := uuid.Parse(orgIDStr); err == nil {
                    ctx = context.WithValue(ctx, OrganizationIDKey, orgID)
                    logger.Debug("Organization context set", 
                        "organization_id", orgID,
                        "path", r.URL.Path)
                }
            }
            
            // Extract organization_type
            if orgType, ok := claims["organization_type"].(string); ok {
                ctx = context.WithValue(ctx, OrganizationTypeKey, orgType)
            }
            
            // Extract user role
            if role, ok := claims["role"].(string); ok {
                ctx = context.WithValue(ctx, UserRoleKey, role)
            }
            
            // Extract permissions
            if perms, ok := claims["permissions"].([]interface{}); ok {
                permissions := make([]string, len(perms))
                for i, p := range perms {
                    if pStr, ok := p.(string); ok {
                        permissions[i] = pStr
                    }
                }
                ctx = context.WithValue(ctx, UserPermissionsKey, permissions)
            }
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// GetOrganizationID extracts organization ID from context
func GetOrganizationID(ctx context.Context) (uuid.UUID, bool) {
    orgID, ok := ctx.Value(OrganizationIDKey).(uuid.UUID)
    return orgID, ok
}

// GetOrganizationType extracts organization type from context
func GetOrganizationType(ctx context.Context) (string, bool) {
    orgType, ok := ctx.Value(OrganizationTypeKey).(string)
    return orgType, ok
}

// GetUserRole extracts user role from context
func GetUserRole(ctx context.Context) (string, bool) {
    role, ok := ctx.Value(UserRoleKey).(string)
    return role, ok
}

// GetUserPermissions extracts user permissions from context
func GetUserPermissions(ctx context.Context) ([]string, bool) {
    perms, ok := ctx.Value(UserPermissionsKey).([]string)
    return perms, ok
}
```

**Register Middleware in main.go:**

```go
// cmd/platform/main.go

import (
    "github.com/ServQR/medical-platform/internal/middleware"
)

// In the main() function, after auth middleware:
router.Use(middleware.OrganizationContextMiddleware(logger))
```

**Testing:**
- [ ] Middleware compiles without errors
- [ ] Middleware extracts organization_id from JWT
- [ ] Context helpers return correct values
- [ ] Middleware logs organization context

**Acceptance Criteria:**
- âœ… Middleware successfully extracts org_id from JWT claims
- âœ… Organization context is available in all downstream handlers
- âœ… System logs show organization_id for each request

---

## Task 1.2: Add OrganizationType to JWT Claims â­

**Objective**: Include organization type in JWT so frontend knows what UI to show.

**Files to Modify:**
- `internal/core/auth/app/jwt_service.go`
- `internal/core/auth/app/auth_service.go`

**Implementation:**

```go
// internal/core/auth/app/jwt_service.go

type TokenRequest struct {
    UserID           uuid.UUID
    Email            string
    Name             string
    OrganizationID   string
    OrganizationType string  // ADD THIS
    Role             string
    Permissions      []string
    DeviceInfo       map[string]interface{}
    IPAddress        string
}

// In GenerateTokenPair method, add to claims:
claims := jwt.MapClaims{
    "user_id":           req.UserID.String(),
    "email":             req.Email,
    "name":              req.Name,
    "organization_id":   req.OrganizationID,
    "organization_type": req.OrganizationType,  // ADD THIS
    "role":              req.Role,
    "permissions":       req.Permissions,
    "exp":               time.Now().Add(s.accessTokenTTL).Unix(),
    "iat":               time.Now().Unix(),
    "iss":               s.issuer,
}
```

**Update auth_service.go to pass org_type:**

```go
// internal/core/auth/app/auth_service.go
// In LoginWithPassword method:

// Get user organizations
userOrgs, _ := s.userRepo.GetUserOrganizations(ctx, user.ID)
var primaryOrg *domain.UserOrganization
if len(userOrgs) > 0 {
    primaryOrg = &userOrgs[0]
}

// Fetch organization details to get org_type
var orgType string
if primaryOrg != nil {
    org, err := s.orgRepo.GetByID(ctx, primaryOrg.OrganizationID)
    if err == nil {
        orgType = org.Type  // Assuming Organization has Type field
    }
}

// Generate tokens
tokenReq := &TokenRequest{
    UserID:           user.ID,
    Email:            *user.Email,
    Name:             user.FullName,
    OrganizationType: orgType,  // ADD THIS
    DeviceInfo:       req.DeviceInfo,
    IPAddress:        req.IPAddress,
}
```

**Testing:**
- [ ] JWT token includes organization_type field
- [ ] Decode JWT at jwt.io and verify org_type is present
- [ ] Frontend can read org_type from token

**Acceptance Criteria:**
- âœ… JWT payload contains organization_type
- âœ… organization_type matches the user's organization
- âœ… Different users have different org_types

---

## Task 1.3: Create Organization Repository Interface

**Objective**: Add methods to fetch organization details.

**Files to Create/Modify:**
- `internal/core/auth/domain/organization_repository.go` (NEW or UPDATE)
- `internal/core/auth/infra/organization_repository.go` (NEW or UPDATE)

**Implementation:**

```go
// internal/core/auth/domain/organization_repository.go
package domain

type OrganizationRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Organization, error)
}

type Organization struct {
    ID       uuid.UUID
    Name     string
    Type     string  // manufacturer, hospital, etc.
    Status   string
    Metadata JSONBMap
}
```

```go
// internal/core/auth/infra/organization_repository.go
package infra

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
    var org domain.Organization
    query := `
        SELECT id, name, org_type, status, metadata
        FROM organizations
        WHERE id = $1
    `
    err := r.db.GetContext(ctx, &org, query, id)
    if err != nil {
        return nil, err
    }
    return &org, nil
}
```

---

# Phase 2: API Data Filtering

## Task 2.1: Update Equipment Registry Repository â­â­

**Objective**: Filter equipment by organization - users only see equipment they own or manufactured.

**Files to Modify:**
- `internal/service-domain/service-equipment/infra/equipment_repository.go`

**Business Rules:**
- **Manufacturers**: See ALL equipment they manufactured (across all organizations)
- **Hospitals**: See ONLY equipment they own
- **Channel Partners/Sub-Sub-sub_sub_SUB_DEALERs**: See equipment they sold/service

**Implementation:**

```go
// Add organization filter to GetAll method
func (r *EquipmentRepository) GetAll(ctx context.Context) ([]*domain.Equipment, error) {
    // Extract organization context
    orgID, ok := middleware.GetOrganizationID(ctx)
    if !ok {
        return nil, fmt.Errorf("organization context not found")
    }
    
    orgType, _ := middleware.GetOrganizationType(ctx)
    
    var query string
    
    switch orgType {
    case "manufacturer":
        // Manufacturers see equipment they manufactured
        query = `
            SELECT * FROM equipment_registry 
            WHERE manufacturer_id = $1
            ORDER BY created_at DESC
        `
    case "hospital", "imaging_center":
        // Hospitals see equipment they own
        query = `
            SELECT * FROM equipment_registry 
            WHERE organization_id = $1
            OR owner_org_id = $1
            ORDER BY created_at DESC
        `
    case "Channel Partner", "Sub-sub_SUB_DEALER":
        // Channel Partners see equipment they sold/service
        query = `
            SELECT * FROM equipment_registry 
            WHERE channel_partner_org_id = $1
            OR service_provider_org_id = $1
            ORDER BY created_at DESC
        `
    default:
        // Default: only owned equipment
        query = `
            SELECT * FROM equipment_registry 
            WHERE organization_id = $1
            ORDER BY created_at DESC
        `
    }
    
    var equipment []*domain.Equipment
    err := r.db.SelectContext(ctx, &equipment, query, orgID)
    return equipment, err
}
```

**Testing:**
- [ ] Manufacturer login â†’ sees only their manufactured equipment
- [ ] Hospital login â†’ sees only their owned equipment
- [ ] No cross-organization equipment visible

---

## Task 2.2: Update Service Tickets Repository â­â­

**Objective**: Filter tickets by organization.

**Files to Modify:**
- `internal/service-domain/service-ticket/infra/ticket_repository.go`

**Business Rules:**
- **Hospitals**: See tickets they created
- **Service Providers**: See tickets assigned to them
- **Manufacturers**: See tickets for their equipment

**Implementation:**

```go
func (r *TicketRepository) GetAll(ctx context.Context) ([]*domain.Ticket, error) {
    orgID, ok := middleware.GetOrganizationID(ctx)
    if !ok {
        return nil, fmt.Errorf("organization context not found")
    }
    
    orgType, _ := middleware.GetOrganizationType(ctx)
    
    var query string
    
    switch orgType {
    case "manufacturer":
        // Manufacturers see tickets for their equipment
        query = `
            SELECT t.* FROM service_tickets t
            JOIN equipment_registry e ON t.equipment_id = e.id
            WHERE e.manufacturer_id = $1
            ORDER BY t.created_at DESC
        `
    case "hospital", "imaging_center":
        // Hospitals see tickets they created
        query = `
            SELECT * FROM service_tickets
            WHERE requester_org_id = $1
            ORDER BY created_at DESC
        `
    case "Channel Partner", "Sub-sub_SUB_DEALER":
        // Service providers see tickets assigned to them
        query = `
            SELECT * FROM service_tickets
            WHERE assigned_org_id = $1
            ORDER BY created_at DESC
        `
    default:
        query = `
            SELECT * FROM service_tickets
            WHERE requester_org_id = $1
            ORDER BY created_at DESC
        `
    }
    
    var tickets []*domain.Ticket
    err := r.db.SelectContext(ctx, &tickets, query, orgID)
    return tickets, err
}
```

---

## Task 2.3: Update Engineers Repository â­

**Objective**: Filter engineers by organization.

**Files to Modify:**
- `internal/service-domain/service-engineer/infra/engineer_repository.go`

**Implementation:**

```go
func (r *EngineerRepository) GetAll(ctx context.Context) ([]*domain.Engineer, error) {
    orgID, ok := middleware.GetOrganizationID(ctx)
    if !ok {
        return nil, fmt.Errorf("organization context not found")
    }
    
    // Engineers belong to the organization through engineer_org_memberships
    query := `
        SELECT e.* FROM engineers e
        JOIN engineer_org_memberships eom ON e.id = eom.engineer_id
        WHERE eom.org_id = $1
        AND eom.status = 'active'
        ORDER BY e.name
    `
    
    var engineers []*domain.Engineer
    err := r.db.SelectContext(ctx, &engineers, query, orgID)
    return engineers, err
}
```

---

## Task 2.4: Create Organization Filter Helper â­

**Objective**: Reusable helper to build organization-filtered queries.

**Files to Create:**
- `internal/pkg/orgfilter/query_builder.go` (NEW)

**Implementation:**

```go
package orgfilter

import (
    "context"
    "fmt"
    
    "github.com/ServQR/medical-platform/internal/middleware"
    "github.com/google/uuid"
)

type OrgContext struct {
    OrgID   uuid.UUID
    OrgType string
}

// GetOrgContext extracts organization context from request context
func GetOrgContext(ctx context.Context) (*OrgContext, error) {
    orgID, ok := middleware.GetOrganizationID(ctx)
    if !ok {
        return nil, fmt.Errorf("organization ID not found in context")
    }
    
    orgType, _ := middleware.GetOrganizationType(ctx)
    
    return &OrgContext{
        OrgID:   orgID,
        OrgType: orgType,
    }, nil
}

// BuildEquipmentFilter returns WHERE clause for equipment queries
func BuildEquipmentFilter(orgType string) string {
    switch orgType {
    case "manufacturer":
        return "manufacturer_id = $1"
    case "hospital", "imaging_center":
        return "(organization_id = $1 OR owner_org_id = $1)"
    case "Channel Partner", "Sub-sub_SUB_DEALER":
        return "(channel_partner_org_id = $1 OR service_provider_org_id = $1)"
    default:
        return "organization_id = $1"
    }
}

// BuildTicketFilter returns WHERE clause for ticket queries
func BuildTicketFilter(orgType string) string {
    switch orgType {
    case "hospital", "imaging_center":
        return "requester_org_id = $1"
    case "Channel Partner", "Sub-sub_SUB_DEALER":
        return "assigned_org_id = $1"
    default:
        return "requester_org_id = $1"
    }
}
```

---

# Phase 3: Frontend Context

## Task 3.1: Decode JWT and Store Organization Context â­

**Objective**: Extract organization information from JWT token in frontend.

**Files to Create/Modify:**
- `admin-ui/src/lib/auth/jwt-decoder.ts` (NEW)
- `admin-ui/src/contexts/AuthContext.tsx` (MODIFY)

**Implementation:**

```typescript
// admin-ui/src/lib/auth/jwt-decoder.ts
export interface JWTPayload {
  user_id: string;
  email: string;
  name: string;
  organization_id: string;
  organization_type: 'manufacturer' | 'hospital' | 'Channel Partner' | 'Sub-sub_SUB_DEALER' | 'supplier' | 'imaging_center';
  role: string;
  permissions: string[];
  exp: number;
  iat: number;
}

export function decodeJWT(token: string): JWTPayload | null {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    
    return JSON.parse(jsonPayload) as JWTPayload;
  } catch (error) {
    console.error('Failed to decode JWT:', error);
    return null;
  }
}

export function isTokenExpired(token: string): boolean {
  const payload = decodeJWT(token);
  if (!payload) return true;
  
  return Date.now() >= payload.exp * 1000;
}
```

```typescript
// admin-ui/src/contexts/AuthContext.tsx
import { decodeJWT, type JWTPayload } from '@/lib/auth/jwt-decoder';

interface AuthContextType {
  user: JWTPayload | null;
  organization: {
    id: string;
    type: string;
  } | null;
  login: (accessToken: string, refreshToken: string) => void;
  logout: () => void;
  isAuthenticated: boolean;
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<JWTPayload | null>(null);
  
  const login = (accessToken: string, refreshToken: string) => {
    localStorage.setItem('access_token', accessToken);
    localStorage.setItem('refresh_token', refreshToken);
    
    const payload = decodeJWT(accessToken);
    setUser(payload);
  };
  
  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (token) {
      const payload = decodeJWT(token);
      setUser(payload);
    }
  }, []);
  
  const organization = user ? {
    id: user.organization_id,
    type: user.organization_type
  } : null;
  
  return (
    <AuthContext.Provider value={{ 
      user, 
      organization,
      login, 
      logout, 
      isAuthenticated: !!user 
    }}>
      {children}
    </AuthContext.Provider>
  );
}
```

**Testing:**
- [ ] JWT is decoded successfully
- [ ] Organization context is available in useAuth()
- [ ] Organization type is correct for each user

---

## Task 3.2: Add Organization Info to API Client â­

**Objective**: Include organization context in API headers.

**Files to Modify:**
- `admin-ui/src/lib/api/client.ts`

**Implementation:**

```typescript
// admin-ui/src/lib/api/client.ts

apiClient.interceptors.request.use(
  (config) => {
    // Add auth token
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
      
      // Decode and add organization header
      const payload = decodeJWT(token);
      if (payload?.organization_id) {
        config.headers['X-Organization-ID'] = payload.organization_id;
        config.headers['X-Organization-Type'] = payload.organization_type;
      }
    }
    
    return config;
  },
  (error) => Promise.reject(error)
);
```

---

# Phase 4: Organization-Specific UI

## Task 4.1: Create Organization-Specific Dashboards â­â­

**Objective**: Different dashboard layouts for each organization type.

**Files to Create:**
- `admin-ui/src/components/dashboards/ManufacturerDashboard.tsx`
- `admin-ui/src/components/dashboards/HospitalDashboard.tsx`
- `admin-ui/src/components/dashboards/ChannelPartnerDashboard.tsx`
- `admin-ui/src/app/dashboard/page.tsx` (MODIFY)

**Implementation:**

```typescript
// admin-ui/src/app/dashboard/page.tsx
'use client';

import { useAuth } from '@/contexts/AuthContext';
import ManufacturerDashboard from '@/components/dashboards/ManufacturerDashboard';
import HospitalDashboard from '@/components/dashboards/HospitalDashboard';
import ChannelPartnerDashboard from '@/components/dashboards/ChannelPartnerDashboard';

export default function DashboardPage() {
  const { organization } = useAuth();
  
  if (!organization) {
    return <div>Loading...</div>;
  }
  
  switch (organization.type) {
    case 'manufacturer':
      return <ManufacturerDashboard />;
    case 'hospital':
    case 'imaging_center':
      return <HospitalDashboard />;
    case 'Channel Partner':
    case 'Sub-sub_SUB_DEALER':
      return <ChannelPartnerDashboard />;
    default:
      return <DefaultDashboard />;
  }
}
```

**Manufacturer Dashboard:**
```typescript
// Show: Equipment catalog, service history, warranty management
export default function ManufacturerDashboard() {
  return (
    <div>
      <h1>Manufacturer Dashboard</h1>
      <StatsGrid>
        <StatCard title="Equipment Manufactured" value="1,234" />
        <StatCard title="Active Service Tickets" value="56" />
        <StatCard title="Warranty Claims" value="12" />
      </StatsGrid>
      
      <RecentServiceTickets />
      <EquipmentCatalog />
    </div>
  );
}
```

**Hospital Dashboard:**
```typescript
// Show: Owned equipment, service requests, maintenance schedules
export default function HospitalDashboard() {
  return (
    <div>
      <h1>Hospital Dashboard</h1>
      <StatsGrid>
        <StatCard title="Equipment Inventory" value="456" />
        <StatCard title="Open Tickets" value="23" />
        <StatCard title="Upcoming Maintenance" value="8" />
      </StatsGrid>
      
      <EquipmentInventory />
      <ServiceRequests />
    </div>
  );
}
```

---

## Task 4.2: Implement Conditional Navigation â­

**Objective**: Show/hide menu items based on organization type.

**Files to Modify:**
- `admin-ui/src/components/Sidebar.tsx` or navigation component

**Implementation:**

```typescript
export default function Navigation() {
  const { organization } = useAuth();
  
  return (
    <nav>
      <NavItem href="/dashboard">Dashboard</NavItem>
      
      {/* Manufacturer-specific */}
      {organization?.type === 'manufacturer' && (
        <>
          <NavItem href="/product-catalog">Product Catalog</NavItem>
          <NavItem href="/warranty-claims">Warranty Claims</NavItem>
        </>
      )}
      
      {/* Hospital-specific */}
      {['hospital', 'imaging_center'].includes(organization?.type) && (
        <>
          <NavItem href="/equipment">Equipment Inventory</NavItem>
          <NavItem href="/service-tickets">Service Requests</NavItem>
          <NavItem href="/maintenance">Maintenance Schedule</NavItem>
        </>
      )}
      
      {/* Channel Partner-specific */}
      {['Channel Partner', 'Sub-sub_SUB_DEALER'].includes(organization?.type) && (
        <>
          <NavItem href="/sales">Sales</NavItem>
          <NavItem href="/installations">Installations</NavItem>
          <NavItem href="/service-assignments">Service Assignments</NavItem>
        </>
      )}
      
      {/* Common for all */}
      <NavItem href="/engineers">Engineers</NavItem>
      <NavItem href="/reports">Reports</NavItem>
    </nav>
  );
}
```

---

## Task 4.3: Create Organization Indicator Component

**Objective**: Show which organization user is currently acting as.

**Files to Create:**
- `admin-ui/src/components/OrganizationBadge.tsx`

**Implementation:**

```typescript
export default function OrganizationBadge() {
  const { user, organization } = useAuth();
  
  const orgTypeColors = {
    manufacturer: 'bg-blue-100 text-blue-800',
    hospital: 'bg-green-100 text-green-800',
    Channel Partner: 'bg-purple-100 text-purple-800',
    Sub-sub_SUB_DEALER: 'bg-orange-100 text-orange-800',
  };
  
  const orgTypeIcons = {
    manufacturer: 'ðŸ­',
    hospital: 'ðŸ¥',
    Channel Partner: 'ðŸšš',
    Sub-sub_SUB_DEALER: 'ðŸª',
  };
  
  return (
    <div className="flex items-center gap-2 px-3 py-2 bg-gray-50 rounded-lg">
      <span className="text-2xl">
        {orgTypeIcons[organization?.type] || 'ðŸ¢'}
      </span>
      <div>
        <div className="text-xs text-gray-500">Logged in as</div>
        <div className="font-medium">{user?.name}</div>
        <span className={`text-xs px-2 py-0.5 rounded ${orgTypeColors[organization?.type]}`}>
          {organization?.type}
        </span>
      </div>
    </div>
  );
}
```

---

# Phase 5: Testing & Validation

## Task 5.1: Backend Integration Tests â­â­

**Objective**: Verify data isolation at API level.

**Files to Create:**
- `tests/integration/multi_tenant_test.go`

**Test Cases:**

```go
func TestEquipmentIsolation(t *testing.T) {
    // Login as manufacturer
    manufacturerToken := loginAsManufacturer()
    
    // Get equipment
    manufacturerEquipment := getEquipment(manufacturerToken)
    
    // Login as hospital
    hospitalToken := loginAsHospital()
    
    // Get equipment
    hospitalEquipment := getEquipment(hospitalToken)
    
    // Verify no overlap
    assert.NoOverlap(manufacturerEquipment, hospitalEquipment)
}

func TestTicketIsolation(t *testing.T) {
    // Similar test for tickets
}

func TestCrossOrganizationAccess(t *testing.T) {
    // Try to access another org's data
    // Should return 403 or empty results
}
```

**Checklist:**
- [ ] Manufacturer sees only their equipment
- [ ] Hospital sees only their equipment
- [ ] No equipment ID overlap between orgs
- [ ] Tickets are properly isolated
- [ ] Engineers are properly isolated
- [ ] Cross-org access returns 403/empty

---

## Task 5.2: Frontend Manual Testing â­

**Test Script:**

1. **Test Manufacturer Login:**
   ```
   - Login: manufacturer@geneqr.com / password
   - Check dashboard shows manufacturer view
   - Check equipment list (should show Siemens equipment)
   - Check navigation shows manufacturer menu items
   - Logout
   ```

2. **Test Hospital Login:**
   ```
   - Login: hospital@geneqr.com / password
   - Check dashboard shows hospital view
   - Check equipment list (should be different from manufacturer)
   - Check navigation shows hospital menu items
   - Logout
   ```

3. **Test Data Isolation:**
   ```
   - Record equipment IDs from manufacturer view
   - Login as hospital
   - Verify those equipment IDs are NOT visible
   - Success! Data isolation working
   ```

**Checklist:**
- [ ] Different dashboards for different org types
- [ ] Different navigation menus
- [ ] Different data sets
- [ ] Organization badge shows correct org
- [ ] JWT token decoded correctly

---

## Task 5.3: Security Testing â­â­

**Objective**: Ensure no data leakage through various attack vectors.

**Test Cases:**

1. **Direct API Call with Wrong Org Token:**
   ```bash
   # Get manufacturer token
   MANUFACTURER_TOKEN=$(curl ...)
   
   # Try to access hospital equipment
   curl -H "Authorization: Bearer $MANUFACTURER_TOKEN" \
        http://localhost:8081/api/v1/equipment?org_id=<hospital_org_id>
   
   # Should return empty or 403
   ```

2. **Token Manipulation:**
   - Try to modify organization_id in JWT
   - Should be rejected (signature invalid)

3. **SQL Injection:**
   - Try org_id with SQL injection payload
   - Should be sanitized by parameterized queries

**Checklist:**
- [ ] Cannot access other org's data with valid token
- [ ] Cannot modify JWT organization_id
- [ ] SQL injection attempts fail
- [ ] Rate limiting works per organization

---

# Implementation Schedule

## Week 1: Backend Foundation
- **Day 1-2**: Task 1.1 - Organization context middleware
- **Day 3**: Task 1.2 - Add org_type to JWT
- **Day 4**: Task 1.3 - Organization repository
- **Day 5**: Testing Phase 1

## Week 2: API Data Filtering
- **Day 1-2**: Task 2.1 - Equipment repository filters
- **Day 3**: Task 2.2 - Ticket repository filters
- **Day 4**: Task 2.3 - Engineer repository filters
- **Day 5**: Task 2.4 - Query builder helper + Testing

## Week 3: Frontend Implementation
- **Day 1-2**: Task 3.1 - JWT decoder & auth context
- **Day 3**: Task 3.2 - API client updates
- **Day 4-5**: Task 4.1 - Organization-specific dashboards

## Week 4: UI Polish & Testing
- **Day 1-2**: Task 4.2 - Conditional navigation
- **Day 3**: Task 4.3 - Organization badge
- **Day 4-5**: Phase 5 - All testing tasks

---

# Success Criteria

## Must Have (P0)
- âœ… Users can only see data belonging to their organization
- âœ… JWT tokens include organization context
- âœ… All repository queries filter by organization
- âœ… Different dashboards for different org types
- âœ… Security testing passes

## Should Have (P1)
- âœ… Organization badge showing current org
- âœ… Conditional navigation based on org type
- âœ… Proper error handling for missing org context
- âœ… Logging of organization access patterns

## Nice to Have (P2)
- Organization switcher for multi-org users
- Organization settings page
- Organization analytics
- Audit trail per organization

---

# Rollout Plan

## Phase A: Development (Weeks 1-4)
- Implement all tasks above
- Test in development environment
- Fix bugs and edge cases

## Phase B: Staging Testing (Week 5)
- Deploy to staging
- Test with real-like data
- Security audit
- Performance testing

## Phase C: Production Rollout (Week 6)
- Deploy backend changes
- Deploy frontend changes
- Monitor for issues
- Rollback plan ready

## Phase D: Validation (Week 7)
- User acceptance testing
- Performance monitoring
- Security monitoring
- Bug fixes

---

# Risk Mitigation

## Risk: Data Leakage
**Mitigation**: 
- Comprehensive testing before deployment
- Code review of all repository changes
- Security audit

## Risk: Performance Degradation
**Mitigation**:
- Add database indexes on organization_id columns
- Query optimization
- Caching strategy

## Risk: Breaking Existing Features
**Mitigation**:
- Incremental rollout
- Feature flags
- Rollback plan

---

# Next Steps

1. **Review this plan** - Ensure all requirements are covered
2. **Prioritize tasks** - Identify must-have vs nice-to-have
3. **Start with Task 1.1** - Organization context middleware
4. **Test incrementally** - Test each task before moving to next
5. **Iterate** - Adjust plan based on learnings

---

# Questions to Answer Before Starting

- [ ] Do we have all required database indexes on org columns?
- [ ] Do we need organization-level rate limiting?
- [ ] Should we support organization switching in same session?
- [ ] What's the rollback strategy if we find issues?
- [ ] Do we need audit logs for organization access?

---

**Let's start with Task 1.1: Organization Context Middleware!**

Ready to begin? ðŸš€
