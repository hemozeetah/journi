package userapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/muxer"
)

const userKey = "user"

func setUser(ctx context.Context, user usercore.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func getUser(ctx context.Context) (usercore.User, error) {
	user, ok := ctx.Value(userKey).(usercore.User)
	if !ok {
		return usercore.User{}, errors.New("user not found")
	}

	return user, nil
}

func (a *api) parseUserMW(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		userID, err := uuid.Parse(r.PathValue("user_id"))
		if err != nil {
			return response.WriteError(w, http.StatusBadRequest, err)
		}

		user, err := a.core.QueryByID(ctx, userID)
		if err != nil {
			if errors.Is(err, usercore.ErrNotFound) {
				return response.WriteError(w, http.StatusNotFound, usercore.ErrNotFound)
			}
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		ctx = setUser(ctx, user)
		return handler(ctx, w, r)
	}
}

func (a *api) adminOrOwner(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		claims, err := jwtauth.GetClaims(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		user, err := getUser(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		if claims.Role == usercore.RoleAdmin || claims.ID == user.ID {
			return handler(ctx, w, r)
		}

		return response.WriteError(w, http.StatusUnauthorized, errors.New("not admin or owner"))
	}
}
