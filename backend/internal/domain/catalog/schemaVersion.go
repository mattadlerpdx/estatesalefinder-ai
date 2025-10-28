package catalog

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SchemaVersion represents a detected CSV schema snapshot
type SchemaVersion struct {
	SchemaVersionID uuid.UUID       `json:"schema_version_id"`
	BusinessID      int             `json:"business_id"`
	IntegrationID   *int            `json:"integration_id,omitempty"`
	SourceName      string          `json:"source_name"`
	SourceNameNorm  string          `json:"source_name_norm"`  // Normalized: lowercase, no extension
	DetectedAt      time.Time       `json:"detected_at"`
	Headers         json.RawMessage `json:"headers"`         // JSON array of column names
	DataTypes       json.RawMessage `json:"data_types"`      // JSON object: {column: type}
	RowCount        *int            `json:"row_count,omitempty"`
	Profile         json.RawMessage `json:"profile,omitempty"` // Statistical profile
}

// SchemaVersionRepository defines the interface for schema version persistence
type SchemaVersionRepository interface {
	// Create saves a new schema version
	Create(sv *SchemaVersion) error

	// GetByID retrieves a schema version by ID
	GetByID(id uuid.UUID) (*SchemaVersion, error)

	// GetBySourceName retrieves the most recent schema version for a source (uses normalized name)
	GetBySourceName(businessID int, sourceName string) (*SchemaVersion, error)

	// GetBySourceNameNorm retrieves the most recent schema version by normalized name
	GetBySourceNameNorm(businessID int, sourceNameNorm string) (*SchemaVersion, error)

	// ListByBusiness retrieves all schema versions for a business
	ListByBusiness(businessID int, limit int) ([]*SchemaVersion, error)

	// GetByIntegration retrieves schema versions for a specific integration
	GetByIntegration(integrationID int) ([]*SchemaVersion, error)
}
