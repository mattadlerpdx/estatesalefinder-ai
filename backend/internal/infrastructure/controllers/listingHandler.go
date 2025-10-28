package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/listing"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/user"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/api"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/middleware"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/validation"
)

// ListingHandler handles HTTP requests for listings
type ListingHandler struct {
	listingService *listing.Service
	userService    *user.Service
	scraperService ScraperService // Will be injected
}

// ScraperService is the interface for the scraper
type ScraperService interface {
	GetListingsByLocation(city, state string) ([]listing.ScrapedListing, error)
}

// NewListingHandler creates a new listing handler
func NewListingHandler(listingService *listing.Service, userService *user.Service) *ListingHandler {
	return &ListingHandler{
		listingService: listingService,
		userService: userService,
	}
}

// SetScraperService sets the scraper service (called after initialization)
func (h *ListingHandler) SetScraperService(scraper ScraperService) {
	h.scraperService = scraper
}

// Create handles POST /api/sales
func (h *ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	uid := r.Context().Value(middleware.ContextKeyUID).(string)
	u, err := h.userService.GetOrCreateUser(uid, "")
	if err != nil {
		api.InternalErrorResponse(w, "Failed to get user")
		return
	}

	var s listing.Listing
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		api.ErrorResponseSingle(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set seller ID from authenticated user
	s.SellerID = &u.ID
	s.ListingType = "owned" // This is an owned listing

	// Validate required fields
	var validationErrors []string
	if err := validation.ValidateRequired(s.Title, "title"); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validation.ValidateRequired(s.City, "city"); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validation.ValidateRequired(s.State, "state"); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}

	if len(validationErrors) > 0 {
		api.ValidationErrorResponse(w, validationErrors)
		return
	}

	// Create the sale
	if err := h.listingService.CreateListing(&s); err != nil {
		api.InternalErrorResponse(w, fmt.Sprintf("Failed to create listing: %v", err))
		return
	}

	api.CreatedResponse(w, s, "Listing created successfully")
}

// GetByID handles GET /api/sales/:id
func (h *ListingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/sales/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.ErrorResponseSingle(w, "Invalid listing ID", http.StatusBadRequest)
		return
	}

	s, err := h.listingService.GetListingByID(id)
	if err != nil {
		api.NotFoundResponse(w, "Listing not found")
		return
	}

	api.OKResponse(w, s, "")
}

// GetAll handles GET /api/sales
func (h *ListingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	filters := listing.ListingFilters{
		City:     query.Get("city"),
		State:    query.Get("state"),
		ZipCode:  query.Get("zip_code"),
		EventType: query.Get("event_type"),
		Status:   "published", // Only show published sales to public
	}

	// Parse pagination
	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil {
			filters.Limit = limit
		}
	}
	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil {
			filters.Offset = offset
		}
	}

	// Parse date filters
	if startDateStr := query.Get("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			filters.StartDate = &startDate
		}
	}
	if endDateStr := query.Get("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			filters.EndDate = &endDate
		}
	}

	sales, err := h.listingService.GetAllListings(filters)
	if err != nil {
		api.InternalErrorResponse(w, "Failed to fetch listings")
		return
	}

	api.OKResponse(w, sales, "")
}

// GetMySales handles GET /api/my-sales (authenticated sellers only)
func (h *ListingHandler) GetMySales(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	uid := r.Context().Value(middleware.ContextKeyUID).(string)
	u, err := h.userService.GetOrCreateUser(uid, "")
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	sales, err := h.listingService.GetSellerListings(u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sales)
}

// Update handles PUT /api/sales/:id
func (h *ListingHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	uid := r.Context().Value(middleware.ContextKeyUID).(string)
	u, err := h.userService.GetOrCreateUser(uid, "")
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Extract ID from path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(pathParts[3])
	if err != nil {
		http.Error(w, "Invalid listing ID", http.StatusBadRequest)
		return
	}

	// Get existing sale to check ownership
	existingListing, err := h.listingService.GetListingByID(id)
	if err != nil {
		http.Error(w, "Sale not found", http.StatusNotFound)
		return
	}

	// Check if user owns this sale
	if existingListing.SellerID == nil || *existingListing.SellerID != u.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var s listing.Listing
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Preserve ID and seller ID
	s.ID = id
	s.SellerID = &u.ID
	s.ListingType = "owned"

	if err := h.listingService.UpdateListing(&s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// Delete handles DELETE /api/sales/:id
func (h *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	uid := r.Context().Value(middleware.ContextKeyUID).(string)
	u, err := h.userService.GetOrCreateUser(uid, "")
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Extract ID from path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(pathParts[3])
	if err != nil {
		http.Error(w, "Invalid listing ID", http.StatusBadRequest)
		return
	}

	// Get existing sale to check ownership
	existingListing, err := h.listingService.GetListingByID(id)
	if err != nil {
		http.Error(w, "Sale not found", http.StatusNotFound)
		return
	}

	// Check if user owns this sale
	if existingListing.SellerID == nil || *existingListing.SellerID != u.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if err := h.listingService.DeleteListing(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddImage handles POST /api/sales/:id/images
func (h *ListingHandler) AddImage(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	uid := r.Context().Value(middleware.ContextKeyUID).(string)
	u, err := h.userService.GetOrCreateUser(uid, "")
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Extract sale ID from path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	listingID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		http.Error(w, "Invalid listing ID", http.StatusBadRequest)
		return
	}

	// Verify ownership
	existingListing, err := h.listingService.GetListingByID(listingID)
	if err != nil {
		http.Error(w, "Sale not found", http.StatusNotFound)
		return
	}
	if existingListing.SellerID == nil || *existingListing.SellerID != u.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var img listing.ListingImage
	if err := json.NewDecoder(r.Body).Decode(&img); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	img.ListingID = listingID

	if err := h.listingService.AddListingImage(&img); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(img)
}

// GetAggregatedSales handles GET /api/sales/aggregated - combines owned + scraped
func (h *ListingHandler) GetAggregatedSales(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	city := query.Get("city")
	state := query.Get("state")

	// Default to Portland, OR if no location provided
	if city == "" && state == "" {
		city = "Portland"
		state = "OR"
	}

	var aggregatedListings []*listing.AggregatedListing

	// 1. Get owned sales from database
	filters := listing.ListingFilters{
		City:   city,
		State:  state,
		Status: "published",
	}

	ownedListings, err := h.listingService.GetAllListings(filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch owned sales: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert owned sales to aggregated format
	for _, s := range ownedListings {
		aggregatedListings = append(aggregatedListings, s.ToAggregatedSale())
	}

	// 2. Get scraped sales (if city and state provided AND scraper is enabled)
	if city != "" && state != "" && h.scraperService != nil {
		scrapedListings, err := h.scraperService.GetListingsByLocation(city, state)
		if err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to fetch scraped sales: %v\n", err)
		} else {
			// Convert scraped sales to aggregated format
			for _, s := range scrapedListings {
				aggregatedListings = append(aggregatedListings, s.ToAggregatedSale())
			}
		}
	}

	// 3. Sort by start_date (most recent/upcoming first)
	sort.Slice(aggregatedListings, func(i, j int) bool {
		return aggregatedListings[i].StartDate.After(aggregatedListings[j].StartDate)
	})

	api.OKResponse(w, map[string]interface{}{
		"sales": aggregatedListings,
		"total": len(aggregatedListings),
	}, "")
}
