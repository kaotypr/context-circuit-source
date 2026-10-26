// Package runner drives one scenario on one host inside an isolated world.
package runner

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

// RunOptions selects what a run will execute.
// Confirm is the --yes flag: without it, and without DryRun, Run refuses
// before any host starts.
type RunOptions struct {
	Checkout     string
	Ref          string
	BaseRef      string
	ScenarioIDs  []string
	HostNames    []string
	Runs         int
	HostsFile    string
	PricesFile   string
	ScenariosDir string
	FixturesDir  string
	OutDir       string
	Keep         bool
	Caller       map[string]string
	ParentTemp   string
	DryRun       bool
	Confirm      bool
	Stdout       io.Writer
}

type planned struct {
	Scenario scenario.Scenario
	Host     string
	Config   hosts.HostConfig
	Side     string
	Ref      string
	RunN     int
}

// Run validates the selection, prints the plan, and unless DryRun or
// Confirm says otherwise refuses to start a host. Combinations run one
// at a time. A failure in one is recorded as harness-error and the rest
// continue.
func Run(ctx context.Context, opts RunOptions) error {
	out := opts.Stdout
	if out == nil {
		out = os.Stdout
	}
	if opts.Runs < 1 {
		return fmt.Errorf("runs must be at least 1")
	}
	if opts.Ref == "" {
		opts.Ref = "HEAD"
	}
	hostFile, err := hosts.LoadHosts(opts.HostsFile)
	if err != nil {
		return fmt.Errorf("hosts: %w", err)
	}
	selectedHosts, err := selectHosts(hostFile, opts.HostNames)
	if err != nil {
		fmt.Fprintln(out, err.Error())
		return err
	}
	loaded, err := scenario.LoadDir(opts.ScenariosDir)
	if err != nil {
		fmt.Fprintln(out, err.Error())
		return err
	}
	selected, err := selectScenarios(loaded, opts.ScenarioIDs)
	if err != nil {
		fmt.Fprintln(out, err.Error())
		return err
	}
	plan := planCombos(selected, selectedHosts, opts)
	fmt.Fprint(out, formatPlan(plan))
	if opts.DryRun {
		return nil
	}
	if !opts.Confirm {
		msg := "this run would use model usage on your account; re-run with --yes to start"
		fmt.Fprintln(out, msg)
		return fmt.Errorf("%s", msg)
	}
	if opts.Checkout == "" {
		return fmt.Errorf("checkout is empty")
	}
	if opts.OutDir == "" {
		return fmt.Errorf("output directory is empty")
	}
	prices, err := hosts.LoadPrices(opts.PricesFile)
	if err != nil {
		return fmt.Errorf("prices: %w", err)
	}
	pin, err := world.Pin(opts.Checkout)
	if err != nil {
		return err
	}
	srcCommit, err := revParse(ctx, opts.Checkout, opts.Ref, callerPath(opts.Caller))
	if err != nil {
		return fmt.Errorf("src commit: %w", err)
	}
	baseCommit := ""
	if opts.BaseRef != "" {
		baseCommit, err = revParse(ctx, opts.Checkout, opts.BaseRef, callerPath(opts.Caller))
		if err != nil {
			return fmt.Errorf("base commit: %w", err)
		}
	}
	digest := ""
	if opts.FixturesDir != "" {
		digest, err = record.Digest(opts.FixturesDir)
		if err != nil {
			return fmt.Errorf("fixture digest: %w", err)
		}
	}
	rt := &runtime{
		opts:          opts,
		prices:        prices,
		id:            time.Now().UTC().Format("20060102T150405Z"),
		pin:           pin,
		srcCommit:     srcCommit,
		baseCommit:    baseCommit,
		fixtureDigest: digest,
		stdout:        out,
	}
	var failures []string
	for _, item := range plan {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := rt.one(ctx, item, false); err != nil {
			fmt.Fprintf(out, "harness-error %s %s %s r%d: %s\n", item.Scenario.ID, item.Host, item.Side, item.RunN, err)
			failures = append(failures, err.Error())
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("%d combination(s) recorded harness-error", len(failures))
	}
	return nil
}

type runtime struct {
	opts          RunOptions
	prices        map[string]hosts.PriceRow
	id            string
	pin           string
	srcCommit     string
	baseCommit    string
	fixtureDigest string
	stdout        io.Writer
}

func selectHosts(all map[string]hosts.HostConfig, names []string) ([]hosts.HostConfig, error) {
	if len(names) == 0 {
		names = make([]string, 0, len(all))
		for name := range all {
			names = append(names, name)
		}
		sort.Strings(names)
	}
	var out []hosts.HostConfig
	var problems []string
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			continue
		}
		seen[name] = true
		cfg, ok := all[name]
		if !ok {
			problems = append(problems, fmt.Sprintf("host %q is not in hosts.yaml", name))
			continue
		}
		if cfg.Model == "" {
			problems = append(problems, fmt.Sprintf("host %q has no pinned model", name))
			continue
		}
		out = append(out, cfg)
	}
	if len(problems) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(problems, "\n"))
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no hosts found")
	}
	return out, nil
}

func selectScenarios(all []scenario.Scenario, ids []string) ([]scenario.Scenario, error) {
	if len(ids) == 0 {
		return all, nil
	}
	byID := map[string]scenario.Scenario{}
	for _, s := range all {
		byID[s.ID] = s
	}
	var out []scenario.Scenario
	var problems []string
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		s, ok := byID[id]
		if !ok {
			problems = append(problems, fmt.Sprintf("scenario %q was not found", id))
			continue
		}
		out = append(out, s)
	}
	if len(problems) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(problems, "\n"))
	}
	return out, nil
}

func planCombos(scenarios []scenario.Scenario, hostList []hosts.HostConfig, opts RunOptions) []planned {
	var plan []planned
	for _, sc := range scenarios {
		for _, cfg := range hostList {
			for runN := 1; runN <= opts.Runs; runN++ {
				for _, side := range sidesFor(runN, opts.BaseRef != "") {
					ref := opts.Ref
					if side == "base" {
						ref = opts.BaseRef
					}
					plan = append(plan, planned{
						Scenario: sc,
						Host:     cfg.Name,
						Config:   cfg,
						Side:     side,
						Ref:      ref,
						RunN:     runN,
					})
				}
			}
		}
	}
	return plan
}

// sidesFor interleaves base and branch. Odd run indexes start with base
// and even ones start with branch, so the same side is not always second.
func sidesFor(runN int, withBase bool) []string {
	if !withBase {
		return []string{"branch"}
	}
	if runN%2 == 1 {
		return []string{"base", "branch"}
	}
	return []string{"branch", "base"}
}

func formatPlan(plan []planned) string {
	var b strings.Builder
	fmt.Fprintf(&b, "plan: %d combination(s)\n", len(plan))
	for _, item := range plan {
		fmt.Fprintf(&b, "%s x %s x %s x r%d\n", item.Scenario.ID, item.Host, item.Side, item.RunN)
	}
	return b.String()
}

func callerPath(caller map[string]string) string {
	if caller != nil && caller["PATH"] != "" {
		return caller["PATH"]
	}
	return os.Getenv("PATH")
}

func revParse(ctx context.Context, checkout, ref, path string) (string, error) {
	res, err := proc.Run(ctx, proc.Spec{
		Argv: []string{"git", "-C", checkout, "rev-parse", ref},
		Env: []string{
			"PATH=" + path,
			"GIT_CONFIG_NOSYSTEM=1",
			"GIT_CONFIG_GLOBAL=/dev/null",
		},
	})
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("git rev-parse %s: exit %d\n%s", ref, res.ExitCode, res.Stderr)
	}
	return strings.TrimSpace(string(res.Stdout)), nil
}

func newSessionID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "00000000-0000-4000-8000-000000000000"
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func mergeMissing(base []string, extra map[string]string) []string {
	have := map[string]bool{}
	for _, item := range base {
		key, _, _ := strings.Cut(item, "=")
		have[key] = true
	}
	out := append([]string{}, base...)
	var keys []string
	for key := range extra {
		if !have[key] {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		out = append(out, key+"="+extra[key])
	}
	sort.Strings(out)
	return out
}

func atoiVersion(v int) string {
	return strconv.Itoa(v)
}
