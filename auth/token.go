package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const JWTLifespan = 3600
const JWTIssuer = "Heimdall"

var key = []byte("so secret")

func generateJWT() (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))

	t.Claims = &jwt.RegisteredClaims{
		Issuer:    JWTIssuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * JWTLifespan)),
	}

	return t.SignedString(key)
}
