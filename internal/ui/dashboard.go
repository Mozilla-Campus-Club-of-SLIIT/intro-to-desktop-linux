package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardModel struct {
	width  int
	height int

	content     *TerminalModel
	leaderboard *LeaderboardModel // Changed to pointer
}

func NewDashboardModel() *DashboardModel {
	term, _ := NewTerminalModel()
	leaderboardModel := NewLeaderboardModel(0, 0) // Initialize with zero dimensions, will be updated by WindowSizeMsg
	return &DashboardModel{
		content:     term,
		leaderboard: &leaderboardModel, // Return a pointer to the initialized LeaderboardModel
	}
}

func (m *DashboardModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	if m.content != nil {
		cmds = append(cmds, m.content.Init())
	}
	// Call Init() for the leaderboard model as well!
	if m.leaderboard != nil {
		cmds = append(cmds, m.leaderboard.Init())
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

	width := m.width - 2*margin
	height := m.height - 2*(margin/4)

	effectiveLayoutWidth := width - gap
	leaderWidth := effectiveLayoutWidth / 5
	termWidth := effectiveLayoutWidth - leaderWidth
	termHeight := height - footerHeight - gap

	// Update TerminalModel
	if m.content != nil {
		// Pass a WindowSizeMsg to TerminalModel if the main WindowSizeMsg was received,
		// or if the sub-model simply needs to update its state with the current calculated dimensions.
		// For robustness, always set the dimensions before calling Update, and also pass msg.
		m.content.width = termWidth
		m.content.height = termHeight
		newModel, cmd = m.content.Update(msg)
		m.content = newModel.(*TerminalModel)
		cmds = append(cmds, cmd)
	}

	// Update LeaderboardModel
	if m.leaderboard != nil {
		m.leaderboard.width = leaderWidth
		m.leaderboard.height = termHeight
		newModel, cmd = m.leaderboard.Update(msg)
		m.leaderboard = newModel.(*LeaderboardModel)
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
	termHeight := height - footerHeight - gap

	// Ensure dimensions are propagated correctly to sub-models for rendering,
	// by setting their internal width/height before calling View().
	// This makes View() itself stateless regarding dimensions, relying on the model's fields.
	if m.leaderboard != nil {
		m.leaderboard.width = leaderWidth
		m.leaderboard.height = termHeight
	}
	if m.content != nil {
		m.content.width = termWidth
		m.content.height = termHeight
	}


	leaderboardView := lipgloss.NewStyle().
		Width(leaderWidth).
		Height(termHeight).
		Render(m.leaderboard.View()) // Correctly call View on pointer

	contentView := lipgloss.NewStyle().
		Width(termWidth).
		Height(termHeight). // Also give content height
		Render(m.content.View())

	footer := m.footerView(width)
	spacer := lipgloss.NewStyle().Width(gap).Render("")
	mainArea := lipgloss.JoinHorizontal(lipgloss.Top, contentView, spacer, leaderboardView)

	ui := lipgloss.JoinVertical(lipgloss.Left, mainArea, footer)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, ui)
}
