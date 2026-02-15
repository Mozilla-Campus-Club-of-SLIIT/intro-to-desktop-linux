package ui

import (
	"strings"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/auth"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AuthModel struct {
	width               int
	height              int
	environment         string
	viewport            viewport.Model
	textinput           textinput.Model
	messages            []string
	senderStyle         lipgloss.Style
	mozStyle            lipgloss.Style
	whiteStyle          lipgloss.Style
	unknownCommandCount int
}

func NewAuthModel(env string) *AuthModel {
	const chatWidth = 80

	ti := textinput.New()
	ti.Placeholder = "Type a command ..."
	ti.Focus()
	ti.Prompt = "> "
	ti.CharLimit = 280
	ti.Width = chatWidth

	vp := viewport.New(chatWidth, 7)

	senderStyle := lipgloss.NewStyle().Foreground(ColorLemonYellow)
	mozStyle := lipgloss.NewStyle().Foreground(ColorWarmRed)
	whiteStyle := lipgloss.NewStyle().Foreground(ColorWhite)

	messages := []string{
		lipgloss.JoinHorizontal(lipgloss.Left, mozStyle.Render("moz: "), whiteStyle.Render("⚠️  SYSTEM BREACH DETECTED ...")),
		lipgloss.JoinHorizontal(lipgloss.Left, mozStyle.Render("moz: "), whiteStyle.Render("Just kidding, Welcome home Mozillian 🦊")),
		lipgloss.JoinHorizontal(lipgloss.Left, senderStyle.Render("you: "), whiteStyle.Render("sudo make me a sandwich 😎")),
		lipgloss.JoinHorizontal(lipgloss.Left, mozStyle.Render("moz: "), whiteStyle.Render("Authenticate first, then we'll talk! 🤭")),
		lipgloss.JoinHorizontal(lipgloss.Left, mozStyle.Render("moz: "), whiteStyle.Render("Just say the magic words... (hint `sudo auth`) 🤫")),
	}
	vp.SetContent(strings.Join(messages, "\n"))
	vp.GotoBottom()

	return &AuthModel{
		environment: env,
		textinput:   ti,
		viewport:    vp,
		messages:    messages,
		senderStyle: senderStyle,
		mozStyle:    mozStyle,
		whiteStyle:  whiteStyle,
	}
}

func (m *AuthModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textinput, tiCmd = m.textinput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			input := m.textinput.Value()
			m.messages = append(m.messages, lipgloss.JoinHorizontal(lipgloss.Left, m.senderStyle.Render("you: "), m.whiteStyle.Render(input)))

			if strings.ToLower(strings.TrimSpace(input)) == "sudo auth" {
				auth.AuthUser()
				m.messages = append(m.messages, lipgloss.JoinHorizontal(lipgloss.Left, m.mozStyle.Render("moz: "), m.whiteStyle.Render("Authentication successful. Shutting down for restart...")))
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
				return m, tea.Quit
			} else {
				m.unknownCommandCount++
				var unknownMsg string
				switch m.unknownCommandCount {
				case 1:
					unknownMsg = "Still waiting for the right words... 🙄"
				case 2:
					unknownMsg = "You're overthinking it! Just type the hint. 🫠"
				case 3:
					unknownMsg = "Are you doing this on purpose just to talk to me? 🤨"
				case 4:
					unknownMsg = "That's it. Go touch grass. 😤"
				default:
					unknownMsg = "Unknown Command"
				}
				m.messages = append(m.messages, lipgloss.JoinHorizontal(lipgloss.Left, m.mozStyle.Render("moz: "), m.whiteStyle.Render(unknownMsg)))
			}
			m.textinput.Reset()
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.viewport.GotoBottom()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *AuthModel) View() string {
	logo := LogoLarge()

	chatContainerStyle := lipgloss.NewStyle().Width(m.viewport.Width)

	chatView := m.viewport.View()
	inputView := m.textinput.View()

	chatContainer := chatContainerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, chatView, "\n", inputView),
	)

	content := lipgloss.JoinVertical(lipgloss.Center, logo, "\n", chatContainer)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
