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
	m     sync.RWMutex
}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) Create(ctx context.Context, user usercore.User) error {
	m.m.Lock()
	defer m.m.Unlock()

	for _, u := range m.Users {
		if u.Email == user.Email {
			return usercore.ErrUniqueEmail
		}
	}
	m.Users = append(m.Users, user)
	return nil
}

func (m *Memory) Update(ctx context.Context, user usercore.User) error {
	m.m.Lock()
	defer m.m.Unlock()

	for i, u := range m.Users {
		if u.ID == user.ID {
			m.Users[i] = user
			return nil
		}
	}

	return usercore.ErrNotFound
}

func (m *Memory) Delete(ctx context.Context, user usercore.User) error {
	m.m.Lock()
	defer m.m.Unlock()

	for i, u := range m.Users {
		if u.ID == user.ID {
			m.Users = slices.Delete(m.Users, i, i+1)
			return nil
		}
	}

	return usercore.ErrNotFound
}

func (m *Memory) QueryByID(ctx context.Context, id uuid.UUID) (usercore.User, error) {
	m.m.RLock()
	defer m.m.RUnlock()

	for _, u := range m.Users {
		if u.ID == id {
			return u, nil
		}
	}

	return usercore.User{}, usercore.ErrNotFound
}

func (m *Memory) QueryByEmail(ctx context.Context, email string) (usercore.User, error) {
	m.m.RLock()
	defer m.m.RUnlock()

	for _, u := range m.Users {
		if u.Email == email {
			return u, nil
		}
	}

	return usercore.User{}, usercore.ErrNotFound
}

func (m *Memory) Query(ctx context.Context) ([]usercore.User, error) {
	m.m.RLock()
	defer m.m.RUnlock()

	users := make([]usercore.User, len(m.Users))
	copy(users, m.Users)

	return users, nil
}
