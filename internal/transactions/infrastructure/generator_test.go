package infrastructure

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCodeGenerator(t *testing.T) {
	fixedTime := time.Date(2026, time.July, 11, 12, 30, 45, 123_000_000, time.UTC)
	generator := &CodeGenerator{now: func() time.Time { return fixedTime }}
	userID := uuid.MustParse("abcdef12-3456-7890-abcd-ef1234567890")

	require.Equal(t, "TXN-20260711-123045123", generator.TransactionCode())
	require.Equal(t, "REF-20260711-ABCDEF-1783773045123", generator.ReferenceNumber(userID))
}
