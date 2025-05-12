package subscriberapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/subscribercore"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/muxer"
)

const subscriberKey = "subscriber"

func setSubscriber(ctx context.Context, subscriber subscribercore.Subscriber) context.Context {
	return context.WithValue(ctx, subscriberKey, subscriber)
}

func getSubscriber(ctx context.Context) (subscribercore.Subscriber, error) {
	subscriber, ok := ctx.Value(subscriberKey).(subscribercore.Subscriber)
	if !ok {
		return subscribercore.Subscriber{}, errors.New("subscriber not found")
	}

	return subscriber, nil
}

func (a *api) parseSubscriberMW(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		referenceID, err := uuid.Parse(r.PathValue("reference_id"))
		if err != nil {
			return response.WriteError(w, http.StatusBadRequest, err)
		}

		subscriber, err := a.core.QueryByReferenceID(ctx, referenceID)
		if err != nil {
			if errors.Is(err, subscribercore.ErrNotFound) {
				return response.WriteError(w, http.StatusNotFound, subscribercore.ErrNotFound)
			}
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		ctx = setSubscriber(ctx, subscriber)
		return handler(ctx, w, r)
	}
}

func (a *api) adminOrOwnerOrCompanyOwner(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		claims, err := jwtauth.GetClaims(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		subscriber, err := getSubscriber(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		program, err := a.programCore.QueryByID(ctx, subscriber.ProgramID)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		if claims.Role == usercore.RoleAdmin || claims.ID == subscriber.UserID || claims.ID == program.CompanyID {
			return handler(ctx, w, r)
		}

		return response.WriteError(w, http.StatusUnauthorized, errors.New("not admin or owner"))
	}
}

func (a *api) adminOrCompanyOwner(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		claims, err := jwtauth.GetClaims(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		subscriber, err := getSubscriber(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		program, err := a.programCore.QueryByID(ctx, subscriber.ProgramID)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		if claims.Role == usercore.RoleAdmin || claims.ID == program.CompanyID {
			return handler(ctx, w, r)
		}

		return response.WriteError(w, http.StatusUnauthorized, errors.New("not admin or owner"))
	}
}
