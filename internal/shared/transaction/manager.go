package transaction

import "context"

// Manager abstracts the transaction boundary so application services do not
// depend on GORM and can be tested without a database.
type Manager interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
}
