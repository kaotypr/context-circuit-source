package hosts

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func parseText(t *testing.T, h Host, text string) Parsed {
	t.Helper()
	p, err := h.Parse(strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func parseFile(t *testing.T, h Host, name string) Parsed {
	t.Helper()
	f, err := os.Open(filepath.Join("..", "..", "testdata", "streams", name))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	p, err := h.Parse(f)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func eqInt(t *testing.T, name string, got *int64, want int64) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s is nil, want %d", name, want)
	}
	if *got != want {
		t.Fatalf("%s = %d, want %d", name, *got, want)
	}
}

func near(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func i64(n int64) *int64 { return &n }
