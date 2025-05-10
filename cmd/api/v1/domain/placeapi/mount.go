package placeapi

import (
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/placecore"
	"github.com/hemozeetah/journi/internal/domain/placecore/stores/placedb"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
	"github.com/jmoiron/sqlx"
)

func Mount(mux *muxer.Mux, log *logger.Logger, db *sqlx.DB, auth *jwtauth.Auth) {
	const version = "v1"

	a := api{
		log:  log,
		core: placecore.New(log, placedb.New(log, db)),
	}

	authen := auth.AuthenticateMW()
	adminOnly := auth.AuthorizeMW(usercore.RoleAdmin)

	mux.HandlerFunc(http.MethodPost, version, "/places", a.create, authen, adminOnly)
	mux.HandlerFunc(http.MethodGet, version, "/places", a.query)
	mux.HandlerFunc(http.MethodGet, version, "/places/{place_id}", a.queryByID, a.parsePlaceMW)
	mux.HandlerFunc(http.MethodPut, version, "/places/{place_id}", a.update, a.parsePlaceMW, authen, adminOnly)
	mux.HandlerFunc(http.MethodDelete, version, "/places/{place_id}", a.delete, a.parsePlaceMW, authen, adminOnly)
}
