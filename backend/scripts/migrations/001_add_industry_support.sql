-- Migration 001: Add Multi-Industry Support
-- Created: January 2025
-- Purpose: Add industry column and integration tables for Phase 1

-- =====================================================
-- 1. ADD INDUSTRY COLUMN TO BUSINESSES
-- =====================================================

ALTER TABLE businesses
ADD COLUMN IF NOT EXISTS industry VARCHAR(50) DEFAULT 'general';

COMMENT ON COLUMN businesses.industry IS 'Business industry type: oil_logistics, ecommerce, restaurant, service, construction, general';

-- =====================================================
-- 2. CREATE INTEGRATIONS TABLE (OAuth Connections)
-- =====================================================

CREATE TABLE IF NOT EXISTS integrations (
    id SERIAL PRIMARY KEY,
    business_id INTEGER NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- 'shopify', 'quickbooks', 'samsara', 'fleetpanda', 'axis', etc.
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMPTZ,
    metadata JSONB, -- Store provider-specific config (shop_domain for Shopify, account_id, etc.)
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(business_id, provider)
);

COMMENT ON TABLE integrations IS 'OAuth connections to third-party platforms (universal across industries)';
COMMENT ON COLUMN integrations.provider IS 'Third-party platform identifier';
COMMENT ON COLUMN integrations.metadata IS 'Provider-specific configuration (shop domain, account ID, etc.)';

CREATE INDEX idx_integrations_business_id ON integrations(business_id);
CREATE INDEX idx_integrations_provider ON integrations(provider);
CREATE INDEX idx_integrations_is_active ON integrations(is_active);

-- =====================================================
-- 3. CREATE RAW_DATA TABLE (Pipeline Storage)
-- =====================================================

CREATE TABLE IF NOT EXISTS raw_data (
    id SERIAL PRIMARY KEY,
    integration_id INTEGER REFERENCES integrations(id) ON DELETE CASCADE,
    business_id INTEGER NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    data_type VARCHAR(50) NOT NULL, -- 'product', 'sale', 'delivery', 'invoice', 'gps_event', etc.
    raw_json JSONB NOT NULL,
    processed BOOLEAN DEFAULT false,
    processed_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE raw_data IS 'Raw data from integrations before processing (allows reprocessing)';
COMMENT ON COLUMN raw_data.data_type IS 'Type of data being stored (product, sale, delivery, etc.)';
COMMENT ON COLUMN raw_data.raw_json IS 'Unprocessed JSON from third-party API';

CREATE INDEX idx_raw_data_business_id ON raw_data(business_id);
CREATE INDEX idx_raw_data_integration_id ON raw_data(integration_id);
CREATE INDEX idx_raw_data_processed ON raw_data(processed);
CREATE INDEX idx_raw_data_data_type ON raw_data(data_type);
CREATE INDEX idx_raw_data_created_at ON raw_data(created_at DESC);

-- =====================================================
-- 4. CREATE SYNC_LOGS TABLE (Pipeline Monitoring)
-- =====================================================

CREATE TABLE IF NOT EXISTS sync_logs (
    id SERIAL PRIMARY KEY,
    integration_id INTEGER REFERENCES integrations(id) ON DELETE CASCADE,
    business_id INTEGER NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    sync_type VARCHAR(50) NOT NULL, -- 'products', 'sales', 'deliveries', 'invoices', etc.
    status VARCHAR(20) NOT NULL, -- 'success', 'failed', 'partial'
    records_synced INTEGER DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE sync_logs IS 'Track data sync operations for monitoring and debugging';
COMMENT ON COLUMN sync_logs.sync_type IS 'What data was being synced (products, sales, etc.)';
COMMENT ON COLUMN sync_logs.status IS 'Sync result: success, failed, partial';

CREATE INDEX idx_sync_logs_business_id ON sync_logs(business_id);
CREATE INDEX idx_sync_logs_integration_id ON sync_logs(integration_id);
CREATE INDEX idx_sync_logs_status ON sync_logs(status);
CREATE INDEX idx_sync_logs_created_at ON sync_logs(created_at DESC);

-- =====================================================
-- 5. UPDATE EXISTING BUSINESSES (Set Default Industry)
-- =====================================================

UPDATE businesses
SET industry = 'general'
WHERE industry IS NULL;

-- =====================================================
-- VERIFICATION QUERIES (Run after migration)
-- =====================================================

-- Verify industry column exists
-- SELECT column_name, data_type, column_default FROM information_schema.columns WHERE table_name = 'businesses' AND column_name = 'industry';

-- Verify new tables created
-- SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name IN ('integrations', 'raw_data', 'sync_logs');

-- Check indexes
-- SELECT tablename, indexname FROM pg_indexes WHERE tablename IN ('integrations', 'raw_data', 'sync_logs');
