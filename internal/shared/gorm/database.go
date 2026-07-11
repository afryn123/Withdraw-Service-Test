package gormdb

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Open(dsn string, loggers ...*slog.Logger) (*gorm.DB, error) {
	var databaseLogger gormlogger.Interface = gormlogger.Default.LogMode(gormlogger.Silent)
	if len(loggers) > 0 && loggers[0] != nil {
		databaseLogger = NewSlogLogger(loggers[0])
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: databaseLogger})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	return db, sqlDB.Ping()
}

type contextKey struct{}

func DB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(contextKey{}).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}
	return fallback.WithContext(ctx)
}

type TransactionManager struct{ db *gorm.DB }

func NewTransactionManager(db *gorm.DB) *TransactionManager { return &TransactionManager{db: db} }

func (m *TransactionManager) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, contextKey{}, tx))
	})
}
