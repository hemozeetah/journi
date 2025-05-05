package userapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/internal/domain/usercore/stores/userdb"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
	"github.com/jmoiron/sqlx"
)

func Mount(mux *muxer.Mux, log *logger.Logger, db *sqlx.DB) {
	a := api{
		log:  log,
		core: usercore.New(log, userdb.New(log, db)),
	}

	const version = "v1"

	mux.HandlerFunc(http.MethodPost, version, "/users", a.create)
	mux.HandlerFunc(http.MethodGet, version, "/users", a.query)
	mux.HandlerFunc(http.MethodGet, version, "/users/{user_id}", a.queryByID, a.parseUserMW)
}

func setUser(ctx context.Context, user usercore.User) context.Context {
	return context.WithValue(ctx, "user", user)
}

func getUser(ctx context.Context) (usercore.User, error) {
	user, ok := ctx.Value("user").(usercore.User)
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
