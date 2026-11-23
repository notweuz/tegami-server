package jwt

import "errors"

var (
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrJWTInvalid         = errors.New("jwt token invalid")
	ErrTokenExpired       = errors.New("jwt expired")
)
