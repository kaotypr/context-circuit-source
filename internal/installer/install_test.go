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
	firstArchive, firstSums := packageFixture(t, root, "2.0.0-install1")
	secondArchive, secondSums := packageFixture(t, root, "2.0.0-install2")
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
	wrongArchive, wrongSums := packageFixture(t, root, "2.0.0-claimed", "2.0.0-actual")
	install("2.0.0-claimed", wrongArchive, wrongSums, false)
	checkVersion("2.0.0-install1")
}

func packageFixture(t *testing.T, root, version string, reported ...string) (string, string) {
	t.Helper()
	dir := filepath.Join(root, version)
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
	payload := map[string][]byte{binaryName: binary, "README.md": []byte("Installer fixture"), "THIRD_PARTY_NOTICES.txt": []byte("Fixture")}
	name := fmt.Sprintf("context-circuit-cli-v%s-%s-%s", version, runtime.GOOS, runtime.GOARCH)
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
	if err := os.WriteFile(sums, []byte(fmt.Sprintf("%x  ./%s\n", sha256.Sum256(data), filepath.Base(archive))), 0600); err != nil {
		t.Fatal(err)
	}
	return archive, sums
}
