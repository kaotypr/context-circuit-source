package hosts

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
)

type claudeHost struct{}

func (claudeHost) Name() string { return "claude" }

func (claudeHost) Build(in LaunchInput) (Launch, error) {
	budget := in.Cfg.MaxBudgetUSD
	if budget == 0 {
		budget = 3
	}
	argv := []string{
		chosenCommand(in, "claude"),
		"-p", in.Prompt,
		"--model", chosenModel(in),
		"--output-format", "stream-json",
		"--verbose",
		"--dangerously-skip-permissions",
		"--strict-mcp-config",
		"--max-budget-usd", fmt.Sprintf("%g", budget),
		"--session-id", in.SessionID,
	}
	if in.Effort != "" {
		argv = append(argv, "--effort", in.Effort)
	}
	configDir := filepath.Join(in.RunRoot, "claude")
	home := filepath.Join(in.RunRoot, "claude-home")
	env := BaseEnv(in, map[string]string{
		"CLAUDE_CONFIG_DIR": configDir,
		"HOME":              home,
	})
	return Launch{
		Argv:   argv,
		Dir:    in.Workspace,
		Env:    env,
		MkDirs: []string{configDir, home},
	}, nil
}

func (claudeHost) Parse(r io.Reader) (Parsed, error) {
	p := Parsed{
		ReadsAvailable: true,
		UsageSource:    "not reported",
	}
	err := eachJSONObject(r, func(m map[string]any) {
		switch m["type"] {
		case "assistant":
			collectClaudeTools(&p, m)
		case "result":
			applyClaudeResult(&p, m)
		}
	})
	return p, err
}

func (claudeHost) AuthReady(caller map[string]string, cfg HostConfig) (bool, string) {
	return authFromEnv(caller, cfg.AuthEnv)
}

func collectClaudeTools(p *Parsed, m map[string]any) {
	var content any
	if msg, ok := m["message"].(map[string]any); ok {
		content = msg["content"]
	}
	if content == nil {
		content = m["content"]
	}
	items, ok := content.([]any)
	if !ok {
		return
	}
	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok || obj["type"] != "tool_use" {
			continue
		}
		name, _ := obj["name"].(string)
		summary := ""
		if in, ok := obj["input"].(map[string]any); ok {
			summary = firstSummary(in)
		}
		p.ToolCalls++
		p.Reads = append(p.Reads, toolLine(name, summary))
	}
}

func applyClaudeResult(p *Parsed, m map[string]any) {
	p.Terminal = true
	if id, ok := m["session_id"].(string); ok {
		p.SessionID = id
	}
	p.IsError = asBool(m["is_error"])
	if s, ok := m["result"].(string); ok {
		p.FinalReply = s
	}
	if n, ok := fieldInt(m, "num_turns"); ok {
		p.Turns = int(n)
	}
	if c, ok := fieldFloat(m, "total_cost_usd"); ok {
		p.CostUSD = &c
		p.CostBasis = "reported"
	}
	if usage, ok := m["usage"].(map[string]any); ok {
		applyClaudeUsage(p, usage)
	} else {
		p.setNull("reasoning", "not reported")
	}
	if mu, ok := m["modelUsage"].(map[string]any); ok && len(mu) > 0 {
		applyModelUsage(p, mu)
	}
}

func applyClaudeUsage(p *Parsed, usage map[string]any) {
	in, inOK := fieldInt(usage, "input_tokens")
	cc, ccOK := fieldInt(usage, "cache_creation_input_tokens")
	cr, crOK := fieldInt(usage, "cache_read_input_tokens")
	out, outOK := fieldInt(usage, "output_tokens")
	if inOK || ccOK || crOK || outOK {
		p.UsageSource = "reported"
	}
	if inOK {
		p.Tokens.InputUncached = i64p(in)
	}
	if ccOK {
		p.Tokens.CacheWrite = i64p(cc)
	}
	if crOK {
		p.Tokens.CacheRead = i64p(cr)
	}
	if outOK {
		p.Tokens.Output = i64p(out)
	}
	if inOK || ccOK || crOK {
		total := int64(0)
		if inOK {
			total += in
		}
		if ccOK {
			total += cc
		}
		if crOK {
			total += cr
		}
		p.Tokens.InputTotal = i64p(total)
	}
	// Only the direct usage.thinking_tokens field counts. Nested
	// output_tokens_details.thinking_tokens is a different key.
	if n, ok := fieldInt(usage, "thinking_tokens"); ok {
		p.Tokens.Reasoning = i64p(n)
		return
	}
	p.setNull("reasoning", "not reported")
}

func applyModelUsage(p *Parsed, mu map[string]any) {
	names := make([]string, 0, len(mu))
	var in, out, cr, cc int64
	for name, raw := range mu {
		names = append(names, name)
		obj, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		if n, ok := fieldInt(obj, "inputTokens"); ok {
			in += n
		}
		if n, ok := fieldInt(obj, "outputTokens"); ok {
			out += n
		}
		if n, ok := fieldInt(obj, "cacheReadInputTokens"); ok {
			cr += n
		}
		if n, ok := fieldInt(obj, "cacheCreationInputTokens"); ok {
			cc += n
		}
	}
	sort.Strings(names)
	p.ModelsUsed = names
	if len(names) > 0 {
		p.Model = names[0]
	}
	if len(names) > 1 {
		p.Subagents = len(names) - 1
	} else {
		p.Subagents = 0
	}
	p.Tokens.InputUncached = i64p(in)
	p.Tokens.CacheRead = i64p(cr)
	p.Tokens.CacheWrite = i64p(cc)
	p.Tokens.Output = i64p(out)
	p.Tokens.InputTotal = i64p(in + cr + cc)
	p.UsageSource = "reported"
}
