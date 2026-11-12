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
	wantHeader := "number,customer,status,total"
	if !strings.Contains(got, wantHeader) || !strings.Contains(got, "7,Ada,issued,100") {
		t.Fatalf("csv %q", got)
	}
}
