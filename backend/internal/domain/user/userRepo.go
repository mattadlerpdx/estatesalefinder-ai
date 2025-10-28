package user

// Repository defines the interface for user data access operations.
type Repository interface {
	EnsureUser(firebaseUID string) error
	FindByFirebaseUID(firebaseUID string) (*User, error)
	GetBusinessIDsByUID(firebaseUID string) ([]int, error)
	LinkUserToBusiness(userID int, businessID int) error
	GetUserIDByFirebaseUID(firebaseUID string) (int, error)
}
