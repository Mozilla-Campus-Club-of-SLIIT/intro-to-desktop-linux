package engine

import (
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
)

type Record struct {
	Timestamp time.Time `json:"timestamp"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
}

var ansiRegexp = regexp.MustCompile(`\x1b(?:\[[0-9;?]*[ -/]*[@-~]|][^\x07]*\x07|[()][A-Za-z]|[>=]|[NOc]|[\x20-\x2F]*[\x30-\x7E])`)

func CleanText(s string) string {
	s = ansiRegexp.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '\n' || ch == '\t' || ch >= 0x20 {
			b.WriteByte(ch)
		}
	}
	return strings.TrimSpace(b.String())
}

func ProcessTerminalOutput(current string, newBytes []byte) string {
	runes := []rune(current)
	for _, b := range newBytes {
		switch b {
		case 8, 127:
			if len(runes) > 0 {
				runes = runes[:len(runes)-1]
			}
		case '\r':
			continue
		default:
			runes = append(runes, rune(b))
		}
	}
	return string(runes)
}

type TerminalEngine struct {
	mu            sync.Mutex
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
	e.mu.Lock()
	defer e.mu.Unlock()

	isControl := len(b) > 0 && b[0] == 27

	for _, char := range b {
		if char == '\r' || char == '\n' {
			// 1. Flush previous record if it exists
			if e.haveCurrent {
				e.flushLocked()
			}

			// 2. Start NEW record for the command just submitted
			cmd := e.inputBuffer.String()
			e.inputBuffer.Reset()

			e.currentRecord = Record{
				Timestamp: time.Now(),
				Input:     CleanText(cmd),
				Output:    "",
			}
			e.haveCurrent = true
		} else if !isControl {
			if char == 0x7f || char == 8 {
				curr := e.inputBuffer.String()
				if len(curr) > 0 {
					r := []rune(curr)
					e.inputBuffer.Reset()
					e.inputBuffer.WriteString(string(r[:len(r)-1]))
				}
			} else if char >= 0x20 || char == '\t' {
				e.inputBuffer.WriteByte(char)
			}
		}
	}
	return e.Ptmx.Write(b)
}

func (e *TerminalEngine) Read(p []byte) (int, error) {
	n, err := e.Ptmx.Read(p)
	if n > 0 {
		e.mu.Lock()
		if e.haveCurrent {
			e.currentRecord.Output = ProcessTerminalOutput(e.currentRecord.Output, p[:n])
		}
		e.mu.Unlock()
	}
	return n, err
}

func (e *TerminalEngine) Flush() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.flushLocked()
}

func (e *TerminalEngine) flushLocked() {
	// Only log if there was actually a command typed
	if e.haveCurrent && strings.TrimSpace(e.currentRecord.Input) != "" {
		e.currentRecord.Output = CleanText(e.currentRecord.Output)
		_ = e.encoder.Encode(e.currentRecord)
		_ = e.transcript.Sync()
	}
	e.haveCurrent = false
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
