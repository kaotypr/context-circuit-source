package hosts

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

// HostConfig is one entry under hosts: in hosts.yaml.
type HostConfig struct {
	Name           string   `yaml:"-"`
	Command        string   `yaml:"command"`
	Model          string   `yaml:"model"`
	Effort         string   `yaml:"effort"`
	HostVersion    string   `yaml:"host_version"`
	Provisional    bool     `yaml:"provisional"`
	TimeoutSeconds int      `yaml:"timeout_seconds"`
	MaxBudgetUSD   float64  `yaml:"max_budget_usd"`
	PriceRow       string   `yaml:"price_row"`
	AuthEnv        []string `yaml:"auth_env"`
	// Auth is "env" or "subscription". Empty in the file means env.
	// subscription is accepted only for the codex host. The runner copies
	// a login file only when the command flag names that host.
	Auth string `yaml:"auth"`
	// SubscriptionAuth is set by the runner from --subscription-auth.
	// It is not loaded from hosts.yaml.
	SubscriptionAuth bool `yaml:"-"`
}

// PriceRow is one entry under prices: in prices.yaml.
type PriceRow struct {
	InputPerM            float64 `yaml:"input_per_m"`
	CachedInputPerM      float64 `yaml:"cached_input_per_m"`
	CacheWritePerM       float64 `yaml:"cache_write_per_m"`
	OutputPerM           float64 `yaml:"output_per_m"`
	LongContextThreshold int64   `yaml:"long_context_threshold"`
	LongInputPerM        float64 `yaml:"long_input_per_m"`
	LongCachedInputPerM  float64 `yaml:"long_cached_input_per_m"`
	LongCacheWritePerM   float64 `yaml:"long_cache_write_per_m"`
	LongOutputPerM       float64 `yaml:"long_output_per_m"`
	Source               string  `yaml:"source"`
	Checked              string  `yaml:"checked"`
	Notes                string  `yaml:"notes"`
}

// LoadHosts reads hosts.yaml. Unknown fields are an error. A host with an
// empty model is an error.
func LoadHosts(path string) (map[string]HostConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc struct {
		Hosts map[string]HostConfig `yaml:"hosts"`
	}
	if err := yaml.UnmarshalWithOptions(b, &doc, yaml.Strict()); err != nil {
		return nil, err
	}
	if doc.Hosts == nil {
		return map[string]HostConfig{}, nil
	}
	for name, cfg := range doc.Hosts {
		if cfg.Model == "" {
			return nil, fmt.Errorf("host %q: model not pinned", name)
		}
		switch cfg.Auth {
		case "", "env":
			cfg.Auth = "env"
		case "subscription":
			if name != "codex" {
				return nil, fmt.Errorf("host %q: auth subscription is only supported for codex", name)
			}
		default:
			return nil, fmt.Errorf("host %q: unknown auth %q", name, cfg.Auth)
		}
		cfg.Name = name
		doc.Hosts[name] = cfg
	}
	return doc.Hosts, nil
}

// LoadPrices reads prices.yaml. Every row needs a source and a YYYY-MM-DD
// checked date. Unknown fields are an error.
func LoadPrices(path string) (map[string]PriceRow, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc struct {
		Prices map[string]PriceRow `yaml:"prices"`
	}
	if err := yaml.UnmarshalWithOptions(b, &doc, yaml.Strict()); err != nil {
		return nil, err
	}
	if doc.Prices == nil {
		return map[string]PriceRow{}, nil
	}
	for name, row := range doc.Prices {
		if row.Source == "" || row.Checked == "" {
			return nil, fmt.Errorf("price %q: source and checked are required", name)
		}
		if _, err := time.Parse("2006-01-02", row.Checked); err != nil {
			return nil, fmt.Errorf("price %q: checked must be YYYY-MM-DD", name)
		}
	}
	return doc.Prices, nil
}
