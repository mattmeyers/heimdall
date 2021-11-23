package client

import (
	"context"
	"errors"
	"net/url"
	"strings"

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

func (s *Service) Register(ctx context.Context, redirectURLs []string) (store.Client, error) {
	err := validateRedirectURLs(redirectURLs)
	if err != nil {
		return store.Client{}, err
	}

	c := store.Client{RedirectURLs: redirectURLs}

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

func validateRedirectURLs(urls []string) error {
	for _, u := range urls {
		parsedU, err := url.Parse(u)
		if err != nil {
			return err
		}

		if parsedU.Fragment != "" || strings.HasSuffix(u, "#") {
			return errors.New("redirect url must not contain a fragment")
		}
	}

	return nil
}
