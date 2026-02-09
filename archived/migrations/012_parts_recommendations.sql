-- Migration: Parts Recommendations Table
-- Description: Stores AI-powered parts recommendations and feedback
-- Dependencies: service_tickets, parts_catalog tables
-- Phase: 2C - AI Services Layer

-- Create parts_recommendations table
CREATE TABLE IF NOT EXISTS parts_recommendations (
    recommendation_id BIGSERIAL PRIMARY KEY,
    
    -- Request identification
    request_id VARCHAR(100) UNIQUE NOT NULL,
    ticket_id BIGINT NOT NULL REFERENCES service_tickets(ticket_id) ON DELETE CASCADE,
    
    -- Recommendations data (full JSON)
    replacement_parts JSONB NOT NULL,
    accessories JSONB NOT NULL,
    preventive_parts JSONB NOT NULL,
    metadata JSONB NOT NULL,
    
    -- Usage tracking
    parts_used BIGINT[], -- Array of part IDs actually used
    accessories_sold BIGINT[], -- Array of accessory IDs sold
    
    -- Feedback
    was_accurate BOOLEAN,
    accuracy_feedback TEXT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create ticket_parts junction table (for tracking what parts were actually used)
CREATE TABLE IF NOT EXISTS ticket_parts (
    ticket_part_id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL REFERENCES service_tickets(ticket_id) ON DELETE CASCADE,
    part_id BIGINT NOT NULL REFERENCES parts_catalog(part_id) ON DELETE CASCADE,
    quantity_used INTEGER DEFAULT 1,
    was_recommended BOOLEAN DEFAULT false,
    cost DECIMAL(10,2),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(ticket_id, part_id)
);

-- Indexes
CREATE INDEX idx_parts_recommendations_ticket ON parts_recommendations(ticket_id);
CREATE INDEX idx_parts_recommendations_created ON parts_recommendations(created_at DESC);
CREATE INDEX idx_parts_recommendations_feedback ON parts_recommendations(was_accurate) WHERE was_accurate IS NOT NULL;

-- GIN indexes for JSONB
CREATE INDEX idx_parts_recommendations_replacement ON parts_recommendations USING GIN (replacement_parts);
CREATE INDEX idx_parts_recommendations_accessories ON parts_recommendations USING GIN (accessories);
CREATE INDEX idx_parts_recommendations_metadata ON parts_recommendations USING GIN (metadata);

-- Indexes for ticket_parts
CREATE INDEX idx_ticket_parts_ticket ON ticket_parts(ticket_id);
CREATE INDEX idx_ticket_parts_part ON ticket_parts(part_id);
CREATE INDEX idx_ticket_parts_recommended ON ticket_parts(was_recommended) WHERE was_recommended = true;

-- Trigger for updated_at
CREATE TRIGGER trg_update_parts_recommendations_timestamp
    BEFORE UPDATE ON parts_recommendations
    FOR EACH ROW
    EXECUTE FUNCTION update_parts_timestamp();

-- Comments
COMMENT ON TABLE parts_recommendations IS 'AI-powered parts recommendations for service tickets';
COMMENT ON COLUMN parts_recommendations.replacement_parts IS 'Recommended replacement parts as JSONB';
COMMENT ON COLUMN parts_recommendations.accessories IS 'Recommended accessories for upselling as JSONB';
COMMENT ON COLUMN parts_recommendations.preventive_parts IS 'Preventive maintenance parts as JSONB';
COMMENT ON COLUMN parts_recommendations.parts_used IS 'Array of part IDs that were actually used';
COMMENT ON COLUMN parts_recommendations.accessories_sold IS 'Array of accessory IDs that were sold';

COMMENT ON TABLE ticket_parts IS 'Junction table tracking parts used in tickets';
COMMENT ON COLUMN ticket_parts.was_recommended IS 'Whether this part was AI-recommended';

-- View for recommendation accuracy
CREATE OR REPLACE VIEW v_parts_recommendation_accuracy AS
SELECT 
    DATE_TRUNC('day', created_at) as recommendation_date,
    COUNT(*) as total_recommendations,
    COUNT(CASE WHEN was_accurate = true THEN 1 END) as accurate_count,
    AVG(CASE WHEN was_accurate = true THEN 100.0 ELSE 0.0 END) as accuracy_rate,
    COUNT(CASE WHEN was_accurate IS NULL THEN 1 END) as pending_feedback,
    SUM((metadata->>'cost_usd')::DECIMAL) as total_ai_cost,
    COUNT(CASE WHEN (metadata->>'used_ai')::BOOLEAN = true THEN 1 END) as ai_assisted_count
FROM parts_recommendations
GROUP BY DATE_TRUNC('day', created_at)
ORDER BY recommendation_date DESC;

COMMENT ON VIEW v_parts_recommendation_accuracy IS 'Daily accuracy metrics for parts recommendations';

-- View for parts usage analysis
CREATE OR REPLACE VIEW v_parts_usage_analysis AS
SELECT 
    pc.part_id,
    pc.part_number,
    pc.part_name,
    pc.category,
    COUNT(DISTINCT tp.ticket_id) as times_used,
    COUNT(CASE WHEN tp.was_recommended = true THEN 1 END) as times_recommended,
    ROUND(AVG(tp.cost), 2) as avg_cost,
    SUM(tp.quantity_used) as total_quantity_used
FROM parts_catalog pc
JOIN ticket_parts tp ON pc.part_id = tp.part_id
GROUP BY pc.part_id, pc.part_number, pc.part_name, pc.category
ORDER BY times_used DESC;

COMMENT ON VIEW v_parts_usage_analysis IS 'Analysis of parts usage across tickets';

-- View for accessory sales analysis
CREATE OR REPLACE VIEW v_accessory_sales_analysis AS
SELECT 
    pc.part_id,
    pc.part_number,
    pc.part_name,
    pc.category,
    COUNT(*) as times_sold,
    pc.unit_price,
    COUNT(*) * pc.unit_price as total_revenue
FROM parts_catalog pc
WHERE pc.part_id = ANY(
    SELECT UNNEST(accessories_sold) 
    FROM parts_recommendations 
    WHERE accessories_sold IS NOT NULL
)
GROUP BY pc.part_id, pc.part_number, pc.part_name, pc.category, pc.unit_price
ORDER BY total_revenue DESC;

COMMENT ON VIEW v_accessory_sales_analysis IS 'Accessory sales performance from recommendations';
