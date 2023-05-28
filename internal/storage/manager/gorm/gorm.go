package gorm

import (
	"context"

	"gorm.io/gorm"
)

type GormStorageManager struct {
	db *gorm.DB
}

func NewGormStorageManager(db *gorm.DB) *GormStorageManager {
	return &GormStorageManager{db}
}

// Wrap the given function with database transaction.
// Any SQL queries that want to use this database transaction should use the gorm client inside the context
// use ExtractClientFromContext(context.Context) to get the gorm client inside the context
//
// Transaction will be rollback if received error from the given function. And will be commited if received no error
func (m *GormStorageManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	db := m.db.Begin()
	ctx = InjectClientToContext(ctx, db)

	err := fn(ctx)
	if err != nil {
		db.Rollback()
		return err
	}

	db.Commit()
	return nil
}
