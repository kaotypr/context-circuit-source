package report

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

// Surface is the part of surface.json the report prints.
type Surface struct {
	BaseSHA   string
	Files     int
	Outputs   int
	Scenarios []string
	Kinds     []string
}

// LoadRuns reads a runs.jsonl file written by record.Append.
func LoadRuns(path string) ([]record.Run, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(make([]byte, 0, 64*1024), 4<<20)
	var out []record.Run
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var run record.Run
		if err := json.Unmarshal([]byte(line), &run); err != nil {
			return nil, fmt.Errorf("runs.jsonl:%d: %w", lineNo, err)
		}
		out = append(out, run)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// LoadBaseline reads runs.jsonl, or summary.json when the log is absent.
// Every returned run is the before side.
func LoadBaseline(dir string) ([]record.Run, error) {
	log := filepath.Join(dir, "runs.jsonl")
	if st, err := os.Stat(log); err == nil && !st.IsDir() {
		runs, err := LoadRuns(log)
		if err != nil {
			return nil, err
		}
		for i := range runs {
			runs[i].Side = "base"
		}
		return runs, nil
	}
	summary := filepath.Join(dir, "summary.json")
	if st, err := os.Stat(summary); err == nil && !st.IsDir() {
		return loadSummary(summary)
	}
	return nil, fmt.Errorf("baseline has no runs.jsonl or summary.json")
}

type summaryFile struct {
	Rows []summaryRow `json:"rows"`
}

type summaryRow struct {
	Scenario     string   `json:"scenario"`
	Host         string   `json:"host"`
	HostVersion  string   `json:"host_version"`
	Model        string   `json:"model"`
	Effort       string   `json:"effort"`
	Verdict      string   `json:"verdict"`
	InputTotal   *int64   `json:"input_total"`
	Output       *int64   `json:"output"`
	CostUSD      *float64 `json:"cost_usd"`
	CostBasis    string   `json:"cost_basis"`
	Runs         int      `json:"runs"`
	FailedChecks []string `json:"failed_checks"`
}

func loadSummary(path string) ([]record.Run, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file summaryFile
	if err := json.Unmarshal(b, &file); err != nil {
		return nil, fmt.Errorf("summary.json: %w", err)
	}
	if file.Rows == nil {
		return nil, fmt.Errorf("summary.json: missing rows")
	}
	var out []record.Run
	for _, row := range file.Rows {
		n := row.Runs
		if n < 1 {
			n = 1
		}
		for i := 1; i <= n; i++ {
			out = append(out, record.Run{
				Scenario:     row.Scenario,
				Host:         row.Host,
				HostVersion:  row.HostVersion,
				Model:        row.Model,
				Effort:       row.Effort,
				Verdict:      row.Verdict,
				Run:          i,
				Side:         "base",
				Summary:      true,
				Tokens:       record.Tokens{InputTotal: row.InputTotal, Output: row.Output},
				CostUSD:      row.CostUSD,
				CostBasis:    row.CostBasis,
				FailedChecks: append([]string{}, row.FailedChecks...),
			})
		}
	}
	return out, nil
}

func loadSurface(runDir string) (Surface, bool, error) {
	for _, rel := range []string{filepath.Join("surface", "surface.json"), "surface.json"} {
		b, err := os.ReadFile(filepath.Join(runDir, rel))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return Surface{}, false, err
		}
		return parseSurface(b)
	}
	return Surface{}, false, nil
}

func parseSurface(b []byte) (Surface, bool, error) {
	var raw struct {
		Base  string `json:"base"`
		Files []struct {
			Status string `json:"status"`
		} `json:"files"`
		Outputs []struct {
			Changed *bool `json:"changed"`
		} `json:"outputs"`
		Scenarios []string `json:"scenarios"`
		Kinds     []string `json:"kinds"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return Surface{}, false, fmt.Errorf("surface.json: %w", err)
	}
	nOut := 0
	for _, out := range raw.Outputs {
		if out.Changed != nil && *out.Changed {
			nOut++
		}
	}
	return Surface{
		BaseSHA:   raw.Base,
		Files:     len(raw.Files),
		Outputs:   nOut,
		Scenarios: append([]string{}, raw.Scenarios...),
		Kinds:     append([]string{}, raw.Kinds...),
	}, true, nil
}
