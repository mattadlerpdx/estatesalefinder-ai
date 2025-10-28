package scraper

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/listing"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/cache"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/db/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

// ScraperIntegrationTestSuite is the test suite for scraper integration tests
type ScraperIntegrationTestSuite struct {
	suite.Suite
	db             *sql.DB
	repo           *postgres.SaleRepository
	redisClient    *cache.RedisClient
	scraperService *ScraperService
}

// SetupSuite runs once before all tests
func (suite *ScraperIntegrationTestSuite) SetupSuite() {
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
	suite.repo = postgres.NewSaleRepository(db)
	suite.redisClient = cache.NewRedisClient()
	suite.scraperService = NewScraperService(suite.redisClient, suite.repo)
}

// TearDownSuite runs once after all tests
func (suite *ScraperIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
	if suite.redisClient != nil {
		suite.redisClient.Close()
	}
}

// SetupTest runs before each test
func (suite *ScraperIntegrationTestSuite) SetupTest() {
	// Clean up test data
	suite.cleanupTestData()
}

// TearDownTest runs after each test
func (suite *ScraperIntegrationTestSuite) TearDownTest() {
	suite.cleanupTestData()
}

// cleanupTestData removes test data from database and cache
func (suite *ScraperIntegrationTestSuite) cleanupTestData() {
	// Delete all external sales from test
	_, err := suite.db.Exec("DELETE FROM listings WHERE listing_type = 'external'")
	if err != nil {
		suite.T().Logf("Warning: Failed to cleanup database: %v", err)
	}

	// Clear Redis cache for Portland
	if suite.redisClient.IsEnabled() {
		err = suite.scraperService.InvalidateCache("Portland", "OR")
		if err != nil {
			suite.T().Logf("Warning: Failed to clear cache: %v", err)
		}
	}

	// Small delay to ensure cleanup completes
	time.Sleep(100 * time.Millisecond)
}

// TestHybridStorage_InitialScrape tests the initial scrape flow
func (suite *ScraperIntegrationTestSuite) TestHybridStorage_InitialScrape() {
	suite.T().Log("=== Test 1: Initial Scrape (Redis empty, DB empty) ===")

	// Verify database is empty
	count := suite.countExternalSales()
	assert.Equal(suite.T(), 0, count, "Database should be empty before test")

	// Make request - should trigger scrape
	sales, err := suite.scraperService.GetSalesByLocation("Portland", "OR")
	require.NoError(suite.T(), err, "Initial scrape should succeed")
	assert.Greater(suite.T(), len(sales), 0, "Should return sales from scrape")

	suite.T().Logf("✓ Scraped %d sales", len(sales))

	// Verify data was persisted to PostgreSQL
	count = suite.countExternalSales()
	assert.Equal(suite.T(), len(sales), count, "All scraped sales should be persisted to database")

	// Verify last_scraped_at is set and recent
	lastScrape := suite.getLastScrapedTime()
	assert.NotNil(suite.T(), lastScrape, "last_scraped_at should be set")
	assert.WithinDuration(suite.T(), time.Now(), *lastScrape, 10*time.Second, "last_scraped_at should be recent")

	suite.T().Logf("✓ Persisted %d sales to PostgreSQL", count)
}

// TestHybridStorage_RedisCacheHit tests Redis cache hit
func (suite *ScraperIntegrationTestSuite) TestHybridStorage_RedisCacheHit() {
	suite.T().Log("=== Test 2: Redis Cache Hit ===")

	// First request to populate cache
	sales1, err := suite.scraperService.GetSalesByLocation("Portland", "OR")
	require.NoError(suite.T(), err)
	initialCount := len(sales1)

	// Second request should hit cache (fast)
	start := time.Now()
	sales2, err := suite.scraperService.GetSalesByLocation("Portland", "OR")
	elapsed := time.Since(start)

	require.NoError(suite.T(), err, "Cached request should succeed")
	assert.Equal(suite.T(), initialCount, len(sales2), "Should return same number of sales")
	assert.Less(suite.T(), elapsed.Milliseconds(), int64(100), "Cache hit should be fast (<100ms)")

	suite.T().Logf("✓ Cache hit in %dms", elapsed.Milliseconds())
}

// TestHybridStorage_PostgreSQLFallback tests PostgreSQL fallback when Redis expires
func (suite *ScraperIntegrationTestSuite) TestHybridStorage_PostgreSQLFallback() {
	suite.T().Log("=== Test 3: PostgreSQL Fallback (Redis expired, DB fresh) ===")

	// First request to populate database
	sales1, err := suite.scraperService.GetSalesByLocation("Portland", "OR")
	require.NoError(suite.T(), err)
	initialCount := len(sales1)

	// Clear Redis cache (simulate expiration)
	err = suite.scraperService.InvalidateCache("Portland", "OR")
	require.NoError(suite.T(), err, "Should clear cache successfully")
	suite.T().Log("✓ Cleared Redis cache")

	// Second request should load from PostgreSQL (not re-scrape)
	start := time.Now()
	sales2, err := suite.scraperService.GetSalesByLocation("Portland", "OR")
	elapsed := time.Since(start)

	require.NoError(suite.T(), err, "PostgreSQL fallback should succeed")
	assert.Equal(suite.T(), initialCount, len(sales2), "Should return same sales from database")
	assert.Less(suite.T(), elapsed.Milliseconds(), int64(1000), "DB query should be faster than scraping (<1s)")

	suite.T().Logf("✓ Loaded from PostgreSQL in %dms", elapsed.Milliseconds())

	// Verify we didn't re-scrape (last_scraped_at should be old)
	lastScrape := suite.getLastScrapedTime()
	assert.NotNil(suite.T(), lastScrape)
	assert.Greater(suite.T(), time.Since(*lastScrape).Seconds(), 1.0, "Should use old scrape data")
}

// TestHybridStorage_SixHourRefresh tests re-scraping after 6 hours
func (suite *ScraperIntegrationTestSuite) TestHybridStorage_SixHourRefresh() {
	suite.T().Log("=== Test 4: 6-Hour Refresh (Stale data triggers re-scrape) ===")

	// Insert old external sale (>6 hours old)
	oldTime := time.Now().Add(-7 * time.Hour)
	suite.insertOldExternalSale(oldTime)

	// Clear Redis cache
	err := suite.scraperService.InvalidateCache("Portland", "OR")
	require.NoError(suite.T(), err)

	// Request should trigger re-scrape (data is stale)
	sales, err := suite.scraperService.GetSalesByLocation("Portland", "OR")
	require.NoError(suite.T(), err, "Re-scrape should succeed")
	assert.Greater(suite.T(), len(sales), 0, "Should return fresh sales")

	// Verify last_scraped_at was updated to now
	lastScrape := suite.getLastScrapedTime()
	assert.NotNil(suite.T(), lastScrape)
	assert.WithinDuration(suite.T(), time.Now(), *lastScrape, 10*time.Second, "Should have fresh scrape time")

	suite.T().Logf("✓ Re-scraped after 6-hour threshold")
}

// TestHybridStorage_ExternalSaleUpsert tests upserting external sales
func (suite *ScraperIntegrationTestSuite) TestHybridStorage_ExternalSaleUpsert() {
	suite.T().Log("=== Test 5: External Sale Upsert (Insert + Update) ===")

	now := time.Now()
	externalID := "test-sale-12345"

	// Create first version
	sale1 := listing.Listing{
		ListingType:    "external",
		ExternalID:     &externalID,
		ExternalSource: stringPtr("TestSource"),
		ExternalURL:    stringPtr("https://example.com/test"),
		LastScrapedAt:  &now,
		Title:          "Original Title",
		Description:    "Original Description",
		AddressLine1:   "123 Test St",
		City:           "Portland",
		State:          "OR",
		ZipCode:        "97201",
		StartDate:      now,
		EndDate:        now.Add(24 * time.Hour),
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Insert
	err := suite.repo.UpsertExternalSale(&sale1)
	require.NoError(suite.T(), err, "First upsert (insert) should succeed")
	assert.Greater(suite.T(), sale1.ID, 0, "Should assign ID")
	insertedID := sale1.ID

	suite.T().Logf("✓ Inserted sale with ID: %d", insertedID)

	// Update with same external_id
	sale2 := listing.Listing{
		ListingType:    "external",
		ExternalID:     &externalID,
		ExternalSource: stringPtr("TestSource"),
		ExternalURL:    stringPtr("https://example.com/test"),
		LastScrapedAt:  &now,
		Title:          "Updated Title",
		Description:    "Updated Description",
		AddressLine1:   "456 New St",
		City:           "Portland",
		State:          "OR",
		ZipCode:        "97202",
		StartDate:      now,
		EndDate:        now.Add(24 * time.Hour),
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Upsert (should update, not insert)
	err = suite.repo.UpsertExternalSale(&sale2)
	require.NoError(suite.T(), err, "Second upsert (update) should succeed")
	assert.Equal(suite.T(), insertedID, sale2.ID, "Should reuse same ID (update, not insert)")

	// Verify update worked
	retrieved, err := suite.repo.GetByID(insertedID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Title", retrieved.Title)
	assert.Equal(suite.T(), "456 New St", retrieved.AddressLine1)

	suite.T().Logf("✓ Updated sale (same external_id)")

	// Verify count (should still be 1)
	count := suite.countExternalSales()
	assert.Equal(suite.T(), 1, count, "Should have exactly 1 sale (upserted, not duplicated)")
}

// TestHybridStorage_GetExternalSalesByLocation tests location filtering
func (suite *ScraperIntegrationTestSuite) TestHybridStorage_GetExternalSalesByLocation() {
	suite.T().Log("=== Test 6: Get External Sales by Location ===")

	now := time.Now()

	// Insert sales in different locations
	portland1 := createTestExternalSale("portland-1", "Portland", "OR", now)
	portland2 := createTestExternalSale("portland-2", "Portland", "OR", now)
	seattle1 := createTestExternalSale("seattle-1", "Seattle", "WA", now)

	require.NoError(suite.T(), suite.repo.UpsertExternalSale(&portland1))
	require.NoError(suite.T(), suite.repo.UpsertExternalSale(&portland2))
	require.NoError(suite.T(), suite.repo.UpsertExternalSale(&seattle1))

	// Query Portland sales
	portlandSales, err := suite.repo.GetExternalSalesByLocation("Portland", "OR")
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(portlandSales), "Should return 2 Portland sales")

	// Query Seattle sales
	seattleSales, err := suite.repo.GetExternalSalesByLocation("Seattle", "WA")
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(seattleSales), "Should return 1 Seattle sale")

	// Query non-existent location
	nySales, err := suite.repo.GetExternalSalesByLocation("New York", "NY")
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, len(nySales), "Should return 0 sales for non-existent location")

	suite.T().Logf("✓ Location filtering works correctly")
}

// Helper functions

func (suite *ScraperIntegrationTestSuite) countExternalSales() int {
	var count int
	err := suite.db.QueryRow("SELECT COUNT(*) FROM listings WHERE listing_type = 'external'").Scan(&count)
	if err != nil {
		suite.T().Fatalf("Failed to count external sales: %v", err)
	}
	return count
}

func (suite *ScraperIntegrationTestSuite) getLastScrapedTime() *time.Time {
	var lastScraped *time.Time
	err := suite.db.QueryRow(`
		SELECT MAX(last_scraped_at)
		FROM listings
		WHERE listing_type = 'external' AND city = 'Portland' AND state = 'OR'
	`).Scan(&lastScraped)
	if err != nil && err != sql.ErrNoRows {
		suite.T().Fatalf("Failed to get last scraped time: %v", err)
	}
	return lastScraped
}

func (suite *ScraperIntegrationTestSuite) insertOldExternalSale(scrapedAt time.Time) {
	externalID := "old-sale-12345"
	sale := listing.Listing{
		ListingType:    "external",
		ExternalID:     &externalID,
		ExternalSource: stringPtr("OldSource"),
		ExternalURL:    stringPtr("https://example.com/old"),
		LastScrapedAt:  &scrapedAt,
		Title:          "Old Sale",
		AddressLine1:   "Old Address",
		City:           "Portland",
		State:          "OR",
		ZipCode:        "97201",
		StartDate:      time.Now(),
		EndDate:        time.Now().Add(24 * time.Hour),
		CreatedAt:      scrapedAt,
		UpdatedAt:      scrapedAt,
	}

	err := suite.repo.UpsertExternalSale(&sale)
	if err != nil {
		suite.T().Fatalf("Failed to insert old sale: %v", err)
	}
}

func createTestExternalSale(externalID, city, state string, scrapedAt time.Time) listing.Listing {
	return listing.Listing{
		ListingType:    "external",
		ExternalID:     &externalID,
		ExternalSource: stringPtr("TestSource"),
		ExternalURL:    stringPtr(fmt.Sprintf("https://example.com/%s", externalID)),
		LastScrapedAt:  &scrapedAt,
		Title:          fmt.Sprintf("Test Sale %s", externalID),
		AddressLine1:   "123 Test St",
		City:           city,
		State:          state,
		ZipCode:        "97201",
		StartDate:      scrapedAt,
		EndDate:        scrapedAt.Add(24 * time.Hour),
		CreatedAt:      scrapedAt,
		UpdatedAt:      scrapedAt,
	}
}

func stringPtr(s string) *string {
	return &s
}

// TestScraperIntegrationSuite runs the integration test suite
func TestScraperIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ScraperIntegrationTestSuite))
}
