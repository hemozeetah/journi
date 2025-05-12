package subscribercore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

var (
	ErrNotFound      = errors.New("subscriber not found")
	ErrAlreadyExists = errors.New("subscriber is already exists")
)

type Storer interface {
	Create(ctx context.Context, subscriber Subscriber) error
	Update(ctx context.Context, subscriber Subscriber) error
	Delete(ctx context.Context, subscriber Subscriber) error
	QueryByID(ctx context.Context, userID uuid.UUID, programID uuid.UUID) (Subscriber, error)
	QueryByReferenceID(ctx context.Context, referenceID uuid.UUID) (Subscriber, error)
	Query(ctx context.Context, query querybuilder.Query) ([]Subscriber, error)
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

func (c *Core) Create(ctx context.Context, p CreateSubscriberParams) (Subscriber, error) {
	now := time.Now()

	subscriber := Subscriber{
		UserID:      p.UserID,
		ProgramID:   p.ProgramID,
		ReferenceID: uuid.New(),
		Accepted:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := c.store.Create(ctx, subscriber); err != nil {
		return Subscriber{}, fmt.Errorf("create: %w", err)
	}

	return subscriber, nil
}

func (c *Core) Update(ctx context.Context, subscriber Subscriber, p UpdateSubscriberParams) (Subscriber, error) {
	subscriber.Accepted = p.Accepted
	subscriber.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, subscriber); err != nil {
		return Subscriber{}, fmt.Errorf("update: %w", err)
	}

	return subscriber, nil
}

func (c *Core) Delete(ctx context.Context, subscriber Subscriber) error {
	if err := c.store.Delete(ctx, subscriber); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) QueryByID(ctx context.Context, userID uuid.UUID, programID uuid.UUID) (Subscriber, error) {
	place, err := c.store.QueryByID(ctx, userID, programID)
	if err != nil {
		return Subscriber{}, fmt.Errorf("querybyid[%s,%s]: %w", userID, programID, err)
	}

	return place, nil
}

func (c *Core) QueryByReferenceID(ctx context.Context, referenceID uuid.UUID) (Subscriber, error) {
	place, err := c.store.QueryByReferenceID(ctx, referenceID)
	if err != nil {
		return Subscriber{}, fmt.Errorf("querybyreferenceid[%s]: %w", referenceID, err)
	}

	return place, nil
}

func (c *Core) Query(ctx context.Context, query querybuilder.Query) ([]Subscriber, error) {
	places, err := c.store.Query(ctx, query)
	if err != nil {
		return []Subscriber{}, fmt.Errorf("query: %w", err)
	}

	return places, nil
}
