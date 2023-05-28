package manager

import "context"

type StorageManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
