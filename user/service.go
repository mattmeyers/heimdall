package user

import (
	"context"
	"errors"

	"github.com/mattmeyers/heimdall/store"
)

type UserStore interface {
	GetByID(ctx context.Context, id int) (store.User, error)
	GetByEmail(ctx context.Context, email string) (store.User, error)
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

func (s *Service) Login(ctx context.Context, email string, password string) error {
	u, err := s.userStore.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	valid, err := validatePassword(password, u.Hash)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid password")
	}

	return nil
}

func (s *Service) Register(ctx context.Context, email, password string) (int, error) {
	hash, err := getPasswordHash(password, DefaultParams)
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
