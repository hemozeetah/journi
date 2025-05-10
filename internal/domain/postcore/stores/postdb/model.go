package postdb

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/postcore"
	"github.com/hemozeetah/journi/pkg/querybuilder"
	"github.com/lib/pq"
)

var fields = map[querybuilder.Field]string{
	postcore.ID:        "post_id",
	postcore.UserID:    "user_id",
	postcore.PlaceID:   "place_id",
	postcore.Caption:   "caption",
	postcore.Images:    "images",
	postcore.CreatedAt: "created_at",
	postcore.UpdatedAt: "updated_at",
}

type post struct {
	ID        uuid.UUID      `db:"post_id"`
	UserID    uuid.UUID      `db:"user_id"`
	PlaceID   uuid.UUID      `db:"place_id"`
	Caption   string         `db:"caption"`
	Images    pq.StringArray `db:"images"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func toPostDB(p postcore.Post) post {
	return post{
		ID:        p.ID,
		UserID:    p.UserID,
		PlaceID:   p.PlaceID,
		Caption:   p.Caption,
		Images:    p.Images,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toPostCore(p post) postcore.Post {
	return postcore.Post{
		ID:        p.ID,
		UserID:    p.UserID,
		PlaceID:   p.PlaceID,
		Caption:   p.Caption,
		Images:    p.Images,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
