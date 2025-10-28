package user

// Service provides user-related operations.
type Service struct {
	repository Repository
}

// NewService creates a new UserService instance.
func NewService(repo Repository) *Service {
	return &Service{repository: repo}
}

// EnsureUserExists checks if a Firebase UID exists in the DB, inserts if not.
func (s *Service) EnsureUserExists(firebaseUID string) error {
	return s.repository.EnsureUser(firebaseUID)
}

// GetBusinessIDsByUID retrieves all business_ids for a given Firebase UID.
func (s *Service) GetBusinessIDsByUID(firebaseUID string) ([]int, error) {
	return s.repository.GetBusinessIDsByUID(firebaseUID)
}

// LinkUserToBusiness creates a relationship between a user and a business.
func (s *Service) LinkUserToBusiness(userID int, businessID int) error {
	return s.repository.LinkUserToBusiness(userID, businessID)
}

// GetUserIDByFirebaseUID retrieves the user ID for a given Firebase UID.
func (s *Service) GetUserIDByFirebaseUID(firebaseUID string) (int, error) {
	return s.repository.GetUserIDByFirebaseUID(firebaseUID)
}

// GetOrCreateUser gets a user by Firebase UID, creating them if they don't exist.
func (s *Service) GetOrCreateUser(firebaseUID string, email string) (*User, error) {
	// Try to find existing user
	user, err := s.repository.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}

	// If user exists, return them
	if user != nil {
		return user, nil
	}

	// User doesn't exist, create them
	err = s.repository.EnsureUser(firebaseUID)
	if err != nil {
		return nil, err
	}

	// Fetch and return the newly created user
	return s.repository.FindByFirebaseUID(firebaseUID)
}
