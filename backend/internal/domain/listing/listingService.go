package listing

import (
	"fmt"
	"time"
)

// Service handles business logic for sales
type Service struct {
	repo Repository
}

// NewService creates a new sale service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateSale creates a new sale listing
func (s *Service) CreateSale(l *Listing) error {
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
	if l.SaleType == "" {
		l.SaleType = "estate_sale"
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

// GetSaleByID retrieves a sale by ID and increments view count
func (s *Service) GetSaleByID(id int) (*Listing, error) {
	l, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Load images
	images, err := s.repo.GetImagesBySaleID(id)
	if err == nil {
		l.Images = images
	}

	// Increment view count asynchronously (don't block on errors)
	go s.repo.IncrementViewCount(id)

	return l, nil
}

// GetAllSales retrieves sales with optional filters
func (s *Service) GetAllSales(filters ListingFilters) ([]Listing, error) {
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
		images, err := s.repo.GetImagesBySaleID(listings[i].ID)
		if err == nil {
			listings[i].Images = images
		}
	}

	return listings, nil
}

// GetSellerSales retrieves all sales for a specific seller
func (s *Service) GetSellerSales(sellerID int) ([]Listing, error) {
	listings, err := s.repo.GetBySellerID(sellerID)
	if err != nil {
		return nil, err
	}

	// Load images for each listing
	for i := range listings {
		images, err := s.repo.GetImagesBySaleID(listings[i].ID)
		if err == nil {
			listings[i].Images = images
		}
	}

	return listings, nil
}

// UpdateSale updates an existing sale
func (s *Service) UpdateSale(l *Listing) error {
	// Validate
	if l.ID == 0 {
		return fmt.Errorf("sale ID is required")
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

// DeleteSale deletes a sale
func (s *Service) DeleteSale(id int) error {
	return s.repo.Delete(id)
}

// AddSaleImage adds an image to a sale
func (s *Service) AddSaleImage(image *ListingImage) error {
	if image.SaleID == 0 {
		return fmt.Errorf("sale ID is required")
	}
	if image.ImageURL == "" {
		return fmt.Errorf("image URL is required")
	}

	image.UploadedAt = time.Now()
	return s.repo.AddImage(image)
}

// DeleteSaleImage deletes an image
func (s *Service) DeleteSaleImage(imageID int) error {
	return s.repo.DeleteImage(imageID)
}

// SetPrimaryImage sets an image as the primary image for a sale
func (s *Service) SetPrimaryImage(imageID int, saleID int) error {
	return s.repo.SetPrimaryImage(imageID, saleID)
}
