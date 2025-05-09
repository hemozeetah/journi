package jwtauth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/usercore"
)

type Claims struct {
	jwt.RegisteredClaims
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	UpdatedAt time.Time `json:"updated_date"`
}

func (a *Auth) newClaims(user usercore.User) Claims {
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			Issuer:    a.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		UpdatedAt: user.UpdatedAt,
	}
}

func (a *Auth) GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(a.method, claims)

	signed, err := token.SignedString(a.key)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (a *Auth) parseToken(signed string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(signed, &claims, func(t *jwt.Token) (any, error) {
		return a.key, nil
	})
	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func (a *Auth) verifyClaims(ctx context.Context, claims Claims) error {
	user, err := a.core.QueryByID(ctx, claims.ID)
	if err != nil {
		return jwt.ErrTokenInvalidId
	}

	if !user.UpdatedAt.Equal(claims.UpdatedAt) {
		return jwt.ErrTokenInvalidClaims
	}

	return nil
}
