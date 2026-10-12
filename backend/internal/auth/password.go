package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    uint32 = 3
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 2
	argonKeyLen  uint32 = 32
	saltLength          = 16
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory,
		argonTime,
		argonThreads,
		encodedSalt,
		encodedHash,
	), nil
}

func VerifyPassword(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false, fmt.Errorf("invalid password hash format")
	}

	var memory uint32
	var iterations uint32
	var threads uint64

	for _, param := range strings.Split(parts[3], ",") {
		keyValue := strings.SplitN(param, "=", 2)
		if len(keyValue) != 2 {
			return false, fmt.Errorf("invalid Argon2 parameters")
		}

		value, err := strconv.ParseUint(keyValue[1], 10, 32)
		if err != nil {
			return false, fmt.Errorf("invalid Argon2 parameter: %w", err)
		}

		switch keyValue[0] {
		case "m":
			memory = uint32(value)
		case "t":
			iterations = uint32(value)
		case "p":
			threads = value
		default:
			return false, fmt.Errorf("unsupported Argon2 parameter")
		}
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("decode password salt: %w", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("decode password hash: %w", err)
	}

	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		uint8(threads),
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(expectedHash, computedHash) != 1 {
		return false, nil
	}

	return true, nil
}
