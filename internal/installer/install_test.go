package installer_test

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	assets "github.com/kaotypr/context-circuit-source"
)

func TestInstallUpdateRollbackAndChecksumProtection(t *testing.T) {
	root := t.TempDir()
	files, err := assets.Files()
	if err != nil {
		t.Fatal(err)
	}
	script := "install.sh"
	if runtime.GOOS == "windows" {
		script = "install.ps1"
	}
	scriptPath := filepath.Join(root, script)
	if err := os.WriteFile(scriptPath, files[".agents/skills/cc-cli/scripts/"+script], 0600); err != nil {
		t.Fatal(err)
	}
	binDir := filepath.Join(root, "bin with spaces")
	install := func(version, archive, sums string, success bool) {
		t.Helper()
		args := []string{scriptPath, "--version", version, "--bin-dir", binDir, "--archive", archive, "--checksums", sums}
		program := "sh"
		if runtime.GOOS == "windows" {
			program = "powershell.exe"
			args = []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath, "-Version", version, "-BinDir", binDir, "-Archive", archive, "-Checksums", sums}
		}
		command := exec.Command(program, args...)
		if runtime.GOOS == "windows" {
			// A PowerShell 7 parent, as CI uses, exports a PSModulePath holding
			// only its own modules. Windows PowerShell then cannot autoload the
			// cmdlets it ships with, and the documented invocation fails for a
			// reason the installer does not control. Clear it so the child
			// rebuilds its own default, which is what a user's shell supplies.
			command.Env = append(os.Environ(), "PSModulePath=")
		}
		data, err := command.CombinedOutput()
		if (err == nil) != success {
			t.Fatalf("installer success=%v: %v\n%s", success, err, data)
		}
	}
	checkVersion := func(expected string) {
		t.Helper()
		command := exec.Command(filepath.Join(binDir, "context-circuit-cli"), "version")
		if runtime.GOOS == "windows" {
			command = exec.Command("cmd.exe", "/d", "/c", filepath.Join(binDir, "context-circuit-cli.cmd"), "version")
		}
		data, err := command.CombinedOutput()
		if err != nil || strings.TrimSpace(string(data)) != expected {
			t.Fatalf("installed command: %v %s", err, data)
		}
	}
	firstArchive, firstSums := packageFixture(t, root, publicCommand, "2.0.0-install1")
	secondArchive, secondSums := packageFixture(t, root, publicCommand, "2.0.0-install2")
	install("2.0.0-install1", firstArchive, firstSums, true)
	checkVersion("2.0.0-install1")
	install("2.0.0-install2", secondArchive, secondSums, true)
	checkVersion("2.0.0-install2")
	install("2.0.0-install1", firstArchive, firstSums, true)
	checkVersion("2.0.0-install1")
	bad := filepath.Join(root, "bad-sums")
	if err := os.WriteFile(bad, []byte(strings.Repeat("0", 64)+"  ./"+filepath.Base(secondArchive)+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	install("2.0.0-install2", secondArchive, bad, false)
	checkVersion("2.0.0-install1")
	// A valid checksum for a package reporting the wrong version is rejected.
	wrongArchive, wrongSums := packageFixture(t, root, publicCommand, "2.0.0-claimed", "2.0.0-actual")
	install("2.0.0-claimed", wrongArchive, wrongSums, false)
	checkVersion("2.0.0-install1")
}

const publicCommand = "context-circuit-cli"

// installer runs the shipped installer for this platform offline. An empty
// binDir leaves the flag out, so the environment decides.
type installer struct {
	t      *testing.T
	script string
}

func newInstaller(t *testing.T, root string) installer {
	t.Helper()
	files, err := assets.Files()
	if err != nil {
		t.Fatal(err)
	}
	script := "install.sh"
	if runtime.GOOS == "windows" {
		script = "install.ps1"
	}
	scriptPath := filepath.Join(root, script)
	if err := os.WriteFile(scriptPath, files[".agents/skills/cc-cli/scripts/"+script], 0600); err != nil {
		t.Fatal(err)
	}
	return installer{t: t, script: scriptPath}
}

func (i installer) run(env []string, version, binDir, archive, sums string, success bool) {
	i.t.Helper()
	args := []string{i.script, "--version", version, "--archive", archive, "--checksums", sums}
	if binDir != "" {
		args = append(args, "--bin-dir", binDir)
	}
	program := "sh"
	if runtime.GOOS == "windows" {
		program = "powershell.exe"
		args = []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File", i.script, "-Version", version, "-Archive", archive, "-Checksums", sums}
		if binDir != "" {
			args = append(args, "-BinDir", binDir)
		}
		// See TestInstallUpdateRollbackAndChecksumProtection.
		env = append(env, "PSModulePath=")
	}
	command := exec.Command(program, args...)
	command.Env = append(append(os.Environ(), "CONTEXT_CIRCUIT_BIN_DIR="), env...)
	data, err := command.CombinedOutput()
	if (err == nil) != success {
		i.t.Fatalf("installer %s success=%v: %v\n%s", version, success, err, data)
	}
}

// installedVersion runs a PATH command the installer wrote and returns what it
// reports, or "" when there is no such command.
func installedVersion(t *testing.T, binDir, name string) string {
	t.Helper()
	path := filepath.Join(binDir, name)
	command := exec.Command(path, "version")
	if runtime.GOOS == "windows" {
		path += ".cmd"
		command = exec.Command("cmd.exe", "/d", "/c", path, "version")
	}
	if _, err := os.Lstat(path); err != nil {
		return ""
	}
	data, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v %s", path, err, data)
	}
	return strings.TrimSpace(string(data))
}

// A variant version installs under its own command, from assets carrying that
// name, and a public and a variant install never repoint each other.
func TestVariantInstallsUnderItsOwnCommand(t *testing.T) {
	root := t.TempDir()
	install := newInstaller(t, root)
	binDir := filepath.Join(root, "bin")
	const variant = "acme-context-circuit-cli"
	stableArchive, stableSums := packageFixture(t, root, publicCommand, "2.5.0")
	install.run(nil, "2.5.0", binDir, stableArchive, stableSums, true)

	// An archive under the public name is not the variant's asset.
	misnamed, misnamedSums := packageFixture(t, root, publicCommand, "2.5.0-acme.1")
	install.run(nil, "2.5.0-acme.1", binDir, misnamed, misnamedSums, false)
	if got := installedVersion(t, binDir, variant); got != "" {
		t.Fatalf("a refused install left a command reporting %q", got)
	}

	variantArchive, variantSums := packageFixture(t, root, variant, "2.5.0-acme.1")
	install.run(nil, "2.5.0-acme.1", binDir, variantArchive, variantSums, true)
	if got := installedVersion(t, binDir, variant); got != "2.5.0-acme.1" {
		t.Fatalf("variant command reports %q", got)
	}
	if got := installedVersion(t, binDir, publicCommand); got != "2.5.0" {
		t.Fatalf("the variant install repointed the public command: %q", got)
	}
	// The installer resolves its directory, which a temporary one may not be.
	resolved, err := filepath.EvalSymlinks(binDir)
	if err != nil {
		t.Fatal(err)
	}
	versions := filepath.Join(resolved, ".context-circuit-versions")
	if runtime.GOOS == "windows" {
		for _, hash := range []string{"launcher.sha256", variant + ".launcher.sha256"} {
			if _, err := os.Stat(filepath.Join(versions, hash)); err != nil {
				t.Fatalf("%s: %v", hash, err)
			}
		}
	} else if link, err := os.Readlink(filepath.Join(binDir, variant)); err != nil || link != filepath.Join(versions, "2.5.0-acme.1-"+runtime.GOOS+"-"+runtime.GOARCH, publicCommand) {
		// The store binary keeps the public name whatever the command is called.
		t.Fatalf("variant link: %q %v", link, err)
	}

	// A public install afterwards accepts its own command and leaves the
	// variant's alone.
	laterArchive, laterSums := packageFixture(t, root, publicCommand, "2.5.1")
	install.run(nil, "2.5.1", binDir, laterArchive, laterSums, true)
	if got := installedVersion(t, binDir, publicCommand); got != "2.5.1" {
		t.Fatalf("public command reports %q", got)
	}
	if got := installedVersion(t, binDir, variant); got != "2.5.0-acme.1" {
		t.Fatalf("the public install repointed the variant command: %q", got)
	}

	// A reserved or uppercase label is no variant: it installs as the public
	// command from an archive under the public name. The uppercase one has its
	// own base, because a case-insensitive file system would read its version
	// folder as the acme.1 one.
	for _, version := range []string{"2.5.0-rc.1", "2.5.2-Acme.1"} {
		archive, sums := packageFixture(t, root, publicCommand, version)
		install.run(nil, version, binDir, archive, sums, true)
		if got := installedVersion(t, binDir, publicCommand); got != version {
			t.Fatalf("%s: public command reports %q", version, got)
		}
	}
	if got := installedVersion(t, binDir, "rc-context-circuit-cli"); got != "" {
		t.Fatalf("a reserved label got its own command, reporting %q", got)
	}
	// On a case-insensitive file system this is also where an Acme command
	// would be, so it still holding acme.1 proves there is none.
	if got := installedVersion(t, binDir, variant); got != "2.5.0-acme.1" {
		t.Fatalf("variant command reports %q", got)
	}
}

// CONTEXT_CIRCUIT_BIN_DIR is the default command directory, and an explicit
// flag still wins over it.
func TestBinDirFromTheEnvironment(t *testing.T) {
	root := t.TempDir()
	install := newInstaller(t, root)
	envDir := filepath.Join(root, "env bin")
	flagDir := filepath.Join(root, "flag bin")
	firstArchive, firstSums := packageFixture(t, root, publicCommand, "2.5.0")
	install.run([]string{"CONTEXT_CIRCUIT_BIN_DIR=" + envDir}, "2.5.0", "", firstArchive, firstSums, true)
	if got := installedVersion(t, envDir, publicCommand); got != "2.5.0" {
		t.Fatalf("CONTEXT_CIRCUIT_BIN_DIR without the flag: %q", got)
	}
	secondArchive, secondSums := packageFixture(t, root, publicCommand, "2.5.1")
	install.run([]string{"CONTEXT_CIRCUIT_BIN_DIR=" + envDir}, "2.5.1", flagDir, secondArchive, secondSums, true)
	if got := installedVersion(t, flagDir, publicCommand); got != "2.5.1" {
		t.Fatalf("the flag beside CONTEXT_CIRCUIT_BIN_DIR: %q", got)
	}
	if got := installedVersion(t, envDir, publicCommand); got != "2.5.0" {
		t.Fatalf("CONTEXT_CIRCUIT_BIN_DIR was used although the flag was given: %q", got)
	}
}

// packageFixture builds a fake release archive, named for prefix, holding a
// binary that reports version (or reported, when given).
func packageFixture(t *testing.T, root, prefix, version string, reported ...string) (string, string) {
	t.Helper()
	dir := filepath.Join(root, prefix+"-v"+version)
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(dir, "main.go")
	text := "package main\nimport \"fmt\"\nvar version string\nfunc main(){fmt.Println(version)}\n"
	if err := os.WriteFile(source, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	binaryName := "context-circuit-cli"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(dir, binaryName)
	runVersion := version
	if len(reported) > 0 {
		runVersion = reported[0]
	}
	command := exec.Command("go", "build", "-ldflags", "-X main.version="+runVersion, "-o", binaryPath, source)
	if data, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build fixture: %v %s", err, data)
	}
	binary, err := os.ReadFile(binaryPath)
	if err != nil {
		t.Fatal(err)
	}
	payload := map[string][]byte{binaryName: binary, "README.md": []byte("Installer fixture"), "LICENSE": []byte("Fixture license"), "THIRD_PARTY_NOTICES.txt": []byte("Fixture")}
	name := fmt.Sprintf("%s-v%s-%s-%s", prefix, version, runtime.GOOS, runtime.GOARCH)
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	archive := filepath.Join(dir, name+ext)
	file, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		writer := zip.NewWriter(file)
		for name, data := range payload {
			entry, err := writer.Create(name)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := entry.Write(data); err != nil {
				t.Fatal(err)
			}
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
	} else {
		compressed := gzip.NewWriter(file)
		writer := tar.NewWriter(compressed)
		for name, data := range payload {
			if err := writer.WriteHeader(&tar.Header{Name: name, Mode: 0755, Size: int64(len(data))}); err != nil {
				t.Fatal(err)
			}
			if _, err := writer.Write(data); err != nil {
				t.Fatal(err)
			}
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		if err := compressed.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(archive)
	if err != nil {
		t.Fatal(err)
	}
	sums := filepath.Join(dir, "SHA256SUMS")
	if err := os.WriteFile(sums, fmt.Appendf(nil, "%x  ./%s\n", sha256.Sum256(data), filepath.Base(archive)), 0600); err != nil {
		t.Fatal(err)
	}
	return archive, sums
}
