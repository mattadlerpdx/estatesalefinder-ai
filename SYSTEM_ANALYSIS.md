# EstateSaleFinder.ai - System Architecture Analysis

**Date**: January 2025
**Question**: Are we implementing the hybrid storage model for scraped sales?

---

## 🔍 INVESTIGATION RESULTS

### ❌ Current State: NOT IMPLEMENTED YET

**What We Have:**
```
✅ Scraper (working)
✅ Redis cache (6-hour TTL)
✅ In-memory aggregation
❌ NO PostgreSQL persistence for scraped sales
❌ NO external_sales table
❌ NO ability to review scraped sales
❌ NO ability to add scraped sales to itineraries
```

### Current Data Flow:

```
┌──────────────────────────────────────────────────────┐
│ 1. Scraper fetches from estatesale-finder.com       │
└──────────────────────────────────────────────────────┘
                        ↓
┌──────────────────────────────────────────────────────┐
│ 2. Stores in Redis (6 hours)                        │
│    Key: "sales:portland:OR"                         │
│    Value: []ScrapedSale (JSON blob)                 │
└──────────────────────────────────────────────────────┘
                        ↓
┌──────────────────────────────────────────────────────┐
│ 3. API fetches from Redis + PostgreSQL              │
│    - Redis: Scraped sales (ephemeral)               │
│    - PostgreSQL: Owned sales (permanent)            │
└──────────────────────────────────────────────────────┘
                        ↓
┌──────────────────────────────────────────────────────┐
│ 4. Combines in-memory                               │
│    aggregatedSales = owned + scraped                │
└──────────────────────────────────────────────────────┘
                        ↓
┌──────────────────────────────────────────────────────┐
│ 5. Returns to frontend                              │
└──────────────────────────────────────────────────────┘
```

### ⚠️ THE PROBLEM:

**After 6 hours, scraped sales disappear from Redis.**

This means:
1. ❌ Can't review a sale you visited yesterday (cache expired)
2. ❌ Can't save scraped sales to itinerary (no permanent ID)
3. ❌ Can't track which scraped sales are high quality
4. ❌ No historical data for AI recommendations

---

## 📊 What's Actually in Our Database

### Existing Tables (from migrations/001_initial_schema.sql):

```sql
✅ users                   -- User accounts
✅ user_profiles          -- Extended user info
✅ estate_sales           -- OWNED sales only
✅ sale_images            -- Images for owned sales
✅ sale_items             -- Items in owned sales
✅ saved_sales            -- User favorites (only works for owned)
✅ professionals          -- Estate sale companies
✅ reviews                -- Reviews for PROFESSIONALS only (not sales!)
✅ subscription_plans     -- Pricing tiers
✅ user_subscriptions    -- User plan assignments

❌ external_sales         -- MISSING!
❌ sale_reviews           -- MISSING!
❌ itineraries            -- MISSING!
❌ itinerary_stops        -- MISSING!
```

### Current Review Table (WRONG):

```sql
-- This is for reviewing PROFESSIONALS, not SALES!
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    professional_id INTEGER REFERENCES professionals(id),  -- ⚠️ Wrong!
    reviewer_id INTEGER REFERENCES users(id),
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    review_text TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Problem**: Can't review individual sales, only companies.

---

## 🏗️ What We NEED to Build

### Phase 1: External Sales Persistence

#### 1. New Database Table:

```sql
-- Store minimal permanent record of scraped sales
CREATE TABLE external_sales (
  id SERIAL PRIMARY KEY,
  external_id VARCHAR(255) UNIQUE NOT NULL,  -- "estatesale-finder-15436"
  source VARCHAR(100) NOT NULL,              -- "EstateSale-Finder.com"
  source_url TEXT NOT NULL,                  -- Deep link

  -- Minimal searchable data
  title TEXT NOT NULL,
  address TEXT,
  city VARCHAR(100),
  state VARCHAR(2),
  zip_code VARCHAR(10),
  start_date TIMESTAMP,
  end_date TIMESTAMP,

  -- Metadata
  first_seen_at TIMESTAMP DEFAULT NOW(),
  last_scraped_at TIMESTAMP DEFAULT NOW(),

  -- Indexes
  CONSTRAINT external_sales_external_id_unique UNIQUE (external_id)
);

CREATE INDEX idx_external_sales_external_id ON external_sales(external_id);
CREATE INDEX idx_external_sales_location ON external_sales(city, state, start_date);
CREATE INDEX idx_external_sales_source ON external_sales(source);
```

#### 2. Repository Layer:

```go
// File: backend/internal/domain/sale/externalSaleRepo.go
package sale

import "time"

type ExternalSaleRepository interface {
    // Upsert creates or updates minimal record
    Upsert(sale *ScrapedSale) error

    // GetByExternalID looks up by external_id
    GetByExternalID(externalID string) (*ExternalSale, error)

    // GetReviewStats gets aggregated review data
    GetReviewStats(externalID string) (*ReviewStats, error)
}

type ExternalSale struct {
    ID           int       `json:"id"`
    ExternalID   string    `json:"external_id"`
    Source       string    `json:"source"`
    SourceURL    string    `json:"source_url"`
    Title        string    `json:"title"`
    Address      string    `json:"address"`
    City         string    `json:"city"`
    State        string    `json:"state"`
    ZipCode      string    `json:"zip_code"`
    StartDate    time.Time `json:"start_date"`
    EndDate      time.Time `json:"end_date"`
    FirstSeenAt  time.Time `json:"first_seen_at"`
    LastScraped  time.Time `json:"last_scraped_at"`
}

type ReviewStats struct {
    AverageRating float64 `json:"average_rating"`
    TotalReviews  int     `json:"total_reviews"`
    WorthItPercent float64 `json:"worth_it_percent"`
}
```

#### 3. Update Scraper:

```go
// File: backend/internal/infrastructure/scraper/scraper.go

func (s *ScraperService) GetSalesByLocation(city, state string) ([]sale.ScrapedSale, error) {
    cacheKey := s.getCacheKey(city, state)

    // 1. Try cache first
    if s.cache.IsEnabled() {
        var cachedSales []sale.ScrapedSale
        err := s.cache.Get(cacheKey, &cachedSales)
        if err == nil {
            log.Printf("✓ Cache HIT: %s (%d sales)", cacheKey, len(cachedSales))
            return cachedSales, nil
        }
    }

    // 2. Cache MISS - scrape
    sales, err := s.scrapeEstateSaleFinder(city, state)
    if err != nil {
        return nil, err
    }

    // 3. 🆕 PERSIST TO DATABASE (NEW!)
    for _, scrapedSale := range sales {
        if err := s.externalSaleRepo.Upsert(&scrapedSale); err != nil {
            log.Printf("Warning: Failed to persist external sale %s: %v",
                scrapedSale.ExternalID, err)
            // Don't fail - just log warning
        }
    }

    // 4. Store in cache
    if s.cache.IsEnabled() {
        s.cache.Set(cacheKey, sales, s.cacheTTL)
    }

    return sales, nil
}
```

### Phase 2: Review System

#### 1. New Sale Reviews Table:

```sql
-- Reviews for INDIVIDUAL SALES (not companies)
CREATE TABLE sale_reviews (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,

  -- Polymorphic reference: ONE will be NULL
  estate_sale_id INTEGER REFERENCES estate_sales(id) ON DELETE CASCADE,
  external_sale_id INTEGER REFERENCES external_sales(id) ON DELETE CASCADE,

  -- Rating data
  rating INTEGER CHECK (rating >= 1 AND rating <= 5) NOT NULL,
  worth_it BOOLEAN,
  quality_rating INTEGER CHECK (quality_rating >= 1 AND quality_rating <= 5),
  pricing_rating INTEGER CHECK (pricing_rating >= 1 AND pricing_rating <= 5),
  organization_rating INTEGER CHECK (organization_rating >= 1 AND organization_rating <= 5),
  crowd_rating INTEGER CHECK (crowd_rating >= 1 AND crowd_rating <= 5),

  comment TEXT,
  photos TEXT[], -- Array of image URLs

  -- Verification
  verified_visit BOOLEAN DEFAULT FALSE,
  attended_at TIMESTAMP,

  created_at TIMESTAMP DEFAULT NOW(),

  -- Constraints
  UNIQUE(user_id, estate_sale_id),
  UNIQUE(user_id, external_sale_id),
  CHECK (
    (estate_sale_id IS NOT NULL AND external_sale_id IS NULL) OR
    (estate_sale_id IS NULL AND external_sale_id IS NOT NULL)
  )
);

CREATE INDEX idx_sale_reviews_estate_sale ON sale_reviews(estate_sale_id);
CREATE INDEX idx_sale_reviews_external_sale ON sale_reviews(external_sale_id);
CREATE INDEX idx_sale_reviews_user ON sale_reviews(user_id);
```

### Phase 3: Itinerary System

#### 1. Itineraries Tables:

```sql
-- User's saved itineraries
CREATE TABLE itineraries (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,                    -- "Saturday Estate Sale Tour"
  date DATE NOT NULL,
  start_time TIME,
  end_time TIME,
  starting_address TEXT,
  starting_lat DECIMAL(10, 8),
  starting_lng DECIMAL(11, 8),
  total_distance_miles DECIMAL(10, 2),
  total_duration_minutes INTEGER,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Individual stops in itinerary
CREATE TABLE itinerary_stops (
  id SERIAL PRIMARY KEY,
  itinerary_id INTEGER REFERENCES itineraries(id) ON DELETE CASCADE,

  -- Polymorphic reference to sale
  estate_sale_id INTEGER REFERENCES estate_sales(id) ON DELETE CASCADE,
  external_sale_id INTEGER REFERENCES external_sales(id) ON DELETE CASCADE,

  stop_order INTEGER NOT NULL,           -- 1, 2, 3...
  estimated_arrival TIME,
  estimated_duration INTEGER,            -- minutes
  notes TEXT,

  created_at TIMESTAMP DEFAULT NOW(),

  CHECK (
    (estate_sale_id IS NOT NULL AND external_sale_id IS NULL) OR
    (estate_sale_id IS NULL AND external_sale_id IS NOT NULL)
  )
);

CREATE INDEX idx_itinerary_stops_itinerary ON itinerary_stops(itinerary_id);
CREATE INDEX idx_itinerary_stops_order ON itinerary_stops(itinerary_id, stop_order);
```

---

## 📝 Implementation Checklist

### ✅ Already Built:
- [x] Scraper service
- [x] Redis cache integration
- [x] ScrapedSale struct
- [x] AggregatedSale struct
- [x] API endpoint for aggregated sales
- [x] Frontend display with source attribution

### ❌ Need to Build:

#### Database:
- [ ] Create migration: 002_external_sales.sql
- [ ] Create migration: 003_sale_reviews.sql
- [ ] Create migration: 004_itineraries.sql

#### Backend:
- [ ] ExternalSaleRepository interface
- [ ] PostgreSQL implementation of ExternalSaleRepository
- [ ] Update scraper to persist external sales
- [ ] SaleReview domain model
- [ ] SaleReviewRepository interface
- [ ] SaleReviewService (business logic)
- [ ] Review API endpoints (POST, GET)
- [ ] Itinerary domain model
- [ ] ItineraryRepository interface
- [ ] ItineraryService (route optimization logic)
- [ ] Itinerary API endpoints

#### Frontend:
- [ ] Review form component
- [ ] Review display component
- [ ] Itinerary builder UI
- [ ] Map integration (Leaflet)
- [ ] Drag-drop route ordering

---

## 🎯 Priority Order (What to Build First)

### Week 1: External Sales Persistence
**Goal**: Stop losing scraped sales data after 6 hours

1. Create `external_sales` table
2. Build `ExternalSaleRepository`
3. Update scraper to persist on every scrape
4. Test: Scrape → Wait 7 hours → Data still exists in DB

**Benefit**: Enables all future features (reviews, itineraries)

### Week 2: Review System
**Goal**: Let users rate sales quality

1. Create `sale_reviews` table
2. Build review API endpoints
3. Build review UI components
4. Display review stats on sale cards

**Benefit**: Builds community, provides quality signals

### Week 3: Itinerary Builder
**Goal**: Let users plan routes

1. Create `itineraries` and `itinerary_stops` tables
2. Build "Add to Route" button
3. Build manual itinerary builder (drag-drop)
4. Add Google Maps integration

**Benefit**: Core differentiator vs. competitors

### Week 4: AI Route Optimization
**Goal**: Auto-optimize routes

1. Integrate Google Directions API
2. Build time-slot scheduling algorithm
3. Add AI suggestions (OpenAI)

**Benefit**: Full "AI-powered" experience

---

## 💾 Storage Cost Analysis

### Current System:
```
Redis only (ephemeral):
- 9 sales × 2KB each = 18KB in Redis
- Expires after 6 hours
- Cost: $0 (included in Upstash free tier)
- Problem: No persistence for reviews/itineraries
```

### After Implementation:
```
Redis (display cache):
- 9 sales × 2KB each = 18KB
- Expires after 6 hours
- Cost: $0

PostgreSQL (minimal persistence):
- 9 sales × 200 bytes each = 1.8KB
- Permanent
- Cost: ~$0.0001/month (negligible)

Total: Still basically free!
```

### At Scale (10,000 scraped sales):
```
Redis: 20MB (cache)
PostgreSQL: 2MB (minimal records)
Reviews: ~100KB (if 10% reviewed)
Total: ~22MB
Cost: <$1/month
```

---

## 🔑 Key Insights

### What We Built Right:
1. ✅ Clean separation (ScrapedSale vs. Sale)
2. ✅ Aggregation layer (AggregatedSale)
3. ✅ Source attribution (always link back)
4. ✅ Cache-through pattern (fast reads)

### What We Missed:
1. ❌ No permanent storage for scraped sales
2. ❌ No way to review scraped sales
3. ❌ No way to add scraped sales to itineraries
4. ❌ No historical data for AI

### Why It Matters:
**Without persistence, scraped sales are "read-only ephemeral data".**
**With persistence, scraped sales become "first-class citizens".**

Users can:
- Review them (build quality database)
- Save them to itineraries (trip planning)
- Favorite them (personalization)
- Get recommendations (AI needs history)

---

## 📋 Next Steps (Immediate)

1. **Create migration file**: `002_external_sales_and_reviews.sql`
2. **Build repository**: `externalSaleRepo.go`
3. **Update scraper**: Add persistence call after scraping
4. **Test**: Verify data survives cache expiry

Then we can build reviews and itineraries on top of this foundation.

---

## 🚀 The Vision

**Today:**
```
User visits site → Sees 13 sales → Leaves
```

**After Implementation:**
```
User visits site → Sees 13 sales with reviews
                 → Adds 7 to route
                 → AI optimizes itinerary
                 → User visits sales
                 → Leaves reviews
                 → Helps future users
                 → Community grows
```

**That's the flywheel. We need persistence to enable it.**

---

**Answer to your question: NO, we're not implementing the hybrid storage model yet. But now you know exactly what's needed! 🎯**
