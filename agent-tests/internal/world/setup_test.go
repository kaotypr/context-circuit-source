package world

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
)

func TestSetupIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("short")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go missing")
	}
	root := repoRoot(t)
	before := porcelain(t, root)

	runRoot, err := NewRunRoot("")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(runRoot) })

	pin, err := Pin(root)
	if err != nil {
		t.Fatal(err)
	}
	cli, err := BuildCLI(context.Background(), root, runRoot, pin)
	if err != nil {
		t.Fatal(err)
	}
	w, err := Setup(context.Background(), SetupInput{
		RunRoot:     runRoot,
		SrcDir:      root,
		CLIPath:     cli,
		FixturesDir: filepath.Join(root, "agent-tests", "fixtures"),
		Pin:         pin,
		PATH:        os.Getenv("PATH"),
	})
	if err != nil {
		t.Fatal(err)
	}
	after := porcelain(t, root)
	if before != after {
		t.Fatalf("checkout status changed\nbefore:\n%s\nafter:\n%s", before, after)
	}
	log, err := os.ReadFile(w.DecoyLog)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if len(bytesTrim(log)) != 0 {
		t.Fatalf("decoy log %q", log)
	}
	if w.Env["CONTEXT_CIRCUIT_BIN_DIR"] != filepath.Join(runRoot, "bin") {
		t.Fatalf("env %+v", w.Env)
	}
	raw, err := os.ReadFile(filepath.Join(runRoot, "out", "check-before.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		OK     bool             `json:"ok"`
		Issues []map[string]any `json:"issues"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("check json: %v\n%s", err, raw)
	}
	if !doc.OK || len(doc.Issues) != 0 {
		t.Fatalf("check not ok: %s", raw)
	}
	if _, err := os.Stat(w.API); err != nil {
		t.Fatal(err)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	// agent-tests/internal/world -> repository root
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "cmd", "context-circuit", "main.go")); err != nil {
		t.Fatal(err)
	}
	return root
}

func porcelain(t *testing.T, root string) string {
	t.Helper()
	res, err := proc.Run(context.Background(), proc.Spec{
		Argv: []string{"git", "-C", root, "status", "--porcelain"},
		Env: []string{
			"PATH=" + os.Getenv("PATH"),
			"GIT_CONFIG_NOSYSTEM=1",
			"GIT_CONFIG_GLOBAL=/dev/null",
		},
	})
	if err != nil || res.ExitCode != 0 {
		t.Fatalf("git status: %v\n%s", err, res.Stderr)
	}
	return string(res.Stdout)
}

func bytesTrim(b []byte) []byte {
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == ' ') {
		b = b[:len(b)-1]
	}
	return b
}
