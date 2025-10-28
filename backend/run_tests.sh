#!/bin/bash

# Test runner script for EstateSaleFinder.ai backend
# This script runs all tests and provides clear feedback

set -e  # Exit on error

echo "======================================"
echo "EstateSaleFinder.ai - Backend Tests"
echo "======================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ ERROR: Go is not installed"
    echo "   Please install Go 1.21 or later"
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

echo "✓ Go version: $(go version)"
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ ERROR: Not in backend directory"
    echo "   Please run this script from the backend directory"
    exit 1
fi

echo "======================================"
echo "Step 1: Check if code compiles"
echo "======================================"
if go build ./...; then
    echo "✓ Code compiles successfully"
else
    echo "❌ Compilation failed"
    exit 1
fi
echo ""

echo "======================================"
echo "Step 2: Run unit tests (no DB needed)"
echo "======================================"
echo ""
if go test -v ./internal/domain/listing -run TestListing; then
    echo ""
    echo "✓ Unit tests passed"
else
    echo ""
    echo "❌ Unit tests failed"
    exit 1
fi
echo ""

echo "======================================"
echo "Step 3: Check for DATABASE_URL"
echo "======================================"
if [ -z "$DATABASE_URL" ]; then
    echo "⚠ WARNING: DATABASE_URL not set"
    echo "   Integration tests will be skipped"
    echo ""
    echo "   To run integration tests, set DATABASE_URL:"
    echo "   export DATABASE_URL=\"postgres://user:pass@localhost:5432/dbname?sslmode=disable\""
    echo ""
    echo "======================================"
    echo "Test Summary"
    echo "======================================"
    echo "✓ Compilation: PASSED"
    echo "✓ Unit Tests: PASSED"
    echo "⚠ Integration Tests: SKIPPED (no DATABASE_URL)"
    exit 0
fi

echo "✓ DATABASE_URL is set"
echo ""

echo "======================================"
echo "Step 4: Run integration tests"
echo "======================================"
echo ""

echo "--- Listing Integration Tests ---"
if go test -v ./internal/domain/listing -run TestListingIntegrationSuite; then
    echo ""
    echo "✓ Listing integration tests passed"
else
    echo ""
    echo "❌ Listing integration tests failed"
    exit 1
fi
echo ""

echo "--- Scraper Integration Tests ---"
if go test -v ./internal/infrastructure/scraper -run TestScraperIntegrationSuite; then
    echo ""
    echo "✓ Scraper integration tests passed"
else
    echo ""
    echo "❌ Scraper integration tests failed"
    exit 1
fi
echo ""

echo "======================================"
echo "Step 5: Run all tests with coverage"
echo "======================================"
echo ""
if go test -cover ./...; then
    echo ""
    echo "✓ All tests passed with coverage"
else
    echo ""
    echo "❌ Some tests failed"
    exit 1
fi
echo ""

echo "======================================"
echo "Test Summary - ALL PASSED! ✓"
echo "======================================"
echo "✓ Compilation: PASSED"
echo "✓ Unit Tests: PASSED"
echo "✓ Listing Integration Tests: PASSED"
echo "✓ Scraper Integration Tests: PASSED"
echo "✓ Coverage Report: PASSED"
echo ""
echo "🎉 All tests completed successfully!"
echo ""
echo "Refactoring verification: Sale → Listing"
echo "✓ All types renamed correctly"
echo "✓ All imports updated"
echo "✓ CRUD operations working"
echo "✓ Conversions working"
echo "✓ Filtering and pagination working"
echo ""
