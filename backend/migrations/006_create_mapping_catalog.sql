-- Migration 006: Create Mapping Catalog (Phase 2A)
-- Purpose: Store approved column mappings, entity links, and relationships
-- Design: Industry-standard Bronze→Silver→Gold with human-in-the-loop
-- Reference: INDUSTRY_STANDARD_IMPLEMENTATION.md

-- ============================================================================
-- SCHEMA VERSIONS: Track detected schemas for each CSV upload
-- ============================================================================
CREATE TABLE schema_versions (
  schema_version_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id INT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  integration_id INT REFERENCES integrations(id) ON DELETE CASCADE,
  source_name TEXT NOT NULL,              -- 'Cannabis Inventory - Jan 2025', 'Shopify Products', etc.
  detected_at TIMESTAMPTZ DEFAULT NOW(),
  headers JSONB NOT NULL,                 -- ['Item Name', 'SKU', 'Current Stock', ...]
  data_types JSONB,                       -- {'Item Name': 'string', 'Current Stock': 'int', ...}
  row_count INT,
  profile JSONB,                          -- {'nulls': 0, 'distinct_skus': 45, 'avg_price': 2.50, ...}
  CONSTRAINT chk_headers_array CHECK (jsonb_typeof(headers) = 'array')
);

CREATE INDEX idx_schema_versions_business ON schema_versions(business_id, detected_at DESC);
CREATE INDEX idx_schema_versions_integration ON schema_versions(integration_id);
CREATE INDEX idx_schema_versions_source_name ON schema_versions(business_id, source_name);

COMMENT ON TABLE schema_versions IS 'Snapshot of detected CSV schema for each upload (Bronze layer metadata)';
COMMENT ON COLUMN schema_versions.source_name IS 'User-friendly name for this data source (e.g., filename or integration name)';
COMMENT ON COLUMN schema_versions.headers IS 'Array of original CSV column names';
COMMENT ON COLUMN schema_versions.profile IS 'Statistical profile: null%, distinct values, min/max, etc.';

-- ============================================================================
-- COLUMN MAPPINGS: Store approved mappings (reusable contract)
-- ============================================================================
CREATE TABLE column_mappings (
  mapping_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id INT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  source_name TEXT NOT NULL,
  schema_version_id UUID REFERENCES schema_versions(schema_version_id) ON DELETE CASCADE,
  source_column TEXT NOT NULL,            -- 'Item Name' (original CSV header)
  canonical_field TEXT NOT NULL,          -- 'product_name' (standardized field name)
  confidence NUMERIC CHECK (confidence BETWEEN 0 AND 1),
  approved_by TEXT,                       -- Firebase UID of user who approved
  approved_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(business_id, source_name, source_column)
);

CREATE INDEX idx_column_mappings_business ON column_mappings(business_id, source_name);
CREATE INDEX idx_column_mappings_schema ON column_mappings(schema_version_id);
CREATE INDEX idx_column_mappings_approved ON column_mappings(business_id, approved_at DESC NULLS LAST);

COMMENT ON TABLE column_mappings IS 'Approved column mappings (contract) - reused for subsequent uploads with same schema';
COMMENT ON COLUMN column_mappings.source_name IS 'Data source identifier - matches schema_versions.source_name';
COMMENT ON COLUMN column_mappings.canonical_field IS 'Standardized field name from canonical data model';
COMMENT ON COLUMN column_mappings.confidence IS 'AI confidence score (0-1) - 1.0 means user-approved or exact match';

-- ============================================================================
-- ENTITY LINKS: Resolve entities across different CSVs
-- ============================================================================
CREATE TABLE entity_links (
  link_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id INT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  entity_type TEXT NOT NULL,              -- 'product', 'vendor', 'customer', 'location'
  source_name TEXT NOT NULL,              -- 'Cannabis Sales - Jan 2025'
  source_key TEXT NOT NULL,               -- 'Glass Jar - 1oz' (as it appears in CSV)
  entity_id UUID NOT NULL,                -- FK to product.product_id (or vendor_id, customer_id, etc.)
  confidence NUMERIC CHECK (confidence BETWEEN 0 AND 1),
  approved_by TEXT,                       -- Firebase UID
  approved_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(business_id, entity_type, source_name, source_key)
);

CREATE INDEX idx_entity_links_business ON entity_links(business_id, entity_type, source_name);
CREATE INDEX idx_entity_links_entity_id ON entity_links(entity_id);
CREATE INDEX idx_entity_links_source_key ON entity_links(business_id, entity_type, source_key);

COMMENT ON TABLE entity_links IS 'Maps source values to canonical entity IDs (e.g., "Glass Jar 1oz" → product_id=123)';
COMMENT ON COLUMN entity_links.entity_type IS 'Type of entity being linked: product, vendor, customer, location';
COMMENT ON COLUMN entity_links.source_key IS 'Original value from CSV (e.g., product name with typos/variations)';
COMMENT ON COLUMN entity_links.entity_id IS 'UUID of the canonical entity in typed Silver table';

-- ============================================================================
-- RELATIONSHIPS: Track discovered relationships between datasets
-- ============================================================================
CREATE TABLE relationships (
  relationship_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id INT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  from_table TEXT NOT NULL,               -- 'sale' (source table)
  from_column TEXT NOT NULL,              -- 'product_name' (source column)
  to_table TEXT NOT NULL,                 -- 'product' (target table)
  to_column TEXT NOT NULL,                -- 'product_name' (target column)
  method TEXT NOT NULL,                   -- 'inclusion_dep', 'semantic', 'manual', 'fuzzy'
  confidence NUMERIC CHECK (confidence BETWEEN 0 AND 1),
  coverage_percent NUMERIC,               -- What % of from values exist in to? (e.g., 85.7)
  sample_matches JSONB,                   -- Example matches for user review: {'Glass Jar 1oz': 'Glass Jar - 1oz', ...}
  approved_by TEXT,                       -- Firebase UID
  approved_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(business_id, from_table, from_column, to_table, to_column)
);

CREATE INDEX idx_relationships_business ON relationships(business_id, from_table);
CREATE INDEX idx_relationships_approved ON relationships(business_id, approved_at DESC NULLS LAST);

COMMENT ON TABLE relationships IS 'Discovered or approved relationships between datasets (enables auto-linking)';
COMMENT ON COLUMN relationships.method IS 'Detection method: inclusion_dep (set overlap), semantic (embedding similarity), manual (user-specified), fuzzy (string similarity)';
COMMENT ON COLUMN relationships.coverage_percent IS 'Percentage of from_column values that exist in to_column (Jaccard similarity)';
COMMENT ON COLUMN relationships.sample_matches IS 'Sample value pairs for user to review before approval';

-- ============================================================================
-- DATA QUALITY RESULTS: Track quality checks per batch
-- ============================================================================
CREATE TABLE quality_results (
  result_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id INT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  integration_id INT REFERENCES integrations(id) ON DELETE CASCADE,
  check_name TEXT NOT NULL,               -- 'row_count_sanity', 'uniqueness_check', 'fk_coverage', 'type_validity'
  check_type TEXT NOT NULL,               -- 'error', 'warning', 'info'
  passed BOOLEAN NOT NULL,
  message TEXT,
  details JSONB,                          -- {'expected': 100, 'actual': 95, 'missing_values': [...]}
  checked_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_quality_results_integration ON quality_results(integration_id, checked_at DESC);
CREATE INDEX idx_quality_results_business ON quality_results(business_id, passed, checked_at DESC);

COMMENT ON TABLE quality_results IS 'Data quality check results (row counts, FK coverage, type validity, etc.)';
COMMENT ON COLUMN quality_results.check_type IS 'Severity: error (blocks Silver), warning (allow but notify), info (logged only)';

-- ============================================================================
-- GRANT PERMISSIONS (if using restricted DB user)
-- ============================================================================
-- GRANT SELECT, INSERT, UPDATE, DELETE ON schema_versions TO cadence_api;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON column_mappings TO cadence_api;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON entity_links TO cadence_api;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON relationships TO cadence_api;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON quality_results TO cadence_api;
