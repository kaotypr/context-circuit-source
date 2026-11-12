package invoice

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteCSVSingleIssuedInvoice(t *testing.T) {
	var buf bytes.Buffer
	inv := Invoice{
		Number:   7,
		Customer: "Ada",
		Status:   "issued",
		Lines:    []Line{{Description: "Work", Cents: 100}, {Description: "Tax", Cents: 25}},
	}
	if err := WriteCSV(&buf, []Invoice{inv}); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}
	want := "number,customer,status,total_cents\n7,Ada,issued,125\n"
	if got := buf.String(); got != want {
		t.Fatalf("csv = %q, want %q", got, want)
	}
}

func TestWriteCSVQuotesCustomer(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteCSV(&buf, []Invoice{{Number: 2, Customer: "Lovelace, Ada", Status: "draft"}}); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}
	if !strings.Contains(buf.String(), `2,"Lovelace, Ada",draft,0`) {
		t.Fatalf("csv %q does not quote the customer", buf.String())
	}
}
