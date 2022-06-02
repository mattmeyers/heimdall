package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ArgonParams holds the configuration used for generating argon2 password hashes.
// Argon2 configurations are dependent upon the host system and must be tweaked to
// maximize the tradeoff between hash speed and resource usage. For additional info,
// see section 4 of the Argon2 RFC (https://datatracker.ietf.org/doc/html/rfc9106).
type ArgonParams struct {
	// Time is the max number of seconds that a hashing can afford to take. This parameter
	// can be used to tune the algorithm independent of memory constraints.
	Time uint32
	// Memory is the max amount of memory (in KiB) that can be used by the hashing algorithm.
	Memory uint32
	// Threads is the number of concurrent (but synchronizing) threads that can be
	// used to compute the hash.
	Threads uint8
	// KeyLen is the length (in bytes) of the final generated hash.
	KeyLen uint32
	// SaltLen is the length (in bytes) of the generated salt.
	SaltLen uint32
}

// DefaultParams is the configuration recommended for all environments . A custom
// configuration should be provided for a production deployment in order to harden
// the service for the hardware it is running on.
var DefaultParams = ArgonParams{
	Time:    1,
	Memory:  32_768, // 32 MiB
	Threads: 4,
	KeyLen:  32,
	SaltLen: 16,
}

// ValidatePassword determines if the provided plain-text password matches the
// encoded hash. Validity is determined by the first return paramter. An error will
// only be returned if the encoded hash is malformed, or the password cannot be hashed.
func ValidatePassword(password, encodedHash string) (bool, error) {
	hash, salt, params, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash, err := hashPassword(password, salt, params)
	if err != nil {
		return false, err
	}

	return hashesAreEqual([]byte(hash), otherHash), nil
}

func hashesAreEqual(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// GetPasswordHash generates an encoded password hash using the argon2id hashing algorith.
// The returned string takes the form
//
//	`$argon2id$v=<argon2 VERISON>$m=<MEMORY>,t=<TIME>,p=<THREADS>$<SALT>$<HASH>`
//
// This encoding provides all of the information required to recompute a hash and validate
// a provided password.
func GetPasswordHash(password string, p ArgonParams) (string, error) {
	salt, err := generateSalt(p.SaltLen)
	if err != nil {
		return "", err
	}

	hash, err := hashPassword(password, salt, p)
	if err != nil {
		return "", err
	}

	return encodeHash(hash, salt, p), nil
}

func hashPassword(password string, salt []byte, p ArgonParams) ([]byte, error) {
	return argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, p.KeyLen), nil
}

func encodeHash(hash, salt []byte, p ArgonParams) string {
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.Memory,
		p.Time,
		p.Threads,
		b64Salt,
		b64Hash,
	)
}

func decodeHash(encodedHash string) ([]byte, []byte, ArgonParams, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[0] != "" {
		return nil, nil, ArgonParams{}, errors.New("malformed hash encoding")
	}

	if parts[1] != "argon2id" {
		return nil, nil, ArgonParams{}, errors.New("unsupported argon2 algorithm")
	}

	if parts[2] != fmt.Sprintf("v=%d", argon2.Version) {
		return nil, nil, ArgonParams{}, errors.New("unsupported argon2 version")
	}

	var p ArgonParams
	n, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Time, &p.Threads)
	if err != nil {
		return nil, nil, ArgonParams{}, err
	} else if n != 3 {
		return nil, nil, ArgonParams{}, errors.New("malformed hash encoding")
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil || len(salt) == 0 {
		return nil, nil, ArgonParams{}, errors.New("malformed salt")
	}

	hash, err := base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil || len(hash) == 0 {
		return nil, nil, ArgonParams{}, errors.New("malformed hash")
	}

	p.KeyLen = uint32(len(hash))
	p.SaltLen = uint32(len(salt))

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
