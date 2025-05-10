// package querybuilder contains filtering, ordering, and pagination.
package querybuilder

// Field represents a query field.
type Field int

// Query represents the query.
type Query struct {
	Constraints []Constraint
	OrderBy     OrderBy
	Page        Page
}

// NewQuery creates a new query
func NewQuery(constraints []Constraint, orderBy OrderBy, page Page) Query {
	return Query{
		Constraints: constraints,
		OrderBy:     orderBy,
		Page:        page,
	}
}
