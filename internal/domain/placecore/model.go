package placecore

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

const (
	ID querybuilder.Field = iota
	CityID
	Name
	Caption
	Type
	Images
	CreatedAt
	UpdatedAt
)

var DefaultOrderBy = querybuilder.NewOrderBy(CreatedAt, querybuilder.ASC)

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
