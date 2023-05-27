package transaction

import (
	"context"
	"time"

	"github.com/defryheryanto/mini-wallet/internal/wallet"
)

type Transaction struct {
	Id           string    `json:"id"`
	Status       string    `json:"status"`
	TransactedAt time.Time `json:"transacted_at"`
	Type         string    `json:"type"`
	Amount       float64   `json:"amount"`
	ReferenceId  string    `json:"reference_id"`
	WalletId     string    `json:"wallet_id"`
}

type TransactionRepository interface {
	FindTransactionsByWalletId(ctx context.Context, walletId string) ([]*Transaction, error)
}

type TransactionIService interface {
	GetTransactionsByCustomerXid(ctx context.Context, xid string) ([]*Transaction, error)
}

type TransactionService struct {
	repository    TransactionRepository
	walletService wallet.WalletIService
}

func NewTransactionService(repository TransactionRepository, walletService wallet.WalletIService) *TransactionService {
	return &TransactionService{repository, walletService}
}

func (s *TransactionService) GetTransactionsByCustomerXid(ctx context.Context, xid string) ([]*Transaction, error) {
	targetWallet, err := s.walletService.GetWalletByXid(ctx, xid)
	if err != nil {
		return nil, err
	}
	if targetWallet == nil {
		return nil, wallet.ErrWalletNotFound
	}
	if targetWallet.Status == wallet.STATUS_DISABLED {
		return nil, wallet.ErrWalletDisabled
	}

	transactions, err := s.repository.FindTransactionsByWalletId(ctx, targetWallet.Id)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
