package client

import (
	"crypto/rand"
	"encoding/hex"
)

const clientIDLength = 32
const clientSecretLength = 64

func generateClientID() (string, error) {
	return generateRandHexString(clientIDLength)
}

func generateClientSecret() (string, error) {
	return generateRandHexString(clientSecretLength)
}

func generateRandHexString(l int) (string, error) {
	buf := make([]byte, l)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
