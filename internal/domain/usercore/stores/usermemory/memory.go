package usermemory

import (
	"context"
	"sync"

	"slices"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/usercore"
)

type Memory struct {
	Users []usercore.User
	mu    sync.RWMutex
}

func New() *Memory {
	return &Memory{}
}

func (mem *Memory) Create(ctx context.Context, user usercore.User) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	for _, u := range mem.Users {
		if u.Email == user.Email {
			return usercore.ErrUniqueEmail
		}
	}
	mem.Users = append(mem.Users, user)
	return nil
}

func (mem *Memory) Update(ctx context.Context, user usercore.User) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	for i, u := range mem.Users {
		if u.ID == user.ID {
			mem.Users[i] = user
			return nil
		}
	}

	return usercore.ErrNotFound
}

func (mem *Memory) Delete(ctx context.Context, user usercore.User) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	for i, u := range mem.Users {
		if u.ID == user.ID {
			mem.Users = slices.Delete(mem.Users, i, i+1)
			return nil
		}
	}

	return usercore.ErrNotFound
}

func (mem *Memory) QueryByID(ctx context.Context, id uuid.UUID) (usercore.User, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	for _, u := range mem.Users {
		if u.ID == id {
			return u, nil
		}
	}

	return usercore.User{}, usercore.ErrNotFound
}

func (mem *Memory) QueryByEmail(ctx context.Context, email string) (usercore.User, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	for _, u := range mem.Users {
		if u.Email == email {
			return u, nil
		}
	}

	return usercore.User{}, usercore.ErrNotFound
}

func (mem *Memory) Query(ctx context.Context) ([]usercore.User, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	users := make([]usercore.User, len(mem.Users))
	copy(users, mem.Users)

	return users, nil
}
