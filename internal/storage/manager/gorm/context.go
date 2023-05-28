package gorm

import (
	"context"

	"gorm.io/gorm"
)

type key string

var gormKey = key("gorm_client_key")

func InjectClientToContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, gormKey, db)
}

func ExtractClientFromContext(ctx context.Context) (*gorm.DB, error) {
	gormClient, ok := ctx.Value(gormKey).(*gorm.DB)
	if !ok {
		return nil, ErrInvalidGormClient
	}
	if gormClient == nil {
		return nil, ErrInvalidGormClient
	}

	return gormClient, nil
}
