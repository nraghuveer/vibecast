package db

import (
	"fmt"
	"time"

	"github.com/nraghuveer/vibecast/lib/models"
)

type Template struct {
	ID        string
	Name      string
	Topic     string
	Persona   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (db *DB) CreateTemplate(t models.Template) error {
	query := `
		INSERT INTO templates (id, name, topic, persona)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			topic = excluded.topic,
			persona = excluded.persona,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := db.Exec(query, t.ID, t.Name, t.Topic, t.Persona)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	return nil
}

func (db *DB) GetTemplate(id string) (*Template, error) {
	query := `
		SELECT id, name, topic, persona, created_at, updated_at
		FROM templates
		WHERE id = ?
	`

	var t Template
	err := db.QueryRow(query, id).Scan(
		&t.ID,
		&t.Name,
		&t.Topic,
		&t.Persona,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return &t, nil
}

func (db *DB) GetAllTemplates() ([]Template, error) {
	query := `
		SELECT id, name, topic, persona, created_at, updated_at
		FROM templates
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get templates: %w", err)
	}
	defer rows.Close()

	var templates []Template
	for rows.Next() {
		var t Template
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Topic,
			&t.Persona,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}
		templates = append(templates, t)
	}

	return templates, nil
}

func (db *DB) UpdateTemplate(t models.Template) error {
	query := `
		UPDATE templates
		SET name = ?, topic = ?, persona = ?
		WHERE id = ?
	`

	result, err := db.Exec(query, t.Name, t.Topic, t.Persona, t.ID)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("template not found")
	}

	return nil
}

func (db *DB) DeleteTemplate(id string) error {
	query := `DELETE FROM templates WHERE id = ?`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("template not found")
	}

	return nil
}

func (db *DB) DeleteAllTemplates() error {
	query := `DELETE FROM templates`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete all templates: %w", err)
	}

	return nil
}

func (db *DB) TemplateExists(id string) (bool, error) {
	query := `SELECT COUNT(*) FROM templates WHERE id = ?`

	var count int
	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check template existence: %w", err)
	}

	return count > 0, nil
}
