package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	mozStyle  = lipgloss.NewStyle().Foreground(ColorWarmRed)
	liveStyle = lipgloss.NewStyle().Foreground(ColorLemonYellow)
)

func LogoLarge() string {
	sliitMozilla := lipgloss.NewStyle().
		Foreground(ColorLemonYellow).
		Render("by SLIIT Mozilla")

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

	live := []string{
		liveStyle.Render("████          █████████████  ███       ███ █████████████"),
		liveStyle.Render("████          █████████████  ███       ███ █████████████"),
		liveStyle.Render("████               ███       ███       ███ ████         "),
		liveStyle.Render("████               ███       ███       ███ ██████████   "),
		liveStyle.Render("████               ███       █████   █████ ██████████   "),
		liveStyle.Render("████               ███         ███   ███   ████         "),
		liveStyle.Render("█████████████ █████████████    █████████   █████████████"),
		liveStyle.Render("█████████████ █████████████       ████     █████████████"),
	}

	var logo strings.Builder
	for i := range moz {
		logo.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, moz[i], live[i]))
		logo.WriteString("\n")
	}

	return lipgloss.JoinVertical(lipgloss.Left, logo.String(), sliitMozilla)
}

func LogoMedium() string {

	sliitMozilla := lipgloss.NewStyle().
		Foreground(ColorLemonYellow).
		Render("SLIIT Mozilla")

	moz := []string{
		mozStyle.Render("███   ███  ██████  █████████"),
		mozStyle.Render("████ ████ ██    ██      ████"),
		mozStyle.Render("█████████ ██    ██    ███   "),
		mozStyle.Render("███ █ ███ ██    ██ ████     "),
		mozStyle.Render("███   ███  ██████  █████████"),
	}

	live := []string{
		liveStyle.Render(" ██       █████████ ██   ██ ██████████"),
		liveStyle.Render(" ██          ███    ██   ██ ███       "),
		liveStyle.Render(" ██          ███    ██   ██ ████████  "),
		liveStyle.Render(" ██          ███     ██ ██  ███       "),
		liveStyle.Render(" ████████ █████████   ███   ██████████"),
	}

	var logo strings.Builder
	for i := range moz {
		logo.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, moz[i], live[i]))
		logo.WriteString("\n")
	}

	return lipgloss.JoinVertical(lipgloss.Left, sliitMozilla, logo.String())
}
