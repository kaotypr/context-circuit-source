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
	want := "number,customer,status,total_cents\n7,Ada,issued,100\n"
	if got != want {
		t.Fatalf("csv %q", got)
	}
	if !strings.Contains(got, "100") {
		t.Fatalf("total omitted: %q", got)
	}
}
