package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#7C3AED") // Purple
	SecondaryColor = lipgloss.Color("#10B981") // Green
	AccentColor    = lipgloss.Color("#3B82F6") // Blue
	MutedColor     = lipgloss.Color("#6B7280") // Gray
	ErrorColor     = lipgloss.Color("#EF4444") // Red

	// Title style for headers
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginBottom(1)

	// Subtitle for descriptions
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			MarginBottom(2)

	// Host message style (right-aligned, blue)
	HostMessageStyle = lipgloss.NewStyle().
				Background(AccentColor).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 1).
				MarginLeft(4).
				MarginBottom(1)

	// Guest message style (left-aligned, green)
	GuestMessageStyle = lipgloss.NewStyle().
				Background(SecondaryColor).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 1).
				MarginRight(4).
				MarginBottom(1)

	// Input prompt style
	InputPromptStyle = lipgloss.NewStyle().
				Foreground(MutedColor).
				Italic(true)

	// Selected item in list
	SelectedStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	// Normal item in list
	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	// Thinking indicator style
	ThinkingStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Italic(true)

	// Header bar style
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(PrimaryColor).
			Padding(0, 2).
			MarginBottom(1)

	// Footer/help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginTop(1)

	// Box style for containers
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2)

	// Logo style
	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor)

	// Voice description style
	VoiceDescStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true)

	// Analog wave style
	WaveStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor)

	// Simple transcript speaker labels (accent color)
	HostLabelStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	GuestLabelStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	// Simple transcript text (normal color)
	TranscriptTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF"))

	// Transcript panel style (no border, uses available space)
	TranscriptPanelStyle = lipgloss.NewStyle().
				PaddingLeft(2)
)

// Logo returns the VibeCast ASCII art logo
func Logo() string {
	logo := `
 ╦  ╦╦╔╗ ╔═╗╔═╗╔═╗╔═╗╔╦╗
 ╚╗╔╝║╠╩╗║╣ ║  ╠═╣╚═╗ ║
  ╚╝ ╩╚═╝╚═╝╚═╝╩ ╩╚═╝ ╩ `
	return LogoStyle.Render(logo)
}

// LogoWithTitle returns the VibeCast logo with the conversation title displayed underneath
func LogoWithTitle(title string) string {
	logo := Logo()
	if title == "" {
		return logo
	}

	// Center the title under the logo
	titleStyle := lipgloss.NewStyle().
		Foreground(MutedColor).
		Italic(true)

	// Logo is roughly 26 chars wide, center title accordingly
	centeredTitle := lipgloss.NewStyle().
		Width(26).
		Align(lipgloss.Center).
		Render(titleStyle.Render(title))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		centeredTitle,
	)
}
