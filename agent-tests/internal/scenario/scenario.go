// Package scenario loads a scenario file and defines the grader interface.
package scenario

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

// idPattern is the scenario id shape: lowercase words separated by hyphens.
var idPattern = regexp.MustCompile(`^[a-z0-9-]+$`)

// Scenario is one host prompt and the checks that grade it.
// Version is the human-facing revision of the text. Hash is the identity
// of a run: it covers the canonical scenario bytes and the fixture digest,
// so two files with the same version can still differ.
type Scenario struct {
	ID              string            `yaml:"id"`
	Version         int               `yaml:"version"`
	Title           string            `yaml:"title"`
	Covers          []string          `yaml:"covers"`
	Setup           []string          `yaml:"setup"`
	Prompt          string            `yaml:"prompt"`
	ExpectedIssues  []string          `yaml:"expected_issues"`
	Checks          []string          `yaml:"checks"`
	Soft            map[string]string `yaml:"soft"`
	DesignArtifacts []string          `yaml:"design_artifacts"`
	// Workspace is the single path element under the run root where the
	// host starts. Empty means ws.
	Workspace string `yaml:"workspace"`
	// Human lists judgements a person still has to make.
	Human []string `yaml:"human"`
	// AllowedBranches are path.Match patterns for branch names this
	// scenario may create or advance. The name is the part after
	// refs/heads/.
	AllowedBranches []string `yaml:"allowed_branches"`
}

// Snapshot is the workspace state a check can read. Maps are empty rather
// than nil when a snapshot was taken and that section had nothing in it.
// Status, Heads, Refs, and Worktrees are keyed by "ws" for the workspace
// and by directory name for each Git checkout under repositories/.
type Snapshot struct {
	DecoyLog   string
	Status     map[string]string
	Heads      map[string]string
	RemoteRefs map[string]string
	Files      []string
	// Refs holds sorted "refname objectname" lines of refs/heads for one repo.
	Refs map[string]string
	// Worktrees holds `git worktree list --porcelain` text for one repo.
	Worktrees map[string]string
	// RunDigest maps a relative path under the run root to the sha256 hex
	// of a regular file, "dir" for a directory, or "link:<target>" for a
	// symlink. Paths the snapshot skips are absent.
	RunDigest map[string]string
}

// CheckContext is the evidence one check sees.
type CheckContext struct {
	World    world.World
	Scenario Scenario
	Parsed   hosts.Parsed
	Before   Snapshot
	After    Snapshot
	// CheckBefore and CheckAfter are the raw `check --json` documents
	// from after setup and after the host exits.
	CheckBefore []byte
	CheckAfter  []byte
	// Stream is the raw host stream text.
	Stream string
	// Exec runs a command with the isolated environment used for
	// snapshots and returns stdout. A non-nil error includes the exit
	// code and stderr. It is nil in tests that do not need a command;
	// checks treat nil as inconclusive.
	Exec func(dir string, argv ...string) (string, error)
}

// CheckResult is one required check.
//
// Inconclusive is set only together with Pass. The result is not a
// failure, it is never counted as a pass, and it makes the run
// inconclusive. Verdict is unchanged: a failed check still beats an
// inconclusive run.
type CheckResult struct {
	Name         string
	Pass         bool
	Inconclusive bool
	Detail       string
}

// Check grades one named property of a run.
type Check func(ctx CheckContext) CheckResult

// Registry holds the shared checks and the scenario checks. Each check
// registers itself.
var Registry = map[string]Check{}

// SoftRegistry holds optional metrics. A registered function's return
// value is what the run record stores. An unregistered name keeps the
// description from the scenario file.
var SoftRegistry = map[string]func(CheckContext) any{}

// Load reads one scenario file. Unknown fields are an error.
func Load(path string) (Scenario, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Scenario{}, err
	}
	var s Scenario
	if err := yaml.UnmarshalWithOptions(b, &s, yaml.Strict()); err != nil {
		return Scenario{}, fmt.Errorf("%s: %w", path, err)
	}
	if err := s.validate(); err != nil {
		return Scenario{}, fmt.Errorf("%s: %w", path, err)
	}
	return s, nil
}

// LoadDir reads every YAML scenario in dir, one file each, not nested.
// A missing directory, or a directory with no scenario files, is an error
// whose text is "no scenarios found". Duplicate ids are an error.
func LoadDir(dir string) ([]Scenario, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no scenarios found")
		}
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			files = append(files, filepath.Join(dir, name))
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no scenarios found")
	}
	sort.Strings(files)
	seen := map[string]string{}
	var out []Scenario
	for _, file := range files {
		s, err := Load(file)
		if err != nil {
			return nil, err
		}
		if prev, ok := seen[s.ID]; ok {
			return nil, fmt.Errorf("duplicate scenario id %q (%s and %s)", s.ID, prev, file)
		}
		seen[s.ID] = file
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (s Scenario) validate() error {
	if !idPattern.MatchString(s.ID) {
		return fmt.Errorf("scenario id %q must match %s", s.ID, idPattern.String())
	}
	if strings.TrimSpace(s.Prompt) == "" {
		return fmt.Errorf("scenario %s has an empty prompt", s.ID)
	}
	if err := validateWorkspace(s.Workspace); err != nil {
		return fmt.Errorf("scenario %s: %w", s.ID, err)
	}
	if err := validateAllowedBranches(s.AllowedBranches); err != nil {
		return fmt.Errorf("scenario %s: %w", s.ID, err)
	}
	return nil
}

// validateWorkspace accepts an empty name, which means ws, or one path
// element with no separator and no ".." .
func validateWorkspace(name string) error {
	if name == "" {
		return nil
	}
	if strings.ContainsAny(name, `/\`) || strings.Contains(name, "..") || name == "." || name != path.Base(name) {
		return fmt.Errorf("workspace %q must be a single path element", name)
	}
	return nil
}

func validateAllowedBranches(patterns []string) error {
	for _, pattern := range patterns {
		if strings.TrimSpace(pattern) == "" {
			return fmt.Errorf("allowed branch pattern is empty")
		}
		if _, err := path.Match(pattern, "branch"); err != nil {
			return fmt.Errorf("allowed branch pattern %q: %w", pattern, err)
		}
	}
	return nil
}

// Hash is the identity of this scenario together with a fixture digest.
// The version field is for people reading the file; callers that compare
// two runs use this hash. The digest is sha256 over the canonical YAML
// of the scenario, a NUL byte, and the fixture digest string.
func (s Scenario) Hash(fixtureDigest string) string {
	b, err := yaml.Marshal(s)
	if err != nil {
		b = []byte(s.ID)
	}
	sum := sha256.New()
	sum.Write(b)
	sum.Write([]byte{0})
	sum.Write([]byte(fixtureDigest))
	return "sha256:" + hex.EncodeToString(sum.Sum(nil))
}

// Verdict returns one of pass, fail, harness-error, or inconclusive.
//
// harness-error means setup, auth, or the host crashed. It is excluded
// from a score and rerun. It wins over every other verdict.
// fail means a required check is contradicted by the evidence.
// inconclusive means the host did not expose evidence a check needs,
// and is never counted as a pass. A check reports that by setting
// Inconclusive with Pass; the runner passes it in as hostIncomplete.
// pass means every required check passed.
func Verdict(results []CheckResult, harnessError error, hostIncomplete bool) string {
	if harnessError != nil {
		return "harness-error"
	}
	for _, result := range results {
		if !result.Pass {
			return "fail"
		}
	}
	if hostIncomplete {
		return "inconclusive"
	}
	return "pass"
}
