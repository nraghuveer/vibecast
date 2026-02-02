package screens

import (
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
	"github.com/nraghuveer/vibecast/lib/config"
)

type NewConversationField int

const (
	FieldTitle NewConversationField = iota
	FieldTopic
	FieldPersona
	FieldProvider
)

type CreateConversationModel struct {
	titleInput   textinput.Model
	topicInput   textinput.Model
	personaInput textinput.Model
	providers    []ProviderInfo
	providerIdx  int
	activeField  NewConversationField
	startedAt    time.Time
	width        int
	height       int
}

// NewConversationCreatedMsg is sent when a new conversation is ready to start
type NewConversationCreatedMsg struct {
	Title    string
	Topic    string
	Persona  string
	Provider string
}

func NewCreateConversationModel() CreateConversationModel {
	// Title input
	ti := textinput.New()
	ti.Placeholder = "e.g., AI Future Discussion, Tech Deep Dive..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	// Topic input
	topic := textinput.New()
	topic.Placeholder = "e.g., The future of artificial intelligence"
	topic.CharLimit = 200
	topic.Width = 50

	// Persona input
	persona := textinput.New()
	persona.Placeholder = "e.g., tech entrepreneur, scientist, chef"
	persona.CharLimit = 100
	persona.Width = 50

	// Get providers from config
	providers := getAvailableProviders()

	return CreateConversationModel{
		titleInput:   ti,
		topicInput:   topic,
		personaInput: persona,
		providers:    providers,
		providerIdx:  0,
		activeField:  FieldTitle,
		startedAt:    time.Now(),
	}
}

func (m CreateConversationModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CreateConversationModel) Update(msg tea.Msg) (CreateConversationModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			return m, func() tea.Msg { return BackToWelcomeMsg{} }

		case key.Matches(msg, key.NewBinding(key.WithKeys("tab", "down"))):
			m = m.nextField()
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab", "up"))):
			m = m.prevField()
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("left"))):
			if m.activeField == FieldProvider && m.providerIdx > 0 {
				m.providerIdx--
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("right"))):
			if m.activeField == FieldProvider && m.providerIdx < len(m.providers)-1 {
				m.providerIdx++
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.activeField == FieldProvider {
				// Submit form if all fields are filled
				if m.isValid() {
					provider := ""
					if len(m.providers) > 0 {
						provider = m.providers[m.providerIdx].Name
					}
					return m, func() tea.Msg {
						return NewConversationCreatedMsg{
							Title:    m.titleInput.Value(),
							Topic:    m.topicInput.Value(),
							Persona:  m.personaInput.Value(),
							Provider: provider,
						}
					}
				}
			} else {
				// Move to next field on Enter
				m = m.nextField()
			}
			return m, nil
		}
	}

	// Update active text input
	switch m.activeField {
	case FieldTitle:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case FieldTopic:
		m.topicInput, cmd = m.topicInput.Update(msg)
	case FieldPersona:
		m.personaInput, cmd = m.personaInput.Update(msg)
	}

	return m, cmd
}

func (m CreateConversationModel) nextField() CreateConversationModel {
	m.titleInput.Blur()
	m.topicInput.Blur()
	m.personaInput.Blur()

	m.activeField = (m.activeField + 1) % 4

	switch m.activeField {
	case FieldTitle:
		m.titleInput.Focus()
	case FieldTopic:
		m.topicInput.Focus()
	case FieldPersona:
		m.personaInput.Focus()
	}

	return m
}

func (m CreateConversationModel) prevField() CreateConversationModel {
	m.titleInput.Blur()
	m.topicInput.Blur()
	m.personaInput.Blur()

	if m.activeField == 0 {
		m.activeField = 3
	} else {
		m.activeField--
	}

	switch m.activeField {
	case FieldTitle:
		m.titleInput.Focus()
	case FieldTopic:
		m.topicInput.Focus()
	case FieldPersona:
		m.personaInput.Focus()
	}

	return m
}

func (m CreateConversationModel) isValid() bool {
	return m.titleInput.Value() != "" &&
		m.topicInput.Value() != "" &&
		m.personaInput.Value() != "" &&
		len(m.providers) > 0
}

func (m CreateConversationModel) View() string {
	// Styles
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.PrimaryColor)
	activeStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.PrimaryColor)
	mutedStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)
	selectedProviderStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.PrimaryColor).Background(lipgloss.Color("#2d2d2d")).Padding(0, 1)
	normalProviderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Padding(0, 1)

	title := styles.TitleStyle.Render("Create New Conversation")
	subtitle := styles.SubtitleStyle.Render("Fill in the details below to start your podcast conversation")

	// Started timestamp (auto-detected)
	timestampLabel := mutedStyle.Render("Started:")
	timestamp := mutedStyle.Render(m.startedAt.Format("Jan 02, 2006 3:04 PM"))

	// Title field
	titleLabel := "  Title"
	if m.activeField == FieldTitle {
		titleLabel = activeStyle.Render("> Title")
	} else {
		titleLabel = labelStyle.Render("  Title")
	}
	titleField := fmt.Sprintf("%s\n  %s", titleLabel, m.titleInput.View())

	// Topic field
	topicLabel := "  Topic"
	if m.activeField == FieldTopic {
		topicLabel = activeStyle.Render("> Topic")
	} else {
		topicLabel = labelStyle.Render("  Topic")
	}
	topicField := fmt.Sprintf("%s\n  %s", topicLabel, m.topicInput.View())

	// Persona field
	personaLabel := "  Persona"
	if m.activeField == FieldPersona {
		personaLabel = activeStyle.Render("> Persona")
	} else {
		personaLabel = labelStyle.Render("  Persona")
	}
	personaField := fmt.Sprintf("%s\n  %s", personaLabel, m.personaInput.View())

	// Provider field
	providerLabel := "  Provider"
	if m.activeField == FieldProvider {
		providerLabel = activeStyle.Render("> Provider")
	} else {
		providerLabel = labelStyle.Render("  Provider")
	}

	var providerOptions string
	for i, p := range m.providers {
		if i == m.providerIdx {
			providerOptions += selectedProviderStyle.Render(p.Display) + " "
		} else {
			providerOptions += normalProviderStyle.Render(p.Display) + " "
		}
	}
	if len(m.providers) == 0 {
		providerOptions = mutedStyle.Render("No providers configured")
	}
	providerField := fmt.Sprintf("%s\n  %s", providerLabel, providerOptions)

	// Help text
	help := styles.HelpStyle.Render("Tab/↓ next field • Shift+Tab/↑ prev field • ←/→ select provider • Enter to continue • Esc to go back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		fmt.Sprintf("  %s %s", timestampLabel, timestamp),
		"",
		titleField,
		"",
		topicField,
		"",
		personaField,
		"",
		providerField,
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

// Helper to get providers sorted alphabetically
func getProvidersFromConfig() []ProviderInfo {
	cfg := config.Get()
	if cfg == nil {
		return []ProviderInfo{}
	}

	var providers []ProviderInfo
	for name, providerCfg := range cfg.Providers {
		displayName := name
		if providerCfg.ChatModel != "" {
			displayName = fmt.Sprintf("%s (%s)", name, providerCfg.ChatModel)
		}
		providers = append(providers, ProviderInfo{
			Name:    name,
			Display: displayName,
			Model:   providerCfg.ChatModel,
		})
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Name < providers[j].Name
	})

	return providers
}
