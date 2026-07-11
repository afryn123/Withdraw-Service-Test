package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLevel(t *testing.T) {
	require.Equal(t, slog.LevelDebug, parseLevel("DEBUG"))
	require.Equal(t, slog.LevelWarn, parseLevel("warning"))
	require.Equal(t, slog.LevelError, parseLevel("error"))
	require.Equal(t, slog.LevelInfo, parseLevel("unknown"))

	logger := New("error")
	require.False(t, logger.Handler().Enabled(context.Background(), slog.LevelInfo))
	require.True(t, logger.Handler().Enabled(context.Background(), slog.LevelError))
}
