package world

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
)

func TestBuildCLITinyTree(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go missing")
	}
	src := t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, "cmd", "context-circuit"), 0o755); err != nil {
		t.Fatal(err)
	}
	gomod := "module example.com/tiny\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(src, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	mainSrc := "package main\n\nimport \"fmt\"\n\nvar version = \"dev\"\n\nfunc main() {\n\tfmt.Println(version)\n}\n"
	if err := os.WriteFile(filepath.Join(src, "cmd", "context-circuit", "main.go"), []byte(mainSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	runRoot := t.TempDir()
	const pin = "9.9.9-test"
	cli, err := BuildCLI(context.Background(), src, runRoot, pin)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(runRoot, "bin", ".context-circuit-versions", pin+"-"+runtime.GOOS+"-"+runtime.GOARCH, "context-circuit-cli")
	if cli != want {
		t.Fatalf("cli %s want %s", cli, want)
	}
	info, err := os.Stat(cli)
	if err != nil {
		t.Fatal(err)
	}
	if info.IsDir() || info.Mode()&0o111 == 0 {
		t.Fatalf("binary mode %o", info.Mode())
	}
	res, err := proc.Run(context.Background(), proc.Spec{
		Argv: []string{cli},
		Env:  []string{"PATH=" + os.Getenv("PATH")},
	})
	if err != nil || res.ExitCode != 0 {
		t.Fatalf("run: %v exit %d\n%s", err, res.ExitCode, res.Stderr)
	}
	if strings.TrimSpace(string(res.Stdout)) != pin {
		t.Fatalf("version %q", res.Stdout)
	}
}
