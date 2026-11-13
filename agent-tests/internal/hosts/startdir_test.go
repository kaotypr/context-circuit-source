package hosts

import (
	"path/filepath"
	"slices"
	"testing"
)

// argAfter is the argument that follows flag, or "" when flag is absent.
func argAfter(argv []string, flag string) string {
	for i, arg := range argv {
		if arg == flag && i+1 < len(argv) {
			return argv[i+1]
		}
	}
	return ""
}

// addDirs is every value passed with --add-dir.
func addDirs(argv []string) []string {
	var out []string
	for i, arg := range argv {
		if arg == "--add-dir" && i+1 < len(argv) {
			out = append(out, argv[i+1])
		}
	}
	return out
}

// Each host starts in the workspace unless the launch names a start
// directory, and then starts there.
func TestBuildStartDir(t *testing.T) {
	const ws = "/work/ws"
	worktree := filepath.Join(ws, ".worktrees", "p1000", "api")
	for _, name := range []string{"claude", "codex", "cursor"} {
		for _, start := range []string{"", worktree} {
			t.Run(name+"/"+map[bool]string{true: "workspace", false: "worktree"}[start == ""], func(t *testing.T) {
				h, err := ForName(name)
				if err != nil {
					t.Fatal(err)
				}
				launch, err := h.Build(LaunchInput{
					Prompt:    "Reply OK",
					Workspace: ws,
					StartDir:  start,
					RunRoot:   "/work/run",
					SessionID: "11111111-1111-1111-1111-111111111111",
					Cfg:       HostConfig{Command: name, Model: "m"},
					Caller:    map[string]string{"PATH": "/usr/bin:/bin"},
				})
				if err != nil {
					t.Fatal(err)
				}
				want := ws
				if start != "" {
					want = start
				}
				if launch.Dir != want {
					t.Fatalf("dir %q, want %q", launch.Dir, want)
				}
				switch name {
				case "codex":
					if got := argAfter(launch.Argv, "-C"); got != want {
						t.Fatalf("-C %q, want %q", got, want)
					}
					dirs := addDirs(launch.Argv)
					// The role directory stays the workspace's, and a run
					// started below the workspace keeps all of it writable.
					if !slices.Contains(dirs, filepath.Join(ws, ".codex")) {
						t.Fatalf("add-dirs %v lack the workspace .codex", dirs)
					}
					if slices.Contains(dirs, ws) != (start != "") {
						t.Fatalf("add-dirs %v for start %q", dirs, start)
					}
					if !slices.Contains(launch.MkDirs, filepath.Join(ws, ".codex")) {
						t.Fatalf("mkdirs %v", launch.MkDirs)
					}
				case "cursor":
					if got := argAfter(launch.Argv, "--workspace"); got != want {
						t.Fatalf("--workspace %q, want %q", got, want)
					}
				}
			})
		}
	}
}
