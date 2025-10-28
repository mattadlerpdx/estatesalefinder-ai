package listing

import (
	"fmt"
	"time"
)

// ScrapedListing represents a minimal estate sale listing from external sources
// We only store metadata + URLs (not full content) to keep storage cheap
type ScrapedListing struct {
	// External identifier (source + their ID)
	ExternalID string `json:"external_id"` // e.g. "estatesales-net-12345"

	// Basic info
	Title       string `json:"title"`
	Description string `json:"description,omitempty"` // Optional short snippet

	// Location (for search/filtering)
	Address  string   `json:"address"`           // Full address string
	City     string   `json:"city"`
	State    string   `json:"state"`
	ZipCode  string   `json:"zip_code"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`

	// Dates
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	// Images (just URLs, we don't download)
	ThumbnailURL string   `json:"thumbnail_url"`           // Primary thumbnail
	ImageURLs    []string `json:"image_urls,omitempty"`    // Additional images

	// Source attribution
	SourceName string `json:"source_name"` // "EstateSales.net"
	SourceURL  string `json:"source_url"`  // Deep link to original listing

	// Metadata
	ScrapedAt time.Time `json:"scraped_at"`
	CachedAt  time.Time `json:"cached_at"`
}

// AggregatedListing combines owned and scraped sales for frontend
type AggregatedListing struct {
	// Union of Sale + ScrapedListing fields
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description,omitempty"`
	Address      string     `json:"address"`
	City         string     `json:"city"`
	State        string     `json:"state"`
	ZipCode      string     `json:"zip_code"`
	Latitude     *float64   `json:"latitude,omitempty"`
	Longitude    *float64   `json:"longitude,omitempty"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      time.Time  `json:"end_date"`
	ThumbnailURL string     `json:"thumbnail_url"`
	ImageURLs    []string   `json:"image_urls,omitempty"`

	// Type indicator
	IsScraped bool `json:"is_scraped"` // true = external, false = owned

	// Source (only for scraped)
	Source *struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"source,omitempty"`

	// Owned sale data (only if IsScraped = false)
	EventType      string `json:"event_type,omitempty"`
	Status        string `json:"status,omitempty"`
	ViewCount     int    `json:"view_count,omitempty"`
}

// ToAggregatedSale converts owned Sale to AggregatedListing
func (s *Listing) ToAggregatedSale() *AggregatedListing {
	var thumbnailURL string
	var imageURLs []string

	// Get primary image as thumbnail
	for _, img := range s.Images {
		if img.IsPrimary {
			thumbnailURL = img.ImageURL
		}
		imageURLs = append(imageURLs, img.ImageURL)
	}

	address := s.AddressLine1
	if s.AddressLine2 != nil && *s.AddressLine2 != "" {
		address += ", " + *s.AddressLine2
	}

	return &AggregatedListing{
		ID:           fmt.Sprintf("%d", s.ID), // Convert int to string
		Title:        s.Title,
		Description:  s.Description,
		Address:      address,
		City:         s.City,
		State:        s.State,
		ZipCode:      s.ZipCode,
		Latitude:     s.Latitude,
		Longitude:    s.Longitude,
		StartDate:    s.StartDate,
		EndDate:      s.EndDate,
		ThumbnailURL: thumbnailURL,
		ImageURLs:    imageURLs,
		IsScraped:    false,
		EventType:     s.EventType,
		Status:       s.Status,
		ViewCount:    s.ViewCount,
	}
}

// ToAggregatedSale converts ScrapedListing to AggregatedListing
func (s *ScrapedListing) ToAggregatedSale() *AggregatedListing {
	return &AggregatedListing{
		ID:           s.ExternalID,
		Title:        s.Title,
		Description:  s.Description,
		Address:      s.Address,
		City:         s.City,
		State:        s.State,
		ZipCode:      s.ZipCode,
		Latitude:     &s.Latitude,
		Longitude:    &s.Longitude,
		StartDate:    s.StartDate,
		EndDate:      s.EndDate,
		ThumbnailURL: s.ThumbnailURL,
		ImageURLs:    s.ImageURLs,
		IsScraped:    true,
		Source: &struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			Name: s.SourceName,
			URL:  s.SourceURL,
		},
	}
}

// ToSale converts ScrapedListing to Sale (for PostgreSQL persistence)
func (s *ScrapedListing) ToSale() Listing {
	now := time.Now()

	return Listing{
		ListingType:    "external",
		ExternalID:     &s.ExternalID,
		ExternalSource: &s.SourceName,
		ExternalURL:    &s.SourceURL,
		LastScrapedAt:  &now,

		Title:        s.Title,
		Description:  s.Description,
		AddressLine1: s.Address,
		City:         s.City,
		State:        s.State,
		ZipCode:      s.ZipCode,
		Latitude:     &s.Latitude,
		Longitude:    &s.Longitude,

		StartDate: s.StartDate,
		EndDate:   s.EndDate,

		ViewCount: 0,
		Featured:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
