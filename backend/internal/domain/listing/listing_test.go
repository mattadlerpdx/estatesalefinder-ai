package listing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestListingCreation tests basic Listing struct creation
func TestListingCreation(t *testing.T) {
	now := time.Now()

	l := Listing{
		ID:           1,
		ListingType:  "owned",
		Title:        "Test Estate Sale",
		Description:  "A wonderful collection",
		AddressLine1: "123 Main St",
		City:         "Portland",
		State:        "OR",
		ZipCode:      "97201",
		StartDate:    now,
		EndDate:      now.Add(24 * time.Hour),
		SaleType:     "estate_sale",
		Status:       "active",
		ViewCount:    0,
		Featured:     false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	assert.Equal(t, "owned", l.ListingType)
	assert.Equal(t, "Test Estate Sale", l.Title)
	assert.True(t, l.IsOwned())
	assert.False(t, l.IsExternal())
}

// TestExternalListing tests external listing creation
func TestExternalListing(t *testing.T) {
	externalID := "test-123"
	source := "TestSource"
	url := "https://example.com/sale/123"
	now := time.Now()

	l := Listing{
		ID:             1,
		ListingType:    "external",
		ExternalID:     &externalID,
		ExternalSource: &source,
		ExternalURL:    &url,
		LastScrapedAt:  &now,
		Title:          "External Sale",
		City:           "Portland",
		State:          "OR",
		ZipCode:        "97201",
		StartDate:      now,
		EndDate:        now.Add(24 * time.Hour),
	}

	assert.True(t, l.IsExternal())
	assert.False(t, l.IsOwned())
	assert.Equal(t, "test-123", *l.ExternalID)
	assert.Equal(t, "TestSource", *l.ExternalSource)
}

// TestListingImage tests ListingImage struct
func TestListingImage(t *testing.T) {
	img := ListingImage{
		ID:           1,
		SaleID:       100,
		ImageURL:     "https://example.com/photo.jpg",
		IsPrimary:    true,
		DisplayOrder: 1,
		UploadedAt:   time.Now(),
	}

	assert.Equal(t, 100, img.SaleID)
	assert.True(t, img.IsPrimary)
	assert.Equal(t, "https://example.com/photo.jpg", img.ImageURL)
}

// TestListingFilters tests ListingFilters struct
func TestListingFilters(t *testing.T) {
	featured := true
	filters := ListingFilters{
		City:     "Portland",
		State:    "OR",
		SaleType: "estate_sale",
		Status:   "active",
		Featured: &featured,
		Limit:    20,
		Offset:   0,
	}

	assert.Equal(t, "Portland", filters.City)
	assert.Equal(t, "OR", filters.State)
	assert.Equal(t, 20, filters.Limit)
	assert.NotNil(t, filters.Featured)
	assert.True(t, *filters.Featured)
}

// TestScrapedListingCreation tests ScrapedListing struct
func TestScrapedListingCreation(t *testing.T) {
	now := time.Now()

	scraped := ScrapedListing{
		ExternalID:   "external-123",
		Title:        "Scraped Estate Sale",
		Description:  "Found online",
		Address:      "456 Oak Ave",
		City:         "Seattle",
		State:        "WA",
		ZipCode:      "98101",
		Latitude:     47.6062,
		Longitude:    -122.3321,
		StartDate:    now,
		EndDate:      now.Add(48 * time.Hour),
		ThumbnailURL: "https://example.com/thumb.jpg",
		ImageURLs:    []string{"https://example.com/img1.jpg"},
		SourceName:   "TestSource",
		SourceURL:    "https://example.com/sale/123",
		ScrapedAt:    now,
		CachedAt:     now,
	}

	assert.Equal(t, "external-123", scraped.ExternalID)
	assert.Equal(t, "Scraped Estate Sale", scraped.Title)
	assert.Equal(t, "Seattle", scraped.City)
	assert.Equal(t, 47.6062, scraped.Latitude)
}

// TestScrapedListingToListing tests conversion from ScrapedListing to Listing
func TestScrapedListingToListing(t *testing.T) {
	now := time.Now()

	scraped := ScrapedListing{
		ExternalID:   "convert-test-123",
		Title:        "Test Conversion",
		Description:  "Testing conversion",
		Address:      "789 Pine St",
		City:         "Portland",
		State:        "OR",
		ZipCode:      "97202",
		Latitude:     45.5155,
		Longitude:    -122.6789,
		StartDate:    now,
		EndDate:      now.Add(24 * time.Hour),
		ThumbnailURL: "https://example.com/thumb.jpg",
		SourceName:   "ConversionTest",
		SourceURL:    "https://example.com/test",
		ScrapedAt:    now,
		CachedAt:     now,
	}

	// Convert to Listing
	converted := scraped.ToSale()

	assert.Equal(t, "external", converted.ListingType)
	assert.NotNil(t, converted.ExternalID)
	assert.Equal(t, "convert-test-123", *converted.ExternalID)
	assert.Equal(t, "ConversionTest", *converted.ExternalSource)
	assert.Equal(t, "Test Conversion", converted.Title)
	assert.Equal(t, "Portland", converted.City)
	assert.Equal(t, "OR", converted.State)
	assert.NotNil(t, converted.Latitude)
	assert.Equal(t, 45.5155, *converted.Latitude)
}

// TestListingToScrapedListing tests conversion from Listing to ScrapedListing
func TestListingToScrapedListing(t *testing.T) {
	externalID := "reverse-test-456"
	source := "ReverseTest"
	url := "https://example.com/reverse"
	now := time.Now()
	lat := 45.5231
	lng := -122.6765

	listing := Listing{
		ID:             100,
		ListingType:    "external",
		ExternalID:     &externalID,
		ExternalSource: &source,
		ExternalURL:    &url,
		LastScrapedAt:  &now,
		Title:          "Reverse Conversion Test",
		Description:    "Testing reverse conversion",
		AddressLine1:   "321 Elm St",
		City:           "Beaverton",
		State:          "OR",
		ZipCode:        "97005",
		Latitude:       &lat,
		Longitude:      &lng,
		StartDate:      now,
		EndDate:        now.Add(24 * time.Hour),
	}

	// Convert to ScrapedListing
	scraped := listing.ToScrapedListing()

	assert.Equal(t, "reverse-test-456", scraped.ExternalID)
	assert.Equal(t, "Reverse Conversion Test", scraped.Title)
	assert.Equal(t, "ReverseTest", scraped.SourceName)
	assert.Equal(t, "https://example.com/reverse", scraped.SourceURL)
	assert.Equal(t, "Beaverton", scraped.City)
	assert.Equal(t, 45.5231, scraped.Latitude)
}

// TestAggregatedListing tests AggregatedListing conversion
func TestAggregatedListing(t *testing.T) {
	now := time.Now()

	// Test owned listing conversion
	lat := 45.5231
	lng := -122.6765
	addr2 := "Apt 2B"

	ownedListing := Listing{
		ID:           1,
		ListingType:  "owned",
		Title:        "Owned Sale",
		Description:  "My own sale",
		AddressLine1: "100 Main St",
		AddressLine2: &addr2,
		City:         "Portland",
		State:        "OR",
		ZipCode:      "97201",
		Latitude:     &lat,
		Longitude:    &lng,
		StartDate:    now,
		EndDate:      now.Add(24 * time.Hour),
		SaleType:     "estate_sale",
		Status:       "active",
		ViewCount:    10,
		Images: []ListingImage{
			{
				ImageURL:  "https://example.com/img1.jpg",
				IsPrimary: true,
			},
			{
				ImageURL:  "https://example.com/img2.jpg",
				IsPrimary: false,
			},
		},
	}

	aggregated := ownedListing.ToAggregatedSale()

	assert.False(t, aggregated.IsScraped)
	assert.Equal(t, "Owned Sale", aggregated.Title)
	assert.Equal(t, "100 Main St, Apt 2B", aggregated.Address)
	assert.Equal(t, "estate_sale", aggregated.SaleType)
	assert.Equal(t, "active", aggregated.Status)
	assert.Equal(t, 10, aggregated.ViewCount)
	assert.Len(t, aggregated.ImageURLs, 2)
	assert.Equal(t, "https://example.com/img1.jpg", aggregated.ThumbnailURL)

	// Test scraped listing conversion
	scrapedListing := ScrapedListing{
		ExternalID:   "agg-test-123",
		Title:        "Scraped Sale",
		Description:  "External sale",
		Address:      "200 Oak Ave",
		City:         "Seattle",
		State:        "WA",
		ZipCode:      "98101",
		Latitude:     47.6062,
		Longitude:    -122.3321,
		StartDate:    now,
		EndDate:      now.Add(24 * time.Hour),
		ThumbnailURL: "https://example.com/thumb.jpg",
		ImageURLs:    []string{"https://example.com/img1.jpg", "https://example.com/img2.jpg"},
		SourceName:   "ExternalSource",
		SourceURL:    "https://example.com/sale/123",
		ScrapedAt:    now,
		CachedAt:     now,
	}

	aggregatedScraped := scrapedListing.ToAggregatedSale()

	assert.True(t, aggregatedScraped.IsScraped)
	assert.Equal(t, "Scraped Sale", aggregatedScraped.Title)
	assert.NotNil(t, aggregatedScraped.Source)
	assert.Equal(t, "ExternalSource", aggregatedScraped.Source.Name)
	assert.Equal(t, "https://example.com/sale/123", aggregatedScraped.Source.URL)
	assert.Len(t, aggregatedScraped.ImageURLs, 2)
	assert.Equal(t, "https://example.com/thumb.jpg", aggregatedScraped.ThumbnailURL)
}
