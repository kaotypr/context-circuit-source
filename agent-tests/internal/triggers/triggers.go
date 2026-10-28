// Package triggers loads the source-path trigger map.
package triggers

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
)

// Kinds are the before/after artifact names surface-diff prints.
var Kinds = map[string]bool{
	"always-loaded":    true,
	"skill-text":       true,
	"skill-route":      true,
	"shipped-docs":     true,
	"structure":        true,
	"cli-action":       true,
	"record-template":  true,
	"check-findings":   true,
	"role-definitions": true,
	"launcher":         true,
	"template-data":    true,
	"host-settings":    true,
}

// Row is one source-path rule.
// Hosts empty means any one host. Larger means this path is a larger
// change on its own; surface-diff can still call other diffs larger.
type Row struct {
	Match     string   `yaml:"match"`
	Scenarios []string `yaml:"scenarios"`
	Hosts     []string `yaml:"hosts"`
	Kind      string   `yaml:"kind"`
	Larger    bool     `yaml:"larger"`
	Note      string   `yaml:"note"`
	glob      *regexp.Regexp
}

// Section maps one ## heading of product/AGENTS.md.in to extra scenarios.
// An empty scenario list is a heading the catalogue does not cover.
type Section struct {
	Heading   string   `yaml:"heading"`
	Scenarios []string `yaml:"scenarios"`
}

// Table is triggers.yaml.
type Table struct {
	Rows     []Row     `yaml:"rows"`
	None     []string  `yaml:"none"`
	Sections []Section `yaml:"sections"`
	none     []*regexp.Regexp
}

// Hit is the union of the rows and section rules that apply to one path.
type Hit struct {
	Scenarios   []string
	Hosts       []string
	Kinds       []string
	Notes       []string
	Larger      bool
	Rows        []Row
	OrientHosts []string
}

// Load reads triggers.yaml.
func Load(path string) (Table, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Table{}, err
	}
	return LoadBytes(b)
}

// LoadBytes parses a trigger map. Globs are compiled. Scenario and host
// names are checked by Validate, once the caller knows those files.
func LoadBytes(b []byte) (Table, error) {
	var table Table
	if err := yaml.UnmarshalWithOptions(b, &table, yaml.Strict()); err != nil {
		return Table{}, err
	}
	if len(table.Rows) == 0 {
		return Table{}, fmt.Errorf("triggers: no rows")
	}
	seen := map[string]bool{}
	for i := range table.Rows {
		row := &table.Rows[i]
		if row.Match == "" {
			return Table{}, fmt.Errorf("triggers: row %d has an empty match", i)
		}
		if seen[row.Match] {
			return Table{}, fmt.Errorf("triggers: duplicate match %q", row.Match)
		}
		seen[row.Match] = true
		if !Kinds[row.Kind] {
			return Table{}, fmt.Errorf("triggers: %s: unknown kind %q", row.Match, row.Kind)
		}
		if len(row.Scenarios) == 0 {
			return Table{}, fmt.Errorf("triggers: %s: no scenarios", row.Match)
		}
		glob, err := compileGlob(row.Match)
		if err != nil {
			return Table{}, fmt.Errorf("triggers: %w", err)
		}
		row.glob = glob
	}
	headings := map[string]bool{}
	for _, sec := range table.Sections {
		if sec.Heading == "" {
			return Table{}, fmt.Errorf("triggers: a section heading is empty")
		}
		if headings[sec.Heading] {
			return Table{}, fmt.Errorf("triggers: duplicate section %q", sec.Heading)
		}
		headings[sec.Heading] = true
	}
	for _, pattern := range table.None {
		glob, err := compileGlob(pattern)
		if err != nil {
			return Table{}, fmt.Errorf("triggers: none: %w", err)
		}
		table.none = append(table.none, glob)
	}
	return table, nil
}

// Validate checks scenario ids and host names against the catalogues.
func (t Table) Validate(scenarios, hosts map[string]bool) error {
	var problems []string
	check := func(where string, names []string, set map[string]bool, label string) {
		seen := map[string]bool{}
		for _, name := range names {
			if seen[name] {
				problems = append(problems, fmt.Sprintf("%s: duplicate %s %q", where, label, name))
				continue
			}
			seen[name] = true
			if !set[name] {
				problems = append(problems, fmt.Sprintf("%s: %s %q is not in the catalogue", where, label, name))
			}
		}
	}
	for _, row := range t.Rows {
		check(row.Match, row.Scenarios, scenarios, "scenario")
		check(row.Match, row.Hosts, hosts, "host")
	}
	for _, sec := range t.Sections {
		check("section "+sec.Heading, sec.Scenarios, scenarios, "scenario")
	}
	if len(problems) > 0 {
		return fmt.Errorf("triggers:\n%s", strings.Join(problems, "\n"))
	}
	return nil
}

// Match returns every row whose glob matches path. path uses slashes.
func (t Table) Match(path string) []Row {
	path = filepath.ToSlash(path)
	var out []Row
	for _, row := range t.Rows {
		if row.glob != nil && row.glob.MatchString(path) {
			out = append(out, row)
		}
	}
	return out
}

// NoScenarios reports whether path is on the none list.
func (t Table) NoScenarios(path string) bool {
	path = filepath.ToSlash(path)
	for _, glob := range t.none {
		if glob.MatchString(path) {
			return true
		}
	}
	return false
}

// SectionScenarios returns the extra scenarios for those headings.
// Unknown headings are ignored; the caller lists them as having no row.
func (t Table) SectionScenarios(headings []string) []string {
	byHeading := map[string][]string{}
	for _, sec := range t.Sections {
		byHeading[sec.Heading] = sec.Scenarios
	}
	var out []string
	for _, heading := range headings {
		out = append(out, byHeading[heading]...)
	}
	return unique(out)
}

// Apply unions the rows for path. Headings are consulted only for
// product/AGENTS.md.in, whose file row is orient and whose sections add
// the rest.
func (t Table) Apply(path string, headings []string) Hit {
	var hit Hit
	for _, row := range t.Match(path) {
		hit.Rows = append(hit.Rows, row)
		hit.Scenarios = append(hit.Scenarios, row.Scenarios...)
		hit.Hosts = append(hit.Hosts, row.Hosts...)
		hit.Kinds = append(hit.Kinds, row.Kind)
		if row.Note != "" {
			hit.Notes = append(hit.Notes, row.Note)
		}
		if row.Larger {
			hit.Larger = true
		}
	}
	if filepath.ToSlash(path) == "product/AGENTS.md.in" {
		hit.Scenarios = append(hit.Scenarios, t.SectionScenarios(headings)...)
	}
	hit.Scenarios = unique(hit.Scenarios)
	hit.Hosts = unique(hit.Hosts)
	hit.Kinds = unique(hit.Kinds)
	hit.Notes = unique(hit.Notes)
	return hit
}

// StructureHit is the union of rows whose kind is structure.
// surface-diff uses it when a shipped file moved or a manifest
// destination was added, removed, or renamed.
func (t Table) StructureHit() Hit {
	var hit Hit
	for _, row := range t.Rows {
		if row.Kind != "structure" {
			continue
		}
		hit.Rows = append(hit.Rows, row)
		hit.Scenarios = append(hit.Scenarios, row.Scenarios...)
		hit.Hosts = append(hit.Hosts, row.Hosts...)
		hit.Kinds = append(hit.Kinds, row.Kind)
		if row.Note != "" {
			hit.Notes = append(hit.Notes, row.Note)
		}
		hit.Larger = true
	}
	hit.Scenarios = unique(hit.Scenarios)
	hit.Hosts = unique(hit.Hosts)
	hit.Kinds = unique(hit.Kinds)
	hit.Notes = unique(hit.Notes)
	// Every scenario runs on any one host. orient also runs on each host.
	hit.OrientHosts = []string{"claude", "codex", "cursor"}
	return hit
}

func unique(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	for _, item := range in {
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	sort.Strings(out)
	return out
}
