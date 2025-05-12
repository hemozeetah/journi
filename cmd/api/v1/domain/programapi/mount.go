package programapi

import (
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/programcore"
	"github.com/hemozeetah/journi/internal/domain/programcore/stores/programdb"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
	"github.com/jmoiron/sqlx"
)

func Mount(mux *muxer.Mux, log *logger.Logger, db *sqlx.DB, auth *jwtauth.Auth) {
	const version = "v1"

	a := api{
		log:  log,
		core: programcore.New(log, programdb.New(log, db)),
	}

	authen := auth.AuthenticateMW()

	mux.HandlerFunc(http.MethodPost, version, "/programs", a.create, authen)
	mux.HandlerFunc(http.MethodGet, version, "/programs", a.query)
	mux.HandlerFunc(http.MethodGet, version, "/programs/{program_id}", a.queryByID, a.parsePostMW)
	mux.HandlerFunc(http.MethodPut, version, "/programs/{program_id}", a.update, a.parsePostMW, authen, a.adminOrOwner)
	mux.HandlerFunc(http.MethodDelete, version, "/programs/{program_id}", a.delete, a.parsePostMW, authen, a.adminOrOwner)
}
