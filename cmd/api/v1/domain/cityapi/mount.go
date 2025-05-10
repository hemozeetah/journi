package cityapi

import (
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/citycore"
	"github.com/hemozeetah/journi/internal/domain/citycore/stores/citydb"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
	"github.com/jmoiron/sqlx"
)

func Mount(mux *muxer.Mux, log *logger.Logger, db *sqlx.DB, auth *jwtauth.Auth) {
	const version = "v1"

	a := api{
		log:  log,
		core: citycore.New(log, citydb.New(log, db)),
	}

	authen := auth.AuthenticateMW()
	adminOnly := auth.AuthorizeMW(usercore.RoleAdmin)

	mux.HandlerFunc(http.MethodPost, version, "/cities", a.create, authen, adminOnly)
	mux.HandlerFunc(http.MethodGet, version, "/cities", a.query)
	mux.HandlerFunc(http.MethodGet, version, "/cities/{city_id}", a.queryByID, a.parseCityMW)
	mux.HandlerFunc(http.MethodPut, version, "/cities/{city_id}", a.update, a.parseCityMW, authen, adminOnly)
	mux.HandlerFunc(http.MethodDelete, version, "/cities/{city_id}", a.delete, a.parseCityMW, authen, adminOnly)
}
