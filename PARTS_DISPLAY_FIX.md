# Parts Display in Tickets - FIXED

## Issue
Parts were not showing in ticket details modal even though we added 22 equipment-parts assignments.

## Root Causes

1. **Wrong Table**: Function used equipment_spare_parts (17 old records) instead of equipment_part_assignments (22 new records)
2. **Missing Links**: Equipment records had NULL catalog_id, so couldn't find parts

## Fixes Applied

### 1. Updated get_parts_for_registry Function
- Now checks equipment_part_assignments table first
- Handles both equipment_registry and equipment tables
- Falls back to equipment_spare_parts for legacy data

### 2. Linked Equipment to Catalog
Updated 6 equipment records to link to their catalog entries:
- 2x Infusion Pump Lite (2 parts each)
- 2x CT Scanner Nova (2 parts each)  
- 2x X-Ray System Alpha (2 parts each)

## Test Results

Ticket 36mNfDnkNPKyqzKYRujS4R4CyC2 (X-Ray System Alpha):
- Flat Panel Detector (45,000 USD) - Critical
- X-Ray Tube Standard (18,000 USD) - Critical

## Equipment with Parts Now

6 equipment records now show parts:
- Infusion Pump Lite: 2 parts
- CT Scanner Nova: 2 parts
- X-Ray System Alpha: 2 parts

## Files Modified
1. database/migrations/fix_parts_function.sql
2. database/migrations/link_equipment_to_catalog.sql

## Status
COMPLETE - Parts now display for 6 equipment types in tickets
