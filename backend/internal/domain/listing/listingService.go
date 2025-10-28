package listing

import (
	"fmt"
	"time"
)

// Service handles business logic for sales
type Service struct {
	repo Repository
}

// NewService creates a new listing service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateListing creates a new listing listing
func (s *Service) CreateListing(l *Listing) error {
	// Validate required fields
	if l.Title == "" {
		return fmt.Errorf("title is required")
	}
	if l.City == "" || l.State == "" {
		return fmt.Errorf("city and state are required")
	}
	if l.StartDate.IsZero() || l.EndDate.IsZero() {
		return fmt.Errorf("start and end dates are required")
	}
	if l.EndDate.Before(l.StartDate) {
		return fmt.Errorf("end date must be after start date")
	}

	// Set defaults
	if l.Status == "" {
		l.Status = "draft"
	}
	if l.EventType == "" {
		l.EventType = "estate_sale"
	}
	if l.ListingTier == "" {
		l.ListingTier = "basic"
	}
	if l.PaymentStatus == "" {
		l.PaymentStatus = "unpaid"
	}

	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()

	return s.repo.Create(l)
}

// GetListingByID retrieves a listing by ID and increments view count
func (s *Service) GetListingByID(id int) (*Listing, error) {
	l, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Load images
	images, err := s.repo.GetImagesByListingID(id)
	if err == nil {
		l.Images = images
	}

	// Increment view count asynchronously (don't block on errors)
	go s.repo.IncrementViewCount(id)

	return l, nil
}

// GetAllListings retrieves sales with optional filters
func (s *Service) GetAllListings(filters ListingFilters) ([]Listing, error) {
	// Set default pagination
	if filters.Limit == 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	listings, err := s.repo.GetAll(filters)
	if err != nil {
		return nil, err
	}

	// Load images for each listing
	for i := range listings {
		images, err := s.repo.GetImagesByListingID(listings[i].ID)
		if err == nil {
			listings[i].Images = images
		}
	}

	return listings, nil
}

// GetSellerListings retrieves all sales for a specific seller
func (s *Service) GetSellerListings(sellerID int) ([]Listing, error) {
	listings, err := s.repo.GetBySellerID(sellerID)
	if err != nil {
		return nil, err
	}

	// Load images for each listing
	for i := range listings {
		images, err := s.repo.GetImagesByListingID(listings[i].ID)
		if err == nil {
			listings[i].Images = images
		}
	}

	return listings, nil
}

// UpdateListing updates an existing sale
func (s *Service) UpdateListing(l *Listing) error {
	// Validate
	if l.ID == 0 {
		return fmt.Errorf("listing ID is required")
	}
	if l.Title == "" {
		return fmt.Errorf("title is required")
	}
	if l.City == "" || l.State == "" {
		return fmt.Errorf("city and state are required")
	}

	l.UpdatedAt = time.Now()
	return s.repo.Update(l)
}

// DeleteListing deletes a sale
func (s *Service) DeleteListing(id int) error {
	return s.repo.Delete(id)
}

// AddListingImage adds an image to a listing
func (s *Service) AddListingImage(image *ListingImage) error {
	if image.ListingID == 0 {
		return fmt.Errorf("listing ID is required")
	}
	if image.ImageURL == "" {
		return fmt.Errorf("image URL is required")
	}

	image.UploadedAt = time.Now()
	return s.repo.AddImage(image)
}

// DeleteListingImage deletes an image
func (s *Service) DeleteListingImage(imageID int) error {
	return s.repo.DeleteImage(imageID)
}

// SetPrimaryImage sets an image as the primary image for a sale
func (s *Service) SetPrimaryImage(imageID int, listingID int) error {
	return s.repo.SetPrimaryImage(imageID, listingID)
}
