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
	var got bytes.Buffer
	err := WriteCSV(&got, []Invoice{{
		Number:   7,
		Customer: "Ada",
		Status:   "issued",
		Lines:    []Line{{Description: "Work", Cents: 100}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	want := "number,customer,status,total_cents\n7,Ada,issued,100\n"
	if got.String() != want {
		t.Fatalf("CSV = %q, want %q", got.String(), want)
	}
}
