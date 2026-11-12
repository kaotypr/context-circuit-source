package invoice

import (
	"bytes"
	"testing"
)

func TestWriteCSV(t *testing.T) {
	var output bytes.Buffer
	err := WriteCSV(&output, []Invoice{{
		Number:   7,
		Customer: "Ada",
		Status:   "issued",
		Lines:    []Line{{Description: "Work", Cents: 100}, {Description: "Tax", Cents: 25}},
	}})
	if err != nil {
		t.Fatalf("write CSV: %v", err)
	}

	const want = "number,customer,status,total_cents\n7,Ada,issued,125\n"
	if got := output.String(); got != want {
		t.Fatalf("CSV output %q, want %q", got, want)
	}
}
