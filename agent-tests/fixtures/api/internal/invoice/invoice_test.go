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
