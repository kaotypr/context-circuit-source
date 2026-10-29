package invoice

import (
	"bytes"
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

func TestWriteCSV(t *testing.T) {
	var buf bytes.Buffer
	err := WriteCSV(&buf, []Invoice{{
		Number:   7,
		Customer: "Ada",
		Status:   "issued",
		Lines:    []Line{{Description: "Work", Cents: 100}, {Description: "Tax", Cents: 25}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	want := "number,customer,status,total_cents\n7,Ada,issued,125\n"
	if got := buf.String(); got != want {
		t.Fatalf("csv %q, want %q", got, want)
	}
}
