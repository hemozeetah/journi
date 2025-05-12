package journeyapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/journeycore"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/muxer"
)

const journeyKey = "journey"

func setJourney(ctx context.Context, journey journeycore.Journey) context.Context {
	return context.WithValue(ctx, journeyKey, journey)
}

func getJourney(ctx context.Context) (journeycore.Journey, error) {
	journey, ok := ctx.Value(journeyKey).(journeycore.Journey)
	if !ok {
		return journeycore.Journey{}, errors.New("journey not found")
	}

	return journey, nil
}

func (a *api) parseJourneyMW(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		journeyID, err := uuid.Parse(r.PathValue("journey_id"))
		if err != nil {
			return response.WriteError(w, http.StatusBadRequest, err)
		}

		journey, err := a.core.QueryByID(ctx, journeyID)
		if err != nil {
			if errors.Is(err, journeycore.ErrNotFound) {
				return response.WriteError(w, http.StatusNotFound, journeycore.ErrNotFound)
			}
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		ctx = setJourney(ctx, journey)
		return handler(ctx, w, r)
	}
}

func (a *api) adminOrOwner(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		claims, err := jwtauth.GetClaims(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		journey, err := getJourney(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		program, err := a.programCore.QueryByID(ctx, journey.ProgramID)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		if claims.Role == usercore.RoleAdmin || claims.ID == program.CompanyID {
			return handler(ctx, w, r)
		}

		return response.WriteError(w, http.StatusUnauthorized, errors.New("not admin or owner"))
	}
}
