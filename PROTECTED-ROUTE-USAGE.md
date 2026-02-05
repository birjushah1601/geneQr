# Protected Route Usage Guide

## Overview
The `ProtectedRoute` component ensures only authenticated users can access certain pages.

## Features
- ✅ Token validation
- ✅ Expiry check (30-second buffer)
- ✅ Auto-redirect to /login
- ✅ Path preservation for post-login redirect

## Usage
```typescript
import ProtectedRoute from '@/components/ProtectedRoute';

export default function MyPage() {
  return (
    <ProtectedRoute>
      <MyContent />
    </ProtectedRoute>
  );
}
```

## Pages That Should Be Protected
- /dashboard
- /equipment
- /tickets
- /engineers

## Pages That Should Stay Public
- /service-request (QR scans - MUST be public!)
- /login
- /register
- /invite/accept
