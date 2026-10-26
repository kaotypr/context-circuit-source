// Package hosts builds headless launches for coding hosts and parses their streams.
package hosts

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

// Parsed is what a host stream reported. Nil token counts and a nil cost
// mean the host did not report that number.
type Parsed struct {
	SessionID      string
	Model          string
	ModelsUsed     []string
	FinalReply     string
	IsError        bool
	Tokens         record.Tokens
	UsageSource    string
	CostUSD        *float64
	CostBasis      string
	NullReasons    map[string]string
	Turns          int
	ToolCalls      int
	Subagents      int
	Reads          []string
	ReadsAvailable bool
	Terminal       bool
	Warnings       []string
}

// LaunchInput is the per-run data an adapter needs.
// Caller is the parent environment, read only for PATH and the auth
// variables named on the host config.
type LaunchInput struct {
	Prompt    string
	Model     string
	Effort    string
	Workspace string
	RunRoot   string
	SessionID string
	Cfg       HostConfig
	Caller    map[string]string
	DecoyDir  string
}

// Launch is a process the runner can start.
// Stdin "" means empty stdin. SetupArgv, when set, runs first with the
// same environment; SetupStdin is one entry per setup command.
type Launch struct {
	Argv       []string
	Dir        string
	Env        []string
	MkDirs     []string
	Stdin      string
	SetupArgv  [][]string
	SetupStdin []string
}

// Host is one coding host.
type Host interface {
	Name() string
	Build(in LaunchInput) (Launch, error)
	Parse(r io.Reader) (Parsed, error)
	AuthReady(caller map[string]string, cfg HostConfig) (bool, string)
}

// ForName returns the adapter for a hosts.yaml name.
func ForName(name string) (Host, error) {
	switch name {
	case "claude":
		return claudeHost{}, nil
	case "codex":
		return codexHost{}, nil
	case "cursor":
		return cursorHost{}, nil
	default:
		return nil, fmt.Errorf("unknown host %q", name)
	}
}

func (p *Parsed) setNull(key, reason string) {
	if p.NullReasons == nil {
		p.NullReasons = make(map[string]string)
	}
	p.NullReasons[key] = reason
}

func eachJSONObject(r io.Reader, fn func(map[string]any)) error {
	br := bufio.NewReader(r)
	for {
		line, err := readLine(br)
		if len(bytes.TrimSpace(line)) > 0 {
			trimmed := bytes.TrimSpace(line)
			if trimmed[0] == '{' {
				var m map[string]any
				if json.Unmarshal(trimmed, &m) == nil && m != nil {
					fn(m)
				}
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func readLine(br *bufio.Reader) ([]byte, error) {
	var buf []byte
	for {
		frag, err := br.ReadSlice('\n')
		buf = append(buf, frag...)
		if err == bufio.ErrBufferFull {
			continue
		}
		return buf, err
	}
}

func fieldInt(m map[string]any, key string) (int64, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			f, ferr := n.Float64()
			if ferr != nil {
				return 0, false
			}
			return int64(f), true
		}
		return i, true
	default:
		return 0, false
	}
}

func fieldFloat(m map[string]any, key string) (float64, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

func asBool(v any) bool {
	b, ok := v.(bool)
	return ok && b
}

func i64p(n int64) *int64 {
	return &n
}

func oneLine(s string, max int) string {
	s = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' {
			return ' '
		}
		return r
	}, s)
	if max > 0 && utf8.RuneCountInString(s) > max {
		s = string([]rune(s)[:max])
	}
	return s
}

func toolLine(name, summary string) string {
	summary = oneLine(summary, 200)
	if summary == "" {
		return "Tool: " + name
	}
	return "Tool: " + name + " " + summary
}

func firstSummary(in map[string]any) string {
	for _, key := range []string{"file_path", "command", "pattern", "path"} {
		if s, ok := in[key].(string); ok && s != "" {
			return s
		}
	}
	return ""
}

func chosenModel(in LaunchInput) string {
	if in.Model != "" {
		return in.Model
	}
	return in.Cfg.Model
}

func chosenCommand(in LaunchInput, fallback string) string {
	if in.Cfg.Command != "" {
		return in.Cfg.Command
	}
	return fallback
}
