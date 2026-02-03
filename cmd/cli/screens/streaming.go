package screens

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/nraghuveer/vibecast/lib/models"
)

type StreamingModel struct {
	isTyping      bool
	streamingText string
	fullResponse  string
	streamIndex   int
}

// TickMsg is kept for legacy typing animation.
// Some screens may still use it for fake streaming.
type TickMsg struct{}

func NewStreamingModel() StreamingModel {
	return StreamingModel{}
}

func (m *StreamingModel) StartResponse(response string) {
	m.isTyping = true
	m.fullResponse = response
	m.streamingText = ""
	m.streamIndex = 0
}

func (m *StreamingModel) Advance() bool {
	if m.streamIndex < len(m.fullResponse) {
		m.streamingText += string(m.fullResponse[m.streamIndex])
		m.streamIndex++
		return true
	}
	m.isTyping = false
	return false
}

func (m *StreamingModel) Complete() Message {
	msg := Message{
		Content:  m.fullResponse,
		Speaker:  models.GUEST,
		Complete: true,
	}
	m.streamingText = ""
	m.fullResponse = ""
	return msg
}

func (m *StreamingModel) IsTyping() bool {
	return m.isTyping
}

func (m *StreamingModel) GetStreamingText() string {
	return m.streamingText
}

func (m *StreamingModel) TickCmd() func() tea.Msg {
	return func() tea.Msg {
		return TickMsg{}
	}
}
