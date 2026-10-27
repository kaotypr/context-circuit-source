// Package record stores one JSON run line, a directory digest, and redacted text.
package record

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// Tokens counts one run. Nil fields marshal as null: a missing report is
// not a zero.
type Tokens struct {
	InputTotal    *int64 `json:"input_total"`
	InputUncached *int64 `json:"input_uncached"`
	CacheWrite    *int64 `json:"cache_write"`
	CacheRead     *int64 `json:"cache_read"`
	Output        *int64 `json:"output"`
	Reasoning     *int64 `json:"reasoning"`
}

// Run is one JSON line in the run log.
// CostBasis is "reported", "price-table", or "none".
// Side is "base" or "branch".
type Run struct {
	Scenario        string `json:"scenario"`
	ScenarioVersion string `json:"scenario_version"`
	FixtureDigest   string `json:"fixture_digest"`
	SrcCommit       string `json:"src_commit"`
	BaseCommit      string `json:"base_commit"`
	Side            string `json:"side"`
	CLIVersion      string `json:"cli_version"`
	Host            string `json:"host"`
	HostVersion     string `json:"host_version"`
	// HostVersionSource is "command" when HostVersion is what the host
	// binary printed for --version, and "config" when it is the value from
	// hosts.yaml because the command could not be read.
	HostVersionSource string `json:"host_version_source,omitempty"`
	// Model is the pinned model id from hosts.yaml, the one sent on the
	// command line, so a baseline row matches its hosts.yaml entry.
	Model string `json:"model"`
	// HostModel is the model name the host itself reported, kept only when
	// it differs from Model: Cursor reports a display name such as
	// "Composer 2.5" for the pin composer-2.5.
	HostModel       string            `json:"host_model,omitempty"`
	Effort          string            `json:"effort"`
	Run             int               `json:"run"`
	Verdict         string            `json:"verdict"`
	FailedChecks    []string          `json:"failed_checks"`
	Soft            map[string]any    `json:"soft"`
	Tokens          Tokens            `json:"tokens"`
	UsageSource     string            `json:"usage_source"`
	CostUSD         *float64          `json:"cost_usd"`
	CostBasis       string            `json:"cost_basis"`
	NullReasons     map[string]string `json:"null_reasons,omitempty"`
	Turns           int               `json:"turns"`
	ToolCalls       int               `json:"tool_calls"`
	Subagents       int               `json:"subagents"`
	WallS           float64           `json:"wall_s"`
	ArtifactsDigest string            `json:"artifacts_digest"`
	Date            string            `json:"date"`
	Summary         bool              `json:"-"`
}

// Append writes r as one JSON line, creating parent directories.
// The write uses O_APPEND.
func Append(path string, r Run) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = f.Write(b)
	return err
}

// Digest is a sha256 over the sorted relative paths and file contents under
// root. The .git directory is skipped. The result is "sha256:" plus hex.
func Digest(root string) (string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("digest root is not a directory")
	}
	var files []string
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(files)
	h := sha256.New()
	for _, rel := range files {
		if _, err := io.WriteString(h, filepath.ToSlash(rel)); err != nil {
			return "", err
		}
		if _, err := h.Write([]byte{0}); err != nil {
			return "", err
		}
		b, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			return "", err
		}
		if _, err := h.Write(b); err != nil {
			return "", err
		}
		if _, err := h.Write([]byte{0}); err != nil {
			return "", err
		}
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}
