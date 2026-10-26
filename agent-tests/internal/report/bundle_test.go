package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBundleCleansPathsAndDropsAuth(t *testing.T) {
	dir := t.TempDir()
	run := filepath.Join(dir, "run")
	root := filepath.Join(run, "runs", "intent", "claude", "branch", "r1")
	files := filepath.Join(root, "artifacts", "files")
	authDir := filepath.Join(root, "claude")
	if err := os.MkdirAll(files, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(authDir, 0o755); err != nil {
		t.Fatal(err)
	}
	home := "/tmp/zzhome/project"
	user := "zzuniqueuser"
	runRoot := "/opt/runroot"
	body := "home " + home + "/a user /opt/" + user + "/bin run " + runRoot + "/ws auto /tmp/agent-test-zz/ws/file\n"
	mustWrite(t, filepath.Join(run, "runs.jsonl"), body)
	mustWrite(t, filepath.Join(root, "stream.jsonl"), body)
	mustWrite(t, filepath.Join(root, "reply.md"), "hello /opt/"+user+"/x\n")
	mustWrite(t, filepath.Join(files, "note.md"), "see "+home+"\n")
	mustWrite(t, filepath.Join(files, "auth.json"), "sk-abcdefghijklmnop\n")
	mustWrite(t, filepath.Join(files, ".credentials.json"), "cred\n")
	mustWrite(t, filepath.Join(files, "id_rsa"), "key\n")
	mustWrite(t, filepath.Join(files, "session.token"), "tok\n")
	mustWrite(t, filepath.Join(files, "key.pem"), "pem\n")
	mustWrite(t, filepath.Join(files, ".netrc"), "netrc\n")
	mustWrite(t, filepath.Join(files, ".env"), "env\n")
	mustWrite(t, filepath.Join(authDir, "session.txt"), "AUTHDIR-MARKER\n")
	outside := filepath.Join(dir, "outside.txt")
	mustWrite(t, outside, "OUTSIDE-MARKER\n")
	if err := os.Symlink(outside, filepath.Join(files, "leak")); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(run, "bundle.zip")
	if err := os.WriteFile(dest, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteBundle(BundleOptions{RunDir: run, Dest: dest, Home: home, User: user, RunRoots: []string{runRoot}}); err != nil {
		t.Fatal(err)
	}
	got, err := readZip(dest)
	if err != nil {
		t.Fatal(err)
	}
	var all strings.Builder
	for name, text := range got {
		if !safeName(name) {
			t.Fatalf("unsafe entry %q", name)
		}
		if strings.Contains(name, "auth.json") || strings.Contains(name, "credentials") || strings.Contains(name, "id_rsa") || strings.HasSuffix(name, ".token") || strings.HasSuffix(name, ".pem") || strings.HasSuffix(name, ".netrc") || strings.HasSuffix(name, ".env") {
			t.Fatalf("auth entry %s", name)
		}
		all.WriteString(name)
		all.WriteString("\n")
		all.WriteString(text)
	}
	joined := all.String()
	for _, banned := range []string{home, user, runRoot, "/tmp/agent-test-zz", "AUTHDIR-MARKER", "OUTSIDE-MARKER", "sk-abcdefghijklmnop"} {
		if strings.Contains(joined, banned) {
			t.Fatalf("zip still has %q\n%s", banned, joined)
		}
	}
	if !strings.Contains(joined, "$HOME") || !strings.Contains(joined, "$USER") || !strings.Contains(joined, "$RUN") {
		t.Fatalf("missing replacements\n%s", joined)
	}
	if _, ok := got["runs/intent/claude/branch/r1/stream.jsonl"]; !ok {
		t.Fatalf("stream missing: %v", keys(got))
	}
}

func TestBundleRefusesSecret(t *testing.T) {
	sk := "sk-abcdefghijklmnop"
	ghp := "ghp_" + strings.Repeat("a", 20)
	bearer := "Bearer " + strings.Repeat("c", 20)
	cases := []struct {
		name string
		rel  string
		body string
	}{
		{name: "stream", rel: "runs/intent/claude/branch/r1/stream.jsonl", body: "token " + sk + "\n"},
		{name: "reply", rel: "runs/intent/claude/branch/r1/reply.md", body: ghp + "\n"},
		{name: "file", rel: "runs/intent/claude/branch/r1/artifacts/files/note.md", body: "Authorization: " + bearer + "\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			mustWrite(t, filepath.Join(dir, tc.rel), tc.body)
			dest := filepath.Join(dir, "bundle.zip")
			if err := os.WriteFile(dest, []byte("stale"), 0o644); err != nil {
				t.Fatal(err)
			}
			err := WriteBundle(BundleOptions{RunDir: dir, Dest: dest})
			if err == nil || !strings.Contains(err.Error(), "key-shaped string") {
				t.Fatalf("err %v", err)
			}
			base := filepath.Base(tc.rel)
			if !strings.Contains(err.Error(), base) {
				t.Fatalf("error does not name the file: %v", err)
			}
			for _, secret := range []string{sk, ghp, bearer} {
				if strings.Contains(err.Error(), secret) {
					t.Fatalf("error contains the secret: %v", err)
				}
			}
			if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
				t.Fatal("archive remains")
			}
			matches, _ := filepath.Glob(filepath.Join(dir, ".bundle-*.zip.tmp"))
			if len(matches) != 0 {
				t.Fatalf("partial zip %v", matches)
			}
		})
	}
}

func TestSafeZipName(t *testing.T) {
	for _, name := range []string{"../outside.txt", "/etc/passwd", "foo/../../etc/passwd", ".."} {
		if _, ok := safeZipName(name); ok {
			t.Fatalf("accepted %q", name)
		}
	}
	got, ok := safeZipName("runs/intent/claude/branch/r1/reply.md")
	if !ok || got != "runs/intent/claude/branch/r1/reply.md" {
		t.Fatalf("%q %v", got, ok)
	}
}

func safeName(name string) bool {
	_, ok := safeZipName(name)
	return ok
}

func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for key := range m {
		out = append(out, key)
	}
	return out
}
