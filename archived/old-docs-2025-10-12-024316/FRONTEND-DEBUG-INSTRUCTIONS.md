# 🐛 Frontend API Debug Instructions

## Problem
Service-request page shows: `Equipment not found for QR code: QR-eq-001`

## Backend Status ✅
```
✓ Backend API is working perfectly!
✓ GET http://localhost:8081/api/v1/equipment/qr/QR-eq-001
  Returns: HTTP 200 OK
  Data: { id: "eq-001", equipment_name: "X-Ray Machine", ... }
```

## Frontend Status ⚠️
```
⚠️ Frontend is running on port 3000
⚠️ API call is failing despite backend working
```

---

## Debug Steps

### Step 1: Check Browser Console
1. Open: http://localhost:3000/service-request?qr=QR-eq-001
2. Open Browser DevTools (F12)
3. Go to **Console** tab
4. Look for error messages (red text)
5. Check for:
   - Network errors
   - CORS errors
   - 404 errors
   - TypeErrors

### Step 2: Check Network Tab
1. Open Browser DevTools (F12)
2. Go to **Network** tab
3. Reload page
4. Find the request to `/api/v1/equipment/qr/QR-eq-001`
5. Check:
   - **Status Code:** Should be 200, is it 404? 500?
   - **Request URL:** Is it correct?
   - **Response:** Click on it, check "Response" tab

### Step 3: Possible Issues & Solutions

#### Issue 1: CORS Error
**Symptom:** Console shows "CORS policy" error  
**Solution:** Backend needs CORS headers

Add to Go backend's handler:
```go
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Tenant-ID")
```

#### Issue 2: Wrong Base URL
**Symptom:** Network tab shows wrong URL like `http://localhost:3000/api/...`  
**Solution:** Check `.env.local` file

Create `admin-ui/.env.local`:
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8081/api
```

Restart frontend:
```bash
cd admin-ui
npm run dev
```

#### Issue 3: Frontend Cache
**Symptom:** Old code still running  
**Solution:** Hard refresh

1. Press: `Ctrl + Shift + R` (Windows) or `Cmd + Shift + R` (Mac)
2. Or: Right-click refresh → "Empty Cache and Hard Reload"

#### Issue 4: API Client Issue
**Symptom:** Response is 200 but data not parsed  
**Solution:** Check response structure

Open `admin-ui/src/lib/api/equipment.ts` line 49:
```typescript
async getByQRCode(qrCode: string) {
  try {
    const response = await apiClient.get<Equipment>(`/v1/equipment/qr/${qrCode}`);
    console.log('🔍 API Response:', response);  // ADD THIS LINE
    console.log('🔍 Response Data:', response.data);  // ADD THIS LINE
    return response.data;
  } catch (error) {
    console.error('🔍 API Error:', error);  // ADD THIS LINE
    throw new Error(handleApiError(error));
  }
},
```

---

## Quick Test

### Test 1: Direct Backend Call
Open new browser tab:
```
http://localhost:8081/api/v1/equipment/qr/QR-eq-001
```

**Expected:** JSON data for X-Ray Machine  
**If fails:** Backend not running or wrong port

### Test 2: Frontend API Test
Open browser console on ANY page of frontend:
```javascript
fetch('http://localhost:8081/api/v1/equipment/qr/QR-eq-001')
  .then(r => r.json())
  .then(data => console.log('✅ Data:', data))
  .catch(err => console.error('❌ Error:', err))
```

**Expected:** Console shows equipment data  
**If CORS error:** Need to add CORS headers to backend  
**If 404:** Backend routing issue

### Test 3: Check API Base URL
Open browser console:
```javascript
console.log('Base URL:', process.env.NEXT_PUBLIC_API_BASE_URL)
```

**Expected:** `http://localhost:8081/api`  
**If undefined:** Need to create `.env.local` file

---

## Most Likely Issues (In Order)

### 1. Missing .env.local File (90% probability)
The frontend might not have the API base URL configured.

**Solution:**
```bash
cd admin-ui
echo "NEXT_PUBLIC_API_BASE_URL=http://localhost:8081/api" > .env.local
npm run dev
```

### 2. CORS Not Configured (60% probability)
Backend might not allow frontend origin.

**Check:** Browser console will show CORS error in red

**Solution:** Update backend CORS middleware

### 3. Frontend Cache (30% probability)
Old code still running with old API client.

**Solution:** Hard refresh (Ctrl + Shift + R)

---

## Expected Working Flow

```
┌────────────────────────────────────────────────┐
│ 1. User opens:                                 │
│    http://localhost:3000/service-request?qr=...│
└────────────────────┬───────────────────────────┘
                     ↓
┌────────────────────────────────────────────────┐
│ 2. Page loads, useEffect runs                  │
│    Calls: equipmentApi.getByQRCode(qrCode)    │
└────────────────────┬───────────────────────────┘
                     ↓
┌────────────────────────────────────────────────┐
│ 3. API client makes request:                   │
│    GET http://localhost:8081/api/v1/equipment/ │
│        qr/QR-eq-001                            │
└────────────────────┬───────────────────────────┘
                     ↓
┌────────────────────────────────────────────────┐
│ 4. Backend returns:                            │
│    { id: "eq-001", equipment_name: "...", ... }│
└────────────────────┬───────────────────────────┘
                     ↓
┌────────────────────────────────────────────────┐
│ 5. Frontend displays equipment details         │
│    Shows form for service request              │
└────────────────────────────────────────────────┘
```

---

## Testing Checklist

- [ ] Backend API returns 200 OK (✅ Already tested - WORKING)
- [ ] Database has QR codes (✅ Already tested - 4/4 equipment)
- [ ] Frontend running on port 3000
- [ ] `.env.local` file exists with correct API URL
- [ ] Browser console shows no CORS errors
- [ ] Network tab shows request to correct URL
- [ ] Network tab shows 200 OK response
- [ ] Response data structure matches Equipment type

---

## Need Help?

**Share these from browser:**
1. Screenshot of Console tab (F12 → Console)
2. Screenshot of Network tab showing the API request
3. Response from this test:
   ```
   http://localhost:8081/api/v1/equipment/qr/QR-eq-001
   ```

---

**Last Updated:** October 11, 2025, 9:30 PM IST  
**Backend:** ✅ WORKING  
**Frontend:** ⚠️ NEEDS DEBUG
