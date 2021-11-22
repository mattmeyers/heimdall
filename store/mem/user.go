package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/mattmeyers/heimdall/store"
)

type UserStore struct {
	db *DB

	lock *sync.RWMutex
}

func NewUserStore(db *DB) (*UserStore, error) {
	return &UserStore{
		db:   db,
		lock: &sync.RWMutex{},
	}, nil
}

func (s *UserStore) GetByID(ctx context.Context, id int) (store.User, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	u, ok := s.db.user.rows[id]
	if !ok {
		return store.User{}, errors.New("user not found")
	}

	return *u, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (store.User, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, u := range s.db.user.rows {
		if u.Email == email {
			return *u, nil
		}
	}

	return store.User{}, errors.New("user not found")
}

func (s *UserStore) Create(ctx context.Context, u store.User) (int, error) {
	if _, err := s.GetByEmail(ctx, u.Email); err == nil {
		return 0, errors.New("user already exists")
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	u.ID = s.db.user.nextID
	s.db.user.nextID++

	s.db.user.rows[u.ID] = &u

	return u.ID, nil
}
