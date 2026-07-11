package domain

import "errors"

var (
	ErrInvalidCredentials      = errors.New("invalid email or password")
	ErrInvalidToken            = errors.New("invalid or expired token")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
	ErrMissingUserIDClaim      = errors.New("user_id claim is missing")
	ErrUnexpectedSigningMethod = errors.New("unexpected token signing method")
)
