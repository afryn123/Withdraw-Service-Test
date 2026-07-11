package domain

import (
	"github.com/google/uuid"
)

type PasswordVerifier interface {
	Verify(password, hash string) bool
}

type TokenManager interface {
	Generate(userID uuid.UUID) (string, error)
	Parse(token string) (uuid.UUID, error)
}
