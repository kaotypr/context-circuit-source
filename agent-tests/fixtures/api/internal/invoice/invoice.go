package invoice

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
