package authapi

import (
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
)

func Mount(mux *muxer.Mux, log *logger.Logger, auth *jwtauth.Auth) {
	const version = "v1"

	a := api{
		log:  log,
		auth: auth,
	}

	basic := a.auth.BasicMW()
	authen := a.auth.AuthenticateMW()
	adminOnly := a.auth.AuthorizeMW(usercore.RoleAdmin)

	mux.HandlerFunc(http.MethodGet, version, "/auth/token", a.token, basic)
	mux.HandlerFunc(http.MethodGet, version, "/auth/claims", a.claims, authen)
	mux.HandlerFunc(http.MethodGet, version, "/auth/admin", a.authorizeAdmin, authen, adminOnly)
}
