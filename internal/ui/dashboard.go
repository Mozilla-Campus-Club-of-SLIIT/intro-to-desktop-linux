package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardModel struct {
	width  int
	height int

	content      *TerminalModel
	leaderboard  *LeaderboardModel
	notification *NotificationModel // New: NotificationModel
}

func NewDashboardModel() *DashboardModel {
	term, _ := NewTerminalModel()
	leaderboardModel := NewLeaderboardModel(0, 0)
	notificationModel := NewNotificationModel() // New: Initialize NotificationModel
	return &DashboardModel{
		content:      term,
		leaderboard:  leaderboardModel,
		notification: notificationModel, // New: Add to struct
	}
}

func (m *DashboardModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	if m.content != nil {
		cmds = append(cmds, m.content.Init())
	}
	if m.leaderboard != nil {
		cmds = append(cmds, m.leaderboard.Init())
	}
	if m.notification != nil { // New: Call Init() for notification model
		cmds = append(cmds, m.notification.Init())
	}
	return tea.Batch(cmds...)
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var newModel tea.Model
	var cmd tea.Cmd

	// Update DashboardModel's dimensions
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
	}

	// Calculate dimensions for sub-models based on current DashboardModel dimensions
	const footerHeight = 1
	const margin = 2
	const gap = 2
	const notificationHeight = 5 // User requested 5 lines for notifications
	const notificationHeader = 1 // 1 line for notification header
	const notificationSpace = 2  // 2 lines for space below leaderboard

	effectiveWidth := m.width - 2*margin
	effectiveHeight := m.height - 2*(margin/4)

	effectiveLayoutWidth := effectiveWidth - gap
	leaderWidth := effectiveLayoutWidth / 5
	termWidth := effectiveLayoutWidth - leaderWidth

	// Calculate heights considering the new notification area
	leaderboardAndTerminalHeight := effectiveHeight - footerHeight - gap - notificationHeader - notificationHeight - notificationSpace
	termHeight := leaderboardAndTerminalHeight
	leaderHeight := leaderboardAndTerminalHeight

	// Update TerminalModel
	if m.content != nil {
		m.content.width = termWidth
		m.content.height = termHeight
		newModel, cmd = m.content.Update(msg)
		m.content = newModel.(*TerminalModel)
		cmds = append(cmds, cmd)
	}

	// Update LeaderboardModel with correct dimensions
	if m.leaderboard != nil {
		// Pass a custom WindowSizeMsg with leaderboard-specific dimensions
		leaderMsg := msg
		if _, ok := msg.(tea.WindowSizeMsg); ok {
			leaderMsg = tea.WindowSizeMsg{Width: leaderWidth, Height: leaderHeight}
		}
		_, cmd = m.leaderboard.Update(leaderMsg)
		cmds = append(cmds, cmd)
	}

	// Update NotificationModel
	if m.notification != nil { // New: Update notification model
		m.notification.width = leaderWidth // Notifications align with leaderboard width
		// The height for notifications is fixed as per requirement
		m.notification.height = notificationHeight + notificationHeader + notificationSpace
		newModel, cmd = m.notification.Update(msg)
		m.notification = newModel.(*NotificationModel)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *DashboardModel) View() string {
	if m.width <= 0 || m.height <= 0 {
		return "Initializing..."
	}
	return m.render()
}

func (m *DashboardModel) render() string {
	const footerHeight = 1
	const margin = 2
	const gap = 2
	const notificationHeight = 5 // User requested 5 lines for notifications content
	const notificationHeader = 1 // 1 line for notification header
	const notificationSpace = 2  // 2 lines for space between leaderboard and notifications

	width := m.width - 2*margin
	height := m.height - 2*(margin/4)

	// Calculate dimensions for leaderboard and content
	effectiveLayoutWidth := width - gap
	leaderWidth := effectiveLayoutWidth / 5
	termWidth := effectiveLayoutWidth - leaderWidth

	// Heights calculation needs to account for notifications
	// Total height available for dynamic content = height - footerHeight - gap (above footer)
	// Remaining height after terminal and leaderboard = total dynamic height - (notificationHeader + notificationHeight + notificationSpace)
	leaderboardAndTerminalHeight := height - footerHeight - gap - notificationHeader - notificationHeight - notificationSpace
	termHeight := leaderboardAndTerminalHeight
	leaderHeight := leaderboardAndTerminalHeight

	// Dimensions are now set in Update() method, so we only use them for rendering here
	// Ensure notification dimensions are set (notification model doesn't receive WindowSizeMsg in Update)
	if m.notification != nil {
		m.notification.width = leaderWidth
		m.notification.height = notificationHeight + notificationHeader
	}

	leaderboardView := lipgloss.NewStyle().
		Width(leaderWidth).
		Height(leaderHeight).
		Render(m.leaderboard.View())

	contentView := lipgloss.NewStyle().
		Width(termWidth).
		Height(termHeight).
		Render(m.content.View())

	// New: Render notification view
	notificationView := lipgloss.NewStyle().
		Width(leaderWidth).
		Height(notificationHeight + notificationHeader). // Actual height rendered by NotificationModel.View()
		Render(m.notification.View())

	footer := m.footerView(width)
	spacer := lipgloss.NewStyle().Width(gap).Render("")

	// Join leaderboard and notification vertically, with space
	leaderboardColumn := lipgloss.JoinVertical(lipgloss.Left,
		leaderboardView,
		lipgloss.NewStyle().Height(notificationSpace).Render(""), // Space between leaderboard and notifications
		notificationView,
	)

	mainArea := lipgloss.JoinHorizontal(lipgloss.Top, contentView, spacer, leaderboardColumn)

	ui := lipgloss.JoinVertical(lipgloss.Left, mainArea, footer)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, ui)
}
