package hash

import "errors"

var (
	ErrFailedGenerateSalt       = errors.New("failed to generate salt")
	ErrInvalidHashFormat        = errors.New("invalid hash format")
	ErrIncompatibleArgonVersion = errors.New("incompatible argon2 version")
	ErrFailedDecryptSalt        = errors.New("failed to decrypt salt")
	ErrFailedDecryptHash        = errors.New("failed to decrypt hash")
)