package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func LogoLarge() string {
	sliitMozilla := lipgloss.NewStyle().
		Foreground(ColorLemonYellow).
		Render("by SLIIT Mozilla")

	mozStyle := lipgloss.NewStyle().Foreground(ColorWarmRed)
	illaStyle := lipgloss.NewStyle().Foreground(ColorLemonYellow)

	moz := []string{
		mozStyle.Render("████      ███    ████████    █████████████ "),
		mozStyle.Render("██████  █████  ████████████  █████████████ "),
		mozStyle.Render("██████  █████  ███      ███         ██████ "),
		mozStyle.Render("█████████████  ███      ███       █████    "),
		mozStyle.Render("████ ████ ███  ███      ███    ██████      "),
		mozStyle.Render("████ ████ ███  ███      ███  ██████        "),
		mozStyle.Render("████      ███  ████████████  █████████████ "),
		mozStyle.Render("████      ███    ████████    █████████████ "),
	}

	illa := []string{
		illaStyle.Render("████          █████████████  ███       ███ █████████████"),
		illaStyle.Render("████          █████████████  ███       ███ █████████████"),
		illaStyle.Render("████               ███       ███       ███ ████         "),
		illaStyle.Render("████               ███       ███       ███ ██████████   "),
		illaStyle.Render("████               ███       █████   █████ ██████████   "),
		illaStyle.Render("████               ███         ███   ███   ████         "),
		illaStyle.Render("█████████████ █████████████    █████████   █████████████"),
		illaStyle.Render("█████████████ █████████████       ████     █████████████"),
	}

	var logo strings.Builder
	for i := range moz {
		logo.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, moz[i], illa[i]))
		logo.WriteString("\n")
	}

	return lipgloss.JoinVertical(lipgloss.Left, logo.String(), sliitMozilla)
}
