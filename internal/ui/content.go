package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ContentModel struct {
	width       int
	height      int
	environment string
}

func (m ContentModel) View() string {
	const diag = "╱"
	const leftFieldWidth = 5 // 5 diagonal slashes on the left

	// Styles
	redStyle := lipgloss.NewStyle().Foreground(ColorWarmRed)
	yellowStyle := lipgloss.NewStyle().Foreground(ColorLemonYellow)

	// Get the logo
	logo := LogoMedium()
	logoLines := strings.Split(strings.TrimSpace(logo), "\n")

	// Build the logo with diagonal fields
	var decoratedLogo strings.Builder

	for i, line := range logoLines {
		// Left field (red diagonal slashes)
		leftField := redStyle.Render(strings.Repeat(diag, leftFieldWidth))

		// Calculate right field width to fill remaining space
		// Use m.width if available, otherwise use a default
		totalWidth := m.width
		if totalWidth <= 0 {
			totalWidth = 100 // default width
		}

		// Account for the width of the line (with ANSI codes stripped for measurement)
		lineWidth := lipgloss.Width(line)
		rightFieldWidth := totalWidth - leftFieldWidth - lineWidth - 2 // 2 for spaces
		rightFieldWidth = max(rightFieldWidth, 5)                      // minimum width

		// Right field (yellow diagonal slashes)
		rightField := yellowStyle.Render(strings.Repeat(diag, rightFieldWidth))

		// Combine: left field + space + logo line + space + right field
		decoratedLogo.WriteString(leftField)
		decoratedLogo.WriteString(" ")
		decoratedLogo.WriteString(line)
		decoratedLogo.WriteString(" ")
		decoratedLogo.WriteString(rightField)

		if i < len(logoLines)-1 {
			decoratedLogo.WriteString("\n")
		}
	}

	// Main content text
	contentText := "\n\nMain Content Area\nEnvironment: " + m.environment +
		"\n[This is where the main content of the application would be displayed.]\n" +
		"lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	// Combine decorated logo with content
	return decoratedLogo.String() + contentText
}
