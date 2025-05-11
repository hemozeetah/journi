package journeydb

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/journeycore"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

var fields = map[querybuilder.Field]string{
	journeycore.ID:            "journey_id",
	journeycore.ProgramID:     "program_id",
	journeycore.PlaceID:       "place_id",
	journeycore.StartDateTime: "start_datetime",
	journeycore.EndDateTime:   "end_datetime",
	journeycore.CreatedAt:     "created_at",
	journeycore.UpdatedAt:     "updated_at",
}

type journey struct {
	ID            uuid.UUID `db:"journey_id"`
	ProgramID     uuid.UUID `db:"program_id"`
	PlaceID       uuid.UUID `db:"place_id"`
	StartDateTime time.Time `db:"start_datetime"`
	EndDateTime   time.Time `db:"end_datetime"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func toJourneyDB(p journeycore.Journey) journey {
	return journey{
		ID:            p.ID,
		ProgramID:     p.ProgramID,
		PlaceID:       p.PlaceID,
		StartDateTime: p.StartDateTime,
		EndDateTime:   p.EndDateTime,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func toJourneyCore(p journey) journeycore.Journey {
	return journeycore.Journey{
		ID:            p.ID,
		ProgramID:     p.ProgramID,
		PlaceID:       p.PlaceID,
		StartDateTime: p.StartDateTime,
		EndDateTime:   p.EndDateTime,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}
