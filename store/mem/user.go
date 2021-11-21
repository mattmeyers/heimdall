package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/mattmeyers/heimdall/store"
)

type Store struct {
	users  map[int]*store.User
	nextID int

	lock *sync.Mutex
}

func NewUserStore() (*Store, error) {
	return &Store{
		users:  make(map[int]*store.User),
		nextID: 1,
		lock:   &sync.Mutex{},
	}, nil
}

func (s *Store) GetByID(ctx context.Context, id int) (store.User, error) {
	u, ok := s.users[id]
	if !ok {
		return store.User{}, errors.New("user not found")
	}

	return *u, nil
}

func (s *Store) GetByEmail(ctx context.Context, email string) (store.User, error) {
	for _, u := range s.users {
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

	u.ID = s.nextID
	s.nextID++

	s.users[u.ID] = &u

	return u.ID, nil
}
