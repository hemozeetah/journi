package querybuilder

// Operation represents the operation of the constraint
type Operation int

// Possible operations
const (
	EQ Operation = iota
	NEQ
	GT
	GTE
	LT
	LTE
)

// Constraint represents a constraint
type Constraint struct {
	Field     Field
	Operation Operation
	Value     any
}

// NewConstraint creates a new constraint
func NewConstraint(field Field, operation Operation, value any) Constraint {
	return Constraint{
		Field:     field,
		Operation: operation,
		Value:     value,
	}
}
