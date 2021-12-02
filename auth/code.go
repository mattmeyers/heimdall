package auth

import "github.com/mattmeyers/heimdall/crypto"

func generateAuthCode() (string, error) {
	return crypto.GenerateRandHexString(32)
}
