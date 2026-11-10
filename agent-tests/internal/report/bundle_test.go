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

func TestBundleUserWordAndDroppedAuthKey(t *testing.T) {
	dir := t.TempDir()
	user := "zzuser"
	note := filepath.Join(dir, "artifacts", "files", user+"-note.md")
	mustWrite(t, note, "zzuserish saw "+user+"\n")
	mustWrite(t, filepath.Join(dir, "artifacts", "files", "auth.json"), "sk-abcdefghijklmnop\n")
	dest := filepath.Join(dir, "bundle.zip")
	if err := WriteBundle(BundleOptions{RunDir: dir, Dest: dest, User: user}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dest); err != nil {
		t.Fatal(err)
	}
	got, err := readZip(dest)
	if err != nil {
		t.Fatal(err)
	}
	body, ok := got["artifacts/files/$USER-note.md"]
	if !ok {
		t.Fatalf("entry names %v", keys(got))
	}
	if body != "zzuserish saw $USER\n" {
		t.Fatalf("%q", body)
	}
	for name, text := range got {
		if strings.Contains(name, "auth.json") || strings.Contains(text, "sk-abcdefghijklmnop") {
			t.Fatalf("auth file kept %s", name)
		}
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
		{name: "jwt", rel: "runs/intent/codex/branch/r1/artifacts/files/note.md", body: "eyJhbGciOiJub25lIn0.eyJzdWIiOiJ0ZXN0In0.signatureXXXX\n"},
		{name: "refresh", rel: "runs/intent/codex/branch/r1/artifacts/files/note.md", body: "rt_" + strings.Repeat("q", 24) + "\n"},
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
			for _, secret := range []string{sk, ghp, bearer, strings.TrimSpace(tc.body)} {
				if secret != "" && strings.Contains(err.Error(), secret) {
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

func TestBundleDropsCodexDirectory(t *testing.T) {
	dir := t.TempDir()
	token := "rt_" + strings.Repeat("b", 20)
	mustWrite(t, filepath.Join(dir, "codex", "session.txt"), token+"\n")
	mustWrite(t, filepath.Join(dir, "codex", "auth.json"), token+"\n")
	mustWrite(t, filepath.Join(dir, "runs.jsonl"), "ok\n")
	dest := filepath.Join(dir, "bundle.zip")
	if err := WriteBundle(BundleOptions{RunDir: dir, Dest: dest}); err != nil {
		t.Fatal(err)
	}
	got, err := readZip(dest)
	if err != nil {
		t.Fatal(err)
	}
	for name, text := range got {
		if strings.Contains(name, "codex/") || strings.Contains(name, "auth.json") || strings.Contains(text, token) {
			t.Fatalf("kept %s %q", name, text)
		}
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

// A run root under a path with a space is found whole, and private words and
// paths are replaced; an archive that would still hold one is refused.
func TestBundlePrivateText(t *testing.T) {
	dir := t.TempDir()
	reply := filepath.Join(dir, "runs", "s", "claude", "branch", "r1", "reply.md")
	mustWrite(t, reply, "Worktree at /Volumes/My Disk/tmp/agent-test-77/ws.\nUpdate with tool@acme-internal.\n")
	dest := filepath.Join(dir, "bundle.zip")
	private := []string{"acme-internal", "/Volumes/My Disk/tmp"}
	if err := WriteBundle(BundleOptions{RunDir: dir, Dest: dest, Private: private}); err != nil {
		t.Fatal(err)
	}
	got, err := readZip(dest)
	if err != nil {
		t.Fatal(err)
	}
	for name, text := range got {
		if strings.Contains(text, "My Disk") || strings.Contains(text, "acme") {
			t.Fatalf("%s: %s", name, text)
		}
	}
	if cleaned := cleanPaths("at /Volumes/My Disk/tmp/agent-test-77/ws", "", "", nil, nil); cleaned != "at $RUN/ws" {
		t.Fatalf("run root with a space: %q", cleaned)
	}
}

func TestBundleRefusesPrivateLeft(t *testing.T) {
	dir := t.TempDir()
	reply := filepath.Join(dir, "runs", "s", "claude", "branch", "r1", "reply.md")
	// Cleaning replaces a private path only where it starts a word; glued
	// to the word before it, the path is left as written, so the archive is
	// refused rather than published with it.
	mustWrite(t, reply, "wrote x/Volumes/My Disk/tmp/y\n")
	dest := filepath.Join(dir, "bundle.zip")
	mustWrite(t, dest, "stale")
	err := WriteBundle(BundleOptions{RunDir: dir, Dest: dest, Private: []string{"/Volumes/My Disk/tmp"}})
	if err == nil || !strings.Contains(err.Error(), "private text left in") {
		t.Fatalf("err %v", err)
	}
	if strings.Contains(err.Error(), "My Disk") {
		t.Fatalf("the error repeats the private path: %v", err)
	}
	if _, serr := os.Stat(dest); serr == nil {
		t.Fatal("an archive was left behind")
	}
}

// A path printed right after a JSON \n escape is found from its slash, so a
// volume name with a space is replaced whole.
func TestRunRootAfterEscape(t *testing.T) {
	in := `"text":"main\n/Volumes/My Disk/tmp/agent-test-9/ws d82563a\n"`
	if got := cleanPaths(in, "", "", nil, nil); got != `"text":"main\n$RUN/ws d82563a\n"` {
		t.Fatalf("got %s", got)
	}
}
