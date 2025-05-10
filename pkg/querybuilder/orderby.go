package querybuilder

type Direction int

const (
	ASC Direction = iota
	DESC
)

type OrderBy struct {
	Field     Field
	Direction Direction
}

func NewOrderBy(field Field, direction Direction) OrderBy {
	return OrderBy{
		Field:     field,
		Direction: direction,
	}
}
