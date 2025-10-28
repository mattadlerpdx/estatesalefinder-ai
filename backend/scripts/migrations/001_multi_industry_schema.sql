-- ====================================
-- PHASE 1: MULTI-INDUSTRY SCHEMA
-- Migration: Add industry support + integration tables
-- ====================================

-- Add industry column to businesses table
ALTER TABLE businesses
ADD COLUMN IF NOT EXISTS industry VARCHAR(50) DEFAULT 'general';

-- Add comment for valid industry values
COMMENT ON COLUMN businesses.industry IS 'Valid values: oil_logistics, ecommerce, restaurant, service, construction, general';

-- ====================================
-- INTEGRATIONS TABLE (OAuth connections)
-- ====================================
CREATE TABLE IF NOT EXISTS integrations (
    id SERIAL PRIMARY KEY,
    business_id INTEGER NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- 'shopify', 'quickbooks', 'samsara', 'fleetpanda', 'toast', 'sysco', etc.
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMPTZ,
    metadata JSONB, -- Store provider-specific config (shop_domain for Shopify, API endpoints, etc.)
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Ensure one active connection per business per provider
CREATE UNIQUE INDEX IF NOT EXISTS idx_integrations_unique_active
ON integrations(business_id, provider)
WHERE is_active = true;

-- ====================================
-- RAW_DATA TABLE (before processing)
-- ====================================
CREATE TABLE IF NOT EXISTS raw_data (
    id SERIAL PRIMARY KEY,
    integration_id INTEGER NOT NULL REFERENCES integrations(id) ON DELETE CASCADE,
    data_type VARCHAR(50) NOT NULL, -- 'product', 'sale', 'delivery', 'invoice', 'vehicle', 'route', etc.
    raw_json JSONB NOT NULL,
    processed BOOLEAN DEFAULT false,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index for querying unprocessed data
CREATE INDEX IF NOT EXISTS idx_raw_data_processed ON raw_data(processed, created_at);
CREATE INDEX IF NOT EXISTS idx_raw_data_integration ON raw_data(integration_id);
CREATE INDEX IF NOT EXISTS idx_raw_data_type ON raw_data(data_type);

-- ====================================
-- SYNC_LOGS TABLE (track data pipeline health)
-- ====================================
CREATE TABLE IF NOT EXISTS sync_logs (
    id SERIAL PRIMARY KEY,
    integration_id INTEGER REFERENCES integrations(id) ON DELETE CASCADE,
    sync_type VARCHAR(50), -- 'products', 'sales', 'deliveries', 'vehicles', 'routes', etc.
    status VARCHAR(20), -- 'success', 'failed', 'partial'
    records_synced INTEGER,
    error_message TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index for querying sync history
CREATE INDEX IF NOT EXISTS idx_sync_logs_integration ON sync_logs(integration_id);
CREATE INDEX IF NOT EXISTS idx_sync_logs_status ON sync_logs(status);
CREATE INDEX IF NOT EXISTS idx_sync_logs_created ON sync_logs(created_at DESC);

-- ====================================
-- INDEXES FOR PERFORMANCE
-- ====================================
CREATE INDEX IF NOT EXISTS idx_integrations_business_id ON integrations(business_id);
CREATE INDEX IF NOT EXISTS idx_integrations_provider ON integrations(provider);
CREATE INDEX IF NOT EXISTS idx_businesses_industry ON businesses(industry);
