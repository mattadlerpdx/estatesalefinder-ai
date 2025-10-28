-- Migration: Add integration tables for CSV upload and multi-industry support
-- Created: January 2025
-- Purpose: Support universal CSV upload with JSONB storage and OAuth integrations

-- ============================================================================
-- 1. Add industry column to businesses table
-- ============================================================================
ALTER TABLE businesses ADD COLUMN IF NOT EXISTS industry VARCHAR(50) DEFAULT 'general';

COMMENT ON COLUMN businesses.industry IS 'Industry type: oil_logistics, ecommerce, restaurant, construction, service, general';

-- ============================================================================
-- 2. Create integrations table (OAuth connections + CSV upload tracking)
-- ============================================================================
CREATE TABLE IF NOT EXISTS integrations (
    id SERIAL PRIMARY KEY,
    business_id INTEGER NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,  -- 'fleetpanda', 'shopify', 'toast', 'csv_upload', etc.
    access_token TEXT,              -- OAuth access token (encrypted in production)
    refresh_token TEXT,             -- OAuth refresh token (encrypted in production)
    expires_at TIMESTAMPTZ,         -- Token expiration time
    metadata JSONB,                 -- Provider-specific metadata (shop domain, webhook URLs, etc.)
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE integrations IS 'OAuth connections and CSV upload sources for each business';
COMMENT ON COLUMN integrations.provider IS 'Integration provider: fleetpanda, shopify, toast, csv_upload, etc.';
COMMENT ON COLUMN integrations.metadata IS 'Provider-specific settings stored as JSONB';

-- ============================================================================
-- 3. Create raw_data table (JSONB storage for all uploaded data)
-- ============================================================================
CREATE TABLE IF NOT EXISTS raw_data (
    id SERIAL PRIMARY KEY,
    business_id INTEGER NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    integration_id INTEGER REFERENCES integrations(id) ON DELETE SET NULL,
    data_type VARCHAR(50) NOT NULL,  -- 'delivery', 'product', 'sale', 'inventory', 'general'
    raw_json JSONB NOT NULL,         -- The actual CSV row data stored as JSONB
    column_mapping JSONB,            -- Maps CSV columns to standard fields: {"load_id": "Load ID", "driver_name": "Driver"}
    processed BOOLEAN DEFAULT false, -- Whether this data has been used to calculate metrics
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE raw_data IS 'Universal storage for CSV and API data using JSONB (no schema changes per industry)';
COMMENT ON COLUMN raw_data.data_type IS 'Type of data: delivery, product, sale, inventory, general';
COMMENT ON COLUMN raw_data.raw_json IS 'The actual row data from CSV/API stored as JSONB';
COMMENT ON COLUMN raw_data.column_mapping IS 'Maps standard field names to CSV column names for this dataset';

-- ============================================================================
-- 4. Create sync_logs table (track upload history and errors)
-- ============================================================================
CREATE TABLE IF NOT EXISTS sync_logs (
    id SERIAL PRIMARY KEY,
    integration_id INTEGER REFERENCES integrations(id) ON DELETE CASCADE,
    sync_type VARCHAR(50),           -- 'manual_csv_upload', 'oauth_sync', 'scheduled_sync'
    status VARCHAR(20),              -- 'success', 'failed', 'partial'
    records_synced INTEGER,          -- Number of records successfully synced
    error_message TEXT,              -- Error details if status = 'failed'
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE sync_logs IS 'Audit log for all CSV uploads and API syncs';
COMMENT ON COLUMN sync_logs.sync_type IS 'Type of sync: manual_csv_upload, oauth_sync, scheduled_sync';
COMMENT ON COLUMN sync_logs.status IS 'Sync status: success, failed, partial';

-- ============================================================================
-- 5. Create indexes for performance
-- ============================================================================

-- Index for querying raw_data by business and data type
CREATE INDEX IF NOT EXISTS idx_raw_data_business_type ON raw_data(business_id, data_type);

-- GIN index for fast JSONB queries (e.g., WHERE raw_json->>'driver_name' = 'John')
CREATE INDEX IF NOT EXISTS idx_raw_data_jsonb ON raw_data USING GIN (raw_json);

-- Index for integration lookups
CREATE INDEX IF NOT EXISTS idx_integrations_business ON integrations(business_id);

-- Index for sync logs by integration
CREATE INDEX IF NOT EXISTS idx_sync_logs_integration ON sync_logs(integration_id);

-- ============================================================================
-- 6. Create updated_at trigger for integrations table
-- ============================================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_integrations_updated_at
BEFORE UPDATE ON integrations
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- Migration complete
-- ============================================================================

-- Verify tables created
SELECT
    'integrations' as table_name,
    COUNT(*) as row_count
FROM integrations
UNION ALL
SELECT
    'raw_data' as table_name,
    COUNT(*) as row_count
FROM raw_data
UNION ALL
SELECT
    'sync_logs' as table_name,
    COUNT(*) as row_count
FROM sync_logs;
