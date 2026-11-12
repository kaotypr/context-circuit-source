package invoice

import (
	"encoding/csv"
	"io"
	"strconv"
)

// WriteCSV writes one invoice as a CSV header and data row.
func WriteCSV(w io.Writer, i Invoice) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"number", "customer", "status", "total_cents"}); err != nil {
		return err
	}
	if err := writer.Write([]string{
		strconv.Itoa(i.Number),
		i.Customer,
		i.Status,
		strconv.FormatInt(i.Total(), 10),
	}); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}
