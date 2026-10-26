package world

import (
	"archive/tar"
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
)

func TestExtractArchiveLeavesCheckoutUntouched(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git missing")
	}
	repo := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		res, err := proc.Run(context.Background(), proc.Spec{
			Argv: append([]string{"git", "-C", repo}, args...),
			Env: []string{
				"PATH=" + os.Getenv("PATH"),
				"GIT_CONFIG_NOSYSTEM=1",
				"GIT_CONFIG_GLOBAL=/dev/null",
				"GIT_AUTHOR_DATE=2026-01-15T09:00:00Z",
				"GIT_COMMITTER_DATE=2026-01-15T09:00:00Z",
			},
		})
		if err != nil || res.ExitCode != 0 {
			t.Fatalf("git %v: %v exit %d\n%s", args, err, res.ExitCode, res.Stderr)
		}
	}
	git("init", "-q", "-b", "main")
	git("-c", "user.name=Agent Test", "-c", "user.email=agent-test@example.invalid", "commit", "--allow-empty", "-qm", "init")
	if err := os.WriteFile(filepath.Join(repo, "hello.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git("add", "hello.txt")
	git("-c", "user.name=Agent Test", "-c", "user.email=agent-test@example.invalid", "commit", "-qm", "add hello")

	beforeList := listFiles(t, repo)
	beforeStatus := gitText(t, repo, "status")
	beforePorcelain := gitText(t, repo, "status", "--porcelain")

	dest := filepath.Join(t.TempDir(), "src")
	if err := ExtractArchive(context.Background(), repo, "HEAD", dest); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dest, "hello.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello\n" {
		t.Fatalf("extracted %q", got)
	}
	if afterList := listFiles(t, repo); afterList != beforeList {
		t.Fatalf("checkout listing changed\nonly before: %s\nonly after: %s", listDiff(beforeList, afterList), listDiff(afterList, beforeList))
	}
	if gitText(t, repo, "status") != beforeStatus || gitText(t, repo, "status", "--porcelain") != beforePorcelain {
		t.Fatal("git status changed")
	}
}

func TestExtractRejectsEscape(t *testing.T) {
	dest := t.TempDir()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{Name: "../outside.txt", Mode: 0o644, Size: 4, Typeflag: tar.TypeReg}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte("evil")); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := extractTar(bytes.NewReader(buf.Bytes()), dest); err == nil {
		t.Fatal("expected refusal")
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(dest), "outside.txt")); err == nil {
		t.Fatal("wrote outside dest")
	}
}

func listDiff(a, b string) string {
	have := map[string]bool{}
	for _, line := range strings.Split(b, "\n") {
		have[line] = true
	}
	var only []string
	for _, line := range strings.Split(a, "\n") {
		if !have[line] {
			only = append(only, line)
		}
	}
	return strings.Join(only, ", ")
}

func listFiles(t *testing.T, root string) string {
	t.Helper()
	var names []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		// Git creates and removes lock files on its own. A maintenance
		// lock that disappears during archive is not a checkout write.
		if strings.HasPrefix(rel, ".git"+string(os.PathSeparator)) && strings.HasSuffix(rel, ".lock") {
			return nil
		}
		names = append(names, rel)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(names)
	return strings.Join(names, "\n")
}

func gitText(t *testing.T, repo string, args ...string) string {
	t.Helper()
	res, err := proc.Run(context.Background(), proc.Spec{
		Argv: append([]string{"git", "-C", repo}, args...),
		Env: []string{
			"PATH=" + os.Getenv("PATH"),
			"GIT_CONFIG_NOSYSTEM=1",
			"GIT_CONFIG_GLOBAL=/dev/null",
		},
	})
	if err != nil || res.ExitCode != 0 {
		t.Fatalf("git %v: %v\n%s", args, err, res.Stderr)
	}
	return string(res.Stdout)
}
