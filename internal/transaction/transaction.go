package transaction

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/defryheryanto/mini-wallet/internal/wallet"
	"github.com/google/uuid"
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
	FindByReferenceId(ctx context.Context, referenceNo string) (*Transaction, error)
	FindById(ctx context.Context, id string) (*Transaction, error)
	Insert(ctx context.Context, data *Transaction) error
	Update(ctx context.Context, data *Transaction) error
}

type TransactionIService interface {
	GetTransactionsByCustomerXid(ctx context.Context, xid string) ([]*Transaction, error)
	CreateDeposit(ctx context.Context, params *CreateDepositParams) (*Transaction, error)
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

func (s *TransactionService) CreateDeposit(ctx context.Context, params *CreateDepositParams) (*Transaction, error) {
	if params.CustomerXid == "" {
		return nil, fmt.Errorf("customer xid is required")
	}
	if params.ReferenceId == "" {
		return nil, fmt.Errorf("reference no is required")
	}

	targetWallet, err := s.walletService.GetWalletByXid(ctx, params.CustomerXid)
	if err != nil {
		return nil, err
	}
	if err = s.walletService.ValidateWallet(targetWallet); err != nil {
		return nil, err
	}

	trx, err := s.repository.FindByReferenceId(ctx, params.ReferenceId)
	if err != nil {
		return nil, err
	}
	if trx != nil {
		return nil, ErrReferenceNoAlreadyExists
	}

	uuidRandom, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	var randomId string
	for {
		randomId = uuidRandom.String()
		existingTrx, err := s.repository.FindById(ctx, randomId)
		if err != nil {
			return nil, err
		}
		if existingTrx == nil {
			break
		}
	}

	err = s.repository.Insert(ctx, &Transaction{
		Id:           randomId,
		Status:       STATUS_PENDING,
		TransactedAt: time.Now(),
		Type:         TYPE_DEPOSIT,
		Amount:       params.Amount,
		ReferenceId:  params.ReferenceId,
		WalletId:     targetWallet.Id,
	})
	if err != nil {
		return nil, err
	}

	trx, err = s.repository.FindByReferenceId(ctx, params.ReferenceId)
	if err != nil {
		return nil, err
	}

	go func() {
		log.Printf("queued %s\n", trx.Id)
		time.Sleep(5 * time.Second)
		log.Printf("disbursing balance to wallet %s, amount: %f\n", targetWallet.Id, params.Amount)
		err := s.walletService.AddBalance(ctx, targetWallet.Id, params.Amount)
		if err != nil {
			log.Printf("error adding balance to wallet %s, amount %f: %v\n", targetWallet.Id, params.Amount, err)
			return
		}

		trx.Status = STATUS_SUCCESS
		log.Printf("updating transaction %s\n", trx.Id)
		err = s.repository.Update(ctx, trx)
		if err != nil {
			log.Printf("error updating transaction %s: %v\n", trx.Id, err)
			return
		}
	}()

	return trx, nil
}
