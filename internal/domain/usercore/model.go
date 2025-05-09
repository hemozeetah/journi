package usercore

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      string
	Profile   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserParams struct {
	Name     string
	Email    string
	Password string
	Profile  string
}

type UpdateUserParams struct {
	Name     *string
	Email    *string
	Password *string
	Role     *string
	Profile  *string
}
