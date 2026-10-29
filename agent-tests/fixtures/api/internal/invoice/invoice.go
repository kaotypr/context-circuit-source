package invoice

import "errors"

// Invoice is one bill. Total is the sum of its line amounts.
type Invoice struct {
	Number   int
	Customer string
	Lines    []Line
	Status   string
}

// Line is one row on an invoice. Cents is the amount in minor units.
type Line struct {
	Description string
	Cents       int64
}

// Total returns the sum of line amounts in cents.
func (i Invoice) Total() int64 {
	var sum int64
	for _, line := range i.Lines {
		sum += line.Cents
	}
	return sum
}

// Void cancels a draft invoice. An issued invoice is never voided; it is
// credited with a credit note instead.
func (i *Invoice) Void() error {
	if i.Status != "draft" {
		return errors.New("only a draft invoice can be voided; credit an issued one")
	}
	i.Status = "void"
	return nil
}

// CreditNote cancels an issued invoice. Cents is the amount credited.
type CreditNote struct {
	Invoice int
	Cents   int64
}

// Credit returns a credit note for the full total of an issued invoice.
func (i Invoice) Credit() (CreditNote, error) {
	if i.Status != "issued" {
		return CreditNote{}, errors.New("only an issued invoice is credited; void a draft")
	}
	return CreditNote{Invoice: i.Number, Cents: i.Total()}, nil
}
