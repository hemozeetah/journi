package citycore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/logger"
)

var ErrNotFound = errors.New("city not found")

type Storer interface {
	Create(ctx context.Context, city City) error
	Update(ctx context.Context, city City) error
	Delete(ctx context.Context, city City) error
	QueryByID(ctx context.Context, cityID uuid.UUID) (City, error)
	Query(ctx context.Context) ([]City, error)
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

func (c *Core) Create(ctx context.Context, p CreateCityParams) (City, error) {
	now := time.Now()

	city := City{
		ID:        uuid.New(),
		Name:      p.Name,
		Caption:   p.Caption,
		Images:    p.Images,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.store.Create(ctx, city); err != nil {
		return City{}, fmt.Errorf("create: %w", err)
	}

	return city, nil
}

func (c *Core) Update(ctx context.Context, city City, p UpdateCityParams) (City, error) {
	if p.Name != nil {
		city.Name = *p.Name
	}
	if p.Caption != nil {
		city.Caption = *p.Caption
	}
	if p.Images != nil {
		city.Images = *p.Images
	}
	city.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, city); err != nil {
		return City{}, fmt.Errorf("update: %w", err)
	}

	return city, nil
}

func (c *Core) Delete(ctx context.Context, city City) error {
	if err := c.store.Delete(ctx, city); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) QueryByID(ctx context.Context, cityID uuid.UUID) (City, error) {
	city, err := c.store.QueryByID(ctx, cityID)
	if err != nil {
		return City{}, fmt.Errorf("querybyid[%s]: %w", cityID, err)
	}

	return city, nil
}

func (c *Core) Query(ctx context.Context) ([]City, error) {
	cities, err := c.store.Query(ctx)
	if err != nil {
		return []City{}, fmt.Errorf("query: %w", err)
	}

	return cities, nil
}
