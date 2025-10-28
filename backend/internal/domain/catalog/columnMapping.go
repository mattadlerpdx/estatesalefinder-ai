package catalog

import (
	"time"

	"github.com/google/uuid"
)

// ColumnMapping represents an approved column mapping (reusable contract)
type ColumnMapping struct {
	MappingID       uuid.UUID  `json:"mapping_id"`
	BusinessID      int        `json:"business_id"`
	SourceName      string     `json:"source_name"`
	SourceNameNorm  string     `json:"source_name_norm"`  // Normalized: lowercase, no extension
	SchemaVersionID *uuid.UUID `json:"schema_version_id,omitempty"`
	SourceColumn    string     `json:"source_column"`
	CanonicalField  string     `json:"canonical_field"`
	Confidence      float64    `json:"confidence"`
	ApprovedBy      *string    `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// ColumnMappingRepository defines the interface for column mapping persistence
type ColumnMappingRepository interface {
	// Create saves a new column mapping
	Create(cm *ColumnMapping) error

	// CreateBatch saves multiple column mappings in one transaction
	CreateBatch(mappings []*ColumnMapping) error

	// GetBySourceName retrieves all approved mappings for a source
	GetBySourceName(businessID int, sourceName string) ([]*ColumnMapping, error)

	// GetBySchemaVersion retrieves all mappings for a schema version
	GetBySchemaVersion(schemaVersionID uuid.UUID) ([]*ColumnMapping, error)

	// Update updates an existing column mapping
	Update(cm *ColumnMapping) error

	// Delete removes a column mapping
	Delete(mappingID uuid.UUID) error

	// ListByBusiness retrieves all mappings for a business
	ListByBusiness(businessID int) ([]*ColumnMapping, error)
}
