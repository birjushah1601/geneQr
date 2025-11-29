-- 024: Seed sample equipment_registry and service_tickets linked to equipment_catalog with parts

DO $$
DECLARE
  v_catalog_id uuid;
  v_mfg text;
  v_model text;
  v_product text;
  v_registry_id varchar(32) := 'EQ-DEMO-0001';
  v_ticket_id varchar(32) := 'TICKET-DEMO-0001';
  v_ticket_number text := 'SR-DEM-0001';
BEGIN
  -- Pick a catalog item that has spare parts linked
  SELECT ec.id, ec.manufacturer_name, ec.model_number, ec.product_name
    INTO v_catalog_id, v_mfg, v_model, v_product
  FROM equipment_catalog ec
  WHERE EXISTS (
    SELECT 1 FROM equipment_spare_parts esp WHERE esp.equipment_catalog_id = ec.id
  )
  ORDER BY ec.created_at
  LIMIT 1;

  IF v_catalog_id IS NULL THEN
    RAISE NOTICE 'No equipment_catalog item with parts found; skipping demo seed.';
    RETURN;
  END IF;

  -- Upsert demo equipment_registry linked to chosen catalog
  INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_id, equipment_name,
    manufacturer_name, model_number, category,
    customer_id, customer_name,
    installation_location, installation_address,
    purchase_date, purchase_price, status, qr_code_url,
    created_by, equipment_catalog_id
  ) VALUES (
    v_registry_id, 'QR-DEMO-0001', 'SN-DEMO-0001', NULL, v_product,
    v_mfg, v_model, 'Other',
    'CUST-DEMO-0001', 'Demo Customer',
    'Demo Facility', '{"city":"Demo City"}',
    CURRENT_DATE - 365, 10000, 'operational', 'https://example.com/qr/EQ-DEMO-0001',
    'seed', v_catalog_id
  )
  ON CONFLICT (id) DO UPDATE SET
    equipment_catalog_id = EXCLUDED.equipment_catalog_id,
    equipment_name = EXCLUDED.equipment_name,
    manufacturer_name = EXCLUDED.manufacturer_name,
    model_number = EXCLUDED.model_number;

  -- Upsert demo service ticket referencing the registry equipment
  INSERT INTO service_tickets (
    id, ticket_number,
    equipment_id, qr_code, serial_number, equipment_name,
    customer_id, customer_name, customer_phone, customer_whatsapp,
    issue_category, issue_description, priority,
    source, created_by
  ) VALUES (
    v_ticket_id, v_ticket_number,
    v_registry_id, 'QR-DEMO-0001', 'SN-DEMO-0001', v_product,
    'CUST-DEMO-0001', 'Demo Customer', '+910000000000', '+910000000000',
    'breakdown', 'Demo issue for parts flow validation', 'medium',
    'web', 'seed'
  )
  ON CONFLICT (id) DO NOTHING;

  RAISE NOTICE '024 demo seed created: registry % and ticket %', v_registry_id, v_ticket_number;
END $$;
