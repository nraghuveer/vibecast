package models

// Template represents a topic/persona combination
type Template struct {
	ID      string
	Name    string
	Topic   string
	Persona string
}

// SpeakerType represents who is speaking in a conversation
type SpeakerType int

const (
	HOST SpeakerType = iota
	GUEST
)

// String returns the string representation of the speaker type
func (s SpeakerType) String() string {
	switch s {
	case HOST:
		return "Host"
	case GUEST:
		return "Guest"
	default:
		return "Unknown"
	}
}

// ParseSpeakerType converts a string to SpeakerType
func ParseSpeakerType(s string) SpeakerType {
	switch s {
	case "Host":
		return HOST
	case "Guest":
		return GUEST
	default:
		return GUEST // Default to guest for unknown speakers
	}
}
