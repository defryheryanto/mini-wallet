package transaction_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/defryheryanto/mini-wallet/internal/storage/manager"
	"github.com/defryheryanto/mini-wallet/internal/transaction"
	transaction_mock "github.com/defryheryanto/mini-wallet/internal/transaction/mocks"
	"github.com/defryheryanto/mini-wallet/internal/wallet"
	wallet_mock "github.com/defryheryanto/mini-wallet/internal/wallet/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_GetTransactionsByCustomerXid(t *testing.T) {
	mockedErr := fmt.Errorf("mocked")
	customerXid := "test"
	t.Run("should return error if failed to get wallet", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, customerXid).Return(nil, mockedErr)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.GetTransactionsByCustomerXid(context.TODO(), customerXid)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})
	t.Run("should return error if wallet not found", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, customerXid).Return(nil, nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.GetTransactionsByCustomerXid(context.TODO(), customerXid)
		assert.Equal(t, wallet.ErrWalletNotFound, err)
		assert.Nil(t, trx)
	})
	t.Run("should return error if wallet is disabled", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, customerXid).Return(&wallet.Wallet{
			Status: wallet.STATUS_DISABLED,
		}, nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.GetTransactionsByCustomerXid(context.TODO(), customerXid)
		assert.Equal(t, wallet.ErrWalletDisabled, err)
		assert.Nil(t, trx)
	})
	t.Run("should return transactions if operations success", func(t *testing.T) {
		targetWallet := &wallet.Wallet{
			Id:     "test-wallet",
			Status: wallet.STATUS_ENABLED,
		}
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindTransactionsByWalletId", mock.Anything, targetWallet.Id).Return([]*transaction.Transaction{
			{
				Id: "first-transaction",
			},
			{
				Id: "second-transaction",
			},
		}, nil)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, customerXid).Return(targetWallet, nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.GetTransactionsByCustomerXid(context.TODO(), customerXid)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(trx))
	})
}

func TestTransactionService_CreateDeposit(t *testing.T) {
	mockedErr := fmt.Errorf("mocked")
	params := &transaction.CreateDepositParams{
		CustomerXid: "test-xid",
		ReferenceId: "test-ref-no",
		Amount:      10_000,
	}

	t.Run("should return error if customer xid is empty", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		walletService := wallet_mock.NewWalletIService(t)
		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateDeposit(context.TODO(), &transaction.CreateDepositParams{
			ReferenceId: "ref-no",
			Amount:      10_000,
		})
		assert.Error(t, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if ref no is empty", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		walletService := wallet_mock.NewWalletIService(t)
		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateDeposit(context.TODO(), &transaction.CreateDepositParams{
			CustomerXid: "test-xid",
			Amount:      10_000,
		})
		assert.Error(t, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if failed to get wallet", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(nil, mockedErr)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateDeposit(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if wallet not active", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(mockedErr)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateDeposit(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if failed to get transaction by ref no", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, mockedErr)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateDeposit(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if ref no already used", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(&transaction.Transaction{}, nil)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateDeposit(context.TODO(), params)
		assert.Equal(t, transaction.ErrReferenceNoAlreadyExists, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if get transaction by id failed", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, nil)
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, mockedErr)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateDeposit(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})
	t.Run("should return error if insert transaction failed", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, nil)
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Return(mockedErr)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateDeposit(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if failed to get created transaction", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, nil).Once()
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Return(nil)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, mockedErr).Once()

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateDeposit(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})

	t.Run("should return created transaction if operation success", func(t *testing.T) {
		createdTransaction := &transaction.Transaction{
			Id:           "test-id",
			Status:       transaction.STATUS_PENDING,
			TransactedAt: time.Now(),
			Type:         transaction.TYPE_DEPOSIT,
			Amount:       10_000,
			ReferenceId:  "ref-no-test",
			WalletId:     "test-wallet-id",
		}
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, nil).Once()
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			insertParams, ok := args.Get(1).(*transaction.Transaction)
			assert.True(t, ok, "params should be *Transaction")
			assert.Equal(t, params.Amount, insertParams.Amount)
			assert.Equal(t, params.ReferenceId, insertParams.ReferenceId)
			assert.Equal(t, transaction.STATUS_PENDING, insertParams.Status)
			assert.Equal(t, transaction.TYPE_DEPOSIT, insertParams.Type)
		}).Return(nil)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(createdTransaction, nil).Once()

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateDeposit(context.TODO(), params)
		assert.Nil(t, err)
		assert.Equal(t, createdTransaction.Id, trx.Id)
		assert.Equal(t, createdTransaction.Status, trx.Status)
		assert.Equal(t, createdTransaction.TransactedAt, trx.TransactedAt)
		assert.Equal(t, createdTransaction.Type, trx.Type)
		assert.Equal(t, createdTransaction.Amount, trx.Amount)
		assert.Equal(t, createdTransaction.ReferenceId, trx.ReferenceId)
		assert.Equal(t, createdTransaction.WalletId, trx.WalletId)
	})
}

func TestTransactionService_CreateWithdrawal(t *testing.T) {
	mockedErr := fmt.Errorf("mocked")
	params := &transaction.CreateWithdrawalParams{
		CustomerXid: "test-xid",
		ReferenceId: "test-ref-no",
		Amount:      10_000,
	}

	t.Run("should return error if customer xid is empty", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		walletService := wallet_mock.NewWalletIService(t)
		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), &transaction.CreateWithdrawalParams{
			ReferenceId: "ref-no",
			Amount:      10_000,
		})
		assert.Error(t, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if ref no is empty", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		walletService := wallet_mock.NewWalletIService(t)
		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), &transaction.CreateWithdrawalParams{
			CustomerXid: "test-xid",
			Amount:      10_000,
		})
		assert.Error(t, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if failed to get wallet", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(nil, mockedErr)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if wallet not active", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(mockedErr)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if wallet balance insufficient", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{
			Balance: 0,
		}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), params)
		assert.Equal(t, wallet.ErrInsufficientBalance, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if failed to get transaction by ref no", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, mockedErr)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{
			Balance: 15_001,
		}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if ref no already used", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(&transaction.Transaction{}, nil)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{
			Balance: 15_001,
		}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), params)
		assert.Equal(t, transaction.ErrReferenceNoAlreadyExists, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if get transaction by id failed", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, nil)
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, mockedErr)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{
			Balance: 15_001,
		}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})
	t.Run("should return error if insert transaction failed", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, nil)
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Return(mockedErr)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{
			Balance: 15_001,
		}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})

	t.Run("should return error if failed to get created transaction", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, nil).Once()
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Return(nil)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, mockedErr).Once()

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{
			Balance: 15_001,
		}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), params)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})

	t.Run("should return created transaction if operation success", func(t *testing.T) {
		createdTransaction := &transaction.Transaction{
			Id:           "test-id",
			Status:       transaction.STATUS_PENDING,
			TransactedAt: time.Now(),
			Type:         transaction.TYPE_DEPOSIT,
			Amount:       10_000,
			ReferenceId:  "ref-no-test",
			WalletId:     "test-wallet-id",
		}
		repository := transaction_mock.NewTransactionRepository(t)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(nil, nil).Once()
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			insertParams, ok := args.Get(1).(*transaction.Transaction)
			assert.True(t, ok, "params should be *Transaction")
			assert.Equal(t, params.Amount, insertParams.Amount)
			assert.Equal(t, params.ReferenceId, insertParams.ReferenceId)
			assert.Equal(t, transaction.STATUS_PENDING, insertParams.Status)
			assert.Equal(t, transaction.TYPE_WITHDRAWAL, insertParams.Type)
		}).Return(nil)
		repository.On("FindByReferenceId", mock.Anything, params.ReferenceId).Return(createdTransaction, nil).Once()

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, params.CustomerXid).Return(&wallet.Wallet{
			Balance: 15_001,
		}, nil)
		walletService.On("ValidateWallet", mock.Anything).Return(nil)

		service := transaction.NewTransactionService(repository, walletService, &manager.MockStorageManager{})

		trx, err := service.CreateWithdrawal(context.TODO(), params)
		assert.Nil(t, err)
		assert.Equal(t, createdTransaction.Id, trx.Id)
		assert.Equal(t, createdTransaction.Status, trx.Status)
		assert.Equal(t, createdTransaction.TransactedAt, trx.TransactedAt)
		assert.Equal(t, createdTransaction.Type, trx.Type)
		assert.Equal(t, createdTransaction.Amount, trx.Amount)
		assert.Equal(t, createdTransaction.ReferenceId, trx.ReferenceId)
		assert.Equal(t, createdTransaction.WalletId, trx.WalletId)
	})
}
