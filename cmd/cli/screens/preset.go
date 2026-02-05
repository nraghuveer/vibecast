package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
	"github.com/nraghuveer/vibecast/lib/data"
	"github.com/nraghuveer/vibecast/lib/db"
	"github.com/nraghuveer/vibecast/lib/logger"
	"github.com/nraghuveer/vibecast/lib/models"
)

// PresetModel represents the preset selection screen with title input
type PresetModel struct {
	db         *db.DB
	templates  []models.Template
	cursor     int
	titleInput textinput.Model
	showError  bool
	selected   models.Template
	width      int
	height     int
	logger     *logger.Logger
}

// NewPresetModel creates a new preset selection screen model
func NewPresetModel(database *db.DB) PresetModel {
	templates := data.GetTemplates(database)

	ti := textinput.New()
	ti.Placeholder = "e.g., AI Future Discussion, Tech Deep Dive..."
	ti.CharLimit = 100
	ti.Focus()

	return PresetModel{
		db:         database,
		templates:  templates,
		cursor:     0,
		titleInput: ti,
		showError:  false,
		logger:     logger.GetInstance(),
	}
}

// Init initializes the preset model
func (m PresetModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the preset screen
func (m PresetModel) Update(msg tea.Msg) (PresetModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.titleInput.Width = m.width - 20
		if m.titleInput.Width > 100 {
			m.titleInput.Width = 100
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
			m.showError = false
		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(m.templates)-1 {
				m.cursor++
			}
			m.showError = false
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.titleInput.Value() == "" {
				m.showError = true
				m.logger.Warn("preset_validation_failed", "reason", "empty_title")
				return m, nil
			}
			m.selected = m.templates[m.cursor]
			m.logger.Info("preset_selected",
				"template_name", m.selected.Name,
				"title", m.titleInput.Value(),
			)
			return m, func() tea.Msg {
				return PresetSelectedMsg{
					Template: m.selected,
					Title:    m.titleInput.Value(),
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			m.logger.Info("preset_back_to_welcome")
			return m, func() tea.Msg { return BackToWelcomeMsg{} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			m.logger.Info("preset_quit")
			return m, tea.Quit
		}
	}

	m.titleInput, cmd = m.titleInput.Update(msg)
	return m, cmd
}

// View renders the preset selection screen
func (m PresetModel) View() string {
	title := styles.TitleStyle.Render("Quick Start - Select a Template")

	description := styles.SubtitleStyle.Render(
		"Enter a title and choose a predefined topic and persona to get started quickly",
	)

	// Title input section
	titleLabel := lipgloss.NewStyle().Bold(true).Foreground(styles.PrimaryColor).Render("Conversation Title:")
	titleInputBox := m.titleInput.View()

	// Error message
	var errorMsg string
	if m.showError {
		errorMsg = lipgloss.NewStyle().Bold(true).Foreground(styles.ErrorColor).Render("⚠ Title is required - please enter a conversation title")
	}

	// Templates section
	templatesLabel := lipgloss.NewStyle().Bold(true).Foreground(styles.PrimaryColor).Render("Select Template:")

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
		titleLabel,
		titleInputBox,
	)

	if m.showError {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			errorMsg,
		)
	}

	content = lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"",
		templatesLabel,
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
func (m PresetModel) SelectedTemplate() models.Template {
	return m.selected
}

// Title returns the entered title
func (m PresetModel) Title() string {
	return m.titleInput.Value()
}

// PresetSelectedMsg signals that a preset has been selected
type PresetSelectedMsg struct {
	Template models.Template
	Title    string
}

// BackToWelcomeMsg signals to go back to the welcome screen
type BackToWelcomeMsg struct{}
