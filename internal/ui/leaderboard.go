package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type LeaderboardModel struct {
	width       int
	height      int
	environment string
}

func (m LeaderboardModel) View() string {
	logoHeight := m.height / 3
	listHeight := m.height * 2 / 3

	list := lipgloss.NewStyle().
		Width(m.width).
		Height(listHeight).
		Padding(1, 0).
		Render("1. Alice\n2. Bob\n3. Charlie" + "lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

	logo := renderLogo(m.width, logoHeight)

	return lipgloss.JoinVertical(lipgloss.Right, logo, list)
}
