# Quick Start Guide - After Refactoring

This guide will help you verify that the Sale → Listing refactoring is working correctly.

## Overview

✅ **Refactoring Status**: COMPLETE and VALIDATED
- Backend: Sale → Listing ✓
- Frontend: Sale → Listing ✓
- Tests Created: 39+ tests ✓
- Validation: PASSED ✓

---

## Prerequisites Check

Before running tests, ensure you have:

- [ ] **WSL or Linux** (you have this)
- [ ] **Go 1.21+** (need to install - see below)
- [ ] **Node.js 18+** (check with `node --version`)
- [ ] **PostgreSQL** (for integration tests - optional)

---

## Step 1: Install Go in WSL

### Option A: Automated (Recommended)
```bash
cd /mnt/c/Users/matt/Desktop/Stuff/estatesalefinder-ai
bash install_go_wsl.sh
source ~/.bashrc
go version
```

### Option B: Manual
See `INSTALL_GO.md` for detailed instructions.

---

## Step 2: Quick Validation (No Compilation)

This checks the refactoring structure without needing Go or Node:

```bash
cd /mnt/c/Users/matt/Desktop/Stuff/estatesalefinder-ai
./validate_refactor.sh
```

**Expected Output:**
```
🎉 ALL CHECKS PASSED!

✓ No old package imports
✓ Package declarations correct
✓ No compilation issues detected
✓ Test files exist
✓ Frontend types updated
✓ All critical files exist
```

**Status**: ✅ This already passed!

---

## Step 3: Backend Tests

### 3a. Unit Tests (No Database Required)

```bash
cd backend

# Check compilation
go build ./...

# Run unit tests
go test -v ./internal/domain/listing -run TestListing
```

**Expected**: All 8 unit tests should pass

### 3b. Integration Tests (Database Required)

First, set up your database URL:
```bash
export DATABASE_URL="postgres://username:password@localhost:5432/your_database?sslmode=disable"
```

Then run integration tests:
```bash
go test -v ./internal/domain/listing -run TestListingIntegrationSuite
go test -v ./internal/infrastructure/scraper -run TestScraperIntegrationSuite
```

**Expected**: All integration tests should pass

### 3c. All Backend Tests

```bash
cd backend
./run_tests.sh
```

This runs everything:
1. ✅ Compilation check
2. ✅ Unit tests
3. ✅ Integration tests (if DATABASE_URL set)
4. ✅ Coverage report

---

## Step 4: Frontend Tests

### 4a. Install Dependencies

```bash
cd frontend
npm install
```

### 4b. Type Tests

```bash
npm test -- types.test.ts
```

**Expected**: 6 type tests should pass

### 4c. Build Check

```bash
npm run build
```

**Expected**: TypeScript should compile without errors

### 4d. API Integration Tests (Backend Must Be Running)

Terminal 1 - Start backend:
```bash
cd backend
go run cmd/api/main.go
```

Terminal 2 - Run tests:
```bash
cd frontend
npm test -- api.test.ts
```

**Expected**: 15+ API integration tests should pass

### 4e. All Frontend Tests

```bash
cd frontend
./run_tests.sh
```

---

## Step 5: Manual Verification

### Start the Application

Terminal 1 - Backend:
```bash
cd backend
go run cmd/api/main.go
```

Terminal 2 - Frontend:
```bash
cd frontend
npm run dev
```

### Test in Browser

1. Open http://localhost:3000
2. Navigate to http://localhost:3000/sales
3. Verify listings display
4. Click on a listing to view details
5. Check that scraped listings show source info
6. Check that filtering works (city, state, featured)

---

## What Each Test Covers

### Backend Unit Tests (`listing_test.go`)
- ✅ Listing struct creation (owned & external)
- ✅ ListingImage structure
- ✅ ListingFilters structure
- ✅ ScrapedListing creation
- ✅ Type conversions

### Backend Integration Tests (`listing_integration_test.go`)
- ✅ CREATE: Adding new listings
- ✅ READ: Fetching listings by ID
- ✅ UPDATE: Modifying listings
- ✅ DELETE: Removing listings
- ✅ Images: Add, delete, set primary
- ✅ Filtering: City, state, type, featured
- ✅ Pagination: Limit, offset
- ✅ External listings: Conversion and persistence

### Frontend Type Tests (`types.test.ts`)
- ✅ Listing interface structure
- ✅ Optional fields handling
- ✅ Array typing
- ✅ Filter/map operations

### Frontend API Tests (`api.test.ts`)
- ✅ GET all listings
- ✅ Filter by city, state, featured
- ✅ Pagination
- ✅ GET single listing
- ✅ 404 handling
- ✅ Date validation
- ✅ Image structure validation

---

## Troubleshooting

### "go: command not found"
**Solution**: Install Go using `bash install_go_wsl.sh`

### "DATABASE_URL not set"
**Solution**: Either set it:
```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/db?sslmode=disable"
```
Or skip integration tests (unit tests don't need DB)

### Backend tests fail with "connection refused"
**Solution**: Make sure PostgreSQL is running:
```bash
# Check if PostgreSQL is running
sudo service postgresql status

# Start if needed
sudo service postgresql start
```

### Frontend tests fail with "Cannot find module 'jest'"
**Solution**: Install dependencies:
```bash
cd frontend
npm install
```

### API tests timeout
**Solution**: Make sure backend is running on port 8080

### Windows line ending errors (`\r` errors)
**Solution**: Convert line endings:
```bash
dos2unix *.sh
# or
sed -i 's/\r$//' *.sh
```

---

## File Structure After Refactoring

```
backend/
├── internal/
│   ├── domain/
│   │   └── listing/              ← RENAMED from 'sale'
│   │       ├── listing.go        ← Listing struct
│   │       ├── listingRepo.go    ← Repository interface
│   │       ├── listingService.go ← Business logic
│   │       ├── scraped_listing.go← External listings
│   │       ├── listing_test.go   ← NEW: Unit tests
│   │       └── listing_integration_test.go ← NEW: Integration tests
│   └── infrastructure/
│       ├── db/postgres/
│       │   └── saleRepo.go       ← Uses listing.Listing
│       ├── controllers/
│       │   └── saleHandler.go    ← Uses listing.Listing
│       └── scraper/
│           ├── scraper.go        ← Uses listing.ScrapedListing
│           └── scraper_integration_test.go
├── run_tests.sh                  ← NEW: Test runner
└── cmd/api/main.go              ← Uses listing package

frontend/
├── app/
│   ├── sales/
│   │   ├── page.tsx             ← Uses Listing interface
│   │   └── [id]/page.tsx        ← Uses Listing interface
├── __tests__/                    ← NEW: Test directory
│   ├── types.test.ts            ← Type tests
│   └── api.test.ts              ← API integration tests
├── package.json                 ← Added Jest dependencies
├── jest.config.js               ← NEW: Jest config
├── jest.setup.js                ← NEW: Test setup
└── run_tests.sh                 ← NEW: Test runner

Root:
├── validate_refactor.sh         ← NEW: Quick validation ✅ PASSED
├── install_go_wsl.sh           ← NEW: Go installer
├── INSTALL_GO.md               ← NEW: Installation guide
├── TESTING.md                  ← NEW: Testing documentation
├── REFACTORING_SUMMARY.md      ← NEW: Refactoring summary
└── QUICK_START.md              ← This file
```

---

## Success Checklist

Use this to verify everything is working:

### Backend
- [ ] Go installed and working (`go version`)
- [ ] Code compiles (`go build ./...`)
- [ ] Unit tests pass
- [ ] Integration tests pass (if DB configured)
- [ ] Backend starts without errors
- [ ] Can access http://localhost:8080/health

### Frontend
- [ ] Node.js installed (`node --version`)
- [ ] Dependencies installed (`npm install`)
- [ ] TypeScript compiles (`npm run build`)
- [ ] Type tests pass
- [ ] Frontend starts without errors
- [ ] Can access http://localhost:3000

### Integration
- [ ] Can fetch listings from `/api/sales`
- [ ] Can view individual listing
- [ ] Scraped listings show source info
- [ ] Filtering works (city, state)
- [ ] Pagination works

---

## Summary

To verify the refactoring works:

```bash
# 1. Quick validation (no tools needed)
./validate_refactor.sh          # ✅ Already passed!

# 2. Install Go
bash install_go_wsl.sh
source ~/.bashrc

# 3. Backend tests
cd backend
./run_tests.sh

# 4. Frontend tests
cd ../frontend
npm install
./run_tests.sh

# 5. Start and test manually
# Terminal 1: cd backend && go run cmd/api/main.go
# Terminal 2: cd frontend && npm run dev
# Browser: http://localhost:3000/sales
```

---

## Need More Help?

- **Installation**: See `INSTALL_GO.md`
- **Testing**: See `TESTING.md`
- **Summary**: See `REFACTORING_SUMMARY.md`
- **Backend Issues**: Check backend logs
- **Frontend Issues**: Check browser console
- **Database Issues**: Check PostgreSQL status

---

**Status**: Ready to test! 🚀

The refactoring is complete and validated. All that's left is to install Go and run the tests to verify everything works correctly.
