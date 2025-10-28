package postgres

import (
	"database/sql"
	"fmt"

	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/listing"
)

// ListingRepository implements the listing.Repository interface
type ListingRepository struct {
	db *sql.DB
}

// NewListingRepository creates a new listing repository
func NewListingRepository(db *sql.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

// Create creates a new listing
func (r *ListingRepository) Create(s *listing.Listing) error {
	query := `
		INSERT INTO listings (
			listing_type, seller_id, title, description, event_type, status,
			address_line1, address_line2, city, state, zip_code, latitude, longitude,
			start_date, end_date, event_hours,
			listing_tier, payment_status, amount_paid,
			view_count, featured, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		s.ListingType, s.SellerID, s.Title, s.Description, s.EventType, s.Status,
		s.AddressLine1, s.AddressLine2, s.City, s.State, s.ZipCode, s.Latitude, s.Longitude,
		s.StartDate, s.EndDate, s.EventHours,
		s.ListingTier, s.PaymentStatus, s.AmountPaid,
		s.ViewCount, s.Featured, s.CreatedAt, s.UpdatedAt,
	).Scan(&s.ID)

	if err != nil {
		return fmt.Errorf("failed to create listing: %w", err)
	}

	return nil
}

// GetByID retrieves a listing by ID
func (r *ListingRepository) GetByID(id int) (*listing.Listing, error) {
	// Increment view count
	updateQuery := `UPDATE listings SET view_count = view_count + 1 WHERE id = $1`
	_, _ = r.db.Exec(updateQuery, id)

	query := `
		SELECT id, seller_id, title, description, event_type, status,
			address_line1, address_line2, city, state, zip_code, latitude, longitude,
			start_date, end_date, event_hours,
			listing_tier, payment_status, amount_paid,
			view_count, featured, created_at, updated_at,
			listing_type, external_id, external_source, external_url, last_scraped_at
		FROM listings
		WHERE id = $1
	`

	s := &listing.Listing{}
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.SellerID, &s.Title, &s.Description, &s.EventType, &s.Status,
		&s.AddressLine1, &s.AddressLine2, &s.City, &s.State, &s.ZipCode, &s.Latitude, &s.Longitude,
		&s.StartDate, &s.EndDate, &s.EventHours,
		&s.ListingTier, &s.PaymentStatus, &s.AmountPaid,
		&s.ViewCount, &s.Featured, &s.CreatedAt, &s.UpdatedAt,
		&s.ListingType, &s.ExternalID, &s.ExternalSource, &s.ExternalURL, &s.LastScrapedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("listing not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get sale: %w", err)
	}

	return s, nil
}

// GetAll retrieves sales with optional filters
func (r *ListingRepository) GetAll(filters listing.ListingFilters) ([]listing.Listing, error) {
	query := `
		SELECT id, seller_id, title, description, event_type, status,
			address_line1, address_line2, city, state, zip_code, latitude, longitude,
			start_date, end_date, event_hours,
			listing_tier, payment_status, amount_paid,
			view_count, featured, created_at, updated_at
		FROM listings
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1

	// Apply filters
	if filters.City != "" {
		query += fmt.Sprintf(" AND LOWER(city) = LOWER($%d)", argPos)
		args = append(args, filters.City)
		argPos++
	}
	if filters.State != "" {
		query += fmt.Sprintf(" AND LOWER(state) = LOWER($%d)", argPos)
		args = append(args, filters.State)
		argPos++
	}
	if filters.ZipCode != "" {
		query += fmt.Sprintf(" AND zip_code = $%d", argPos)
		args = append(args, filters.ZipCode)
		argPos++
	}
	if filters.EventType != "" {
		query += fmt.Sprintf(" AND event_type = $%d", argPos)
		args = append(args, filters.EventType)
		argPos++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, filters.Status)
		argPos++
	}
	if filters.Featured != nil {
		query += fmt.Sprintf(" AND featured = $%d", argPos)
		args = append(args, *filters.Featured)
		argPos++
	}
	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND start_date >= $%d", argPos)
		args = append(args, *filters.StartDate)
		argPos++
	}
	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND end_date <= $%d", argPos)
		args = append(args, *filters.EndDate)
		argPos++
	}

	// Order by featured first, then by start date
	query += " ORDER BY featured DESC, start_date DESC"

	// Pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sales: %w", err)
	}
	defer rows.Close()

	sales := []listing.Listing{}
	for rows.Next() {
		s := listing.Listing{}
		err := rows.Scan(
			&s.ID, &s.SellerID, &s.Title, &s.Description, &s.EventType, &s.Status,
			&s.AddressLine1, &s.AddressLine2, &s.City, &s.State, &s.ZipCode, &s.Latitude, &s.Longitude,
			&s.StartDate, &s.EndDate, &s.EventHours,
			&s.ListingTier, &s.PaymentStatus, &s.AmountPaid,
			&s.ViewCount, &s.Featured, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}
		sales = append(sales, s)
	}

	return sales, nil
}

// GetBySellerID retrieves all sales for a seller
func (r *ListingRepository) GetBySellerID(sellerID int) ([]listing.Listing, error) {
	query := `
		SELECT id, seller_id, title, description, event_type, status,
			address_line1, address_line2, city, state, zip_code, latitude, longitude,
			start_date, end_date, event_hours,
			listing_tier, payment_status, amount_paid,
			view_count, featured, created_at, updated_at
		FROM listings
		WHERE seller_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sales by seller: %w", err)
	}
	defer rows.Close()

	sales := []listing.Listing{}
	for rows.Next() {
		s := listing.Listing{}
		err := rows.Scan(
			&s.ID, &s.SellerID, &s.Title, &s.Description, &s.EventType, &s.Status,
			&s.AddressLine1, &s.AddressLine2, &s.City, &s.State, &s.ZipCode, &s.Latitude, &s.Longitude,
			&s.StartDate, &s.EndDate, &s.EventHours,
			&s.ListingTier, &s.PaymentStatus, &s.AmountPaid,
			&s.ViewCount, &s.Featured, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}
		sales = append(sales, s)
	}

	return sales, nil
}

// Update updates an existing sale
func (r *ListingRepository) Update(s *listing.Listing) error {
	query := `
		UPDATE listings SET
			title = $1, description = $2, event_type = $3, status = $4,
			address_line1 = $5, address_line2 = $6, city = $7, state = $8, zip_code = $9,
			latitude = $10, longitude = $11,
			start_date = $12, end_date = $13, event_hours = $14,
			listing_tier = $15, payment_status = $16, amount_paid = $17,
			featured = $18, updated_at = $19
		WHERE id = $20
	`

	result, err := r.db.Exec(
		query,
		s.Title, s.Description, s.EventType, s.Status,
		s.AddressLine1, s.AddressLine2, s.City, s.State, s.ZipCode,
		s.Latitude, s.Longitude,
		s.StartDate, s.EndDate, s.EventHours,
		s.ListingTier, s.PaymentStatus, s.AmountPaid,
		s.Featured, s.UpdatedAt,
		s.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("listing not found")
	}

	return nil
}

// Delete deletes a sale
func (r *ListingRepository) Delete(id int) error {
	query := `DELETE FROM listings WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("listing not found")
	}

	return nil
}

// IncrementViewCount increments the view count for a sale
func (r *ListingRepository) IncrementViewCount(id int) error {
	query := `UPDATE listings SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// AddImage adds an image to a listing
func (r *ListingRepository) AddImage(img *listing.ListingImage) error {
	query := `
		INSERT INTO listing_images (listing_id, image_url, thumbnail_url, is_primary, display_order, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		img.ListingID, img.ImageURL, img.ThumbnailURL, img.IsPrimary, img.DisplayOrder, img.UploadedAt,
	).Scan(&img.ID)

	if err != nil {
		return fmt.Errorf("failed to add image: %w", err)
	}

	return nil
}

// GetImagesByListingID retrieves all images for a listing
func (r *ListingRepository) GetImagesByListingID(listingID int) ([]listing.ListingImage, error) {
	query := `
		SELECT id, listing_id, image_url, thumbnail_url, is_primary, display_order, uploaded_at
		FROM listing_images
		WHERE listing_id = $1
		ORDER BY is_primary DESC, display_order ASC
	`

	rows, err := r.db.Query(query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query images: %w", err)
	}
	defer rows.Close()

	images := []listing.ListingImage{}
	for rows.Next() {
		img := listing.ListingImage{}
		err := rows.Scan(
			&img.ID, &img.ListingID, &img.ImageURL, &img.ThumbnailURL,
			&img.IsPrimary, &img.DisplayOrder, &img.UploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}
		images = append(images, img)
	}

	return images, nil
}

// DeleteImage deletes an image
func (r *ListingRepository) DeleteImage(imageID int) error {
	query := `DELETE FROM listing_images WHERE id = $1`

	result, err := r.db.Exec(query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("image not found")
	}

	return nil
}

// SetPrimaryImage sets an image as primary and unsets all others for that listing
func (r *ListingRepository) SetPrimaryImage(imageID int, listingID int) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Unset all primary images for this sale
	_, err = tx.Exec(`UPDATE listing_images SET is_primary = false WHERE listing_id = $1`, listingID)
	if err != nil {
		return fmt.Errorf("failed to unset primary images: %w", err)
	}

	// Set the new primary image
	result, err := tx.Exec(
		`UPDATE listing_images SET is_primary = true WHERE id = $1 AND listing_id = $2`,
		imageID, listingID,
	)
	if err != nil {
		return fmt.Errorf("failed to set primary image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("image not found or doesn't belong to this listing")
	}

	return tx.Commit()
}

// UpsertExternalSale inserts or updates an external sale (uses external_id for conflict detection)
func (r *ListingRepository) UpsertExternalSale(s *listing.Listing) error {
	query := `
		INSERT INTO listings (
			listing_type, external_id, external_source, external_url,
			title, description,
			address_line1, address_line2, city, state, zip_code, latitude, longitude,
			start_date, end_date, event_hours,
			view_count, featured, last_scraped_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		ON CONFLICT (external_id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			address_line1 = EXCLUDED.address_line1,
			address_line2 = EXCLUDED.address_line2,
			city = EXCLUDED.city,
			state = EXCLUDED.state,
			zip_code = EXCLUDED.zip_code,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			start_date = EXCLUDED.start_date,
			end_date = EXCLUDED.end_date,
			event_hours = EXCLUDED.event_hours,
			last_scraped_at = EXCLUDED.last_scraped_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		s.ListingType, s.ExternalID, s.ExternalSource, s.ExternalURL,
		s.Title, s.Description,
		s.AddressLine1, s.AddressLine2, s.City, s.State, s.ZipCode, s.Latitude, s.Longitude,
		s.StartDate, s.EndDate, s.EventHours,
		s.ViewCount, s.Featured, s.LastScrapedAt, s.CreatedAt, s.UpdatedAt,
	).Scan(&s.ID)

	if err != nil {
		return fmt.Errorf("failed to upsert external listing: %w", err)
	}

	return nil
}

// GetExternalSalesByLocation retrieves external sales for a city/state
func (r *ListingRepository) GetExternalSalesByLocation(city, state string) ([]listing.Listing, error) {
	query := `
		SELECT id, listing_type, external_id, external_source, external_url,
			title, description,
			address_line1, address_line2, city, state, zip_code, latitude, longitude,
			start_date, end_date, event_hours,
			view_count, featured, last_scraped_at, created_at, updated_at
		FROM listings
		WHERE listing_type = 'external'
			AND LOWER(city) = LOWER($1)
			AND LOWER(state) = LOWER($2)
		ORDER BY start_date DESC
	`

	rows, err := r.db.Query(query, city, state)
	if err != nil {
		return nil, fmt.Errorf("failed to query external sales: %w", err)
	}
	defer rows.Close()

	sales := []listing.Listing{}
	for rows.Next() {
		s := listing.Listing{}
		err := rows.Scan(
			&s.ID, &s.ListingType, &s.ExternalID, &s.ExternalSource, &s.ExternalURL,
			&s.Title, &s.Description,
			&s.AddressLine1, &s.AddressLine2, &s.City, &s.State, &s.ZipCode, &s.Latitude, &s.Longitude,
			&s.StartDate, &s.EndDate, &s.EventHours,
			&s.ViewCount, &s.Featured, &s.LastScrapedAt, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan external sale: %w", err)
		}
		sales = append(sales, s)
	}

	return sales, nil
}

// GetLastScrapedTime retrieves the last scraped timestamp for a location
func (r *ListingRepository) GetLastScrapedTime(city, state string) (*listing.Listing, error) {
	query := `
		SELECT id, last_scraped_at
		FROM listings
		WHERE listing_type = 'external'
			AND LOWER(city) = LOWER($1)
			AND LOWER(state) = LOWER($2)
		ORDER BY last_scraped_at DESC NULLS LAST
		LIMIT 1
	`

	s := &listing.Listing{}
	err := r.db.QueryRow(query, city, state).Scan(&s.ID, &s.LastScrapedAt)

	if err == sql.ErrNoRows {
		// No external sales found for this location - needs initial scrape
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get last scraped time: %w", err)
	}

	return s, nil
}
