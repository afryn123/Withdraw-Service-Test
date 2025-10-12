package utils

import (
	"log"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	txnLogger  *log.Logger
	systemLogger *log.Logger
	once       sync.Once
)

// InitLoggers 
func InitLoggers() {
	once.Do(func() {
		// circular logging for transaction logs
		txnLogger = log.New(&lumberjack.Logger{
			Filename:   "./logs/transaction.log",
			MaxSize:    10, // MB
			MaxBackups: 5, // number of backups
			MaxAge:     30, // days
			Compress:   true,
		}, "TXN: ", log.LstdFlags)

		// circular logging for system logs
		systemLogger = log.New(&lumberjack.Logger{
			Filename:   "./logs/system.log",
			MaxSize:    20, // MB
			MaxBackups: 10, // number of backups
			MaxAge:     30, // days
			Compress:   true,
		}, "SYSTEM: ", log.LstdFlags)
	})
}

// Returns the logger based on category
func GetLogger(category string) *log.Logger {
	switch category {
	case "txn":
		return txnLogger
	case "system":
		return systemLogger
	default:
		return log.Default()
	}
}
