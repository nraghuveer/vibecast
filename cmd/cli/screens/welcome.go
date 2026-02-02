package screens

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
)

// WelcomeOption represents a menu option on the welcome screen
type WelcomeOption int

const (
	OptionNewConversation WelcomeOption = iota
	OptionContinueConversation
	OptionQuickStart
	OptionCreateTemplate
)

// WelcomeModel represents the welcome screen
type WelcomeModel struct {
	width   int
	height  int
	cursor  int
	options []string
}

// NewWelcomeModel creates a new welcome screen model
func NewWelcomeModel() WelcomeModel {
	return WelcomeModel{
		cursor: 0,
		options: []string{
			"Create new conversation",
			"Continue conversation",
			"Quick start with preset",
			"Create new template",
		},
	}
}

// Init initializes the welcome model
func (m WelcomeModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the welcome screen
func (m WelcomeModel) Update(msg tea.Msg) (WelcomeModel, tea.Cmd) {
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
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			return m, func() tea.Msg {
				return WelcomeOptionSelectedMsg{Option: WelcomeOption(m.cursor)}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the welcome screen
func (m WelcomeModel) View() string {
	logo := styles.Logo()

	description := styles.SubtitleStyle.Render(
		"Your AI-powered podcast companion\n" +
			"Create dynamic conversations with AI personas",
	)

	// Render options with cursor
	var options string
	for i, opt := range m.options {
		cursor := "  "
		itemStyle := styles.NormalStyle
		if i == m.cursor {
			cursor = "> "
			itemStyle = styles.SelectedStyle
		}
		options += cursor + itemStyle.Render(opt) + "\n"
	}

	help := styles.HelpStyle.Render("↑/↓ to navigate • Enter to select • q to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		description,
		"",
		options,
		"",
		help,
	)

	// Center the content
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// NextScreenMsg signals transition to the next screen
type NextScreenMsg struct{}

// WelcomeOptionSelectedMsg signals which option was selected on the welcome screen
type WelcomeOptionSelectedMsg struct {
	Option WelcomeOption
}
