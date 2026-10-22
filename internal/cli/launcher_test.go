package cli_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// launcherStub stands in for an installed CLI: it records how it was started
// into the file CC_LAUNCHER_STUB_RECORD names, and exits 3 so a test can tell
// the exit code was passed through rather than made up.
var launcherStub string

// builtCLI is this repository's CLI built as builtVersion, which the switch
// tests install into a version store and run as a separate process: a switch
// inside the test process would replace the test.
var builtCLI string

const builtVersion = "2.5.1"

const launcherStubSource = `package main

import (
	"encoding/json"
	"io"
	"os"
)

func main() {
	stdin, _ := io.ReadAll(os.Stdin)
	cwd, _ := os.Getwd()
	data, _ := json.Marshal(map[string]any{"args": os.Args[1:], "cwd": cwd, "stdin": string(stdin), "switched_from": os.Getenv("CC_CLI_SWITCHED_FROM")})
	if err := os.WriteFile(os.Getenv("CC_LAUNCHER_STUB_RECORD"), data, 0600); err != nil {
		os.Exit(100)
	}
	os.Exit(3)
}
`

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "cc-launcher-stub-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	source := filepath.Join(dir, "main.go")
	launcherStub = filepath.Join(dir, "stub"+exeSuffix())
	if err := os.WriteFile(source, []byte(launcherStubSource), 0600); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if data, err := exec.Command("go", "build", "-o", launcherStub, source).CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build launcher stub: %v %s\n", err, data)
		os.Exit(1)
	}
	builtCLI = filepath.Join(dir, "context-circuit-cli"+exeSuffix())
	build := exec.Command("go", "build", "-ldflags", "-X main.version="+builtVersion, "-o", builtCLI, "github.com/kaotypr/context-circuit-source/cmd/context-circuit")
	if data, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build CLI: %v %s\n", err, data)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

type stubRecord struct {
	Args         []string `json:"args"`
	Cwd          string   `json:"cwd"`
	Stdin        string   `json:"stdin"`
	SwitchedFrom string   `json:"switched_from"`
}

// launcherFixture is an initialized workspace carrying the seed's launcher, and a
// version store the launcher is pointed at.
type launcherFixture struct {
	fixture
	kind  string
	bin   string
	store string
}

func newLauncherFixture(t *testing.T, kind string) launcherFixture {
	t.Helper()
	f := launcherFixture{fixture: setup(t), kind: kind}
	f.bin = filepath.Join(f.home, "bin")
	f.store = filepath.Join(f.bin, ".context-circuit-versions")
	f.pin("2.5.0")
	f.install("2.5.0")
	return f
}

func (f launcherFixture) pin(version string) {
	f.t.Helper()
	write(f.t, filepath.Join(f.root, ".context-circuit", "CLI_VERSION"), version+"\n")
}

// target is where the launcher looks for a version on this machine.
func (f launcherFixture) target(version string) string {
	return filepath.Join(f.store, version+"-"+runtime.GOOS+"-"+runtime.GOARCH, "context-circuit-cli"+exeSuffix())
}

func (f launcherFixture) install(version string) {
	f.t.Helper()
	data, err := os.ReadFile(launcherStub)
	if err != nil {
		f.t.Fatal(err)
	}
	target := f.target(version)
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		f.t.Fatal(err)
	}
	if err := os.WriteFile(target, data, 0755); err != nil {
		f.t.Fatal(err)
	}
}

func (f launcherFixture) uninstall(version string) {
	f.t.Helper()
	if err := os.RemoveAll(filepath.Dir(f.target(version))); err != nil {
		f.t.Fatal(err)
	}
}

// path is how a caller of this launcher names a directory. Git Bash's dirname
// does not split on a backslash, so its callers use forward slashes.
func (f launcherFixture) path(path string) string {
	if f.kind == "sh" {
		return filepath.ToSlash(path)
	}
	return path
}

// launch runs the launcher the documented way from dir, and reports its exit
// code, its stderr, and what the stub saw, nil when the stub never ran.
func (f launcherFixture) launch(dir, stdin string, args ...string) (int, string, *stubRecord) {
	f.t.Helper()
	record := filepath.Join(f.t.TempDir(), "record.json")
	var command *exec.Cmd
	env := append(os.Environ(), "CONTEXT_CIRCUIT_BIN_DIR="+f.bin, "CC_LAUNCHER_STUB_RECORD="+record)
	switch f.kind {
	case "sh":
		command = exec.Command("sh", append([]string{f.path(filepath.Join(f.root, ".context-circuit", "cli"))}, args...)...)
	case "powershell":
		script := filepath.Join(f.root, ".context-circuit", "cli.ps1")
		command = exec.Command("powershell.exe", append([]string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File", script}, args...)...)
		// See TestInstallUpdateRollbackAndChecksumProtection in internal/installer.
		env = append(env, "PSModulePath=")
	}
	command.Dir = dir
	command.Env = env
	command.Stdin = strings.NewReader(stdin)
	var stderr strings.Builder
	command.Stderr = &stderr
	err := command.Run()
	code := 0
	if exit := (*exec.ExitError)(nil); errors.As(err, &exit) {
		code = exit.ExitCode()
	} else if err != nil {
		f.t.Fatalf("%s launcher did not start: %v", f.kind, err)
	}
	data, err := os.ReadFile(record)
	if os.IsNotExist(err) {
		return code, stderr.String(), nil
	} else if err != nil {
		f.t.Fatal(err)
	}
	var seen stubRecord
	if err := json.Unmarshal(data, &seen); err != nil {
		f.t.Fatal(err)
	}
	return code, stderr.String(), &seen
}

// ran asserts the stub ran with --workspace <root> first, then args unchanged.
func (f launcherFixture) ran(code int, stderr string, seen *stubRecord, dir string, args ...string) {
	f.t.Helper()
	if seen == nil {
		f.t.Fatalf("%v: the pinned CLI never ran (exit %d): %s", args, code, stderr)
	}
	if code != 3 {
		f.t.Fatalf("%v: exit code %d, want the CLI's 3: %s", args, code, stderr)
	}
	if len(seen.Args) < 2 || seen.Args[0] != "--workspace" || !samePath(f.t, seen.Args[1], f.root) {
		f.t.Fatalf("%v: the CLI must get --workspace %s first: %q", args, f.root, seen.Args)
	}
	if !slices.Equal(seen.Args[2:], args) {
		f.t.Fatalf("original arguments changed: got %q, want %q", seen.Args[2:], args)
	}
	if !samePath(f.t, seen.Cwd, dir) {
		f.t.Fatalf("working directory: got %s, want %s", seen.Cwd, dir)
	}
}

// refused asserts the launcher refused and the stub never ran.
func (f launcherFixture) refused(code int, stderr string, seen *stubRecord, want string) {
	f.t.Helper()
	if seen != nil {
		f.t.Fatalf("a refused launch ran the CLI with %q", seen.Args)
	}
	if code != 1 || !strings.Contains(stderr, "refused") || !strings.Contains(stderr, want) {
		f.t.Fatalf("want a refusal naming %q, got exit %d: %s", want, code, stderr)
	}
}

func samePath(t *testing.T, a, b string) bool {
	t.Helper()
	resolve := func(path string) string {
		resolved, err := filepath.EvalSymlinks(filepath.Clean(path))
		if err != nil {
			t.Fatalf("resolve %s: %v", path, err)
		}
		return resolved
	}
	a, b = resolve(a), resolve(b)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(a, b)
	}
	return a == b
}

func launcherKinds(t *testing.T) []string {
	t.Helper()
	// The gate runs this on Windows under Git Bash's sh, which is the only
	// place the sh launcher's Windows branch runs at all; a missing sh there
	// must fail rather than pass by skipping it.
	if _, err := exec.LookPath("sh"); err != nil {
		t.Fatalf("sh is needed to run the sh launcher: %v", err)
	}
	if runtime.GOOS == "windows" {
		return []string{"sh", "powershell"}
	}
	return []string{"sh"}
}

func TestLauncherRunsExactlyThePinnedCLI(t *testing.T) {
	for _, kind := range launcherKinds(t) {
		t.Run(kind, func(t *testing.T) {
			f := newLauncherFixture(t, kind)
			args := []string{"record", "list", "--json", "two words"}

			code, stderr, seen := f.launch(f.root, "piped input\n", args...)
			f.ran(code, stderr, seen, f.root, args...)
			if seen.Stdin != "piped input\n" {
				t.Fatalf("stdin did not reach the CLI: %q", seen.Stdin)
			}

			// The pin and the root are the launcher's own, not the working
			// directory's.
			elsewhere := t.TempDir()
			write(t, filepath.Join(elsewhere, ".context-circuit", "CLI_VERSION"), "2.4.9\n")
			code, stderr, seen = f.launch(elsewhere, "", "check")
			f.ran(code, stderr, seen, elsewhere, "check")
		})
	}
}

func TestLauncherAcceptsOnlyItsOwnWorkspace(t *testing.T) {
	for _, kind := range launcherKinds(t) {
		t.Run(kind, func(t *testing.T) {
			f := newLauncherFixture(t, kind)
			parent := filepath.Dir(f.root)
			name := filepath.Base(f.root)
			absolute := f.path(f.root)
			for _, args := range [][]string{
				{"--workspace", name, "check"},
				{"--workspace=" + name, "check"},
				{"--workspace", absolute, "check"},
				{"check", "--workspace=" + absolute},
				{"check", "-workspace", name},
				{"-workspace=" + absolute, "check", "--json"},
			} {
				code, stderr, seen := f.launch(parent, "", args...)
				f.ran(code, stderr, seen, parent, args...)
			}

			// A later --workspace overrides an earlier one, so another
			// workspace named anywhere would run under this one's pin.
			other := t.TempDir()
			for _, args := range [][]string{
				{"--workspace", f.path(other), "check"},
				{"check", "--workspace=" + f.path(other)},
				{"check", "-workspace", f.path(other)},
				{"--workspace", name, "check", "--workspace", f.path(other)},
			} {
				code, stderr, seen := f.launch(parent, "", args...)
				f.refused(code, stderr, seen, "another workspace")
			}
			code, stderr, seen := f.launch(parent, "", "--workspace", f.path(filepath.Join(other, "missing")), "check")
			f.refused(code, stderr, seen, "missing")
		})
	}
}

func TestLauncherRefusesWhatItCannotTrust(t *testing.T) {
	for _, kind := range launcherKinds(t) {
		t.Run(kind, func(t *testing.T) {
			f := newLauncherFixture(t, kind)

			// Another version in the same store is never a stand-in.
			f.uninstall("2.5.0")
			f.install("2.4.9")
			code, stderr, seen := f.launch(f.root, "", "check")
			f.refused(code, stderr, seen, "Install exactly 2.5.0")
			if !strings.Contains(stderr, "CONTEXT_CIRCUIT_BIN_DIR") {
				t.Fatalf("the refusal must name the bin directory override: %s", stderr)
			}

			// A directory where the binary belongs.
			if err := os.MkdirAll(f.target("2.5.0"), 0755); err != nil {
				t.Fatal(err)
			}
			code, stderr, seen = f.launch(f.root, "", "check")
			f.refused(code, stderr, seen, "Install exactly 2.5.0")
			f.uninstall("2.5.0")

			// A symlink to a binary is not a binary the installer verified.
			if err := os.MkdirAll(filepath.Dir(f.target("2.5.0")), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(launcherStub, f.target("2.5.0")); err != nil {
				t.Logf("symlink case skipped, this system cannot create one: %v", err)
			} else {
				code, stderr, seen = f.launch(f.root, "", "check")
				f.refused(code, stderr, seen, "Install exactly 2.5.0")
			}
			f.uninstall("2.5.0")

			// A pin the installer would refuse is refused before any path is
			// built from it, even where a binary sits at that path.
			for _, pin := range []string{"2.5.0/x", `2.5.0\x`, "2..5", "3.0.0", "v2.5.0", "2.5.0;x"} {
				f.pin(pin)
				target := f.target(pin)
				if os.MkdirAll(filepath.Dir(target), 0755) == nil {
					if data, err := os.ReadFile(launcherStub); err == nil {
						_ = os.WriteFile(target, data, 0755)
					}
				}
				code, stderr, seen = f.launch(f.root, "", "check")
				if seen != nil || code != 1 || !strings.Contains(stderr, "refused") {
					t.Fatalf("pin %q: want a refusal, got exit %d, ran=%v: %s", pin, code, seen != nil, stderr)
				}
			}
		})
	}
}

// The setup remedy a dispatch reports runs the version the workspace pins when
// the workspace carries a launcher, and the running CLI's own command otherwise.
func TestDispatchSetupNamesTheLauncher(t *testing.T) {
	f := setup(t)
	if err := os.RemoveAll(filepath.Join(f.root, ".claude", "agents")); err != nil {
		t.Fatal(err)
	}
	root, err := filepath.EvalSymlinks(f.root)
	if err != nil {
		t.Fatal(err)
	}
	setupRequired := func() string {
		t.Helper()
		var spec workspace.Dispatch
		out := f.ok("agent", "dispatch", "--host", "claude-code", "--role", "explorer",
			"--path", f.root, "--task", "Investigate.")
		if err := json.Unmarshal([]byte(out), &spec); err != nil {
			t.Fatal(err)
		}
		return spec.SetupRequired
	}

	launcher := filepath.Join(root, ".context-circuit", "cli")
	want := "sh " + launcher + " agent setup --host claude-code"
	if runtime.GOOS == "windows" {
		want = "powershell -NoProfile -ExecutionPolicy Bypass -File " + launcher + ".ps1 agent setup --host claude-code"
	}
	if got := setupRequired(); got != want {
		t.Fatalf("setup remedy with a launcher:\n got %q\nwant %q", got, want)
	}

	if err := os.Remove(launcher); err != nil {
		t.Fatal(err)
	}
	want = "context-circuit-cli --workspace " + root + " agent setup --host claude-code"
	if got := setupRequired(); got != want {
		t.Fatalf("setup remedy without a launcher:\n got %q\nwant %q", got, want)
	}
}
