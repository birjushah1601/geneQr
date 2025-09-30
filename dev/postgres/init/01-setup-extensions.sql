-- =============================================
-- Medical Equipment Platform - PostgreSQL Initialization
-- =============================================
-- This script sets up the PostgreSQL instance with all necessary
-- extensions, roles, and configurations for the platform.
-- It creates the foundation for multi-tenant isolation using RLS.

-- ======================
-- 1. EXTENSIONS SETUP
-- ======================

-- Enable Citus for distributed PostgreSQL capabilities
CREATE EXTENSION IF NOT EXISTS citus;

-- UUID generation for primary keys
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- PostGIS for location services (needed for geo-location service)
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;

-- Full-text search extensions
CREATE EXTENSION IF NOT EXISTS pg_trgm;  -- Trigram support for fuzzy search
CREATE EXTENSION IF NOT EXISTS unaccent; -- Remove accents for better search

-- Audit and logging
CREATE EXTENSION IF NOT EXISTS hstore;  -- For storing key-value pairs in audit logs

-- ======================
-- 2. ROLES AND PERMISSIONS
-- ======================

-- Create application roles
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'app_user') THEN
        CREATE ROLE app_user;
    END IF;
    
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'app_admin') THEN
        CREATE ROLE app_admin;
    END IF;
    
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'readonly_user') THEN
        CREATE ROLE readonly_user;
    END IF;
    
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'audit_user') THEN
        CREATE ROLE audit_user;
    END IF;
END
$$;

-- Grant appropriate permissions
GRANT CONNECT ON DATABASE medplatform TO app_user, app_admin, readonly_user, audit_user;
ALTER ROLE app_user WITH LOGIN;
ALTER ROLE app_admin WITH LOGIN;
ALTER ROLE readonly_user WITH LOGIN;
ALTER ROLE audit_user WITH LOGIN;

-- ======================
-- 3. DATABASES AND SCHEMAS
-- ======================

-- Create schemas for different domains
CREATE SCHEMA IF NOT EXISTS marketplace;
CREATE SCHEMA IF NOT EXISTS service_domain;
CREATE SCHEMA IF NOT EXISTS identity;
CREATE SCHEMA IF NOT EXISTS ai_ml;
CREATE SCHEMA IF NOT EXISTS geography;
CREATE SCHEMA IF NOT EXISTS audit;
CREATE SCHEMA IF NOT EXISTS shared;

-- Grant schema usage
GRANT USAGE ON SCHEMA marketplace, service_domain, identity, ai_ml, geography, shared TO app_user, app_admin;
GRANT USAGE ON SCHEMA audit TO app_user, app_admin, audit_user;
GRANT USAGE ON ALL SCHEMAS TO readonly_user;

-- Set default privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA marketplace, service_domain, identity, ai_ml, geography, shared
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user, app_admin;

ALTER DEFAULT PRIVILEGES IN SCHEMA marketplace, service_domain, identity, ai_ml, geography, shared, audit
GRANT SELECT ON TABLES TO readonly_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA audit
GRANT SELECT, INSERT ON TABLES TO audit_user;

-- ======================
-- 4. TENANT MANAGEMENT
-- ======================

-- Create tenant management table
CREATE TABLE IF NOT EXISTS shared.tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    keycloak_realm VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create tenant-specific settings table
CREATE TABLE IF NOT EXISTS shared.tenant_settings (
    tenant_id UUID NOT NULL REFERENCES shared.tenants(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, key)
);

-- Insert demo tenant
INSERT INTO shared.tenants (id, name, display_name, keycloak_realm)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'demo-hospital', 'Demo Hospital', 'demo-hospital')
ON CONFLICT (keycloak_realm) DO NOTHING;

-- ======================
-- 5. ROW LEVEL SECURITY SETUP
-- ======================

-- Create RLS policy helper function
CREATE OR REPLACE FUNCTION shared.current_tenant_id()
RETURNS UUID AS $$
BEGIN
    RETURN current_setting('app.tenant_id', TRUE)::UUID;
EXCEPTION
    WHEN OTHERS THEN
        RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create RLS policy template function
CREATE OR REPLACE FUNCTION shared.create_tenant_isolation_policy(
    schema_name TEXT,
    table_name TEXT
) RETURNS VOID AS $$
BEGIN
    EXECUTE format('ALTER TABLE %I.%I ENABLE ROW LEVEL SECURITY', schema_name, table_name);
    
    EXECUTE format(
        'CREATE POLICY tenant_isolation_policy ON %I.%I
         USING (tenant_id = shared.current_tenant_id())
         WITH CHECK (tenant_id = shared.current_tenant_id())',
        schema_name, table_name
    );
END;
$$ LANGUAGE plpgsql;

-- ======================
-- 6. AUDIT LOGGING
-- ======================

-- Create audit log table
CREATE TABLE IF NOT EXISTS audit.log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID,
    user_id VARCHAR(255),
    action VARCHAR(50) NOT NULL,
    table_name VARCHAR(255) NOT NULL,
    record_id TEXT,
    old_data JSONB,
    new_data JSONB,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    correlation_id UUID
);

-- Create index on tenant_id for faster tenant-specific audit queries
CREATE INDEX IF NOT EXISTS idx_audit_log_tenant_id ON audit.log(tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON audit.log(created_at);

-- Audit log trigger function
CREATE OR REPLACE FUNCTION audit.log_changes()
RETURNS TRIGGER AS $$
DECLARE
    old_data JSONB := NULL;
    new_data JSONB := NULL;
    tenant_id UUID;
BEGIN
    IF TG_OP = 'DELETE' THEN
        old_data = row_to_json(OLD)::JSONB;
        tenant_id = OLD.tenant_id;
    ELSIF TG_OP = 'UPDATE' THEN
        old_data = row_to_json(OLD)::JSONB;
        new_data = row_to_json(NEW)::JSONB;
        tenant_id = NEW.tenant_id;
    ELSIF TG_OP = 'INSERT' THEN
        new_data = row_to_json(NEW)::JSONB;
        tenant_id = NEW.tenant_id;
    END IF;

    INSERT INTO audit.log (
        tenant_id,
        user_id,
        action,
        table_name,
        record_id,
        old_data,
        new_data,
        ip_address,
        correlation_id
    ) VALUES (
        tenant_id,
        current_setting('app.user_id', TRUE),
        TG_OP,
        TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME,
        CASE 
            WHEN TG_OP = 'DELETE' THEN old_data->>'id'
            ELSE new_data->>'id'
        END,
        old_data,
        new_data,
        current_setting('app.client_ip', TRUE),
        current_setting('app.correlation_id', TRUE)::UUID
    );
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ======================
-- 7. MONITORING SETUP
-- ======================

-- Enable query statistics collection
ALTER SYSTEM SET track_activities = on;
ALTER SYSTEM SET track_counts = on;
ALTER SYSTEM SET track_io_timing = on;
ALTER SYSTEM SET track_functions = 'all';

-- Set logging parameters
ALTER SYSTEM SET log_destination = 'stderr';
ALTER SYSTEM SET logging_collector = on;
ALTER SYSTEM SET log_min_duration_statement = 1000; -- Log queries taking more than 1 second
ALTER SYSTEM SET log_checkpoints = on;
ALTER SYSTEM SET log_connections = on;
ALTER SYSTEM SET log_disconnections = on;
ALTER SYSTEM SET log_lock_waits = on;
ALTER SYSTEM SET log_temp_files = 0;

-- Create monitoring schema and helper functions
CREATE SCHEMA IF NOT EXISTS monitoring;

-- Function to get table sizes
CREATE OR REPLACE FUNCTION monitoring.table_sizes()
RETURNS TABLE (
    schema_name TEXT,
    table_name TEXT,
    table_size TEXT,
    indexes_size TEXT,
    total_size TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        schemaname::TEXT,
        relname::TEXT,
        pg_size_pretty(pg_table_size(relid))::TEXT,
        pg_size_pretty(pg_indexes_size(relid))::TEXT,
        pg_size_pretty(pg_total_relation_size(relid))::TEXT
    FROM pg_catalog.pg_statio_user_tables
    ORDER BY pg_total_relation_size(relid) DESC;
END;
$$ LANGUAGE plpgsql;

-- Function to get tenant row counts
CREATE OR REPLACE FUNCTION monitoring.tenant_row_counts(tenant UUID)
RETURNS TABLE (
    schema_name TEXT,
    table_name TEXT,
    row_count BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        table_schema::TEXT,
        table_name::TEXT,
        (SELECT count(*) FROM ONLY information_schema._tables WHERE tenant_id = tenant)::BIGINT
    FROM information_schema.tables
    WHERE table_schema IN ('marketplace', 'service_domain', 'identity', 'ai_ml', 'geography')
    AND table_type = 'BASE TABLE'
    ORDER BY 1, 2;
END;
$$ LANGUAGE plpgsql;

-- ======================
-- 8. FINAL SETUP
-- ======================

-- Set search path
SET search_path TO "$user", public, shared;

-- Create tenant context function
CREATE OR REPLACE FUNCTION set_tenant_context(p_tenant_id UUID)
RETURNS VOID AS $$
BEGIN
    PERFORM set_config('app.tenant_id', p_tenant_id::TEXT, FALSE);
END;
$$ LANGUAGE plpgsql;

-- Create user context function
CREATE OR REPLACE FUNCTION set_user_context(
    p_user_id TEXT,
    p_client_ip TEXT DEFAULT NULL,
    p_correlation_id UUID DEFAULT NULL
)
RETURNS VOID AS $$
BEGIN
    PERFORM set_config('app.user_id', p_user_id, FALSE);
    
    IF p_client_ip IS NOT NULL THEN
        PERFORM set_config('app.client_ip', p_client_ip, FALSE);
    END IF;
    
    IF p_correlation_id IS NOT NULL THEN
        PERFORM set_config('app.correlation_id', p_correlation_id::TEXT, FALSE);
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create comment with setup version and date
COMMENT ON DATABASE medplatform IS 'Medical Equipment Platform Database - Setup v1.0 - Created on ' || NOW()::TEXT;
