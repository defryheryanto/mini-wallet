package gorm

import (
	"context"

	gorm_manager "github.com/defryheryanto/mini-wallet/internal/storage/manager/gorm"
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

	db := r.getGormClient(ctx)
	err := db.Create(&payload).Error
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

func (r *WalletRepository) FindByCustomerXid(ctx context.Context, xid string) (*wallet.Wallet, error) {
	result := &Wallet{}

	err := r.db.Where("owned_by = ?", xid).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return result.ToServiceModel(), nil
}

func (r *WalletRepository) Update(ctx context.Context, data *wallet.Wallet) error {
	result := Wallet{}.FromServiceModel(data)

	db := r.getGormClient(ctx)
	err := db.Where("id = ?", data.Id).Select("*").Updates(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	return nil
}

func (r *WalletRepository) getGormClient(ctx context.Context) *gorm.DB {
	db, err := gorm_manager.ExtractClientFromContext(ctx)
	if err != nil {
		return r.db
	}

	return db
}
