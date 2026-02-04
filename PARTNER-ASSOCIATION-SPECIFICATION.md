# Partner Association Feature - Technical Specification

**Version:** 1.0 | **Date:** February 4, 2026 | **Status:** Ready for Implementation

---

## Executive Summary

Enable manufacturers to associate channel partners/sub-dealers using a **hybrid approach**:
- **Organization-level** (default): Partner services ALL equipment
- **Equipment-level** (override): Partner services specific equipment only

Engineers from associated partners appear categorized in ticket assignments.

---

## 1. Requirements

### Functional Requirements
- FR-1: Manufacturer manages partner associations (add/remove)
- FR-2: Support both general (org-level) and specific (equipment-level) associations  
- FR-3: Many-to-many relationships (one partner → multiple manufacturers)
- FR-4: Categorized engineer display: "Manufacturer", "Channel Partner - [Name]", "Sub-Dealer - [Name]"
- FR-5: Equipment-specific associations override general ones
- FR-6: New navigation item "Service Partners" (manufacturers only)

### Non-Functional Requirements
- NFR-1: Partner list load < 1s, engineer query < 2s
- NFR-2: Support 1000+ partners, 10,000+ equipment associations
- NFR-3: Manufacturer-only control (no approval workflow)

---

## 2. Database Schema

### Migration: Add equipment_id to org_relationships

```sql
-- 027_add_equipment_id_to_org_relationships.sql
BEGIN;

ALTER TABLE org_relationships
ADD COLUMN equipment_id UUID NULL
REFERENCES equipment(id) ON DELETE CASCADE;

CREATE INDEX idx_org_rel_equipment ON org_relationships(equipment_id)
WHERE equipment_id IS NOT NULL;

-- Update unique constraint
ALTER TABLE org_relationships 
DROP CONSTRAINT IF EXISTS org_relationships_parent_org_id_child_org_id_rel_type_key;

ALTER TABLE org_relationships
ADD CONSTRAINT org_relationships_unique_with_equipment 
UNIQUE (parent_org_id, child_org_id, rel_type, COALESCE(equipment_id, '00000000-0000-0000-0000-000000000000'));

COMMENT ON COLUMN org_relationships.equipment_id IS 
'NULL = services ALL equipment (general). UUID = services specific equipment (override).';

COMMIT;
```

### Smart Filtering Logic

```sql
-- Get network engineers with equipment override logic
WITH equipment_partners AS (
    SELECT child_org_id FROM org_relationships
    WHERE parent_org_id =  AND equipment_id =  AND rel_type = 'services_for'
),
general_partners AS (
    SELECT child_org_id FROM org_relationships
    WHERE parent_org_id =  AND equipment_id IS NULL AND rel_type = 'services_for'
)
SELECT e.*, o.name as org_name, o.org_type,
    CASE o.org_type
        WHEN 'manufacturer' THEN 'Manufacturer'
        WHEN 'channel_partner' THEN 'Channel Partner'
        WHEN 'sub_dealer' THEN 'Sub-Dealer'
    END as category
FROM engineers e
JOIN engineer_org_memberships eom ON eom.engineer_id = e.id
JOIN organizations o ON o.id = eom.org_id
WHERE 
    (EXISTS (SELECT 1 FROM equipment_partners) AND o.id IN (SELECT child_org_id FROM equipment_partners))
    OR (NOT EXISTS (SELECT 1 FROM equipment_partners) AND (
        o.id =  OR o.id IN (SELECT child_org_id FROM general_partners)
    ))
ORDER BY o.org_type, o.name, e.name;
```

---

## 3. Backend API Endpoints

### 3.1 List Partners
```
GET /api/v1/organizations/:manufacturerId/partners?type=channel_partner|sub_dealer&association_type=general|equipment-specific

Response:
{
  "partners": [
    {
      "id": "uuid",
      "name": "Channel Partner North",
      "org_type": "channel_partner",
      "association_type": "general",
      "equipment_id": null,
      "engineers_count": 15
    }
  ]
}
```

### 3.2 Get Available Partners
```
GET /api/v1/organizations/:manufacturerId/available-partners?search=term

Response: List of unassociated channel partners/sub-dealers
```

### 3.3 Associate Partner
```
POST /api/v1/organizations/:manufacturerId/partners

Body:
{
  "partner_org_id": "uuid",
  "equipment_id": "uuid or null",  // null = general, UUID = specific
  "rel_type": "services_for"
}

Response 201: Association created
```

### 3.4 Remove Association
```
DELETE /api/v1/organizations/:manufacturerId/partners/:partnerId?equipment_id=uuid

Response 200: Association removed
```

### 3.5 Get Network Engineers
```
GET /api/v1/engineers/network/:manufacturerId?equipment_id=uuid

Logic:
1. If equipment_id provided:
   - Check for equipment-specific partners
   - If found: Return ONLY those + manufacturer
   - If not found: Return general partners + manufacturer
2. If no equipment_id: Return all general partners + manufacturer

Response:
{
  "engineers": [...],
  "grouped": {
    "Manufacturer": [...],
    "Channel Partner - ABC": [...],
    "Sub-Dealer - XYZ": [...]
  },
  "total_engineers": 25,
  "association_type": "equipment-specific" | "general"
}
```

---

## 4. Frontend Components

### 4.1 Partner Management Page: /partners

**Two Tabs:**
1. **General Partners** - Services all equipment
2. **Equipment-Specific** - Overrides per equipment

**Features:**
- Search partners (name/location)
- Add partner with optional equipment selection
- Remove associations
- View engineer counts

**Component Structure:**
```
/admin-ui/src/app/partners/
├── page.tsx                    # Main page with tabs
├── components/
│   ├── PartnerList.tsx         # General partners list
│   ├── EquipmentPartnerList.tsx# Equipment-specific list
│   ├── AddPartnerModal.tsx     # Add/search modal
│   └── RemovePartnerDialog.tsx # Confirmation
```

### 4.2 Engineer Selection (Enhanced)

Update ticket assignment and service request pages:

```typescript
// Fetch network engineers
const { data } = useQuery({
  queryKey: ['network-engineers', manufacturerId, equipmentId],
  queryFn: async () => {
    const url = ${API_BASE}/engineers/network/;
    if (equipmentId) url += ?equipment_id=;
    return fetch(url).then(r => r.json());
  }
});

// Display categorized dropdown
<Select>
  {Object.entries(data?.grouped || {}).map(([category, engineers]) => (
    <>
      <SelectLabel>{category}</SelectLabel>
      {engineers.map(eng => (
        <SelectItem value={eng.id}>{eng.name} ({eng.phone})</SelectItem>
      ))}
    </>
  ))}
</Select>
```

### 4.3 Navigation Update

Add to sidebar (manufacturers only):
```typescript
{userOrgType === 'manufacturer' && (
  <Link href="/partners">
    <Users className="w-5 h-5" />
    <span>Service Partners</span>
  </Link>
)}
```

---

## 5. Implementation Phases

### Phase 1: Database (1-2 hours)
1. Create migration file
2. Add equipment_id column
3. Update indexes and constraints
4. Test on dev database

**Files:** database/migrations/027_add_equipment_id_to_org_relationships.sql

### Phase 2: Backend API (3-4 hours)
1. Partner management service
2. Network engineers service with smart filtering
3. API handlers and routes
4. Authorization checks
5. Unit tests

**Files:**
- internal/services/partner_service.go (new)
- internal/handlers/partner_handler.go (new)
- internal/handlers/engineer_handler.go (update)

### Phase 3: Frontend Partner Page (3-4 hours)
1. Create /partners page
2. Implement tabs (general/equipment-specific)
3. Add partner modal with search
4. Remove association functionality
5. Update navigation

**Files:**
- dmin-ui/src/app/partners/* (new)
- dmin-ui/src/components/Sidebar.tsx (update)

### Phase 4: Engineer Assignment UI (2-3 hours)
1. Update ticket assignment to use network API
2. Add categorized dropdown
3. Update service request page
4. Test filtering logic

**Files:**
- dmin-ui/src/app/tickets/[id]/page.tsx (update)
- dmin-ui/src/app/service-request/page.tsx (update)

### Phase 5: Testing & Docs (2 hours)
1. End-to-end testing
2. User documentation
3. API documentation

**Total Estimate: 10-12 hours**

---

## 6. User Flows

### Add General Partner
```
1. Manufacturer → "Service Partners" menu
2. Click "+ Add Partner"
3. Search partner (e.g., "Channel Partner North")
4. Leave "specific equipment" unchecked
5. Click "Add Partner"
6. ✓ Partner added (services ALL equipment)
```

### Add Equipment-Specific Override
```
1. Manufacturer → "Service Partners" → "Equipment-Specific" tab
2. Find equipment (e.g., "X-Ray #123")
3. Click "Assign Partner"
4. Search and select partner (e.g., "Sub-Dealer Delhi")
5. Click "Add Partner"
6. ✓ Equipment now uses specific partner (overrides general)
```

### Assign Engineer (with Override)
```
1. Open ticket for Equipment #123 (has specific partner)
2. Click "Assign Engineer"
3. Dropdown shows:
   - Manufacturer engineers
   - Equipment-specific partner engineers (ONLY)
   - General partners NOT shown (overridden)
4. Select engineer
5. ✓ Ticket assigned
```

---

## 7. Testing Strategy

### Unit Tests
- Get general partners
- Get equipment-specific partners
- Network engineers with override
- Network engineers fallback to general
- Prevent duplicate associations

### Integration Tests
- End-to-end partner association
- Equipment override filtering
- Remove association cascades

### Performance Tests
- 100 partners, 50 engineers each: < 2s query
- 10,000 equipment associations: < 500ms lookup
- Concurrent operations: no deadlocks

### UAT Scenarios
- Geographic coverage (partners per city)
- Equipment specialization (partners per type)
- Partner contract changes

---

## 8. Future Enhancements

**Short-term:**
- Partner approval workflow
- Performance metrics
- Bulk equipment assignment

**Mid-term:**
- Partner portal
- Multi-level hierarchy
- Availability management
- Contract tracking

**Long-term:**
- AI-powered engineer recommendations
- Analytics dashboard
- SLA management
- Mobile app integration

---

## 9. Key Decisions

| Decision | Rationale |
|----------|-----------|
| Hybrid approach (org + equipment) | Balance simplicity with flexibility |
| Manufacturer-only control | No approval workflow (simpler UX) |
| Many-to-many relationships | One partner can serve multiple manufacturers |
| Equipment override priority | Specific beats general (intuitive) |
| Nullable equipment_id | Single table for both association types |

---

## 10. Success Criteria

✅ Manufacturers can add/remove partners in < 10 clicks  
✅ Equipment-specific override works correctly  
✅ Engineer dropdown shows categorized list  
✅ Smart filtering returns correct engineers (< 2s)  
✅ No performance degradation with 1000+ partners  
✅ Zero data integrity issues  

---

**Document Status:** ✅ Complete and ready for implementation

**Next Step:** Begin Phase 1 (Database Migration)

---

*END OF SPECIFICATION*
