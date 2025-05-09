package postcore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/logger"
)

var ErrNotFound = errors.New("post not found")

type Storer interface {
	Create(ctx context.Context, post Post) error
	Update(ctx context.Context, post Post) error
	Delete(ctx context.Context, post Post) error
	QueryByID(ctx context.Context, postID uuid.UUID) (Post, error)
	Query(ctx context.Context) ([]Post, error)
}

type Core struct {
	log   *logger.Logger
	store Storer
}

func New(log *logger.Logger, store Storer) *Core {
	return &Core{
		log:   log,
		store: store,
	}
}

func (c *Core) Create(ctx context.Context, p CreatePostParams) (Post, error) {
	now := time.Now()

	post := Post{
		ID:        uuid.New(),
		UserID:    p.UserID,
		PlaceID:   p.PlaceID,
		Caption:   p.Caption,
		Images:    p.Images,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.store.Create(ctx, post); err != nil {
		return Post{}, fmt.Errorf("create: %w", err)
	}

	return post, nil
}

func (c *Core) Update(ctx context.Context, post Post, p UpdatePostParams) (Post, error) {
	if p.PlaceID != nil {
		post.PlaceID = *p.PlaceID
	}
	if p.Caption != nil {
		post.Caption = *p.Caption
	}
	if p.Images != nil {
		post.Images = *p.Images
	}
	post.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, post); err != nil {
		return Post{}, fmt.Errorf("update: %w", err)
	}

	return post, nil
}

func (c *Core) Delete(ctx context.Context, post Post) error {
	if err := c.store.Delete(ctx, post); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) QueryByID(ctx context.Context, postID uuid.UUID) (Post, error) {
	post, err := c.store.QueryByID(ctx, postID)
	if err != nil {
		return Post{}, fmt.Errorf("querybyid[%s]: %w", postID, err)
	}

	return post, nil
}

func (c *Core) Query(ctx context.Context) ([]Post, error) {
	posts, err := c.store.Query(ctx)
	if err != nil {
		return []Post{}, fmt.Errorf("query: %w", err)
	}

	return posts, nil
}
