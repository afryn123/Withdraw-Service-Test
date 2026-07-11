package gormdb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type SlogLogger struct {
	logger        *slog.Logger
	level         gormlogger.LogLevel
	slowThreshold time.Duration
}

func NewSlogLogger(logger *slog.Logger) *SlogLogger {
	return &SlogLogger{logger: logger.With("component", "database"), level: gormlogger.Warn, slowThreshold: time.Second}
}

func (l *SlogLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	copy := *l
	copy.level = level
	return &copy
}

func (l *SlogLogger) Info(ctx context.Context, message string, data ...any) {
	if l.level >= gormlogger.Info {
		l.logger.InfoContext(ctx, fmt.Sprintf(message, data...))
	}
}

func (l *SlogLogger) Warn(ctx context.Context, message string, data ...any) {
	if l.level >= gormlogger.Warn {
		l.logger.WarnContext(ctx, fmt.Sprintf(message, data...))
	}
}

func (l *SlogLogger) Error(ctx context.Context, message string, data ...any) {
	if l.level >= gormlogger.Error {
		l.logger.ErrorContext(ctx, fmt.Sprintf(message, data...))
	}
}

func (l *SlogLogger) Trace(ctx context.Context, startedAt time.Time, query func() (string, int64), err error) {
	if l.level == gormlogger.Silent {
		return
	}
	elapsed := time.Since(startedAt)
	sql, rows := query()
	attributes := []any{"latency_ms", elapsed.Milliseconds(), "rows", rows, "sql", sql}

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && l.level >= gormlogger.Error:
		l.logger.ErrorContext(ctx, "database query failed", append(attributes, "error", err)...)
	case elapsed > l.slowThreshold && l.level >= gormlogger.Warn:
		l.logger.WarnContext(ctx, "slow database query", attributes...)
	case l.level >= gormlogger.Info:
		l.logger.DebugContext(ctx, "database query completed", attributes...)
	}
}
