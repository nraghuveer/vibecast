package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nraghuveer/vibecast/cmd/cli/mock"
	"github.com/nraghuveer/vibecast/cmd/cli/screens"
)

// Screen represents the current screen state
type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenNewConversation
	ScreenConversationList
	ScreenTopic
	ScreenPersona
	ScreenVoice
	ScreenProvider
	ScreenConversation
	ScreenPreset
	ScreenTemplateName
	ScreenTemplateTopic
	ScreenTemplatePersona
)

// Model is the main application model
type Model struct {
	screen           Screen
	welcome          screens.WelcomeModel
	newConversation  screens.CreateConversationModel
	conversationList screens.ConversationListModel
	topic            screens.TopicModel
	persona          screens.PersonaModel
	voice            screens.VoiceModel
	provider         screens.ProviderModel
	conversation     screens.ConversationModel
	preset           screens.PresetModel
	templateName     screens.TemplateNameModel

	// Collected data
	selectedTitle    string
	selectedTopic    string
	selectedPersona  string
	selectedVoice    mock.Voice
	selectedProvider string

	// Template creation data
	newTemplateName string

	width  int
	height int
}

// NewModel creates a new application model
func NewModel() Model {
	return Model{
		screen:       ScreenWelcome,
		welcome:      screens.NewWelcomeModel(),
		topic:        screens.NewTopicModel(),
		persona:      screens.NewPersonaModel(),
		voice:        screens.NewVoiceModel(),
		provider:     screens.NewProviderModel(),
		preset:       screens.NewPresetModel(),
		templateName: screens.NewTemplateNameModel(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.welcome.Init()
}

// Update handles all messages and routes to the appropriate screen
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window size globally
	if wsm, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = wsm.Width
		m.height = wsm.Height
	}

	switch m.screen {
	case ScreenWelcome:
		return m.updateWelcome(msg)
	case ScreenNewConversation:
		return m.updateNewConversation(msg)
	case ScreenConversationList:
		return m.updateConversationList(msg)
	case ScreenTopic:
		return m.updateTopic(msg)
	case ScreenPersona:
		return m.updatePersona(msg)
	case ScreenVoice:
		return m.updateVoice(msg)
	case ScreenProvider:
		return m.updateProvider(msg)
	case ScreenConversation:
		return m.updateConversation(msg)
	case ScreenPreset:
		return m.updatePreset(msg)
	case ScreenTemplateName:
		return m.updateTemplateName(msg)
	case ScreenTemplateTopic:
		return m.updateTemplateTopic(msg)
	case ScreenTemplatePersona:
		return m.updateTemplatePersona(msg)
	}

	return m, nil
}

func (m Model) updateWelcome(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.welcome, cmd = m.welcome.Update(msg)

	// Check for option selection
	if wsm, ok := msg.(screens.WelcomeOptionSelectedMsg); ok {
		switch wsm.Option {
		case screens.OptionNewConversation:
			m.screen = ScreenNewConversation
			m.newConversation = screens.NewCreateConversationModel()
			return m, m.newConversation.Init()
		case screens.OptionContinueConversation:
			m.screen = ScreenConversationList
			m.conversationList = screens.NewConversationListModel()
			return m, m.conversationList.Init()
		case screens.OptionQuickStart:
			m.screen = ScreenPreset
			m.preset = screens.NewPresetModel()
			return m, m.preset.Init()
		case screens.OptionCreateTemplate:
			m.screen = ScreenTemplateName
			m.templateName = screens.NewTemplateNameModel()
			return m, m.templateName.Init()
		}
	}

	return m, cmd
}

func (m Model) updateNewConversation(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.newConversation, cmd = m.newConversation.Update(msg)

	// Check for new conversation created
	if ncm, ok := msg.(screens.NewConversationCreatedMsg); ok {
		m.selectedTitle = ncm.Title
		m.selectedTopic = ncm.Topic
		m.selectedPersona = ncm.Persona
		m.selectedProvider = ncm.Provider
		m.screen = ScreenVoice
		m.voice = screens.NewVoiceModel()
		return m, m.voice.Init()
	}

	// Check for back navigation
	if _, ok := msg.(screens.BackToWelcomeMsg); ok {
		m.screen = ScreenWelcome
		m.welcome = screens.NewWelcomeModel()
		return m, m.welcome.Init()
	}

	return m, cmd
}

func (m Model) updateConversationList(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.conversationList, cmd = m.conversationList.Update(msg)

	// Check for conversation selected
	if csm, ok := msg.(screens.ConversationSelectedMsg); ok {
		m.selectedTitle = csm.Conversation.Title
		m.selectedTopic = csm.Conversation.Topic
		m.selectedPersona = csm.Conversation.Persona
		m.selectedProvider = csm.Conversation.Provider
		m.selectedVoice = mock.Voice{
			ID:   csm.Conversation.VoiceID,
			Name: csm.Conversation.VoiceName,
		}
		m.screen = ScreenConversation
		m.conversation = screens.NewConversationModelWithTitle(
			m.selectedTitle,
			m.selectedTopic,
			m.selectedPersona,
			m.selectedVoice,
			m.selectedProvider,
			m.width,
			m.height,
		)
		return m, m.conversation.Init()
	}

	// Check for back navigation
	if _, ok := msg.(screens.BackToWelcomeMsg); ok {
		m.screen = ScreenWelcome
		m.welcome = screens.NewWelcomeModel()
		return m, m.welcome.Init()
	}

	return m, cmd
}

func (m Model) updateTopic(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.topic, cmd = m.topic.Update(msg)

	// Check for topic selection
	if tsm, ok := msg.(screens.TopicSelectedMsg); ok {
		m.selectedTopic = tsm.Topic
		m.screen = ScreenPersona
		return m, m.persona.Init()
	}

	return m, cmd
}

func (m Model) updatePersona(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.persona, cmd = m.persona.Update(msg)

	// Check for persona selection
	if psm, ok := msg.(screens.PersonaSelectedMsg); ok {
		m.selectedPersona = psm.Persona
		m.screen = ScreenVoice
		return m, m.voice.Init()
	}

	return m, cmd
}

func (m Model) updateVoice(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.voice, cmd = m.voice.Update(msg)

	// Check for voice selection
	if vsm, ok := msg.(screens.VoiceSelectedMsg); ok {
		m.selectedVoice = vsm.Voice

		// If provider is already selected (from NewConversation screen), go directly to conversation
		if m.selectedProvider != "" {
			m.screen = ScreenConversation
			m.conversation = screens.NewConversationModelWithTitle(
				m.selectedTitle,
				m.selectedTopic,
				m.selectedPersona,
				m.selectedVoice,
				m.selectedProvider,
				m.width,
				m.height,
			)
			return m, m.conversation.Init()
		}

		// Otherwise go to provider selection
		m.screen = ScreenProvider
		m.provider = screens.NewProviderModel()
		return m, m.provider.Init()
	}

	return m, cmd
}

func (m Model) updateProvider(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.provider, cmd = m.provider.Update(msg)

	if psm, ok := msg.(screens.ProviderSelectedMsg); ok {
		m.selectedProvider = psm.Provider
		m.screen = ScreenConversation
		m.conversation = screens.NewConversationModelWithTitle(
			m.selectedTitle,
			m.selectedTopic,
			m.selectedPersona,
			m.selectedVoice,
			m.selectedProvider,
			m.width,
			m.height,
		)
		return m, m.conversation.Init()
	}

	return m, cmd
}

func (m Model) updateConversation(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.conversation, cmd = m.conversation.Update(msg)
	return m, cmd
}

func (m Model) updatePreset(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.preset, cmd = m.preset.Update(msg)

	// Check for preset selection
	if psm, ok := msg.(screens.PresetSelectedMsg); ok {
		m.selectedTopic = psm.Template.Topic
		m.selectedPersona = psm.Template.Persona
		m.screen = ScreenVoice
		m.voice = screens.NewVoiceModel()
		return m, m.voice.Init()
	}

	// Check for back navigation
	if _, ok := msg.(screens.BackToWelcomeMsg); ok {
		m.screen = ScreenWelcome
		m.welcome = screens.NewWelcomeModel()
		return m, m.welcome.Init()
	}

	return m, cmd
}

func (m Model) updateTemplateName(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.templateName, cmd = m.templateName.Update(msg)

	// Check for template name entered
	if tnm, ok := msg.(screens.TemplateNameEnteredMsg); ok {
		m.newTemplateName = tnm.Name
		m.screen = ScreenTemplateTopic
		m.topic = screens.NewTopicModel()
		return m, m.topic.Init()
	}

	// Check for back navigation
	if _, ok := msg.(screens.BackToWelcomeMsg); ok {
		m.screen = ScreenWelcome
		m.welcome = screens.NewWelcomeModel()
		return m, m.welcome.Init()
	}

	return m, cmd
}

func (m Model) updateTemplateTopic(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.topic, cmd = m.topic.Update(msg)

	// Check for topic selection
	if tsm, ok := msg.(screens.TopicSelectedMsg); ok {
		m.selectedTopic = tsm.Topic
		m.screen = ScreenTemplatePersona
		m.persona = screens.NewPersonaModel()
		return m, m.persona.Init()
	}

	return m, cmd
}

func (m Model) updateTemplatePersona(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.persona, cmd = m.persona.Update(msg)

	// Check for persona selection - save template and go back to welcome
	if psm, ok := msg.(screens.PersonaSelectedMsg); ok {
		m.selectedPersona = psm.Persona

		// Save the new template
		mock.AddTemplate(mock.Template{
			ID:      m.newTemplateName,
			Name:    m.newTemplateName,
			Topic:   m.selectedTopic,
			Persona: m.selectedPersona,
		})

		// Reset and go back to welcome
		m.newTemplateName = ""
		m.selectedTopic = ""
		m.selectedPersona = ""
		m.screen = ScreenWelcome
		m.welcome = screens.NewWelcomeModel()
		return m, m.welcome.Init()
	}

	return m, cmd
}

// View renders the current screen
func (m Model) View() string {
	switch m.screen {
	case ScreenWelcome:
		return m.welcome.View()
	case ScreenNewConversation:
		return m.newConversation.View()
	case ScreenConversationList:
		return m.conversationList.View()
	case ScreenTopic, ScreenTemplateTopic:
		return m.topic.View()
	case ScreenPersona, ScreenTemplatePersona:
		return m.persona.View()
	case ScreenVoice:
		return m.voice.View()
	case ScreenProvider:
		return m.provider.View()
	case ScreenConversation:
		return m.conversation.View()
	case ScreenPreset:
		return m.preset.View()
	case ScreenTemplateName:
		return m.templateName.View()
	}
	return ""
}
