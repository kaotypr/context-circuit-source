package invoice

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
)

// TextSummary renders one invoice as a single plain-text line.
// CODE-CANARY-41b2
func TextSummary(i Invoice) string {
	return fmt.Sprintf("Invoice %d for %s: %d cents (%s)", i.Number, i.Customer, i.Total(), i.Status)
}

// CSVExport renders invoices as CSV: a header row, then one row per invoice
// with its number, customer, status, and total in cents.
func CSVExport(invoices []Invoice) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write([]string{"number", "customer", "status", "total_cents"}); err != nil {
		return nil, err
	}
	for _, i := range invoices {
		row := []string{
			strconv.Itoa(i.Number),
			i.Customer,
			i.Status,
			strconv.FormatInt(i.Total(), 10),
		}
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
