package querybuilder

type Page struct {
	Number int
	Rows   int
}

func NewPage(number int, rows int) Page {
	return Page{
		Number: number,
		Rows:   rows,
	}
}
