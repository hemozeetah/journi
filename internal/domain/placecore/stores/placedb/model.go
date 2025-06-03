package placedb

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/placecore"
	"github.com/hemozeetah/journi/pkg/querybuilder"
	"github.com/lib/pq"
)

var fields = map[querybuilder.Field]string{
	placecore.ID:        "place_id",
	placecore.CityID:    "city_id",
	placecore.Name:      "name",
	placecore.Caption:   "caption",
	placecore.Type:      "type",
	placecore.Images:    "images",
	placecore.CreatedAt: "created_at",
	placecore.UpdatedAt: "updated_at",
}

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
		Type:      p.Type,
		Images:    p.Images,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toPlaceCore(p place) placecore.Place {
	return placecore.Place{
		ID:        p.ID,
		CityID:    p.CityID,
		Name:      p.Name,
		Caption:   p.Caption,
		Type:      p.Type,
		Images:    p.Images,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
