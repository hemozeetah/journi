package journeycore

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

const (
	ID querybuilder.Field = iota
	ProgramID
	PlaceID
	StartDateTime
	EndDateTime
	CreatedAt
	UpdatedAt
)

var DefaultOrderBy = querybuilder.NewOrderBy(CreatedAt, querybuilder.ASC)

type Journey struct {
	ID            uuid.UUID
	ProgramID     uuid.UUID
	PlaceID       uuid.UUID
	StartDateTime time.Time
	EndDateTime   time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateJourneyParams struct {
	ProgramID     uuid.UUID
	PlaceID       uuid.UUID
	StartDateTime time.Time
	EndDateTime   time.Time
}

type UpdateJourneyParams struct {
	PlaceID       *uuid.UUID
	StartDateTime *time.Time
	EndDateTime   *time.Time
}
