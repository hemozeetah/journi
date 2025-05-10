package querybuilder

type Operation int

const (
	EQ Operation = iota
	NEQ
	GT
	GTE
	LT
	LTE
)

type Constraint struct {
	Field     Field
	Operation Operation
	Value     any
}

func NewConstraint(field Field, operation Operation, value any) Constraint {
	return Constraint{
		Field:     field,
		Operation: operation,
		Value:     value,
	}
}
