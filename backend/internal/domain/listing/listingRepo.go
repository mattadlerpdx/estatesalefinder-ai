package listing

// Repository defines the interface for listing data operations
type Repository interface {
	// Listing CRUD
	Create(listing *Listing) error
	GetByID(id int) (*Listing, error)
	GetAll(filters ListingFilters) ([]Listing, error)
	GetBySellerID(sellerID int) ([]Listing, error)
	Update(listing *Listing) error
	Delete(id int) error
	IncrementViewCount(id int) error

	// Image operations
	AddImage(image *ListingImage) error
	GetImagesByListingID(listingID int) ([]ListingImage, error)
	DeleteImage(imageID int) error
	SetPrimaryImage(imageID int, listingID int) error

	// External listing operations
	UpsertExternalSale(listing *Listing) error
	GetExternalSalesByLocation(city, state string) ([]Listing, error)
	GetLastScrapedTime(city, state string) (*Listing, error)
}
