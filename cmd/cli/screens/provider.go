package screens

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
	"github.com/nraghuveer/vibecast/lib/config"
	"github.com/nraghuveer/vibecast/lib/logger"
)

type ProviderInfo struct {
	Name    string
	Display string
	Model   string
}

type ProviderModel struct {
	providers []ProviderInfo
	cursor    int
	selected  string
	width     int
	height    int
	logger    *logger.Logger
}

func NewProviderModel() ProviderModel {
	providers := getAvailableProviders()
	return ProviderModel{
		providers: providers,
		cursor:    0,
		logger:    logger.GetInstance(),
	}
}

func getAvailableProviders() []ProviderInfo {
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

func (m ProviderModel) Init() tea.Cmd {
	return nil
}

func (m ProviderModel) Update(msg tea.Msg) (ProviderModel, tea.Cmd) {
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
			if m.cursor < len(m.providers)-1 {
				m.cursor++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			m.selected = m.providers[m.cursor].Name
			m.logger.Info("provider_selected", "provider", m.selected)
			return m, func() tea.Msg { return ProviderSelectedMsg{Provider: m.selected} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			m.logger.Info("provider_selection_quit")
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ProviderModel) View() string {
	title := styles.TitleStyle.Render("Select AI Provider")

	description := styles.SubtitleStyle.Render(
		"Choose the AI provider for this conversation",
	)

	var items string
	for i, provider := range m.providers {
		cursor := "  "
		itemStyle := styles.NormalStyle
		if i == m.cursor {
			cursor = "> "
			itemStyle = styles.SelectedStyle
		}

		item := fmt.Sprintf("%s%s", cursor, itemStyle.Render(provider.Display))
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

func (m ProviderModel) SelectedProvider() string {
	return m.selected
}

type ProviderSelectedMsg struct {
	Provider string
}
