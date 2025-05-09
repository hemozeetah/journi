package usercore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is already exists")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

type Storer interface {
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, user User) error
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email string) (User, error)
	Query(ctx context.Context) ([]User, error)
}

type Core struct {
	log   *logger.Logger
	store Storer
}

func New(log *logger.Logger, store Storer) *Core {
	return &Core{
		log:   log,
		store: store,
	}
}

func (c *Core) Create(ctx context.Context, p CreateUserParams) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	user := User{
		ID:        uuid.New(),
		Name:      p.Name,
		Email:     p.Email,
		Password:  string(hash),
		Role:      "user",
		Profile:   p.Profile,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.store.Create(ctx, user); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return user, nil
}

func (c *Core) Update(ctx context.Context, user User, p UpdateUserParams) (User, error) {
	if p.Name != nil {
		user.Name = *p.Name
	}
	if p.Email != nil {
		user.Email = *p.Email
	}
	if p.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*p.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generatefrompassword: %w", err)
		}
		user.Password = string(hash)
	}
	if p.Role != nil {
		user.Role = *p.Role
	}
	if p.Profile != nil {
		user.Profile = *p.Profile
	}
	user.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, user); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}

	return user, nil
}

func (c *Core) Delete(ctx context.Context, user User) error {
	if err := c.store.Delete(ctx, user); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := c.store.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("querybyid[%s]: %w", userID, err)
	}

	return user, nil
}

func (c *Core) QueryByEmail(ctx context.Context, email string) (User, error) {
	user, err := c.store.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("querybyemail[%s]: %w", email, err)
	}

	return user, nil
}

func (c *Core) Query(ctx context.Context) ([]User, error) {
	users, err := c.store.Query(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return users, nil
}

func (c *Core) Authenticate(ctx context.Context, email, password string) (User, error) {
	usr, err := c.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("querybyemail[%s]: %w", email, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		return User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return usr, nil
}
