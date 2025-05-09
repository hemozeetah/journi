package citycore

import (
	"time"

	"github.com/google/uuid"
)

type City struct {
	ID        uuid.UUID
	Name      string
	Caption   string
	Images    []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateCityParams struct {
	Name    string
	Caption string
	Images  []string
}

type UpdateCityParams struct {
	Name    *string
	Caption *string
	Images  *[]string
}
