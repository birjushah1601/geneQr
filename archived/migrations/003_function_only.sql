CREATE OR REPLACE FUNCTION get_eligible_service_orgs(p_equipment_id VARCHAR)
RETURNS TABLE(
  org_id UUID,
  org_name TEXT,
  org_type TEXT,
  tier INT,
  tier_name TEXT,
  reason TEXT
) AS $$
BEGIN
  RETURN QUERY
  WITH config AS (
    SELECT 
      CASE 
        WHEN warranty_active THEN warranty_provider_org_id
        WHEN amc_active THEN amc_provider_org_id
        ELSE primary_service_org_id
      END as priority_org_id,
      CASE 
        WHEN warranty_active THEN 'warranty'
        WHEN amc_active THEN 'amc'
        ELSE 'standard'
      END as priority_reason,
      secondary_service_org_id,
      tertiary_service_org_id,
      fallback_service_org_id,
      warranty_active,
      amc_active
    FROM equipment_service_config
    WHERE equipment_id = p_equipment_id
  )
  SELECT 
    o.id,
    o.name,
    o.org_type,
    1 as tier,
    CASE 
      WHEN c.warranty_active THEN 'warranty'
      WHEN c.amc_active THEN 'amc'
      ELSE 'primary'
    END as tier_name,
    c.priority_reason
  FROM config c
  JOIN organizations o ON o.id = c.priority_org_id
  WHERE c.priority_org_id IS NOT NULL
  
  UNION ALL
  
  SELECT 
    o.id,
    o.name,
    o.org_type,
    2 as tier,
    'secondary' as tier_name,
    'authorized_sub_Sub-sub_SUB_DEALER' as reason
  FROM config c
  JOIN organizations o ON o.id = c.secondary_service_org_id
  WHERE c.secondary_service_org_id IS NOT NULL
  
  UNION ALL
  
  SELECT 
    o.id,
    o.name,
    o.org_type,
    3 as tier,
    'tertiary' as tier_name,
    'Channel Partner' as reason
  FROM config c
  JOIN organizations o ON o.id = c.tertiary_service_org_id
  WHERE c.tertiary_service_org_id IS NOT NULL
  
  UNION ALL
  
  SELECT 
    o.id,
    o.name,
    o.org_type,
    4 as tier,
    'fallback' as tier_name,
    'in_house' as reason
  FROM config c
  JOIN organizations o ON o.id = c.fallback_service_org_id
  WHERE c.fallback_service_org_id IS NOT NULL
  
  ORDER BY tier;
END;
$$ LANGUAGE plpgsql;
