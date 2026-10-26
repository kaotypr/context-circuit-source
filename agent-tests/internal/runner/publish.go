package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

type textRedactor struct {
	home    string
	user    string
	runRoot string
}

func (rt *runtime) redact() textRedactor {
	home, user := "", ""
	if rt.opts.Caller != nil {
		home = rt.opts.Caller["HOME"]
		user = rt.opts.Caller["USER"]
		if user == "" {
			user = rt.opts.Caller["LOGNAME"]
		}
	}
	return textRedactor{home: home, user: user}
}

func (r textRedactor) text(s string) string {
	s = record.Redact(s, r.home, r.user)
	if r.runRoot != "" {
		s = record.RedactRunRoot(s, r.runRoot)
	}
	return s
}

func (rt *runtime) publishHost(runDir string, w world.World, res proc.Result, streamPath string, parsed *hosts.Parsed, drops *[]string) error {
	red := rt.redact()
	red.runRoot = w.RunRoot
	stream, err := os.ReadFile(streamPath)
	if err != nil {
		stream = res.Stdout
	}
	if err := writeBytes(filepath.Join(runDir, "stream.jsonl"), stream, red, drops); err != nil {
		return err
	}
	if err := writeBytes(filepath.Join(runDir, "stderr.txt"), res.Stderr, red, drops); err != nil {
		return err
	}
	reply := ""
	reads := "reads list not available"
	if parsed != nil {
		reply = parsed.FinalReply
		if parsed.ReadsAvailable {
			reads = strings.Join(parsed.Reads, "\n")
			if reads != "" {
				reads += "\n"
			}
		}
	}
	if err := writeText(filepath.Join(runDir, "reply.md"), reply, red); err != nil {
		return err
	}
	if err := writeText(filepath.Join(runDir, "reads.txt"), reads, red); err != nil {
		return err
	}
	decoy, err := os.ReadFile(w.DecoyLog)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return writeBytes(filepath.Join(runDir, "decoy.log"), decoy, red, drops)
}

func (rt *runtime) copyChanged(ctx context.Context, w world.World, dest string, drops *[]string) error {
	red := rt.redact()
	red.runRoot = w.RunRoot
	env := snippetEnv(w, callerPath(rt.opts.Caller))
	repos := []string{w.WS}
	if w.API != "" {
		repos = append(repos, w.API)
	}
	for _, repo := range repos {
		if repo == "" {
			continue
		}
		if _, statErr := os.Lstat(filepath.Join(repo, ".git")); statErr != nil {
			continue
		}
		paths, err := changedPaths(ctx, env, repo)
		if err != nil {
			return err
		}
		relRoot, err := filepath.Rel(w.WS, repo)
		if err != nil {
			return err
		}
		if relRoot == "." {
			relRoot = ""
		}
		for _, rel := range paths {
			src := filepath.Join(repo, rel)
			info, err := os.Lstat(src)
			if err != nil || !info.Mode().IsRegular() {
				continue
			}
			target := filepath.Join(dest, relRoot, rel)
			b, err := os.ReadFile(src)
			if err != nil {
				return err
			}
			if err := writeBytes(target, b, red, drops); err != nil {
				return err
			}
		}
	}
	return nil
}

func changedPaths(ctx context.Context, env []string, repo string) ([]string, error) {
	res, err := proc.Run(ctx, proc.Spec{
		Argv: []string{"git", "-C", repo, "status", "--porcelain", "-z"},
		Env:  env,
	})
	if err != nil {
		return nil, err
	}
	if res.ExitCode != 0 {
		return nil, fmt.Errorf("git status: exit %d\n%s", res.ExitCode, res.Stderr)
	}
	parts := strings.Split(string(res.Stdout), "\x00")
	var paths []string
	for i := 0; i < len(parts); i++ {
		entry := parts[i]
		if len(entry) < 4 {
			continue
		}
		status := entry[:2]
		path := entry[3:]
		if status[0] == 'R' || status[0] == 'C' || status[1] == 'R' || status[1] == 'C' {
			i++
			if i < len(parts) {
				path = parts[i]
			}
		}
		if path != "" {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

func (rt *runtime) copyDesign(w world.World, sc scenario.Scenario, dest string, drops *[]string) error {
	red := rt.redact()
	red.runRoot = w.RunRoot
	for _, rel := range sc.DesignArtifacts {
		src := filepath.Join(w.WS, rel)
		b, err := os.ReadFile(src)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		if err := writeBytes(filepath.Join(dest, rel), b, red, drops); err != nil {
			return err
		}
	}
	return nil
}

func writeText(path, text string, red textRedactor) error {
	return writeBytes(path, []byte(text), red, nil)
}

func writeBytes(path string, data []byte, red textRedactor, drops *[]string) error {
	base := filepath.Base(path)
	if secretName(base) {
		if drops != nil {
			*drops = append(*drops, path+" name")
		}
		return nil
	}
	// Judge the original bytes. Redaction would hide the same patterns, and
	// a file that still looks like a secret must not be copied at all.
	if record.LooksLikeSecret(string(data)) {
		if drops != nil {
			*drops = append(*drops, path+" content")
		}
		return nil
	}
	text := red.text(string(data))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(text), 0o644)
}

func secretName(name string) bool {
	lower := strings.ToLower(name)
	if lower == ".credentials.json" {
		return true
	}
	return strings.HasPrefix(lower, "auth") || strings.HasPrefix(lower, "credentials")
}
