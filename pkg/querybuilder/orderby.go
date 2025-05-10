package querybuilder

import (
	"fmt"
	"strings"
)

// Direction represents the direction of the order.
type Direction int

// Possible directions.
const (
	ASC Direction = iota
	DESC
)

// OrderBy represents the order of field with a direction.
type OrderBy struct {
	Field     Field
	Direction Direction
}

// NewOrderBy creates a new order by.
func NewOrderBy(field Field, direction Direction) OrderBy {
	return OrderBy{
		Field:     field,
		Direction: direction,
	}
}

// ParseOrderBy constructs a OrderBy value by parsing a string
func ParseOrderBy(orderBy string, fieldsMapping map[string]Field, defaultOrder OrderBy) (OrderBy, error) {
	if orderBy == "" {
		return defaultOrder, nil
	}

	direction := ASC

	if strings.HasPrefix(orderBy, "-") {
		orderBy = orderBy[1:]
		direction = DESC
	}

	field, exists := fieldsMapping[orderBy]
	if !exists {
		return OrderBy{}, fmt.Errorf("unknown order: %s", orderBy)
	}

	return NewOrderBy(field, direction), nil
}
