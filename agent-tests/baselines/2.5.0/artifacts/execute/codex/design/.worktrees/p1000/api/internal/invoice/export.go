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

// WriteCSV writes invoices as CSV rows with totals expressed in cents.
func WriteCSV(w io.Writer, invoices []Invoice) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"number", "customer", "status", "total_cents"}); err != nil {
		return err
	}
	for _, invoice := range invoices {
		if err := writer.Write([]string{
			strconv.Itoa(invoice.Number),
			invoice.Customer,
			invoice.Status,
			strconv.FormatInt(invoice.Total(), 10),
		}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
