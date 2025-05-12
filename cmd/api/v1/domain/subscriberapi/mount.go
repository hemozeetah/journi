package subscriberapi

import (
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/programcore"
	"github.com/hemozeetah/journi/internal/domain/programcore/stores/programdb"
	"github.com/hemozeetah/journi/internal/domain/subscribercore"
	"github.com/hemozeetah/journi/internal/domain/subscribercore/stores/subscriberdb"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
	"github.com/jmoiron/sqlx"
)

func Mount(mux *muxer.Mux, log *logger.Logger, db *sqlx.DB, auth *jwtauth.Auth) {
	const version = "v1"

	a := api{
		log:         log,
		core:        subscribercore.New(log, subscriberdb.New(log, db)),
		programCore: programcore.New(log, programdb.New(log, db)),
	}

	authen := auth.AuthenticateMW()

	mux.HandlerFunc(http.MethodPost, version, "/subscribers", a.create, authen)
	mux.HandlerFunc(http.MethodGet, version, "/subscribers", a.query)
	mux.HandlerFunc(http.MethodGet, version, "/subscribers/{reference_id}", a.queryByID, a.parseSubscriberMW)
	mux.HandlerFunc(http.MethodPut, version, "/subscribers/{reference_id}", a.update, a.parseSubscriberMW, authen, a.adminOrCompanyOwner)
	mux.HandlerFunc(http.MethodDelete, version, "/subscribers/{reference_id}", a.delete, a.parseSubscriberMW, authen, a.adminOrOwnerOrCompanyOwner)
}
