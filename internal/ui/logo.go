package ui

import (
	_ "embed"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine"
	"github.com/charmbracelet/lipgloss"
)

//go:embed assets/logo.png
var logoBytes []byte

func renderLogo(width, height int) string {
	if width <= 0 {
		return ""
	}

	logoStr, logoHeight, err := engine.RenderImage(logoBytes, width, height)
	if err != nil {
		logoStr = "Error"
		logoHeight = 1
	}

	imageBlock := lipgloss.NewStyle().
		Width(width).
		Height(logoHeight).
		Render(logoStr)

	return imageBlock + "\n"
}
