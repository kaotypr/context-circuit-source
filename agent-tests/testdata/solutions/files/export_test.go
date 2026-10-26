package invoice

import (
	"strings"
	"testing"
)

func TestTextSummary(t *testing.T) {
	got := TextSummary(Invoice{
		Number:   7,
		Customer: "Ada",
		Status:   "issued",
		Lines:    []Line{{Description: "Work", Cents: 100}},
	})
	if !strings.Contains(got, "Ada") || !strings.Contains(got, "100") || !strings.Contains(got, "issued") {
		t.Fatalf("summary %q", got)
	}
}

func TestCSV(t *testing.T) {
	var buf strings.Builder
	err := CSV(&buf, []Invoice{{
		Number:   7,
		Customer: "Ada",
		Status:   "issued",
		Lines:    []Line{{Description: "Work", Cents: 100}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "number,customer,status,cents") || !strings.Contains(got, "7,Ada,issued,100") {
		t.Fatalf("csv %q", got)
	}
}
