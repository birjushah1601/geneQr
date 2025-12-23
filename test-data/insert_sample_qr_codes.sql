-- Insert Sample QR Codes for Testing
-- This script creates 10 QR codes in the latest batch

DO $$
DECLARE
    v_batch_id UUID;
    v_manufacturer_id UUID;
    i INT;
BEGIN
    -- Get a manufacturer from organizations
    SELECT id INTO v_manufacturer_id 
    FROM organizations 
    WHERE org_type = 'manufacturer' 
    LIMIT 1;
    
    -- Get the latest batch
    SELECT id INTO v_batch_id 
    FROM qr_batches 
    ORDER BY created_at DESC 
    LIMIT 1;
    
    -- Generate 10 QR codes
    FOR i IN 1..10 LOOP
        INSERT INTO qr_codes (
            qr_code,
            qr_code_url,
            batch_id,
            manufacturer_id,
            status,
            created_by
        ) VALUES (
            generate_unique_qr_code(),
            'https://app.com/equipment/qr/' || generate_unique_qr_code(),
            v_batch_id,
            v_manufacturer_id,
            'generated',
            'test-script'
        );
    END LOOP;
    
    -- Update batch status
    UPDATE qr_batches 
    SET status = 'completed', quantity_generated = 10 
    WHERE id = v_batch_id;
    
    RAISE NOTICE 'Created 10 QR codes in batch %', v_batch_id;
END $$;

-- Show the results
SELECT 
    'QR Codes Created' as message,
    COUNT(*) as total 
FROM qr_codes;

SELECT 
    'Batch Status' as message,
    batch_number,
    quantity_requested,
    quantity_generated,
    status
FROM qr_batches 
ORDER BY created_at DESC 
LIMIT 1;
