# Hybrid Storage Strategy - Industry Standard Approach

**Analysis Date**: January 2025

---

## 🏛️ Current Schema Analysis

### What You Have (Owned Sales):

```sql
estate_sales                  -- Full marketplace listings
├── seller_id                 -- Who listed it
├── title, description        -- Rich content
├── address, city, state      -- Full location
├── start_date, end_date      -- Dates
├── listing_tier              -- basic/featured/premium
├── payment_status            -- Revenue tracking
├── view_count, featured      -- Analytics
└── created_at, updated_at    -- Timestamps

sale_images (1-to-many)       -- Full image storage
sale_items (1-to-many)        -- Detailed inventory
saved_sales (many-to-many)    -- User favorites
reviews → professionals       -- Company reviews (wrong!)
```

**Total storage per owned sale**: ~50KB (with images/items)

---

## 🏗️ Industry Standard: Polymorphic "Sales" Pattern

### The Problem with Separate Tables:

❌ **Anti-Pattern (What NOT to do):**
```sql
-- Separate everything = duplication hell
estate_sales           (owned)
external_sales         (scraped)
reviews_for_estate
reviews_for_external
favorites_for_estate
favorites_for_external
itinerary_estate_stops
itinerary_external_stops
-- 2x the code, 2x the bugs
```

### ✅ Industry Standard: Single Polymorphic Table

**How Airbnb, Zillow, and aggregators do it:**

```sql
-- ONE unified "listings" table
-- Source-agnostic design
CREATE TABLE listings (
  id SERIAL PRIMARY KEY,

  -- Discriminator (tells us which type)
  listing_type VARCHAR(20) NOT NULL, -- 'owned' or 'external'

  -- External identifier (for scraped)
  external_id VARCHAR(255) UNIQUE,   -- "estatesale-finder-15436"
  external_source VARCHAR(100),       -- "EstateSale-Finder.com"
  external_url TEXT,                  -- Link to original

  -- Ownership (for owned listings)
  seller_id INTEGER REFERENCES users(id),

  -- SHARED fields (all listings have these)
  title VARCHAR(255) NOT NULL,
  description TEXT,
  address TEXT,
  city VARCHAR(100) NOT NULL,
  state VARCHAR(50) NOT NULL,
  zip_code VARCHAR(20),
  latitude DECIMAL(10, 8),
  longitude DECIMAL(11, 8),
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,

  -- Owned-only fields (NULL for external)
  sale_type VARCHAR(50),
  listing_tier VARCHAR(50),
  payment_status VARCHAR(50),
  amount_paid DECIMAL(10, 2),

  -- Metadata
  view_count INTEGER DEFAULT 0,
  featured BOOLEAN DEFAULT FALSE,
  first_seen_at TIMESTAMPTZ DEFAULT NOW(),
  last_scraped_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),

  -- Constraints
  CHECK (
    (listing_type = 'owned' AND seller_id IS NOT NULL AND external_id IS NULL) OR
    (listing_type = 'external' AND external_id IS NOT NULL AND seller_id IS NULL)
  )
);

CREATE UNIQUE INDEX idx_listings_external_id ON listings(external_id) WHERE external_id IS NOT NULL;
CREATE INDEX idx_listings_type ON listings(listing_type);
CREATE INDEX idx_listings_seller ON listings(seller_id) WHERE seller_id IS NOT NULL;
CREATE INDEX idx_listings_location ON listings(city, state, start_date);
CREATE INDEX idx_listings_external_source ON listings(external_source) WHERE external_source IS NOT NULL;
```

### Why This is Better:

1. **One codebase** - All features work for both types
2. **One query** - `SELECT * FROM listings WHERE city = 'Portland'` (no JOINs)
3. **One foreign key** - Reviews, favorites, itineraries all reference `listing_id`
4. **Easy migration** - External listings can become owned (if company signs up)

---

## 📊 Storage Comparison

### Approach 1: Separate Tables (What I Initially Suggested)

```sql
estate_sales (50KB per sale)
external_sales (200 bytes per sale)

10,000 owned sales:  500 MB
10,000 scraped sales: 2 MB
Total: 502 MB
```

**Pros**: Clear separation
**Cons**: Duplicate code, complex queries

### Approach 2: Unified Table (Industry Standard)

```sql
listings (mixed)

10,000 owned:  500 MB (full data)
10,000 scraped: 2 MB (minimal fields, rest NULL)
Total: 502 MB (same storage!)

But with SPARSE COLUMNS optimization:
PostgreSQL doesn't allocate space for NULL columns
Actual storage: ~502 MB (negligible overhead)
```

**Pros**: One codebase, simple queries, industry standard
**Cons**: Some NULL columns (minor)

---

## 🎯 Recommended Approach: Hybrid Table

### Migration: 002_unified_listings.sql

```sql
-- Rename existing table to keep owned sales
ALTER TABLE estate_sales RENAME TO listings;

-- Add discriminator and external fields
ALTER TABLE listings
  ADD COLUMN listing_type VARCHAR(20) NOT NULL DEFAULT 'owned',
  ADD COLUMN external_id VARCHAR(255),
  ADD COLUMN external_source VARCHAR(100),
  ADD COLUMN external_url TEXT,
  ADD COLUMN first_seen_at TIMESTAMPTZ DEFAULT NOW(),
  ADD COLUMN last_scraped_at TIMESTAMPTZ;

-- Make seller_id optional (NULL for external)
ALTER TABLE listings ALTER COLUMN seller_id DROP NOT NULL;

-- Add constraints
ALTER TABLE listings ADD CONSTRAINT listings_type_check CHECK (
  (listing_type = 'owned' AND seller_id IS NOT NULL AND external_id IS NULL) OR
  (listing_type = 'external' AND external_id IS NOT NULL AND seller_id IS NULL)
);

-- Add indexes for external lookups
CREATE UNIQUE INDEX idx_listings_external_id ON listings(external_id)
  WHERE external_id IS NOT NULL;

CREATE INDEX idx_listings_type ON listings(listing_type);

CREATE INDEX idx_listings_external_source ON listings(external_source)
  WHERE external_source IS NOT NULL;

-- Update foreign key references (keep backward compatible)
-- sale_images, sale_items, saved_sales still reference listings(id)
-- No changes needed!
```

### How It Works:

```go
// Owned sale insert
INSERT INTO listings (
  listing_type, seller_id, title, description,
  city, state, start_date, end_date
) VALUES (
  'owned', 123, 'Vintage Furniture Sale', 'Amazing...',
  'Portland', 'OR', '2025-11-01', '2025-11-03'
);

// External sale insert (minimal data)
INSERT INTO listings (
  listing_type, external_id, external_source, external_url,
  title, address, city, state, start_date, end_date
) VALUES (
  'external', 'estatesale-finder-15436', 'EstateSale-Finder.com',
  'https://www.estatesale-finder.com/viewsale.php?saleid=15436',
  'David Johnson Estate Sales', 'TBA', 'Portland', 'OR',
  '2025-11-04', '2025-11-06'
);
```

---

## 🔗 Unified Foreign Key References

### Reviews (No Polymorphic Hacks!)

```sql
CREATE TABLE sale_reviews (
  id SERIAL PRIMARY KEY,
  listing_id INTEGER REFERENCES listings(id) ON DELETE CASCADE,  -- ONE FK!
  user_id INTEGER REFERENCES users(id),

  rating INTEGER CHECK (rating >= 1 AND rating <= 5) NOT NULL,
  worth_it BOOLEAN,
  quality_rating INTEGER,
  pricing_rating INTEGER,
  organization_rating INTEGER,
  crowd_rating INTEGER,
  comment TEXT,
  photos TEXT[],

  verified_visit BOOLEAN DEFAULT FALSE,
  attended_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(user_id, listing_id)  -- One review per user per listing
);

CREATE INDEX idx_sale_reviews_listing ON sale_reviews(listing_id);
CREATE INDEX idx_sale_reviews_user ON sale_reviews(user_id);
```

### Itineraries

```sql
CREATE TABLE itineraries (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  date DATE NOT NULL,
  start_time TIME,
  end_time TIME,
  starting_address TEXT,
  total_distance_miles DECIMAL(10, 2),
  total_duration_minutes INTEGER,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE itinerary_stops (
  id SERIAL PRIMARY KEY,
  itinerary_id INTEGER REFERENCES itineraries(id) ON DELETE CASCADE,
  listing_id INTEGER REFERENCES listings(id) ON DELETE CASCADE,  -- ONE FK!

  stop_order INTEGER NOT NULL,
  estimated_arrival TIME,
  estimated_duration INTEGER,
  notes TEXT,

  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_itinerary_stops_itinerary ON itinerary_stops(itinerary_id);
CREATE INDEX idx_itinerary_stops_listing ON itinerary_stops(listing_id);
```

### Favorites (Update existing)

```sql
-- Already exists as saved_sales, just rename FK
ALTER TABLE saved_sales RENAME TO favorite_listings;
-- The sale_id column still references listings(id) - works perfectly!
```

---

## 🔍 Query Examples

### Get All Sales (Owned + External):

```sql
SELECT
  id, listing_type, title, city, state, start_date,
  CASE
    WHEN listing_type = 'external' THEN external_url
    ELSE '/sales/' || id::text
  END as view_url
FROM listings
WHERE city = 'Portland'
  AND start_date >= NOW()
ORDER BY start_date ASC;
```

### Get Sale with Review Stats:

```sql
SELECT
  l.*,
  COUNT(r.id) as review_count,
  AVG(r.rating) as avg_rating,
  SUM(CASE WHEN r.worth_it THEN 1 ELSE 0 END)::float / NULLIF(COUNT(r.id), 0) as worth_it_percent
FROM listings l
LEFT JOIN sale_reviews r ON r.listing_id = l.id
WHERE l.id = $1
GROUP BY l.id;
```

### User's Itinerary with Mixed Listings:

```sql
SELECT
  i.name as itinerary_name,
  s.stop_order,
  l.listing_type,
  l.title,
  l.city,
  l.external_url,
  CASE
    WHEN l.listing_type = 'owned' THEN true
    ELSE false
  END as is_our_listing
FROM itineraries i
JOIN itinerary_stops s ON s.itinerary_id = i.id
JOIN listings l ON l.id = s.listing_id
WHERE i.user_id = $1
ORDER BY s.stop_order;
```

---

## 📈 Data Flow (Unified Approach)

### Scraper Update:

```go
func (s *ScraperService) GetSalesByLocation(city, state string) ([]sale.Sale, error) {
    // 1. Check Redis cache
    cached, err := s.cache.Get(cacheKey)
    if err == nil {
        return cached, nil
    }

    // 2. Scrape external source
    scrapedSales, err := s.scrapeEstateSaleFinder(city, state)
    if err != nil {
        return nil, err
    }

    // 3. UPSERT to unified listings table
    for _, scraped := range scrapedSales {
        _, err := s.db.Exec(`
            INSERT INTO listings (
                listing_type, external_id, external_source, external_url,
                title, address, city, state, zip_code,
                start_date, end_date, last_scraped_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
            ON CONFLICT (external_id)
            DO UPDATE SET
                title = EXCLUDED.title,
                start_date = EXCLUDED.start_date,
                end_date = EXCLUDED.end_date,
                last_scraped_at = NOW()
        `, "external", scraped.ExternalID, scraped.SourceName, scraped.SourceURL,
           scraped.Title, scraped.Address, scraped.City, scraped.State, scraped.ZipCode,
           scraped.StartDate, scraped.EndDate)

        if err != nil {
            log.Printf("Warning: Failed to upsert external sale: %v", err)
        }
    }

    // 4. Cache for 6 hours
    s.cache.Set(cacheKey, scrapedSales, 6*time.Hour)

    return scrapedSales, nil
}
```

### API Query (Aggregated):

```go
func (h *SaleHandler) GetSales(w http.ResponseWriter, r *http.Request) {
    city := r.URL.Query().Get("city")
    state := r.URL.Query().Get("state")

    // Default to Portland
    if city == "" && state == "" {
        city, state = "Portland", "OR"
    }

    // ONE query gets both owned + external!
    sales, err := h.saleRepo.GetByLocation(city, state)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // Trigger scraper in background (async)
    go h.scraperService.GetSalesByLocation(city, state)

    json.NewEncoder(w).Encode(sales)
}
```

---

## 🎯 Migration Path (Non-Breaking)

### Step 1: Add Columns (Backward Compatible)

```sql
-- 002_add_external_support.sql
ALTER TABLE estate_sales
  ADD COLUMN listing_type VARCHAR(20) DEFAULT 'owned',
  ADD COLUMN external_id VARCHAR(255),
  ADD COLUMN external_source VARCHAR(100),
  ADD COLUMN external_url TEXT,
  ADD COLUMN last_scraped_at TIMESTAMPTZ;

-- Existing data automatically gets listing_type = 'owned'
-- No data loss!
```

### Step 2: Rename (Optional, for clarity)

```sql
-- 003_rename_to_listings.sql
ALTER TABLE estate_sales RENAME TO listings;
ALTER TABLE sale_images RENAME COLUMN sale_id TO listing_id;
-- Update your Go code to use "listings"
```

### Step 3: Update Code

```go
// Old (still works!)
type Sale struct {
    ID       int
    SellerID int
    Title    string
    // ...
}

// New (supports both)
type Listing struct {
    ID          int
    ListingType string  // "owned" or "external"

    // Owned fields
    SellerID    *int

    // External fields
    ExternalID     *string
    ExternalSource *string
    ExternalURL    *string

    // Shared fields
    Title       string
    City        string
    // ...
}
```

---

## 🏆 Industry Examples

### Zillow:
```
listings table:
├── listing_type: "mls" | "fsbo" | "zillow_owned"
├── mls_id (external)
├── owner_id (owned)
└── shared: address, price, beds, baths
```

### Airbnb:
```
listings table:
├── listing_type: "airbnb" | "vrbo_scraped" | "homeaway_scraped"
├── external_id
├── host_id (owned)
└── shared: title, location, price
```

### Indeed (Job Aggregator):
```
jobs table:
├── source: "direct" | "scraped_linkedin" | "scraped_monster"
├── external_id
├── employer_id (direct postings)
└── shared: title, description, salary
```

**Pattern**: Everyone uses unified tables with discriminator columns.

---

## ✅ Recommendation

**Use the unified `listings` table approach:**

### Benefits:
1. ✅ Industry standard (Airbnb, Zillow, Indeed all do this)
2. ✅ One codebase for all features
3. ✅ Simple queries (no polymorphic JOINs)
4. ✅ Easy data migration (external → owned if company signs up)
5. ✅ Backward compatible (add columns, don't break existing)
6. ✅ Same storage cost as separate tables
7. ✅ Future-proof (add more sources later)

### Drawbacks:
- Some NULL columns (minor, PostgreSQL handles well)
- Need discriminator checks in code (easy with helper methods)

---

## 📋 Implementation Checklist

- [ ] Create migration `002_add_external_support.sql`
- [ ] Update `Sale` struct to `Listing` with type discriminator
- [ ] Update scraper to UPSERT to listings table
- [ ] Create `sale_reviews` table (references `listing_id`)
- [ ] Create `itineraries` and `itinerary_stops` tables
- [ ] Update API queries to filter by `listing_type` when needed
- [ ] Test: Insert owned sale, insert external sale, query both

---

**This is the industry-standard way. Let's implement it! 🚀**
