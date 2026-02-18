package ui

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/auth"
	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/leaderboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MsgLeaderboardUpdate is a message sent when the leaderboard updates.
type MsgLeaderboardUpdate []leaderboard.Entry

// MsgLeaderboardError is a message sent when an error occurs in the leaderboard stream.
type MsgLeaderboardError error

type LeaderboardModel struct {
	width  int
	height int

	leaderboardClient  *leaderboard.LeaderboardClient
	currentLeaderboard []leaderboard.Entry
	leaderboardErr     error

	ctx    context.Context
	cancel context.CancelFunc // For cancelling the gRPC stream

	userID      string
	accessToken string
	currentUser string // For display in header
}

// NewLeaderboardModel creates and initializes a new LeaderboardModel.
func NewLeaderboardModel(width, height int) LeaderboardModel {
	// Initialize with dummy values, actual client and stream started in Init()
	return LeaderboardModel{
		width:              width,
		height:             height,
		currentLeaderboard: []leaderboard.Entry{},
	}
}

// Init initializes the model and starts the leaderboard stream.
func (m LeaderboardModel) Init() tea.Cmd {
	authClient := auth.GetClient()

	// Get currentUser (username) for display
	var err error
	m.currentUser, _, _, err = authClient.GetUserInfo()
	if err != nil {
		m.leaderboardErr = fmt.Errorf("failed to get user info: %w", err)
		return nil
	}

	m.userID, err = authClient.GetUserID()
	if err != nil {
		m.leaderboardErr = fmt.Errorf("failed to get user ID: %w", err)
		return nil
	}

	m.accessToken, err = authClient.GetAccessToken()
	if err != nil {
		m.leaderboardErr = fmt.Errorf("failed to get access token: %w", err)
		return nil
	}

	// Initialize context for the gRPC stream
	m.ctx, m.cancel = context.WithCancel(context.Background())

	m.leaderboardClient, err = leaderboard.NewLeaderboardClient(m.ctx)
	if err != nil {
		m.leaderboardErr = fmt.Errorf("failed to create leaderboard client: %w", err)
		return nil // No command, as client failed
	}

	return m.listenForLeaderboardUpdates()
}

// listenForLeaderboardUpdates starts a goroutine to listen for leaderboard updates
// and sends them as tea.Msg to the main Update loop.
func (m LeaderboardModel) listenForLeaderboardUpdates() tea.Cmd {
	updates, errs := m.leaderboardClient.GetLeaderboardStream(m.ctx, m.userID, m.accessToken)

	return func() tea.Msg {
		for {
			select {
			case <-m.ctx.Done():
				// Context cancelled, clean up client connection
				if m.leaderboardClient != nil {
					m.leaderboardClient.Close()
				}
				return nil // Stop sending messages
			case update, ok := <-updates:
				if !ok {
					log.Println("Leaderboard update channel closed.")
					if m.leaderboardClient != nil {
						m.leaderboardClient.Close()
					}
					return MsgLeaderboardError(fmt.Errorf("leaderboard stream closed unexpectedly"))
				}
				return MsgLeaderboardUpdate(update)
			case err, ok := <-errs:
				if !ok {
					log.Println("Leaderboard error channel closed.")
					if m.leaderboardClient != nil {
						m.leaderboardClient.Close()
					}
					return MsgLeaderboardError(fmt.Errorf("leaderboard error stream closed unexpectedly"))
				}
				if err != nil {
					log.Printf("Error in leaderboard stream: %v", err)
					// Don't close client here, stream might recover or next update will fail
					return MsgLeaderboardError(err)
				}
			}
		}
	}
}

// Update handles messages.
func (m LeaderboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.cancel != nil {
				m.cancel() // Cancel context to stop gRPC stream
			}
			return m, tea.Quit
		}
	case MsgLeaderboardUpdate:
		m.currentLeaderboard = msg
		m.leaderboardErr = nil                    // Clear any previous errors
		return m, m.listenForLeaderboardUpdates() // Keep listening for more updates
	case MsgLeaderboardError:
		m.leaderboardErr = msg
		// Maybe retry or just display error. For now, just display.
		return m, nil
	}
	return m, nil
}

// View renders the leaderboard UI.
func (m LeaderboardModel) View() string {
	const diag = "#"

	// 1. Header Rendering
	headerText := "Live Leaderboard " + m.currentUser
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

	// Display error if present
	if m.leaderboardErr != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			Width(m.width - 4)
		return lipgloss.JoinVertical(lipgloss.Left, header, errorStyle.Render(fmt.Sprintf("Error: %v", m.leaderboardErr)))
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

		scoreStr := fmt.Sprintf("%04.0f", entry.Score)                        // Format float64 score
		if entry.Score == 0 && rank >= 5 && entry.Username != m.currentUser { // Don't show actual score if it's 0 and not current user and not top 5
			scoreStr = "...."
		}

		// Build columns
		sRank := lipgloss.NewStyle().Width(rankWidth).Render(rankStr)
		sName := lipgloss.NewStyle().Width(nameWidth).Render(entry.Username)
		sScore := lipgloss.NewStyle().Width(scoreWidth).Align(lipgloss.Right).Render(fmt.Sprintf("[%s]", scoreStr))

		rowContent := lipgloss.JoinHorizontal(lipgloss.Left, sRank, sName, sScore)

		style := lipgloss.NewStyle().Width(m.width).Padding(0, 1) // Corrected: NewInStyle -> NewStyle()
		if entry.Username == m.currentUser {
			style = style.Background(ColorWarmRed).Foreground(ColorWhite)
		}
		return style.Render(rowContent)
	}

	// formatLoadingRow creates a row with "..." placeholders for loading state.
	formatLoadingRow := func() string {
		sRank := lipgloss.NewStyle().Width(rankWidth).Render("...")
		sName := lipgloss.NewStyle().Width(nameWidth).Render("...")
		sScore := lipgloss.NewStyle().Width(scoreWidth).Align(lipgloss.Right).Render("[...]")
		rowContent := lipgloss.JoinHorizontal(lipgloss.Left, sRank, sName, sScore)
		return lipgloss.NewStyle().Width(m.width).Padding(0, 1).Render(rowContent)
	}

	// 2. Data Preparation - now using m.currentLeaderboard
	entries := m.currentLeaderboard
	if len(entries) == 0 {
		var loadingRows []string
		contentHeight := m.height - 3
		if contentHeight > 0 {
			for i := 0; i < contentHeight; i++ {
				loadingRows = append(loadingRows, formatLoadingRow())
			}
		}
		list := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height-1).
			Padding(1, 0).
			Render(strings.Join(loadingRows, "\n"))
		return lipgloss.JoinVertical(lipgloss.Left, header, list)
	}

	// 3. Layout Calculation
	// Total height - 1 (header) - 1 (top space) - 1 (bottom space)
	contentHeight := m.height - 3
	if contentHeight <= 0 {
		return header
	}

	userRank := -1
	for i, entry := range entries {
		if entry.Username == m.currentUser {
			userRank = i
			break
		}
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
			if len(entries) > startIdx+remainingSlots-2 { // Add ... only if there are more entries
				rows = append(rows, lipgloss.NewStyle().PaddingLeft(2).Render("..."))
			}
			if userRank != -1 { // Only add user's row if found
				rows = append(rows, formatRow(userRank, entries[userRank]))
			}
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
	list := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-1).
		Padding(1, 0).
		Render(strings.Join(rows, "\n"))

	return lipgloss.JoinVertical(lipgloss.Left, header, list)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
