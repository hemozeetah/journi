package subscriberdb

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/subscribercore"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

var fields = map[querybuilder.Field]string{
	subscribercore.UserID:      "user_id",
	subscribercore.ProgramID:   "program_id",
	subscribercore.ReferenceID: "reference_id",
	subscribercore.Accepted:    "accepted",
	subscribercore.CreatedAt:   "created_at",
	subscribercore.UpdatedAt:   "updated_at",
}

type subscriber struct {
	UserID      uuid.UUID `db:"user_id"`
	ProgramID   uuid.UUID `db:"program_id"`
	ReferenceID uuid.UUID `db:"reference_id"`
	Accepted    bool      `db:"accepted"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func toSubscriberDB(p subscribercore.Subscriber) subscriber {
	return subscriber{
		UserID:      p.UserID,
		ProgramID:   p.ProgramID,
		ReferenceID: p.ReferenceID,
		Accepted:    p.Accepted,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func toSubscriberCore(p subscriber) subscribercore.Subscriber {
	return subscribercore.Subscriber{
		UserID:      p.UserID,
		ProgramID:   p.ProgramID,
		ReferenceID: p.ReferenceID,
		Accepted:    p.Accepted,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
