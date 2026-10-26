package hosts

import (
	"io"
	"path/filepath"
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
	env := BaseEnv(in, map[string]string{
		"HOME":            home,
		"XDG_CONFIG_HOME": config,
	})
	return Launch{
		Argv:   argv,
		Dir:    in.Workspace,
		Env:    env,
		MkDirs: []string{home, config},
	}, nil
}

func (cursorHost) Parse(r io.Reader) (Parsed, error) {
	p := Parsed{UsageSource: "not reported"}
	sawTool := false
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
			sawTool = true
			name, _ := m["name"].(string)
			summary := ""
			if in, ok := m["input"].(map[string]any); ok {
				summary = firstSummary(in)
			}
			if name == "" {
				if tc, ok := m["tool_call"].(map[string]any); ok {
					name, _ = tc["name"].(string)
					if summary == "" {
						if args, ok := tc["args"].(map[string]any); ok {
							summary = firstSummary(args)
						}
					}
				}
			}
			p.ToolCalls++
			p.Reads = append(p.Reads, toolLine(name, summary))
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
