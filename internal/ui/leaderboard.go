package ui

import (
	"fmt"
	"strings"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/auth"
	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/leaderboard"
	"github.com/charmbracelet/lipgloss"
)

type LeaderboardModel struct {
	width  int
	height int
}

func (m LeaderboardModel) View() string {
	const diag = "#"
	currentUser, _, _, _ := auth.GetClient().GetUserInfo()

	// 1. Header Rendering
	headerText := "Live Leaderboard " + currentUser
	headerTextWidth := lipgloss.Width(headerText)
	paddingWidth := 2
	remainingWidth := max(m.width-headerTextWidth-paddingWidth, 0)
	slashes := strings.Repeat(diag, remainingWidth)

	headerStyle := lipgloss.NewStyle().
		Foreground(ColorWhite).
		Background(ColorDarkGreen).
		Width(m.width).
		Padding(0, 1)

	header := headerStyle.Render(headerText + slashes)

	// 2. Data Preparation
	entries := leaderboard.GetLeaderboard()
	userRank := -1
	for i, entry := range entries {
		if entry.Username == currentUser {
			userRank = i
			break
		}
	}

	// 3. Layout Calculation
	// Total height - 1 (header) - 1 (top space) - 1 (bottom space)
	contentHeight := m.height - 3
	if contentHeight <= 0 {
		return header
	}

	// Column Widths
	rankWidth := 4
	scoreWidth := 8                                   // " [0000]"
	nameWidth := m.width - rankWidth - scoreWidth - 2 // -2 for internal padding

	// Helper to format a single entry as a table row
	formatRow := func(rank int, entry leaderboard.Entry) string {
		var rankStr string
		switch rank {
		case 0:
			rankStr = "🥇"
		case 1:
			rankStr = "🥈"
		case 2:
			rankStr = "🥉"
		default:
			rankStr = fmt.Sprintf("%d.", rank+1)
		}

		score := "...."
		if rank < 5 || entry.Username == currentUser {
			score = fmt.Sprintf("%04d", entry.Score)
		}

		// Build columns
		sRank := lipgloss.NewStyle().Width(rankWidth).Render(rankStr)
		sName := lipgloss.NewStyle().Width(nameWidth).Render(entry.Username)
		sScore := lipgloss.NewStyle().Width(scoreWidth).Align(lipgloss.Right).Render(fmt.Sprintf("[%s]", score))

		rowContent := lipgloss.JoinHorizontal(lipgloss.Left, sRank, sName, sScore)

		style := lipgloss.NewStyle().Width(m.width).Padding(0, 1)
		if entry.Username == currentUser {
			style = style.Background(ColorWarmRed).Foreground(ColorWhite)
		}
		return style.Render(rowContent)
	}

	// 4. Build Row List
	var rows []string

	// Add Top 3
	for i := 0; i < 3 && i < len(entries); i++ {
		rows = append(rows, formatRow(i, entries[i]))
	}

	// Add space between 3rd and 4th
	if len(entries) > 3 && len(rows) < contentHeight {
		rows = append(rows, "")
	}

	// Fill remaining slots
	remainingSlots := contentHeight - len(rows)
	if remainingSlots > 0 {
		startIdx := 3

		// Logic to ensure current user is visible if they are further down
		if userRank >= startIdx+remainingSlots {
			// Show as many as possible before the "..."
			for i := startIdx; i < startIdx+remainingSlots-2 && i < len(entries); i++ {
				rows = append(rows, formatRow(i, entries[i]))
			}
			rows = append(rows, lipgloss.NewStyle().PaddingLeft(2).Render("..."))
			rows = append(rows, formatRow(userRank, entries[userRank]))
		} else {
			// Sequential fill
			for i := startIdx; i < len(entries) && len(rows) < contentHeight; i++ {
				rows = append(rows, formatRow(i, entries[i]))
			}
		}
	}

	// Truncate if logic overflowed
	if len(rows) > contentHeight {
		rows = rows[:contentHeight]
	}

	// 5. Final Assembly
	// We use Padding(1, 0) to ensure the empty line at top and bottom
	list := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-1).
		Padding(1, 0).
		Render(strings.Join(rows, "\n"))

	return lipgloss.JoinVertical(lipgloss.Left, header, list)
}
