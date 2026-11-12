package invoice

import (
	"encoding/csv"
	"io"
	"strconv"
)

// WriteCSV writes invoices as CSV with their totals in cents.
func WriteCSV(w io.Writer, invoices []Invoice) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"number", "customer", "status", "total_cents"}); err != nil {
		return err
	}
	for _, invoice := range invoices {
		row := []string{
			strconv.Itoa(invoice.Number),
			invoice.Customer,
			invoice.Status,
			strconv.FormatInt(invoice.Total(), 10),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
