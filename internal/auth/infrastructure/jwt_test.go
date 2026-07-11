package infrastructure_test

import (
	"testing"
	"time"

	authdomain "github.com/afryn123/withdraw-service-test/internal/auth/domain"
	"github.com/afryn123/withdraw-service-test/internal/auth/infrastructure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestJWTManager(t *testing.T) {
	manager := infrastructure.NewJWTManager("test-secret", time.Hour)
	userID := uuid.New()

	token, err := manager.Generate(userID)
	require.NoError(t, err)
	parsedID, err := manager.Parse(token)
	require.NoError(t, err)
	require.Equal(t, userID, parsedID)

	_, err = manager.Parse("not-a-token")
	require.ErrorIs(t, err, authdomain.ErrInvalidToken)

	otherManager := infrastructure.NewJWTManager("other-secret", time.Hour)
	_, err = otherManager.Parse(token)
	require.ErrorIs(t, err, authdomain.ErrInvalidToken)
}
