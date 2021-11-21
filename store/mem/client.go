package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/mattmeyers/heimdall/store"
)

type ClientStore struct {
	db *DB

	lock *sync.RWMutex
}

func NewClientStore(db *DB) (*ClientStore, error) {
	return &ClientStore{
		db:   db,
		lock: &sync.RWMutex{},
	}, nil
}

func (s *ClientStore) GetByClientID(ctx context.Context, clientID string) (store.Client, error) {
	for _, c := range s.db.client.rows {
		if c.ClientID == clientID {
			return *c, nil
		}
	}

	return store.Client{}, errors.New("client not found")
}

func (s *ClientStore) Create(ctx context.Context, c store.Client) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, err := s.GetByClientID(ctx, c.ClientID); err == nil {
		return 0, errors.New("client already exists")
	}

	c.ID = s.db.client.nextID
	s.db.client.nextID++

	s.db.client.rows[c.ID] = &c

	return c.ID, nil
}
