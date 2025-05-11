package placecore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/logger"
)

var ErrNotFound = errors.New("place not found")

type Storer interface {
	Create(ctx context.Context, place Place) error
	Update(ctx context.Context, place Place) error
	Delete(ctx context.Context, place Place) error
	QueryByID(ctx context.Context, placeID uuid.UUID) (Place, error)
	Query(ctx context.Context) ([]Place, error)
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

func (c *Core) Create(ctx context.Context, p CreatePlaceParams) (Place, error) {
	now := time.Now()

	place := Place{
		ID:        uuid.New(),
		CityID:    p.CityID,
		Name:      p.Name,
		Caption:   p.Caption,
		Type:      p.Type,
		Images:    p.Images,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.store.Create(ctx, place); err != nil {
		return Place{}, fmt.Errorf("create: %w", err)
	}

	return place, nil
}

func (c *Core) Update(ctx context.Context, place Place, p UpdatePlaceParams) (Place, error) {
	if p.CityID != nil {
		place.CityID = *p.CityID
	}
	if p.Name != nil {
		place.Name = *p.Name
	}
	if p.Caption != nil {
		place.Caption = *p.Caption
	}
	if p.Type != nil {
		place.Type = *p.Type
	}
	if p.Images != nil {
		place.Images = *p.Images
	}
	place.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, place); err != nil {
		return Place{}, fmt.Errorf("update: %w", err)
	}

	return place, nil
}

func (c *Core) Delete(ctx context.Context, place Place) error {
	if err := c.store.Delete(ctx, place); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) QueryByID(ctx context.Context, placeID uuid.UUID) (Place, error) {
	place, err := c.store.QueryByID(ctx, placeID)
	if err != nil {
		return Place{}, fmt.Errorf("querybyid[%s]: %w", placeID, err)
	}

	return place, nil
}

func (c *Core) Query(ctx context.Context) ([]Place, error) {
	places, err := c.store.Query(ctx)
	if err != nil {
		return []Place{}, fmt.Errorf("query: %w", err)
	}

	return places, nil
}
