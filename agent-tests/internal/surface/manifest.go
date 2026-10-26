package surface

import (
	"fmt"
	"strings"
)

// ParseManifest reads release-manifest.txt. The result maps a workspace
// path to the source path. A source may be listed more than once; a
// destination may not.
func ParseManifest(text string) (map[string]string, error) {
	out := map[string]string{}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid release manifest line %q", line)
		}
		source, dest := fields[0], fields[1]
		if _, exists := out[dest]; exists {
			return nil, fmt.Errorf("duplicate release path %s", dest)
		}
		out[dest] = source
	}
	return out, nil
}

// WithReadme adds the workspace front page. assets.Readme serves
// product/workspace-README.md and the release manifest does not list it,
// because init writes that file once rather than exporting it.
func WithReadme(manifest map[string]string) map[string]string {
	if manifest == nil {
		manifest = map[string]string{}
	}
	if _, ok := manifest["README.md"]; !ok {
		manifest["README.md"] = "product/workspace-README.md"
	}
	return manifest
}
