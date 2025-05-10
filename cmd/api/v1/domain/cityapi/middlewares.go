package cityapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/citycore"
	"github.com/hemozeetah/journi/pkg/muxer"
)

const cityKey = "city"

func setCity(ctx context.Context, city citycore.City) context.Context {
	return context.WithValue(ctx, cityKey, city)
}

func getCity(ctx context.Context) (citycore.City, error) {
	city, ok := ctx.Value(cityKey).(citycore.City)
	if !ok {
		return citycore.City{}, errors.New("city not found")
	}

	return city, nil
}

func (a *api) parseCityMW(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		cityID, err := uuid.Parse(r.PathValue("city_id"))
		if err != nil {
			return response.WriteError(w, http.StatusBadRequest, err)
		}

		city, err := a.core.QueryByID(ctx, cityID)
		if err != nil {
			if errors.Is(err, citycore.ErrNotFound) {
				return response.WriteError(w, http.StatusNotFound, citycore.ErrNotFound)
			}
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		ctx = setCity(ctx, city)
		return handler(ctx, w, r)
	}
}
