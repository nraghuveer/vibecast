package screens

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/nraghuveer/vibecast/cmd/cli/mock"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
	"github.com/nraghuveer/vibecast/lib/audio"
	"github.com/nraghuveer/vibecast/lib/config"
	"github.com/nraghuveer/vibecast/lib/db"
	"github.com/nraghuveer/vibecast/lib/llm"
	"github.com/nraghuveer/vibecast/lib/logger"
	"github.com/nraghuveer/vibecast/lib/models"
	"github.com/nraghuveer/vibecast/lib/storage"
)

// Message represents a chat message
type Message struct {
	Content  string
	Speaker  models.SpeakerType
	Complete bool
}

// ConversationModel represents the main conversation screen
type ConversationModel struct {
	db            *db.DB
	textInput     textinput.Model
	messages      []Message
	width         int
	height        int
	title         string
	topic         string
	persona       string
	voice         mock.Voice
	provider      string
	isTyping      bool
	streamingText string
	ttsQueue      []string
	ttsInFlight   bool
	ttsLastChunk  string
	llmRawBuffer  string
	llmInSpeech   bool
	llmSpeechBuf  string
	llmClient     *llm.Client
	llmStream     <-chan llm.StreamEvent
	llmCancel     context.CancelFunc
	id            string
	dotFrame      int  // For flowing dots animation
	showDetails   bool // Toggle for showing topic/persona (Ctrl+I)
	inputMode     string
	isMuted       bool
	sttDraft      string
	logger        *logger.Logger
	toastModel    ToastModel
}

// NewConversationModelWithTitle creates a new conversation screen model with a title
func NewConversationModelWithTitle(database *db.DB, title, topic, persona string, voice mock.Voice, provider string, width, height int) ConversationModel {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()
	ti.Width = width - 4

	conversationID := uuid.New().String()

	storage.CreateConversationDir(conversationID)
	storage.CreateTranscript(conversationID)

	// Auto-generate title if not provided
	if title == "" {
		title = "Conversation"
	}

	conv := models.Conversation{
		ID:        conversationID,
		Title:     title,
		Topic:     topic,
		Persona:   persona,
		VoiceID:   voice.ID,
		VoiceName: voice.Name,
		Provider:  provider,
		CreatedAt: time.Now(),
	}
	database.CreateConversation(conv)

	return ConversationModel{
		db:          database,
		textInput:   ti,
		messages:    []Message{},
		width:       width,
		height:      height,
		title:       title,
		topic:       topic,
		persona:     persona,
		voice:       voice,
		provider:    provider,
		id:          conversationID,
		dotFrame:    0,
		showDetails: false,
		inputMode:   "text",
		isMuted:     false,
		sttDraft:    "",
		ttsQueue:    []string{},
		llmClient:   llm.New(),
		logger:      logger.GetInstance(),
		toastModel:  NewToastModel(),
	}
}

// NewConversationModelFromExisting creates a conversation screen model from an existing conversation
func NewConversationModelFromExisting(database *db.DB, conversation db.Conversation, width, height int) ConversationModel {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()
	ti.Width = width - 4

	// Load existing messages from storage
	loadedMessages, err := storage.LoadMessages(conversation.ID)
	if err != nil {
		loadedMessages = []storage.Message{}
	}

	// Convert storage messages to screen messages
	var messages []Message
	for _, msg := range loadedMessages {
		messages = append(messages, Message{
			Content:  msg.Content,
			Speaker:  msg.Speaker,
			Complete: true,
		})
	}

	return ConversationModel{
		db:          database,
		textInput:   ti,
		messages:    messages,
		width:       width,
		height:      height,
		title:       conversation.Title,
		topic:       conversation.Topic,
		persona:     conversation.Persona,
		voice:       mock.Voice{ID: conversation.VoiceID, Name: conversation.VoiceName},
		provider:    conversation.Provider,
		id:          conversation.ID,
		dotFrame:    0,
		showDetails: false,
		inputMode:   "text",
		isMuted:     false,
		sttDraft:    "",
		ttsQueue:    []string{},
		llmClient:   llm.New(),
		logger:      logger.GetInstance(),
		toastModel:  NewToastModel(),
	}
}

// Init initializes the conversation model
func (m ConversationModel) Init() tea.Cmd {
	// Check last message to determine whose turn it is
	// If no messages or last message is from Guest → trigger Guest greeting (new conversation)
	// If last message is from Host → trigger Guest response (continue conversation)
	// If last message is from Guest → wait for Host input (continue conversation, Guest already spoke)
	if len(m.messages) == 0 {
		// New conversation: Guest starts with greeting
		return tea.Batch(
			textinput.Blink,
			m.startGuestResponse(true),
		)
	}

	// Check last message
	lastMsg := m.messages[len(m.messages)-1]
	if lastMsg.Speaker == models.HOST {
		// Host spoke last, it's Guest's turn to respond
		return tea.Batch(
			textinput.Blink,
			m.startGuestResponse(false),
		)
	}

	// Guest spoke last, wait for Host input
	return textinput.Blink
}

// DotAnimationMsg is sent for flowing dots animation
type DotAnimationMsg struct{}

// ResponseCompleteMsg signals the response is done streaming
type ResponseCompleteMsg struct{}

// LLMStreamMsg delivers streamed LLM deltas to the UI.
type LLMStreamMsg struct {
	Event llm.StreamEvent
}

// TTSSavedMsg indicates synthesized audio has been saved.
type TTSSavedMsg struct {
	Filename string
	Path     string
	Err      error
}

// SttDraftMsg updates the live speech-to-text draft text.
// Future STT streaming should emit this message with partial transcripts.
type SttDraftMsg struct {
	Text string
}

func (m ConversationModel) startGuestResponse(isFirst bool) tea.Cmd {
	return func() tea.Msg {
		return StartResponseMsg{IsFirst: isFirst}
	}
}

// StartResponseMsg signals to start a guest response
type StartResponseMsg struct {
	IsFirst bool
}

// Update handles messages for the conversation screen
func (m ConversationModel) Update(msg tea.Msg) (ConversationModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = m.width - 4 // Full width minus small padding

	case StartResponseMsg:
		// Cancel any in-flight request.
		if m.llmCancel != nil {
			m.llmCancel()
			m.llmCancel = nil
		}

		m.isTyping = true
		m.streamingText = ""
		m.dotFrame = 0
		m.resetLLMParser()
		m.ttsQueue = nil
		m.ttsInFlight = false
		m.ttsLastChunk = ""

		ctx, cancel := context.WithCancel(context.Background())
		m.llmCancel = cancel

		history := m.toChatHistory(msg.IsFirst)
		stream, err := m.llmClient.StreamGuestResponse(ctx, m.provider, m.persona, m.topic, history)
		if err != nil {
			m.isTyping = false
			m.logger.LogError("llm_stream_init", err)
			m.toastModel.AddError("AI connection failed. Check your settings.")
			notice := "Sorry—I'm having trouble connecting to the AI provider right now. Give me a moment and try again."
			m.messages = append(m.messages, Message{Content: notice, Speaker: models.GUEST, Complete: true})
			_ = storage.AppendMessage(m.id, "Guest", notice)
			return m, DismissToastCmd(len(m.toastModel.GetToasts())-1, 5*time.Second)
		}

		m.llmStream = stream
		return m, tea.Batch(m.dotAnimationCmd(), m.waitLLMEventCmd())

	case LLMStreamMsg:
		if msg.Event.Err != nil {
			m.isTyping = false
			m.llmStream = nil
			if m.llmCancel != nil {
				m.llmCancel()
				m.llmCancel = nil
			}
			m.resetLLMParser()

			m.logger.LogError("llm_stream_error", msg.Event.Err)
			m.toastModel.AddError("AI stream error. Please try again.")
			notice := "Sorry—looks like I'm having trouble reaching the AI right now. Want to try that again in a second?"
			m.messages = append(m.messages, Message{Content: notice, Speaker: models.GUEST, Complete: true})
			_ = storage.AppendMessage(m.id, "Guest", notice)
			return m, DismissToastCmd(len(m.toastModel.GetToasts())-1, 5*time.Second)
		}

		if msg.Event.Delta != "" {
			speechDelta, blocks := m.consumeLLMDelta(msg.Event.Delta)
			if speechDelta != "" {
				m.streamingText += speechDelta
			}
			var ttsCmd tea.Cmd
			m, ttsCmd = m.enqueueTTSBlocks(blocks)
			return m, tea.Batch(m.waitLLMEventCmd(), ttsCmd)
		}

		if msg.Event.Done {
			m.isTyping = false
			m.llmStream = nil
			if m.llmCancel != nil {
				m.llmCancel()
				m.llmCancel = nil
			}

			finalDelta, blocks := m.finalizeLLMStream()
			if finalDelta != "" {
				m.streamingText += finalDelta
			}

			final := strings.TrimSpace(m.streamingText)
			if final == "" {
				final = "Um—I'm blanking for a second. Could you rephrase that?"
			}

			m.messages = append(m.messages, Message{Content: final, Speaker: models.GUEST, Complete: true})
			_ = storage.AppendMessage(m.id, "Guest", final)
			m.streamingText = ""
			m, ttsCmd := m.enqueueTTSBlocks(blocks)
			return m, ttsCmd
		}

		return m, m.waitLLMEventCmd()

	case DotAnimationMsg:
		if m.isTyping {
			m.dotFrame = (m.dotFrame + 1) % 20
			return m, m.dotAnimationCmd()
		}
		return m, nil

	case TTSSavedMsg:
		if msg.Err != nil || msg.Filename == "" {
			if msg.Err != nil {
				m.logger.LogError("tts_save", msg.Err)
			}
			m.ttsInFlight = false
			return m.startNextTTS()
		}
		if msg.Path != "" {
			if err := audio.Enqueue(msg.Path); err != nil {
				m.logger.LogError("audio_enqueue", err)
			}
		}
		m.ttsInFlight = false
		return m.startNextTTS()

	case ToastDismissMsg:
		m.toastModel.RemoveToast(msg.Index)
		return m, nil

	case SttDraftMsg:
		if m.inputMode == "voice" {
			m.sttDraft = msg.Text
			return m, nil
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			if m.inputMode == "text" {
				m.inputMode = "voice"
				m.textInput.Blur()
			} else {
				m.inputMode = "text"
				m.textInput.Focus()
			}
			return m, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("m"))):
			if m.inputMode == "voice" {
				m.isMuted = !m.isMuted
				return m, nil
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			m.logger.Info("conversation_quit", "conversation_id", m.id)
			m.cancelInflightLLM()
			audio.Drain()
			m.EndConversation()
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+i"))):
			m.showDetails = !m.showDetails
			return m, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("q"))):
			if !m.isTyping && m.textInput.Value() == "" {
				m.cancelInflightLLM()
				audio.Drain()
				m.EndConversation()
				return m, tea.Quit
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if !m.isTyping && m.textInput.Value() != "" {
				// Add host message
				hostMsg := m.textInput.Value()
				m.logger.Info("host_message_sent", "conversation_id", m.id, "message_length", len(hostMsg))
				m.messages = append(m.messages, Message{
					Content:  hostMsg,
					Speaker:  models.HOST,
					Complete: true,
				})
				if err := storage.AppendMessage(m.id, "Host", hostMsg); err != nil {
					m.logger.LogError("storage_append_message", err)
				}
				m.textInput.Reset()
				return m, m.startGuestResponse(false)
			}
		}
	}

	if !m.isTyping && m.inputMode == "text" {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m ConversationModel) waitLLMEventCmd() tea.Cmd {
	stream := m.llmStream
	return func() tea.Msg {
		if stream == nil {
			return LLMStreamMsg{Event: llm.StreamEvent{Done: true}}
		}
		ev, ok := <-stream
		if !ok {
			return LLMStreamMsg{Event: llm.StreamEvent{Done: true}}
		}
		return LLMStreamMsg{Event: ev}
	}
}

func (m *ConversationModel) resetLLMParser() {
	m.llmRawBuffer = ""
	m.llmInSpeech = false
	m.llmSpeechBuf = ""
}

func (m *ConversationModel) consumeLLMDelta(delta string) (string, []string) {
	m.llmRawBuffer += delta
	return m.parseLLMBuffer(false)
}

func (m *ConversationModel) finalizeLLMStream() (string, []string) {
	return m.parseLLMBuffer(true)
}

func (m *ConversationModel) parseLLMBuffer(final bool) (string, []string) {
	const startTag = "<speech>"
	const endTag = "</speech>"
	buffer := m.llmRawBuffer
	inSpeech := m.llmInSpeech
	speechBuf := m.llmSpeechBuf
	var out strings.Builder
	var blocks []string

	for {
		if inSpeech {
			idx := strings.Index(buffer, endTag)
			if idx == -1 {
				keep := partialTagSuffix(buffer, endTag)
				if len(buffer) > keep {
					segment := buffer[:len(buffer)-keep]
					out.WriteString(segment)
					speechBuf += segment
					buffer = buffer[len(buffer)-keep:]
				}
				break
			}
			segment := buffer[:idx]
			out.WriteString(segment)
			speechBuf += segment
			buffer = buffer[idx+len(endTag):]
			block := strings.TrimSpace(speechBuf)
			if block != "" {
				blocks = append(blocks, block)
			}
			speechBuf = ""
			inSpeech = false
			continue
		}

		idx := strings.Index(buffer, startTag)
		if idx == -1 {
			keep := partialTagSuffix(buffer, startTag)
			if len(buffer) > keep {
				buffer = buffer[len(buffer)-keep:]
			}
			break
		}
		buffer = buffer[idx+len(startTag):]
		inSpeech = true
	}

	if final {
		if inSpeech {
			keep := partialTagSuffix(buffer, endTag)
			content := buffer
			if keep > 0 && len(buffer) >= keep {
				content = buffer[:len(buffer)-keep]
			}
			content = strings.TrimSpace(content)
			if content != "" {
				out.WriteString(content)
				speechBuf += content
			}
			block := strings.TrimSpace(speechBuf)
			if block != "" {
				blocks = append(blocks, block)
			}
			speechBuf = ""
			inSpeech = false
		}
		buffer = ""
	}

	m.llmRawBuffer = buffer
	m.llmInSpeech = inSpeech
	m.llmSpeechBuf = speechBuf
	return out.String(), blocks
}

func partialTagSuffix(buffer string, tag string) int {
	max := len(tag) - 1
	if max > len(buffer) {
		max = len(buffer)
	}
	for i := max; i > 0; i-- {
		if strings.HasSuffix(buffer, tag[:i]) {
			return i
		}
	}
	return 0
}

func (m *ConversationModel) cancelInflightLLM() {
	if m.llmCancel != nil {
		m.llmCancel()
		m.llmCancel = nil
	}
	m.llmStream = nil
}

func (m ConversationModel) toChatHistory(isFirst bool) []llm.ChatMessage {
	if isFirst {
		return []llm.ChatMessage{{
			Role:    "user",
			Content: fmt.Sprintf("Start the episode with a brief warm greeting to the host, then invite the first question about the topic: %s.", m.topic),
		}}
	}

	history := make([]llm.ChatMessage, 0, len(m.messages))
	for _, msg := range m.messages {
		switch msg.Speaker {
		case models.HOST:
			history = append(history, llm.ChatMessage{Role: "user", Content: msg.Content})
		case models.GUEST:
			history = append(history, llm.ChatMessage{Role: "assistant", Content: msg.Content})
		}
	}
	return history
}

func (m ConversationModel) ttsCmd(text string) tea.Cmd {
	// Keep TTS best-effort; conversation should work without it.
	ttsProvider := config.GetTextToSpeechProvider()
	if strings.TrimSpace(ttsProvider) == "" {
		ttsProvider = "openai"
	}

	if url, err := config.GetProviderTTSURL(ttsProvider); err != nil || strings.TrimSpace(url) == "" {
		if _, err := config.GetProviderConfig("openai"); err == nil {
			ttsProvider = "openai"
		}
	}

	voiceID := m.voice.ID
	persona := m.persona
	topic := m.topic
	conversationID := m.id
	client := m.llmClient

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		cleanText := normalizeTTSInput(text)
		audioData, _, err := client.SynthesizeGuestSpeech(ctx, "", ttsProvider, persona, topic, voiceID, cleanText)
		if err != nil {
			return TTSSavedMsg{Err: err}
		}
		filename, err := storage.SaveAudio(conversationID, audioData)
		if err != nil {
			return TTSSavedMsg{Err: err}
		}
		audioDir, err := storage.GetAudioDir(conversationID)
		if err != nil {
			return TTSSavedMsg{Filename: filename, Err: err}
		}
		return TTSSavedMsg{Filename: filename, Path: filepath.Join(audioDir, filename)}
	}
}

func (m ConversationModel) enqueueTTSBlocks(blocks []string) (ConversationModel, tea.Cmd) {
	if len(blocks) == 0 {
		return m, nil
	}
	for _, block := range blocks {
		trimmed := strings.TrimSpace(block)
		if trimmed == "" {
			continue
		}
		m.ttsQueue = append(m.ttsQueue, trimmed)
	}
	return m.startNextTTS()
}

func normalizeTTSInput(text string) string {
	clean := stripHTMLTags(text)
	clean = ensureSentenceSpacing(clean)
	clean = strings.Join(strings.Fields(clean), " ")
	return strings.TrimSpace(clean)
}

func stripHTMLTags(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

func ensureSentenceSpacing(s string) string {
	var b strings.Builder
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		b.WriteRune(runes[i])
		if i == len(runes)-1 {
			continue
		}
		if isPunct(runes[i]) && !isSpace(runes[i+1]) {
			b.WriteRune(' ')
		}
	}
	return b.String()
}

func isPunct(r rune) bool {
	switch r {
	case '.', '!', '?', ',', ':', ';':
		return true
	default:
		return false
	}
}

func isSpace(r rune) bool {
	switch r {
	case ' ', '\n', '\t', '\r':
		return true
	default:
		return false
	}
}

func renderTranscriptMessage(labelStyle lipgloss.Style, label, content string, width int) string {
	labelWidth := lipgloss.Width(label)
	textWidth := width - labelWidth - 2
	if textWidth < 10 {
		textWidth = 10
	}

	lines := wrapText(content, textWidth)
	if len(lines) == 0 {
		lines = []string{""}
	}

	var b strings.Builder
	firstPrefix := labelStyle.Render(label) + "  "
	pad := strings.Repeat(" ", labelWidth+2)
	for i, line := range lines {
		prefix := pad
		if i == 0 {
			prefix = firstPrefix
		}
		b.WriteString(prefix)
		b.WriteString(styles.TranscriptTextStyle.Render(line))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}

func renderTranscriptStreaming(labelStyle lipgloss.Style, label, content string, width int) string {
	labelWidth := lipgloss.Width(label)
	textWidth := width - labelWidth - 2
	if textWidth < 10 {
		textWidth = 10
	}

	lines := wrapText(content, textWidth)
	if len(lines) == 0 {
		lines = []string{""}
	}

	var b strings.Builder
	firstPrefix := labelStyle.Render(label) + "  "
	pad := strings.Repeat(" ", labelWidth+2)
	for i, line := range lines {
		prefix := pad
		if i == 0 {
			prefix = firstPrefix
		}
		b.WriteString(prefix)
		b.WriteString(styles.TranscriptTextStyle.Render(line))
		b.WriteString("\n")
	}
	return b.String()
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	var lines []string
	paragraphs := strings.Split(text, "\n")
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			lines = append(lines, "")
			continue
		}
		words := strings.Fields(p)
		var line string
		for _, word := range words {
			if line == "" {
				line = word
				continue
			}
			candidate := line + " " + word
			if lipgloss.Width(candidate) <= width {
				line = candidate
				continue
			}
			lines = append(lines, line)
			line = word
		}
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func (m ConversationModel) startNextTTS() (ConversationModel, tea.Cmd) {
	if m.ttsInFlight || len(m.ttsQueue) == 0 {
		return m, nil
	}
	chunk := strings.TrimSpace(m.ttsQueue[0])
	m.ttsQueue = m.ttsQueue[1:]
	if chunk == "" {
		return m.startNextTTS()
	}
	if chunk == m.ttsLastChunk {
		return m.startNextTTS()
	}
	m.ttsInFlight = true
	m.ttsLastChunk = chunk
	return m, m.ttsCmd(chunk)
}

func (m ConversationModel) dotAnimationCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return DotAnimationMsg{}
	})
}

// View renders the conversation screen
func (m ConversationModel) View() string {
	mainContent := m.renderConversation()

	// Render toasts in upper left if any
	if m.toastModel.HasToasts() {
		toasts := RenderToasts(m.toastModel.GetToasts())
		if toasts != "" {
			// Position toasts at the top with proper spacing
			toastContainer := lipgloss.NewStyle().
				Padding(1, 2).
				Render(toasts)

			// Overlay toasts on top of main content
			return lipgloss.JoinVertical(
				lipgloss.Left,
				toastContainer,
				mainContent,
			)
		}
	}

	return mainContent
}

// Pre-rendered animation frames
var animationFrames []string

func init() {
	animationFrames = preRenderAnimationFrames()
}

func preRenderAnimationFrames() []string {
	trackWidth := 10
	maxHeadPos := trackWidth - 3  // head can go from 0 to 7 (leaves room for 3-block snake)
	cycleLength := maxHeadPos * 2 // 14 frames: 0->7 then 7->0

	// Styles for rendering
	mutedDotStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3f3f46"))
	primaryStyle := lipgloss.NewStyle().Foreground(styles.PrimaryColor).Bold(true)
	lightGreyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#71717a"))
	veryLightGreyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#52525b"))

	frames := make([]string, cycleLength)

	for frame := 0; frame < cycleLength; frame++ {
		var headPos int
		var movingRight bool

		if frame < maxHeadPos {
			// Moving right: 0, 1, 2, ... maxHeadPos-1
			headPos = frame
			movingRight = true
		} else {
			// Moving left: maxHeadPos-1, maxHeadPos-2, ... 0
			headPos = cycleLength - frame - 1
			movingRight = false
		}

		var result []string
		for i := 0; i < trackWidth; i++ {
			if movingRight {
				// Head faces right: [tail2][tail1][head] →
				if i == headPos+2 {
					result = append(result, primaryStyle.Render("█"))
				} else if i == headPos+1 {
					result = append(result, lightGreyStyle.Render("█"))
				} else if i == headPos {
					result = append(result, veryLightGreyStyle.Render("█"))
				} else {
					result = append(result, mutedDotStyle.Render("·"))
				}
			} else {
				// Head faces left: ← [head][tail1][tail2]
				if i == headPos {
					result = append(result, primaryStyle.Render("█"))
				} else if i == headPos+1 {
					result = append(result, lightGreyStyle.Render("█"))
				} else if i == headPos+2 {
					result = append(result, veryLightGreyStyle.Render("█"))
				} else {
					result = append(result, mutedDotStyle.Render("·"))
				}
			}
		}
		frames[frame] = strings.Join(result, "")
	}

	return frames
}

// renderFlowingDots returns pre-rendered animation frame
func (m ConversationModel) renderFlowingDots() string {
	return animationFrames[m.dotFrame%len(animationFrames)]
}

func (m ConversationModel) renderConversation() string {
	// Fixed heights for bottom elements
	logoHeight := 3
	inputHeight := 1
	inputMetaHeight := 1
	animationHeight := 1
	helpHeight := 1
	padding := 4 // spacing between elements

	// Calculate transcript height to fill remaining space
	bottomHeight := inputHeight + inputMetaHeight + helpHeight + padding
	if m.isTyping {
		bottomHeight += animationHeight + 1
	}
	transcriptHeight := m.height - logoHeight - bottomHeight - 2

	// Adjust height if showing details
	detailsHeight := 0
	if m.showDetails {
		detailsHeight = 2 // Topic and Persona lines
		transcriptHeight -= detailsHeight
	}

	// Transcript area with HOST and GUEST labels
	var transcriptView strings.Builder

	// Render completed messages with HOST/GUEST labels
	contentWidth := m.width - 4
	if contentWidth < 20 {
		contentWidth = 20
	}
	for _, msg := range m.messages {
		if msg.Speaker == models.HOST {
			transcriptView.WriteString(renderTranscriptMessage(styles.HostLabelStyle, "HOST", msg.Content, contentWidth))
		} else {
			transcriptView.WriteString(renderTranscriptMessage(styles.GuestLabelStyle, "GUEST", msg.Content, contentWidth))
		}
	}

	// Add streaming message if typing
	if m.isTyping && m.streamingText != "" {
		transcriptView.WriteString(renderTranscriptStreaming(styles.GuestLabelStyle, "GUEST", m.streamingText+"▌", contentWidth))
	}

	// Create transcript container
	transcriptContent := transcriptView.String()
	transcriptContainer := styles.TranscriptPanelStyle.
		Height(transcriptHeight).
		Render(transcriptContent)

	// Input area (left column)
	var inputArea string
	if m.inputMode == "voice" {
		muteLabel := "[ON]"
		if m.isMuted {
			muteLabel = "[MUTE]"
		}
		muteStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)
		placeholder := m.sttDraft
		if strings.TrimSpace(placeholder) == "" {
			placeholder = "Listening for speech..."
		}
		voiceBlockStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
		inputArea = fmt.Sprintf("  %s %s", muteStyle.Render(muteLabel), voiceBlockStyle.Render(placeholder))
	} else {
		inputArea = "  " + m.textInput.View()
	}

	modeStyle := lipgloss.NewStyle().Foreground(styles.PrimaryColor).Bold(true)
	providerStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)
	modeMeta := fmt.Sprintf("%s %s", modeStyle.Render(strings.ToUpper(m.inputMode)), providerStyle.Render(m.provider))
	metaLine := fmt.Sprintf("  %s", modeMeta)

	// Help text
	help := styles.HelpStyle.Render("  Tab toggle input | m mute | Enter to send | Ctrl+I show/hide details | q or Ctrl+C to exit")

	// Build bottom section: input first, then animation below (both anchored to bottom)
	var bottomSection string
	inputLine := inputArea
	inputColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		inputLine,
		metaLine,
	)
	if m.isTyping {
		dotsAnimation := m.renderFlowingDots()
		animStyle := lipgloss.NewStyle().PaddingLeft(2)
		animationArea := animStyle.Render(dotsAnimation)

		bottomSection = lipgloss.JoinVertical(
			lipgloss.Left,
			inputColumn,
			animationArea,
			help,
		)
	} else {
		bottomSection = lipgloss.JoinVertical(
			lipgloss.Left,
			inputColumn,
			help,
		)
	}

	// Build the top section (logo with optional details)
	var topSection string
	if m.showDetails {
		mutedStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)
		textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
		details := lipgloss.JoinVertical(
			lipgloss.Left,
			fmt.Sprintf("%s %s", mutedStyle.Render("Topic:"), textStyle.Render(m.topic)),
			fmt.Sprintf("%s %s", mutedStyle.Render("Persona:"), textStyle.Render(m.persona)),
		)
		topSection = lipgloss.JoinVertical(
			lipgloss.Left,
			styles.LogoWithTitle(m.title),
			"",
			details,
		)
	} else {
		topSection = styles.LogoWithTitle(m.title)
	}

	// Combine: logo (with optional details), transcript, then bottom section anchored to bottom
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		topSection,
		"",
		transcriptContainer,
		"",
		bottomSection,
	)

	return content
}

// Messages returns all conversation messages
func (m ConversationModel) Messages() []Message {
	return m.messages
}

// ID returns the conversation ID
func (m ConversationModel) ID() string {
	return m.id
}

// EndConversation marks the conversation as ended
func (m ConversationModel) EndConversation() error {
	m.logger.Info("conversation_ended", "conversation_id", m.id)
	now := time.Now()
	err := m.db.UpdateConversationEndedAt(m.id, now)
	if err != nil {
		m.logger.LogError("conversation_end", err)
	}
	return err
}

// ExitConversationMsg signals to exit the conversation
type ExitConversationMsg struct{}
