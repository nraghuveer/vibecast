package screens

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
)

// TemplateNameModel represents the template name input screen
type TemplateNameModel struct {
	textInput textinput.Model
	width     int
	height    int
}

// NewTemplateNameModel creates a new template name input screen model
func NewTemplateNameModel() TemplateNameModel {
	ti := textinput.New()
	ti.Placeholder = "e.g., Tech Interview, Cooking Show..."
	ti.Focus()

	return TemplateNameModel{
		textInput: ti,
	}
}

// Init initializes the template name model
func (m TemplateNameModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the template name screen
func (m TemplateNameModel) Update(msg tea.Msg) (TemplateNameModel, tea.Cmd) {
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
				return m, func() tea.Msg {
					return TemplateNameEnteredMsg{Name: m.textInput.Value()}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			return m, func() tea.Msg { return BackToWelcomeMsg{} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the template name input screen
func (m TemplateNameModel) View() string {
	title := styles.TitleStyle.Render("Create New Template")

	description := styles.SubtitleStyle.Render(
		"Give your template a memorable name",
	)

	input := m.textInput.View()

	help := styles.HelpStyle.Render("Enter to continue • Esc to go back • Ctrl+C to quit")

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

// TemplateNameEnteredMsg signals that the template name has been entered
type TemplateNameEnteredMsg struct {
	Name string
}
