package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/internal/domain/usercore/stores/userdb"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Log    *logger.Logger
	DB     *sqlx.DB
	JwtKey string
	Issuer string
}

type Auth struct {
	key    []byte
	issuer string
	method jwt.SigningMethod
	core   *usercore.Core
}

func New(cfg Config) *Auth {
	return &Auth{
		key:    []byte(cfg.JwtKey),
		issuer: cfg.Issuer,
		method: jwt.GetSigningMethod(jwt.SigningMethodHS256.Name),
		core:   usercore.New(cfg.Log, userdb.New(cfg.Log, cfg.DB)),
	}
}
