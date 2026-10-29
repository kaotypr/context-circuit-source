package invoice

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteCSVIssuedInvoice(t *testing.T) {
	var buf bytes.Buffer
	err := WriteCSV(&buf, []Invoice{{
		Number:   7,
		Customer: "Ada",
		Status:   "issued",
		Lines:    []Line{{Description: "Work", Cents: 100}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "number,customer,status,total") {
		t.Fatalf("header missing: %q", got)
	}
	if !strings.Contains(got, "7,Ada,issued,100") {
		t.Fatalf("row missing: %q", got)
	}
}
