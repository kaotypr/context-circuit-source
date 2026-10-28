package surface

import (
	"slices"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/triggers"
)

func TestSelectSources(t *testing.T) {
	table := testTable(t)
	manifest := map[string]string{
		".agents/skills/cc-intent/SKILL.md": "product/skills/cc-intent/SKILL.md",
	}
	skill := SelectSources(table, SourceInput{
		GitPaths:       []string{"product/skills/cc-intent/SKILL.md", "internal/cli/topics.go", "assets.go"},
		Files:          []FileChange{{Status: "changed", Path: ".agents/skills/cc-intent/SKILL.md"}},
		ManifestBase:   manifest,
		ManifestBranch: manifest,
	})
	if !slices.Equal(skill.Paths, []string{"product/skills/cc-intent/SKILL.md"}) {
		t.Fatalf("skill paths %v", skill.Paths)
	}
	if skill.EveryScenario {
		t.Fatal("skill edit is not every scenario")
	}

	help := SelectSources(table, SourceInput{
		GitPaths:       []string{"internal/cli/topics.go"},
		OutputsChanged: true,
	})
	if !slices.Equal(help.Paths, []string{"internal/cli/topics.go"}) {
		t.Fatalf("help paths %v", help.Paths)
	}

	baseMan := map[string]string{".context-circuit/docs/how-it-works.md": "product/docs/how-it-works.md"}
	branchMan := map[string]string{".context-circuit/docs/how-it-works-renamed.md": "product/docs/how-it-works.md"}
	moved := SelectSources(table, SourceInput{
		GitPaths: []string{"scripts/release-manifest.txt"},
		Files: []FileChange{{
			Status: "moved",
			From:   ".context-circuit/docs/how-it-works.md",
			Path:   ".context-circuit/docs/how-it-works-renamed.md",
		}},
		ManifestBase:   baseMan,
		ManifestBranch: branchMan,
	})
	if !slices.Equal(moved.Paths, []string{"scripts/release-manifest.txt"}) || !moved.EveryScenario {
		t.Fatalf("move %+v", moved)
	}

	ignored := SelectSources(table, SourceInput{GitPaths: []string{"VERSION", "product/skills/cc-intent/SKILL.md"}})
	if !slices.Equal(ignored.Ignored, []string{"VERSION"}) {
		t.Fatalf("ignored %v", ignored.Ignored)
	}
}

func testTable(t *testing.T) triggers.Table {
	t.Helper()
	table, err := triggers.LoadBytes([]byte(`
rows:
  - match: product/skills/cc-intent/SKILL.md
    scenarios: [intent]
    kind: skill-text
  - match: internal/cli/topics.go
    scenarios: [orient]
    kind: cli-action
  - match: scripts/release-manifest.txt
    scenarios: [orient]
    kind: structure
    larger: true
none:
  - VERSION
`))
	if err != nil {
		t.Fatal(err)
	}
	return table
}
