package hosts

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodexSubscriptionBuild(t *testing.T) {
	h, err := ForName("codex")
	if err != nil {
		t.Fatal(err)
	}
	const secret = "codex-secret-value"
	run := filepath.Join(t.TempDir(), "run")
	in := LaunchInput{
		Prompt:           "Reply OK",
		Model:            "gpt-6-sol",
		Workspace:        filepath.Join(run, "ws"),
		RunRoot:          run,
		DecoyDir:         filepath.Join(run, "decoy"),
		SubscriptionAuth: true,
		Cfg: HostConfig{
			Command: "codex",
			Model:   "gpt-6-sol",
			Auth:    "env",
			AuthEnv: []string{"OPENAI_API_KEY"},
		},
		Caller: map[string]string{
			"PATH":           "/usr/bin:/bin",
			"HOME":           filepath.Join(run, "caller-home"),
			"CODEX_HOME":     filepath.Join(run, "caller-codex"),
			"OPENAI_API_KEY": secret,
		},
	}
	launch, err := h.Build(in)
	if err != nil {
		t.Fatal(err)
	}
	if launch.SetupArgv != nil || launch.SetupStdin != nil {
		t.Fatalf("setup %#v %#v", launch.SetupArgv, launch.SetupStdin)
	}
	if len(launch.CopyIn) != 1 {
		t.Fatalf("copies %#v", launch.CopyIn)
	}
	copy := launch.CopyIn[0]
	if copy.Src != filepath.Join(run, "caller-codex", "auth.json") || copy.Dst != filepath.Join(run, "codex", "auth.json") || copy.Mode != 0o600 {
		t.Fatalf("copy %+v", copy)
	}
	for _, e := range launch.Env {
		if strings.Contains(e, secret) || strings.Contains(e, "OPENAI_API_KEY") || strings.Contains(e, "caller-home") || strings.Contains(e, "caller-codex") {
			t.Fatalf("env leaked: %s", e)
		}
	}
	env := envMap(t, launch.Env)
	if env["HOME"] != filepath.Join(run, "codex-home") || env["CODEX_HOME"] != filepath.Join(run, "codex") {
		t.Fatalf("env %v", env)
	}

	in.Caller = map[string]string{"HOME": filepath.Join(run, "caller-home"), "PATH": "/usr/bin:/bin"}
	homeLaunch, err := h.Build(in)
	if err != nil {
		t.Fatal(err)
	}
	if homeLaunch.CopyIn[0].Src != filepath.Join(run, "caller-home", ".codex", "auth.json") {
		t.Fatalf("src %s", homeLaunch.CopyIn[0].Src)
	}
}

func TestCodexSubscriptionAuthReady(t *testing.T) {
	h, err := ForName("codex")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	cfg := HostConfig{AuthEnv: []string{"OPENAI_API_KEY"}, SubscriptionAuth: true}
	missing := filepath.Join(dir, "missing-login-src")
	caller := map[string]string{"CODEX_HOME": missing, "OPENAI_API_KEY": "not-a-login"}
	ok, reason := h.AuthReady(caller, cfg)
	if ok || reason != "no Codex login file found" {
		t.Fatalf("missing ok %v reason %q", ok, reason)
	}
	if strings.Contains(reason, "not-a-login") || strings.Contains(reason, "missing-login-src") {
		t.Fatalf("reason leaked %q", reason)
	}

	codexHome := filepath.Join(dir, "codex")
	if err := os.MkdirAll(codexHome, 0o700); err != nil {
		t.Fatal(err)
	}
	empty := filepath.Join(codexHome, "auth.json")
	if err := os.WriteFile(empty, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	ok, reason = h.AuthReady(map[string]string{"CODEX_HOME": codexHome}, cfg)
	if ok || reason != "no Codex login file found" {
		t.Fatalf("empty ok %v reason %q", ok, reason)
	}

	if err := os.WriteFile(empty, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	ok, reason = h.AuthReady(map[string]string{"CODEX_HOME": codexHome}, cfg)
	if !ok || reason != "" {
		t.Fatalf("present ok %v reason %q", ok, reason)
	}

	linkHome := filepath.Join(dir, "linkhome")
	if err := os.MkdirAll(linkHome, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "no-such-target"), filepath.Join(linkHome, "auth.json")); err != nil {
		t.Fatal(err)
	}
	ok, reason = h.AuthReady(map[string]string{"CODEX_HOME": linkHome}, cfg)
	if ok || reason != "the Codex login file is a symlink" {
		t.Fatalf("dangling ok %v reason %q", ok, reason)
	}

	// A symlink to a real login is refused too, as the copy refuses it.
	goodLink := filepath.Join(dir, "goodlink")
	if err := os.MkdirAll(goodLink, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(empty, filepath.Join(goodLink, "auth.json")); err != nil {
		t.Fatal(err)
	}
	ok, reason = h.AuthReady(map[string]string{"CODEX_HOME": goodLink}, cfg)
	if ok || reason != "the Codex login file is a symlink" {
		t.Fatalf("symlink ok %v reason %q", ok, reason)
	}

	// With neither CODEX_HOME nor HOME there is no login, not a file
	// relative to the working directory.
	ok, reason = h.AuthReady(map[string]string{}, cfg)
	if ok || reason != "no Codex login file found" {
		t.Fatalf("no home ok %v reason %q", ok, reason)
	}
}

func TestLoginSecretsRereadsTheFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	if got := LoginSecrets(path); len(got) != 0 {
		t.Fatalf("missing file gave %q", got)
	}
	refreshed := "rt_" + strings.Repeat("r", 30)
	if err := os.WriteFile(path, []byte(`{"tokens":{"refresh_token":"`+refreshed+`"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := LoginSecrets(path); !containsString(got, refreshed) {
		t.Fatalf("refreshed token not found: %q", got)
	}
}

func TestApplyCopies(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "source")
	if err := os.MkdirAll(srcDir, 0o700); err != nil {
		t.Fatal(err)
	}
	const marker = "MARKER-subscription-login-9f3c-zzzz"
	jwt := "eyJhbGciOiJub25lIn0.eyJzdWIiOiJ0ZXN0In0.signatureXXXX"
	refresh := "rt_" + strings.Repeat("q", 24)
	body := []byte(`{"access_token":"` + jwt + `","refresh_token":"` + refresh + `","note":"` + marker + `"}` + "\n")
	src := filepath.Join(srcDir, "auth.json")
	if err := os.WriteFile(src, body, 0o644); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "run", "codex", "auth.json")
	secrets, err := ApplyCopies([]FileCopy{{Src: src, Dst: dst, Mode: 0o600}})
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, body) {
		t.Fatalf("copy bytes changed")
	}
	info, err := os.Lstat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode %o", info.Mode().Perm())
	}
	srcInfo, err := os.Lstat(src)
	if err != nil {
		t.Fatal(err)
	}
	if srcInfo.Mode().Perm() != 0o644 {
		t.Fatalf("source mode %o", srcInfo.Mode().Perm())
	}
	again, err := os.ReadFile(src)
	if err != nil || !bytes.Equal(again, body) {
		t.Fatal("source bytes changed")
	}
	for _, want := range []string{marker, jwt, refresh} {
		if !containsString(secrets, want) {
			t.Fatalf("secrets %q missing %s", secrets, want)
		}
	}

	link := filepath.Join(dir, "link.json")
	if err := os.Symlink(src, link); err != nil {
		t.Fatal(err)
	}
	_, err = ApplyCopies([]FileCopy{{Src: link, Dst: filepath.Join(dir, "other", "auth.json"), Mode: 0o600}})
	if err == nil || !strings.Contains(err.Error(), "subscription login file") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("err %v", err)
	}
	if strings.Contains(err.Error(), marker) || strings.Contains(err.Error(), jwt) {
		t.Fatalf("symlink error leaked contents: %v", err)
	}
	if _, statErr := os.Lstat(filepath.Join(dir, "other", "auth.json")); !os.IsNotExist(statErr) {
		t.Fatal("symlink copy created a destination")
	}

	_, err = ApplyCopies([]FileCopy{{Src: filepath.Join(dir, "missing-login-src", "auth.json"), Dst: filepath.Join(dir, "nope"), Mode: 0o600}})
	if err == nil || !strings.Contains(err.Error(), "subscription login file") || strings.Contains(err.Error(), "missing-login-src") {
		t.Fatalf("err %v", err)
	}
}

func TestSecretsFromLoginParseError(t *testing.T) {
	jwt := "eyJhbGciOiJub25lIn0.eyJzdWIiOiJ0ZXN0In0.signatureXXXX"
	refresh := "rt_" + strings.Repeat("q", 24)
	body := "not-json " + jwt + " " + refresh + " ORDINARYWORDLONGXX"
	got := secretsFromLogin([]byte(body))
	if !containsString(got, jwt) || !containsString(got, refresh) {
		t.Fatalf("%q", got)
	}
	if containsString(got, "ORDINARYWORDLONGXX") {
		t.Fatalf("plain token extracted %q", got)
	}
}

func containsString(list []string, want string) bool {
	for _, item := range list {
		if item == want {
			return true
		}
	}
	return false
}
