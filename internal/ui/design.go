package ui

import "github.com/charmbracelet/lipgloss"

var (
	ColorBlack       = lipgloss.Color("#000000")
	ColorWhite       = lipgloss.Color("#FFFFFF")
	ColorLemonYellow = lipgloss.Color("#FFF44F")
	ColorWarmRed     = lipgloss.Color("#FF4F5E")
	ColorDarkGreen   = lipgloss.Color("#005E5E")

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorBlack).
			Background(ColorLemonYellow).
			Padding(0, 1)
)
