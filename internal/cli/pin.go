package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// A workspace pins its CLI in .context-circuit/CLI_VERSION. Versions install
// side by side, but a PATH command runs whichever was installed last, so the
// running CLI is often not the pinned one. Warning and carrying on is how a
// newer CLI used to move a schema 2 workspace to schema 3 and lock out every
// teammate still on the pin. So a CLI that is not the pin hands off to it when
// it is installed beside this one, and otherwise refuses, before anything is
// written, a change that would migrate the schema or cross release lines.
const (
	pinFile = ".context-circuit/CLI_VERSION"
	// switchedFrom is set on a CLI started by the hand-off, naming the one that
	// handed off. A pinned binary reporting another version than its folder
	// name would otherwise hand off again, forever.
	switchedFrom = "CC_CLI_SWITCHED_FROM"
	// cliChoice set to local keeps the running CLI. It is the only value, and
	// it never lifts a refusal.
	cliChoice = "CC_CLI_VERSION"
	// versionStore is the folder the installers keep every version in, one
	// <version>-<os>-<arch> folder each, beside the PATH commands.
	versionStore  = ".context-circuit-versions"
	storeBinary   = "context-circuit-cli"
	defaultSource = "kaotypr/context-circuit-source"
	defaultAPI    = "https://api.github.com"
)

// Seams for tests: where this binary is, and the hand-off itself, which on
// Unix replaces the test process. handOff and its own seam are per platform.
var executable = os.Executable

// installablePin is the installers' rule for a version they accept, which is
// also all that keeps a pin from naming a path outside the version store.
var installablePin = regexp.MustCompile(`^2\.[A-Za-z0-9.+-]*$`)

func validPin(pin string) bool {
	return installablePin.MatchString(pin) && !strings.Contains(pin, "..")
}

// differentLines reports two versions on different release lines. A version
// of an unknown line is never one by itself, so developer builds and release
// candidates are held back only by the schema rule.
func differentLines(a, b string) bool {
	kindA, labelA := releaseLine(a)
	kindB, labelB := releaseLine(b)
	if kindA == "unknown" || kindB == "unknown" {
		return false
	}
	return kindA != kindB || labelA != labelB
}

func lineName(version string) string {
	kind, label := releaseLine(version)
	if kind == "variant" {
		return "variant " + label
	}
	return kind
}

func binaryName(goos string) string {
	if goos == "windows" {
		return storeBinary + ".exe"
	}
	return storeBinary
}

// findStore reports the version store this binary runs from, and the command
// directory holding it. It is found only where the installers put a binary,
// <bin>/.context-circuit-versions/<running>-<os>-<arch>/context-circuit-cli;
// a source build, a test binary, or a folder that names another version is in
// no store, and nothing is switched to from it.
func findStore(running string) (store, bin string) {
	path, err := executable()
	if err != nil {
		return "", ""
	}
	// macOS may report the PATH symlink rather than the file it points at.
	if path, err = filepath.EvalSymlinks(path); err != nil {
		return "", ""
	}
	folder := filepath.Dir(path)
	store = filepath.Dir(folder)
	if filepath.Base(path) != binaryName(runtime.GOOS) ||
		filepath.Base(folder) != running+"-"+runtime.GOOS+"-"+runtime.GOARCH ||
		filepath.Base(store) != versionStore {
		return "", ""
	}
	return store, filepath.Dir(store)
}

// installed reports the pinned binary in store, or "" when it is not there.
// Only a regular file counts: the installer verified that file, so a cloned
// workspace can choose among binaries this person installed and nothing else.
func installed(store, pin string) string {
	if store == "" || !validPin(pin) {
		return ""
	}
	target := filepath.Join(store, pin+"-"+runtime.GOOS+"-"+runtime.GOARCH, binaryName(runtime.GOOS))
	if info, err := os.Lstat(target); err != nil || !info.Mode().IsRegular() {
		return ""
	}
	return target
}

// writes reports whether a call can change the workspace: its command is not
// read-only, and the call is not a preview.
func writes(cmd command, in *input) bool {
	return !cmd.readOnly && (cmd.previews == nil || !cmd.previews(in))
}

// migration names the change to the workspace's layout a call would make,
// which only the pinned CLI makes, or "" when it would make none. A check that
// fails counts as nothing pending, and the command then reports that failure
// itself.
func migration(s *workspace.Store, name string, cmd command, in *input) string {
	split := func() bool {
		pending, err := s.IDSplitPending()
		return err == nil && pending
	}
	// The schema stays where it is for this one: a top-level archive/ reads as
	// archived before and after, but a CLI that does not know _archive/ must
	// not be the one to move the records there.
	rename := func() bool {
		pending, err := s.ArchiveRenamePending()
		return err == nil && pending
	}
	schema := fmt.Sprintf("this change would move the workspace to schema %d, which the pinned CLI must make", workspace.CurrentSchema)
	// Initialization writes the current schema itself, whatever its seed
	// records, and a blank seed fails the strict decode IDSplitPending uses.
	if name == "init" {
		if recorded, err := s.RecordedSchema(); err == nil && recorded > 0 && recorded < workspace.CurrentSchema || split() {
			return schema
		}
	} else if !splitEnabled(cmd, in) {
		return ""
	} else if split() {
		return schema
	}
	if rename() {
		return "this change would rename the workspace's archive/ folders to _archive/, which the pinned CLI must make"
	}
	return ""
}

// pinState is what the messages about one pin are written from.
type pinState struct {
	root, pin, running, goos string
	store, bin               string   // "" when this binary is in no version store
	target                   string   // the pinned binary in that store, "" when not installed
	skipped                  string   // why the switch was skipped although the pin is installed
	args                     []string // the command line, as given
}

// applyPin runs the pin step for a command that reads a workspace. stop is true
// when the command must not run here, with the exit code to return: the pinned
// CLI ran instead, or the call was refused and nothing was changed.
func applyPin(s *workspace.Store, cmd command, in *input, original []string, running string, out, errOut io.Writer) (code int, stop bool) {
	// Only the pin is read before a switch, so a workspace.yaml this CLI cannot
	// decode never stops the hand-off to one that can.
	data, err := s.Read(pinFile)
	if err != nil {
		return 0, false
	}
	pin := strings.TrimSpace(string(data))
	if pin == "" || pin == running {
		return 0, false
	}
	choice := os.Getenv(cliChoice)
	if choice != "" && choice != "local" {
		fmt.Fprintf(errOut, "%s accepts only `local`, which keeps the running CLI; nothing was run\n", cliChoice)
		return 2, true
	}
	p := pinState{root: s.Root, pin: pin, running: running, goos: runtime.GOOS, args: original}
	p.store, p.bin = findStore(running)
	p.target = installed(p.store, pin)
	if p.target != "" {
		switch from := os.Getenv(switchedFrom); {
		case from != "":
			p.skipped = fmt.Sprintf("%s=%s says this run was already handed off once, so it does not switch again", switchedFrom, from)
		case choice == "local":
			p.skipped = cliChoice + "=local keeps the running CLI"
		default:
			fmt.Fprintf(errOut, "note: this workspace pins CLI %s; running it from the version store instead of %s\n", pin, running)
			code, err := handOff(p.target, original, switchedEnv(running), out, errOut)
			if err != nil {
				fmt.Fprintf(errOut, "the pinned CLI %s could not be started: %v; nothing was changed\n", p.target, err)
				return 1, true
			}
			return code, true
		}
	}
	if writes(cmd, in) {
		reason := migration(s, in.name, cmd, in)
		if reason == "" && differentLines(pin, running) {
			reason = "a CLI from a different release line"
		}
		if reason != "" {
			fmt.Fprint(errOut, p.refusal(s, reason))
			return 1, true
		}
	}
	fmt.Fprint(errOut, p.warning(s))
	return 0, false
}

// switchedEnv is this process's environment for the pinned CLI, naming the
// version that handed off in place of any earlier value.
func switchedEnv(running string) []string {
	var env []string
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, switchedFrom+"=") {
			env = append(env, entry)
		}
	}
	return append(env, switchedFrom+"="+running)
}

func (p pinState) warning(s *workspace.Store) string {
	var head string
	switch {
	case p.skipped != "":
		head = fmt.Sprintf("warning: this workspace pins CLI %s but %s is running; %s.", p.pin, p.running, p.skipped)
	case p.store == "":
		head = fmt.Sprintf("warning: this workspace pins CLI %s but %s is running, and %s is not installed beside this binary.", p.pin, p.running, p.pin)
	default:
		head = fmt.Sprintf("warning: this workspace pins CLI %s but %s is running, and %s is not installed in %s.", p.pin, p.running, p.pin, p.store)
	}
	return head + "\n" + p.nextSteps(s)
}

func (p pinState) refusal(s *workspace.Store, reason string) string {
	return fmt.Sprintf("refused: this workspace pins CLI %s (%s) but %s (%s) is running: %s; nothing was changed.\n",
		p.pin, lineName(p.pin), p.running, lineName(p.running), reason) + p.nextSteps(s)
}

// nextSteps names how to run the pinned CLI: the workspace's launcher first,
// which runs the pin whatever else is installed, then how to install the pin
// where it is missing, or where it already is when the switch was skipped.
func (p pinState) nextSteps(s *workspace.Store) string {
	var lines []string
	if launcher := filepath.Join(p.root, ".context-circuit", "cli"); regular(launcher) {
		command := "sh " + shellQuote(launcher)
		if p.goos == "windows" {
			command = "powershell -NoProfile -ExecutionPolicy Bypass -File " + powerShellQuote(launcher+".ps1")
		}
		lines = append(lines, "run it through this workspace's launcher: `"+command+p.quoted(p.args)+"`")
	}
	if p.target != "" {
		lines = append(lines, "the pinned CLI is installed at `"+p.target+"`")
	} else {
		lines = append(lines, p.install(s)...)
	}
	return "  " + strings.Join(lines, "\n  ") + "\n"
}

// install names the installer call that puts the pin beside this binary, from
// the registry the workspace records. A token is named, never printed.
func (p pinState) install(s *workspace.Store) []string {
	type option struct{ name, value string }
	var options []option
	var notes []string
	cfg, err := s.Config()
	switch registry := cfg.CLIRegistry; {
	case err != nil:
		notes = append(notes, "read cli_registry in workspace.yaml first: it names where this workspace's CLI is published, and workspace.yaml could not be read here")
	case registry == nil:
	case registry.Source == "gitlab":
		options = append(options, option{"source", "gitlab"}, option{"repo", registry.Repository}, option{"api", registry.API})
		notes = append(notes, "(a gitlab registry needs a token with read access: set CONTEXT_CIRCUIT_TOKEN for that one command)")
	case registry.Repository != defaultSource || registry.API != "" && registry.API != defaultAPI:
		options = append(options, option{"source", "github"}, option{"repo", registry.Repository})
		if registry.API != "" && registry.API != defaultAPI {
			options = append(options, option{"api", registry.API})
		}
	}
	if p.bin != "" && !sameDir(p.bin, defaultBinDir(p.goos), p.goos) {
		options = append(options, option{"bin-dir", p.bin})
	}
	script := filepath.Join(p.root, ".agents", "skills", "cc-cli", "scripts", "install.sh")
	command := "sh " + shellQuote(script) + " --version " + shellQuote(p.pin)
	if p.goos == "windows" {
		script = strings.TrimSuffix(script, ".sh") + ".ps1"
		command = "powershell -NoProfile -ExecutionPolicy Bypass -File " + powerShellQuote(script) + " -Version " + powerShellQuote(p.pin)
	}
	for _, o := range options {
		if p.goos == "windows" {
			command += " -" + powerShellName(o.name) + " " + powerShellQuote(o.value)
		} else {
			command += " --" + o.name + " " + shellQuote(o.value)
		}
	}
	return append([]string{"install the pinned CLI: `" + command + "`"}, notes...)
}

// defaultBinDir is where the installer puts commands when it is given no
// --bin-dir: CONTEXT_CIRCUIT_BIN_DIR, else its platform default.
func defaultBinDir(goos string) string {
	if dir := os.Getenv("CONTEXT_CIRCUIT_BIN_DIR"); dir != "" {
		return dir
	}
	if goos == "windows" {
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "ContextCircuit", "bin")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "bin")
}

// sameDir compares directories the way the installer resolves them, through
// symlinks, and without case on Windows.
func sameDir(a, b, goos string) bool {
	if b == "" {
		return false
	}
	resolve := func(path string) string {
		if resolved, err := filepath.EvalSymlinks(path); err == nil {
			return resolved
		}
		return filepath.Clean(path)
	}
	a, b = resolve(a), resolve(b)
	if goos == "windows" {
		return strings.EqualFold(a, b)
	}
	return a == b
}

func regular(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular()
}

// quoted renders arguments for the shell the next step names.
func (p pinState) quoted(args []string) string {
	var b strings.Builder
	for _, arg := range args {
		b.WriteByte(' ')
		if p.goos == "windows" {
			b.WriteString(powerShellQuote(arg))
		} else {
			b.WriteString(shellQuote(arg))
		}
	}
	return b.String()
}

var shellSafe = regexp.MustCompile(`^[A-Za-z0-9_@%+=:,./-]+$`)

// shellQuote quotes a word for a POSIX shell where it needs it.
func shellQuote(word string) string {
	if shellSafe.MatchString(word) {
		return word
	}
	return "'" + strings.ReplaceAll(word, "'", `'\''`) + "'"
}

var powerShellSafe = regexp.MustCompile(`^[A-Za-z0-9_+=:,./\\-]+$`)

// powerShellQuote quotes a word for PowerShell where it needs it.
func powerShellQuote(word string) string {
	if powerShellSafe.MatchString(word) {
		return word
	}
	return "'" + strings.ReplaceAll(word, "'", "''") + "'"
}

// powerShellName is install.ps1's parameter for an install.sh option.
func powerShellName(option string) string {
	switch option {
	case "bin-dir":
		return "BinDir"
	case "api":
		return "Api"
	}
	return strings.ToUpper(option[:1]) + option[1:]
}
