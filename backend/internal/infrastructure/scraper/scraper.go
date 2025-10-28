package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/listing"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/cache"
)

// ScraperService handles web scraping with cache-through pattern
type ScraperService struct {
	cache          *cache.RedisClient
	repo           listing.Repository // Database for persistent storage
	cacheTTL       time.Duration
	httpClient     *http.Client
	scrapeInFlight sync.Map // Prevents duplicate scrapes
	esFinderScraper *EstateSaleFinderScraper
}

// NewScraperService creates a new scraper service
func NewScraperService(redisClient *cache.RedisClient, repo listing.Repository) *ScraperService {
	return &ScraperService{
		cache:    redisClient,
		repo:     repo,
		cacheTTL: 6 * time.Hour, // Cache for 6 hours
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		esFinderScraper: NewEstateSaleFinderScraper(),
	}
}

// GetSalesByLocation returns sales for a city/state (cached or scraped)
// Implements 3-tier strategy: Redis → PostgreSQL → Scrape
func (s *ScraperService) GetSalesByLocation(city, state string) ([]listing.ScrapedListing, error) {
	cacheKey := s.getCacheKey(city, state)

	// 1. Try Redis cache first (fastest)
	if s.cache.IsEnabled() {
		var cachedSales []listing.ScrapedListing
		err := s.cache.Get(cacheKey, &cachedSales)
		if err == nil {
			log.Printf("✓ Cache HIT (Redis): %s (%d sales)", cacheKey, len(cachedSales))
			return cachedSales, nil
		}
	}

	log.Printf("✗ Cache MISS (Redis): %s", cacheKey)

	// 2. Check PostgreSQL - is data less than 6 hours old?
	lastScrape, err := s.repo.GetLastScrapedTime(city, state)
	if err != nil {
		log.Printf("Warning: Failed to check PostgreSQL last scrape time: %v", err)
	}

	needsRescrape := true
	if lastScrape != nil && lastScrape.LastScrapedAt != nil {
		age := time.Since(*lastScrape.LastScrapedAt)
		log.Printf("→ PostgreSQL data age: %v (threshold: 6h)", age)

		if age < 6*time.Hour {
			needsRescrape = false
			log.Printf("✓ PostgreSQL data is fresh (<6h), loading from DB...")

			// Load from PostgreSQL
			dbSales, err := s.repo.GetExternalSalesByLocation(city, state)
			if err != nil {
				log.Printf("Warning: Failed to load from PostgreSQL: %v", err)
				needsRescrape = true // Fallback to scraping
			} else if len(dbSales) > 0 {
				// Convert Listing to ScrapedListing
				scrapedSales := make([]listing.ScrapedListing, len(dbSales))
				for i, s := range dbSales {
					scrapedSales[i] = s.ToScrapedListing()
				}

				// Re-cache in Redis
				if s.cache.IsEnabled() {
					if err := s.cache.Set(cacheKey, scrapedSales, s.cacheTTL); err != nil {
						log.Printf("Warning: Failed to re-cache from PostgreSQL: %v", err)
					} else {
						log.Printf("✓ Re-cached %d sales from PostgreSQL", len(scrapedSales))
					}
				}

				return scrapedSales, nil
			} else {
				log.Printf("Warning: PostgreSQL returned 0 sales, will re-scrape")
				needsRescrape = true
			}
		} else {
			log.Printf("→ PostgreSQL data is stale (>6h), re-scraping...")
		}
	} else {
		log.Printf("→ No PostgreSQL data found, initial scrape needed")
	}

	// 3. Needs scraping - check if scrape already in progress (request deduplication)
	if needsRescrape {
		ch := make(chan scrapeResult, 1)
		if existing, loaded := s.scrapeInFlight.LoadOrStore(cacheKey, ch); loaded {
			log.Printf("⏳ Scrape in progress for %s, waiting...", cacheKey)
			result := <-existing.(chan scrapeResult)
			return result.sales, result.err
		}

		// 4. We're first - do the scrape
		log.Printf("🌐 Scraping %s, %s...", city, state)

		sales, err := s.scrapeEstateSaleFinder(city, state)
		result := scrapeResult{sales: sales, err: err}

		// Notify all waiting goroutines
		ch <- result
		close(ch)
		s.scrapeInFlight.Delete(cacheKey)

		if err != nil {
			return nil, err
		}

		// 5. Persist to PostgreSQL (converts ScrapedListing → Sale)
		if s.repo != nil {
			log.Printf("→ Persisting %d sales to PostgreSQL...", len(sales))
			for _, scraped := range sales {
				saleEntity := scraped.ToSale()
				if err := s.repo.UpsertExternalSale(&saleEntity); err != nil {
					log.Printf("Warning: Failed to persist sale %s: %v", scraped.ExternalID, err)
				}
			}
			log.Printf("✓ Persisted %d sales to PostgreSQL", len(sales))
		}

		// 6. Store in Redis cache
		if s.cache.IsEnabled() {
			if err := s.cache.Set(cacheKey, sales, s.cacheTTL); err != nil {
				log.Printf("Warning: Failed to cache results: %v", err)
			} else {
				log.Printf("✓ Cached %d sales in Redis (TTL: %v)", len(sales), s.cacheTTL)
			}
		}

		return sales, nil
	}

	// Shouldn't reach here
	return []listing.ScrapedListing{}, nil
}

// scrapeEstateSaleFinder scrapes estatesale-finder.com for Portland area
func (s *ScraperService) scrapeEstateSaleFinder(city, state string) ([]listing.ScrapedListing, error) {
	// For now, only supports Portland area (estatesale-finder.com is regional)
	if strings.ToUpper(state) == "OR" || strings.ToLower(city) == "portland" {
		return s.esFinderScraper.ScrapePortlandSales()
	}

	// Fallback: return empty for other locations
	log.Printf("Note: estatesale-finder.com only covers Portland area. %s, %s not supported yet.", city, state)
	return []listing.ScrapedListing{}, nil
}

// scrapeEstateSalesNet scrapes estatesales.net for a city/state (legacy - not used)
func (s *ScraperService) scrapeEstateSalesNet(city, state string) ([]listing.ScrapedListing, error) {
	// Normalize inputs
	city = strings.ReplaceAll(city, " ", "-")
	state = strings.ToUpper(state)

	// Build URL: https://www.estatesales.net/OR/Portland
	url := fmt.Sprintf("https://www.estatesales.net/%s/%s", state, city)

	log.Printf("→ Scraping: %s", url)

	// Rate limiting: Sleep 1 second between requests
	time.Sleep(1 * time.Second)

	// Fetch HTML
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("got status code %d for %s", resp.StatusCode, url)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract sales
	var sales []listing.ScrapedListing

	// Find listing cards (you'll need to inspect HTML to get correct selectors)
	// This is a placeholder - actual selectors depend on estatesales.net structure
	doc.Find(".sale-item, .listing-card, article").Each(func(i int, sel *goquery.Selection) {
		scraped := s.parseSaleListing(sel, city, state)
		if scraped != nil {
			sales = append(sales, *scraped)
		}
	})

	log.Printf("✓ Scraped %d sales from %s", len(sales), url)
	return sales, nil
}

// parseSaleListing extracts data from a single listing HTML element
func (s *ScraperService) parseSaleListing(sel *goquery.Selection, city, state string) *listing.ScrapedListing {
	// Extract data (these selectors are placeholders - inspect actual HTML)
	title := strings.TrimSpace(sel.Find("h2, h3, .title").Text())
	if title == "" {
		return nil // Skip if no title
	}

	// Get source URL
	linkHref, exists := sel.Find("a").Attr("href")
	if !exists {
		return nil
	}

	// Make absolute URL
	sourceURL := linkHref
	if !strings.HasPrefix(linkHref, "http") {
		sourceURL = "https://www.estatesales.net" + linkHref
	}

	// Extract external ID from URL (e.g., /OR/Portland/12345 -> estatesales-net-12345)
	parts := strings.Split(strings.Trim(linkHref, "/"), "/")
	externalID := fmt.Sprintf("estatesales-net-%s", parts[len(parts)-1])

	// Get thumbnail
	thumbnailURL, _ := sel.Find("img").Attr("src")
	if !strings.HasPrefix(thumbnailURL, "http") {
		thumbnailURL = "https://www.estatesales.net" + thumbnailURL
	}

	// Get address (placeholder)
	address := strings.TrimSpace(sel.Find(".address, .location").Text())

	// Get dates (placeholder - parse from text)
	dateText := strings.TrimSpace(sel.Find(".date, .dates, time").Text())
	startDate, endDate := s.parseDates(dateText)

	return &listing.ScrapedListing{
		ExternalID:   externalID,
		Title:        title,
		Address:      address,
		City:         city,
		State:        state,
		StartDate:    startDate,
		EndDate:      endDate,
		ThumbnailURL: thumbnailURL,
		SourceName:   "EstateSales.net",
		SourceURL:    sourceURL,
		ScrapedAt:    time.Now(),
		CachedAt:     time.Now(),
	}
}

// parseDates attempts to parse date strings (placeholder implementation)
func (s *ScraperService) parseDates(dateText string) (time.Time, time.Time) {
	// This is a simplified version - you'll need to parse actual date formats
	now := time.Now()

	// Look for patterns like "Feb 1-2, 2025"
	// For now, return placeholder dates
	startDate := now.AddDate(0, 0, 7)  // 7 days from now
	endDate := now.AddDate(0, 0, 8)    // 8 days from now

	return startDate, endDate
}

// getCacheKey generates a cache key for city/state
func (s *ScraperService) getCacheKey(city, state string) string {
	return fmt.Sprintf("sales:%s:%s", strings.ToLower(city), strings.ToUpper(state))
}

// scrapeResult is used for request deduplication
type scrapeResult struct {
	sales []listing.ScrapedListing
	err   error
}

// InvalidateCache clears the cache for a city/state
func (s *ScraperService) InvalidateCache(city, state string) error {
	key := s.getCacheKey(city, state)
	return s.cache.Delete(key)
}
