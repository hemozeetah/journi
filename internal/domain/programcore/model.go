package programcore

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

const (
	ID querybuilder.Field = iota
	CompanyID
	StartDate
	EndDate
	CreatedAt
	UpdatedAt
)

var DefaultOrderBy = querybuilder.NewOrderBy(CreatedAt, querybuilder.ASC)

type Program struct {
	ID        uuid.UUID
	CompanyID uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateProgramParams struct {
	CompanyID uuid.UUID
	StartDate time.Time
	EndDate   time.Time
}

type UpdateProgramParams struct {
	StartDate *time.Time
	EndDate   *time.Time
}
