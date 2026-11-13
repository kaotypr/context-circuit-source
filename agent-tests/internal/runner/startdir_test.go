package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

// A start_dir resolves under the active workspace, after workspace: has
// moved it, and must exist once setup has run.
func TestHostStartDir(t *testing.T) {
	root := t.TempDir()
	w := world.World{RunRoot: root, WS: filepath.Join(root, "ws")}
	worktree := filepath.Join(w.WS, ".worktrees", "p1000", "api")
	if err := os.MkdirAll(worktree, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.WS, "file"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := hostStartDir(w, scenario.Scenario{ID: "s"})
	if err != nil || got != "" {
		t.Fatalf("no start_dir: %q %v", got, err)
	}
	got, err = hostStartDir(w, scenario.Scenario{ID: "s", StartDir: ".worktrees/p1000/api"})
	if err != nil || got != worktree {
		t.Fatalf("worktree: %q %v", got, err)
	}
	for _, bad := range []string{".worktrees/p1001/api", "file", "../ws", "/abs"} {
		if got, err := hostStartDir(w, scenario.Scenario{ID: "s", StartDir: bad}); err == nil {
			t.Errorf("%q resolved to %q", bad, got)
		}
	}

	ws2 := filepath.Join(root, "ws2", "repositories", "api")
	if err := os.MkdirAll(ws2, 0o755); err != nil {
		t.Fatal(err)
	}
	sc := scenario.Scenario{ID: "s", Workspace: "ws2", StartDir: "repositories/api"}
	applyWorkspace(&w, sc)
	got, err = hostStartDir(w, sc)
	if err != nil || got != ws2 {
		t.Fatalf("ws2: %q %v", got, err)
	}
}
