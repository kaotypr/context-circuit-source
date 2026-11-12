package invoice

import (
	"encoding/csv"
	"io"
	"strconv"
)

// WriteCSV writes invoice details as CSV, with totals in cents.
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
