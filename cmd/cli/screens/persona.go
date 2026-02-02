package screens

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
)

// PersonaModel represents the persona input screen
type PersonaModel struct {
	textInput textinput.Model
	width     int
	height    int
	persona   string
}

// NewPersonaModel creates a new persona input screen model
func NewPersonaModel() PersonaModel {
	ti := textinput.New()
	ti.Placeholder = "e.g., tech entrepreneur, celebrity chef, scientist"
	ti.Focus()

	return PersonaModel{
		textInput: ti,
	}
}

// Init initializes the persona model
func (m PersonaModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the persona screen
func (m PersonaModel) Update(msg tea.Msg) (PersonaModel, tea.Cmd) {
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
				m.persona = m.textInput.Value()
				return m, func() tea.Msg { return PersonaSelectedMsg{Persona: m.persona} }
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the persona input screen
func (m PersonaModel) View() string {
	title := styles.TitleStyle.Render("Who is your AI guest?")

	description := styles.SubtitleStyle.Render(
		"Describe the persona of your AI podcast guest\n" +
			"This shapes how they respond and what expertise they bring",
	)

	examples := styles.InputPromptStyle.Render(
		"Examples: tech expert, fitness coach, history professor, startup founder",
	)

	input := m.textInput.View()

	help := styles.HelpStyle.Render("Enter to continue | Ctrl+C to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		description,
		"",
		examples,
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

// Persona returns the selected persona
func (m PersonaModel) Persona() string {
	return m.persona
}

// PersonaSelectedMsg signals that a persona has been selected
type PersonaSelectedMsg struct {
	Persona string
}
