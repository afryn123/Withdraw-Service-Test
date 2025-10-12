package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Format: TXN-20251012-172856123456 (date + miliseconds)
func GenerateTransactionCode() string {
	timestamp := time.Now().Format("20060102-150405.000") // 20251012-172856.123
	cleanTime := strings.ReplaceAll(timestamp, ".", "")   // hapus titik
	return fmt.Sprintf("TXN-%s", cleanTime)
}

// Format: REF-20251012-<userShort>-<epoch>
func GenerateReferenceNumber(userID string) string {
	now := time.Now()
	shortUser := ""
	if len(userID) >= 6 {
		shortUser = userID[:6]
	} else {
		shortUser = fmt.Sprintf("%-6s", userID)
	}

	epochMilli := now.UnixMilli() // epoch dalam milidetik
	return fmt.Sprintf("REF-%s-%s-%d", now.Format("20060102"), strings.ToUpper(shortUser), epochMilli)
}

// Contoh hasil: "9F7A12C3"
func GenerateCompactUUID() string {
	id := uuid.New()
	return strings.ToUpper(id.String()[:8])
}
