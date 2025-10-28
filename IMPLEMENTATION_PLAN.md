# Implementation Plan - Current Sprint

**Last Updated**: January 2025
**Status**: Active Development
**Goal**: Launch MVP with hybrid storage, auto-location, and map view

---

## 🎯 Current Sprint Goals

### Core Features to Build:
1. ✅ Hybrid storage model (Redis + PostgreSQL)
2. ✅ Auto-location detection (no user input needed)
3. ✅ Intelligent cache refresh (only scrape when stale)
4. ✅ Interactive map view with user auto-centering
5. ✅ Sale cards with source attribution

---

## 📋 Detailed Implementation Checklist

### Phase 1: Hybrid Storage Foundation (Week 1) 🔥 IN PROGRESS

#### Database Migration
- [ ] **Create migration file**: `002_add_external_listing_support.sql`
  - Add `listing_type` column (ENUM: 'owned' | 'external')
  - Add `external_id` column (VARCHAR 255, UNIQUE)
  - Add `external_source` column (VARCHAR 100)
  - Add `external_url` column (TEXT)
  - Add `last_scraped_at` column (TIMESTAMPTZ)
  - Make `seller_id` nullable (only owned listings have seller)
  - Add CHECK constraint for type validation
  - Add indexes for external lookups
  - **File**: `backend/migrations/002_add_external_listing_support.sql`
  - **Est. Time**: 30 min

- [ ] **Run migration**
  ```bash
  docker compose exec postgres psql -U postgres -d estatesale_db -f /docker-entrypoint-initdb.d/002_add_external_listing_support.sql
  ```
  - **Est. Time**: 5 min

#### Backend Repository Layer
- [ ] **Update Sale domain model**
  - Add new fields to `Sale` struct
  - Add helper methods: `IsExternal()`, `IsOwned()`
  - **File**: `backend/internal/domain/sale/sale.go`
  - **Est. Time**: 20 min

- [ ] **Add repository methods**
  - `GetLastScrapedTime(city, state string) (time.Time, error)`
  - `UpsertExternalListing(sale *ScrapedSale) error`
  - `GetExternalSalesByLocation(city, state string) ([]Sale, error)`
  - **File**: `backend/internal/infrastructure/db/postgres/saleRepo.go`
  - **Est. Time**: 45 min

#### Scraper Service Updates
- [ ] **Implement intelligent refresh logic**
  - Check Redis first (fastest)
  - If miss, check PostgreSQL for last_scraped_at
  - Only scrape if data is stale (>6 hours) or missing
  - Persist to PostgreSQL asynchronously (don't block)
  - **File**: `backend/internal/infrastructure/scraper/scraper.go`
  - **Est. Time**: 1 hour

- [ ] **Add error handling**
  - Graceful degradation if PostgreSQL fails
  - Log warnings, don't crash
  - Fall back to scraping if DB check fails
  - **Est. Time**: 20 min

#### Testing
- [ ] **Test scraper flow**
  - First request: Should scrape + persist
  - Second request (within 6 hours): Should load from Redis
  - Redis expired + DB fresh: Should load from DB → Redis
  - Redis expired + DB stale: Should re-scrape
  - **Est. Time**: 30 min

**Phase 1 Total Time**: ~4 hours

---

### Phase 2: Auto-Location Detection (Week 1-2)

#### Frontend Location Service
- [ ] **Create location utility**
  - Browser Geolocation API (primary)
  - IP geolocation fallback (ipapi.co)
  - LocalStorage caching (7-day TTL)
  - Default to Portland, OR
  - **File**: `frontend/lib/location.ts`
  - **Est. Time**: 1 hour

- [ ] **Add reverse geocoding**
  - Convert lat/lng → city/state
  - Use Nominatim (OpenStreetMap, free)
  - Cache results in localStorage
  - **File**: `frontend/lib/geocoding.ts`
  - **Est. Time**: 30 min

#### Landing Page Integration
- [ ] **Update homepage to auto-fetch location**
  - Call `getUserLocation()` on mount
  - Show loading state while detecting
  - Display detected location with "Change" option
  - **File**: `frontend/app/page.tsx`
  - **Est. Time**: 45 min

- [ ] **Add location permission prompt**
  - Friendly UI explaining why we need location
  - "Allow" → Precise results
  - "Deny" → Falls back to IP location
  - **Est. Time**: 30 min

#### Sales Page Integration
- [ ] **Auto-load sales based on location**
  - Remove default empty filters
  - Auto-populate city/state from geolocation
  - Show "Sales near you" header
  - **File**: `frontend/app/sales/page.tsx`
  - **Est. Time**: 30 min

**Phase 2 Total Time**: ~3.5 hours

---

### Phase 3: Interactive Map View (Week 2)

#### Map Component Setup
- [ ] **Install dependencies**
  ```bash
  cd frontend
  npm install leaflet react-leaflet @types/leaflet
  ```
  - **Est. Time**: 5 min

- [ ] **Create base map component**
  - MapContainer with Stadia Maps tiles
  - Auto-center on user location
  - Responsive height/width
  - **File**: `frontend/components/SalesMap.tsx`
  - **Est. Time**: 45 min

- [ ] **Add sale markers**
  - Pin for each sale location
  - Color-coded by type (owned = blue, scraped = green)
  - Clustered when zoomed out (react-leaflet-cluster)
  - **Est. Time**: 1 hour

- [ ] **Add marker popups**
  - Sale title, address, date
  - Quick preview card
  - "View Details" button
  - Link to external source if scraped
  - **Est. Time**: 45 min

#### Map Features
- [ ] **Add map controls**
  - Zoom buttons
  - "Center on me" button
  - Layer toggle (list view ↔ map view)
  - **Est. Time**: 30 min

- [ ] **Add search on map**
  - Drag map → "Search this area" button appears
  - Updates results based on visible bounds
  - **Est. Time**: 1 hour

- [ ] **Mobile optimization**
  - Touch gestures
  - Full-screen map option
  - Bottom sheet for sale details
  - **Est. Time**: 1.5 hours

**Phase 3 Total Time**: ~5.5 hours

---

### Phase 4: Review System (Week 3)

#### Database Schema
- [ ] **Create sale_reviews table**
  - `listing_id` references `listings(id)`
  - Rating fields (overall, quality, pricing, etc.)
  - `worth_it` boolean
  - Comment text
  - Photo URLs array
  - GPS verification fields
  - **File**: `backend/migrations/003_sale_reviews.sql`
  - **Est. Time**: 30 min

#### Backend API
- [ ] **Review endpoints**
  - POST `/api/sales/:id/reviews` - Create review
  - GET `/api/sales/:id/reviews` - List reviews
  - GET `/api/reviews/stats/:id` - Aggregate stats
  - **File**: `backend/internal/infrastructure/controllers/reviewHandler.go`
  - **Est. Time**: 2 hours

- [ ] **Review service logic**
  - Validate user attended sale (GPS check)
  - Prevent duplicate reviews
  - Calculate aggregate ratings
  - **File**: `backend/internal/domain/review/reviewService.go`
  - **Est. Time**: 1.5 hours

#### Frontend Components
- [ ] **Review form**
  - Star ratings (5 categories)
  - "Worth it?" toggle
  - Comment textarea
  - Photo upload (optional)
  - **File**: `frontend/components/ReviewForm.tsx`
  - **Est. Time**: 2 hours

- [ ] **Review display**
  - List of reviews with pagination
  - Average rating display
  - "Worth it" percentage
  - Helpful votes
  - **File**: `frontend/components/ReviewList.tsx`
  - **Est. Time**: 1.5 hours

**Phase 4 Total Time**: ~7.5 hours

---

### Phase 5: Itinerary Builder (Week 4)

#### Database Schema
- [ ] **Create itineraries tables**
  - `itineraries` (user's saved routes)
  - `itinerary_stops` (individual sales in route)
  - Foreign keys to `listings(id)` (works for both types!)
  - **File**: `backend/migrations/004_itineraries.sql`
  - **Est. Time**: 30 min

#### Backend API
- [ ] **Itinerary endpoints**
  - POST `/api/itineraries` - Create itinerary
  - GET `/api/itineraries` - List user's itineraries
  - PUT `/api/itineraries/:id` - Update route order
  - DELETE `/api/itineraries/:id` - Delete itinerary
  - POST `/api/itineraries/:id/stops` - Add sale to route
  - **Est. Time**: 2 hours

- [ ] **Route optimization service**
  - Integrate Google Directions API
  - Calculate optimal order
  - Estimate travel times
  - Account for sale hours
  - **File**: `backend/internal/domain/itinerary/optimizer.go`
  - **Est. Time**: 3 hours

#### Frontend Components
- [ ] **"Add to Route" button**
  - On each sale card
  - Shows current route count
  - Quick add without leaving page
  - **Est. Time**: 30 min

- [ ] **Itinerary builder page**
  - Drag-and-drop route ordering
  - Map with route line
  - Time estimates for each stop
  - Export to Google Maps
  - **File**: `frontend/app/itinerary/page.tsx`
  - **Est. Time**: 4 hours

**Phase 5 Total Time**: ~10 hours

---

## 🗓️ Timeline Summary

| Phase | Focus | Duration | Status |
|-------|-------|----------|--------|
| Phase 1 | Hybrid Storage | 4 hours | 🔥 Next |
| Phase 2 | Auto-Location | 3.5 hours | ⏳ Pending |
| Phase 3 | Map View | 5.5 hours | ⏳ Pending |
| Phase 4 | Reviews | 7.5 hours | ⏳ Pending |
| Phase 5 | Itineraries | 10 hours | ⏳ Pending |

**Total**: ~30 hours of development (1 week full-time or 2 weeks part-time)

---

## 🎯 Success Criteria

### Phase 1 Complete When:
- [ ] External sales persist in PostgreSQL
- [ ] Scraper checks DB before scraping
- [ ] Only re-scrapes after 6 hours
- [ ] Redis + PostgreSQL work together
- [ ] Performance: <10ms for cached requests

### Phase 2 Complete When:
- [ ] Browser detects user location automatically
- [ ] Falls back gracefully if denied
- [ ] Location cached in localStorage
- [ ] Sales auto-filtered to user's area
- [ ] "Near you" messaging displays

### Phase 3 Complete When:
- [ ] Map displays with modern tiles
- [ ] Auto-centers on user location
- [ ] All sales shown as markers
- [ ] Popups show sale details
- [ ] Works on mobile
- [ ] Clustering handles 100+ markers

### Phase 4 Complete When:
- [ ] Users can leave reviews
- [ ] Review stats aggregate correctly
- [ ] "Worth it" percentage displays
- [ ] Sale cards show rating stars
- [ ] GPS verification works

### Phase 5 Complete When:
- [ ] Users can create itineraries
- [ ] Route optimization works
- [ ] Map shows route line
- [ ] Time estimates display
- [ ] Export to Google Maps works

---

## 🚀 Quick Start (Next Session)

### To Begin Implementation:

```bash
# 1. Start your environment
cd /mnt/c/Users/matt/Desktop/Stuff/estatesalefinder-ai
make dev

# 2. Create migration file
touch backend/migrations/002_add_external_listing_support.sql

# 3. Start coding!
# Files to edit:
# - backend/migrations/002_add_external_listing_support.sql
# - backend/internal/domain/sale/sale.go
# - backend/internal/infrastructure/db/postgres/saleRepo.go
# - backend/internal/infrastructure/scraper/scraper.go
```

---

## 📝 Notes & Decisions

### Architecture Decisions:
- **Why unified table?** Industry standard (Airbnb, Zillow), simpler queries
- **Why 6-hour cache?** Balance freshness vs. scraping load
- **Why async persistence?** Don't block user requests
- **Why localStorage?** Reduce API calls, better UX

### API Costs (All Free Tier):
- IP Geolocation: ipapi.co (30k/month free)
- Reverse Geocoding: Nominatim OSM (unlimited, attribution required)
- Map Tiles: Stadia Maps (free with attribution)
- Google Directions: 2,500 requests/day free

### Technical Debt to Address Later:
- Add rate limiting to prevent scraper abuse
- Implement exponential backoff for failed scrapes
- Add monitoring/alerting for scraper health
- Optimize PostgreSQL indexes for review queries
- Add full-text search for sale descriptions

---

## 🔗 Related Documents

- **FEATURES_ROADMAP.md** - High-level feature strategy
- **HYBRID_STORAGE_STRATEGY.md** - Database architecture deep-dive
- **SYSTEM_ANALYSIS.md** - Current vs. needed state analysis
- **PROJECT_STATUS.md** - Overall project progress
- **ROADMAP.md** - 16-week development plan

---

**Ready to build! Start with Phase 1. 🚀**
