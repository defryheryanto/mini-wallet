package wallet_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/defryheryanto/mini-wallet/internal/wallet"
	"github.com/defryheryanto/mini-wallet/internal/wallet/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWalletService_Create(t *testing.T) {
	mockedErr := fmt.Errorf("mocked")

	t.Run("should return error when params invalid", func(t *testing.T) {
		service := wallet.NewWalletService(mocks.NewWalletRepository(t))

		err := service.Create(context.TODO(), &wallet.CreateWalletParams{})
		assert.Equal(t, wallet.ErrOwnedByRequired, err)
	})

	t.Run("should return error when failed to get wallet by id", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, mockedErr)
		service := wallet.NewWalletService(repository)

		err := service.Create(context.TODO(), &wallet.CreateWalletParams{
			OwnedBy: "test",
		})
		assert.Equal(t, mockedErr, err)
	})

	t.Run("should return error when failed to insert", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Return(mockedErr)
		service := wallet.NewWalletService(repository)

		err := service.Create(context.TODO(), &wallet.CreateWalletParams{
			OwnedBy: "test",
		})
		assert.Equal(t, mockedErr, err)
	})

	t.Run("should not return error if operation success", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			insertParams, ok := args.Get(1).(*wallet.Wallet)
			assert.True(t, ok, "second argument of insert should be *Wallet")
			assert.Equal(t, "test", insertParams.OwnedBy)
		}).Return(nil)

		service := wallet.NewWalletService(repository)
		err := service.Create(context.TODO(), &wallet.CreateWalletParams{
			OwnedBy: "test",
		})
		assert.Nil(t, err)
	})
}

func TestWalletService_EnableWallet(t *testing.T) {
	customerXid := "test"
	mockedErr := fmt.Errorf("mocked")
	t.Run("should return error if failed to find wallet by customer xid", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindByCustomerXid", mock.Anything, customerXid).Return(nil, mockedErr)

		service := wallet.NewWalletService(repository)
		result, err := service.EnableWallet(context.TODO(), customerXid)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, result)
	})
	t.Run("should return error if wallet not found", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindByCustomerXid", mock.Anything, customerXid).Return(nil, nil)

		service := wallet.NewWalletService(repository)
		result, err := service.EnableWallet(context.TODO(), customerXid)
		assert.Equal(t, wallet.ErrWalletNotFound, err)
		assert.Nil(t, result)
	})
	t.Run("should return error if wallet already enabled", func(t *testing.T) {
		now := time.Now()
		repository := mocks.NewWalletRepository(t)
		repository.On("FindByCustomerXid", mock.Anything, customerXid).Return(&wallet.Wallet{
			EnabledAt: &now,
		}, nil)

		service := wallet.NewWalletService(repository)
		result, err := service.EnableWallet(context.TODO(), customerXid)
		assert.Equal(t, wallet.ErrWalletAlreadyEnabled, err)
		assert.Nil(t, result)
	})
	t.Run("should return error if update wallet failed", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindByCustomerXid", mock.Anything, customerXid).Return(&wallet.Wallet{}, nil)
		repository.On("Update", mock.Anything, mock.Anything).Return(mockedErr)

		service := wallet.NewWalletService(repository)
		result, err := service.EnableWallet(context.TODO(), customerXid)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, result)
	})
	t.Run("should return wallet data if operation success", func(t *testing.T) {
		now := time.Now()
		repository := mocks.NewWalletRepository(t)
		repository.On("FindByCustomerXid", mock.Anything, customerXid).Return(&wallet.Wallet{
			DisabledAt: &now,
			EnabledAt:  nil,
		}, nil)
		repository.On("Update", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			updateParams, ok := args.Get(1).(*wallet.Wallet)
			assert.True(t, ok, "second argument of update should be *Wallet")
			assert.Nil(t, updateParams.DisabledAt)
			assert.NotNil(t, updateParams.EnabledAt)
			assert.Equal(t, wallet.STATUS_ENABLED, updateParams.Status)
		}).Return(nil)

		service := wallet.NewWalletService(repository)
		result, err := service.EnableWallet(context.TODO(), customerXid)
		assert.NotNil(t, result)
		assert.Nil(t, err)
	})
}

func TestWalletService_GetWalletByXid(t *testing.T) {
	customerXid := "test"
	mockedErr := fmt.Errorf("mocked")
	t.Run("should return error if failed to find wallet by customer xid", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindByCustomerXid", mock.Anything, customerXid).Return(nil, mockedErr)

		service := wallet.NewWalletService(repository)
		result, err := service.GetWalletByXid(context.TODO(), customerXid)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, result)
	})
	t.Run("should return error if wallet not found", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindByCustomerXid", mock.Anything, customerXid).Return(nil, nil)

		service := wallet.NewWalletService(repository)
		result, err := service.GetWalletByXid(context.TODO(), customerXid)
		assert.Equal(t, wallet.ErrWalletNotFound, err)
		assert.Nil(t, result)
	})
	t.Run("should return error if wallet is disabled", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindByCustomerXid", mock.Anything, customerXid).Return(&wallet.Wallet{
			Status: wallet.STATUS_DISABLED,
		}, nil)

		service := wallet.NewWalletService(repository)
		result, err := service.GetWalletByXid(context.TODO(), customerXid)
		assert.Equal(t, wallet.ErrWalletDisabled, err)
		assert.Nil(t, result)
	})
	t.Run("should return wallet if operation success", func(t *testing.T) {
		now := time.Now()
		targetWallet := &wallet.Wallet{
			Id:        "test",
			OwnedBy:   "test-owned",
			Status:    wallet.STATUS_ENABLED,
			EnabledAt: &now,
			Balance:   0,
		}
		repository := mocks.NewWalletRepository(t)
		repository.On("FindByCustomerXid", mock.Anything, customerXid).Return(targetWallet, nil)

		service := wallet.NewWalletService(repository)
		result, err := service.GetWalletByXid(context.TODO(), customerXid)
		assert.Nil(t, err)
		assert.Equal(t, targetWallet.Id, result.Id)
		assert.Equal(t, targetWallet.OwnedBy, result.OwnedBy)
		assert.Equal(t, targetWallet.Status, result.Status)
		assert.Equal(t, targetWallet.EnabledAt, result.EnabledAt)
		assert.Equal(t, targetWallet.Balance, result.Balance)
	})
}

func TestWalletService_AddBalance(t *testing.T) {
	mockedErr := fmt.Errorf("mocked")
	walletId := "test-wallet-id"
	amount := float64(10_000)

	t.Run("should return error if failed to get wallet", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindById", mock.Anything, walletId).Return(nil, mockedErr)

		service := wallet.NewWalletService(repository)
		err := service.AddBalance(context.TODO(), walletId, amount)
		assert.Equal(t, mockedErr, err)
	})

	t.Run("should return error if wallet not active", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindById", mock.Anything, walletId).Return(&wallet.Wallet{
			Status: wallet.STATUS_DISABLED,
		}, nil)

		service := wallet.NewWalletService(repository)
		err := service.AddBalance(context.TODO(), walletId, amount)
		assert.Equal(t, wallet.ErrWalletDisabled, err)
	})

	t.Run("should return error if failed to update balance", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindById", mock.Anything, walletId).Return(&wallet.Wallet{
			Status: wallet.STATUS_ENABLED,
		}, nil)
		repository.On("Update", mock.Anything, mock.Anything).Return(mockedErr)

		service := wallet.NewWalletService(repository)
		err := service.AddBalance(context.TODO(), walletId, amount)
		assert.Equal(t, mockedErr, err)
	})

	t.Run("should update balance if operations success", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindById", mock.Anything, walletId).Return(&wallet.Wallet{
			Status:  wallet.STATUS_ENABLED,
			Balance: 100_000,
		}, nil)
		repository.On("Update", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			updateParams, ok := args.Get(1).(*wallet.Wallet)
			assert.True(t, ok, "params should be *Wallet")
			assert.Equal(t, float64(110_000), updateParams.Balance)
		}).Return(nil)

		service := wallet.NewWalletService(repository)
		err := service.AddBalance(context.TODO(), walletId, amount)
		assert.Nil(t, err)
	})
}
