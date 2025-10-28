-- Migration 002: Add External Listing Support
-- Adds columns to estate_sales to support scraped/external listings
-- Creates unified "listings" table (owned + external)

-- 1. Add new columns to existing estate_sales table
ALTER TABLE estate_sales
  ADD COLUMN IF NOT EXISTS listing_type VARCHAR(20) DEFAULT 'owned',
  ADD COLUMN IF NOT EXISTS external_id VARCHAR(255),
  ADD COLUMN IF NOT EXISTS external_source VARCHAR(100),
  ADD COLUMN IF NOT EXISTS external_url TEXT,
  ADD COLUMN IF NOT EXISTS last_scraped_at TIMESTAMPTZ;

-- 2. Make seller_id optional (external listings don't have a seller)
ALTER TABLE estate_sales ALTER COLUMN seller_id DROP NOT NULL;

-- 3. Add constraint: listing_type must be 'owned' or 'external'
ALTER TABLE estate_sales
  ADD CONSTRAINT listings_type_check CHECK (
    listing_type IN ('owned', 'external')
  );

-- 4. Add constraint: owned listings must have seller_id, external must have external_id
ALTER TABLE estate_sales
  ADD CONSTRAINT listings_ownership_check CHECK (
    (listing_type = 'owned' AND seller_id IS NOT NULL AND external_id IS NULL) OR
    (listing_type = 'external' AND external_id IS NOT NULL AND seller_id IS NULL)
  );

-- 5. Add unique constraint on external_id (prevent duplicates)
-- NOTE: Must be non-partial index for ON CONFLICT to work
CREATE UNIQUE INDEX IF NOT EXISTS idx_estate_sales_external_id
  ON estate_sales(external_id);

-- 6. Add index on listing_type for filtering
CREATE INDEX IF NOT EXISTS idx_estate_sales_listing_type
  ON estate_sales(listing_type);

-- 7. Add index on external_source for analytics
CREATE INDEX IF NOT EXISTS idx_estate_sales_external_source
  ON estate_sales(external_source)
  WHERE external_source IS NOT NULL;

-- 8. Add index on last_scraped_at for refresh logic
CREATE INDEX IF NOT EXISTS idx_estate_sales_last_scraped
  ON estate_sales(last_scraped_at)
  WHERE last_scraped_at IS NOT NULL;

-- 9. Add composite index for location-based queries on external listings
CREATE INDEX IF NOT EXISTS idx_estate_sales_external_location
  ON estate_sales(city, state, start_date)
  WHERE listing_type = 'external';

-- 10. Update existing rows to have listing_type = 'owned'
UPDATE estate_sales SET listing_type = 'owned' WHERE listing_type IS NULL;

-- 11. Optional: Rename table to 'listings' (uncomment if desired)
-- This is optional - estate_sales still works, just less clear
-- ALTER TABLE estate_sales RENAME TO listings;
-- ALTER TABLE sale_images RENAME COLUMN sale_id TO listing_id;
-- ALTER TABLE sale_items RENAME COLUMN sale_id TO listing_id;
-- ALTER TABLE saved_sales RENAME COLUMN sale_id TO listing_id;

-- Migration complete!
-- Summary:
-- ✅ estate_sales table now supports both owned and external listings
-- ✅ external_id is unique and indexed
-- ✅ Constraints ensure data integrity
-- ✅ Indexes optimize queries for external listings
-- ✅ Backward compatible (existing data marked as 'owned')
