package placecore

import (
	"time"

	"github.com/google/uuid"
)

type Place struct {
	ID        uuid.UUID
	CityID    uuid.UUID
	Name      string
	Caption   string
	Type      string
	Images    []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatePlaceParams struct {
	CityID  uuid.UUID
	Name    string
	Caption string
	Type    string
	Images  []string
}

type UpdatePlaceParams struct {
	CityID  *uuid.UUID
	Name    *string
	Caption *string
	Type    *string
	Images  *[]string
}
