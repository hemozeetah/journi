package journeycore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

var ErrNotFound = errors.New("journey not found")

type Storer interface {
	Create(ctx context.Context, journey Journey) error
	Update(ctx context.Context, journey Journey) error
	Delete(ctx context.Context, journey Journey) error
	QueryByID(ctx context.Context, journeyID uuid.UUID) (Journey, error)
	Query(ctx context.Context, query querybuilder.Query) ([]Journey, error)
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

func (c *Core) Create(ctx context.Context, p CreateJourneyParams) (Journey, error) {
	now := time.Now()

	journey := Journey{
		ID:            uuid.New(),
		ProgramID:     p.ProgramID,
		PlaceID:       p.PlaceID,
		StartDateTime: p.StartDateTime,
		EndDateTime:   p.EndDateTime,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := c.store.Create(ctx, journey); err != nil {
		return Journey{}, fmt.Errorf("create: %w", err)
	}

	return journey, nil
}

func (c *Core) Update(ctx context.Context, journey Journey, p UpdateJourneyParams) (Journey, error) {
	if p.PlaceID != nil {
		journey.PlaceID = *p.PlaceID
	}
	if p.StartDateTime != nil {
		journey.StartDateTime = *p.StartDateTime
	}
	if p.EndDateTime != nil {
		journey.EndDateTime = *p.EndDateTime
	}
	journey.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, journey); err != nil {
		return Journey{}, fmt.Errorf("update: %w", err)
	}

	return journey, nil
}

func (c *Core) Delete(ctx context.Context, journey Journey) error {
	if err := c.store.Delete(ctx, journey); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) QueryByID(ctx context.Context, journeyID uuid.UUID) (Journey, error) {
	journey, err := c.store.QueryByID(ctx, journeyID)
	if err != nil {
		return Journey{}, fmt.Errorf("querybyid[%s]: %w", journeyID, err)
	}

	return journey, nil
}

func (c *Core) Query(ctx context.Context, query querybuilder.Query) ([]Journey, error) {
	journeys, err := c.store.Query(ctx, query)
	if err != nil {
		return []Journey{}, fmt.Errorf("query: %w", err)
	}

	return journeys, nil
}
