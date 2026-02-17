package engine

import (
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/creack/pty"
)

type Record struct {
	Timestamp time.Time `json:"timestamp"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
}

var ansiRegexp = regexp.MustCompile(`\x1b(?:\[[0-9;?]*[ -/]*[@-~]|][^\x07]*\x07|[()][A-Za-z]|[>=]|[NOc]|[\x20-\x2F]*[\x30-\x7E])`)

// ProcessTerminalOutput handles backspaces and carriage returns to keep the buffer clean.
func ProcessTerminalOutput(current string, newBytes []byte) string {
	runes := []rune(current)
	for _, b := range newBytes {
		switch b {
		case 8, 127: // Backspace or Delete
			if len(runes) > 0 {
				runes = runes[:len(runes)-1]
			}
		case '\r':
			// Carriage returns usually mean the cursor moves to start of line,
			// but in our simple buffer we'll treat it as a newline or ignore it.
			continue
		default:
			runes = append(runes, rune(b))
		}
	}
	return string(runes)
}

func CleanText(s string) string {
	// 1. Strip ANSI sequences first
	s = ansiRegexp.ReplaceAllString(s, "")
	// 2. We keep the string as-is because ProcessTerminalOutput handles the logic
	return s
}

type TerminalEngine struct {
	Ptmx          *os.File
	cmd           *exec.Cmd
	transcript    *os.File
	encoder       *json.Encoder
	currentRecord Record
	haveCurrent   bool
	inputBuffer   strings.Builder
}

func NewTerminalEngine(transcriptPath string) (*TerminalEngine, error) {
	c := exec.Command("zsh", "-i")
	ptmx, err := pty.Start(c)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(transcriptPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		_ = ptmx.Close()
		return nil, err
	}

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)

	return &TerminalEngine{
		Ptmx:       ptmx,
		cmd:        c,
		transcript: f,
		encoder:    enc,
	}, nil
}

func (e *TerminalEngine) Write(b []byte) (int, error) {
	for _, char := range b {
		if char == '\r' || char == '\n' {
			e.Flush()
			cmd := e.inputBuffer.String()
			e.inputBuffer.Reset()
			e.currentRecord = Record{
				Timestamp: time.Now(),
				Input:     CleanText(cmd),
				Output:    "",
			}
			e.haveCurrent = true
		} else if char == 0x7f {
			curr := e.inputBuffer.String()
			if len(curr) > 0 {
				e.inputBuffer.Reset()
				e.inputBuffer.WriteString(curr[:len(curr)-1])
			}
		} else if char >= 0x20 || char == '\t' {
			e.inputBuffer.WriteByte(char)
		}
	}
	return e.Ptmx.Write(b)
}

func (e *TerminalEngine) Read(p []byte) (int, error) {
	n, err := e.Ptmx.Read(p)
	if n > 0 && e.haveCurrent {
		// This is where logging happens
		e.currentRecord.Output += CleanText(string(p[:n]))
	}
	return n, err
}

func (e *TerminalEngine) Flush() {
	if e.haveCurrent {
		_ = e.encoder.Encode(e.currentRecord)
		_ = e.transcript.Sync()
		e.haveCurrent = false
	}
}

func (e *TerminalEngine) Resize(w, h int) {
	if w <= 0 || h <= 0 {
		return
	}
	_ = pty.Setsize(e.Ptmx, &pty.Winsize{Rows: uint16(h), Cols: uint16(w)})
}

func (e *TerminalEngine) Close() error {
	e.Flush()
	_ = e.transcript.Close()
	return e.Ptmx.Close()
}
