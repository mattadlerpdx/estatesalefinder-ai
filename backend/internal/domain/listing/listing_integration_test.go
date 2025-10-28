package listing_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/listing"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/db/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

// ListingIntegrationTestSuite tests the listing domain with real database
type ListingIntegrationTestSuite struct {
	suite.Suite
	db      *sql.DB
	repo    *postgres.ListingRepository
	service *listing.Service
}

// SetupSuite runs once before all tests
func (suite *ListingIntegrationTestSuite) SetupSuite() {
	// Get database connection from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		suite.T().Skip("DATABASE_URL not set, skipping integration tests")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	require.NoError(suite.T(), err, "Failed to connect to database")
	require.NoError(suite.T(), db.Ping(), "Failed to ping database")

	suite.db = db
	suite.repo = postgres.NewListingRepository(db)
	suite.service = listing.NewService(suite.repo)
}

// TearDownSuite runs once after all tests
func (suite *ListingIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest runs before each test
func (suite *ListingIntegrationTestSuite) SetupTest() {
	suite.cleanupTestData()
}

// TearDownTest runs after each test
func (suite *ListingIntegrationTestSuite) TearDownTest() {
	suite.cleanupTestData()
}

// cleanupTestData removes test data from database
func (suite *ListingIntegrationTestSuite) cleanupTestData() {
	// Delete test listings (with title starting with "Test:")
	_, err := suite.db.Exec("DELETE FROM listings WHERE title LIKE 'Test:%'")
	if err != nil {
		suite.T().Logf("Warning: Failed to cleanup database: %v", err)
	}

	// Small delay to ensure cleanup completes
	time.Sleep(100 * time.Millisecond)
}

// TestListingCRUD tests Create, Read, Update, Delete operations
func (suite *ListingIntegrationTestSuite) TestListingCRUD() {
	suite.T().Log("=== Test: Listing CRUD Operations ===")

	// Create a test listing
	now := time.Now()
	sellerID := 10 // Test seller created in setup
	testListing := listing.Listing{
		ListingType:   "owned",
		SellerID:      &sellerID,
		Title:         "Test: Estate Sale - Antiques & Collectibles",
		Description:   "Beautiful collection of vintage items, furniture, and collectibles",
		AddressLine1:  "123 Main Street",
		City:          "Portland",
		State:         "OR",
		ZipCode:       "97201",
		StartDate:     now.Add(7 * 24 * time.Hour),
		EndDate:       now.Add(9 * 24 * time.Hour),
		EventType:      "estate_sale",
		Status:        "draft",
		ListingTier:   "basic",
		PaymentStatus: "unpaid",
	}

	// Test CREATE
	err := suite.service.CreateListing(&testListing)
	require.NoError(suite.T(), err, "CreateListing should succeed")
	assert.Greater(suite.T(), testListing.ID, 0, "Listing should have ID assigned")
	suite.T().Logf("✓ Created listing with ID: %d", testListing.ID)

	createdID := testListing.ID

	// Test READ (GetByID)
	retrieved, err := suite.service.GetListingByID(createdID)
	require.NoError(suite.T(), err, "GetListingByID should succeed")
	assert.Equal(suite.T(), testListing.Title, retrieved.Title, "Title should match")
	assert.Equal(suite.T(), testListing.City, retrieved.City, "City should match")
	assert.Equal(suite.T(), "draft", retrieved.Status, "Status should be draft")
	assert.GreaterOrEqual(suite.T(), retrieved.ViewCount, 1, "View count should be incremented")
	suite.T().Logf("✓ Retrieved listing: %s (views: %d)", retrieved.Title, retrieved.ViewCount)

	// Test UPDATE
	retrieved.Title = "Test: UPDATED - Estate Sale"
	retrieved.Status = "active"
	retrieved.Description = "Updated description with more details"
	err = suite.service.UpdateListing(retrieved)
	require.NoError(suite.T(), err, "UpdateListing should succeed")
	suite.T().Logf("✓ Updated listing")

	// Verify update
	updated, err := suite.service.GetListingByID(createdID)
	require.NoError(suite.T(), err, "GetListingByID after update should succeed")
	assert.Equal(suite.T(), "Test: UPDATED - Estate Sale", updated.Title, "Title should be updated")
	assert.Equal(suite.T(), "active", updated.Status, "Status should be updated")
	suite.T().Logf("✓ Verified update: %s", updated.Title)

	// Test DELETE
	err = suite.service.DeleteListing(createdID)
	require.NoError(suite.T(), err, "DeleteListing should succeed")
	suite.T().Logf("✓ Deleted listing")

	// Verify deletion
	_, err = suite.service.GetListingByID(createdID)
	assert.Error(suite.T(), err, "GetListingByID after delete should fail")
	suite.T().Logf("✓ Verified deletion")
}

// TestListingWithImages tests image operations
func (suite *ListingIntegrationTestSuite) TestListingWithImages() {
	suite.T().Log("=== Test: Listing with Images ===")

	// Create a listing
	now := time.Now()
	sellerID := 10 // Test seller created in setup
	testListing := listing.Listing{
		ListingType:  "owned",
		SellerID:     &sellerID,
		Title:        "Test: Moving Sale with Photos",
		Description:  "Everything must go!",
		AddressLine1: "456 Oak Avenue",
		City:         "Portland",
		State:        "OR",
		ZipCode:      "97202",
		StartDate:    now.Add(5 * 24 * time.Hour),
		EndDate:      now.Add(7 * 24 * time.Hour),
		EventType:     "moving_sale",
		Status:       "active",
	}

	err := suite.service.CreateListing(&testListing)
	require.NoError(suite.T(), err)
	suite.T().Logf("✓ Created listing: %s (ID: %d)", testListing.Title, testListing.ID)

	// Add images
	image1 := listing.ListingImage{
		ListingID:    testListing.ID,
		ImageURL:     "https://example.com/images/photo1.jpg",
		IsPrimary:    false,
		DisplayOrder: 1,
	}
	err = suite.service.AddListingImage(&image1)
	require.NoError(suite.T(), err, "AddListingImage should succeed")
	suite.T().Logf("✓ Added image 1")

	image2 := listing.ListingImage{
		ListingID:    testListing.ID,
		ImageURL:     "https://example.com/images/photo2.jpg",
		IsPrimary:    false,
		DisplayOrder: 2,
	}
	err = suite.service.AddListingImage(&image2)
	require.NoError(suite.T(), err)
	suite.T().Logf("✓ Added image 2")

	// Set primary image
	err = suite.service.SetPrimaryImage(image1.ID, testListing.ID)
	require.NoError(suite.T(), err, "SetPrimaryImage should succeed")
	suite.T().Logf("✓ Set image 1 as primary")

	// Retrieve listing with images
	retrieved, err := suite.service.GetListingByID(testListing.ID)
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), retrieved.Images, 2, "Should have 2 images")

	// Verify primary image
	var primaryCount int
	for _, img := range retrieved.Images {
		if img.IsPrimary {
			primaryCount++
			assert.Equal(suite.T(), image1.ImageURL, img.ImageURL, "Primary image should be image1")
		}
	}
	assert.Equal(suite.T(), 1, primaryCount, "Should have exactly 1 primary image")
	suite.T().Logf("✓ Verified images (primary: %s)", image1.ImageURL)

	// Delete image
	err = suite.service.DeleteListingImage(image2.ID)
	require.NoError(suite.T(), err, "DeleteListingImage should succeed")
	suite.T().Logf("✓ Deleted image 2")

	// Verify deletion
	retrieved, err = suite.service.GetListingByID(testListing.ID)
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), retrieved.Images, 1, "Should have 1 image after deletion")
	suite.T().Logf("✓ Verified image deletion")

	// Cleanup
	suite.service.DeleteListing(testListing.ID)
}

// TestListingFilters tests filtering and pagination
func (suite *ListingIntegrationTestSuite) TestListingFilters() {
	suite.T().Log("=== Test: Listing Filters and Pagination ===")

	now := time.Now()
	sellerID := 10 // Test seller created in setup

	// Create multiple test listings
	listings := []listing.Listing{
		{
			ListingType:  "owned",
			SellerID:     &sellerID,
			Title:        "Test: Portland Estate Sale",
			AddressLine1: "100 SW Broadway",
			City:         "Portland",
			State:        "OR",
			ZipCode:      "97205",
			StartDate:    now.Add(1 * 24 * time.Hour),
			EndDate:      now.Add(2 * 24 * time.Hour),
			EventType:     "estate_sale",
			Status:       "active",
			Featured:     true,
		},
		{
			ListingType:  "owned",
			SellerID:     &sellerID,
			Title:        "Test: Portland Moving Sale",
			AddressLine1: "200 NW Couch",
			City:         "Portland",
			State:        "OR",
			ZipCode:      "97209",
			StartDate:    now.Add(3 * 24 * time.Hour),
			EndDate:      now.Add(4 * 24 * time.Hour),
			EventType:     "moving_sale",
			Status:       "active",
			Featured:     false,
		},
		{
			ListingType:  "owned",
			SellerID:     &sellerID,
			Title:        "Test: Seattle Estate Sale",
			AddressLine1: "300 Pike Street",
			City:         "Seattle",
			State:        "WA",
			ZipCode:      "98101",
			StartDate:    now.Add(5 * 24 * time.Hour),
			EndDate:      now.Add(6 * 24 * time.Hour),
			EventType:     "estate_sale",
			Status:       "active",
			Featured:     false,
		},
	}

	for i := range listings {
		err := suite.service.CreateListing(&listings[i])
		require.NoError(suite.T(), err)
		suite.T().Logf("✓ Created: %s", listings[i].Title)
	}

	// Test 1: Filter by city
	filters := listing.ListingFilters{
		City:  "Portland",
		Limit: 10,
	}
	results, err := suite.service.GetAllListings(filters)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(results), 2, "Should find at least 2 Portland listings")
	for _, l := range results {
		if l.Title[:5] == "Test:" {
			assert.Equal(suite.T(), "Portland", l.City, "All results should be from Portland")
		}
	}
	suite.T().Logf("✓ City filter: Found %d Portland listings", len(results))

	// Test 2: Filter by state
	filters = listing.ListingFilters{
		State: "OR",
		Limit: 10,
	}
	results, err = suite.service.GetAllListings(filters)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(results), 2, "Should find at least 2 OR listings")
	suite.T().Logf("✓ State filter: Found %d Oregon listings", len(results))

	// Test 3: Filter by sale type
	filters = listing.ListingFilters{
		EventType: "estate_sale",
		Limit:    10,
	}
	results, err = suite.service.GetAllListings(filters)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(results), 2, "Should find at least 2 estate sales")
	suite.T().Logf("✓ Sale type filter: Found %d estate sales", len(results))

	// Test 4: Filter by featured
	featured := true
	filters = listing.ListingFilters{
		Featured: &featured,
		Limit:    10,
	}
	results, err = suite.service.GetAllListings(filters)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(results), 1, "Should find at least 1 featured listing")
	suite.T().Logf("✓ Featured filter: Found %d featured listings", len(results))

	// Test 5: Pagination
	filters = listing.ListingFilters{
		Limit:  1,
		Offset: 0,
	}
	page1, err := suite.service.GetAllListings(filters)
	require.NoError(suite.T(), err)
	assert.LessOrEqual(suite.T(), len(page1), 1, "Should return at most 1 result")
	suite.T().Logf("✓ Pagination: Page 1 has %d results", len(page1))

	filters.Offset = 1
	page2, err := suite.service.GetAllListings(filters)
	require.NoError(suite.T(), err)
	assert.LessOrEqual(suite.T(), len(page2), 1, "Should return at most 1 result")
	suite.T().Logf("✓ Pagination: Page 2 has %d results", len(page2))

	// Cleanup
	for _, l := range listings {
		suite.service.DeleteListing(l.ID)
	}
}

// TestExternalListingConversion tests ScrapedListing conversion
func (suite *ListingIntegrationTestSuite) TestExternalListingConversion() {
	suite.T().Log("=== Test: External Listing Conversion ===")

	now := time.Now()

	// Create a ScrapedListing
	scraped := listing.ScrapedListing{
		ExternalID:   "test-external-123",
		Title:        "Test: External Estate Sale",
		Description:  "Scraped from external source",
		Address:      "789 Test Lane",
		City:         "Portland",
		State:        "OR",
		ZipCode:      "97210",
		Latitude:     45.5231,
		Longitude:    -122.6765,
		StartDate:    now.Add(10 * 24 * time.Hour),
		EndDate:      now.Add(12 * 24 * time.Hour),
		ThumbnailURL: "https://example.com/thumb.jpg",
		ImageURLs:    []string{"https://example.com/img1.jpg", "https://example.com/img2.jpg"},
		SourceName:   "TestSource",
		SourceURL:    "https://example.com/sale/123",
		ScrapedAt:    now,
		CachedAt:     now,
	}

	// Convert to Listing
	convertedListing := scraped.ToSale()
	assert.Equal(suite.T(), "external", convertedListing.ListingType, "ListingType should be external")
	assert.Equal(suite.T(), "test-external-123", *convertedListing.ExternalID, "ExternalID should match")
	assert.Equal(suite.T(), "TestSource", *convertedListing.ExternalSource, "Source should match")
	assert.Equal(suite.T(), scraped.Title, convertedListing.Title, "Title should match")
	assert.Equal(suite.T(), scraped.City, convertedListing.City, "City should match")
	suite.T().Logf("✓ Converted ScrapedListing to Listing")

	// Persist to database
	err := suite.repo.UpsertExternalSale(&convertedListing)
	require.NoError(suite.T(), err, "UpsertExternalSale should succeed")
	assert.Greater(suite.T(), convertedListing.ID, 0, "Should have ID assigned")
	suite.T().Logf("✓ Persisted external listing (ID: %d)", convertedListing.ID)

	// Convert Listing back to ScrapedListing
	retrieved, err := suite.repo.GetByID(convertedListing.ID)
	require.NoError(suite.T(), err)

	convertedBack := retrieved.ToScrapedListing()
	assert.Equal(suite.T(), scraped.ExternalID, convertedBack.ExternalID, "ExternalID should match")
	assert.Equal(suite.T(), scraped.SourceName, convertedBack.SourceName, "Source name should match")
	assert.Equal(suite.T(), scraped.Title, convertedBack.Title, "Title should match")
	suite.T().Logf("✓ Converted Listing back to ScrapedListing")

	// Test ToAggregatedSale
	aggregated := convertedBack.ToAggregatedSale()
	assert.True(suite.T(), aggregated.IsScraped, "Should be marked as scraped")
	assert.NotNil(suite.T(), aggregated.Source, "Should have source info")
	assert.Equal(suite.T(), scraped.SourceName, aggregated.Source.Name, "Source name should match")
	suite.T().Logf("✓ Converted to AggregatedListing (scraped: %v)", aggregated.IsScraped)

	// Cleanup
	suite.db.Exec("DELETE FROM listings WHERE external_id = 'test-external-123'")
}

// TestListingIntegrationSuite runs the integration test suite
func TestListingIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ListingIntegrationTestSuite))
}
