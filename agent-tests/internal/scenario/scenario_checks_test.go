package scenario

import (
	"os"
	"path/filepath"
	"testing"

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
