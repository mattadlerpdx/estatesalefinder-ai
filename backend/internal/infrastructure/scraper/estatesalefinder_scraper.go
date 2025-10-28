package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/listing"
)

// EstateSaleFinderScraper scrapes estatesale-finder.com
type EstateSaleFinderScraper struct {
	httpClient *http.Client
}

// NewEstateSaleFinderScraper creates a new scraper
func NewEstateSaleFinderScraper() *EstateSaleFinderScraper {
	return &EstateSaleFinderScraper{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ScrapePortlandSales scrapes sales from Portland area
// Region IDs: 1=N Portland, 2=NW Portland, 3=NE Portland, 4=SE Portland, 5=SW Portland
func (s *EstateSaleFinderScraper) ScrapePortlandSales() ([]listing.ScrapedListing, error) {
	// All Portland regions
	regions := "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15"
	// All sale types: 1=Estate, 2=Moving, 4=Garage/Yard
	saleTypes := "1,2,4,5,7,8,9,10,11,12,13"

	url := fmt.Sprintf("https://www.estatesale-finder.com/all_sales_list.php?saletypeshow=%s&regionsshow=%s",
		saleTypes, regions)

	log.Printf("→ Scraping: %s", url)

	// Rate limiting
	time.Sleep(1 * time.Second)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("got status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var sales []listing.ScrapedListing

	// Find each sale row (from both "This Week's Sales" AND "Upcoming Sales" sections)
	doc.Find(".salerow").Each(func(i int, sel *goquery.Selection) {
		scraped := s.parseSaleRow(sel)
		if scraped != nil {
			sales = append(sales, *scraped)
		}
	})

	log.Printf("✓ Scraped %d sales from estatesale-finder.com (current + upcoming)", len(sales))
	return sales, nil
}

// parseSaleRow extracts data from a single sale row
func (s *EstateSaleFinderScraper) parseSaleRow(sel *goquery.Selection) *listing.ScrapedListing {
	// Get sale ID from id attribute (e.g., "sale15436")
	saleID, exists := sel.Attr("id")
	if !exists || !strings.HasPrefix(saleID, "sale") {
		return nil
	}
	externalID := fmt.Sprintf("estatesale-finder-%s", strings.TrimPrefix(saleID, "sale"))

	// Get provider/title
	title := strings.TrimSpace(sel.Find("h5 a").First().Text())
	if title == "" {
		title = "Estate Sale" // Fallback
	}

	// Get view link to construct source URL
	viewLink, _ := sel.Find("a.view").Attr("href")
	sourceURL := "https://www.estatesale-finder.com/" + viewLink

	// Get address (in .columns p elements)
	address := ""
	city := ""
	state := "OR"
	zipCode := ""

	sel.Find(".columns p").Each(func(i int, p *goquery.Selection) {
		text := strings.TrimSpace(p.Text())

		// Look for address pattern (contains city names or OR and numbers)
		if strings.Contains(text, "Portland") || strings.Contains(text, "Beaverton") ||
		   strings.Contains(text, "Gresham") || strings.Contains(text, "Mulino") ||
		   strings.Contains(text, "Milwaukie") || strings.Contains(text, "OR") {

			// Extract zip code (5 digits)
			parts := strings.Fields(text)
			for _, part := range parts {
				if len(part) == 5 && isNumeric(part) {
					zipCode = part
					break
				}
			}

			// Extract city (look for known cities)
			cities := []string{"Portland", "Beaverton", "Gresham", "Mulino", "Milwaukie", "Lake Oswego"}
			for _, c := range cities {
				if strings.Contains(text, c) {
					city = c
					break
				}
			}

			// Extract address - everything before city/state/zip
			// Remove "TBA" prefix if present
			addressText := strings.TrimPrefix(text, "TBA")
			addressText = strings.TrimSpace(addressText)

			// If we have a city, try to extract street address
			if city != "" {
				// Find the position of the city in the text
				cityIdx := strings.Index(addressText, city)
				if cityIdx > 0 {
					address = strings.TrimSpace(addressText[:cityIdx])
				} else if addressText != "" && !strings.HasPrefix(addressText, city) {
					// If city not found but we have text, use it
					address = addressText
				}
			}

			// If address is still empty or is just the city, mark as TBA
			if address == "" || address == city {
				address = "TBA"
			}
		}
	})

	// Get dates/times
	saleInfo := ""
	hours := ""

	sel.Find(".columns p").Each(func(i int, p *goquery.Selection) {
		text := strings.TrimSpace(p.Text())
		if strings.Contains(text, "Opens") || strings.Contains(text, "day") {
			saleInfo = text
		}
		if strings.Contains(text, "am") || strings.Contains(text, "pm") {
			if !strings.Contains(text, "Opens") {
				hours = text
			}
		}
	})

	// Parse dates (simplified - just use today + 7 days as placeholder)
	startDate := time.Now().AddDate(0, 0, 7)
	endDate := startDate.AddDate(0, 0, 2) // 3-day sale typically

	// Try to extract actual open date from saleInfo
	if strings.Contains(saleInfo, "Opens") {
		// Format: "Opens 31st Oct 10:00am"
		parts := strings.Fields(saleInfo)
		for j, part := range parts {
			if part == "Opens" && j+2 < len(parts) {
				// parts[j+1] = "31st", parts[j+2] = "Oct"
				dateStr := strings.Join(parts[j+1:j+3], " ")
				parsedDate := s.parseDate(dateStr)
				if !parsedDate.IsZero() {
					startDate = parsedDate
					endDate = startDate.AddDate(0, 0, 2)
				}
				break
			}
		}
	}

	description := fmt.Sprintf("%s\n\n%s\n%s", title, saleInfo, hours)

	return &listing.ScrapedListing{
		ExternalID:   externalID,
		Title:        title,
		Description:  description,
		Address:      address,
		City:         city,
		State:        state,
		ZipCode:      zipCode,
		StartDate:    startDate,
		EndDate:      endDate,
		ThumbnailURL: "", // No images in list view
		SourceName:   "EstateSale-Finder.com",
		SourceURL:    sourceURL,
		ScrapedAt:    time.Now(),
		CachedAt:     time.Now(),
	}
}

// parseDate parses dates like "31st Oct" into current year
func (s *EstateSaleFinderScraper) parseDate(dateStr string) time.Time {
	// Remove ordinal suffixes (st, nd, rd, th)
	dateStr = strings.ReplaceAll(dateStr, "st", "")
	dateStr = strings.ReplaceAll(dateStr, "nd", "")
	dateStr = strings.ReplaceAll(dateStr, "rd", "")
	dateStr = strings.ReplaceAll(dateStr, "th", "")

	// Add current year
	currentYear := time.Now().Year()
	dateStr = fmt.Sprintf("%s %d", dateStr, currentYear)

	// Try parsing
	layouts := []string{
		"2 Jan 2006",
		"02 Jan 2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t
		}
	}

	return time.Time{}
}

// isNumeric checks if string is all digits
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
