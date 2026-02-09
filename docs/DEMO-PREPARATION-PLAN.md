# Demo Preparation Implementation Plan

## Overview
Implementing feature flags and fixes for client demo on Feb 9, 2026.

## Feature Flag System
Created: `admin-ui/src/lib/featureFlags.ts`
- URL-based feature toggling
- Usage: `?enable=FeatureName` or `?enable=Feature1,Feature2`

## Requirements & Implementation

### 1. Hide Phone/OTP Login ✅
**Status:** Ready to implement
**Files:** `admin-ui/src/app/login/page.tsx`
**Implementation:**
- Hide "Use OTP instead" button by default
- Show only with `?enable=PhoneLogin`
- Default to password-only login

### 2. Hide AI Onboarding
**Status:** Ready to implement
**Files:**
- `admin-ui/src/app/dashboard/page.tsx` (or layout)
- Navigation component (sidebar/header)
**Implementation:**
- Hide "AI Onboarding" link from navigation
- Show only with `?enable=AIOnboarding`
- Works for both manufacturer and admin dashboards

### 3. Hide Equipment Add/Import Buttons
**Status:** Ready to implement
**Files:** `admin-ui/src/app/equipment/page.tsx`
**Implementation:**
- Hide "Add Equipment" button
- Hide "Import CSV" button
- Show with `?enable=AddNew`

### 4. Hide Service Tickets Create Button
**Status:** Ready to implement
**Files:** `admin-ui/src/app/tickets/page.tsx`
**Implementation:**
- Hide "Create Ticket" button
- Show with `?enable=AddNew`

### 5. Hide Engineer Add/Import Buttons
**Status:** Ready to implement
**Files:** `admin-ui/src/app/engineers/page.tsx`
**Implementation:**
- Hide "Add Engineer" button
- Hide "Import CSV" button
- Show with `?enable=AddNew`

### 6. Fix QR PDF Preview URL
**Status:** Requires investigation
**Files:** QR code generation/PDF preview component
**Implementation:**
- Change equipment details URL to service request URL
- Update message: "Use this URL to create a Service Request if QR Code is not functioning"

### 7. Hide AI Diagnosis on Service Request Page
**Status:** Ready to implement
**Files:** Service request creation page
**Implementation:**
- Hide "AI Powered Diagnosis" section
- Could keep hidden permanently or use feature flag

### 8. Update Track URL to HTTPS
**Status:** Ready to implement
**Files:** Tracking URL generation
**Implementation:**
- Change from http:// to https://servqr.com

### 9. Fix SLA Business Hours
**Status:** Requires backend changes
**Files:** `internal/service-domain/service-ticket/app/timeline_service.go`
**Implementation:**
- Business hours: 9 AM - 6 PM IST
- SLA: 2 days (48 business hours)
- If ticket created after hours, start SLA from next 9 AM
- Acknowledgment can happen anytime, SLA starts from acknowledgment

### 10. Fix Milestone Adjustment Logic
**Status:** Requires backend changes
**Files:** `admin-ui/src/components/TimelineEditModal.tsx`, backend timeline service
**Implementation:**
- When intermediate milestone pushed later
- Auto-adjust remaining milestones (3-5 hours each)
- Ensure target completion is after last milestone

### 11. Fix Team Invite Email Link
**Status:** Requires backend changes
**Files:** Email template, invitation service
**Implementation:**
- Use direct link (not token-based redirect)
- Use https://servqr.com domain

## Priority Order
1. **High Priority (Quick Wins):**
   - Hide Phone/OTP Login
   - Hide AI Onboarding
   - Hide Add/Import Buttons (Equipment, Tickets, Engineers)
   - Hide AI Diagnosis
   - Update Track URL

2. **Medium Priority:**
   - Fix QR PDF Preview URL
   - Fix Team Invite Email Link

3. **Low Priority (Time-Consuming):**
   - Fix SLA Business Hours Logic
   - Fix Milestone Adjustment Logic

## Testing Checklist
- [ ] Login page shows only password option
- [ ] Navigation hides AI Onboarding
- [ ] Equipment page hides Add/Import
- [ ] Tickets page hides Create button
- [ ] Engineers page hides Add/Import
- [ ] Feature flags work with URL parameters
- [ ] Track URL uses HTTPS
- [ ] QR PDF shows correct URL
- [ ] Invite emails have correct link

## Rollout Strategy
1. Implement all feature flags first (frontend only)
2. Test with ?enable parameters
3. Implement URL fixes (quick backend changes)
4. Leave complex logic (SLA, milestones) for post-demo if needed
