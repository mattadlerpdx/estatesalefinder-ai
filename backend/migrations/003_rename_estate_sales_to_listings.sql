-- Migration 003: Rename estate_sales to listings
-- Makes naming clearer: "listings" can be owned OR external (scraped)

-- 1. Rename the table
ALTER TABLE estate_sales RENAME TO listings;

-- 2. Rename foreign key columns in related tables
ALTER TABLE sale_images RENAME COLUMN sale_id TO listing_id;
ALTER TABLE sale_items RENAME COLUMN sale_id TO listing_id;
ALTER TABLE saved_sales RENAME COLUMN sale_id TO listing_id;

-- 3. Rename the saved_sales table to saved_listings for consistency
ALTER TABLE saved_sales RENAME TO saved_listings;

-- 4. Rename all indexes (PostgreSQL doesn't auto-rename them)
ALTER INDEX IF EXISTS estate_sales_pkey RENAME TO listings_pkey;
ALTER INDEX IF EXISTS idx_estate_sales_city_state RENAME TO idx_listings_city_state;
ALTER INDEX IF EXISTS idx_estate_sales_start_date RENAME TO idx_listings_start_date;
ALTER INDEX IF EXISTS idx_estate_sales_seller_id RENAME TO idx_listings_seller_id;
ALTER INDEX IF EXISTS idx_estate_sales_status RENAME TO idx_listings_status;
ALTER INDEX IF EXISTS idx_estate_sales_featured RENAME TO idx_listings_featured;
ALTER INDEX IF EXISTS idx_estate_sales_external_id RENAME TO idx_listings_external_id;
ALTER INDEX IF EXISTS idx_estate_sales_listing_type RENAME TO idx_listings_listing_type;
ALTER INDEX IF EXISTS idx_estate_sales_external_source RENAME TO idx_listings_external_source;
ALTER INDEX IF EXISTS idx_estate_sales_last_scraped RENAME TO idx_listings_last_scraped;
ALTER INDEX IF EXISTS idx_estate_sales_external_location RENAME TO idx_listings_external_location;

-- Migration complete!
-- Summary:
-- ✅ estate_sales → listings
-- ✅ sale_images.sale_id → listing_id
-- ✅ sale_items.sale_id → listing_id
-- ✅ saved_sales → saved_listings (with listing_id column)
-- ✅ All indexes renamed
