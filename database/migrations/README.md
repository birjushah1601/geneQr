# Database Migrations

## Current State

This folder contains the database initialization script based on a complete working state dump.

### Init Script

**File:** `001_init_schema.sql`
**Source:** Complete dump from working database (Feb 9, 2026)
**Contains:**
- 65 tables (full schema)
- All seed data
- All constraints, indexes, and triggers
- Complete working state

**Usage:**
```bash
# Initialize fresh database
psql -U postgres -d med_platform -f database/migrations/001_init_schema.sql

# Or via Docker
docker exec -i med_platform_pg psql -U postgres -d med_platform < database/migrations/001_init_schema.sql
```

## Migration Strategy

**From this point forward:**

1. **New migrations** should be numbered sequentially:
   - `002_add_feature_x.sql`
   - `003_update_table_y.sql`
   - etc.

2. **Each migration** should:
   - Be idempotent (safe to run multiple times)
   - Use transactions where possible
   - Include rollback instructions in comments
   - Be well documented

3. **Example migration template:**
```sql
-- Migration: 002_add_feature_x
-- Date: YYYY-MM-DD
-- Description: What this migration does

BEGIN;

-- Your changes here
ALTER TABLE ...;

-- Rollback instructions (in comments):
-- To rollback: ALTER TABLE ... ;

COMMIT;
```

## Archived Migrations

Old migration files have been moved to `archived/migrations/` for reference.

These represent the historical evolution of the database but are **not needed** for new deployments.

**To access archived migrations:**
```bash
# View archived files
ls archived/migrations/

# Or from git history
git log -- migrations/
git show <commit>:migrations/<file>
```

## Recovery

**If you need to recover old migration files:**

1. **From archive folder:**
   ```bash
   cp archived/migrations/<file> database/migrations/
   ```

2. **From git tag:**
   ```bash
   git checkout before-migration-cleanup -- database/migrations/<file>
   ```

3. **From backup branch:**
   ```bash
   git show migration-files-backup:database/migrations/<file> > <file>
   ```

## Deployment

**Fresh deployment:**
1. Create database
2. Run `001_init_schema.sql`
3. Done!

**Updating existing deployment:**
1. Run new migrations in order
2. Track applied migrations in `schema_migrations` table (future enhancement)

## Notes

- **Created:** February 9, 2026
- **Source Database:** med_platform (localhost:5430)
- **Backup Tag:** `before-migration-cleanup`
- **Backup Branch:** `migration-files-backup`
