package manager

import (
	"context"
)

type MockStorageManager struct {
}

func (m *MockStorageManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
