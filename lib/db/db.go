package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nraghuveer/vibecast/lib/config"
)

var db *sql.DB

func Init() error {
	dbPath := config.GetDBPath()

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	var err error
	db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func createTables() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS templates (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			topic TEXT NOT NULL,
			persona TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TRIGGER IF NOT EXISTS update_templates_timestamp 
		AFTER UPDATE ON templates
		BEGIN
			UPDATE templates SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
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
	`)

	return err
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return db
}
