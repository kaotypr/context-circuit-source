package hosts

import (
	"io"
	"path/filepath"
	"strings"
)

type codexHost struct{}

func (codexHost) Name() string { return "codex" }

func (codexHost) Build(in LaunchInput) (Launch, error) {
	cmd := chosenCommand(in, "codex")
	gocache := filepath.Join(in.RunRoot, "gocache")
	argv := []string{
		cmd,
		"exec",
		"--json",
		"--ignore-user-config",
		"-C", in.Workspace,
		"-m", chosenModel(in),
		"--sandbox", "workspace-write",
		"--add-dir", gocache,
	}
	if in.Effort != "" {
		argv = append(argv, "-c", `model_reasoning_effort="`+in.Effort+`"`)
	}
	argv = append(argv, in.Prompt)

	// The API key is passed only on the login command's stdin, never in
	// argv or in the environment.
	caller := cloneMap(in.Caller)
	var setup [][]string
	var setupIn []string
	for _, name := range in.Cfg.AuthEnv {
		v := caller[name]
		delete(caller, name)
		if v != "" && len(setup) == 0 {
			setup = [][]string{{cmd, "login", "--with-api-key"}}
			setupIn = []string{v}
		}
	}
	in.Caller = caller

	codexHome := filepath.Join(in.RunRoot, "codex")
	home := filepath.Join(in.RunRoot, "codex-home")
	env := BaseEnv(in, map[string]string{
		"CODEX_HOME": codexHome,
		"HOME":       home,
	})
	return Launch{
		Argv:       argv,
		Dir:        in.Workspace,
		Env:        env,
		MkDirs:     []string{codexHome, home},
		Stdin:      "",
		SetupArgv:  setup,
		SetupStdin: setupIn,
	}, nil
}

func (codexHost) Parse(r io.Reader) (Parsed, error) {
	p := Parsed{
		ReadsAvailable: true,
		UsageSource:    "not reported",
	}
	p.setNull("model", "not in the stream; pinned model is used")

	var in, cached, out, reason, cacheWrite counted
	err := eachJSONObject(r, func(m map[string]any) {
		switch m["type"] {
		case "thread.started":
			if id, ok := m["thread_id"].(string); ok {
				p.SessionID = id
			}
		case "item.completed":
			applyCodexItem(&p, m)
		case "turn.completed":
			p.Turns++
			p.Terminal = true
			usage, ok := m["usage"].(map[string]any)
			if !ok {
				return
			}
			in.add(usage, "input_tokens")
			cached.add(usage, "cached_input_tokens")
			out.add(usage, "output_tokens")
			reason.add(usage, "reasoning_output_tokens")
			cacheWrite.add(usage, "cache_write_input_tokens")
		}
	})
	if in.saw || cached.saw || out.saw || reason.saw || cacheWrite.saw {
		p.UsageSource = "reported"
	}
	if in.saw {
		// cached_input_tokens are inside input_tokens (to be confirmed in the pilot).
		uncached := in.n
		if cached.saw {
			uncached -= cached.n
		}
		if cacheWrite.saw {
			uncached -= cacheWrite.n
		}
		p.Tokens.InputTotal = i64p(in.n)
		p.Tokens.InputUncached = i64p(uncached)
	}
	if cached.saw {
		p.Tokens.CacheRead = i64p(cached.n)
	}
	if cacheWrite.saw {
		p.Tokens.CacheWrite = i64p(cacheWrite.n)
	}
	if out.saw {
		p.Tokens.Output = i64p(out.n)
	}
	if reason.saw {
		p.Tokens.Reasoning = i64p(reason.n)
	}
	return p, err
}

func (codexHost) AuthReady(caller map[string]string, cfg HostConfig) (bool, string) {
	return authFromEnv(caller, cfg.AuthEnv)
}

type counted struct {
	n   int64
	saw bool
}

func (c *counted) add(m map[string]any, key string) {
	n, ok := fieldInt(m, key)
	if !ok {
		return
	}
	c.n += n
	c.saw = true
}

func applyCodexItem(p *Parsed, m map[string]any) {
	item, ok := m["item"].(map[string]any)
	if !ok {
		return
	}
	typ, _ := item["type"].(string)
	switch typ {
	case "agent_message":
		if text, ok := item["text"].(string); ok {
			p.FinalReply = text
		}
	case "error":
		if msg, ok := item["message"].(string); ok {
			p.Warnings = append(p.Warnings, msg)
		}
	default:
		if codexToolType(typ) {
			p.ToolCalls++
			detail, _ := item["command"].(string)
			if detail == "" {
				if path, ok := item["path"].(string); ok {
					detail = path
				}
			}
			p.Reads = append(p.Reads, toolLine(typ, detail))
		}
	}
}

func codexToolType(typ string) bool {
	switch typ {
	case "command_execution", "file_change":
		return true
	default:
		return strings.Contains(typ, "file_change")
	}
}
