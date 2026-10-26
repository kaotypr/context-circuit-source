package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The skills state facts about what the CLI returns and what the installers
// accept. These tests take each fact from the real command or script, so a
// skill that drifts from it fails here rather than misleading an agent.

// cc-complete names the key record complete actually returns, and keeps
// completed_at for the frontmatter field it is.
func TestTheCompletionSkillNamesTheResultKey(t *testing.T) {
	f := setup(t)
	f.repository("api")
	p := record(t, f.ok("record", "create", "--kind", "plan", "--slug", "alone", "--title", "Alone", "--repo", "api"))
	var result map[string]any
	if err := json.Unmarshal([]byte(f.ok("record", "complete", "--id", p.ID, "--text", "Done.")), &result); err != nil {
		t.Fatal(err)
	}
	if result["completed"] != p.ID {
		t.Fatalf("record complete returned %v, want completed=%s", result, p.ID)
	}
	if _, ok := result["completed_at"]; ok {
		t.Fatal("record complete now returns completed_at; update cc-complete")
	}
	if !strings.Contains(read(t, filepath.Join(f.root, filepath.FromSlash(p.Path))), "completed_at:") {
		t.Fatal("record complete no longer stamps completed_at in the plan")
	}
	skill := strings.Join(strings.Fields(productFile(t, filepath.Join("skills", "cc-complete", "SKILL.md"))), " ")
	if !strings.Contains(skill, "carries `completed`") {
		t.Error("cc-complete does not name the completed result key")
	}
	if strings.Contains(skill, "result carries `completed_at`") {
		t.Error("cc-complete says the result carries completed_at")
	}
}

// Running a plan does not edit it, so no shipped text may tell the agent to
// write into a plan during a run.
func TestNothingRecordsIntoAPlanDuringARun(t *testing.T) {
	var files []string
	for _, name := range shippedSkills(t) {
		files = append(files, filepath.Join("skills", name, "SKILL.md"))
	}
	for _, name := range shippedDocs(t) {
		files = append(files, filepath.Join("docs", name))
	}
	for _, rel := range append(files, "AGENTS.md.in") {
		text := strings.Join(strings.Fields(productFile(t, rel)), " ")
		if strings.Contains(text, "in the plan so it can be audited") {
			t.Errorf("%s tells the agent to record into a plan during a run", rel)
		}
	}
}

// cc-cli explains every reserved label and the version range the installers
// enforce. Both are read from the installer, so a label or a range added there
// without the skill fails here.
func TestTheCLISkillMatchesTheInstallers(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "product", "skills", "cc-cli", "scripts", "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	script := string(data)
	labels := regexp.MustCompile(`reserved_labels='([^']+)'`).FindStringSubmatch(script)
	if labels == nil {
		t.Fatal("install.sh no longer lists reserved_labels")
	}
	if !strings.Contains(script, `case "$version" in 2.*)`) {
		t.Fatal("install.sh no longer accepts only 2.* versions; update this test and cc-cli")
	}
	skill := strings.Join(strings.Fields(productFile(t, filepath.Join("skills", "cc-cli", "SKILL.md"))), " ")
	for _, label := range strings.Fields(labels[1]) {
		if !strings.Contains(skill, "`"+label+"`") {
			t.Errorf("cc-cli does not name the reserved label %q", label)
		}
	}
	for _, fact := range []string{"`2.*`", "prerelease of the stable line", "`context-circuit-cli`"} {
		if !strings.Contains(skill, fact) {
			t.Errorf("cc-cli does not say %s", fact)
		}
	}
}
