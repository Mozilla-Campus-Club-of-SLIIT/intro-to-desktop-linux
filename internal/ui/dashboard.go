package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardModel struct {
	width  int
	height int

	content     *TerminalModel
	leaderboard *LeaderboardModel
}

func NewDashboardModel() *DashboardModel {
	term, _ := NewTerminalModel()
	leaderboardModel := NewLeaderboardModel(0, 0)
	return &DashboardModel{
		content:     term,
		leaderboard: leaderboardModel,
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
	return tea.Batch(cmds...)
}

func (m *DashboardModel) Stop() {
	if m.leaderboard != nil {
		m.leaderboard.Stop()
	}
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

	effectiveWidth := m.width - 2*margin
	effectiveHeight := m.height - 2*(margin/4)

	effectiveLayoutWidth := effectiveWidth - gap
	leaderWidth := effectiveLayoutWidth / 5
	termWidth := effectiveLayoutWidth - leaderWidth

	// Calculate heights
	leaderboardAndTerminalHeight := effectiveHeight - footerHeight - gap
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

	width := m.width - 2*margin
	height := m.height - 2*(margin/4)

	// Calculate dimensions for leaderboard and content
	effectiveLayoutWidth := width - gap
	leaderWidth := effectiveLayoutWidth / 5
	termWidth := effectiveLayoutWidth - leaderWidth

	// Heights calculation
	leaderboardAndTerminalHeight := height - footerHeight - gap
	termHeight := leaderboardAndTerminalHeight
	leaderHeight := leaderboardAndTerminalHeight

	leaderboardView := lipgloss.NewStyle().
		Width(leaderWidth).
		Height(leaderHeight).
		Render(m.leaderboard.View())

	contentView := lipgloss.NewStyle().
		Width(termWidth).
		Height(termHeight).
		Render(m.content.View())

	footer := m.footerView(width)
	spacer := lipgloss.NewStyle().Width(gap).Render("")

	mainArea := lipgloss.JoinHorizontal(lipgloss.Top, contentView, spacer, leaderboardView)

	ui := lipgloss.JoinVertical(lipgloss.Left, mainArea, footer)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, ui)
}
