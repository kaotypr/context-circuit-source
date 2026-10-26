package hosts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadHosts(t *testing.T) {
	hosts, err := LoadHosts(filepath.Join("..", "..", "hosts.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	claude := hosts["claude"]
	if claude.Name != "claude" || claude.Model != "claude-opus-5-5" || !claude.Provisional {
		t.Fatalf("%+v", claude)
	}
	if claude.Command != "claude" || claude.HostVersion != "2.1.304" || claude.MaxBudgetUSD != 3 {
		t.Fatalf("%+v", claude)
	}
	if !reflectAuth(claude.AuthEnv, []string{"ANTHROPIC_API_KEY", "CLAUDE_CODE_OAUTH_TOKEN"}) {
		t.Fatalf("auth %v", claude.AuthEnv)
	}
	codex := hosts["codex"]
	if codex.Model != "gpt-6-sol" || codex.Effort != "medium" || codex.PriceRow != "gpt-6-sol" || codex.HostVersion != "0.162.1" {
		t.Fatalf("%+v", codex)
	}
	cursor := hosts["cursor"]
	if cursor.Model != "auto" || cursor.PriceRow != "" || cursor.HostVersion != "2026.10.22-27dcf5f" {
		t.Fatalf("%+v", cursor)
	}
}

func TestLoadHostsUnknownField(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hosts.yaml")
	body := "hosts:\n  claude:\n    command: claude\n    model: m\n    not_a_field: true\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadHosts(path)
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestLoadHostsModelNotPinned(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hosts.yaml")
	body := "hosts:\n  claude:\n    command: claude\n    model: \"\"\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadHosts(path)
	if err == nil || !strings.Contains(err.Error(), "model not pinned") {
		t.Fatalf("err %v", err)
	}
}

func TestLoadPrices(t *testing.T) {
	prices, err := LoadPrices(filepath.Join("..", "..", "prices.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	row := prices["gpt-6-sol"]
	if row.Source != "https://developers.openai.com/api/docs/pricing" || row.Checked != "2026-10-26" {
		t.Fatalf("%+v", row)
	}
	if row.InputPerM != 2 || row.CachedInputPerM != 0.2 || row.CacheWritePerM != 2.5 || row.OutputPerM != 10 {
		t.Fatalf("%+v", row)
	}
	if row.LongContextThreshold != 272000 || row.LongInputPerM != 4 || row.LongOutputPerM != 15 {
		t.Fatalf("%+v", row)
	}
}

func TestLoadPricesRequiresSource(t *testing.T) {
	path := filepath.Join(t.TempDir(), "prices.yaml")
	body := "prices:\n  x:\n    input_per_m: 1\n    checked: \"2026-10-26\"\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadPrices(path)
	if err == nil || !strings.Contains(err.Error(), "source and checked") {
		t.Fatalf("err %v", err)
	}
}

func reflectAuth(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
