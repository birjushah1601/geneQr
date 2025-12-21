-- ============================================================================
-- Create ticket_parts Table
-- ============================================================================
-- Junction table to store which parts are assigned to which tickets
-- This is the correct place to store part assignments, not in JSONB field

BEGIN;

-- Drop table if exists (for clean creation)
DROP TABLE IF EXISTS ticket_parts CASCADE;

-- Create ticket_parts table
CREATE TABLE ticket_parts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Foreign keys
    ticket_id VARCHAR(50) NOT NULL,
    spare_part_id UUID NOT NULL,
    
    -- Assignment details
    quantity_required INTEGER NOT NULL DEFAULT 1,
    quantity_used INTEGER,
    is_critical BOOLEAN DEFAULT false,
    
    -- Status tracking
    status VARCHAR(50) DEFAULT 'pending', -- pending, ordered, received, installed, returned
    
    -- Pricing
    unit_price NUMERIC(12,2),
    total_price NUMERIC(12,2),
    currency VARCHAR(3) DEFAULT 'USD',
    
    -- Notes and tracking
    notes TEXT,
    assigned_by VARCHAR(100),
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    installed_at TIMESTAMP WITH TIME ZONE,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT fk_ticket_parts_ticket FOREIGN KEY (ticket_id) REFERENCES service_tickets(id) ON DELETE CASCADE,
    CONSTRAINT fk_ticket_parts_part FOREIGN KEY (spare_part_id) REFERENCES spare_parts_catalog(id) ON DELETE RESTRICT,
    CONSTRAINT ticket_parts_quantity_check CHECK (quantity_required > 0),
    CONSTRAINT ticket_parts_status_check CHECK (status IN ('pending', 'ordered', 'received', 'installed', 'returned', 'cancelled'))
);

-- Indexes for performance
CREATE INDEX idx_ticket_parts_ticket ON ticket_parts(ticket_id);
CREATE INDEX idx_ticket_parts_part ON ticket_parts(spare_part_id);
CREATE INDEX idx_ticket_parts_status ON ticket_parts(status);
CREATE INDEX idx_ticket_parts_assigned_at ON ticket_parts(assigned_at DESC);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_ticket_parts_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_ticket_parts_timestamp
    BEFORE UPDATE ON ticket_parts
    FOR EACH ROW
    EXECUTE FUNCTION update_ticket_parts_timestamp();

-- Comments for documentation
COMMENT ON TABLE ticket_parts IS 'Junction table storing which spare parts are assigned to which service tickets';
COMMENT ON COLUMN ticket_parts.ticket_id IS 'Reference to service ticket';
COMMENT ON COLUMN ticket_parts.spare_part_id IS 'Reference to spare part from catalog';
COMMENT ON COLUMN ticket_parts.quantity_required IS 'Number of parts needed';
COMMENT ON COLUMN ticket_parts.quantity_used IS 'Actual number of parts used (may differ from required)';
COMMENT ON COLUMN ticket_parts.is_critical IS 'Whether this part is critical for the repair';
COMMENT ON COLUMN ticket_parts.status IS 'Part assignment status: pending, ordered, received, installed, returned, cancelled';
COMMENT ON COLUMN ticket_parts.unit_price IS 'Price per unit at time of assignment';
COMMENT ON COLUMN ticket_parts.total_price IS 'Total price for all units';

-- Create view for easy querying
CREATE OR REPLACE VIEW v_ticket_parts_detail AS
SELECT 
    tp.id as assignment_id,
    tp.ticket_id,
    t.ticket_number,
    t.equipment_name,
    tp.spare_part_id,
    sp.part_number,
    sp.part_name,
    sp.category as part_category,
    tp.quantity_required,
    tp.quantity_used,
    tp.is_critical,
    tp.status,
    tp.unit_price,
    tp.total_price,
    tp.currency,
    tp.notes,
    tp.assigned_by,
    tp.assigned_at,
    tp.installed_at,
    sp.stock_status,
    sp.lead_time_days,
    t.status as ticket_status,
    t.priority as ticket_priority
FROM ticket_parts tp
JOIN service_tickets t ON tp.ticket_id = t.id
JOIN spare_parts_catalog sp ON tp.spare_part_id = sp.id
ORDER BY tp.assigned_at DESC;

COMMENT ON VIEW v_ticket_parts_detail IS 'Detailed view of ticket parts with related ticket and part information';

-- Show summary
SELECT 
    '=== ticket_parts Table Created ===' as status;

SELECT 
    'Table: ticket_parts' as object_type,
    COUNT(*) as column_count
FROM information_schema.columns 
WHERE table_name = 'ticket_parts' AND table_schema = 'public';

SELECT 
    'Indexes created: ' || COUNT(*) as index_info
FROM pg_indexes 
WHERE tablename = 'ticket_parts';

SELECT 
    'Ready to store part assignments' as ready_status;

COMMIT;

-- Verification
SELECT '✅ ticket_parts table created successfully. Parts can now be assigned to tickets.' as final_status;
