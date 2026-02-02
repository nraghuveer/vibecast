package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/mock"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
)

// VoiceModel represents the voice selection screen
type VoiceModel struct {
	voices   []mock.Voice
	cursor   int
	selected mock.Voice
	width    int
	height   int
}

// NewVoiceModel creates a new voice selection screen model
func NewVoiceModel() VoiceModel {
	voices := mock.GetVoices()
	return VoiceModel{
		voices: voices,
		cursor: 0,
	}
}

// Init initializes the voice model
func (m VoiceModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the voice screen
func (m VoiceModel) Update(msg tea.Msg) (VoiceModel, tea.Cmd) {
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
			if m.cursor < len(m.voices)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			m.selected = m.voices[m.cursor]
			return m, func() tea.Msg { return VoiceSelectedMsg{Voice: m.selected} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the voice selection screen
func (m VoiceModel) View() string {
	title := styles.TitleStyle.Render("Select a voice for your AI guest")

	description := styles.SubtitleStyle.Render(
		"Choose the voice that best fits your guest's persona",
	)

	var items string
	for i, voice := range m.voices {
		cursor := "  "
		itemStyle := styles.NormalStyle
		if i == m.cursor {
			cursor = "> "
			itemStyle = styles.SelectedStyle
		}

		voiceDesc := styles.VoiceDescStyle.Render(fmt.Sprintf("(%s)", voice.Description))
		item := fmt.Sprintf("%s%s %s", cursor, itemStyle.Render(voice.Name), voiceDesc)
		items += item + "\n"
	}

	help := styles.HelpStyle.Render("↑/↓ or j/k to navigate | Enter to select | Ctrl+C to quit")

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

// SelectedVoice returns the selected voice
func (m VoiceModel) SelectedVoice() mock.Voice {
	return m.selected
}

// VoiceSelectedMsg signals that a voice has been selected
type VoiceSelectedMsg struct {
	Voice mock.Voice
}
