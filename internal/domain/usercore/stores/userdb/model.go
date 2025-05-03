package userdb

import (
	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/usercore"
)

type user struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
}

func toUserDB(u usercore.User) user {
	return user{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func toUserCore(u user) usercore.User {
	return usercore.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}
