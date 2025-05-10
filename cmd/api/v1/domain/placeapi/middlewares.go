package placeapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/placecore"
	"github.com/hemozeetah/journi/pkg/muxer"
)

const placeKey = "place"

func setPlace(ctx context.Context, place placecore.Place) context.Context {
	return context.WithValue(ctx, placeKey, place)
}

func getPlace(ctx context.Context) (placecore.Place, error) {
	place, ok := ctx.Value(placeKey).(placecore.Place)
	if !ok {
		return placecore.Place{}, errors.New("place not found")
	}

	return place, nil
}

func (a *api) parsePlaceMW(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		placeID, err := uuid.Parse(r.PathValue("place_id"))
		if err != nil {
			return response.WriteError(w, http.StatusBadRequest, err)
		}

		place, err := a.core.QueryByID(ctx, placeID)
		if err != nil {
			if errors.Is(err, placecore.ErrNotFound) {
				return response.WriteError(w, http.StatusNotFound, placecore.ErrNotFound)
			}
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		ctx = setPlace(ctx, place)
		return handler(ctx, w, r)
	}
}
