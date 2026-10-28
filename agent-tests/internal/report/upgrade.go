package report

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/surface"
)

// Change is one file in a template delta.
type Change struct {
	Path string
	Op   string
	Body []byte
}

// Finding is one check issue.
type Finding struct {
	Issue   string
	Resolve string
}

// CheckResult is the branch CLI's check --json.
type CheckResult struct {
	OK       bool
	Findings []Finding
}

// UpgradeOptions selects the trees and the branch CLI.
type UpgradeOptions struct {
	RunDir    string
	Checkout  string
	Path      string
	Home      string
	BaseRef   string
	BranchRef string
}

// FileDelta is the file-level difference from base to branch.
// Paths a workspace never receives, the id ledger file, and a version
// change in workspace.yaml are omitted.
func FileDelta(base, branch string) ([]Change, error) {
	baseFiles, err := walkFiles(base)
	if err != nil {
		return nil, err
	}
	branchFiles, err := walkFiles(branch)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var changes []Change
	for rel, body := range branchFiles {
		seen[rel] = true
		if skipRel(rel) {
			continue
		}
		if rel == "workspace.yaml" {
			body = keepVersion(baseFiles[rel], body)
		}
		old, ok := baseFiles[rel]
		if !ok {
			changes = append(changes, Change{Path: rel, Op: "add", Body: body})
			continue
		}
		if !bytes.Equal(old, body) {
			changes = append(changes, Change{Path: rel, Op: "change", Body: append([]byte{}, body...)})
		}
	}
	for rel := range baseFiles {
		if seen[rel] || skipRel(rel) {
			continue
		}
		changes = append(changes, Change{Path: rel, Op: "remove"})
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	return changes, nil
}

func skipRel(rel string) bool {
	switch rel {
	case "README.md", "LICENSE", "CHANGELOG.md", "CODE_OF_CONDUCT.md", "CONTRIBUTING.md", "SECURITY.md", ".context-circuit/ids.yaml":
		return true
	}
	return rel == "assets" || strings.HasPrefix(rel, "assets/")
}

func keepVersion(base, branch []byte) []byte {
	if len(base) == 0 {
		return branch
	}
	baseLine, ok := versionLine(base)
	if !ok {
		return branch
	}
	text := strings.ReplaceAll(string(branch), "\r\n", "\n")
	trailing := strings.HasSuffix(text, "\n")
	lines := strings.Split(text, "\n")
	if trailing && len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	replaced := false
	out := make([]string, 0, len(lines)+1)
	for _, line := range lines {
		if isVersionLine(line) {
			if !replaced {
				out = append(out, baseLine)
				replaced = true
			}
			continue
		}
		out = append(out, line)
	}
	if !replaced {
		out = append(out, baseLine)
	}
	joined := strings.Join(out, "\n")
	if trailing || strings.HasSuffix(string(base), "\n") {
		joined += "\n"
	}
	return []byte(joined)
}

func versionLine(b []byte) (string, bool) {
	text := strings.ReplaceAll(string(b), "\r\n", "\n")
	for _, line := range strings.Split(text, "\n") {
		if isVersionLine(line) {
			return strings.TrimRight(line, "\r"), true
		}
	}
	return "", false
}

func isVersionLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "version:")
}

func walkFiles(root string) (map[string][]byte, error) {
	out := map[string][]byte{}
	if !isDir(root) {
		return nil, fmt.Errorf("upgrade: %s is not a directory", filepath.Base(root))
	}
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 || !d.Type().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !safeRel(rel) {
			return nil
		}
		body, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out[rel] = body
		return nil
	})
	return out, err
}

// ApplyDelta writes changes onto a copy of a workspace.
func ApplyDelta(dst string, changes []Change) error {
	for _, change := range changes {
		if !safeRel(change.Path) {
			return fmt.Errorf("upgrade: unsafe path %q", change.Path)
		}
		target := filepath.Join(dst, filepath.FromSlash(change.Path))
		if change.Op == "remove" {
			if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, change.Body, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// CopyTree copies a workspace, including .git, and skips symlinks.
func CopyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.Type()&os.ModeSymlink != 0 {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		body, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, body, 0o644)
	})
}

// CheckWorkspace runs the branch CLI's check --json.
func CheckWorkspace(ctx context.Context, cli, ws, pathEnv string) (CheckResult, error) {
	if pathEnv == "" {
		pathEnv = os.Getenv("PATH")
	}
	parent := filepath.Dir(ws)
	home := filepath.Join(parent, "check-home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		return CheckResult{}, err
	}
	gitconfig := filepath.Join(parent, "check-gitconfig")
	if err := os.WriteFile(gitconfig, []byte("[user]\n\tname = Agent Test\n\temail = agent-test@example.invalid\n"), 0o644); err != nil {
		return CheckResult{}, err
	}
	res, err := proc.Run(ctx, proc.Spec{
		Argv:    []string{cli, "--workspace", ws, "--json", "check"},
		Dir:     ws,
		Timeout: 2 * time.Minute,
		Env: []string{
			"PATH=" + pathEnv,
			"HOME=" + home,
			"GIT_CONFIG_GLOBAL=" + gitconfig,
			"GIT_CONFIG_NOSYSTEM=1",
			"GIT_TERMINAL_PROMPT=0",
			"LANG=C.UTF-8",
			"LC_ALL=C.UTF-8",
			"NO_COLOR=1",
		},
	})
	if err != nil {
		return CheckResult{}, err
	}
	if res.ExitCode != 0 && res.ExitCode != 1 {
		return CheckResult{}, fmt.Errorf("check exited %d", res.ExitCode)
	}
	return parseCheck(res.Stdout)
}

func parseCheck(stdout []byte) (CheckResult, error) {
	text := strings.TrimSpace(string(stdout))
	if i := strings.Index(text, "{"); i > 0 {
		text = text[i:]
	}
	var raw struct {
		OK     bool `json:"ok"`
		Issues []struct {
			Issue   string `json:"issue"`
			Resolve string `json:"resolve"`
		} `json:"issues"`
	}
	if err := json.Unmarshal([]byte(text), &raw); err != nil {
		return CheckResult{}, fmt.Errorf("check json: %w", err)
	}
	out := CheckResult{OK: raw.OK}
	for _, issue := range raw.Issues {
		out.Findings = append(out.Findings, Finding{Issue: issue.Issue, Resolve: issue.Resolve})
	}
	return out, nil
}

// UpgradeSection applies the template delta and returns the design.md section.
func UpgradeSection(ctx context.Context, opts UpgradeOptions) (string, error) {
	baseTree, branchTree, cli, cleanup, err := resolveUpgrade(ctx, opts)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return renderUpgrade(nil, nil, err), err
	}
	changes, err := FileDelta(baseTree, branchTree)
	if err != nil {
		return renderUpgrade(nil, nil, err), err
	}
	parent, err := os.MkdirTemp("", "agent-test-upgrade-")
	if err != nil {
		return renderUpgrade(changes, nil, err), err
	}
	defer os.RemoveAll(parent)
	dest := filepath.Join(parent, "ws")
	if err := CopyTree(baseTree, dest); err != nil {
		return renderUpgrade(changes, nil, err), err
	}
	if err := ApplyDelta(dest, changes); err != nil {
		return renderUpgrade(changes, nil, err), err
	}
	result, err := CheckWorkspace(ctx, cli, dest, opts.Path)
	if err != nil {
		return renderUpgrade(changes, nil, err), err
	}
	return renderUpgrade(changes, &result, nil), nil
}

func resolveUpgrade(ctx context.Context, opts UpgradeOptions) (baseTree, branchTree, cli string, cleanup func(), err error) {
	base := filepath.Join(opts.RunDir, "surface", "base")
	branch := filepath.Join(opts.RunDir, "surface", "branch")
	if isDir(base) && isDir(branch) {
		ref := opts.BranchRef
		if ref == "" {
			ref = "HEAD"
		}
		if opts.Checkout == "" {
			return "", "", "", nil, fmt.Errorf("upgrade check needs the source checkout to build the branch CLI")
		}
		_, cli, cleanup, err = surface.OpenInitWorkspace(ctx, opts.Checkout, ref, opts.Path, opts.Home)
		if err != nil {
			return "", "", "", cleanup, err
		}
		return base, branch, cli, cleanup, nil
	}
	if opts.Checkout == "" || opts.BaseRef == "" || opts.BranchRef == "" {
		return "", "", "", nil, fmt.Errorf("upgrade check needs surface/base and surface/branch, or a checkout and both refs")
	}
	baseWS, _, cleanBase, err := surface.OpenInitWorkspace(ctx, opts.Checkout, opts.BaseRef, opts.Path, opts.Home)
	if err != nil {
		return "", "", "", nil, err
	}
	branchWS, cli, cleanBranch, err := surface.OpenInitWorkspace(ctx, opts.Checkout, opts.BranchRef, opts.Path, opts.Home)
	if err != nil {
		cleanBase()
		return "", "", "", nil, err
	}
	cleanup = func() {
		cleanBase()
		cleanBranch()
	}
	return baseWS, branchWS, cli, cleanup, nil
}

func renderUpgrade(changes []Change, result *CheckResult, runErr error) string {
	var b strings.Builder
	b.WriteString("## Upgrade check\n\n")
	b.WriteString("This is a file-level delta applied to a copy of the base workspace, not the host's judgement.\n\n")
	b.WriteString("Not applied: `README.md`, `LICENSE`, `CHANGELOG.md`, `CODE_OF_CONDUCT.md`, `CONTRIBUTING.md`, `SECURITY.md`, `assets/`, `.context-circuit/ids.yaml`, and a `version` change in `workspace.yaml`.\n\n")
	b.WriteString(appliedLine(changes))
	b.WriteString("\n### Findings\n\n")
	if runErr != nil {
		fmt.Fprintf(&b, "check could not be run: %s\n", runErr)
		return b.String()
	}
	if result == nil || len(result.Findings) == 0 {
		b.WriteString("check reported no findings.\n")
		return b.String()
	}
	for _, finding := range result.Findings {
		fmt.Fprintf(&b, "- %s\n", finding.Issue)
		if finding.Resolve != "" {
			fmt.Fprintf(&b, "  - %s\n", finding.Resolve)
		}
	}
	return b.String()
}

func appliedLine(changes []Change) string {
	if len(changes) == 0 {
		return "Applied no files.\n"
	}
	if len(changes) > 20 {
		return fmt.Sprintf("Applied %d files.\n", len(changes))
	}
	parts := make([]string, len(changes))
	for i, change := range changes {
		parts[i] = "`" + change.Path + "` (" + change.Op + ")"
	}
	return "Applied " + strings.Join(parts, ", ") + ".\n"
}
