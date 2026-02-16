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
	const leftFieldWidth = 5

	redStyle := lipgloss.NewStyle().Foreground(ColorWarmRed)
	yellowStyle := lipgloss.NewStyle().Foreground(ColorLemonYellow)

	logo := LogoMedium()
	logoLines := strings.Split(strings.TrimSpace(logo), "\n")

	var decoratedLogo strings.Builder

	for i, line := range logoLines {
		leftField := redStyle.Render(strings.Repeat(diag, leftFieldWidth))

		totalWidth := m.width
		if totalWidth <= 0 {
			totalWidth = 100
		}

		lineWidth := lipgloss.Width(line)
		rightFieldWidth := totalWidth - leftFieldWidth - lineWidth - 2
		rightFieldWidth = max(rightFieldWidth, 5)

		rightField := yellowStyle.Render(strings.Repeat(diag, rightFieldWidth))

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

	return decoratedLogo.String() + contentText
}
