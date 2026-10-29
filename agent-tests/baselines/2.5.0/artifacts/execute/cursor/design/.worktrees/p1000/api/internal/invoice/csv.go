package invoice

import (
	"encoding/csv"
	"io"
	"strconv"
)

// WriteCSV writes a header and one row per invoice: number, customer, status,
// and total in cents.
func WriteCSV(w io.Writer, invoices []Invoice) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"number", "customer", "status", "total"}); err != nil {
		return err
	}
	for _, inv := range invoices {
		row := []string{
			strconv.Itoa(inv.Number),
			inv.Customer,
			inv.Status,
			strconv.FormatInt(inv.Total(), 10),
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
