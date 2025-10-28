package user

// User represents an authenticated application user.
type User struct {
	ID          int    `json:"id"`
	FirebaseUID string `json:"firebase_uid"`
	Email       string `json:"email,omitempty"`      // Optional if you want to store
	CreatedAt   string `json:"created_at,omitempty"` // Optional
}
