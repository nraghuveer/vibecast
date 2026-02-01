# VibeCast Implementation Plan

## Project Structure

```
vibecast/
├── cmd/
│   └── vibecast/
│       └── main.go              # Entry point
├── internal/
│   ├── ai/
│   │   ├── provider.go          # AI provider interface
│   │   ├── context.go           # Hybrid context manager
│   │   ├── openai/
│   │   │   └── client.go        # OpenAI implementation
│   │   └── anthropic/
│   │       └── client.go        # Anthropic implementation (optional)
│   ├── audio/
│   │   ├── player.go            # Audio playback interface
│   │   ├── tts.go               # Text-to-speech interface
│   │   └── voices.go            # Voice definitions
│   ├── podcast/
│   │   ├── session.go           # Podcast session management
│   │   ├── transcript.go        # Transcript handling
│   │   ├── search.go            # Transcript search/lookup
│   │   └── persona.go           # Persona definitions
│   ├── ui/
│   │   ├── app.go               # Main Bubble Tea application
│   │   ├── views/
│   │   │   ├── welcome.go       # Welcome screen
│   │   │   ├── topic.go         # Topic selection
│   │   │   ├── persona.go       # Persona selection
│   │   │   ├── voice.go         # Voice selection
│   │   │   ├── conversation.go  # Main conversation view
│   │   │   └── transcript.go    # End-of-session transcript view
│   │   ├── components/
│   │   │   ├── input.go         # Text input component
│   │   │   ├── spinner.go       # Loading spinner
│   │   │   └── typewriter.go    # Typewriter animation effect
│   │   └── styles/
│   │       └── styles.go        # Lip Gloss styles
│   └── config/
│       └── config.go            # Configuration management
├── go.mod
├── go.sum
└── README.md
```

---

## Phase 1: Project Setup & Core Interfaces

### Phase 1.5: Context Management (Hybrid Approach)
**Critical for long conversations (30+ minutes)**

```go
type ContextManager struct {
    recentMessages  []Message  // Last N messages (verbatim)
    summary         string     // Rolling summary of older content
    maxRecent       int        // e.g., 20 messages
    summaryInterval int        // Summarize every N messages
}
```

Strategy:
- Keep last 20 messages verbatim
- When messages exceed threshold, summarize oldest batch
- Prepend summary to context on each API call
- Summary updated incrementally (not full re-summarization)

**File to create:** `internal/ai/context.go` - hybrid context manager

---

### 1.1 Initialize Go Module
- Create `go.mod` with module path
- Add dependencies:
  - `github.com/charmbracelet/bubbletea` - TUI framework
  - `github.com/charmbracelet/lipgloss` - Styling
  - `github.com/charmbracelet/bubbles` - UI components
  - `github.com/sashabaranov/go-openai` - OpenAI SDK

### 1.2 Define AI Provider Interface
```go
type AIProvider interface {
    StreamResponse(ctx context.Context, messages []Message, opts StreamOptions) (<-chan StreamChunk, error)
    GenerateSpeech(ctx context.Context, text string, voice Voice) ([]byte, error)
}

type Message struct {
    Role    string // "system", "user", "assistant"
    Content string
}

type StreamChunk struct {
    Text  string
    Done  bool
    Error error
}
```

### 1.3 Define TTS Interface
```go
type TTSProvider interface {
    Speak(ctx context.Context, text string, voice Voice) error
    ListVoices() []Voice
}

type Voice struct {
    ID          string
    Name        string
    Description string
}
```

---

## Phase 2: UI Foundation (Bubble Tea)

### 2.1 Main Application Model
```go
type Model struct {
    state       AppState
    topic       string
    persona     Persona
    voice       Voice
    session     *podcast.Session
    // Sub-models for each view
    welcomeView    welcome.Model
    topicView      topic.Model
    personaView    persona.Model
    voiceView      voice.Model
    conversationView conversation.Model
}

type AppState int
const (
    StateWelcome AppState = iota
    StateTopicSelection
    StatePersonaSelection
    StateVoiceSelection
    StateConversation
    StateTranscript
)
```

### 2.2 Welcome Screen
- Display VibeCast logo/banner
- Brief instructions
- Press Enter to continue

### 2.3 Topic Selection View
- Text input for custom topic
- Optional: predefined topic suggestions
- Validation (non-empty topic)

### 2.4 Persona Selection View
- List of predefined personas (expert, celebrity, fictional character)
- Option to create custom persona
- Each persona includes: name, description, system prompt

### 2.5 Voice Selection View
- List available TTS voices
- Preview voice (optional)
- Display voice characteristics

---

## Phase 3: Conversation Engine

### 3.1 Session Management
```go
type Session struct {
    Topic      string
    Persona    Persona
    Voice      Voice
    Messages   []Message
    StartTime  time.Time
    Transcript []TranscriptEntry
}

type TranscriptEntry struct {
    Speaker   string // "Host" or persona name
    Text      string
    Timestamp time.Time
}
```

### 3.2 System Prompt Construction
Build a system prompt that:
- Defines the podcast context
- Sets persona characteristics
- Embeds all 12 guidelines from spec
- Instructs conversational, engaging responses

### 3.3 Conversation View
- Split view: transcript history (scrollable) + input area
- Real-time streaming with typewriter animation
- Visual indicator when AI is "speaking"
- Interrupt handling (cancel current stream)
- Input disabled while AI is responding (or allow interrupt)

### 3.4 Interruption Handling
- User can press a key (e.g., Ctrl+C or Esc) to interrupt
- Cancel current AI stream
- AI acknowledges interruption in next response
- Track interruption state in context

---

## Phase 4: AI Integration

### 4.1 OpenAI Implementation
- Implement `AIProvider` interface
- Use streaming chat completions API
- Configure model (gpt-4o recommended)
- Handle rate limits and errors

### 4.2 Streaming Handler
- Process SSE chunks
- Send chunks to UI channel
- Handle connection errors gracefully
- Support cancellation via context

### 4.3 TTS Integration (OpenAI TTS)
- Use OpenAI TTS API
- Stream audio playback
- Support multiple voices (alloy, echo, fable, onyx, nova, shimmer)

---

## Phase 5: Typewriter Animation

### 5.1 Typewriter Component
```go
type TypewriterModel struct {
    fullText    string
    displayed   string
    charIndex   int
    speed       time.Duration
    done        bool
}
```
- Tick-based character reveal
- Configurable speed
- Smooth animation feel

### 5.2 Integration with Conversation
- Feed streamed chunks to typewriter
- Buffer and reveal at natural pace
- Handle rapid chunk arrival

---

## Phase 6: Transcription & Live Search

### 6.1 Real-time Transcription
- Log each message as it completes
- Store timestamps
- Track speaker attribution
- Maintain live, searchable transcript throughout session

### 6.2 Transcript Search (Keyword Lookup)
```go
type TranscriptSearch struct {
    entries []TranscriptEntry
    index   map[string][]int  // keyword -> entry indices
}

func (ts *TranscriptSearch) Search(keyword string) []TranscriptEntry
func (ts *TranscriptSearch) SearchBySpeaker(speaker, keyword string) []TranscriptEntry
func (ts *TranscriptSearch) GetContext(keyword string, windowSize int) string
```

Features:
- User can reference past conversation ("What did you say about X?")
- AI can query transcript to recall if a topic was discussed
- Case-insensitive keyword matching
- Returns relevant transcript entries with timestamps
- Context window around matches for better recall

Use cases:
- User: "Earlier you mentioned something about neural networks..."
- AI checks transcript for "neural networks" to provide accurate callback
- Avoids AI hallucinating past conversation content

**File to create:** `internal/podcast/search.go` - transcript search/lookup

### 6.3 Transcript Export
- End-of-session summary view
- Export to markdown/text file
- Include metadata (topic, persona, duration)

---

## Phase 7: Configuration & Polish

### 7.1 Configuration
- Environment variables for API keys
- Config file support (optional)
- Voice preferences persistence

### 7.2 Error Handling
- Graceful API error messages
- Retry logic for transient failures
- User-friendly error display

### 7.3 Styling
- Consistent Lip Gloss theme
- Color scheme for different speakers
- Responsive terminal sizing

---

## Phase 8: Testing

### 8.1 Unit Tests
- AI provider interface mocking
- Session management logic
- Transcript generation
- System prompt construction

### 8.2 Integration Tests
- Full conversation flow (with mocked AI)
- UI state transitions

---

## Key Implementation Considerations

| Concern | Approach |
|---------|----------|
| **Streaming** | Use channels + goroutines; context cancellation for interrupts |
| **UI responsiveness** | Bubble Tea's non-blocking Update/View pattern |
| **Persona adherence** | Strong system prompt + message history context |
| **Guidelines compliance** | Embed guidelines in system prompt as explicit rules |
| **Modular AI** | Interface-based design allows swapping providers |
| **Transcript recall** | Live searchable transcript; AI queries before referencing past content |

---

## Suggested Implementation Order

1. **Project scaffold** - go.mod, directory structure
2. **AI interface + OpenAI client** - get streaming working
3. **Basic Bubble Tea app** - welcome → topic → persona flow
4. **Conversation view** - text input + streaming display
5. **Typewriter animation** - polish the AI response display
6. **Voice selection + TTS** - audio playback
7. **Transcript export** - end-of-session feature
8. **Tests + error handling** - production readiness
