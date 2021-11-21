package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argonParams struct {
	iterations uint32
	memory     uint32
	threads    uint8
	keyLen     uint32
	saltLen    uint32
}

var DefaultParams = argonParams{
	iterations: 3,
	memory:     64 * 1024,
	threads:    4,
	keyLen:     32,
	saltLen:    16,
}

func validatePassword(password, encodedHash string) (bool, error) {
	hash, salt, params, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash, err := hashPassword(password, salt, params)
	if err != nil {
		return false, err
	}

	if !hashesAreEqual([]byte(hash), otherHash) {
		return false, nil
	}

	return true, nil
}

func hashesAreEqual(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

func getPasswordHash(password string, p argonParams) (string, error) {
	salt, err := generateSalt(p.saltLen)
	if err != nil {
		return "", err
	}

	hash, err := hashPassword(password, salt, p)
	if err != nil {
		return "", err
	}

	return encodeHash(hash, salt, p), nil
}

func hashPassword(password string, salt []byte, p argonParams) ([]byte, error) {
	return argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.threads, p.keyLen), nil
}

func encodeHash(hash, salt []byte, p argonParams) string {
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.memory,
		p.iterations,
		p.threads,
		b64Salt,
		b64Hash,
	)
}

func decodeHash(encodedHash string) ([]byte, []byte, argonParams, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, argonParams{}, errors.New("malformed hash encoding")
	}

	if parts[1] != "argon2id" {
		return nil, nil, argonParams{}, errors.New("unsupported argon2 algorithm")
	}

	if parts[2] != fmt.Sprintf("v=%d", argon2.Version) {
		return nil, nil, argonParams{}, errors.New("unsupported argon2 version")
	}

	var p argonParams
	n, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.threads)
	if err != nil {
		return nil, nil, argonParams{}, err
	} else if n != 3 {
		return nil, nil, argonParams{}, errors.New("malformed hash encoding")
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return nil, nil, argonParams{}, errors.New("malformed salt")
	}

	hash, err := base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil {
		return nil, nil, argonParams{}, errors.New("malformed hash")
	}

	p.keyLen = uint32(len(hash))

	return hash, salt, p, nil
}

func generateSalt(length uint32) ([]byte, error) {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
