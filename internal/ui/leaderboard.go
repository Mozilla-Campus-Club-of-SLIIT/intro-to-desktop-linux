package ui

import (
	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/auth"
	"github.com/charmbracelet/lipgloss"
)

type LeaderboardModel struct {
	width  int
	height int
}

func (m LeaderboardModel) View() string {
	listHeight := m.height * 2 / 3

	content := "1. Alice\n2. Bob\n3. Charlie"

	// Example: Show different content for session hosts
	if auth.IsSessionHost() {
		content += "\n\n[Session Host Mode] 👑"
	}

	list := lipgloss.NewStyle().
		Width(m.width).
		Height(listHeight).
		Padding(1, 0).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Right, list)
}
