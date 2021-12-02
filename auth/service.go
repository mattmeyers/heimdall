package auth

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"html/template"
	"time"

	"github.com/mattmeyers/heimdall/crypto"
	"github.com/mattmeyers/heimdall/store"
)

//go:embed templates/*
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*"))

type Service struct {
	userStore     store.UserStore
	clientStore   store.ClientStore
	authCodeStore store.AuthCodeStore
	jwtSettings   JWTSettings
}

func NewService(userStore store.UserStore,
	clientStore store.ClientStore,
	authCodeStore store.AuthCodeStore,
	jwtSettings JWTSettings) (*Service, error) {
	return &Service{
		userStore:     userStore,
		clientStore:   clientStore,
		authCodeStore: authCodeStore,
		jwtSettings:   jwtSettings}, nil
}

func (s *Service) Login(ctx context.Context, email, password, clientID, redirectURL string) (string, error) {
	u, err := s.userStore.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	valid, err := crypto.ValidatePassword(password, u.Hash)
	if err != nil {
		return "", err
	}

	if !valid {
		return "", errors.New("invalid password")
	}

	err = s.validateRedirectURL(ctx, clientID, redirectURL)
	if err != nil {
		return "", err
	}

	code, err := generateAuthCode()
	if err != nil {
		return "", err
	}

	_, err = s.authCodeStore.Insert(ctx, store.AuthCode{Code: code, UserID: u.ID})
	if err != nil {
		return "", err
	}

	return code, nil
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

func (s *Service) AuthCodeFlow(ctx context.Context, clientID, redirectURL string) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := templates.ExecuteTemplate(
		buf,
		"auth_code_flow.html",
		map[string]interface{}{"clientID": clientID, "redirectURL": redirectURL},
	)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Service) ConvertCodeToToken(ctx context.Context, code, clientID, clientSecret, redirectURL string) (Token, error) {
	client, err := s.clientStore.GetByClientID(ctx, clientID)
	if err != nil {
		return Token{}, err
	}

	if client.ClientSecret != clientSecret {
		return Token{}, errors.New("incorrect secret")
	}

	codeObj, err := s.authCodeStore.GetByCode(ctx, code)
	if err != nil {
		return Token{}, err
	} else if time.Now().After(codeObj.CreatedAt.Add(3600 * time.Second)) {
		return Token{}, errors.New("access code has expired")
	}

	return generateJWT(s.jwtSettings)
}
