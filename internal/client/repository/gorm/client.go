package gorm

import (
	"context"

	"github.com/defryheryanto/mini-wallet/internal/client"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{db}
}

func (r *ClientRepository) Insert(ctx context.Context, data *client.Client) error {
	payload := Client{}.FromServiceModel(data)

	err := r.db.Create(&payload).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientRepository) FindByXid(ctx context.Context, xid string) (*client.Client, error) {
	result := &Client{}

	err := r.db.Where("xid = ?", xid).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return result.ToServiceModel(), nil
}

func (r *ClientRepository) FindByToken(ctx context.Context, token string) (*client.Client, error) {
	result := &Client{}

	err := r.db.Where("token = ?", token).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return result.ToServiceModel(), nil
}
