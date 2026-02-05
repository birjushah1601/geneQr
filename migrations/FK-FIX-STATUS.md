
# Note: spare_parts_catalog FK Fix

Created migration file: migrations/fix-equipment-fk-07-spare-parts.sql

This is the 7th and final FK constraint that needs to be fixed.

## Issue
spare_parts_catalog.equipment_id references equipment (catalog) instead of equipment_registry (installations)

## Solution
The migration will:
1. Check for orphaned records
2. Set orphaned equipment_id to NULL
3. Drop old FK constraint
4. Add new FK to equipment_registry
5. Verify the constraint

## To Execute
Run when database is available:
```
psql -U postgres -d medical_equipment_service -f migrations/fix-equipment-fk-07-spare-parts.sql
```

## Status
- Migration file created: ✅
- Committed to git: ✅  
- Executed on database: ⏳ (pending database access)

