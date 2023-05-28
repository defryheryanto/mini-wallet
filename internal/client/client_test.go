package client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/defryheryanto/mini-wallet/internal/client"
	client_mock "github.com/defryheryanto/mini-wallet/internal/client/mocks"
	"github.com/defryheryanto/mini-wallet/internal/storage/manager"
	wallet_mock "github.com/defryheryanto/mini-wallet/internal/wallet/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientService_Create(t *testing.T) {
	mockedErr := fmt.Errorf("mocked")
	xid := "random-string-xid"
	t.Run("should return error if failed to get by xid", func(t *testing.T) {
		walletService := wallet_mock.NewWalletIService(t)

		repository := client_mock.NewClientRepository(t)
		repository.On("FindByXid", mock.Anything, xid).Return(nil, mockedErr)

		storageManager := &manager.MockStorageManager{}
		service := client.NewClientService(repository, walletService, storageManager)

		res, err := service.Create(context.TODO(), xid)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, res)
	})

	t.Run("should return error if xid already taken", func(t *testing.T) {
		walletService := wallet_mock.NewWalletIService(t)

		repository := client_mock.NewClientRepository(t)
		repository.On("FindByXid", mock.Anything, xid).Return(&client.Client{}, nil)

		storageManager := &manager.MockStorageManager{}
		service := client.NewClientService(repository, walletService, storageManager)

		res, err := service.Create(context.TODO(), xid)
		assert.Equal(t, client.ErrXidAlreadyTaken, err)
		assert.Nil(t, res)
	})

	t.Run("should return error failed to find by token", func(t *testing.T) {
		walletService := wallet_mock.NewWalletIService(t)

		repository := client_mock.NewClientRepository(t)
		repository.On("FindByXid", mock.Anything, xid).Return(nil, nil)
		repository.On("FindByToken", mock.Anything, mock.Anything).Return(nil, mockedErr)

		storageManager := &manager.MockStorageManager{}
		service := client.NewClientService(repository, walletService, storageManager)

		res, err := service.Create(context.TODO(), xid)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, res)
	})

	t.Run("should return error if failed to insert client data", func(t *testing.T) {
		walletService := wallet_mock.NewWalletIService(t)

		repository := client_mock.NewClientRepository(t)
		repository.On("FindByXid", mock.Anything, xid).Return(nil, nil)
		repository.On("FindByToken", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Return(mockedErr)

		storageManager := &manager.MockStorageManager{}
		service := client.NewClientService(repository, walletService, storageManager)

		res, err := service.Create(context.TODO(), xid)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, res)
	})

	t.Run("should return error if failed to create wallet", func(t *testing.T) {
		repository := client_mock.NewClientRepository(t)
		repository.On("FindByXid", mock.Anything, xid).Return(nil, nil)
		repository.On("FindByToken", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Return(nil)

		walletService := wallet_mock.NewWalletIService(t)
		walletService.On("Create", mock.Anything, mock.Anything).Return(mockedErr)

		storageManager := &manager.MockStorageManager{}
		service := client.NewClientService(repository, walletService, storageManager)

		res, err := service.Create(context.TODO(), xid)
		assert.Equal(t, mockedErr, err)
		assert.Nil(t, res)
	})
}
