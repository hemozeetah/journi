package userapi

import (
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/internal/domain/usercore/stores/userdb"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
	"github.com/jmoiron/sqlx"
)

func Mount(mux *muxer.Mux, log *logger.Logger, db *sqlx.DB, auth *jwtauth.Auth) {
	const version = "v1"

	a := api{
		log:  log,
		core: usercore.New(log, userdb.New(log, db)),
	}

	authen := auth.AuthenticateMW()

	mux.HandlerFunc(http.MethodPost, version, "/users", a.create)
	mux.HandlerFunc(http.MethodGet, version, "/users", a.query)
	mux.HandlerFunc(http.MethodGet, version, "/users/{user_id}", a.queryByID, a.parseUserMW)
	mux.HandlerFunc(http.MethodPut, version, "/users/{user_id}", a.update, a.parseUserMW, authen, a.adminOrOwner)
	mux.HandlerFunc(http.MethodDelete, version, "/users/{user_id}", a.delete, a.parseUserMW, authen, a.adminOrOwner)
}
