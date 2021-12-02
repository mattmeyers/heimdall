package client

import (
	"github.com/mattmeyers/heimdall/crypto"
)

const clientIDLength = 32
const clientSecretLength = 64

func generateClientID() (string, error) {
	return crypto.GenerateRandHexString(clientIDLength)
}

func generateClientSecret() (string, error) {
	return crypto.GenerateRandHexString(clientSecretLength)
}
