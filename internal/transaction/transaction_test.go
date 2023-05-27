package transaction_test

import (
	"context"
	"fmt"
	"testing"

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

		service := transaction.NewTransactionService(repository, walletService)

		trx, err := service.GetTransactionsByCustomerXid(context.TODO(), customerXid)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, trx)
	})
	t.Run("should return error if wallet not found", func(t *testing.T) {
		repository := transaction_mock.NewTransactionRepository(t)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("GetWalletByXid", mock.Anything, customerXid).Return(nil, nil)

		service := transaction.NewTransactionService(repository, walletService)

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

		service := transaction.NewTransactionService(repository, walletService)

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

		service := transaction.NewTransactionService(repository, walletService)

		trx, err := service.GetTransactionsByCustomerXid(context.TODO(), customerXid)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(trx))
	})
}
