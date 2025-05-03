package usercore

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
}

type CreateUserParams struct {
	Name     string
	Email    string
	Password string
}

type UpdateUserParams struct {
	Name     *string
	Email    *string
	Password *string
}
