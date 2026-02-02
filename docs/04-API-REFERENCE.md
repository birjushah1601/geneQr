# ServQR API Reference

Quick reference for all REST API endpoints.

## ðŸŒ Base URL

- **Development:** `http://localhost:8081`
- **Production:** `https://api.ServQR.com` (when deployed)

## ðŸ”‘ Authentication

Most endpoints require authentication via JWT token in cookies or Authorization header.

```http
Authorization: Bearer <jwt_token>
X-User-Role: admin
X-User-ID: user-uuid
```

---

## ðŸ“‹ API Modules

### 1. Tickets API

**Base:** `/api/v1/tickets`

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/` | Create ticket | Required |
| GET | `/` | List tickets | Required |
| GET | `/{id}` | Get ticket details | Required |
| GET | `/number/{number}` | Get by ticket number | Required |
| POST | `/{id}/assign` | Assign engineer | Admin |
| PATCH | `/{id}/priority` | Update priority | Admin only |
| POST | `/{id}/start` | Start work | Engineer |
| POST | `/{id}/resolve` | Resolve ticket | Engineer |
| POST | `/{id}/close` | Close ticket | Admin |
| GET | `/{id}/comments` | Get comments | Required |
| POST | `/{id}/comments` | Add comment | Required |

**Details:** See [Postman collection](../postman/)

### 2. Equipment API

**Base:** `/api/v1/equipment`

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/` | Register equipment | Required |
| GET | `/` | List equipment | Required |
| GET | `/{id}` | Get details | Required |
| GET | `/qr/{qrCode}` | Get by QR code | Public |

### 3. Organizations API

**Base:** `/api/v1/organizations`

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/` | Create organization | Admin |
| GET | `/` | List organizations | Required |
| GET | `/{id}` | Get details | Required |
| POST | `/import` | Bulk CSV import | Admin |

### 4. Engineers API

**Base:** `/api/v1/engineers`

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/` | List engineers | Required |
| GET | `/{id}` | Get details | Required |
| PUT | `/{id}/level` | Update level | Admin |

### 5. WhatsApp API

**Base:** `/api/v1/whatsapp`

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/webhook` | Verify webhook | Public |
| POST | `/webhook` | Handle messages | Public |

### 6. Parts API

**Base:** `/api/v1/parts`

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/catalog` | List parts | Required |
| POST | `/catalog` | Add part | Admin |
| GET | `/{id}` | Get part details | Required |

### 7. AI Diagnosis API

**Base:** `/api/v1/diagnosis`

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/` | Create diagnosis | Required |
| GET | `/{id}` | Get diagnosis | Required |
| POST | `/{id}/feedback` | Submit feedback | Required |

---

## ðŸ“ Example Requests

### Create Ticket
```http
POST /api/v1/tickets
Content-Type: application/json

{
  "equipment_id": "uuid",
  "qr_code": "QR-20251223-001",
  "issue_description": "Equipment not starting",
  "priority": "medium",
  "source": "web"
}
```

### Update Priority (Admin Only)
```http
PATCH /api/v1/tickets/{id}/priority
Content-Type: application/json
X-User-Role: admin

{
  "priority": "critical"
}
```

### Bulk Import Organizations
```http
POST /api/v1/organizations/import
Content-Type: multipart/form-data

csv_file: <file>
dry_run: false
```

---

## ðŸ“š Full Documentation

See subdirectories:
- `/docs/api/` - Detailed API specs
- `/postman/` - Postman collections

**Last Updated:** December 23, 2025
