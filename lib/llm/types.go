package llm

// ChatMessage is an OpenAI-compatible chat message.
// Role should be one of: system, user, assistant.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StreamEvent represents a streamed token (delta) or terminal event.
type StreamEvent struct {
	Delta string
	Done  bool
	Err   error
}
