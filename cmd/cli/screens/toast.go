package screens

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nraghuveer/vibecast/cmd/cli/styles"
)

// Toast represents a toast notification
type Toast struct {
	Message   string
	Type      ToastType
	CreatedAt time.Time
}

// ToastType represents the type of toast notification
type ToastType int

const (
	ToastError ToastType = iota
	ToastWarning
	ToastInfo
	ToastSuccess
)

// ToastModel manages toast notifications
type ToastModel struct {
	toasts []Toast
}

// NewToastModel creates a new toast model
func NewToastModel() ToastModel {
	return ToastModel{
		toasts: []Toast{},
	}
}

// AddToast adds a new toast notification
func (m *ToastModel) AddToast(message string, toastType ToastType) {
	m.toasts = append(m.toasts, Toast{
		Message:   message,
		Type:      toastType,
		CreatedAt: time.Now(),
	})
}

// AddError adds an error toast
func (m *ToastModel) AddError(message string) {
	m.AddToast(message, ToastError)
}

// AddWarning adds a warning toast
func (m *ToastModel) AddWarning(message string) {
	m.AddToast(message, ToastWarning)
}

// AddInfo adds an info toast
func (m *ToastModel) AddInfo(message string) {
	m.AddToast(message, ToastInfo)
}

// AddSuccess adds a success toast
func (m *ToastModel) AddSuccess(message string) {
	m.AddToast(message, ToastSuccess)
}

// ClearToasts removes all toasts
func (m *ToastModel) ClearToasts() {
	m.toasts = []Toast{}
}

// RemoveToast removes a toast at the given index
func (m *ToastModel) RemoveToast(index int) {
	if index >= 0 && index < len(m.toasts) {
		m.toasts = append(m.toasts[:index], m.toasts[index+1:]...)
	}
}

// GetToasts returns all current toasts
func (m *ToastModel) GetToasts() []Toast {
	return m.toasts
}

// HasToasts returns true if there are any toasts
func (m *ToastModel) HasToasts() bool {
	return len(m.toasts) > 0
}

// ToastDismissMsg is sent when a toast should be dismissed
type ToastDismissMsg struct {
	Index int
}

// DismissToastCmd creates a command to dismiss a toast after a duration
func DismissToastCmd(index int, duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return ToastDismissMsg{Index: index}
	})
}

// RenderToast renders a single toast notification
func RenderToast(toast Toast) string {
	var style lipgloss.Style
	var icon string

	switch toast.Type {
	case ToastError:
		style = lipgloss.NewStyle().
			Background(lipgloss.Color("#EF4444")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2).
			Bold(true)
		icon = "✗"
	case ToastWarning:
		style = lipgloss.NewStyle().
			Background(lipgloss.Color("#F59E0B")).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 2).
			Bold(true)
		icon = "⚠"
	case ToastSuccess:
		style = lipgloss.NewStyle().
			Background(lipgloss.Color("#10B981")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2).
			Bold(true)
		icon = "✓"
	default:
		style = lipgloss.NewStyle().
			Background(styles.PrimaryColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2).
			Bold(true)
		icon = "ℹ"
	}

	content := icon + " " + toast.Message
	return style.Render(content)
}

// RenderToasts renders all toasts in the upper left position
func RenderToasts(toasts []Toast) string {
	if len(toasts) == 0 {
		return ""
	}

	var rendered []string
	for _, toast := range toasts {
		rendered = append(rendered, RenderToast(toast))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}
