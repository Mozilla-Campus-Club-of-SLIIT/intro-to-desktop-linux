package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type LeaderboardModel struct {
	width  int
	height int
}

func (m LeaderboardModel) View() string {
	const diag = "#"

	// Header: "Live Leaderboard" + diagonal slashes
	headerText := "Live Leaderboard "
	headerTextWidth := lipgloss.Width(headerText)

	// Account for padding (1 on each side = 2 total)
	paddingWidth := 2

	// Calculate how many slashes to fill the width
	remainingWidth := m.width - headerTextWidth - paddingWidth
	remainingWidth = max(remainingWidth, 0)
	slashes := strings.Repeat(diag, remainingWidth)

	// Style the header with white text on dark green background
	headerStyle := lipgloss.NewStyle().
		Foreground(ColorWhite).
		Background(ColorDarkGreen).
		Width(m.width).
		Padding(0, 1) // vertical padding 0, horizontal padding 1

	header := headerStyle.Render(headerText + slashes)

	// Leaderboard content
	listHeight := m.height - lipgloss.Height(header) // account for header height

	content := "1. Alice\n2. Bob\n3. Charlie"

	list := lipgloss.NewStyle().
		Width(m.width).
		Height(listHeight).
		Padding(1, 0).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, header, list)
}
