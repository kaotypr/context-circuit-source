package surface

import (
	"encoding/json"
	"sort"
)

func jsonKeyList(text string) ([]string, bool) {
	var v any
	if err := json.Unmarshal([]byte(text), &v); err != nil {
		return nil, false
	}
	var keys []string
	collectKeys(v, "", &keys)
	sort.Strings(keys)
	return keys, true
}

func collectKeys(v any, prefix string, keys *[]string) {
	switch t := v.(type) {
	case map[string]any:
		names := make([]string, 0, len(t))
		for name := range t {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			next := name
			if prefix != "" {
				next = prefix + "." + name
			}
			*keys = append(*keys, next)
			collectKeys(t[name], next, keys)
		}
	case []any:
		if len(t) == 0 {
			return
		}
		collectKeys(t[0], prefix+"[]", keys)
	}
}
