package auth

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"html/template"

	"github.com/mattmeyers/heimdall/crypto"
	"github.com/mattmeyers/heimdall/store"
)

//go:embed templates/*
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*"))

type Service struct {
	userStore   store.UserStore
	clientStore store.ClientStore
	jwtSettings JWTSettings
}

func NewService(userStore store.UserStore, clientStore store.ClientStore, jwtSettings JWTSettings) (*Service, error) {
	return &Service{userStore: userStore, clientStore: clientStore, jwtSettings: jwtSettings}, nil
}

func (s *Service) Login(ctx context.Context, email, password, clientID, redirectURL string) (Token, error) {
	u, err := s.userStore.GetByEmail(ctx, email)
	if err != nil {
		return Token{}, err
	}

	valid, err := crypto.ValidatePassword(password, u.Hash)
	if err != nil {
		return Token{}, err
	}

	if !valid {
		return Token{}, errors.New("invalid password")
	}

	err = s.validateRedirectURL(ctx, clientID, redirectURL)
	if err != nil {
		return Token{}, err
	}

	return generateJWT(s.jwtSettings)
}

func (s *Service) validateRedirectURL(ctx context.Context, clientID, redirectURL string) error {
	c, err := s.clientStore.GetByClientID(ctx, clientID)
	if err != nil {
		return err
	}

	for _, u := range c.RedirectURLs {
		if redirectURL == u {
			return nil
		}
	}

	return errors.New("invalid redirect URL")
}

func (s *Service) ImplicitFlow(ctx context.Context, clientID, redirectURL string) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := templates.ExecuteTemplate(
		buf,
		"implicit_flow.html",
		map[string]interface{}{"clientID": clientID, "redirectURL": redirectURL},
	)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
