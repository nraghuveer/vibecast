package models

import "time"

type Conversation struct {
	ID        string
	Title     string
	Topic     string
	Persona   string
	VoiceID   string
	VoiceName string
	Provider  string
	CreatedAt time.Time
	EndedAt   *time.Time
}
