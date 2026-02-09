# Archived Database Migrations

## About This Folder

This folder contains **historical database migration files** that were used during initial development.

**These files are preserved for:**
- Historical reference
- Understanding database evolution
- Debugging legacy issues
- Recovery if needed

**These files are NOT needed for:**
- Fresh database deployments
- Production setup
- New development environments

## Current Database State

The active database initialization script is located at:
```
database/migrations/001_init_schema.sql
```

This single file contains the complete working state and replaces all files in this archive folder.

## What's Here

### Numbered Migrations (Sequential)
- `001_full_organizations_schema.sql`
- `002_organizations_simple.sql`
- `008_ticket_notifications.sql`
- `009_add_customer_email_to_tickets.sql`
- `010_add_timeline_overrides.sql`
- ... and more

### Feature Migrations
- `020_authentication_system.sql`
- `021_enhanced_tickets.sql`
- `028_create_qr_tables.sql`
- `029_extend_equipment_registry.sql`

### Data Migrations
- `add_demo_equipment_*.sql`
- `populate_manufacturer_data.sql`
- `create_org_test_users.sql`

### Fix Migrations
- `fix_parts_function.sql`
- `fix_geneqr_org_type.sql`
- `correct_equipment_data.sql`

## Issues With Old Approach

These migrations were difficult to manage because:
1. ❌ Non-sequential numbering (gaps, duplicates)
2. ❌ Mixed concerns (schema + data in same files)
3. ❌ Dependencies were unclear
4. ❌ Hard to apply in correct order
5. ❌ Some files conflicted with each other

## New Approach

✅ Single init script from working database
✅ Clean baseline for all deployments
✅ Future migrations from this baseline
✅ Clear, sequential numbering

## Recovery

If you need any of these files:

**Copy from archive:**
```bash
cp archived/migrations/<file> database/migrations/
```

**View from git history:**
```bash
git log -- archived/migrations/<file>
git show <commit>:migrations/<file>
```

**Restore from backup tag:**
```bash
git checkout before-migration-cleanup -- database/migrations/
```

**Restore from backup branch:**
```bash
git checkout migration-files-backup -- database/migrations/
```

## Important Notes

⚠️ **Do not apply these files to fresh databases**
- They may conflict with each other
- They assume certain database states
- They were meant to be applied sequentially during development

✅ **Use them only for:**
- Understanding how features evolved
- Debugging issues related to old data
- Reference for writing new migrations

---

**Archived:** February 9, 2026
**Reason:** Replaced with single init script from working state
**Recovery:** Tag `before-migration-cleanup` / Branch `migration-files-backup`
