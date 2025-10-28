#!/bin/bash

# Migration Runner Script
# Usage: ./run-migration.sh <migration_file>

set -e

MIGRATION_FILE=$1

if [ -z "$MIGRATION_FILE" ]; then
    echo "Usage: ./run-migration.sh <migration_file>"
    echo "Example: ./run-migration.sh migrations/001_add_industry_support.sql"
    exit 1
fi

if [ ! -f "$MIGRATION_FILE" ]; then
    echo "Error: Migration file '$MIGRATION_FILE' not found"
    exit 1
fi

# Load environment variables from .env
if [ -f ../.env ]; then
    export $(cat ../.env | grep -v '^#' | xargs)
fi

# Database connection details from environment or defaults
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-cadence}

echo "========================================="
echo "Running Migration: $MIGRATION_FILE"
echo "Database: $DB_NAME@$DB_HOST:$DB_PORT"
echo "========================================="

# Run migration using docker exec (if using Docker) or psql directly
if docker ps | grep -q cadence-db; then
    echo "Using Docker container 'cadence-db'..."
    docker exec -i cadence-db psql -U $DB_USER -d $DB_NAME < $MIGRATION_FILE
else
    echo "Using local psql..."
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME < $MIGRATION_FILE
fi

echo ""
echo "✅ Migration completed successfully!"
echo ""
echo "To verify, run:"
echo "  docker exec -it cadence-db psql -U postgres -d cadence"
echo "  Then: \\dt (list tables)"
echo "       \\d integrations (describe table)"
