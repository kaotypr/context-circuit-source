package scenario

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadErrors(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("bad-id.yaml", "id: Not_Valid\nversion: 1\nprompt: hi\n")
	if _, err := Load(filepath.Join(dir, "bad-id.yaml")); err == nil {
		t.Fatal("expected bad id")
	}
	write("empty.yaml", "id: empty-prompt\nversion: 1\nprompt: \"\"\n")
	if _, err := Load(filepath.Join(dir, "empty.yaml")); err == nil {
		t.Fatal("expected empty prompt")
	}
	write("unknown.yaml", "id: ok-id\nversion: 1\nprompt: hi\nextra: 1\n")
	if _, err := Load(filepath.Join(dir, "unknown.yaml")); err == nil {
		t.Fatal("expected unknown field")
	}
	if err := os.Remove(filepath.Join(dir, "bad-id.yaml")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, "empty.yaml")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, "unknown.yaml")); err != nil {
		t.Fatal(err)
	}
	write("a.yaml", "id: same\nversion: 1\nprompt: one\n")
	write("b.yaml", "id: same\nversion: 1\nprompt: two\n")
	if _, err := LoadDir(dir); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("duplicate: %v", err)
	}
}

func TestLoadDirMissing(t *testing.T) {
	if _, err := LoadDir(filepath.Join(t.TempDir(), "missing")); err == nil || !strings.Contains(err.Error(), "no scenarios found") {
		t.Fatal(err)
	}
	empty := t.TempDir()
	if _, err := LoadDir(empty); err == nil || !strings.Contains(err.Error(), "no scenarios found") {
		t.Fatal(err)
	}
}

func TestHashStable(t *testing.T) {
	s := Scenario{
		ID:      "orient",
		Version: 3,
		Title:   "Orientation",
		Prompt:  "What is $WS for?",
		Soft:    map[string]string{"b": "second", "a": "first"},
		Checks:  []string{"u1"},
	}
	again := s
	again.Soft = map[string]string{"a": "first", "b": "second"}
	first := s.Hash("sha256:abc")
	second := again.Hash("sha256:abc")
	if first != second {
		t.Fatalf("unstable\n%s\n%s", first, second)
	}
	if !strings.HasPrefix(first, "sha256:") || len(first) != len("sha256:")+64 {
		t.Fatalf("hash %s", first)
	}
	if s.Hash("sha256:def") == first {
		t.Fatal("fixture digest did not change the hash")
	}
	changed := s
	changed.Prompt = "other"
	if changed.Hash("sha256:abc") == first {
		t.Fatal("prompt did not change the hash")
	}
	changed = s
	changed.Version = 4
	if changed.Hash("sha256:abc") == first {
		t.Fatal("version did not change the hash")
	}
}

func TestVerdict(t *testing.T) {
	pass := []CheckResult{{Name: "u1", Pass: true}}
	fail := []CheckResult{{Name: "u1", Pass: true}, {Name: "u2", Pass: false, Detail: "pushed"}}
	cases := []struct {
		name       string
		results    []CheckResult
		harness    error
		incomplete bool
		want       string
	}{
		{"pass", pass, nil, false, "pass"},
		{"empty pass", nil, nil, false, "pass"},
		{"fail", fail, nil, false, "fail"},
		{"fail beats incomplete", fail, nil, true, "fail"},
		{"incomplete", pass, nil, true, "inconclusive"},
		{"harness beats fail", fail, errSentinel{}, false, "harness-error"},
		{"harness beats incomplete", pass, errSentinel{}, true, "harness-error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var harness error
			if tc.harness != nil {
				harness = tc.harness
			}
			if got := Verdict(tc.results, harness, tc.incomplete); got != tc.want {
				t.Fatalf("got %s want %s", got, tc.want)
			}
		})
	}
}

type errSentinel struct{}

func (errSentinel) Error() string { return "boom" }
