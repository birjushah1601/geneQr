# Backend Metadata JSON Fix

## Issue
The organizations API was returning metadata as base64-encoded string instead of JSON object.

## Root Cause
The `Organization` struct in the repository was using `[]byte` for the metadata field, which causes Go's JSON marshaller to base64-encode the bytes.

## Fix Applied

### File: `internal/core/organizations/infra/repository.go`

**Before:**
```go
type Organization struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    OrgType  string `json:"org_type"`
    Status   string `json:"status"`
    Metadata []byte `json:"metadata"`
}
```

**After:**
```go
import (
    "encoding/json"
    // ... other imports
)

type Organization struct {
    ID       string          `json:"id"`
    Name     string          `json:"name"`
    OrgType  string          `json:"org_type"`
    Status   string          `json:"status"`
    Metadata json.RawMessage `json:"metadata"`
}
```

## Result

**Before:**
```json
{
  "id": "uuid",
  "name": "Siemens Healthineers India",
  "metadata": "eyJlbWFpbCI6ICJyYWplc2gua3VtYXJAc2llbWVucy1oZWFsdGhpbmVlcnMuY29tIiwgInBob25lIjogIis5MS04MC00MTQxLTQxNDEiLCAiYWRkcmVzcyI6IHsiY2l0eSI6ICJHdWluZHksIENoZW5uYWkiLCAic3RhdGUiOiAiVGFtaWwgTmFkdSIsICJzdHJlZXQiOiAiT2x5bXBpYSBUZWNobm9sb2d5IFBhcmssIDEtQSwgU0lEQ08gSW5kdXN0cmlhbCBFc3RhdGUiLCAiY291bnRyeSI6ICJJbmRpYSIsICJwb3N0YWxfY29kZSI6ICI2MDAwMzIifSwgIndlYnNpdGUiOiAiaHR0cHM6Ly93d3cuc2llbWVucy1oZWFsdGhpbmVlcnMuY29tL2VuLWluIiwgInN1cHBvcnRfaW5mbyI6IHsic3VwcG9ydF9lbWFpbCI6ICJzZXJ2aWNlLmluZGlhQHNpZW1lbnMtaGVhbHRoaW5lZXJzLmNvbSIsICJzdXBwb3J0X2hvdXJzIjogIjI0LzcgQXZhaWxhYmxlIiwgInN1cHBvcnRfcGhvbmUiOiAiKzkxLTgwLTQxNDEtNDIwMCIsICJyZXNwb25zZV90aW1lX3NsYSI6ICI0IGhvdXJzIn0sICJidXNpbmVzc19pbmZvIjogeyJnc3RfbnVtYmVyIjogIjMzQUFDQ1MxMTE5RjFaNSIsICJwYW5fbnVtYmVyIjogIkFBQ0NTMTExOUYiLCAiaGVhZHF1YXJ0ZXJzIjogIk11bWJhaSwgTWFoYXJhc2h0cmEiLCAiZW1wbG95ZWVfY291bnQiOiA1MDAwLCAiZXN0YWJsaXNoZWRfeWVhciI6IDE5OTJ9LCAiY29udGFjdF9wZXJzb24iOiAiRHIuIFJhamVzaCBLdW1hciJ9"
}
```

**After:**
```json
{
  "id": "uuid",
  "name": "Siemens Healthineers India",
  "metadata": {
    "contact_person": "Dr. Rajesh Kumar",
    "email": "rajesh.kumar@siemens-healthineers.com",
    "phone": "+91-80-4141-4141",
    "website": "https://www.siemens-healthineers.com/en-in",
    "address": {
      "street": "Olympia Technology Park, 1-A, SIDCO Industrial Estate",
      "city": "Guindy, Chennai",
      "state": "Tamil Nadu",
      "postal_code": "600032",
      "country": "India"
    },
    "business_info": {
      "gst_number": "33AACCS1119F1Z5",
      "pan_number": "AACCS1119F",
      "established_year": 1992,
      "employee_count": 5000,
      "headquarters": "Mumbai, Maharashtra"
    },
    "support_info": {
      "support_email": "service.india@siemens-healthineers.com",
      "support_phone": "+91-80-4141-4200",
      "support_hours": "24/7 Available",
      "response_time_sla": "4 hours"
    }
  }
}
```

## Frontend Impact

**Frontend can now directly access metadata fields:**
```typescript
const contact = organization.metadata?.contact_person;
const email = organization.metadata?.email;
const city = organization.metadata?.address?.city;
const supportEmail = organization.metadata?.support_info?.support_email;
```

**No need for base64 decoding!**

## Status
✅ Backend rebuilt and restarted
✅ API returns metadata as JSON object
✅ Frontend can access metadata directly
✅ All 4 manufacturers returning proper data
