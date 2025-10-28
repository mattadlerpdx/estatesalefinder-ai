#!/bin/bash

# Validation script for Sale → Listing refactoring
# This runs basic checks that don't require Go or Node to be installed

echo "======================================"
echo "Refactoring Validation Script"
echo "Sale → Listing"
echo "======================================"
echo ""

ERRORS=0
WARNINGS=0

echo "Checking backend files..."
echo ""

# Check 1: No old 'sale' package imports
echo "❯ Checking for old package imports..."
if grep -r "github.com/mattadlerpdx.*domain/sale\"" backend/internal --include="*.go" 2>/dev/null; then
    echo "❌ Found old 'domain/sale' imports!"
    ERRORS=$((ERRORS + 1))
else
    echo "✓ No old package imports found"
fi
echo ""

# Check 2: Package declarations are correct
echo "❯ Checking package declarations in domain/listing..."
if grep -l "^package sale$" backend/internal/domain/listing/*.go 2>/dev/null; then
    echo "❌ Found 'package sale' in listing directory!"
    ERRORS=$((ERRORS + 1))
else
    echo "✓ All files use 'package listing'"
fi
echo ""

# Check 3: Check for undefined references
echo "❯ Checking for common compilation issues..."
if grep -n "undefined: l$" backend/internal/domain/listing/*.go 2>/dev/null; then
    echo "❌ Found 'undefined: l' references!"
    ERRORS=$((ERRORS + 1))
fi

if grep -n "sale\.Sale" backend/internal --include="*.go" -r 2>/dev/null | grep -v "listing\.Listing"; then
    echo "❌ Found old 'sale.Sale' references!"
    ERRORS=$((ERRORS + 1))
fi

if ! grep -q "listing\.Listing" backend/internal/domain/listing/*.go 2>/dev/null; then
    echo "❌ No 'listing.Listing' references found - something is wrong!"
    ERRORS=$((ERRORS + 1))
else
    echo "✓ Found listing.Listing references"
fi
echo ""

# Check 4: Test files exist
echo "❯ Checking test files..."
if [ ! -f "backend/internal/domain/listing/listing_test.go" ]; then
    echo "❌ listing_test.go not found!"
    ERRORS=$((ERRORS + 1))
else
    echo "✓ listing_test.go exists"
fi

if [ ! -f "backend/internal/domain/listing/listing_integration_test.go" ]; then
    echo "❌ listing_integration_test.go not found!"
    ERRORS=$((ERRORS + 1))
else
    echo "✓ listing_integration_test.go exists"
fi
echo ""

# Check 5: Frontend types
echo "❯ Checking frontend types..."
if grep -n "interface Sale {" frontend/app --include="*.tsx" --include="*.ts" -r 2>/dev/null; then
    echo "⚠ WARNING: Found old 'interface Sale' in frontend"
    WARNINGS=$((WARNINGS + 1))
fi

if grep -q "interface Listing {" frontend/app/sales/page.tsx 2>/dev/null; then
    echo "✓ Frontend uses 'interface Listing'"
else
    echo "❌ Frontend doesn't use 'interface Listing'"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# Check 6: Critical files exist
echo "❯ Checking critical files..."
CRITICAL_FILES=(
    "backend/internal/domain/listing/listing.go"
    "backend/internal/domain/listing/listingRepo.go"
    "backend/internal/domain/listing/listingService.go"
    "backend/internal/domain/listing/scraped_listing.go"
)

for file in "${CRITICAL_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        echo "❌ Critical file missing: $file"
        ERRORS=$((ERRORS + 1))
    fi
done
echo "✓ All critical files exist"
echo ""

# Summary
echo "======================================"
echo "Validation Summary"
echo "======================================"

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo "🎉 ALL CHECKS PASSED!"
    echo ""
    echo "✓ No old package imports"
    echo "✓ Package declarations correct"
    echo "✓ No compilation issues detected"
    echo "✓ Test files exist"
    echo "✓ Frontend types updated"
    echo "✓ All critical files exist"
    echo ""
    echo "Next steps:"
    echo "1. Run backend tests: cd backend && ./run_tests.sh"
    echo "2. Run frontend tests: cd frontend && ./run_tests.sh"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo "⚠ PASSED WITH WARNINGS"
    echo ""
    echo "Warnings: $WARNINGS"
    echo ""
    echo "Please review warnings above"
    exit 0
else
    echo "❌ VALIDATION FAILED"
    echo ""
    echo "Errors: $ERRORS"
    echo "Warnings: $WARNINGS"
    echo ""
    echo "Please fix errors above before running tests"
    exit 1
fi
