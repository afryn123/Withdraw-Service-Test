package infrastructure

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CodeGenerator struct{ now func() time.Time }

func NewCodeGenerator() *CodeGenerator { return &CodeGenerator{now: time.Now} }

func (g *CodeGenerator) TransactionCode() string {
	return "TXN-" + strings.ReplaceAll(g.now().Format("20060102-150405.000"), ".", "")
}

func (g *CodeGenerator) ReferenceNumber(userID uuid.UUID) string {
	now := g.now()
	shortUser := strings.ToUpper(strings.ReplaceAll(userID.String(), "-", "")[:6])
	return fmt.Sprintf("REF-%s-%s-%d", now.Format("20060102"), shortUser, now.UnixMilli())
}
