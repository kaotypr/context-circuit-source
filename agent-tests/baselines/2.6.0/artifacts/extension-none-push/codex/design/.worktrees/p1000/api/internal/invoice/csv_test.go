package invoice

import (
	"bytes"
	"testing"
)

func TestWriteCSV(t *testing.T) {
	var got bytes.Buffer
	err := WriteCSV(&got, Invoice{
		Number:   7,
		Customer: "Ada",
		Status:   "issued",
		Lines:    []Line{{Description: "Work", Cents: 100}, {Description: "Tax", Cents: 25}},
	})
	if err != nil {
		t.Fatalf("write CSV: %v", err)
	}
	want := "number,customer,status,total_cents\n7,Ada,issued,125\n"
	if got.String() != want {
		t.Fatalf("CSV = %q, want %q", got.String(), want)
	}
}
