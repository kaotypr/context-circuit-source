package report

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

// Options selects an existing run directory. No host is started.
type Options struct {
	OutDir   string
	Run      string
	Baseline string
	Design   bool
	Upgrade  bool
	NoZip    bool
	Home     string
	User     string
	RunRoots []string
	Checkout string
	Path     string
	Stdout   io.Writer
}

// Run writes report.md, and when asked design.md and bundle.zip.
func Run(ctx context.Context, opts Options) error {
	out := opts.Stdout
	if out == nil {
		out = os.Stdout
	}
	dir, err := ResolveRun(opts.OutDir, opts.Run)
	if err != nil {
		return err
	}
	runs, err := LoadRuns(filepath.Join(dir, "runs.jsonl"))
	if err != nil {
		return err
	}
	var before []record.Run
	baselineName := ""
	if opts.Baseline != "" {
		before, err = LoadBaseline(opts.Baseline)
		if err != nil {
			return err
		}
		baselineName = filepath.Base(opts.Baseline)
	}
	branch, baseSide := splitSides(runs, before)
	if opts.Baseline == "" {
		before = baseSide
	}
	surface, hasSurface, err := loadSurface(dir)
	if err != nil {
		return err
	}
	model := Build(Input{
		Branch:       branch,
		Before:       before,
		BaselineName: baselineName,
		Surface:      surface,
		HasSurface:   hasSurface,
	})
	if err := os.WriteFile(filepath.Join(dir, "report.md"), []byte(Render(model)), 0o644); err != nil {
		return err
	}
	fmt.Fprintln(out, "wrote report.md")
	if !opts.NoZip {
		if err := WriteBundle(BundleOptions{
			RunDir:   dir,
			Dest:     filepath.Join(dir, "bundle.zip"),
			Home:     opts.Home,
			User:     opts.User,
			RunRoots: opts.RunRoots,
		}); err != nil {
			return err
		}
		fmt.Fprintln(out, "wrote bundle.zip")
	}
	if !opts.Design && !opts.Upgrade {
		return nil
	}
	body, err := RenderDesign(DesignInput{
		RunDir:      dir,
		BaselineDir: opts.Baseline,
		Branch:      branch,
		Before:      before,
		Home:        opts.Home,
		User:        opts.User,
		RunRoots:    opts.RunRoots,
	})
	if err != nil {
		return err
	}
	if opts.Upgrade {
		baseRef, branchRef := loadRefs(dir, runs)
		section, uerr := UpgradeSection(ctx, UpgradeOptions{
			RunDir:    dir,
			Checkout:  opts.Checkout,
			Path:      opts.Path,
			Home:      opts.Home,
			BaseRef:   baseRef,
			BranchRef: branchRef,
		})
		if !strings.HasSuffix(body, "\n") {
			body += "\n"
		}
		body += "\n" + section
		if werr := os.WriteFile(filepath.Join(dir, "design.md"), []byte(body), 0o644); werr != nil {
			return werr
		}
		fmt.Fprintln(out, "wrote design.md")
		return uerr
	}
	if err := os.WriteFile(filepath.Join(dir, "design.md"), []byte(body), 0o644); err != nil {
		return err
	}
	fmt.Fprintln(out, "wrote design.md")
	return nil
}

// ResolveRun finds the run directory. An empty run id is the newest
// directory under outDir that contains runs.jsonl.
func ResolveRun(outDir, run string) (string, error) {
	if strings.TrimSpace(run) != "" {
		if filepath.IsAbs(run) || strings.Contains(run, string(filepath.Separator)) {
			if isDir(run) {
				return run, nil
			}
		}
		cand := filepath.Join(outDir, run)
		if isDir(cand) {
			return cand, nil
		}
		if isDir(run) {
			return run, nil
		}
		return "", fmt.Errorf("report: run %q was not found", filepath.Base(run))
	}
	if outDir == "" {
		return "", fmt.Errorf("report: output directory is empty")
	}
	return newestRun(outDir)
}

func newestRun(outDir string) (string, error) {
	entries, err := os.ReadDir(outDir)
	if err != nil {
		return "", err
	}
	var best string
	var bestTime time.Time
	var bestName string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(outDir, entry.Name())
		if _, err := os.Stat(filepath.Join(dir, "runs.jsonl")); err != nil {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return "", err
		}
		if best == "" || info.ModTime().After(bestTime) || (info.ModTime().Equal(bestTime) && entry.Name() > bestName) {
			best = dir
			bestTime = info.ModTime()
			bestName = entry.Name()
		}
	}
	if best == "" {
		return "", fmt.Errorf("report: no run with runs.jsonl under the output directory")
	}
	return best, nil
}

func splitSides(runs, baseline []record.Run) (branch, base []record.Run) {
	if len(baseline) > 0 {
		for _, run := range runs {
			if run.Side == "base" {
				continue
			}
			branch = append(branch, run)
		}
		return branch, nil
	}
	for _, run := range runs {
		if run.Side == "base" {
			base = append(base, run)
			continue
		}
		branch = append(branch, run)
	}
	return branch, base
}

func loadRefs(runDir string, runs []record.Run) (baseRef, branchRef string) {
	b, err := os.ReadFile(filepath.Join(runDir, "surface", "surface.json"))
	if err == nil {
		var raw struct {
			Base   string `json:"base"`
			Branch string `json:"branch"`
		}
		if json.Unmarshal(b, &raw) == nil && raw.Base != "" && raw.Branch != "" {
			return raw.Base, raw.Branch
		}
	}
	for _, run := range runs {
		if run.Side == "base" && baseRef == "" {
			baseRef = run.SrcCommit
		}
		if run.Side != "base" && branchRef == "" {
			branchRef = run.SrcCommit
		}
		if baseRef == "" && run.BaseCommit != "" {
			baseRef = run.BaseCommit
		}
	}
	return baseRef, branchRef
}
