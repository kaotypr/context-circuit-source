package invoice

import "testing"

func TestTotalSumsLines(t *testing.T) {
	inv := Invoice{
		Number:   1,
		Customer: "Ada",
		Status:   "draft",
		Lines: []Line{
			{Description: "Work", Cents: 100},
			{Description: "Tax", Cents: 25},
		},
	}
	if got := inv.Total(); got != 125 {
		t.Fatalf("total %d", got)
	}
}

func TestTotalEmpty(t *testing.T) {
	if got := (Invoice{}).Total(); got != 0 {
		t.Fatalf("total %d", got)
	}
}

func TestVoidOnlyDraft(t *testing.T) {
	draft := Invoice{Status: "draft"}
	if err := draft.Void(); err != nil || draft.Status != "void" {
		t.Fatalf("void draft: %v, status %s", err, draft.Status)
	}
	issued := Invoice{Status: "issued"}
	if err := issued.Void(); err == nil {
		t.Fatal("voided an issued invoice")
	}
}

func TestCreditOnlyIssued(t *testing.T) {
	issued := Invoice{Number: 7, Status: "issued", Lines: []Line{{Cents: 300}}}
	note, err := issued.Credit()
	if err != nil || note.Invoice != 7 || note.Cents != 300 {
		t.Fatalf("credit issued: %v, %+v", err, note)
	}
	if _, err := (Invoice{Status: "draft"}).Credit(); err == nil {
		t.Fatal("credited a draft invoice")
	}
}
