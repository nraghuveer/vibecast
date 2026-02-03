package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/nraghuveer/vibecast/cmd/cli/mock"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
	"github.com/nraghuveer/vibecast/lib/db"
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
	fullResponse  string
	streamIndex   int
	id            string
	dotFrame      int  // For flowing dots animation
	showDetails   bool // Toggle for showing topic/persona (Ctrl+I)
	inputMode     string
	isMuted       bool
	sttDraft      string
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

// TickMsg is sent for streaming animation
type TickMsg struct{}

// DotAnimationMsg is sent for flowing dots animation
type DotAnimationMsg struct{}

// ResponseCompleteMsg signals the response is done streaming
type ResponseCompleteMsg struct{}

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
		m.isTyping = true
		m.fullResponse = mock.GetResponse("", msg.IsFirst)
		m.streamingText = ""
		m.streamIndex = 0
		m.dotFrame = 0
		return m, tea.Batch(m.tickCmd(), m.dotAnimationCmd())

	case TickMsg:
		if m.streamIndex < len(m.fullResponse) {
			m.streamingText += string(m.fullResponse[m.streamIndex])
			m.streamIndex++
			return m, m.tickCmd()
		}
		// Response complete
		m.isTyping = false
		m.messages = append(m.messages, Message{
			Content:  m.fullResponse,
			Speaker:  models.GUEST,
			Complete: true,
		})
		storage.AppendMessage(m.id, "Guest", m.fullResponse)
		m.streamingText = ""
		m.fullResponse = ""
		return m, nil

	case DotAnimationMsg:
		if m.isTyping {
			m.dotFrame = (m.dotFrame + 1) % 20
			return m, m.dotAnimationCmd()
		}
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
			m.EndConversation()
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+i"))):
			m.showDetails = !m.showDetails
			return m, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("q"))):
			if !m.isTyping && m.textInput.Value() == "" {
				m.EndConversation()
				return m, tea.Quit
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if !m.isTyping && m.textInput.Value() != "" {
				// Add host message
				hostMsg := m.textInput.Value()
				m.messages = append(m.messages, Message{
					Content:  hostMsg,
					Speaker:  models.HOST,
					Complete: true,
				})
				storage.AppendMessage(m.id, "Host", hostMsg)
				m.textInput.Reset()

				// Start guest response
				m.isTyping = true
				m.fullResponse = mock.GetResponse(hostMsg, false)
				m.streamingText = ""
				m.streamIndex = 0
				m.dotFrame = 0
				return m, tea.Batch(m.tickCmd(), m.dotAnimationCmd())
			}
		}
	}

	if !m.isTyping && m.inputMode == "text" {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m ConversationModel) tickCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m ConversationModel) dotAnimationCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return DotAnimationMsg{}
	})
}

// View renders the conversation screen
func (m ConversationModel) View() string {
	return m.renderConversation()
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
	for _, msg := range m.messages {
		if msg.Speaker == models.HOST {
			label := styles.HostLabelStyle.Render("HOST")
			content := styles.TranscriptTextStyle.Render(msg.Content)
			transcriptView.WriteString(fmt.Sprintf("%s  %s\n\n", label, content))
		} else {
			label := styles.GuestLabelStyle.Render("GUEST")
			content := styles.TranscriptTextStyle.Render(msg.Content)
			transcriptView.WriteString(fmt.Sprintf("%s %s\n\n", label, content))
		}
	}

	// Add streaming message if typing
	if m.isTyping && m.streamingText != "" {
		label := styles.GuestLabelStyle.Render("GUEST")
		content := styles.TranscriptTextStyle.Render(m.streamingText + "▌")
		transcriptView.WriteString(fmt.Sprintf("%s %s\n", label, content))
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
	now := time.Now()
	return m.db.UpdateConversationEndedAt(m.id, now)
}

// ExitConversationMsg signals to exit the conversation
type ExitConversationMsg struct{}
