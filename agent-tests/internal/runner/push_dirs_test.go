package runner

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/scenario"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

// Only a scenario that lets the run push gets the workspace's .git and its
// remote as writable directories; every other scenario keeps the host
// sandbox guarding them.
func TestPushWritableDirsOnlyForPushScenarios(t *testing.T) {
	w := world.World{WS: "/run/ws", RunRoot: "/run"}
	if got := pushWritableDirs(scenario.Scenario{}, w); got != nil {
		t.Fatalf("plain scenario got %v", got)
	}
	got := pushWritableDirs(scenario.Scenario{PushWorkspaceBase: true}, w)
	want := []string{filepath.Join("/run/ws", ".git"), filepath.Join("/run", "remotes", scenario.WorkspaceRemote)}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
