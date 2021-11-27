package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Token struct {
	SignedString string
	Lifespan     int
}

type SigningAlgorithm string

const (
	HMAC256Algorithm SigningAlgorithm = "HS256"
)

func (a SigningAlgorithm) isValid() bool {
	return a == HMAC256Algorithm
}

type JWTSettings struct {
	Issuer     string
	Lifespan   int
	SigningKey string
	Algorithm  SigningAlgorithm
}

func (s JWTSettings) validate() error {
	if strings.TrimSpace(s.Issuer) == "" {
		return errors.New("JWT issuer cannot be empty")
	}

	if s.Lifespan <= 0 {
		return errors.New("JWT lifetime must be a positive integer")
	}

	if strings.TrimSpace(s.SigningKey) == "" {
		return errors.New("JWT signing key required")
	}

	if !s.Algorithm.isValid() {
		return errors.New("unknown signing algorithm")
	}

	return nil
}

func generateJWT(settings JWTSettings) (Token, error) {
	if err := settings.validate(); err != nil {
		return Token{}, err
	}

	t := jwt.New(jwt.GetSigningMethod(string(settings.Algorithm)))

	now := time.Now()
	t.Claims = &jwt.RegisteredClaims{
		Issuer:    settings.Issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(settings.Lifespan))),
	}

	signed, err := t.SignedString(settings.SigningKey)
	if err != nil {
		return Token{}, err
	}

	return Token{
		SignedString: signed,
		Lifespan:     settings.Lifespan,
	}, nil
}
