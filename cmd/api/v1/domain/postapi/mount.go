package postapi

import (
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/postcore"
	"github.com/hemozeetah/journi/internal/domain/postcore/stores/postdb"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
	"github.com/jmoiron/sqlx"
)

func Mount(mux *muxer.Mux, log *logger.Logger, db *sqlx.DB, auth *jwtauth.Auth) {
	const version = "v1"

	a := api{
		log:  log,
		core: postcore.New(log, postdb.New(log, db)),
	}

	authen := auth.AuthenticateMW()

	mux.HandlerFunc(http.MethodPost, version, "/posts", a.create, authen)
	mux.HandlerFunc(http.MethodGet, version, "/posts", a.query)
	mux.HandlerFunc(http.MethodGet, version, "/posts/{post_id}", a.queryByID, a.parsePostMW)
	mux.HandlerFunc(http.MethodPut, version, "/posts/{post_id}", a.update, a.parsePostMW, authen, a.adminOrOwner)
	mux.HandlerFunc(http.MethodDelete, version, "/posts/{post_id}", a.delete, a.parsePostMW, authen, a.adminOrOwner)
}
