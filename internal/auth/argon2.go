package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	Memory, Time, SaltLen, KeyLen uint32
	Threads                       uint8
}

var params = Argon2Params{
	Memory:  16 * 1024,
	Time:    3,
	SaltLen: 16,
	KeyLen:  32,
	Threads: 4,
}

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, ErrFailedGenerateSalt
	}
	return salt, nil
}

func HashPassword(password string) (string, error) {
	salt, err := generateSalt(params.SaltLen)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, params.Memory, params.Time, params.Threads, b64Salt, b64Hash)

	return encodedHash, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, fmt.Errorf("verify password: %w", ErrInvalidHashFormat)
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, fmt.Errorf("verify password: %w", ErrIncompatibleArgonVersion)
	}
	if version != argon2.Version {
		return false, fmt.Errorf("verify password: %w", ErrIncompatibleArgonVersion)
	}

	p := Argon2Params{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Time, &p.Threads); err != nil {
		return false, fmt.Errorf("verify password: %w", err)
	}

	if p.Time == 0 || p.Threads == 0 || p.Memory == 0 {
		return false, fmt.Errorf("verify password: %w", ErrInvalidHashFormat)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("verify password: %w", ErrFailedDecryptSalt)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("verify password: %w", ErrFailedDecryptHash)
	}
	keyLength := uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, keyLength)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
