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

func (r *TransactionRepository) FindByReferenceId(ctx context.Context, referenceId string) (*transaction.Transaction, error) {
	transaction := &Transaction{}

	err := r.db.Where("reference_id = ?", referenceId).First(&transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return transaction.ToServiceModel(), nil
}

func (r *TransactionRepository) FindById(ctx context.Context, id string) (*transaction.Transaction, error) {
	transaction := &Transaction{}

	err := r.db.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return transaction.ToServiceModel(), nil
}

func (r *TransactionRepository) Insert(ctx context.Context, data *transaction.Transaction) error {
	trx := Transaction{}.FromServiceModel(data)

	err := r.db.Create(&trx).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) Update(ctx context.Context, data *transaction.Transaction) error {
	trx := Transaction{}.FromServiceModel(data)

	err := r.db.Where("id = ?", trx.Id).Select("*").Updates(&trx).Error
	if err != nil {
		return err
	}

	return nil
}
