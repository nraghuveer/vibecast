package screens

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
)

// TopicModel represents the topic input screen
type TopicModel struct {
	textInput textinput.Model
	width     int
	height    int
	topic     string
}

// NewTopicModel creates a new topic input screen model
func NewTopicModel() TopicModel {
	ti := textinput.New()
	ti.Placeholder = "e.g., The future of artificial intelligence"
	ti.Focus()

	return TopicModel{
		textInput: ti,
	}
}

// Init initializes the topic model
func (m TopicModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the topic screen
func (m TopicModel) Update(msg tea.Msg) (TopicModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = m.width - 20
		if m.textInput.Width > 100 {
			m.textInput.Width = 100
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.textInput.Value() != "" {
				m.topic = m.textInput.Value()
				return m, func() tea.Msg { return TopicSelectedMsg{Topic: m.topic} }
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the topic input screen
func (m TopicModel) View() string {
	title := styles.TitleStyle.Render("What's your podcast topic?")

	description := styles.SubtitleStyle.Render(
		"Enter the main topic or theme for your podcast episode",
	)

	input := m.textInput.View()

	help := styles.HelpStyle.Render("Enter to continue | Ctrl+C to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		description,
		"",
		input,
		"",
		help,
	)

	box := styles.BoxStyle.Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

// Topic returns the selected topic
func (m TopicModel) Topic() string {
	return m.topic
}

// TopicSelectedMsg signals that a topic has been selected
type TopicSelectedMsg struct {
	Topic string
}
