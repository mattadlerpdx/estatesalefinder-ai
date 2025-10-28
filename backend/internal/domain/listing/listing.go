package listing

import "time"

// Listing represents an estate sale listing (both owned and external)
type Listing struct {
	ID int `json:"id"`

	// Type discriminator
	ListingType string `json:"listing_type"` // 'owned' | 'external'

	// Owned listing fields
	SellerID *int `json:"seller_id,omitempty"` // NULL for external listings

	// External listing fields
	ExternalID     *string    `json:"external_id,omitempty"`      // e.g. "estatesale-finder-15436"
	ExternalSource *string    `json:"external_source,omitempty"`  // e.g. "EstateSale-Finder.com"
	ExternalURL    *string    `json:"external_url,omitempty"`     // Deep link to original
	LastScrapedAt  *time.Time `json:"last_scraped_at,omitempty"` // When we last scraped

	// Shared fields (all listings have these)
	Title       string `json:"title"`
	Description string `json:"description"`

	// Location
	AddressLine1 string   `json:"address_line1"`
	AddressLine2 *string  `json:"address_line2,omitempty"`
	City         string   `json:"city"`
	State        string   `json:"state"`
	ZipCode      string   `json:"zip_code"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`

	// Sale dates/times
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	SaleHours *string   `json:"sale_hours,omitempty"` // e.g., "Fri 9am-5pm, Sat 9am-3pm"

	// Owned-only fields (NULL for external)
	SaleType      string   `json:"sale_type,omitempty"`      // 'estate_sale', 'auction', 'moving_sale'
	Status        string   `json:"status,omitempty"`         // 'draft', 'published', 'completed', 'cancelled'
	ListingTier   string   `json:"listing_tier,omitempty"`   // 'basic', 'featured', 'premium'
	PaymentStatus string   `json:"payment_status,omitempty"` // 'unpaid', 'paid', 'refunded'
	AmountPaid    *float64 `json:"amount_paid,omitempty"`

	// Metadata
	ViewCount int       `json:"view_count"`
	Featured  bool      `json:"featured"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Related data (loaded separately)
	Images []ListingImage `json:"images,omitempty"`
}

// IsExternal returns true if this is an external/scraped listing
func (l *Listing) IsExternal() bool {
	return l.ListingType == "external"
}

// IsOwned returns true if this is an owned listing (user-created)
func (l *Listing) IsOwned() bool {
	return l.ListingType == "owned"
}

// ListingImage represents an image associated with a listing
type ListingImage struct {
	ID           int       `json:"id"`
	SaleID       int       `json:"sale_id"`
	ImageURL     string    `json:"image_url"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	IsPrimary    bool      `json:"is_primary"`
	DisplayOrder int       `json:"display_order"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

// ListingFilters represents query parameters for filtering listings
type ListingFilters struct {
	City      string
	State     string
	ZipCode   string
	SaleType  string // 'estate_sale', 'auction', 'moving_sale'
	Status    string // Only 'published' for public, all for sellers
	StartDate *time.Time
	EndDate   *time.Time
	Featured  *bool
	Limit     int
	Offset    int
}

// ToScrapedListing converts a Listing (external) to ScrapedListing for display
func (l *Listing) ToScrapedListing() ScrapedListing {
	scraped := ScrapedListing{
		ExternalID: "",
		Title:      l.Title,
		Address:    l.AddressLine1,
		City:       l.City,
		State:      l.State,
		ZipCode:    l.ZipCode,
		StartDate:  l.StartDate,
		EndDate:    l.EndDate,
		SourceName: "",
		SourceURL:  "",
	}

	// Populate external fields if present
	if l.ExternalID != nil {
		scraped.ExternalID = *l.ExternalID
	}
	if l.ExternalSource != nil {
		scraped.SourceName = *l.ExternalSource
	}
	if l.ExternalURL != nil {
		scraped.SourceURL = *l.ExternalURL
	}
	if l.LastScrapedAt != nil {
		scraped.ScrapedAt = *l.LastScrapedAt
		scraped.CachedAt = *l.LastScrapedAt
	}
	if l.Latitude != nil {
		scraped.Latitude = *l.Latitude
	}
	if l.Longitude != nil {
		scraped.Longitude = *l.Longitude
	}

	return scraped
}
