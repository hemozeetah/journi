package placedb

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/placecore"
	"github.com/lib/pq"
)

type place struct {
	ID        uuid.UUID      `db:"place_id"`
	CityID    uuid.UUID      `db:"city_id"`
	Name      string         `db:"name"`
	Caption   string         `db:"caption"`
	Type      string         `db:"type"`
	Images    pq.StringArray `db:"images"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func toPlaceDB(p placecore.Place) place {
	return place{
		ID:        p.ID,
		CityID:    p.CityID,
		Name:      p.Name,
		Caption:   p.Caption,
		Images:    p.Images,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toPlaceCore(p place) placecore.Place {
	return placecore.Place{
		ID:        p.ID,
		Name:      p.Name,
		Caption:   p.Caption,
		Images:    p.Images,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
