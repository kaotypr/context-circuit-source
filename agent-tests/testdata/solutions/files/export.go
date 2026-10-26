package invoice

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// TextSummary renders one invoice as a single plain-text line.
// CODE-CANARY-41b2
func TextSummary(i Invoice) string {
	return fmt.Sprintf("Invoice %d for %s: %d cents (%s)", i.Number, i.Customer, i.Total(), i.Status)
}

// CSV writes invoices as CSV text. The header is number, customer, status,
// and cents, and each invoice is one row.
func CSV(w io.Writer, invoices []Invoice) error {
	out := csv.NewWriter(w)
	if err := out.Write([]string{"number", "customer", "status", "cents"}); err != nil {
		return err
	}
	for _, item := range invoices {
		row := []string{
			strconv.Itoa(item.Number),
			item.Customer,
			item.Status,
			strconv.FormatInt(item.Total(), 10),
		}
		if err := out.Write(row); err != nil {
			return err
		}
	}
	out.Flush()
	return out.Error()
}
