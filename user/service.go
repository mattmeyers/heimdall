package user

import (
	"context"

	"github.com/mattmeyers/heimdall/crypto"
	"github.com/mattmeyers/heimdall/store"
)

type UserStore interface {
	GetByID(ctx context.Context, id int) (store.User, error)
	Create(ctx context.Context, u store.User) (int, error)
}

type Service struct {
	userStore UserStore
}

func NewService(s UserStore) (*Service, error) {
	return &Service{userStore: s}, nil
}

func (s *Service) Get(ctx context.Context, id int) (store.User, error) {
	return s.userStore.GetByID(ctx, id)
}

func (s *Service) Register(ctx context.Context, email, password string) (int, error) {
	hash, err := crypto.GetPasswordHash(password, crypto.DefaultParams)
	if err != nil {
		return 0, err
	}

	u := store.User{Email: email, Hash: hash}
	id, err := s.userStore.Create(ctx, u)
	if err != nil {
		return 0, err
	}

	return id, nil
}
