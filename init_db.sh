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
sqlite3 "$DB_PATH" <<EOF
-- Templates table: Stores predefined and custom templates
CREATE TABLE IF NOT EXISTS templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    topic TEXT NOT NULL,
    persona TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to automatically update updated_at timestamp for templates
CREATE TRIGGER IF NOT EXISTS update_templates_timestamp
AFTER UPDATE ON templates
BEGIN
    UPDATE templates SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Conversations table: Stores conversation metadata and index
CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    topic TEXT NOT NULL,
    persona TEXT NOT NULL,
    voice_id TEXT NOT NULL,
    voice_name TEXT NOT NULL,
    provider TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    ended_at DATETIME
);

-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON conversations(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_conversations_topic ON conversations(topic);
EOF

echo "Database created successfully at: $DB_PATH"
echo ""
echo "To verify the schema, run:"
echo "  sqlite3 \"$DB_PATH\" \".schema\""
