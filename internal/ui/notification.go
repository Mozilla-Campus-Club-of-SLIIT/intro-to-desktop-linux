package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MsgNotification is a message for the notification system.
type MsgNotification struct {
	Timestamp time.Time
	Message   string
	IsError   bool // To style error messages differently
}

type NotificationModel struct {
	width  int
	height int // Actual rendered height (header + content + spacing)

	notifications []MsgNotification
	headerStyle   lipgloss.Style
	messageStyle  lipgloss.Style
	errorStyle    lipgloss.Style

	maxNotifications int // Number of message lines to display
}

// NewNotificationModel creates a new notification model.
func NewNotificationModel() *NotificationModel {
	return &NotificationModel{
		notifications:    make([]MsgNotification, 0),
		maxNotifications: 5, // User requested 5 lines for notifications
		headerStyle: lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorDarkGreen).
			Padding(0, 1),
		messageStyle: lipgloss.NewStyle().
			Foreground(ColorWhite). // Default message color
			Padding(0, 1),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4F5E")). // Error message color (similar to ColorWarmRed)
			Padding(0, 1),
	}
}

func (m *NotificationModel) Init() tea.Cmd {
	return nil
}

func (m *NotificationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // Update dimensions if needed
		m.width = msg.Width
		m.height = msg.Height // This will be the actual allocated height for the whole notification area
	case MsgNotification:
		// Add new notification and trim old ones
		m.notifications = append(m.notifications, msg)
		if len(m.notifications) > m.maxNotifications {
			m.notifications = m.notifications[len(m.notifications)-m.maxNotifications:]
		}
	}
	return m, nil
}

func (m NotificationModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "" // Don't render if dimensions not set
	}

	header := m.headerStyle.Width(m.width).Render("Notifications")

	var content strings.Builder
	// Render messages, newest at the bottom
	for i := 0; i < len(m.notifications); i++ {
		notif := m.notifications[i]
		style := m.messageStyle
		if notif.IsError {
			style = m.errorStyle
		}
		content.WriteString(style.Width(m.width).Render(fmt.Sprintf("%s %s", notif.Timestamp.Format("15:04:05"), notif.Message)))
		content.WriteString("\n")
	}

	renderedMessages := strings.Split(strings.TrimSuffix(content.String(), "\n"), "\n")
	if len(renderedMessages) == 1 && renderedMessages[0] == "" { // Handle case with no actual messages
		renderedMessages = []string{}
	}

	messageLinesToRender := make([]string, 0, m.maxNotifications)
	// Fill from newest messages up to maxNotifications
	startIdx := 0
	if len(renderedMessages) > m.maxNotifications {
		startIdx = len(renderedMessages) - m.maxNotifications
	}
	for i := startIdx; i < len(renderedMessages); i++ {
		messageLinesToRender = append(messageLinesToRender, renderedMessages[i])
	}


	var finalContent strings.Builder
	// Add empty lines at the top if there are fewer than maxNotifications
	for i := 0; i < m.maxNotifications-len(messageLinesToRender); i++ {
		finalContent.WriteString(m.messageStyle.Width(m.width).Render("...")) // Placeholder for empty lines
		finalContent.WriteString("\n")
	}
	finalContent.WriteString(strings.Join(messageLinesToRender, "\n"))


	return lipgloss.JoinVertical(lipgloss.Left, header, finalContent.String())
}
