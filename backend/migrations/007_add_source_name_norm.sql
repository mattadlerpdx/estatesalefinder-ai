-- Migration 007: Add source_name_norm for case-insensitive lookups
-- Purpose: Normalize source names for reliable mapping reuse (cannabis_inventory.csv == Cannabis_Inventory.CSV)

-- Add normalized source name columns
ALTER TABLE schema_versions
ADD COLUMN source_name_norm TEXT;

ALTER TABLE column_mappings
ADD COLUMN source_name_norm TEXT;

-- Backfill normalized values from existing source_name
UPDATE schema_versions
SET source_name_norm = LOWER(TRIM(REGEXP_REPLACE(source_name, '\.(csv|xlsx|xls)$', '', 'i')));

UPDATE column_mappings
SET source_name_norm = LOWER(TRIM(REGEXP_REPLACE(source_name, '\.(csv|xlsx|xls)$', '', 'i')));

-- Create indexes for fast normalized lookups
CREATE INDEX idx_schema_versions_bz_srcnorm ON schema_versions(business_id, source_name_norm);
CREATE INDEX idx_column_mappings_bz_srcnorm ON column_mappings(business_id, source_name_norm);
CREATE INDEX idx_schema_versions_detected_at ON schema_versions(business_id, detected_at DESC);

-- Add comments
COMMENT ON COLUMN schema_versions.source_name_norm IS 'Normalized source name for case-insensitive matching (lowercase, no extension)';
COMMENT ON COLUMN column_mappings.source_name_norm IS 'Normalized source name for case-insensitive matching (lowercase, no extension)';
