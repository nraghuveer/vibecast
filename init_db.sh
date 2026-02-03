#!/bin/bash

# VibeCast Database Initialization Script
# This script creates a new VibeCast database with all required tables

set -e

DEFAULT_DB_PATH="${HOME}/.vibecast/data.sqlite"

if [ -z "$1" ]; then
    DB_PATH="$DEFAULT_DB_PATH"
else
    DB_PATH="$1"
fi

echo "Creating VibeCast database at: $DB_PATH"

# Create the database directory if it doesn't exist
DB_DIR=$(dirname "$DB_PATH")
mkdir -p "$DB_DIR"

# Create database and apply schema
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
SCHEMA_FILE="$SCRIPT_DIR/schema/v0.sql"

if [ ! -f "$SCHEMA_FILE" ]; then
    echo "Schema file not found: $SCHEMA_FILE"
    exit 1
fi

sqlite3 "$DB_PATH" < "$SCHEMA_FILE"

echo "Database created successfully at: $DB_PATH"
echo ""
echo "To verify the schema, run:"
echo "  sqlite3 \"$DB_PATH\" \".schema\""
