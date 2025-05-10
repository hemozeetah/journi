package querybuilder

import (
	"fmt"
	"strconv"
)

// Page represents the page number and the number of rows per page.
type Page struct {
	Number int
	Rows   int
}

// NewPage creates a new page.
func NewPage(number int, rows int) Page {
	return Page{
		Number: number,
		Rows:   rows,
	}
}

// ParsePage parses the strings and validates the values are in reason.
func ParsePage(page string, rowsPerPage string) (Page, error) {
	p := NewPage(1, 10)

	if page != "" {
		number, err := strconv.Atoi(page)
		if err != nil {
			return Page{}, fmt.Errorf("page conversion: %w", err)
		}

		p.Number = number
	}

	if rowsPerPage != "" {
		rows, err := strconv.Atoi(rowsPerPage)
		if err != nil {
			return Page{}, fmt.Errorf("rows conversion: %w", err)
		}

		p.Rows = rows
	}

	if p.Number <= 0 {
		return Page{}, fmt.Errorf("page value too small, must be larger than 0")
	}

	if p.Rows <= 0 {
		return Page{}, fmt.Errorf("rows value too small, must be larger than 0")
	}

	if p.Rows > 100 {
		return Page{}, fmt.Errorf("rows value too large, must be less than 100")
	}

	return p, nil
}
