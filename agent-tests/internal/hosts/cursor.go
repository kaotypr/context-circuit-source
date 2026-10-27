package hosts

import (
	"io"
	"path/filepath"
	"sort"
	"strings"
)

type cursorHost struct{}

func (cursorHost) Name() string { return "cursor" }

func (cursorHost) Build(in LaunchInput) (Launch, error) {
	argv := []string{
		chosenCommand(in, "cursor-agent"),
		"-p", in.Prompt,
		"--model", chosenModel(in),
		"--output-format", "stream-json",
		"--force",
		"--trust",
		"--workspace", in.Workspace,
	}
	home := filepath.Join(in.RunRoot, "cursor-home")
	config := filepath.Join(home, ".config")
	configDir := filepath.Join(config, "cursor")
	dataDir := filepath.Join(home, ".cursor")
	// cursor-agent resolves its config, data, plugin, hook and project
	// directories from HOME, XDG_CONFIG_HOME, CURSOR_CONFIG_DIR and
	// CURSOR_DATA_DIR; all four point into the run root. On macOS the
	// default credential store is the login keychain, written through
	// /usr/bin/security. The memory store keeps the API key's session in
	// the process instead, so no keychain is read or written.
	env := BaseEnv(in, map[string]string{
		"HOME":                       home,
		"XDG_CONFIG_HOME":            config,
		"CURSOR_CONFIG_DIR":          configDir,
		"CURSOR_DATA_DIR":            dataDir,
		"AGENT_CLI_CREDENTIAL_STORE": "memory",
	})
	return Launch{
		Argv:   argv,
		Dir:    in.Workspace,
		Env:    env,
		MkDirs: []string{home, config, configDir, dataDir},
	}, nil
}

func (cursorHost) Parse(r io.Reader) (Parsed, error) {
	p := Parsed{UsageSource: "not reported"}
	sawTool := false
	seen := map[string]bool{}
	err := eachJSONObject(r, func(m map[string]any) {
		switch m["type"] {
		case "system":
			if m["subtype"] != "init" {
				return
			}
			if s, ok := m["model"].(string); ok {
				p.Model = s
			}
			if s, ok := m["session_id"].(string); ok {
				p.SessionID = s
			}
		case "tool_call":
			if line, ok := cursorToolLine(m, seen); ok {
				sawTool = true
				p.ToolCalls++
				p.Reads = append(p.Reads, line)
			}
		case "result":
			applyCursorResult(&p, m)
		}
	})
	p.ReadsAvailable = sawTool
	if !sawTool {
		p.setNull("reads", "reads list not available")
	}
	return p, err
}

func (cursorHost) AuthReady(caller map[string]string, cfg HostConfig) (bool, string) {
	return authFromEnv(caller, cfg.AuthEnv)
}

func applyCursorResult(p *Parsed, m map[string]any) {
	p.Terminal = true
	p.IsError = asBool(m["is_error"])
	if s, ok := m["result"].(string); ok {
		p.FinalReply = s
	}
	if s, ok := m["session_id"].(string); ok && s != "" {
		p.SessionID = s
	}
	usage, ok := m["usage"].(map[string]any)
	if !ok {
		return
	}
	in, inOK := fieldInt(usage, "inputTokens")
	out, outOK := fieldInt(usage, "outputTokens")
	cr, crOK := fieldInt(usage, "cacheReadTokens")
	cw, cwOK := fieldInt(usage, "cacheWriteTokens")
	if inOK || outOK || crOK || cwOK {
		p.UsageSource = "reported"
	}
	if inOK {
		p.Tokens.InputUncached = i64p(in)
	}
	if outOK {
		p.Tokens.Output = i64p(out)
	}
	if crOK {
		p.Tokens.CacheRead = i64p(cr)
	}
	if cwOK {
		p.Tokens.CacheWrite = i64p(cw)
	}
	if inOK || crOK || cwOK {
		total := int64(0)
		if inOK {
			total += in
		}
		if crOK {
			total += cr
		}
		if cwOK {
			total += cw
		}
		p.Tokens.InputTotal = i64p(total)
	}
}

// cursorToolLine turns one tool_call event into a reads line. cursor-agent
// sends a started and a completed event per call, sharing call_id, with the
// call under tool_call.<kind>ToolCall.args. A call is counted once: on the
// first of its events seen. An event without a call_id is counted on its
// own. The older flat shape, name and input on the event, is still read.
func cursorToolLine(m map[string]any, seen map[string]bool) (string, bool) {
	if id, ok := m["call_id"].(string); ok && id != "" {
		if seen[id] {
			return "", false
		}
		seen[id] = true
	}
	if name, ok := m["name"].(string); ok && name != "" {
		summary := ""
		if in, ok := m["input"].(map[string]any); ok {
			summary = firstSummary(in)
		}
		return toolLine(name, summary), true
	}
	tc, ok := m["tool_call"].(map[string]any)
	if !ok {
		return toolLine("", ""), true
	}
	if name, ok := tc["name"].(string); ok && name != "" {
		summary := ""
		if args, ok := tc["args"].(map[string]any); ok {
			summary = firstSummary(args)
		}
		return toolLine(name, summary), true
	}
	keys := make([]string, 0, len(tc))
	for key, val := range tc {
		if _, ok := val.(map[string]any); ok && strings.HasSuffix(key, "ToolCall") {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return toolLine("", ""), true
	}
	sort.Strings(keys)
	key := keys[0]
	name := strings.TrimSuffix(key, "ToolCall")
	call := tc[key].(map[string]any)
	args, _ := call["args"].(map[string]any)
	return toolLine(name, cursorSummary(name, args)), true
}

// cursorSummary is the path, command or pattern of one Cursor tool call.
func cursorSummary(name string, args map[string]any) string {
	if args == nil {
		return ""
	}
	str := func(key string) string {
		s, _ := args[key].(string)
		return s
	}
	join := func(a, b string) string {
		switch {
		case a == "":
			return b
		case b == "":
			return a
		default:
			return a + " " + b
		}
	}
	switch name {
	case "shell":
		return str("command")
	case "grep":
		return join(str("pattern"), str("path"))
	case "glob":
		return join(str("globPattern"), str("targetDirectory"))
	}
	if s := firstSummary(args); s != "" {
		return s
	}
	return join(str("globPattern"), str("targetDirectory"))
}
