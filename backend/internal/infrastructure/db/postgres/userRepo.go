package postgres

import (
	"database/sql"
	"fmt"

	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/user"
)

// UserRepository provides access to the user data stored in PostgreSQL.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository initializes a new UserRepository with a PostgreSQL connection.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// EnsureUser inserts a user if they don't already exist (based on Firebase UID).
func (r *UserRepository) EnsureUser(firebaseUID string) error {
	// Check if user already exists
	existingUser, err := r.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return nil // Already exists
	}

	// Step 1: Create new business
	var businessID int
	err = r.db.QueryRow(`INSERT INTO businesses (name) VALUES ($1) RETURNING id`, "New Business").Scan(&businessID)
	if err != nil {
		return fmt.Errorf("failed to create new business: %w", err)
	}

	// Step 2: Insert user
	var userID int
	err = r.db.QueryRow(`INSERT INTO users (firebase_uid) VALUES ($1) RETURNING id`, firebaseUID).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to insert new user: %w", err)
	}

	// Step 3: Link user to business
	_, err = r.db.Exec(`INSERT INTO user_businesses (user_id, business_id, role) VALUES ($1, $2, $3)`, userID, businessID, "owner")
	if err != nil {
		return fmt.Errorf("failed to link user to business: %w", err)
	}

	return nil
}

// FindByFirebaseUID retrieves a user by their Firebase UID.
func (r *UserRepository) FindByFirebaseUID(firebaseUID string) (*user.User, error) {
	query := `SELECT id, firebase_uid, email, created_at FROM users WHERE firebase_uid = $1`
	var u user.User
	var email, createdAt sql.NullString
	err := r.db.QueryRow(query, firebaseUID).Scan(&u.ID, &u.FirebaseUID, &email, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to retrieve user by firebase_uid: %w", err)
	}
	// Handle NULL values
	if email.Valid {
		u.Email = email.String
	}
	if createdAt.Valid {
		u.CreatedAt = createdAt.String
	}
	return &u, nil
}

// GetUserIDByFirebaseUID retrieves the user ID for a given Firebase UID.
func (r *UserRepository) GetUserIDByFirebaseUID(firebaseUID string) (int, error) {
	var userID int
	query := `SELECT id FROM users WHERE firebase_uid = $1`
	err := r.db.QueryRow(query, firebaseUID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user not found for firebase_uid: %s", firebaseUID)
		}
		return 0, fmt.Errorf("failed to retrieve user_id for uid %s: %w", firebaseUID, err)
	}
	return userID, nil
}

// GetBusinessIDsByUID retrieves all business IDs for a given Firebase UID.
func (r *UserRepository) GetBusinessIDsByUID(firebaseUID string) ([]int, error) {
	query := `
		SELECT ub.business_id
		FROM user_businesses ub
		JOIN users u ON u.id = ub.user_id
		WHERE u.firebase_uid = $1
	`
	rows, err := r.db.Query(query, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve business_ids for uid %s: %w", firebaseUID, err)
	}
	defer rows.Close()

	var businessIDs []int
	for rows.Next() {
		var businessID int
		if err := rows.Scan(&businessID); err != nil {
			return nil, fmt.Errorf("failed to scan business_id: %w", err)
		}
		businessIDs = append(businessIDs, businessID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over business_ids: %w", err)
	}

	return businessIDs, nil
}

// LinkUserToBusiness creates a relationship between a user and a business.
func (r *UserRepository) LinkUserToBusiness(userID int, businessID int) error {
	query := `INSERT INTO user_businesses (user_id, business_id, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, userID, businessID, "owner")
	if err != nil {
		return fmt.Errorf("failed to link user %d to business %d: %w", userID, businessID, err)
	}
	return nil
}
