package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardModel struct {
	width       int
	height      int
	environment string

	content     ContentModel
	leaderboard LeaderboardModel
}

func (m *DashboardModel) Init() tea.Cmd {
	return nil
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
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

	width := m.width - 2*margin
	height := m.height - 2*(margin/4)

	m.leaderboard.width = (width - 2*gap) / 4
	m.content.width = (width - 2*gap) - m.leaderboard.width

	m.content.height = height - footerHeight - gap
	m.leaderboard.height = m.content.height

	leaderboard := lipgloss.NewStyle().
		Width(m.leaderboard.width).
		Height(m.leaderboard.height).
		Render(m.leaderboard.View())

	content := lipgloss.NewStyle().
		Width(m.content.width).
		Height(m.content.height).
		Render(m.content.View())

	footer := m.footerView(width)

	spacer := lipgloss.NewStyle().Width(gap).Render("")
	mainArea := lipgloss.JoinHorizontal(lipgloss.Top, content, spacer, leaderboard)

	ui := lipgloss.JoinVertical(lipgloss.Left, mainArea, footer)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		ui,
	)
}
