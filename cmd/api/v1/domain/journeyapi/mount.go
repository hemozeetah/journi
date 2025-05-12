package journeyapi

import (
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/journeycore"
	"github.com/hemozeetah/journi/internal/domain/journeycore/stores/journeydb"
	"github.com/hemozeetah/journi/internal/domain/programcore"
	"github.com/hemozeetah/journi/internal/domain/programcore/stores/programdb"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
	"github.com/jmoiron/sqlx"
)

func Mount(mux *muxer.Mux, log *logger.Logger, db *sqlx.DB, auth *jwtauth.Auth) {
	const version = "v1"

	a := api{
		log:         log,
		core:        journeycore.New(log, journeydb.New(log, db)),
		programCore: programcore.New(log, programdb.New(log, db)),
	}

	authen := auth.AuthenticateMW()
	adminOrCompany := auth.AuthorizeMW(usercore.RoleAdmin, usercore.RoleCompany)

	mux.HandlerFunc(http.MethodPost, version, "/journeys", a.create, authen, adminOrCompany)
	mux.HandlerFunc(http.MethodGet, version, "/journeys", a.query)
	mux.HandlerFunc(http.MethodGet, version, "/journeys/{journey_id}", a.queryByID, a.parseJourneyMW)
	mux.HandlerFunc(http.MethodPut, version, "/journeys/{journey_id}", a.update, a.parseJourneyMW, authen, a.adminOrOwner)
	mux.HandlerFunc(http.MethodDelete, version, "/journeys/{journey_id}", a.delete, a.parseJourneyMW, authen, a.adminOrOwner)
}
