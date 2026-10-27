package hosts

import (
	"path/filepath"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

func TestCostTiersAndNil(t *testing.T) {
	row := PriceRow{
		InputPerM:            2,
		CachedInputPerM:      0.2,
		CacheWritePerM:       2.5,
		OutputPerM:           10,
		LongContextThreshold: 272000,
		LongInputPerM:        4,
		LongCachedInputPerM:  0.4,
		LongCacheWritePerM:   5,
		LongOutputPerM:       15,
	}
	short := record.Tokens{
		InputTotal:    i64(1250),
		InputUncached: i64(1000),
		CacheRead:     i64(200),
		CacheWrite:    i64(50),
		Output:        i64(10),
		Reasoning:     i64(999999),
	}
	got, ok, reason := Cost(row, short)
	if !ok {
		t.Fatal(reason)
	}
	// 1000*2 + 200*0.20 + 50*2.50 + 10*10, per million. Reasoning is not added.
	want := 1000*2.0/1e6 + 200*0.20/1e6 + 50*2.50/1e6 + 10*10.0/1e6
	near(t, got, want)
	withoutReasoning := short
	withoutReasoning.Reasoning = nil
	again, ok, reason := Cost(row, withoutReasoning)
	if !ok {
		t.Fatal(reason)
	}
	near(t, again, got)

	atThreshold := record.Tokens{
		InputTotal:    i64(272000),
		InputUncached: i64(272000),
		CacheRead:     i64(0),
		CacheWrite:    i64(0),
		Output:        i64(0),
	}
	got, ok, reason = Cost(row, atThreshold)
	if !ok {
		t.Fatal(reason)
	}
	near(t, got, 272000*2.0/1e6)

	long := record.Tokens{
		InputTotal:    i64(272001),
		InputUncached: i64(272001),
		CacheRead:     i64(0),
		CacheWrite:    i64(0),
		Output:        i64(10),
	}
	got, ok, reason = Cost(row, long)
	if !ok {
		t.Fatal(reason)
	}
	// Whole request uses the long rates: 272001*4 + 10*15, per million.
	near(t, got, float64(272001*4+10*15)/1e6)

	missing := long
	missing.CacheRead = nil
	if _, ok, reason = Cost(row, missing); ok || reason != "cache_read is null" {
		t.Fatalf("ok %v reason %q", ok, reason)
	}

	zeroPrice := row
	zeroPrice.LongContextThreshold = 0
	zeroPrice.InputPerM = 0
	priced, ok, reason := Cost(zeroPrice, record.Tokens{
		InputUncached: i64(5),
		CacheRead:     i64(0),
		CacheWrite:    i64(0),
		Output:        i64(0),
	})
	if ok || priced != 0 || reason != "input price is zero" {
		t.Fatalf("ok %v usd %v reason %q", ok, priced, reason)
	}
}

func TestCostCodexRecordedStream(t *testing.T) {
	prices, err := LoadPrices(filepath.Join("..", "..", "prices.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	row := prices["gpt-5.6-luna"]
	h, err := ForName("codex")
	if err != nil {
		t.Fatal(err)
	}
	p := parseFile(t, h, "codex-exec.jsonl")
	// 15974 uncached is under the long-context threshold.
	// 15974 * $0.20/M + 0 * $0.02/M + 0 * $0.25/M + 5 * $1.20/M = 0.0032008
	got, ok, reason := Cost(row, p.Tokens)
	if !ok {
		t.Fatal(reason)
	}
	near(t, got, 0.0032008)
}

func TestCostCursorRecordedStream(t *testing.T) {
	prices, err := LoadPrices(filepath.Join("..", "..", "prices.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	h, err := ForName("cursor")
	if err != nil {
		t.Fatal(err)
	}
	p := parseFile(t, h, "cursor-stream.jsonl")
	got, ok, reason := Cost(prices["composer-2.5"], p.Tokens)
	if !ok {
		t.Fatal(reason)
	}
	// Cache writes are charged at the input rate, the row's stated estimate.
	want := (float64(*p.Tokens.InputUncached)*0.5 + float64(*p.Tokens.CacheRead)*0.2 +
		float64(*p.Tokens.CacheWrite)*0.5 + float64(*p.Tokens.Output)*2.5) / 1e6
	near(t, got, want)
	if *p.Tokens.CacheWrite == 0 {
		t.Fatal("the recorded Cursor stream no longer writes cache; this test no longer covers the estimated rate")
	}
}
