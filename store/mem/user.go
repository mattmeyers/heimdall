package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/mattmeyers/heimdall/store"
)

type Store struct {
	db *DB

	lock *sync.Mutex
}

func NewUserStore(db *DB) (*Store, error) {
	return &Store{
		db:   db,
		lock: &sync.Mutex{},
	}, nil
}

func (s *Store) GetByID(ctx context.Context, id int) (store.User, error) {
	u, ok := s.db.user.rows[id]
	if !ok {
		return store.User{}, errors.New("user not found")
	}

	return *u, nil
}

func (s *Store) GetByEmail(ctx context.Context, email string) (store.User, error) {
	for _, u := range s.db.user.rows {
		if u.Email == email {
			return *u, nil
		}
	}

	return store.User{}, errors.New("user not found")
}

func (s *Store) Create(ctx context.Context, u store.User) (int, error) {
	if _, err := s.GetByEmail(ctx, u.Email); err == nil {
		return 0, errors.New("user already exists")
	}

	u.ID = s.db.user.nextID
	s.db.user.nextID++

	s.db.user.rows[u.ID] = &u

	return u.ID, nil
}
