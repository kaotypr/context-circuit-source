package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFinalObjectReadsTheLastJSONObject(t *testing.T) {
	for _, tc := range []struct {
		out     string
		status  string
		summary string
		found   bool
	}{
		{`{"status":"ok","summary":"14 behaviors match"}`, "ok", "14 behaviors match", true},
		{"checking…\nall good\n{\"status\":\"skipped\",\"summary\":\"no parity list\"}\n", "skipped", "no parity list", true},
		{"{\"status\": \"ok\",\n \"summary\": \"pretty\",\n \"nested\": {\"a\": 1}\n}\n", "ok", "pretty", true},
		{"progress {not json}\n", "", "", false},
		{"no object at all", "", "", false},
		{"", "", "", false},
	} {
		object, found := finalObject([]byte(tc.out))
		if found != tc.found || object.Status != tc.status || object.Summary != tc.summary {
			t.Errorf("%q: got %+v %v", tc.out, object, found)
		}
	}
}

func TestCappedKeepsTheFirstBytesAndSaysItCut(t *testing.T) {
	c := &capped{limit: 4}
	for _, chunk := range []string{"ab", "cdef", "gh"} {
		if n, err := c.Write([]byte(chunk)); n != len(chunk) || err != nil {
			t.Fatalf("write %q: %d %v", chunk, n, err)
		}
	}
	if c.buf.String() != "abcd" || !c.truncated {
		t.Fatalf("kept %q truncated=%v", c.buf.String(), c.truncated)
	}
}

// A hook with no Windows command is skipped on Windows, and says so.
func TestScriptCommandChoosesThePlatformsCommand(t *testing.T) {
	script := HookScript{Run: []string{"sh", "scripts/a.sh"}, OnFailure: "stop"}
	if command, why := scriptCommand(script, "linux"); why != "" || command[0] != "sh" {
		t.Fatalf("linux: %v %q", command, why)
	}
	if _, why := scriptCommand(script, "windows"); !strings.Contains(why, "no Windows command") {
		t.Fatalf("windows without run_windows: %q", why)
	}
	script.RunWindows = []string{"powershell", "-File", `scripts\a.ps1`}
	if command, why := scriptCommand(script, "windows"); why != "" || command[0] != "powershell" {
		t.Fatalf("windows: %v %q", command, why)
	}
}

func TestScriptArgumentsResolveInsideTheFolder(t *testing.T) {
	folder := t.TempDir()
	if err := os.MkdirAll(filepath.Join(folder, "scripts"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, "scripts", "a.sh"), []byte("exit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	args, err := resolveScriptArgs(folder, []string{"sh", "scripts/a.sh", "--flag", "plain", "not/there", "https://example.invalid/x"})
	if err != nil {
		t.Fatal(err)
	}
	real, _ := filepath.EvalSymlinks(filepath.Join(folder, "scripts", "a.sh"))
	if args[0] != "sh" || args[1] != real || args[2] != "--flag" || args[3] != "plain" || args[4] != "not/there" || args[5] != "https://example.invalid/x" {
		t.Fatalf("resolved: %v", args)
	}
	if _, err := resolveScriptArgs(folder, []string{"sh", "../outside.sh"}); err == nil {
		t.Fatal("a path leaving the folder was accepted")
	}
	if _, err := resolveScriptArgs(folder, []string{"scripts/missing"}); err == nil {
		t.Fatal("a missing program in the folder was accepted")
	}
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "x.sh"), []byte("exit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "x.sh"), filepath.Join(folder, "scripts", "link.sh")); err == nil {
		if _, err := resolveScriptArgs(folder, []string{"sh", "scripts/link.sh"}); err == nil {
			t.Fatal("a symlink out of the folder was accepted")
		}
	}
}

// A folder's version is the same whether its files end lines in LF or CRLF,
// and changes when any file, its name included, changes.
func TestFolderHashNormalizesLineEndings(t *testing.T) {
	writeTree := func(files map[string]string) string {
		dir := t.TempDir()
		for name, text := range files {
			path := filepath.Join(dir, filepath.FromSlash(name))
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(text), 0644); err != nil {
				t.Fatal(err)
			}
		}
		return dir
	}
	hash := func(dir string) string {
		h, err := FolderHash(dir)
		if err != nil {
			t.Fatal(err)
		}
		return h
	}
	lf := hash(writeTree(map[string]string{"extension.yaml": "id: a\nmembers: opt-in\n", "scripts/run.sh": "echo hi\n"}))
	crlf := hash(writeTree(map[string]string{"extension.yaml": "id: a\r\nmembers: opt-in\r\n", "scripts/run.sh": "echo hi\r\n"}))
	if lf != crlf || !strings.HasPrefix(lf, "sha256:") {
		t.Fatalf("line endings changed the version: %s %s", lf, crlf)
	}
	edited := hash(writeTree(map[string]string{"extension.yaml": "id: a\nmembers: opt-in\n", "scripts/run.sh": "echo bye\n"}))
	renamed := hash(writeTree(map[string]string{"extension.yaml": "id: a\nmembers: opt-in\n", "scripts/go.sh": "echo hi\n"}))
	added := hash(writeTree(map[string]string{"extension.yaml": "id: a\nmembers: opt-in\n", "scripts/run.sh": "echo hi\n", "README.md": ""}))
	if edited == lf || renamed == lf || added == lf {
		t.Fatal("a changed folder kept its version")
	}
}
