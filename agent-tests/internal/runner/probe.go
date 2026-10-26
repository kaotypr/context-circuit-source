package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

// ProbeOptions selects hosts for a version and overhead probe.
// Confirm is the --yes flag. Without it, Probe prints the measured
// overhead and does not start a prompt.
type ProbeOptions struct {
	Checkout    string
	Ref         string
	HostNames   []string
	HostsFile   string
	PricesFile  string
	FixturesDir string
	OutDir      string
	Caller      map[string]string
	ParentTemp  string
	Confirm     bool
	Stdout      io.Writer
}

// overhead is the fixed cost measured with a loaded user config.
// The isolated figure is lower; S0 is what measures it.
var overhead = []struct {
	host string
	text string
}{
	{"claude", "claude: 42500 input tokens, $0.2519 list"},
	{"codex", "codex: 15974 input tokens"},
	{"cursor", "cursor: 21366 input tokens (2 uncached + 21364 cache write), 73-75 output"},
}

// Probe prints each host's version, pinned model, and auth variable names.
// With Confirm it also runs the overhead prompt in an empty directory and
// in the common world. Without Confirm it prints the measured token cost
// and returns.
func Probe(ctx context.Context, opts ProbeOptions) error {
	out := opts.Stdout
	if out == nil {
		out = os.Stdout
	}
	if opts.Ref == "" {
		opts.Ref = "HEAD"
	}
	hostFile, err := hosts.LoadHosts(opts.HostsFile)
	if err != nil {
		return fmt.Errorf("hosts: %w", err)
	}
	selected, err := selectHosts(hostFile, opts.HostNames)
	if err != nil {
		return err
	}
	path := callerPath(opts.Caller)
	for _, cfg := range selected {
		fmt.Fprintf(out, "host: %s\n", cfg.Name)
		fmt.Fprintf(out, "configured version: %s\n", cfg.HostVersion)
		ver, verr := commandVersion(ctx, cfg.Command, path)
		if verr != nil {
			fmt.Fprintf(out, "command version: unavailable (%s)\n", verr)
		} else {
			fmt.Fprintf(out, "command version: %s\n", ver)
		}
		if cfg.Model == "" {
			fmt.Fprintln(out, "model: not pinned")
		} else {
			fmt.Fprintf(out, "model: %s pinned\n", cfg.Model)
		}
		adapter, err := hosts.ForName(cfg.Name)
		if err != nil {
			return err
		}
		ready, detail := adapter.AuthReady(opts.Caller, cfg)
		fmt.Fprintf(out, "auth env: %s\n", strings.Join(cfg.AuthEnv, ", "))
		if ready {
			fmt.Fprintln(out, "auth: ready")
		} else {
			fmt.Fprintln(out, "auth: not ready")
			if detail != "" {
				fmt.Fprintln(out, detail)
			}
		}
	}
	fmt.Fprintln(out, "measured with the maintainer's own config loaded; the isolated figure is lower")
	if !opts.Confirm {
		fmt.Fprintln(out, "would run S0 prompt \"Reply with exactly: OK\" in an empty directory (overhead-empty) and in the common-world workspace (overhead-world), each fully isolated")
		for _, row := range overhead {
			if hostSelected(selected, row.host) {
				fmt.Fprintln(out, row.text)
			}
		}
		return nil
	}
	if opts.Checkout == "" || opts.OutDir == "" || opts.PricesFile == "" || opts.FixturesDir == "" {
		return fmt.Errorf("probe --yes needs checkout, out, prices, and fixtures")
	}
	prices, err := hosts.LoadPrices(opts.PricesFile)
	if err != nil {
		return err
	}
	pin, err := worldPin(opts.Checkout)
	if err != nil {
		return err
	}
	srcCommit, err := revParse(ctx, opts.Checkout, opts.Ref, path)
	if err != nil {
		return err
	}
	digest := ""
	if opts.FixturesDir != "" {
		// Digest is best-effort here; a missing tree is reported by Setup.
		digest, err = digestFixtures(opts.FixturesDir)
		if err != nil {
			return err
		}
	}
	rt := &runtime{
		opts: RunOptions{
			Checkout:    opts.Checkout,
			Ref:         opts.Ref,
			FixturesDir: opts.FixturesDir,
			OutDir:      opts.OutDir,
			Caller:      opts.Caller,
			ParentTemp:  opts.ParentTemp,
			Confirm:     true,
		},
		prices:        prices,
		id:            time.Now().UTC().Format("20060102T150405Z"),
		pin:           pin,
		srcCommit:     srcCommit,
		fixtureDigest: digest,
		stdout:        out,
	}
	var failures []string
	for _, cfg := range selected {
		for _, spec := range []struct {
			id    string
			empty bool
		}{
			{"overhead-empty", true},
			{"overhead-world", false},
		} {
			item := planned{
				Scenario: scenario.Scenario{
					ID:      spec.id,
					Version: 0,
					Title:   "Overhead",
					Prompt:  "Reply with exactly: OK",
				},
				Host:   cfg.Name,
				Config: cfg,
				Side:   "branch",
				Ref:    opts.Ref,
				RunN:   1,
			}
			if err := rt.one(ctx, item, spec.empty); err != nil {
				fmt.Fprintf(out, "harness-error %s %s: %s\n", spec.id, cfg.Name, err)
				failures = append(failures, err.Error())
			}
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("%d probe run(s) recorded harness-error", len(failures))
	}
	return nil
}

func commandVersion(ctx context.Context, command, path string) (string, error) {
	res, err := proc.Run(ctx, proc.Spec{
		Argv:    []string{command, "--version"},
		Env:     []string{"PATH=" + path},
		Timeout: 30 * time.Second,
	})
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("exit %d", res.ExitCode)
	}
	text := strings.TrimSpace(string(res.Stdout))
	if line, _, ok := strings.Cut(text, "\n"); ok {
		return strings.TrimSpace(line), nil
	}
	return text, nil
}

func hostSelected(list []hosts.HostConfig, name string) bool {
	for _, cfg := range list {
		if cfg.Name == name {
			return true
		}
	}
	return false
}

func digestFixtures(dir string) (string, error) {
	return record.Digest(dir)
}

func worldPin(checkout string) (string, error) {
	return world.Pin(checkout)
}
