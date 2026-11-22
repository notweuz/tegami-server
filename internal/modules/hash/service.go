package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"tegami-server/internal/config"

	"golang.org/x/crypto/argon2"
)

const semaphoreSize = 8

type Argon2Params struct {
	Memory, Time, SaltLen, KeyLen uint32
	Threads                       uint8
}

type service struct {
	cfg       *config.Config
	semaphore chan struct{}
}

func NewService(cfg *config.Config) *service {
	return &service{
		cfg:       cfg,
		semaphore: make(chan struct{}, semaphoreSize),
	}
}

func (s *service) HashPassword(password string) (string, error) {
	addToSemaphore(s.semaphore)
	defer releaseFromSemaphore(s.semaphore)

	salt, err := generateSalt(s.cfg.Argon2SaltLen)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	aTime := s.cfg.Argon2Time
	aMemory := s.cfg.Argon2Memory
	aThreads := s.cfg.Argon2Threads

	hash := argon2.IDKey([]byte(password), salt, aTime, aMemory, aThreads, s.cfg.Argon2KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, aMemory, aTime, aThreads, b64Salt, b64Hash)

	return encodedHash, nil
}

func (s *service) VerifyPassword(password, encodedHash string) (bool, error) {
	addToSemaphore(s.semaphore)
	defer releaseFromSemaphore(s.semaphore)

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

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, ErrFailedGenerateSalt
	}
	return salt, nil
}

func addToSemaphore(s chan<- struct{}) {
	s <- struct{}{}
}

func releaseFromSemaphore(s <-chan struct{}) {
	<-s
}
