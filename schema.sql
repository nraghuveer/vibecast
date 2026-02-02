-- VibeCast Database Schema
-- This file contains the SQL schema for the VibeCast application database

-- Templates table: Stores predefined and custom templates
CREATE TABLE IF NOT EXISTS templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    topic TEXT NOT NULL,
    persona TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to automatically update the updated_at timestamp for templates
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

-- Enable foreign keys (if not already enabled)
PRAGMA foreign_keys = ON;
