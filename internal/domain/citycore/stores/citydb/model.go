package citydb

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/citycore"
	"github.com/lib/pq"
)

type city struct {
	ID        uuid.UUID      `db:"city_id"`
	Name      string         `db:"name"`
	Caption   string         `db:"caption"`
	Images    pq.StringArray `db:"images"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func toCityDB(c citycore.City) city {
	return city{
		ID:        c.ID,
		Name:      c.Name,
		Caption:   c.Caption,
		Images:    c.Images,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCityCore(c city) citycore.City {
	return citycore.City{
		ID:        c.ID,
		Name:      c.Name,
		Caption:   c.Caption,
		Images:    c.Images,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
