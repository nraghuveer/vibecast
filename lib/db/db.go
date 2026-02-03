package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nraghuveer/vibecast/lib/config"
)

// DB wraps sql.DB and provides database operations
type DB struct {
	*sql.DB
}

// NewDB creates and initializes a new database instance
func NewDB() (*DB, error) {
	dbPath := config.GetDBPath()

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	sqlDB, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{DB: sqlDB}

	if err := db.createTables(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func (db *DB) createTables() error {
	schemaSQL, err := os.ReadFile("schema/v0.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	_, err = db.Exec(string(schemaSQL))
	return err
}
