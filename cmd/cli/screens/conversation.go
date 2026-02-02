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
	IsHost   bool
	Complete bool
}

// ConversationModel represents the main conversation screen
type ConversationModel struct {
	textInput     textinput.Model
	messages      []Message
	width         int
	height        int
	topic         string
	persona       string
	voice         mock.Voice
	provider      string
	isTyping      bool
	streamingText string
	fullResponse  string
	streamIndex   int
	id            string
	dotFrame      int // For flowing dots animation
}

// NewConversationModel creates a new conversation screen model
func NewConversationModel(topic, persona string, voice mock.Voice, provider string, width, height int) ConversationModel {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()
	ti.Width = width - 4

	conversationID := uuid.New().String()

	storage.CreateConversationDir(conversationID)
	storage.CreateTranscript(conversationID)

	conv := models.Conversation{
		ID:        conversationID,
		Topic:     topic,
		Persona:   persona,
		VoiceID:   voice.ID,
		VoiceName: voice.Name,
		Provider:  provider,
		CreatedAt: time.Now(),
	}
	db.CreateConversation(conv)

	return ConversationModel{
		textInput: ti,
		messages:  []Message{},
		width:     width,
		height:    height,
		topic:     topic,
		persona:   persona,
		voice:     voice,
		provider:  provider,
		id:        conversationID,
		dotFrame:  0,
	}
}

// Init initializes the conversation model
func (m ConversationModel) Init() tea.Cmd {
	// Start with an initial greeting from the AI guest
	return tea.Batch(
		textinput.Blink,
		m.startGuestResponse(true),
	)
}

// TickMsg is sent for streaming animation
type TickMsg struct{}

// DotAnimationMsg is sent for flowing dots animation
type DotAnimationMsg struct{}

// ResponseCompleteMsg signals the response is done streaming
type ResponseCompleteMsg struct{}

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
			IsHost:   false,
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

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			m.EndConversation()
			return m, tea.Quit
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
					IsHost:   true,
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

	if !m.isTyping {
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
	maxHeadPos := trackWidth - 3 // head can go from 0 to 7 (leaves room for 3-block snake)
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
	animationHeight := 1
	helpHeight := 1
	padding := 4 // spacing between elements

	// Calculate transcript height to fill remaining space
	bottomHeight := inputHeight + helpHeight + padding
	if m.isTyping {
		bottomHeight += animationHeight + 1
	}
	transcriptHeight := m.height - logoHeight - bottomHeight - 2

	// Transcript area with HOST and GUEST labels
	var transcriptView strings.Builder

	// Render completed messages with HOST/GUEST labels
	for _, msg := range m.messages {
		if msg.IsHost {
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

	// Input area (full width)
	inputArea := "  " + m.textInput.View()

	// Help text
	help := styles.HelpStyle.Render("  Enter to send | q or Ctrl+C to exit")

	// Build bottom section: input first, then animation below (both anchored to bottom)
	var bottomSection string
	if m.isTyping {
		dotsAnimation := m.renderFlowingDots()
		animStyle := lipgloss.NewStyle().PaddingLeft(2)
		animationArea := animStyle.Render(dotsAnimation)

		bottomSection = lipgloss.JoinVertical(
			lipgloss.Left,
			inputArea,
			animationArea,
			help,
		)
	} else {
		bottomSection = lipgloss.JoinVertical(
			lipgloss.Left,
			inputArea,
			help,
		)
	}

	// Combine: logo, transcript, then bottom section anchored to bottom
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.Logo(),
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
	return db.UpdateConversationEndedAt(m.id, now)
}

// ExitConversationMsg signals to exit the conversation
type ExitConversationMsg struct{}
