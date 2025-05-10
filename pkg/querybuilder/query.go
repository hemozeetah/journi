package querybuilder

type Field int

type Query struct {
	Constraints []Constraint
	OrderBy     OrderBy
	Page        Page
}

func NewQuery(constraints []Constraint, orderBy OrderBy, page Page) Query {
	return Query{
		Constraints: constraints,
		OrderBy:     orderBy,
		Page:        page,
	}
}
