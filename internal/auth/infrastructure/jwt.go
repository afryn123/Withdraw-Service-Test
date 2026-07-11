package infrastructure

import (
	"fmt"
	"time"

	authdomain "github.com/afryn123/withdraw-service-test/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: ttl}
}

func (m *JWTManager) Generate(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{"user_id": userID.String(), "exp": time.Now().Add(m.ttl).Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

func (m *JWTManager) Parse(raw string) (uuid.UUID, error) {
	token, err := jwt.Parse(raw, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("%w: %s", authdomain.ErrUnexpectedSigningMethod, token.Method.Alg())
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, authdomain.ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, authdomain.ErrInvalidTokenClaims
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, authdomain.ErrMissingUserIDClaim
	}
	return uuid.Parse(userID)
}
