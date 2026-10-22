package cli_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/cli"
	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// pinned writes a workspace's pin and clears the variables that steer the pin
// step, so the host's own never reach a test.
func pinned(t *testing.T, root, pin string) {
	t.Helper()
	t.Setenv("CC_CLI_VERSION", "")
	t.Setenv("CC_CLI_SWITCHED_FROM", "")
	write(t, filepath.Join(root, ".context-circuit", "CLI_VERSION"), pin+"\n")
}

// runAs runs the CLI in this process as the running version given.
func runAs(root, running string, args ...string) (int, string, string) {
	var out, errOut bytes.Buffer
	code := cli.Run(context.Background(), append([]string{"--workspace", root}, args...), &out, &errOut, running)
	return code, out.String(), errOut.String()
}

// versionStore is an installer's layout: <bin>/.context-circuit-versions, one
// <version>-<os>-<arch> folder per installed version.
type versionStore struct{ bin, store string }

func newStore(t *testing.T) versionStore {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "bin")
	return versionStore{bin, filepath.Join(bin, ".context-circuit-versions")}
}

func (v versionStore) binary(version string) string {
	return filepath.Join(v.store, version+"-"+runtime.GOOS+"-"+runtime.GOARCH, "context-circuit-cli"+exeSuffix())
}

// put installs data as version's binary and returns its path.
func (v versionStore) put(t *testing.T, version string, data []byte) string {
	t.Helper()
	path := v.binary(version)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0755); err != nil {
		t.Fatal(err)
	}
	return path
}

// runFrom makes the pin step see this process as running's binary in v.
func (v versionStore) runFrom(t *testing.T, running string) {
	t.Helper()
	t.Cleanup(cli.SetExecutable(v.put(t, running, []byte("running"))))
}

// noHandOff fails the test if the pin step hands off.
func noHandOff(t *testing.T) {
	t.Helper()
	t.Cleanup(cli.StubHandOff(func(target string, _, _ []string) {
		t.Errorf("handed off to %s", target)
	}))
}

func intentFiles(t *testing.T, root string) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(root, "intent", "i*.md"))
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func TestDifferentReleaseLines(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want bool
	}{
		{"2.5.0", "2.6.1", false},
		{"2.5.0-acme.1", "2.6.0-acme.4", false},
		{"2.5.0", "2.5.0-acme.1", true},
		{"2.5.0-acme.1", "2.5.0", true},
		{"2.5.0-acme.1", "2.5.0-corp.1", true},
		// An unknown line is never a different line by itself.
		{"2.5.0", "2.0.0-dev", false},
		{"2.5.0-acme.1", "2.6.0-rc.1", false},
		{"2.5.0", "2.5.0+x", false},
		{"2.5.0-acme.1", "2.5.0-Acme.1", false},
		{"2.0.0-dev", "2.0.0-test", false},
	} {
		if got := cli.DifferentLines(c.a, c.b); got != c.want {
			t.Errorf("differentLines(%q, %q) = %v; want %v", c.a, c.b, got, c.want)
		}
	}
	for _, word := range strings.Fields(cli.ReservedLabels) {
		if kind, _ := cli.ReleaseLine("2.5.0-" + word + ".1"); kind != "unknown" {
			t.Errorf("reserved label %q is a %s line", word, kind)
		}
	}
}

func TestFindStore(t *testing.T) {
	v := newStore(t)
	binary := v.put(t, "2.5.0", []byte("cli"))
	store, bin := resolved(t, v.store), resolved(t, v.bin)
	found := func(path, running string) (string, string) {
		t.Helper()
		defer cli.SetExecutable(path)()
		return cli.FindStore(running)
	}

	if s, b := found(binary, "2.5.0"); s != store || b != bin {
		t.Fatalf("in a store: got %q, %q; want %q, %q", s, b, store, bin)
	}
	// A folder naming another version is not where the installer put this one.
	if s, _ := found(binary, "2.5.1"); s != "" {
		t.Fatalf("a folder named for another version was taken as a store: %s", s)
	}
	loose := filepath.Join(t.TempDir(), "context-circuit-cli"+exeSuffix())
	write(t, loose, "cli")
	if s, _ := found(loose, "2.5.0"); s != "" {
		t.Fatalf("a binary outside a store was taken as one: %s", s)
	}
	if s, _ := found(filepath.Join(t.TempDir(), "missing"), "2.5.0"); s != "" {
		t.Fatalf("a binary that is not there was taken as a store: %s", s)
	}

	// The PATH commands are symlinks into the store, the label one included.
	variant := v.put(t, "2.5.0-acme.1", []byte("cli"))
	for link, target := range map[string]string{"context-circuit-cli": binary, "acme-context-circuit-cli": variant} {
		path := filepath.Join(v.bin, link)
		if err := os.Symlink(target, path); err != nil {
			t.Logf("symlink case skipped, this system cannot create one: %v", err)
			continue
		}
		running := strings.TrimSuffix(filepath.Base(filepath.Dir(target)), "-"+runtime.GOOS+"-"+runtime.GOARCH)
		if s, b := found(path, running); s != store || b != bin {
			t.Fatalf("through %s: got %q, %q; want %q, %q", link, s, b, store, bin)
		}
	}
}

func TestPinValidation(t *testing.T) {
	for pin, want := range map[string]bool{
		"2.5.0": true, "2.5.0-acme.1": true, "2.6.0-rc.1": true, "2.5.0+x": true,
		"2.5.0/x": false, `2.5.0\x`: false, "2..5": false, "2.5.0/../x": false,
		"3.0.0": false, "v2.5.0": false, "": false, "2.5.0;x": false,
	} {
		if got := cli.ValidPin(pin); got != want {
			t.Errorf("validPin(%q) = %v; want %v", pin, got, want)
		}
	}

	v := newStore(t)
	if got := cli.Installed(v.store, "2.5.0"); got != "" {
		t.Fatalf("a pin that is not installed was found at %s", got)
	}
	target := v.put(t, "2.5.0", []byte("cli"))
	if got := cli.Installed(v.store, "2.5.0"); got != target {
		t.Fatalf("installed pin: got %q, want %q", got, target)
	}
	// A directory where the binary belongs.
	if err := os.Remove(target); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(target, 0755); err != nil {
		t.Fatal(err)
	}
	if got := cli.Installed(v.store, "2.5.0"); got != "" {
		t.Fatalf("a directory was taken as the pinned binary: %s", got)
	}
	if err := os.Remove(target); err != nil {
		t.Fatal(err)
	}
	// A symlink is not the file the installer verified.
	elsewhere := filepath.Join(t.TempDir(), "cli")
	write(t, elsewhere, "cli")
	if err := os.Symlink(elsewhere, target); err != nil {
		t.Logf("symlink case skipped, this system cannot create one: %v", err)
	} else if got := cli.Installed(v.store, "2.5.0"); got != "" {
		t.Fatalf("a symlink was taken as the pinned binary: %s", got)
	}
	// A pin the installer refuses never becomes a path, even where one exists.
	escape := filepath.Join(v.store, "2..5-"+runtime.GOOS+"-"+runtime.GOARCH, "context-circuit-cli"+exeSuffix())
	write(t, escape, "cli")
	if got := cli.Installed(v.store, "2..5"); got != "" {
		t.Fatalf("an invalid pin was used: %s", got)
	}
}

// Row 1: no pin, an empty pin, or the pin itself runs quietly.
func TestAMatchingPinRunsQuietly(t *testing.T) {
	f := setup(t)
	pin := filepath.Join(f.root, ".context-circuit", "CLI_VERSION")
	pinned(t, f.root, "2.5.0")
	if code, out, errOut := runAs(f.root, "2.5.0", "status"); code != 0 || errOut != "" {
		t.Fatalf("the pinned version: exit=%d %s %s", code, out, errOut)
	}
	write(t, pin, " \n")
	if code, out, errOut := runAs(f.root, "2.5.1", "record", "create", "--kind", "intent", "--slug", "a", "--title", "A"); code != 0 || errOut != "" {
		t.Fatalf("an empty pin: exit=%d %s %s", code, out, errOut)
	}
	if err := os.Remove(pin); err != nil {
		t.Fatal(err)
	}
	if code, out, errOut := runAs(f.root, "2.5.0-acme.1", "record", "create", "--kind", "intent", "--slug", "b", "--title", "B"); code != 0 || errOut != "" {
		t.Fatalf("no pin: exit=%d %s %s", code, out, errOut)
	}
}

// Row 2: CC_CLI_VERSION names no version; anything but local is a mistake.
func TestTheChoiceAcceptsOnlyLocal(t *testing.T) {
	f := setup(t)
	pinned(t, f.root, "2.5.0")
	t.Setenv("CC_CLI_VERSION", "2.5.0")
	noHandOff(t)
	code, out, errOut := runAs(f.root, "2.5.1", "record", "create", "--kind", "intent", "--slug", "a", "--title", "A")
	if code != 2 || out != "" || !strings.Contains(errOut, "CC_CLI_VERSION accepts only `local`") {
		t.Fatalf("exit=%d %s %s", code, out, errOut)
	}
	if files := intentFiles(t, f.root); len(files) != 0 {
		t.Fatalf("a refused call created %v", files)
	}
}

// Row 4: the pinned binary in this binary's store runs instead, with the line
// unchanged and the version that handed off named.
func TestTheSwitchHandsOffToThePinnedBinary(t *testing.T) {
	f := setup(t)
	pinned(t, f.root, "2.5.0")
	v := newStore(t)
	v.runFrom(t, "2.5.1")
	target := resolved(t, v.put(t, "2.5.0", []byte("pinned")))
	var calls [][]string
	var env []string
	t.Cleanup(cli.StubHandOff(func(path string, argv, environment []string) {
		calls = append(calls, argv)
		env = environment
		if !samePath(t, path, target) {
			t.Errorf("handed off to %s, want %s", path, target)
		}
	}))
	args := []string{"--workspace", f.root, "record", "list", "--json"}
	var out, errOut bytes.Buffer
	code := cli.Run(context.Background(), args, &out, &errOut, "2.5.1")
	if code != 0 || out.Len() != 0 || len(calls) != 1 {
		t.Fatalf("exit=%d calls=%d stdout=%q", code, len(calls), out.String())
	}
	if !samePath(t, calls[0][0], target) || !slices.Equal(calls[0][1:], args) {
		t.Fatalf("the pinned CLI got %q, want %s then %q", calls[0], target, args)
	}
	var marks []string
	for _, entry := range env {
		if strings.HasPrefix(entry, "CC_CLI_SWITCHED_FROM=") {
			marks = append(marks, entry)
		}
	}
	if !slices.Equal(marks, []string{"CC_CLI_SWITCHED_FROM=2.5.1"}) {
		t.Fatalf("the hand-off must name the version that handed off, once: %q", marks)
	}
	if want := "note: this workspace pins CLI 2.5.0; running it from the version store instead of 2.5.1\n"; errOut.String() != want {
		t.Fatalf("stderr: got %q, want %q", errOut.String(), want)
	}
}

// Row 3: the loop guard and CC_CLI_VERSION=local skip the switch, and local
// still never lifts a refusal.
func TestTheSwitchCanBeSkippedButNotTheGuard(t *testing.T) {
	f := setup(t)
	pinned(t, f.root, "2.5.0")
	v := newStore(t)
	v.runFrom(t, "2.5.1")
	target := resolved(t, v.put(t, "2.5.0", []byte("pinned")))
	noHandOff(t)

	t.Setenv("CC_CLI_SWITCHED_FROM", "2.4.9")
	code, out, errOut := runAs(f.root, "2.5.1", "status")
	if code != 0 || !strings.Contains(errOut, "already handed off") || !strings.Contains(errOut, "the pinned CLI is installed at `"+target+"`") {
		t.Fatalf("loop guard: exit=%d %s %s", code, out, errOut)
	}
	t.Setenv("CC_CLI_SWITCHED_FROM", "")

	t.Setenv("CC_CLI_VERSION", "local")
	if code, out, errOut := runAs(f.root, "2.5.1", "status"); code != 0 || !strings.Contains(errOut, "CC_CLI_VERSION=local keeps the running CLI") {
		t.Fatalf("local: exit=%d %s %s", code, out, errOut)
	}
	schema2(t, f, []string{"i050"}, nil)
	config, legacy := filepath.Join(f.root, "workspace.yaml"), filepath.Join(f.root, ".context-circuit", "ids.yaml")
	wantConfig, wantLegacy := read(t, config), read(t, legacy)
	code, out, errOut = runAs(f.root, "2.5.1", "record", "create", "--kind", "intent", "--slug", "a", "--title", "A")
	if code != 1 || !strings.Contains(errOut, "refused:") || !strings.Contains(errOut, "schema 3") || !strings.Contains(errOut, "nothing was changed") {
		t.Fatalf("local lifted a refusal: exit=%d %s %s", code, out, errOut)
	}
	if read(t, config) != wantConfig || read(t, legacy) != wantLegacy || exists(filepath.Join(f.root, ".context-circuit", "ids")) {
		t.Fatal("a refused call changed the workspace")
	}
}

// Rows 5 to 8, where the pinned CLI is not installed.
func TestTheGuard(t *testing.T) {
	t.Run("reads, previews and order only warn", func(t *testing.T) {
		f := setup(t)
		f.repository("api")
		intent := f.intent("billing")
		plan := f.plan(intent.ID, "billing-api", "api")
		schema2(t, f, []string{intent.ID}, []string{plan.ID})
		pinned(t, f.root, "2.5.0-acme.1")
		config, legacy := filepath.Join(f.root, "workspace.yaml"), filepath.Join(f.root, ".context-circuit", "ids.yaml")
		wantConfig, wantLegacy := read(t, config), read(t, legacy)
		for _, args := range [][]string{
			{"status"},
			{"record", "list"},
			{"record", "order", "--plan", plan.ID},
			{"record", "delete", "--member", "maya"},
		} {
			code, out, errOut := runAs(f.root, "2.5.1", args...)
			if code != 0 || strings.Contains(errOut, "refused") || !strings.Contains(errOut, "warning: this workspace pins CLI 2.5.0-acme.1 but 2.5.1 is running") {
				t.Fatalf("%v: exit=%d %s %s", args, code, out, errOut)
			}
		}
		if read(t, config) != wantConfig || read(t, legacy) != wantLegacy {
			t.Fatal("a call that writes nothing changed the workspace")
		}
	})

	t.Run("a pending split is refused", func(t *testing.T) {
		f := setup(t)
		schema2(t, f, []string{"i050"}, nil)
		pinned(t, f.root, "2.5.0")
		config, legacy := filepath.Join(f.root, "workspace.yaml"), filepath.Join(f.root, ".context-circuit", "ids.yaml")
		wantConfig, wantLegacy := read(t, config), read(t, legacy)
		code, out, errOut := runAs(f.root, "2.5.1", "record", "create", "--kind", "intent", "--slug", "a", "--title", "A")
		want := "refused: this workspace pins CLI 2.5.0 (stable) but 2.5.1 (stable) is running: this change would move the workspace to schema 3, which the pinned CLI must make; nothing was changed.\n"
		if code != 1 || out != "" || !strings.HasPrefix(errOut, want) {
			t.Fatalf("exit=%d %s %s", code, out, errOut)
		}
		if read(t, config) != wantConfig || read(t, legacy) != wantLegacy || exists(filepath.Join(f.root, ".context-circuit", "ids")) {
			t.Fatal("a refused call changed workspace.yaml or the ID file")
		}
	})

	t.Run("init on a schema 2 seed is refused", func(t *testing.T) {
		seed := filepath.Join(t.TempDir(), "seed")
		if code, out, errOut := run("template", "export", "--path", seed); code != 0 {
			t.Fatalf("export: exit=%d %s %s", code, out, errOut)
		}
		config := filepath.Join(seed, "workspace.yaml")
		write(t, config, strings.Replace(read(t, config), "version: 3", "version: 2", 1))
		pinned(t, seed, "2.5.0")
		want := read(t, config)
		code, out, errOut := runAs(seed, "2.5.1", "init", "--name", "Acme", "--purpose", "p", "--member", "maya", "--member-name", "Maya")
		if code != 1 || !strings.Contains(errOut, "refused:") || !strings.Contains(errOut, "schema 3") {
			t.Fatalf("exit=%d %s %s", code, out, errOut)
		}
		if read(t, config) != want || exists(filepath.Join(seed, "member.local.yaml")) {
			t.Fatal("a refused init changed the seed")
		}
	})

	t.Run("a change from another release line is refused", func(t *testing.T) {
		f := setup(t)
		pinned(t, f.root, "2.5.0-acme.1")
		code, out, errOut := runAs(f.root, "2.5.1", "record", "create", "--kind", "intent", "--slug", "a", "--title", "A")
		want := "refused: this workspace pins CLI 2.5.0-acme.1 (variant acme) but 2.5.1 (stable) is running: a CLI from a different release line; nothing was changed.\n"
		if code != 1 || out != "" || !strings.HasPrefix(errOut, want) {
			t.Fatalf("exit=%d %s %s", code, out, errOut)
		}
		if files := intentFiles(t, f.root); len(files) != 0 {
			t.Fatalf("a refused call created %v", files)
		}
	})

	t.Run("a same-line change warns and runs", func(t *testing.T) {
		f := setup(t)
		for _, c := range []struct{ pin, running string }{
			{"2.5.0", "2.5.1"},
			// A developer build is on no known line.
			{"2.5.0-acme.1", "2.0.0-dev"},
		} {
			pinned(t, f.root, c.pin)
			code, out, errOut := runAs(f.root, c.running, "record", "create", "--kind", "intent", "--slug", "a"+strings.ReplaceAll(c.running, ".", ""), "--title", "A")
			want := "warning: this workspace pins CLI " + c.pin + " but " + c.running + " is running, and " + c.pin + " is not installed beside this binary.\n"
			if code != 0 || !strings.HasPrefix(errOut, want) || !strings.Contains(errOut, "run it through this workspace's launcher") {
				t.Fatalf("%v: exit=%d %s %s", c, code, out, errOut)
			}
		}
	})

	t.Run("a fetch on a schema 2 workspace warns and runs", func(t *testing.T) {
		f := setup(t)
		api := f.repository("api")
		f.ok("repo", "clone", "--id", "copy", "--path", filepath.Join(f.home, "copy"), "--base", "main", "--url", api)
		schema2(t, f, nil, nil)
		pinned(t, f.root, "2.5.0")
		config := filepath.Join(f.root, "workspace.yaml")
		want := read(t, config)
		code, out, errOut := runAs(f.root, "2.5.1", "repo", "fetch", "--id", "copy")
		if code != 0 || strings.Contains(errOut, "refused") || !strings.Contains(errOut, "warning:") {
			t.Fatalf("exit=%d %s %s", code, out, errOut)
		}
		if read(t, config) != want {
			t.Fatal("a fetch moved the schema")
		}
	})
}

// Row 0: what answers without a workspace, or never opens one, never reads
// the pin, so even a mistaken CC_CLI_VERSION does not stop it.
func TestWhatNeverReachesThePin(t *testing.T) {
	f := setup(t)
	pinned(t, f.root, "2.5.0")
	t.Setenv("CC_CLI_VERSION", "mistake")
	for _, args := range [][]string{
		{"version"},
		{"--version"},
		{"help"},
		{"help", "record"},
		{"template", "export", "--path", filepath.Join(t.TempDir(), "seed")},
	} {
		if code, out, errOut := runAs(f.root, "2.5.1", args...); code != 0 || errOut != "" {
			t.Fatalf("%v: exit=%d %s %s", args, code, out, errOut)
		}
	}
	for _, args := range [][]string{
		{"record", "nosuch"},
		{"record", "create"},
		{"status", "--nosuch"},
	} {
		if code, out, errOut := runAs(f.root, "2.5.1", args...); code != 2 || strings.Contains(errOut, "CC_CLI_VERSION") {
			t.Fatalf("%v: exit=%d %s %s", args, code, out, errOut)
		}
	}
}

// The real switch, from a built CLI in a version store: Unix replaces the
// process and Windows passes the child's exit code through, so the stub's 3
// comes back either way.
func TestTheSwitchRunsThePinnedCLI(t *testing.T) {
	f := setup(t)
	pinned(t, f.root, "2.5.0")
	v := newStore(t)
	data, err := os.ReadFile(builtCLI)
	if err != nil {
		t.Fatal(err)
	}
	running := v.put(t, builtVersion, data)
	if data, err = os.ReadFile(launcherStub); err != nil {
		t.Fatal(err)
	}
	v.put(t, "2.5.0", data)
	dir := t.TempDir()
	args := []string{"--workspace", f.root, "record", "create", "--kind", "intent", "--slug", "a", "--title", "Two words"}

	invoke := func(command string) {
		t.Helper()
		record := filepath.Join(t.TempDir(), "record.json")
		var env []string
		for _, entry := range os.Environ() {
			if !strings.HasPrefix(entry, "CC_CLI_") {
				env = append(env, entry)
			}
		}
		child := exec.Command(command, args...)
		child.Dir = dir
		child.Env = append(env, "CC_LAUNCHER_STUB_RECORD="+record)
		var out, errOut bytes.Buffer
		child.Stdout, child.Stderr = &out, &errOut
		err := child.Run()
		code := 0
		if exit := (*exec.ExitError)(nil); errors.As(err, &exit) {
			code = exit.ExitCode()
		} else if err != nil {
			t.Fatal(err)
		}
		want := "note: this workspace pins CLI 2.5.0; running it from the version store instead of " + builtVersion + "\n"
		if code != 3 || out.Len() != 0 || errOut.String() != want {
			t.Fatalf("%s: exit=%d stdout=%q stderr=%q", command, code, out.String(), errOut.String())
		}
		var seen stubRecord
		if err := json.Unmarshal([]byte(read(t, record)), &seen); err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(seen.Args, args) || !samePath(t, seen.Cwd, dir) || seen.SwitchedFrom != builtVersion {
			t.Fatalf("%s: the pinned CLI saw %+v", command, seen)
		}
	}

	invoke(running)
	link := filepath.Join(v.bin, "context-circuit-cli")
	if err := os.Symlink(running, link); err != nil {
		t.Logf("PATH symlink case skipped, this system cannot create one: %v", err)
		return
	}
	invoke(link)
	if files := intentFiles(t, f.root); len(files) != 0 {
		t.Fatalf("the switching CLI created %v itself", files)
	}
}

func TestPinMessages(t *testing.T) {
	f := setup(t)
	s, err := workspace.Open(f.root)
	if err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(f.root, "workspace.yaml")
	seed := read(t, config)
	args := []string{"record", "create", "--title", "Two words", "it's"}
	steps := func(goos, bin, target string) string {
		return cli.PinSteps(s, "2.5.0", "2.5.1", goos, bin, target, args)
	}
	has := func(text string, want ...string) {
		t.Helper()
		for _, w := range want {
			if !strings.Contains(text, w) {
				t.Fatalf("missing %q in:\n%s", w, text)
			}
		}
	}
	lacks := func(text string, unwanted ...string) {
		t.Helper()
		for _, w := range unwanted {
			if strings.Contains(text, w) {
				t.Fatalf("unexpected %q in:\n%s", w, text)
			}
		}
	}
	t.Setenv("CONTEXT_CIRCUIT_BIN_DIR", "")

	// The launcher comes first, with the same arguments, quoted for its shell.
	unix := steps("linux", "", "")
	has(unix, "  run it through this workspace's launcher: `sh ", " record create --title 'Two words' 'it'\\''s'`",
		"\n  install the pinned CLI: `sh ", "install.sh", " --version 2.5.0`")
	if !strings.HasPrefix(unix, "  run it through") {
		t.Fatalf("the launcher is not named first:\n%s", unix)
	}
	windows := steps("windows", "", "")
	has(windows, "`powershell -NoProfile -ExecutionPolicy Bypass -File ", "cli.ps1", " record create --title 'Two words' 'it''s'`",
		"install.ps1", " -Version 2.5.0`")

	// Where the pin is installed, it is named instead of an install.
	installed := steps("linux", "", "/opt/store/2.5.0/context-circuit-cli")
	has(installed, "the pinned CLI is installed at `/opt/store/2.5.0/context-circuit-cli`")
	lacks(installed, "install the pinned CLI")

	// A gitlab registry, and a token that is named but never printed.
	t.Setenv("CONTEXT_CIRCUIT_TOKEN", "glpat-secretvalue")
	write(t, config, seed+"cli_registry:\n  source: gitlab\n  repository: group/cli\n  api: https://gitlab.example.com/api/v4\n")
	gitlab := steps("linux", "", "")
	has(gitlab, " --version 2.5.0 --source gitlab --repo group/cli --api https://gitlab.example.com/api/v4`",
		"(a gitlab registry needs a token with read access: set CONTEXT_CIRCUIT_TOKEN for that one command)")
	lacks(gitlab, "glpat-secretvalue")
	has(steps("windows", "", ""), " -Version 2.5.0 -Source gitlab -Repo group/cli -Api https://gitlab.example.com/api/v4`")

	// A GitHub registry other than the installer's default.
	write(t, config, seed+"cli_registry:\n  source: github\n  repository: acme/cli\n")
	has(steps("linux", "", ""), " --version 2.5.0 --source github --repo acme/cli`")

	// A store outside the installer's default directory is named.
	write(t, config, seed)
	bin := filepath.Join(t.TempDir(), "custom-bin")
	t.Setenv("HOME", t.TempDir())
	t.Setenv("LOCALAPPDATA", t.TempDir())
	has(steps("linux", bin, ""), " --bin-dir ", filepath.Base(bin))
	has(steps("windows", bin, ""), " -BinDir ", filepath.Base(bin))
	t.Setenv("CONTEXT_CIRCUIT_BIN_DIR", bin)
	lacks(steps("linux", bin, ""), "--bin-dir")

	// Without a launcher, the install is the first step.
	launcher := filepath.Join(f.root, ".context-circuit", "cli")
	saved := read(t, launcher)
	if err := os.Remove(launcher); err != nil {
		t.Fatal(err)
	}
	bare := steps("linux", "", "")
	lacks(bare, "launcher")
	if !strings.HasPrefix(bare, "  install the pinned CLI: `sh ") {
		t.Fatalf("the install is not first:\n%s", bare)
	}
	write(t, launcher, saved)

	// A workspace.yaml this CLI cannot read gives the plain install, and says
	// where the registry would have come from.
	write(t, config, "version: [\n")
	broken := steps("linux", "", "")
	has(broken, " --version 2.5.0`", "read cli_registry in workspace.yaml first")
	lacks(broken, "--source")
}

// The release-line rule is written once in Go and once in each script that
// names a build; these lines are where the scripts keep it.
func TestReleaseLineMatchesScripts(t *testing.T) {
	sh := [2]string{"reserved_labels='", "label_pattern='"}
	for file, prefixes := range map[string][2]string{
		"../../scripts/release-version.sh":                {sh[0], sh[1]},
		"../../scripts/build-cli.sh":                      {sh[0], sh[1]},
		"../../product/skills/cc-cli/scripts/install.sh":  {sh[0], sh[1]},
		"../../product/skills/cc-cli/scripts/install.ps1": {"$reservedLabels = '", "$labelPattern = '"},
	} {
		data, err := os.ReadFile(filepath.FromSlash(file))
		if err != nil {
			t.Fatal(err)
		}
		for i, want := range []string{cli.ReservedLabels, cli.LabelPattern} {
			var found []string
			scanner := bufio.NewScanner(bytes.NewReader(data))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if value, ok := strings.CutPrefix(line, prefixes[i]); ok && strings.HasSuffix(value, "'") {
					found = append(found, strings.TrimSuffix(value, "'"))
				}
			}
			if !slices.Equal(found, []string{want}) {
				t.Errorf("%s: %s…' lines hold %q; want exactly one holding %q", file, prefixes[i], found, want)
			}
		}
	}
}
