package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/mock"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
)

// PresetModel represents the preset selection screen
type PresetModel struct {
	templates []mock.Template
	cursor    int
	selected  mock.Template
	width     int
	height    int
}

// NewPresetModel creates a new preset selection screen model
func NewPresetModel() PresetModel {
	templates := mock.GetTemplates()
	return PresetModel{
		templates: templates,
		cursor:    0,
	}
}

// Init initializes the preset model
func (m PresetModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the preset screen
func (m PresetModel) Update(msg tea.Msg) (PresetModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(m.templates)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			m.selected = m.templates[m.cursor]
			return m, func() tea.Msg { return PresetSelectedMsg{Template: m.selected} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			return m, func() tea.Msg { return BackToWelcomeMsg{} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the preset selection screen
func (m PresetModel) View() string {
	title := styles.TitleStyle.Render("Quick Start - Select a Template")

	description := styles.SubtitleStyle.Render(
		"Choose a predefined topic and persona to get started quickly",
	)

	var items string
	for i, tmpl := range m.templates {
		cursor := "  "
		nameStyle := styles.NormalStyle
		if i == m.cursor {
			cursor = "> "
			nameStyle = styles.SelectedStyle
		}

		topicPreview := tmpl.Topic
		if len(topicPreview) > 40 {
			topicPreview = topicPreview[:37] + "..."
		}

		desc := styles.VoiceDescStyle.Render(fmt.Sprintf("(%s)", topicPreview))
		item := fmt.Sprintf("%s%s %s", cursor, nameStyle.Render(tmpl.Name), desc)
		items += item + "\n"
	}

	help := styles.HelpStyle.Render("↑/↓ or j/k to navigate • Enter to select • Esc to go back • Ctrl+C to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		description,
		"",
		items,
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

// SelectedTemplate returns the selected template
func (m PresetModel) SelectedTemplate() mock.Template {
	return m.selected
}

// PresetSelectedMsg signals that a preset has been selected
type PresetSelectedMsg struct {
	Template mock.Template
}

// BackToWelcomeMsg signals to go back to the welcome screen
type BackToWelcomeMsg struct{}
