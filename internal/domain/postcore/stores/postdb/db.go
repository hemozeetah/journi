package postdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/postcore"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/postgres"
	"github.com/hemozeetah/journi/pkg/querybuilder"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log *logger.Logger
	db  *sqlx.DB
}

func New(log *logger.Logger, db *sqlx.DB) *DB {
	return &DB{
		log: log,
		db:  db,
	}
}

func (db *DB) Create(ctx context.Context, post postcore.Post) error {
	const q = `
INSERT INTO posts
  (post_id, user_id, place_id, caption, images, created_at, updated_at)
VALUES
  (:post_id, :user_id, :place_id, :caption, :images, :created_at, :updated_at)`

	if err := postgres.ExecContext(ctx, db.db, q, toPostDB(post)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Update(ctx context.Context, post postcore.Post) error {
	const q = `
UPDATE posts
SET
  place_id = :place_id,
  caption = :caption,
  images = :images,
  updated_at = :updated_at
WHERE
  post_id = :post_id`

	if err := postgres.ExecContext(ctx, db.db, q, toPostDB(post)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Delete(ctx context.Context, post postcore.Post) error {
	const q = `
DELETE FROM posts
WHERE
  post_id = :post_id`

	if err := postgres.ExecContext(ctx, db.db, q, toPostDB(post)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) QueryByID(ctx context.Context, postID uuid.UUID) (postcore.Post, error) {
	const q = `
SELECT
  post_id,
  user_id,
  place_id,
  caption,
  images,
  created_at,
  updated_at
FROM
  posts
WHERE
  post_id = :post_id`

	data := struct {
		ID string `db:"post_id"`
	}{
		ID: postID.String(),
	}

	var p post
	err := postgres.QueryOneContext(ctx, db.db, q, data, &p)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return postcore.Post{}, fmt.Errorf("queryonecontext: %w", postcore.ErrNotFound)
		}
		return postcore.Post{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toPostCore(p), nil
}

func (db *DB) Query(ctx context.Context, query querybuilder.Query) ([]postcore.Post, error) {
	const q = `
SELECT
  post_id,
  user_id,
  place_id,
  caption,
  images,
  created_at,
  updated_at
FROM
  posts`

	data := map[string]any{}

	qq := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		q,
		postgres.WhereCluase(fields, query.Constraints),
		postgres.OrderByCluase(fields, query.OrderBy),
		postgres.OffsetCluase(fields, query.Page),
	)

	var ps []post
	err := postgres.QueryContext(ctx, db.db, qq, data, &ps)
	if err != nil {
		return []postcore.Post{}, fmt.Errorf("querycontext: %w", err)
	}

	posts := make([]postcore.Post, len(ps))
	for i, p := range ps {
		posts[i] = toPostCore(p)
	}

	return posts, nil
}
