package scenario

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/proc"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

func TestExecNilIsInconclusive(t *testing.T) {
	for _, name := range []string{"execute-worktree", "execute-tests-pass"} {
		got := Registry[name](CheckContext{})
		if !got.Pass || !got.Inconclusive {
			t.Fatalf("%s %+v", name, got)
		}
	}
	if softCommitSubjects(CheckContext{}) != "inconclusive" {
		t.Fatal(softCommitSubjects(CheckContext{}))
	}
}

func TestOrientReplyShape(t *testing.T) {
	good := CheckContext{
		Parsed: hosts.Parsed{FinalReply: "The api repository. Use sh .context-circuit/cli. A draft invoice is voided. An issued invoice is credited. See domains/invoices.md."},
		Stream: `{"file_path":"context/INDEX.md"}`,
	}
	if !Registry[nameOrientReply](good).Pass {
		t.Fatal(Registry[nameOrientReply](good))
	}
	bad := good
	bad.Parsed.FinalReply = "No."
	bad.Stream = ""
	if Registry[nameOrientReply](bad).Pass {
		t.Fatal("short reply passed")
	}
	canary := good
	canary.Stream += codeCanary
	if Registry[nameOrientReply](canary).Pass {
		t.Fatal("canary passed")
	}
}

func TestLanguageHeuristic(t *testing.T) {
	id := `
# Ekspor CSV faktur
## Goal
Pelanggan perlu mengunduh faktur sebagai berkas CSV dari layanan api.
## Non-goals
Pembatalan faktur draf dan nota kredit faktur terbit tetap seperti sekarang.
## Constraints
Perubahan berada di repositori api pada paket internal/invoice tanpa dependensi baru.
## Success criteria
- Pemanggil dapat menuliskan faktur sebagai CSV beserta nomor, pelanggan, status, dan total dalam sen.
## Repository scope
Repositori api.
## Open questions
1. Tidak ada pertanyaan terbuka.
`
	words := proseWords(id)
	if len(words) < 20 {
		t.Fatalf("words %d", len(words))
	}
	if englishRatio(words) >= notEnglishRatio {
		t.Fatalf("ratio %v", englishRatio(words))
	}
	english, err := os.ReadFile(filepath.Join("..", "..", "fixtures", "records", "intent-csv.md"))
	if err != nil {
		t.Fatal(err)
	}
	if englishRatio(proseWords(string(english))) < notEnglishRatio {
		t.Fatalf("english ratio %v", englishRatio(proseWords(string(english))))
	}
}

func TestDirectReadmePorcelain(t *testing.T) {
	dir := t.TempDir()
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("Customers receive a text summary.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := CheckContext{
		World:  world.World{API: dir},
		Before: Snapshot{Heads: map[string]string{"api": "abc"}},
		After: Snapshot{
			Heads:  map[string]string{"api": "abc"},
			Status: map[string]string{"api": " M README.md\n"},
		},
	}
	if !Registry[nameDirectReadme](ctx).Pass {
		t.Fatal(Registry[nameDirectReadme](ctx))
	}
	if err := os.WriteFile(readme, []byte("Customers recieve a text summary.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if Registry[nameDirectReadme](ctx).Pass {
		t.Fatal("typo passed")
	}
}

func TestMissingCheckoutFails(t *testing.T) {
	empty := CheckContext{}
	for _, name := range []string{nameOrientClean, nameIntentNoPlan, namePlanLinked} {
		got := Registry[name](empty)
		if got.Pass || !strings.Contains(got.Detail, "HEAD is missing") {
			t.Fatalf("%s %+v", name, got)
		}
	}
	recorded := func() CheckContext {
		heads := map[string]string{"ws": "aaa", "api": "bbb"}
		status := map[string]string{"ws": "", "api": ""}
		return CheckContext{
			Before: Snapshot{Heads: heads, Status: status},
			After:  Snapshot{Heads: heads, Status: status},
		}
	}
	if !Registry[nameOrientClean](recorded()).Pass {
		t.Fatal(Registry[nameOrientClean](recorded()))
	}
	if !Registry[nameIntentNoPlan](recorded()).Pass {
		t.Fatal(Registry[nameIntentNoPlan](recorded()))
	}
	noStatus := recorded()
	delete(noStatus.After.Status, "api")
	if Registry[nameIntentNoPlan](noStatus).Pass {
		t.Fatal("missing api status passed")
	}
}

func TestJoinBindingAndRoles(t *testing.T) {
	ws := t.TempDir()
	api := filepath.Join(ws, "repositories", "api")
	if err := os.MkdirAll(api, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ws, "member.local.yaml"), []byte("member: dana\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeRoles := func(host string) {
		t.Helper()
		dir, ext, ok := roleLocation(host)
		if !ok {
			t.Fatal(host)
		}
		base := filepath.Join(ws, filepath.FromSlash(dir))
		if err := os.MkdirAll(base, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, name := range roleFileNames {
			if err := os.WriteFile(filepath.Join(base, name+ext), []byte(name+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	bindings := func(workspace, apiPath string) {
		t.Helper()
		body := fmt.Sprintf("workspace:\n  path: %q\nbindings:\n  api:\n    path: %q\n", workspace, apiPath)
		if err := os.WriteFile(filepath.Join(ws, "repositories.local.yaml"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	ctx := CheckContext{World: world.World{WS: ws}, Host: "claude"}
	bindings(".", api)
	writeRoles("claude")
	if !Registry[nameJoinLocal](ctx).Pass {
		t.Fatal(Registry[nameJoinLocal](ctx))
	}
	bindings("/tmp/elsewhere", api)
	if got := Registry[nameJoinLocal](ctx); got.Pass || !strings.Contains(got.Detail, "active workspace") {
		t.Fatalf("%+v", got)
	}
	bindings(".", "/tmp/elsewhere")
	if got := Registry[nameJoinLocal](ctx); got.Pass || !strings.Contains(got.Detail, "repositories/api") {
		t.Fatalf("%+v", got)
	}
	bindings(".", api)
	if err := os.Remove(filepath.Join(ws, ".claude", "agents", "cc-worker.md")); err != nil {
		t.Fatal(err)
	}
	if got := Registry[nameJoinLocal](ctx); got.Pass || !strings.Contains(got.Detail, "claude") {
		t.Fatalf("%+v", got)
	}
	writeRoles("codex")
	if Registry[nameJoinLocal](ctx).Pass {
		t.Fatal("claude run accepted another host's roles")
	}
	ctx.Host = ""
	if !Registry[nameJoinLocal](ctx).Pass {
		t.Fatal(Registry[nameJoinLocal](ctx))
	}
	for _, name := range roleFileNames {
		if err := os.Remove(filepath.Join(ws, ".codex", "agents", name+".toml")); err != nil {
			t.Fatal(err)
		}
	}
	if Registry[nameJoinLocal](ctx).Pass {
		t.Fatal("partial role sets passed")
	}
}

func TestCSVExportMustCallEncodingCSV(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "internal", "invoice")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	comment := "package invoice\n\n// CSV export is mentioned here.\nfunc Note() string { return \"CSV\" }\n"
	if err := os.WriteFile(filepath.Join(dir, "export.go"), []byte(comment), 0o644); err != nil {
		t.Fatal(err)
	}
	if treeHasCSV(root) {
		t.Fatal("a comment counted as a CSV export")
	}
	real := "package invoice\n\nimport \"encoding/csv\"\n\nfunc CSV() { csv.NewWriter(nil) }\n"
	if err := os.WriteFile(filepath.Join(dir, "export.go"), []byte(real), 0o644); err != nil {
		t.Fatal(err)
	}
	if !treeHasCSV(root) {
		t.Fatal("encoding/csv function was not seen")
	}
}

func TestCompleteNoteRejectsUnsupported(t *testing.T) {
	ws := t.TempDir()
	if err := os.MkdirAll(filepath.Join(ws, "context", "domains"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(note, date string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(ws, "context", "domains", "invoices.md"), []byte(note), 0o644); err != nil {
			t.Fatal(err)
		}
		index := "- [Invoices](domains/invoices.md) — voiding · anchor · reviewed " + date + "\n"
		if err := os.WriteFile(filepath.Join(ws, "context", "INDEX.md"), []byte(index), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	ctx := CheckContext{World: world.World{WS: ws}}
	today := time.Now().UTC().Format("2006-01-02")
	write("Invoices export as a plain-text summary. CSV is not supported.\n", today)
	if got := Registry[nameCompleteNote](ctx); got.Pass {
		t.Fatal("unsupported note passed")
	}
	write("Invoices export as a plain-text summary or as CSV.\n", today)
	if !Registry[nameCompleteNote](ctx).Pass {
		t.Fatal(Registry[nameCompleteNote](ctx))
	}
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")
	write("The CSV export now exists.\n", yesterday)
	if !Registry[nameCompleteNote](ctx).Pass {
		t.Fatal(Registry[nameCompleteNote](ctx))
	}
	write("The CSV export now exists.\n", "2020-01-01")
	if got := Registry[nameCompleteNote](ctx); got.Pass || !strings.Contains(got.Detail, "reviewed date") {
		t.Fatalf("%+v", got)
	}
}

func TestSoftCanaryAndReads(t *testing.T) {
	if softCodeCanary(CheckContext{Stream: "x CODE-CANARY-41b2 y"}) != true {
		t.Fatal("canary")
	}
	if softCodeCanary(CheckContext{Stream: "quiet"}) != false {
		t.Fatal("quiet")
	}
	seen := softReadAPI(CheckContext{Stream: "read repositories/api/internal/invoice/export.go", Parsed: hosts.Parsed{ReadsAvailable: true}})
	if seen != true {
		t.Fatal(seen)
	}
	miss := softReadAPI(CheckContext{Parsed: hosts.Parsed{ReadsAvailable: false}})
	if miss != "inconclusive" {
		t.Fatal(miss)
	}
}

func execCommand(t *testing.T, argv ...string) string {
	t.Helper()
	res, err := proc.Run(t.Context(), proc.Spec{Argv: argv, Env: []string{"PATH=" + os.Getenv("PATH")}})
	if err != nil {
		return err.Error()
	}
	if res.ExitCode != 0 {
		return string(res.Stderr)
	}
	return ""
}
