// Package assets contains only the allowlisted blank workspace and instructions.
package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed product/AGENTS.md.in product/CLAUDE.md.in product/CURSOR.md.in product/cursor.mdc.in product/README.md product/docs/commands.md product/docs/workspace.md product/docs/worktrees.md product/docs/working.md template/gitignore.in template/workspace.yaml template/members.yaml template/ids.yaml template/context/INDEX.md template/intent/README.md template/plans/README.md template/sources/README.md scripts/release-manifest.txt VERSION
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
