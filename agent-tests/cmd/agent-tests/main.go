// Command agent-tests runs contributor scenarios against coding hosts.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

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
	case "export-baseline":
		os.Exit(exportBaselineCmd(os.Args[2:], caller))
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
agent-tests export-baseline --run DIR [--run DIR ...] --out DIR [--deny-file PATH [--redact-deny]]

run    execute scenarios. Without --yes the plan is printed and no host starts.
probe  print host version, pinned model, and auth variable names.
surface-diff
       compare two refs offline. It does not start a host and needs no key.
       Exit 0 whether or not the surface differs.
report read a finished run and write report.md. It does not start a host.
       --design writes design.md. --upgrade-check runs check on a file-level delta.
export-baseline
       merge finished run directories into a baseline directory for
       report --baseline. It does not start a host.

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
  --subscription-auth
                comma-separated hosts that copy a subscription login file.
                Only codex is accepted. Empty by default.

probe flags:
  --hosts --ref --checkout --out --yes --subscription-auth
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

export-baseline flags:
  --run        a finished run directory (repeatable)
  --out        baseline directory; must not exist or must be empty
  --deny-file  local file with one term per line; an exported file that
               contains a term refuses the export. Terms are never printed.
  --redact-deny
               with --deny-file, replace each case-insensitive occurrence of a
               term in file contents with [private] instead of refusing, and
               count the replacements per file in summary.json. A file name
               that contains a term still refuses the export.
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
	subscription := fs.String("subscription-auth", "", "hosts that copy a subscription login file (only codex)")
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
	subs, err := subscriptionNames(*subscription)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
	ctx, stop := interruptContext()
	defer stop()
	err = runner.Run(ctx, runner.RunOptions{
		Checkout:         root,
		Ref:              *ref,
		BaseRef:          *base,
		ScenarioIDs:      splitList(*scenarios),
		HostNames:        splitList(*hosts),
		Runs:             *runs,
		HostsFile:        hostsFile,
		PricesFile:       pricesFile,
		ScenariosDir:     scenariosDir,
		FixturesDir:      fixturesDir,
		OutDir:           outDir,
		Keep:             *keep,
		Caller:           caller,
		DryRun:           *dry,
		Confirm:          *yes,
		SubscriptionAuth: subs,
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
	subscription := fs.String("subscription-auth", "", "hosts that copy a subscription login file (only codex)")
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
	subs, err := subscriptionNames(*subscription)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
	ctx, stop := interruptContext()
	defer stop()
	err = runner.Probe(ctx, runner.ProbeOptions{
		Checkout:         root,
		Ref:              *ref,
		HostNames:        splitList(*hosts),
		HostsFile:        hostsFile,
		PricesFile:       pricesFile,
		FixturesDir:      fixturesDir,
		OutDir:           outDir,
		Caller:           caller,
		Confirm:          *yes,
		SubscriptionAuth: subs,
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

func exportBaselineCmd(args []string, caller map[string]string) int {
	fs := flag.NewFlagSet("export-baseline", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var runs listFlag
	fs.Var(&runs, "run", "finished run directory (repeatable)")
	out := fs.String("out", "", "baseline directory; must not exist or must be empty")
	deny := fs.String("deny-file", "", "local file with one deny-listed term per line")
	redactDeny := fs.Bool("redact-deny", false, "replace deny-listed terms in file contents with [private] instead of refusing")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: agent-tests export-baseline --run DIR [--run DIR ...] --out DIR [--deny-file PATH [--redact-deny]]")
		fmt.Fprintln(os.Stderr, "merge finished runs into a baseline directory. It does not start a host.")
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
	if len(runs) == 0 || strings.TrimSpace(*out) == "" {
		fmt.Fprintln(os.Stderr, "export-baseline: --run and --out are required")
		fs.Usage()
		return 2
	}
	if *redactDeny && strings.TrimSpace(*deny) == "" {
		fmt.Fprintln(os.Stderr, "export-baseline: --redact-deny needs --deny-file")
		fs.Usage()
		return 2
	}
	user := caller["USER"]
	if user == "" {
		user = caller["LOGNAME"]
	}
	err := report.ExportBaseline(report.ExportOptions{
		Runs:       runs,
		Out:        *out,
		DenyFile:   *deny,
		RedactDeny: *redactDeny,
		Home:       caller["HOME"],
		User:       user,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// listFlag collects a repeatable string flag.
type listFlag []string

func (l *listFlag) String() string { return strings.Join(*l, ",") }

func (l *listFlag) Set(v string) error {
	*l = append(*l, v)
	return nil
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

func subscriptionNames(s string) ([]string, error) {
	names := splitList(s)
	if err := runner.ValidateSubscriptionAuth(names); err != nil {
		return nil, err
	}
	return names, nil
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

// interruptContext is cancelled by Ctrl-C or SIGTERM. The host runs in its own
// process group and does not get the terminal's signal, so the runner has to
// stop it and let its cleanup run, which removes a copied login file.
func interruptContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
