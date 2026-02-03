package screens

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
	"github.com/nraghuveer/vibecast/lib/models"
)

type TranscriptModel struct {
	messages       []Message
	showTranscript bool
}

func NewTranscriptModel() TranscriptModel {
	return TranscriptModel{
		messages: []Message{},
	}
}

func (m *TranscriptModel) AddMessage(msg Message) {
	m.messages = append(m.messages, msg)
}

func (m *TranscriptModel) GetMessages() []Message {
	return m.messages
}

func (m *TranscriptModel) SetShowTranscript(show bool) {
	m.showTranscript = show
}

func (m *TranscriptModel) ShowTranscript() bool {
	return m.showTranscript
}

func (m *TranscriptModel) RenderPanel(width, height int, isTyping bool, streamingText string) string {
	var transcript strings.Builder
	for _, msg := range m.messages {
		if msg.Speaker == models.HOST {
			transcript.WriteString(styles.HostLabelStyle.Render("HOST: "))
			transcript.WriteString(styles.TranscriptTextStyle.Render(msg.Content))
		} else {
			transcript.WriteString(styles.GuestLabelStyle.Render("GUEST: "))
			transcript.WriteString(styles.TranscriptTextStyle.Render(msg.Content))
		}
		transcript.WriteString("\n\n")
	}

	if isTyping && streamingText != "" {
		transcript.WriteString(styles.GuestLabelStyle.Render("GUEST: "))
		transcript.WriteString(styles.TranscriptTextStyle.Render(streamingText + "▌"))
	} else if isTyping {
		transcript.WriteString(styles.ThinkingStyle.Render("●●●"))
	}

	content := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(transcript.String())

	title := styles.TitleStyle.Render("Transcript")

	return styles.TranscriptPanelStyle.Width(width + 2).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			content,
		),
	)
}
