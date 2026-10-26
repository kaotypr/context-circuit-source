package surface

import "strings"

// helpScenariosFor adds the scenario of a help topic whose text changed.
// topics.go itself maps to orient. A known topic adds its command's scenarios.
func helpScenariosFor(source string, outputs []OutputChange) []string {
	if source != "internal/cli/topics.go" {
		return nil
	}
	var out []string
	for _, item := range outputs {
		if !item.Changed {
			continue
		}
		topic, ok := strings.CutPrefix(item.Name, "help-")
		if !ok {
			continue
		}
		out = append(out, scenariosForHelpTopic(topic)...)
	}
	return out
}

func scenariosForHelpTopic(topic string) []string {
	switch topic {
	case "record":
		return []string{"intent", "plan"}
	case "worktree":
		return []string{"execute"}
	case "knowledge":
		return []string{"borrowed"}
	case "workspace", "repo", "member", "agent":
		return []string{"join"}
	case "context":
		return []string{"complete"}
	default:
		return nil
	}
}
