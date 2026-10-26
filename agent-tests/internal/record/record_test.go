package record

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "runs.jsonl")
	n := int64(3)
	if err := Append(path, Run{Scenario: "one", Tokens: Tokens{InputTotal: &n}}); err != nil {
		t.Fatal(err)
	}
	if err := Append(path, Run{Scenario: "two"}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("lines %d: %s", len(lines), b)
	}
	var first Run
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatal(err)
	}
	if first.Scenario != "one" || first.Tokens.InputTotal == nil || *first.Tokens.InputTotal != 3 {
		t.Fatalf("%+v", first)
	}
	var second Run
	if err := json.Unmarshal([]byte(lines[1]), &second); err != nil {
		t.Fatal(err)
	}
	if second.Scenario != "two" {
		t.Fatalf("%+v", second)
	}
}

func TestNilTokensMarshalNull(t *testing.T) {
	b, err := json.Marshal(Run{})
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}
	var tok map[string]json.RawMessage
	if err := json.Unmarshal(raw["tokens"], &tok); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"input_total", "input_uncached", "cache_write", "cache_read", "output", "reasoning"} {
		if string(tok[key]) != "null" {
			t.Fatalf("%s marshaled as %s", key, tok[key])
		}
	}
	if string(raw["cost_usd"]) != "null" {
		t.Fatalf("cost_usd %s", raw["cost_usd"])
	}
	if _, ok := raw["null_reasons"]; ok {
		t.Fatalf("null_reasons present: %s", raw["null_reasons"])
	}
}

func TestDigestStableAndSensitive(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sub", "a.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitFile := filepath.Join(root, ".git", "HEAD")
	if err := os.WriteFile(gitFile, []byte("ref: refs/heads/main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	first, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("unstable %s %s", first, second)
	}
	if !strings.HasPrefix(first, "sha256:") || len(first) != len("sha256:")+64 {
		t.Fatalf("digest %s", first)
	}
	if err := os.WriteFile(gitFile, []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	afterGit, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if afterGit != first {
		t.Fatal(".git changed the digest")
	}
	if err := os.WriteFile(filepath.Join(root, "sub", "a.txt"), []byte("hellp"), 0o644); err != nil {
		t.Fatal(err)
	}
	changed, err := Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	if changed == first {
		t.Fatal("byte change did not change the digest")
	}
}

func TestRedact(t *testing.T) {
	longSK := "sk-abcdefghijklmnop"
	longGHP := "ghp_" + strings.Repeat("a", 20)
	longPat := "github_pat_" + strings.Repeat("b", 20)
	longAKIA := "AKIA" + strings.Repeat("A", 16)
	longBearer := "Bearer " + strings.Repeat("c", 20)
	tests := []struct {
		name string
		in   string
		home string
		user string
		want string
	}{
		{name: "home", in: "file /tmp/home/a", home: "/tmp/home", want: "file $HOME/a"},
		{name: "home boundary", in: "file /tmp/homework", home: "/tmp/home", want: "file /tmp/homework"},
		{name: "user element", in: "/opt/ada/bin", user: "ada", want: "/opt/$USER/bin"},
		{name: "user substring", in: "/opt/adapter/bin", user: "ada", want: "/opt/adapter/bin"},
		{name: "home and user", in: "/tmp/home/ada/x", home: "/tmp/home", user: "ada", want: "$HOME/$USER/x"},
		{name: "sk", in: "key " + longSK + " xx", want: "key [redacted] xx"},
		{name: "short sk", in: "sk-short", want: "sk-short"},
		{name: "ghp", in: longGHP, want: "[redacted]"},
		{name: "github pat", in: "token " + longPat, want: "token [redacted]"},
		{name: "akia", in: longAKIA, want: "[redacted]"},
		{name: "bearer", in: "Authorization: " + longBearer, want: "Authorization: Bearer [redacted]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Redact(tt.in, tt.home, tt.user); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestRedactPathsLeavesSecrets(t *testing.T) {
	secret := "sk-abcdefghijklmnop"
	in := "/tmp/home/ada/" + secret
	got := RedactPaths(in, "/tmp/home", "ada")
	if got != "$HOME/$USER/"+secret {
		t.Fatalf("%q", got)
	}
	if !LooksLikeSecret(got) {
		t.Fatal("path cleaning removed the key shape")
	}
}

func TestRedactRunRoot(t *testing.T) {
	if got := RedactRunRoot("log /tmp/run/abc/out", "/tmp/run/abc"); got != "log $RUN/out" {
		t.Fatalf("%q", got)
	}
	if got := RedactRunRoot("log /tmp/runaway", "/tmp/run"); got != "log /tmp/runaway" {
		t.Fatalf("%q", got)
	}
}

func TestLooksLikeSecret(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{in: "sk-abcdefghijklmnop", want: true},
		{in: "sk-short", want: false},
		{in: "ghp_" + strings.Repeat("a", 20), want: true},
		{in: "ghp_short", want: false},
		{in: "github_pat_" + strings.Repeat("b", 20), want: true},
		{in: "AKIA" + strings.Repeat("A", 16), want: true},
		{in: "AKIA_short", want: false},
		{in: "Bearer " + strings.Repeat("c", 20), want: true},
		{in: "Bearer short", want: false},
		{in: "hello", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := LooksLikeSecret(tt.in); got != tt.want {
				t.Fatalf("got %v", got)
			}
		})
	}
}
