package programcore

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

const (
	ID querybuilder.Field = iota
	CompanyID
	Caption
	StartDate
	EndDate
	CreatedAt
	UpdatedAt
)

var DefaultOrderBy = querybuilder.NewOrderBy(CreatedAt, querybuilder.ASC)

type Program struct {
	ID        uuid.UUID
	CompanyID uuid.UUID
	Caption   string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateProgramParams struct {
	CompanyID uuid.UUID
	Caption   string
	StartDate time.Time
	EndDate   time.Time
}

type UpdateProgramParams struct {
	Caption   *string
	StartDate *time.Time
	EndDate   *time.Time
}
