package gorm

import (
	"context"

	"github.com/defryheryanto/mini-wallet/internal/transaction"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (r *TransactionRepository) FindTransactionsByWalletId(ctx context.Context, walletId string) ([]*transaction.Transaction, error) {
	transactions := []*Transaction{}

	err := r.db.Where("wallet_id = ?", walletId).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return SliceToServiceModel(transactions), nil
}
