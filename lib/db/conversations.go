package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nraghuveer/vibecast/lib/models"
)

type Conversation struct {
	ID        string
	Title     string
	Topic     string
	Persona   string
	VoiceID   string
	VoiceName string
	Provider  string
	CreatedAt time.Time
	EndedAt   sql.NullTime
}

func (db *DB) CreateConversation(c models.Conversation) error {
	query := `
		INSERT INTO conversations (id, title, topic, persona, voice_id, voice_name, provider, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, c.ID, c.Title, c.Topic, c.Persona, c.VoiceID, c.VoiceName, c.Provider, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	return nil
}

func (db *DB) GetConversation(id string) (*Conversation, error) {
	query := `
		SELECT id, title, topic, persona, voice_id, voice_name, provider, created_at, ended_at
		FROM conversations
		WHERE id = ?
	`

	var c Conversation
	err := db.QueryRow(query, id).Scan(
		&c.ID,
		&c.Title,
		&c.Topic,
		&c.Persona,
		&c.VoiceID,
		&c.VoiceName,
		&c.Provider,
		&c.CreatedAt,
		&c.EndedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &c, nil
}

func (db *DB) GetAllConversations() ([]Conversation, error) {
	query := `
		SELECT id, title, topic, persona, voice_id, voice_name, provider, created_at, ended_at
		FROM conversations
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var c Conversation
		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Topic,
			&c.Persona,
			&c.VoiceID,
			&c.VoiceName,
			&c.Provider,
			&c.CreatedAt,
			&c.EndedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, c)
	}

	return conversations, nil
}

func (db *DB) UpdateConversationEndedAt(id string, endedAt time.Time) error {
	query := `
		UPDATE conversations
		SET ended_at = ?
		WHERE id = ?
	`

	result, err := db.Exec(query, endedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update conversation ended_at: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

func (db *DB) DeleteConversation(id string) error {
	query := `DELETE FROM conversations WHERE id = ?`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

func (db *DB) DeleteAllConversations() error {
	query := `DELETE FROM conversations`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete all conversations: %w", err)
	}

	return nil
}

func (db *DB) ConversationExists(id string) (bool, error) {
	query := `SELECT COUNT(*) FROM conversations WHERE id = ?`

	var count int
	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check conversation existence: %w", err)
	}

	return count > 0, nil
}
