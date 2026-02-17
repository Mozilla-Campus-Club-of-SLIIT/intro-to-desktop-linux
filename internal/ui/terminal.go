package ui

import (
	"strings"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TerminalMsg string

type TerminalModel struct {
	width  int
	height int
	engine *engine.TerminalEngine
	output string
}

func NewTerminalModel() (*TerminalModel, error) {
	eng, err := engine.NewTerminalEngine("transcript.ndjson")
	if err != nil {
		return nil, err
	}
	return &TerminalModel{
		engine: eng,
	}, nil
}

func (m *TerminalModel) Init() tea.Cmd {
	return m.listenToPTY()
}

func (m *TerminalModel) listenToPTY() tea.Cmd {
	return func() tea.Msg {
		buf := make([]byte, 2048)
		n, err := m.engine.Read(buf)
		if n > 0 {
			return TerminalMsg(string(buf[:n]))
		}
		if err != nil {
			return nil
		}
		return nil
	}
}

func (m *TerminalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		logoHeight := 10 // Approximation for the shell's internal wrapping
		m.engine.Resize(m.width, m.height-logoHeight)
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.engine.Close()
			return m, tea.Quit
		}

		var b []byte
		switch msg.Type {
		case tea.KeyEnter:
			b = []byte{'\r'}
		case tea.KeyBackspace:
			b = []byte{127}
		case tea.KeyTab:
			b = []byte{'\t'}
		case tea.KeyEsc:
			b = []byte{27}
		// Arrow Keys mappings for Shell History/Navigation
		case tea.KeyUp:
			b = []byte("\x1b[A")
		case tea.KeyDown:
			b = []byte("\x1b[B")
		case tea.KeyRight:
			b = []byte("\x1b[C")
		case tea.KeyLeft:
			b = []byte("\x1b[D")
		case tea.KeySpace:
			b = []byte{' '}
		case tea.KeyRunes:
			b = []byte(string(msg.Runes))
		default:
			s := msg.String()
			if len(s) == 1 {
				b = []byte(s)
			}
		}

		if len(b) > 0 {
			_, _ = m.engine.Write(b)
		}
		return m, nil

	case TerminalMsg:
		if string(msg) == "" {
			return m, m.listenToPTY()
		}
		// Process backspaces and raw bytes from the shell
		m.output = engine.ProcessTerminalOutput(m.output, []byte(string(msg)))
		return m, m.listenToPTY()
	}
	return m, nil
}

func (m *TerminalModel) View() string {
	const diag = "╱"
	const leftFieldWidth = 5
	redStyle := lipgloss.NewStyle().Foreground(ColorWarmRed)
	yellowStyle := lipgloss.NewStyle().Foreground(ColorLemonYellow)

	// 1. Render Logo
	logo := LogoMedium()
	logoLines := strings.Split(strings.TrimSpace(logo), "\n")
	var decoratedLogo strings.Builder
	for i, line := range logoLines {
		leftField := redStyle.Render(strings.Repeat(diag, leftFieldWidth))
		totalWidth := m.width
		if totalWidth <= 0 {
			totalWidth = 80
		}
		lineWidth := lipgloss.Width(line)
		rightFieldWidth := max(totalWidth-leftFieldWidth-lineWidth-2, 5)
		rightField := yellowStyle.Render(strings.Repeat(diag, rightFieldWidth))
		decoratedLogo.WriteString(leftField + " " + line + " " + rightField)
		if i < len(logoLines)-1 {
			decoratedLogo.WriteString("\n")
		}
	}

	// 2. Available Height
	logoHeight := lipgloss.Height(decoratedLogo.String())
	terminalHeight := m.height - logoHeight - 2 // 2 lines for spacers
	if terminalHeight < 0 {
		terminalHeight = 0
	}

	// 3. Clean and Wrap Text
	// We MUST clean ANSI escapes so the shell doesn't wipe the Logo
	cleanContent := engine.CleanText(m.output)
	wrapped := lipgloss.NewStyle().Width(m.width).Render(cleanContent)
	lines := strings.Split(wrapped, "\n")

	// 4. Truncate from the TOP (Keep bottom lines)
	if len(lines) > terminalHeight && terminalHeight > 0 {
		lines = lines[len(lines)-terminalHeight:]
	} else if terminalHeight == 0 {
		lines = []string{}
	}

	// Fill remaining space with empty lines to keep UI stable
	for len(lines) < terminalHeight {
		lines = append(lines, "")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		decoratedLogo.String(),
		" ",
		strings.Join(lines, "\n"),
		" ",
	)
}
