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
	if claude.Name != "claude" || claude.Model != "claude-sonnet-5-5" || !claude.Provisional {
		t.Fatalf("%+v", claude)
	}
	if claude.Command != "claude" || claude.HostVersion != "2.1.304" || claude.MaxBudgetUSD != 3 {
		t.Fatalf("%+v", claude)
	}
	if !reflectAuth(claude.AuthEnv, []string{"ANTHROPIC_API_KEY", "CLAUDE_CODE_OAUTH_TOKEN"}) {
		t.Fatalf("auth %v", claude.AuthEnv)
	}
	codex := hosts["codex"]
	if codex.Model != "gpt-5.6-luna" || codex.Effort != "medium" || codex.PriceRow != "" || codex.HostVersion != "0.162.1" {
		t.Fatalf("%+v", codex)
	}
	cursor := hosts["cursor"]
	if cursor.Model != "composer-2.5" || cursor.PriceRow != "" || cursor.HostVersion != "2026.10.22-27dcf5f" {
		t.Fatalf("%+v", cursor)
	}
	if claude.Auth != "env" || codex.Auth != "env" || cursor.Auth != "env" {
		t.Fatalf("auth mode claude %q codex %q cursor %q", claude.Auth, codex.Auth, cursor.Auth)
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

func TestLoadHostsAuth(t *testing.T) {
	t.Run("subscription codex", func(t *testing.T) {
		cfg := loadHostYAML(t, "hosts:\n  codex:\n    command: codex\n    model: m\n    auth: subscription\n")
		if cfg["codex"].Auth != "subscription" {
			t.Fatalf("%+v", cfg["codex"])
		}
	})
	t.Run("default env", func(t *testing.T) {
		cfg := loadHostYAML(t, "hosts:\n  codex:\n    command: codex\n    model: m\n")
		if cfg["codex"].Auth != "env" {
			t.Fatalf("%+v", cfg["codex"])
		}
	})
	t.Run("other host", func(t *testing.T) {
		_, err := loadHostYAMLErr(t, "hosts:\n  claude:\n    command: claude\n    model: m\n    auth: subscription\n")
		if err == nil || !strings.Contains(err.Error(), "only supported for codex") {
			t.Fatalf("err %v", err)
		}
	})
	t.Run("unknown", func(t *testing.T) {
		_, err := loadHostYAMLErr(t, "hosts:\n  codex:\n    command: codex\n    model: m\n    auth: token\n")
		if err == nil || !strings.Contains(err.Error(), "unknown auth") {
			t.Fatalf("err %v", err)
		}
	})
}

func loadHostYAML(t *testing.T, body string) map[string]HostConfig {
	t.Helper()
	cfg, err := loadHostYAMLErr(t, body)
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func loadHostYAMLErr(t *testing.T, body string) (map[string]HostConfig, error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "hosts.yaml")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return LoadHosts(path)
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
