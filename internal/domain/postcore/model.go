package postcore

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

const (
	ID querybuilder.Field = iota
	UserID
	PlaceID
	Caption
	Images
	CreatedAt
	UpdatedAt
)

type Post struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PlaceID   uuid.UUID
	Caption   string
	Images    []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatePostParams struct {
	UserID  uuid.UUID
	PlaceID uuid.UUID
	Caption string
	Images  []string
}

type UpdatePostParams struct {
	PlaceID *uuid.UUID
	Caption *string
	Images  *[]string
}
