package hosts

import (
	"path/filepath"
	"sort"
	"strings"
)

// PassThroughDenied is the caller environment BaseEnv never copies.
// A trailing * matches a name prefix. HOME is set only when a host passes
// it in extra, never from the caller.
var PassThroughDenied = []string{
	"HOME",
	"XDG_*",
	"SSH_AUTH_SOCK",
	"SSH_AGENT_PID",
	"GITHUB_TOKEN",
	"GH_TOKEN",
	"AWS_*",
	"KUBECONFIG",
}

// BaseEnv builds the complete child environment, sorted KEY=VAL.
// PATH is DecoyDir, then the caller's PATH when that value is non-empty.
// A trailing colon is not added: an empty PATH entry is the current
// directory. Nothing else is copied from Caller except the auth names in
// Cfg.AuthEnv, and those are still dropped when they are on
// PassThroughDenied. extra is applied last, so a host can set HOME there.
func BaseEnv(in LaunchInput, extra map[string]string) []string {
	caller := in.Caller
	if caller == nil {
		caller = map[string]string{}
	}
	path := in.DecoyDir
	if rest := caller["PATH"]; rest != "" {
		path += ":" + rest
	}
	env := map[string]string{
		"PATH":                path,
		"TMPDIR":              filepath.Join(in.RunRoot, "tmp"),
		"GOCACHE":             filepath.Join(in.RunRoot, "gocache"),
		"LANG":                "C.UTF-8",
		"GIT_CONFIG_NOSYSTEM": "1",
		"GIT_CONFIG_GLOBAL":   filepath.Join(in.RunRoot, "gitconfig"),
		"GIT_TERMINAL_PROMPT": "0",
		"NO_COLOR":            "1",
	}
	for _, name := range in.Cfg.AuthEnv {
		if deniedPassThrough(name) {
			continue
		}
		if v, ok := caller[name]; ok && v != "" {
			env[name] = v
		}
	}
	for k, v := range extra {
		env[k] = v
	}
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, len(keys))
	for i, k := range keys {
		out[i] = k + "=" + env[k]
	}
	return out
}

func deniedPassThrough(name string) bool {
	for _, pattern := range PassThroughDenied {
		if strings.HasSuffix(pattern, "*") {
			if strings.HasPrefix(name, strings.TrimSuffix(pattern, "*")) {
				return true
			}
			continue
		}
		if name == pattern {
			return true
		}
	}
	return false
}

func cloneMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func authFromEnv(caller map[string]string, names []string) (bool, string) {
	if caller == nil {
		caller = map[string]string{}
	}
	if len(names) == 0 {
		return false, "no auth env configured"
	}
	for _, name := range names {
		if caller[name] != "" {
			return true, ""
		}
	}
	return false, "missing " + strings.Join(names, " or ")
}
