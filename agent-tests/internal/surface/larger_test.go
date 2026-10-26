package surface

import (
	"strings"
	"testing"
)

func TestLayoutReasons(t *testing.T) {
	files := []FileChange{
		{Status: "changed", Path: "AGENTS.md"},
		{Status: "added", Path: "intent/i100-demo.md"},
		{Status: "moved", From: ".context-circuit/docs/old.md", Path: ".context-circuit/docs/new.md"},
		{Status: "removed", Path: ".context-circuit/ids/maya.yaml"},
	}
	got := LayoutReasons(files)
	if len(got) != 1 || got[0].Kind != "layout" || !strings.Contains(got[0].Detail, "old.md ->") {
		t.Fatalf("%#v", got)
	}
}

func TestSkillReasons(t *testing.T) {
	body := "---\nname: cc-intent\n---\n# Title\n\nBody.\n"
	base := map[string][]byte{".agents/skills/cc-intent/SKILL.md": []byte(body)}
	branch := map[string][]byte{".agents/skills/cc-intent/SKILL.md": []byte(body + "Contributors read this extra sentence.\n")}
	if reasons := SkillReasons(base, branch); len(reasons) != 0 {
		t.Fatalf("body edit %#v", reasons)
	}
	headed := map[string][]byte{".agents/skills/cc-intent/SKILL.md": []byte(body + "# Extra\n")}
	if reasons := SkillReasons(base, headed); len(reasons) != 1 || reasons[0].Kind != "skill" {
		t.Fatalf("heading %#v", reasons)
	}
	added := map[string][]byte{".agents/skills/cc-new/SKILL.md": []byte(body)}
	if reasons := SkillReasons(nil, added); len(reasons) != 1 || !strings.Contains(reasons[0].Detail, "added") {
		t.Fatalf("added %#v", reasons)
	}
}

func TestCLIReasons(t *testing.T) {
	help := CLIReasons([]OutputChange{{Name: "help", Changed: true, Base: "a", Branch: "b"}}, nil, nil)
	if len(help) != 1 || !strings.Contains(help[0].Detail, "help") {
		t.Fatalf("%#v", help)
	}
	keys := CLIReasons([]OutputChange{{
		Name:         "status-json",
		Changed:      true,
		StdoutBase:   `{"ok":true}`,
		StdoutBranch: `{"ok":true,"extra":1}`,
	}}, nil, nil)
	if len(keys) != 1 || !strings.Contains(keys[0].Detail, "JSON keys") {
		t.Fatalf("%#v", keys)
	}
	check := CLIReasons([]OutputChange{{Name: "check", Changed: true}}, nil, nil)
	if len(check) != 1 || !strings.Contains(check[0].Detail, "check findings") {
		t.Fatalf("%#v", check)
	}
}

func TestGateReasonsHeuristic(t *testing.T) {
	if got := GateReasons([]string{"commit_recommended stays"}); len(got) != 0 {
		t.Fatalf("underscore word %#v", got)
	}
	got := GateReasons([]string{"do not push the branch"})
	if len(got) != 1 || got[0].Kind != "gate" || !strings.Contains(got[0].Detail, "heuristic") || !strings.Contains(got[0].Detail, "push") {
		t.Fatalf("%#v", got)
	}
}

func TestMigrationReasons(t *testing.T) {
	base := map[string][]byte{
		"workspace.yaml":           []byte("name: Acme\nversion: 3\n"),
		".context-circuit/VERSION": []byte("1\n"),
	}
	branch := map[string][]byte{
		"workspace.yaml":           []byte("name: Acme\nversion: 3\n"),
		".context-circuit/VERSION": []byte("2\n"),
	}
	if got := MigrationReasons(base, branch); len(got) != 0 {
		t.Fatalf("version file %#v", got)
	}
	branch["workspace.yaml"] = []byte("name: Acme\nversion: 4\n")
	got := MigrationReasons(base, branch)
	if len(got) != 1 || got[0].Kind != "migration" || !strings.Contains(got[0].Detail, "3") {
		t.Fatalf("schema %#v", got)
	}
	branch["workspace.yaml"] = base["workspace.yaml"]
	branch[".context-circuit/ids/maya.yaml"] = []byte("band: 1\n")
	got = MigrationReasons(base, branch)
	if len(got) != 1 || !strings.Contains(got[0].Detail, "ids file layout") {
		t.Fatalf("ids %#v", got)
	}
}

func TestMeasureSizes(t *testing.T) {
	base := map[string][]byte{"AGENTS.md": []byte("12345678")}
	branch := map[string][]byte{"AGENTS.md": []byte("12345678abcd")}
	sizes := measureSizes(base, branch)
	var found SizeChange
	for _, size := range sizes {
		if size.Path == "AGENTS.md" {
			found = size
		}
	}
	if found.BaseBytes != 8 || found.BranchBytes != 12 || found.DeltaBytes != 4 {
		t.Fatalf("%#v", found)
	}
	if found.BaseTokens != 2 || found.BranchTokens != 3 || found.DeltaTokens != 1 {
		t.Fatalf("tokens %#v", found)
	}
}
