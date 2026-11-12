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

// WriteCSV writes invoices with their number, customer, status, and total in cents.
func WriteCSV(w io.Writer, invoices []Invoice) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"number", "customer", "status", "total_cents"}); err != nil {
		return err
	}
	for _, i := range invoices {
		if err := writer.Write([]string{
			strconv.Itoa(i.Number),
			i.Customer,
			i.Status,
			strconv.FormatInt(i.Total(), 10),
		}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
