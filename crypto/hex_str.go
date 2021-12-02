package crypto

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandHexString(l int) (string, error) {
	buf := make([]byte, l)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
