# Testing Documentation - EstateSaleFinder.ai

This document describes how to run tests for the EstateSaleFinder.ai application after the Sale → Listing refactoring.

## Overview

The application has comprehensive test coverage for both backend (Go) and frontend (TypeScript/React):

- **Backend Tests**: Integration tests using testify/suite with real PostgreSQL database
- **Frontend Tests**: Type tests and API integration tests using Jest

## Backend Testing

### Prerequisites

1. PostgreSQL database running (required for integration tests)
2. Set `DATABASE_URL` environment variable
3. Go 1.21+ installed

### Running Backend Tests

#### All Tests
```bash
cd backend
go test -v ./...
```

#### Listing Domain Tests Only
```bash
cd backend
go test -v ./internal/domain/listing/...
```

#### Specific Test Suites
```bash
# Unit tests only (no database required)
cd backend
go test -v ./internal/domain/listing -run TestListing

# Integration tests (requires database)
cd backend
export DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
go test -v ./internal/domain/listing -run TestListingIntegrationSuite

# Scraper integration tests
cd backend
export DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
go test -v ./internal/infrastructure/scraper -run TestScraperIntegrationSuite
```

#### With Coverage
```bash
cd backend
go test -v -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Backend Test Structure

#### Unit Tests (`listing_test.go`)
Tests basic struct creation and conversions without database:
- ✅ Listing creation (owned and external)
- ✅ ListingImage creation
- ✅ ListingFilters structure
- ✅ ScrapedListing creation
- ✅ ScrapedListing ↔ Listing conversions
- ✅ AggregatedListing conversions

#### Integration Tests (`listing_integration_test.go`)
Tests full CRUD operations with real database:
- ✅ **TestListingCRUD**: Create, Read, Update, Delete operations
- ✅ **TestListingWithImages**: Image operations (add, delete, set primary)
- ✅ **TestListingFilters**: Filtering by city, state, sale type, featured status
- ✅ **TestListingFilters**: Pagination (limit, offset)
- ✅ **TestExternalListingConversion**: ScrapedListing conversion and persistence

#### Scraper Integration Tests (`scraper_integration_test.go`)
Tests hybrid storage (Redis + PostgreSQL):
- ✅ Initial scrape (empty cache, empty DB)
- ✅ Redis cache hit (fast retrieval)
- ✅ PostgreSQL fallback (when Redis expires)
- ✅ 6-hour refresh (re-scrape stale data)
- ✅ External sale upsert (insert + update)
- ✅ Location filtering

### Test Data Cleanup

All tests use the `cleanupTestData()` method to:
- Delete test listings (titles starting with "Test:")
- Delete external listings with test IDs
- Prevent test data pollution

## Frontend Testing

### Prerequisites

1. Node.js 18+ installed
2. npm or pnpm installed
3. Backend API running (for integration tests)

### Running Frontend Tests

#### Install Dependencies
```bash
cd frontend
npm install
```

#### All Tests
```bash
cd frontend
npm test
```

#### Watch Mode (for development)
```bash
cd frontend
npm run test:watch
```

#### With Coverage
```bash
cd frontend
npm run test:coverage
```

#### Specific Tests
```bash
# Type tests only
cd frontend
npm test -- types.test.ts

# API integration tests only
cd frontend
npm test -- api.test.ts
```

### Frontend Test Structure

#### Type Tests (`__tests__/types.test.ts`)
Validates TypeScript type definitions:
- ✅ Listing interface structure (owned listings)
- ✅ Listing interface structure (scraped listings)
- ✅ ListingDetail interface extension
- ✅ Array of Listings typing
- ✅ Optional fields handling
- ✅ Filtering and mapping operations

#### API Integration Tests (`__tests__/api.test.ts`)
Tests API endpoints with real backend:
- ✅ **GET /api/sales**: Fetch all listings
- ✅ **GET /api/sales?city=X**: Filter by city
- ✅ **GET /api/sales?state=X**: Filter by state
- ✅ **GET /api/sales?featured=true**: Filter featured
- ✅ **GET /api/sales?limit=X**: Pagination
- ✅ **GET /api/sales/:id**: Fetch single listing
- ✅ **404 handling**: Non-existent listing
- ✅ **Type validation**: is_scraped field
- ✅ **Date validation**: ISO string parsing
- ✅ **Image structure**: Correct image format
- ✅ **Error handling**: Network errors

**Note**: API integration tests require the backend to be running:
```bash
# Terminal 1: Start backend
cd backend
go run cmd/api/main.go

# Terminal 2: Run frontend tests
cd frontend
npm test -- api.test.ts
```

## Refactoring Verification Checklist

Use this checklist to verify the Sale → Listing refactoring is working:

### Backend Verification
- [ ] Backend compiles without errors: `go build ./...`
- [ ] Unit tests pass: `go test -v ./internal/domain/listing -run TestListing`
- [ ] Integration tests pass: `go test -v ./internal/domain/listing -run TestListingIntegrationSuite`
- [ ] Scraper tests pass: `go test -v ./internal/infrastructure/scraper`
- [ ] All package imports use `listing` package
- [ ] All structs use `Listing`, `ListingImage`, `ListingFilters` types
- [ ] Repository interface uses correct types
- [ ] Service layer uses consistent variable names (`l` for listing)

### Frontend Verification
- [ ] TypeScript compiles: `npm run build`
- [ ] Type tests pass: `npm test -- types.test.ts`
- [ ] API tests pass (with backend running): `npm test -- api.test.ts`
- [ ] All interfaces use `Listing` type
- [ ] Components compile with new types
- [ ] No TypeScript errors in IDE

### End-to-End Verification
- [ ] Backend starts successfully
- [ ] Frontend starts successfully
- [ ] Can fetch listings from `/api/sales`
- [ ] Can view individual listing at `/sales/:id`
- [ ] Scraped listings show source information
- [ ] Owned listings show seller information
- [ ] Filtering works (city, state, featured)
- [ ] Pagination works

## Test Data Examples

### Sample Owned Listing
```json
{
  "id": 1,
  "listing_type": "owned",
  "title": "Estate Sale - Vintage Furniture",
  "description": "Beautiful mid-century modern pieces",
  "address_line1": "123 Main Street",
  "city": "Portland",
  "state": "OR",
  "zip_code": "97201",
  "start_date": "2025-11-01T09:00:00Z",
  "end_date": "2025-11-03T17:00:00Z",
  "sale_type": "estate_sale",
  "status": "active",
  "is_scraped": false,
  "featured": true,
  "view_count": 42
}
```

### Sample Scraped Listing
```json
{
  "id": "external-123",
  "listing_type": "external",
  "external_id": "estatesales-net-12345",
  "external_source": "EstateSales.net",
  "external_url": "https://estatesales.net/OR/Portland/12345",
  "title": "Moving Sale - Everything Must Go",
  "address": "456 Oak Avenue",
  "city": "Seattle",
  "state": "WA",
  "zip_code": "98101",
  "start_date": "2025-11-05T08:00:00Z",
  "end_date": "2025-11-06T16:00:00Z",
  "is_scraped": true,
  "source": {
    "name": "EstateSales.net",
    "url": "https://estatesales.net/WA/Seattle/12345"
  }
}
```

## Continuous Integration

### GitHub Actions Example
```yaml
name: Test

on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: testdb
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable
        run: |
          cd backend
          go test -v ./...

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: |
          cd frontend
          npm install
      - name: Run type tests
        run: |
          cd frontend
          npm test -- types.test.ts
```

## Troubleshooting

### Backend Tests Failing

**Problem**: `DATABASE_URL not set`
**Solution**: Set the environment variable:
```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
```

**Problem**: `Failed to connect to database`
**Solution**: Ensure PostgreSQL is running and accessible:
```bash
psql -h localhost -U user -d dbname
```

**Problem**: `undefined: l`
**Solution**: Variable naming issue in service layer - check listingService.go

### Frontend Tests Failing

**Problem**: `Cannot find module 'jest'`
**Solution**: Install dependencies:
```bash
cd frontend
npm install
```

**Problem**: API tests timeout
**Solution**: Ensure backend is running on port 8080:
```bash
cd backend
go run cmd/api/main.go
```

**Problem**: Type errors in tests
**Solution**: Check that Listing interface matches backend response structure

## Summary

After refactoring from Sale → Listing:

✅ **Backend Tests Created**:
- 8 unit tests (listing_test.go)
- 4 integration test suites (listing_integration_test.go)
- 6 scraper integration tests (scraper_integration_test.go)

✅ **Frontend Tests Created**:
- 6 type validation tests (types.test.ts)
- 15+ API integration tests (api.test.ts)

✅ **Test Coverage**:
- Listing CRUD operations
- Image management
- Filtering and pagination
- External listing conversion
- Type safety verification
- API endpoint validation

Run all tests to verify the refactoring is successful!
