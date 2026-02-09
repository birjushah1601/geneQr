-- Migration: Add timeline override columns to service_tickets
-- Date: 2026-02-08
-- Purpose: Store admin-adjusted timelines and milestone data

-- Add timeline_overrides column to store custom milestone adjustments
ALTER TABLE service_tickets 
ADD COLUMN IF NOT EXISTS timeline_overrides JSONB;

-- Add parts_override column to store custom parts status/ETA
ALTER TABLE service_tickets 
ADD COLUMN IF NOT EXISTS parts_override JSONB;

-- Add index for faster queries
CREATE INDEX IF NOT EXISTS idx_service_tickets_timeline_overrides 
ON service_tickets USING GIN (timeline_overrides);

-- Add comment
COMMENT ON COLUMN service_tickets.timeline_overrides IS 'Admin-adjusted milestone data (JSON array of PublicMilestone)';
COMMENT ON COLUMN service_tickets.parts_override IS 'Admin-adjusted parts status and ETA (JSON object)';
