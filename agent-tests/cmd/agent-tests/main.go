// Command agent-tests runs contributor scenarios against coding hosts.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/report"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/runner"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/surface"
)

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(2)
	}
	// The process environment is read once, here, and passed down as a map.
	caller := callerEnv()
	switch os.Args[1] {
	case "run":
		os.Exit(runCmd(os.Args[2:], caller))
	case "probe":
		os.Exit(probeCmd(os.Args[2:], caller))
	case "surface-diff":
		os.Exit(surfaceDiffCmd(os.Args[2:], caller))
	case "report":
		os.Exit(reportCmd(os.Args[2:], caller))
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n", os.Args[1])
		usage(os.Stderr)
		os.Exit(2)
	}
}

func usage(w io.Writer) {
	fmt.Fprint(w, `agent-tests run [flags]
agent-tests probe [flags]
agent-tests surface-diff --base <ref> [flags]
agent-tests report [flags]

run    execute scenarios. Without --yes the plan is printed and no host starts.
probe  print host version, pinned model, and auth variable names.
surface-diff
       compare two refs offline. It does not start a host and needs no key.
       Exit 0 whether or not the surface differs.
report read a finished run and write report.md. It does not start a host.
       --design writes design.md. --upgrade-check runs check on a file-level delta.

run flags:
  --scenarios   comma-separated scenario ids (default: all in scenarios/)
  --hosts       comma-separated host names (default: all in hosts.yaml)
  --runs        repetitions of each scenario and host (default: 1)
  --base        also run this ref, alternating with --ref
  --ref         git ref to archive (default: HEAD)
  --checkout    source checkout (default: git top-level of the current directory)
  --out         output directory (default: agent-tests/out under the checkout)
  --keep        keep the run root after the run
  --yes         start hosts; a live run uses model usage on your account
  --dry-run     print the plan and exit

probe flags:
  --hosts --ref --checkout --out --yes
  Without --yes, probe prints the measured overhead and does not start a prompt.

surface-diff flags:
  --base        git ref of the earlier side (required)
  --ref         git ref of the later side (default: HEAD)
  --checkout    source checkout (default: git top-level of the current directory)
  --out         output directory (default: agent-tests/out under the checkout)

report flags:
  --run            run id under --out, or a path (default: newest run)
  --out            output directory (default: agent-tests/out under the checkout)
  --baseline       directory with runs.jsonl or summary.json, used as before
  --design         write design.md, pairing base and branch runs
  --upgrade-check  apply the template delta and run the branch CLI check
  --no-zip         do not write bundle.zip
`)
}

func runCmd(args []string, caller map[string]string) int {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	scenarios := fs.String("scenarios", "", "comma-separated scenario ids")
	hosts := fs.String("hosts", "", "comma-separated host names")
	runs := fs.Int("runs", 1, "repetitions")
	base := fs.String("base", "", "also run this ref")
	ref := fs.String("ref", "HEAD", "git ref to archive")
	checkout := fs.String("checkout", "", "source checkout")
	out := fs.String("out", "", "output directory")
	keep := fs.Bool("keep", false, "keep the run root")
	yes := fs.Bool("yes", false, "start hosts")
	dry := fs.Bool("dry-run", false, "print the plan and exit")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: agent-tests run [flags]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return 2
	}
	root, err := resolveCheckout(*checkout, caller["PATH"])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	hostsFile, pricesFile, scenariosDir, fixturesDir, outDir := layout(root)
	if *out != "" {
		outDir = *out
	}
	err = runner.Run(context.Background(), runner.RunOptions{
		Checkout:     root,
		Ref:          *ref,
		BaseRef:      *base,
		ScenarioIDs:  splitList(*scenarios),
		HostNames:    splitList(*hosts),
		Runs:         *runs,
		HostsFile:    hostsFile,
		PricesFile:   pricesFile,
		ScenariosDir: scenariosDir,
		FixturesDir:  fixturesDir,
		OutDir:       outDir,
		Keep:         *keep,
		Caller:       caller,
		DryRun:       *dry,
		Confirm:      *yes,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func probeCmd(args []string, caller map[string]string) int {
	fs := flag.NewFlagSet("probe", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	hosts := fs.String("hosts", "", "comma-separated host names")
	ref := fs.String("ref", "HEAD", "git ref to archive")
	checkout := fs.String("checkout", "", "source checkout")
	out := fs.String("out", "", "output directory")
	yes := fs.Bool("yes", false, "run the overhead prompt")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: agent-tests probe [flags]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return 2
	}
	root, err := resolveCheckout(*checkout, caller["PATH"])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	hostsFile, pricesFile, _, fixturesDir, outDir := layout(root)
	if *out != "" {
		outDir = *out
	}
	err = runner.Probe(context.Background(), runner.ProbeOptions{
		Checkout:    root,
		Ref:         *ref,
		HostNames:   splitList(*hosts),
		HostsFile:   hostsFile,
		PricesFile:  pricesFile,
		FixturesDir: fixturesDir,
		OutDir:      outDir,
		Caller:      caller,
		Confirm:     *yes,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func surfaceDiffCmd(args []string, caller map[string]string) int {
	fs := flag.NewFlagSet("surface-diff", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	base := fs.String("base", "", "git ref of the earlier side")
	ref := fs.String("ref", "HEAD", "git ref of the later side")
	checkout := fs.String("checkout", "", "source checkout")
	outDirFlag := fs.String("out", "", "output directory")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: agent-tests surface-diff --base <ref> [--ref HEAD] [--checkout PATH] [--out DIR]")
		fmt.Fprintln(os.Stderr, "compare two refs offline. It does not start a host and needs no key.")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return 2
	}
	if strings.TrimSpace(*base) == "" {
		fmt.Fprintln(os.Stderr, "surface-diff: --base is required")
		fs.Usage()
		return 2
	}
	root, err := resolveCheckout(*checkout, caller["PATH"])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	hostsFile, _, scenariosDir, _, outDir := layout(root)
	if *outDirFlag != "" {
		outDir = *outDirFlag
	}
	_, err = surface.Run(context.Background(), surface.Options{
		Checkout:     root,
		BaseRef:      *base,
		Ref:          *ref,
		OutDir:       outDir,
		TriggersFile: filepath.Join(root, "agent-tests", "triggers.yaml"),
		HostsFile:    hostsFile,
		ScenariosDir: scenariosDir,
		Path:         caller["PATH"],
		Home:         caller["HOME"],
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func reportCmd(args []string, caller map[string]string) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	runID := fs.String("run", "", "run id or path (default: newest under --out)")
	outDirFlag := fs.String("out", "", "output directory")
	baseline := fs.String("baseline", "", "before runs: a directory with runs.jsonl or summary.json")
	design := fs.Bool("design", false, "write design.md")
	upgrade := fs.Bool("upgrade-check", false, "apply the template delta and run check")
	noZip := fs.Bool("no-zip", false, "do not write bundle.zip")
	checkout := fs.String("checkout", "", "source checkout, used by --upgrade-check")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: agent-tests report [--run ID] [--out DIR] [--baseline DIR] [--design] [--upgrade-check] [--no-zip]")
		fmt.Fprintln(os.Stderr, "read a finished run. It does not start a host.")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return 2
	}
	outDir := *outDirFlag
	var root string
	if outDir == "" || *upgrade || *checkout != "" {
		resolved, err := resolveCheckout(*checkout, caller["PATH"])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		root = resolved
		if outDir == "" {
			_, _, _, _, outDir = layout(root)
		}
	}
	user := caller["USER"]
	if user == "" {
		user = caller["LOGNAME"]
	}
	err := report.Run(context.Background(), report.Options{
		OutDir:   outDir,
		Run:      *runID,
		Baseline: *baseline,
		Design:   *design,
		Upgrade:  *upgrade,
		NoZip:    *noZip,
		Home:     caller["HOME"],
		User:     user,
		Checkout: root,
		Path:     caller["PATH"],
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func layout(checkout string) (hostsFile, pricesFile, scenariosDir, fixturesDir, outDir string) {
	base := filepath.Join(checkout, "agent-tests")
	return filepath.Join(base, "hosts.yaml"),
		filepath.Join(base, "prices.yaml"),
		filepath.Join(base, "scenarios"),
		filepath.Join(base, "fixtures"),
		filepath.Join(base, "out")
}

func resolveCheckout(checkout, path string) (string, error) {
	if checkout != "" {
		return checkout, nil
	}
	res, err := proc.Run(context.Background(), proc.Spec{
		Argv: []string{"git", "rev-parse", "--show-toplevel"},
		Env: []string{
			"PATH=" + path,
			"GIT_CONFIG_NOSYSTEM=1",
			"GIT_CONFIG_GLOBAL=/dev/null",
		},
	})
	if err != nil {
		return "", fmt.Errorf("checkout: %w", err)
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("checkout: git rev-parse: exit %d\n%s", res.ExitCode, res.Stderr)
	}
	top := strings.TrimSpace(string(res.Stdout))
	if top == "" {
		return "", fmt.Errorf("checkout: git rev-parse returned an empty path")
	}
	return top, nil
}

func splitList(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func callerEnv() map[string]string {
	out := make(map[string]string)
	for _, item := range os.Environ() {
		key, val, ok := strings.Cut(item, "=")
		if !ok || key == "" {
			continue
		}
		out[key] = val
	}
	return out
}
