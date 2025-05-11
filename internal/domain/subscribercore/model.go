package subscribercore

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

const (
	UserID querybuilder.Field = iota
	ProgramID
	ReferenceID
  Accepted
	CreatedAt
	UpdatedAt
)

var DefaultOrderBy = querybuilder.NewOrderBy(CreatedAt, querybuilder.ASC)

type Subscriber struct {
	UserID      uuid.UUID
	ProgramID   uuid.UUID
	ReferenceID uuid.UUID
	Accepted    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateSubscriberParams struct {
	UserID    uuid.UUID
	ProgramID uuid.UUID
}

type UpdateSubscriberParams struct {
	Accepted bool
}
