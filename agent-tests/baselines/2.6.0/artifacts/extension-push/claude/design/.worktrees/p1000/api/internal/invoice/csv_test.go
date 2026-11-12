package invoice

import (
	"bytes"
	"testing"
)

func TestWriteCSVIssuedInvoice(t *testing.T) {
	inv := Invoice{
		Number:   7,
		Customer: "Ada",
		Status:   "issued",
		Lines:    []Line{{Description: "Work", Cents: 100}, {Description: "Tax", Cents: 25}},
	}
	var buf bytes.Buffer
	if err := WriteCSV(&buf, []Invoice{inv}); err != nil {
		t.Fatalf("write: %v", err)
	}
	want := "number,customer,status,total_cents\n7,Ada,issued,125\n"
	if got := buf.String(); got != want {
		t.Fatalf("csv %q, want %q", got, want)
	}
}
