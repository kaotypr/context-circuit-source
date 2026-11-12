package invoice

import (
	"bytes"
	"testing"
)

func TestWriteCSV(t *testing.T) {
	var got bytes.Buffer
	err := WriteCSV(&got, []Invoice{{
		Number:   7,
		Customer: "Ada",
		Status:   "issued",
		Lines:    []Line{{Description: "Work", Cents: 1250}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	want := "number,customer,status,total_cents\n7,Ada,issued,1250\n"
	if got.String() != want {
		t.Fatalf("CSV %q, want %q", got.String(), want)
	}
}
