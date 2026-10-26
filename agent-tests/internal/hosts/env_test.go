package hosts

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestBaseEnvDropsDeniedCallerNames(t *testing.T) {
	caller := map[string]string{
		"PATH":                 "/usr/bin:/bin",
		"UNRELATED_CALLER_VAR": "nope-unrelated",
	}
	var auth []string
	for _, pattern := range PassThroughDenied {
		name := strings.TrimSuffix(pattern, "*")
		if strings.HasSuffix(pattern, "*") {
			name += "SECRET"
		}
		caller[name] = "nope-" + name
		auth = append(auth, name)
	}
	in := LaunchInput{
		DecoyDir: "/work/run/decoy",
		RunRoot:  "/work/run",
		Caller:   caller,
		Cfg:      HostConfig{AuthEnv: auth},
	}
	env := BaseEnv(in, nil)
	sorted := append([]string(nil), env...)
	sort.Strings(sorted)
	if !reflect.DeepEqual(sorted, env) {
		t.Fatalf("not sorted: %v", env)
	}
	for _, e := range env {
		if strings.Contains(e, "nope-") {
			t.Fatalf("passed through %s", e)
		}
		key, _, _ := strings.Cut(e, "=")
		if key == "HOME" {
			t.Fatalf("HOME set from caller: %s", e)
		}
	}
	got := envMap(t, env)
	if got["PATH"] != "/work/run/decoy:/usr/bin:/bin" {
		t.Fatalf("PATH %s", got["PATH"])
	}
	if got["LANG"] != "C.UTF-8" || got["NO_COLOR"] != "1" || got["GIT_CONFIG_NOSYSTEM"] != "1" || got["GIT_TERMINAL_PROMPT"] != "0" {
		t.Fatalf("%v", got)
	}
	if got["TMPDIR"] != "/work/run/tmp" || got["GOCACHE"] != "/work/run/gocache" || got["GIT_CONFIG_GLOBAL"] != "/work/run/gitconfig" {
		t.Fatalf("%v", got)
	}
}

func TestForName(t *testing.T) {
	for _, name := range []string{"claude", "codex", "cursor"} {
		h, err := ForName(name)
		if err != nil || h.Name() != name {
			t.Fatalf("%s: %v %v", name, h, err)
		}
	}
	if _, err := ForName("nope"); err == nil {
		t.Fatal("expected unknown host")
	}
}
