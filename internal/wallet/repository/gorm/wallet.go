package gorm

import (
	"context"

	"github.com/defryheryanto/mini-wallet/internal/wallet"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db}
}

func (r *WalletRepository) Insert(ctx context.Context, data *wallet.Wallet) error {
	payload := Wallet{}.FromServiceModel(data)

	err := r.db.Create(&payload).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *WalletRepository) FindById(ctx context.Context, id string) (*wallet.Wallet, error) {
	result := &Wallet{}

	err := r.db.Where("id = ?", id).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return result.ToServiceModel(), nil
}
