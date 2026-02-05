package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
	"github.com/nraghuveer/vibecast/lib/db"
	"github.com/nraghuveer/vibecast/lib/logger"
)

// ConversationListModel displays a list of existing conversations
type ConversationListModel struct {
	db            *db.DB
	conversations []db.Conversation
	cursor        int
	showDetails   bool // Toggle for showing topic/persona (Ctrl+I)
	width         int
	height        int
	err           error
	logger        *logger.Logger
}

// ConversationSelectedMsg is sent when a conversation is selected
type ConversationSelectedMsg struct {
	Conversation db.Conversation
}

func NewConversationListModel(database *db.DB) ConversationListModel {
	log := logger.GetInstance()
	conversations, err := database.GetAllConversations()
	if err != nil {
		log.LogError("conversation_list_load", err)
	} else {
		log.Info("conversation_list_loaded", "count", len(conversations))
	}

	return ConversationListModel{
		db:            database,
		conversations: conversations,
		cursor:        0,
		showDetails:   false,
		err:           err,
		logger:        log,
	}
}

func (m ConversationListModel) Init() tea.Cmd {
	return nil
}

func (m ConversationListModel) Update(msg tea.Msg) (ConversationListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			m.logger.Info("conversation_list_quit")
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			m.logger.Info("conversation_list_back_to_welcome")
			return m, func() tea.Msg { return BackToWelcomeMsg{} }

		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+i"))):
			m.showDetails = !m.showDetails
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(m.conversations)-1 {
				m.cursor++
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if len(m.conversations) > 0 {
				selected := m.conversations[m.cursor]
				m.logger.Info("conversation_selected",
					"id", selected.ID,
					"title", selected.Title,
				)
				return m, func() tea.Msg {
					return ConversationSelectedMsg{Conversation: selected}
				}
			}
		}
	}

	return m, nil
}

func (m ConversationListModel) View() string {
	title := styles.TitleStyle.Render("Continue Conversation")
	subtitle := styles.SubtitleStyle.Render("Select a conversation to continue")

	if m.err != nil {
		errMsg := styles.HelpStyle.Render(fmt.Sprintf("Error loading conversations: %v", m.err))
		content := lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", errMsg)
		box := styles.BoxStyle.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}

	if len(m.conversations) == 0 {
		emptyMsg := styles.HelpStyle.Render("No conversations found. Create a new one first!")
		help := styles.HelpStyle.Render("Esc to go back")
		content := lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "", emptyMsg, "", help)
		box := styles.BoxStyle.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}

	// Styles
	titleStylePrimary := lipgloss.NewStyle().Foreground(styles.PrimaryColor).Bold(true)
	timestampStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)
	detailLabelStyle := lipgloss.NewStyle().Foreground(styles.MutedColor).Italic(true)
	detailTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	var items string
	for i, conv := range m.conversations {
		cursor := "  "
		itemTitleStyle := titleStylePrimary
		if i == m.cursor {
			cursor = "> "
			itemTitleStyle = titleStylePrimary.Bold(true).Underline(true)
		}

		// Title in primary color
		convTitle := conv.Title
		if convTitle == "" {
			convTitle = "Untitled Conversation"
		}
		titleLine := itemTitleStyle.Render(convTitle)

		// Timestamp in muted color
		timestamp := conv.CreatedAt.Format("Jan 02, 2006 3:04 PM")
		timestampLine := timestampStyle.Render(timestamp)

		item := fmt.Sprintf("%s%s\n   %s", cursor, titleLine, timestampLine)

		// Show topic and persona if Ctrl+I is toggled
		if m.showDetails {
			topicLine := fmt.Sprintf("   %s %s",
				detailLabelStyle.Render("Topic:"),
				detailTextStyle.Render(truncate(conv.Topic, 40)))
			personaLine := fmt.Sprintf("   %s %s",
				detailLabelStyle.Render("Persona:"),
				detailTextStyle.Render(truncate(conv.Persona, 40)))
			item += "\n" + topicLine + "\n" + personaLine
		}

		items += item + "\n\n"
	}

	// Help text
	detailsHint := "Ctrl+I to show details"
	if m.showDetails {
		detailsHint = "Ctrl+I to hide details"
	}
	help := styles.HelpStyle.Render(fmt.Sprintf("↑/↓ or j/k to navigate | Enter to select | %s | Esc to go back", detailsHint))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
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

// truncate shortens a string to maxLen and adds "..." if truncated
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
