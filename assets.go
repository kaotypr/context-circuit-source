// Package assets contains only the allowlisted blank workspace and instructions.
package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed product/AGENTS.md.in product/CLAUDE.md.in product/CURSOR.md.in product/cursor.mdc.in product/README.md product/assets/readme/context-circuit-mark.svg product/assets/readme/context-circuit-mark.png product/assets/readme/context-circuit-logo.png product/assets/readme/knowledge-circuit.png product/assets/readme/workflow-overview.png product/assets/readme/shared-vs-local.png product/assets/readme/social-preview.png product/docs/commands.md product/docs/workspace.md product/docs/worktrees.md product/docs/working.md product/docs/agents.md template/gitignore.in template/workspace.yaml template/members.yaml template/ids.yaml template/role-tiering.yaml template/context/INDEX.md template/context/glossary.md template/intent/README.md template/plans/README.md template/sources/README.md scripts/release-manifest.txt VERSION CLI_VERSION
//go:embed product/skills/cc-cli/SKILL.md product/skills/cc-cli/route.md product/skills/cc-cli/scripts/install.sh product/skills/cc-cli/scripts/install.ps1 product/skills/cc-dispatch/SKILL.md product/skills/cc-dispatch/route.md product/skills/cc-workspace/SKILL.md product/skills/cc-workspace/route.md product/skills/cc-intent/SKILL.md product/skills/cc-intent/route.md product/skills/cc-plan/SKILL.md product/skills/cc-plan/route.md product/skills/cc-direct/SKILL.md product/skills/cc-direct/route.md product/skills/cc-knowledge/SKILL.md product/skills/cc-knowledge/route.md product/skills/cc-worktree/SKILL.md product/skills/cc-worktree/route.md product/skills/cc-stacked/SKILL.md product/skills/cc-stacked/route.md product/skills/cc-review/SKILL.md product/skills/cc-review/route.md product/skills/cc-complete/SKILL.md product/skills/cc-complete/route.md product/skills/cc-deliver/SKILL.md product/skills/cc-deliver/route.md
var bundled embed.FS

// Files resolves the same explicit source-to-destination manifest used by releases.
func Files() (map[string][]byte, error) {
	manifest, err := bundled.ReadFile("scripts/release-manifest.txt")
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	for _, line := range strings.Split(string(manifest), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 || !fs.ValidPath(fields[0]) || !fs.ValidPath(fields[1]) {
			return nil, fmt.Errorf("invalid release manifest line: %q", line)
		}
		if _, exists := files[fields[1]]; exists {
			return nil, fmt.Errorf("duplicate release path: %s", fields[1])
		}
		data, err := bundled.ReadFile(fields[0])
		if err != nil {
			return nil, err
		}
		files[fields[1]] = data
	}
	return files, nil
}
