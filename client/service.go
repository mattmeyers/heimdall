package client

import (
	"context"

	"github.com/mattmeyers/heimdall/store"
)

type Service struct {
	clientStore store.ClientStore
}

func NewService(s store.ClientStore) (*Service, error) {
	return &Service{clientStore: s}, nil
}

func (s *Service) Get(ctx context.Context, clientID string) (store.Client, error) {
	return s.clientStore.GetByClientID(ctx, clientID)
}

func (s *Service) Register(ctx context.Context) (c store.Client, err error) {
	if c.ClientID, err = generateClientID(); err != nil {
		return store.Client{}, err
	}

	if c.ClientSecret, err = generateClientSecret(); err != nil {
		return store.Client{}, err
	}

	if c.ID, err = s.clientStore.Create(ctx, c); err != nil {
		return store.Client{}, err
	}

	return c, nil
}
