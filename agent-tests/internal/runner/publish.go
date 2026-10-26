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
	// secrets are literal values of the caller's key and token variables.
	secrets []string
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
	return textRedactor{home: home, user: user, secrets: secretValues(rt.opts.Caller)}
}

// secretValues returns the values of caller variables that look like keys or
// tokens, so the exact string is removed wherever a host echoed it.
func secretValues(caller map[string]string) []string {
	var out []string
	for name, val := range caller {
		upper := strings.ToUpper(name)
		if len(val) < 8 {
			continue
		}
		if strings.Contains(upper, "KEY") || strings.Contains(upper, "TOKEN") || strings.Contains(upper, "SECRET") {
			out = append(out, val)
		}
	}
	return out
}

func (r textRedactor) text(s string) string {
	for _, secret := range r.secrets {
		s = strings.ReplaceAll(s, secret, "[redacted]")
	}
	s = record.Redact(s, r.home, r.user)
	if r.runRoot != "" {
		s = record.RedactRunRoot(s, r.runRoot)
		// Git prints the symlink-resolved path, which on some systems
		// is not the path the run root was created with.
		if resolved, err := filepath.EvalSymlinks(r.runRoot); err == nil {
			s = record.RedactRunRoot(s, resolved)
		}
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
			if !cleanRel(rel) {
				continue
			}
			src := filepath.Join(repo, rel)
			info, err := os.Lstat(src)
			if err != nil || !info.Mode().IsRegular() || !within(repo, src) {
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
		if !cleanRel(rel) {
			return fmt.Errorf("design artifact %q must be a relative path inside the workspace", rel)
		}
		matches, err := filepath.Glob(filepath.Join(w.WS, filepath.FromSlash(rel)))
		if err != nil {
			return fmt.Errorf("design artifact %q: %w", rel, err)
		}
		if len(matches) == 0 {
			if drops != nil {
				*drops = append(*drops, "design artifact "+rel+": no match")
			}
			continue
		}
		for _, match := range matches {
			if err := rt.copyDesignMatch(w, match, dest, red, drops); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rt *runtime) copyDesignMatch(w world.World, match, dest string, red textRedactor, drops *[]string) error {
	info, err := os.Lstat(match)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		rel, err := filepath.Rel(w.WS, match)
		if err != nil {
			return err
		}
		if !cleanRel(rel) || !within(w.WS, match) {
			if drops != nil {
				*drops = append(*drops, rel+" not a regular file inside the workspace")
			}
			return nil
		}
		return filepath.WalkDir(match, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path == match || d.IsDir() {
				return nil
			}
			return rt.copyDesignFile(w, path, dest, red, drops)
		})
	}
	return rt.copyDesignFile(w, match, dest, red, drops)
}

func (rt *runtime) copyDesignFile(w world.World, src, dest string, red textRedactor, drops *[]string) error {
	rel, err := filepath.Rel(w.WS, src)
	if err != nil {
		return err
	}
	info, lerr := os.Lstat(src)
	if lerr != nil || !cleanRel(rel) || !info.Mode().IsRegular() || !within(w.WS, src) {
		if drops != nil {
			*drops = append(*drops, rel+" not a regular file inside the workspace")
		}
		return nil
	}
	b, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return writeBytes(filepath.Join(dest, rel), b, red, drops)
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
	if record.LooksLikeSecret(string(data)) || containsAny(string(data), red.secrets) {
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

func containsAny(s string, needles []string) bool {
	for _, n := range needles {
		if n != "" && strings.Contains(s, n) {
			return true
		}
	}
	return false
}

// within reports whether path, after resolving symlinks, is root or below it.
func within(root, path string) bool {
	r, err := filepath.EvalSymlinks(root)
	if err != nil {
		return false
	}
	p, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(r, p)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

// cleanRel accepts only a relative path that stays below its base.
func cleanRel(rel string) bool {
	if rel == "" || filepath.IsAbs(rel) {
		return false
	}
	c := filepath.Clean(rel)
	return c != ".." && !strings.HasPrefix(c, ".."+string(filepath.Separator))
}
