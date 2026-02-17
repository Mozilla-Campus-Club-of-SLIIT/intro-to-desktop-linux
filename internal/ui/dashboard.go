package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardModel struct {
	width  int
	height int

	content     *TerminalModel
	leaderboard LeaderboardModel
}

func NewDashboardModel() *DashboardModel {
	term, _ := NewTerminalModel()
	return &DashboardModel{
		content:     term,
		leaderboard: LeaderboardModel{},
	}
}

func (m *DashboardModel) Init() tea.Cmd {
	if m.content != nil {
		return m.content.Init()
	}
	return nil
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.content != nil {
		const margin = 2
		const gap = 2
		layoutWidth := m.width - 2*margin
		leaderWidth := (layoutWidth - 2*gap) / 5
		termWidth := (layoutWidth - 2*gap) - leaderWidth
		termHeight := m.height - (margin / 4) - 1 - gap

		var cmd tea.Cmd
		var newModel tea.Model

		if _, ok := msg.(tea.WindowSizeMsg); ok {
			newModel, cmd = m.content.Update(tea.WindowSizeMsg{
				Width:  termWidth,
				Height: termHeight,
			})
		} else {
			newModel, cmd = m.content.Update(msg)
		}

		m.content = newModel.(*TerminalModel)
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

	m.leaderboard.width = (width - 2*gap) / 5
	m.content.width = (width - 2*gap) - m.leaderboard.width
	m.content.height = height - footerHeight - gap
	m.leaderboard.height = m.content.height

	leaderboard := lipgloss.NewStyle().
		Width(m.leaderboard.width).
		Height(m.leaderboard.height).
		Render(m.leaderboard.View())

	content := lipgloss.NewStyle().
		Width(m.content.width).
		Render(m.content.View())

	footer := m.footerView(width)
	spacer := lipgloss.NewStyle().Width(gap).Render("")
	mainArea := lipgloss.JoinHorizontal(lipgloss.Top, content, spacer, leaderboard)

	ui := lipgloss.JoinVertical(lipgloss.Left, mainArea, footer)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, ui)
}
