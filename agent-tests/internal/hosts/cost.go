package hosts

import "github.com/kaotypr/context-circuit-source/agent-tests/internal/record"

// Cost prices t against row.
//
// Uncached input is t.InputUncached, which the adapters set explicitly.
// When row.LongContextThreshold is greater than zero and InputTotal is above
// that threshold, the long-context rates apply to the whole request.
//
// Reasoning tokens are not added. They are assumed to be inside Output
// (to be confirmed in the pilot).
//
// ok is false when a needed count is nil, or a rate is zero while the
// matching token count is non-zero. reason explains that case.
func Cost(row PriceRow, t record.Tokens) (usd float64, ok bool, reason string) {
	if row.LongContextThreshold > 0 && t.InputTotal == nil {
		return 0, false, "input_total is null"
	}
	if t.InputUncached == nil {
		return 0, false, "input_uncached is null"
	}
	if t.CacheRead == nil {
		return 0, false, "cache_read is null"
	}
	if t.CacheWrite == nil {
		return 0, false, "cache_write is null"
	}
	if t.Output == nil {
		return 0, false, "output is null"
	}

	inRate := row.InputPerM
	cachedRate := row.CachedInputPerM
	writeRate := row.CacheWritePerM
	outRate := row.OutputPerM
	if row.LongContextThreshold > 0 && *t.InputTotal > row.LongContextThreshold {
		inRate = row.LongInputPerM
		cachedRate = row.LongCachedInputPerM
		writeRate = row.LongCacheWritePerM
		outRate = row.LongOutputPerM
	}
	if *t.InputUncached != 0 && inRate == 0 {
		return 0, false, "input price is zero"
	}
	if *t.CacheRead != 0 && cachedRate == 0 {
		return 0, false, "cached input price is zero"
	}
	if *t.CacheWrite != 0 && writeRate == 0 {
		return 0, false, "cache write price is zero"
	}
	if *t.Output != 0 && outRate == 0 {
		return 0, false, "output price is zero"
	}

	// Reasoning is intentionally absent from this sum.
	const million = 1000000.0
	usd = float64(*t.InputUncached)*inRate/million +
		float64(*t.CacheRead)*cachedRate/million +
		float64(*t.CacheWrite)*writeRate/million +
		float64(*t.Output)*outRate/million
	return usd, true, ""
}
